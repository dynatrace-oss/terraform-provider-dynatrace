package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/version"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/dynatrace/dynatrace-configuration-as-code-core/api/auth"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
)

var eligiblePlatformRequests = map[string]string{
	"/api/v2/settings/objects":   "/platform/classic/environment-api/v2/settings/objects",
	"/api/v2/settings/schemas":   "/platform/classic/environment-api/v2/settings/schemas",
	"/api/v2/networkZones":       "/platform/classic/environment-api/v2/networkZones",
	"/api/v2/entities":           "/platform/classic/environment-api/v2/entities",
	"/api/v2/entityTypes":        "/platform/classic/environment-api/v2/entityTypes",
	"/api/v2/tags":               "/platform/classic/environment-api/v2/tags",
	"/api/v2/slo":                "/platform/classic/environment-api/v2/slo",
	"/api/v2/extensions":         "/platform/classic/environment-api/v2/extensions",
	"/api/v2/apiTokens":          "/platform/classic/environment-api/v2/apiTokens",
	"/api/v2/credentials":        "/platform/classic/environment-api/v2/credentials",
	"/api/v2/synthetic/monitors": "/platform/classic/environment-api/v2/synthetic/monitors",
	"/api/v2/activeGateTokens":   "/platform/classic/environment-api/v2/activeGateTokens",

	"/api/v1/synthetic/monitors": "/platform/classic/environment-api/v1/synthetic/monitors",
	"/api/v1/deployment":         "/platform/classic/environment-api/v1/deployment",
	"/api/v1/synthetic/nodes":    "/platform/classic/environment-api/v1/synthetic/nodes",
}

type platform_request request

var platformClientCache = map[string]*rest.Client{}

var platformClientCacheMutex sync.Mutex
var NoPlatformCredentialsErr = errors.New("neither oauth credentials nor platform token present")

func configureCommonRestClient(restClient *rest.Client) (*rest.Client, error) {
	restClient.SetHeader("User-Agent", version.UserAgent())
	return restClient, nil
}

func setHttpClientContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, oauth2.HTTPClient, &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
	})
}

func CreatePlatformOAuthClient(ctx context.Context, u string, credentials *Credentials) (*rest.Client, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL %q: %w", u, err)
	}
	oauthConfig := clientcredentials.Config{
		ClientID:     credentials.OAuth.ClientID,
		ClientSecret: credentials.OAuth.ClientSecret,
		TokenURL:     evalTokenURL(parsedURL.String()),
	}
	ctx = setHttpClientContext(ctx)
	httpClient := auth.NewOAuthClient(ctx, &oauthConfig)

	opts := []rest.Option{
		rest.WithHTTPListener(logging.HTTPListener("platform")),
		rest.WithRateLimiter(),
		rest.WithRetryOptions(defaultRetryOptions),
	}

	return configureCommonRestClient(rest.NewClient(parsedURL, httpClient, opts...))
}

func CreatePlatformTokenClient(u string, credentials *Credentials) (*rest.Client, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL %q: %w", u, err)
	}

	opts := []rest.Option{
		rest.WithHTTPListener(logging.HTTPListener("plat/tok")),
		rest.WithRateLimiter(),
		rest.WithRetryOptions(defaultRetryOptions),
	}

	return configureCommonRestClient(rest.NewClient(parsedURL, NewBearerTokenBasedClient(credentials.OAuth.PlatformToken), opts...))
}

func CreatePlatformClient(ctx context.Context, platformURL string, credentials *Credentials) (*rest.Client, error) {
	platformClientCacheMutex.Lock()
	defer platformClientCacheMutex.Unlock()

	if client, found := platformClientCache[platformURL]; found {
		return client, nil
	}
	PreRequest()

	var client *rest.Client
	var err error
	if credentials.ContainsPlatformToken() {
		client, err = CreatePlatformTokenClient(platformURL, credentials)
	} else if credentials.ContainsOAuth() {
		client, err = CreatePlatformOAuthClient(ctx, platformURL, credentials)
	} else {
		return nil, NoPlatformCredentialsErr
	}
	if err != nil {
		return nil, err
	}
	platformClientCache[platformURL] = client

	return client, err
}

func (me *platform_request) Finish(optionalTarget ...any) error {
	var target any
	if len(optionalTarget) > 0 {
		target = optionalTarget[0]
	}

	PreRequest()

	platformURL := me.evalPlatformURL()

	client, err := CreatePlatformClient(me.ctx, platformURL, me.client.Credentials())
	if err != nil {
		return err
	}

	meURL := me.url
	for eligiblePlatformURL, replacement := range eligiblePlatformRequests {
		if strings.HasPrefix(meURL, eligiblePlatformURL) {
			meURL = strings.ReplaceAll(meURL, eligiblePlatformURL, replacement)
			// We can exit early here
			// Once we have found a URL eligible to replace with its platform counterpart
			// we shouldn't expect a second one
			break
		}
	}

	fullURL, err := url.Parse(platformURL + meURL)
	if err != nil {
		return err
	}
	return request(*me).HandleResponse(client, fullURL, target)
}

func (me *platform_request) evalPlatformURL() string {
	envURL := strings.TrimSuffix(strings.TrimSpace(me.client.Credentials().URL), "/")
	if len(envURL) == 0 {
		return envURL
	}
	envURL = strings.ReplaceAll(envURL, ".dev.dynatracelabs.", ".dev.apps.dynatracelabs.")
	envURL = strings.ReplaceAll(envURL, ".live.dynatrace.", ".apps.dynatrace.")
	envURL = strings.ReplaceAll(envURL, ".sprint.dynatracelabs.", ".sprint.apps.dynatracelabs.")
	return envURL

}

var regexpSaasTenant = regexp.MustCompile(`https:\/\/(.*).(live|apps).dynatrace.com`)
var regexpSprintTenant = regexp.MustCompile(`https:\/\/(.*).sprint(?:\.apps)?.dynatracelabs.com`)
var regexpDevTenant = regexp.MustCompile(`https:\/\/(.*).dev(?:\.apps)?.dynatracelabs.com`)

func evalTokenURL(dtEnvURL string) string {
	if match := regexpSaasTenant.FindStringSubmatch(dtEnvURL); len(match) > 0 {
		return ProdTokenURL
	}
	if match := regexpSprintTenant.FindStringSubmatch(dtEnvURL); len(match) > 0 {
		return SprintTokenURL
	}
	if match := regexpDevTenant.FindStringSubmatch(dtEnvURL); len(match) > 0 {
		return DevTokenURL
	}
	return ""
}

// NewBearerTokenBasedClient creates a new HTTP client with token-based authentication.
// It takes a token string as an argument and returns an instance of *http.Client.
func NewBearerTokenBasedClient(token string) *http.Client {
	// Create a new tokenAuthTransport and initialize it with the provided token.
	return &http.Client{Transport: newBearerTokenAuthTransport(nil, token)}
}

// bearerTokenAuthTransport is a custom transport that adds token-based authentication headers to HTTP requests.
type bearerTokenAuthTransport struct {
	http.RoundTripper
	header http.Header
}

// newBearerTokenAuthTransport creates a new instance of tokenAuthTransport.
// It takes a baseTransport (an existing HTTP transport) and a token string as arguments,
// and returns a pointer to the newly created tokenAuthTransport instance.
func newBearerTokenAuthTransport(baseTransport http.RoundTripper, token string) *bearerTokenAuthTransport {
	// If no baseTransport is provided, use the default HTTP transport.
	if baseTransport == nil {
		baseTransport = http.DefaultTransport
	}

	// Create a new tokenAuthTransport instance and initialize it.
	t := &bearerTokenAuthTransport{
		RoundTripper: baseTransport,
		header:       http.Header{},
	}

	// Set the "Authorization" header with the provided token.
	t.header.Set("Authorization", "Bearer "+token)
	return t
}

// RoundTrip implements the http.RoundTripper interface's RoundTrip method.
// It adds the authentication headers from the tokenAuthTransport instance to the request's headers
// and delegates the actual round trip to the underlying transport.
func (t *bearerTokenAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Copy authentication headers from tokenAuthTransport to the request.
	for k, v := range t.header {
		req.Header[k] = v
	}

	// Perform the actual HTTP request using the underlying transport.
	return t.RoundTripper.RoundTrip(req)
}
