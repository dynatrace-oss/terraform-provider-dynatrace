package rest

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/version"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/dynatrace/dynatrace-configuration-as-code-core/api/auth"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
	crest "github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/clients"
)

var NoClassicURLDefinedErr = errors.New("no Environment URL has been specified. Use either the environment variable `DYNATRACE_ENV_URL` or the configuration attribute `dt_env_url` of the provider for that")

func APITokenClient(credentials *Credentials) Client {
	return &api_token_client{credentials: credentials}
}

type api_token_client struct {
	credentials *Credentials
}

func (me *api_token_client) Credentials() *Credentials {
	return me.credentials
}

func (me *api_token_client) Get(ctx context.Context, url string, expectedStatusCodes ...int) Request {
	req := &api_token_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: http.MethodGet}
	if len(expectedStatusCodes) > 0 {
		req.expect = statuscodes(expectedStatusCodes)
	}
	return req
}

func (me *api_token_client) Post(ctx context.Context, url string, payload any, expectedStatusCodes ...int) Request {
	req := &api_token_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: http.MethodPost, payload: payload, headers: map[string]string{"Content-Type": "application/json"}}
	if len(expectedStatusCodes) > 0 {
		req.expect = statuscodes(expectedStatusCodes)
	}
	return req
}

func (me *api_token_client) Put(ctx context.Context, url string, payload any, expectedStatusCodes ...int) Request {
	req := &api_token_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: http.MethodPut, payload: payload, headers: map[string]string{"Content-Type": "application/json"}}
	if len(expectedStatusCodes) > 0 {
		req.expect = statuscodes(expectedStatusCodes)
	}
	return req
}

func (me *api_token_client) Delete(ctx context.Context, url string, expectedStatusCodes ...int) Request {
	req := &api_token_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: http.MethodDelete}
	if len(expectedStatusCodes) > 0 {
		req.expect = statuscodes(expectedStatusCodes)
	}
	return req
}

func (me *api_token_request) Finish(vs ...any) error {
	credentials := me.client.Credentials()
	if !credentials.ContainsAPIToken() {
		return NoAPITokenError
	}
	if DYNATRACE_HTTP_LEGACY {
		legacyRequest := legacy_request(*me)
		if credentials.URL == TestCaseEnvURL {
			return errors.New("legacy")
		}
		return legacyRequest.Finish(vs...)
	}
	classicRequest := classic_request(*me)
	if credentials.URL == TestCaseEnvURL {
		return errors.New("classic")
	}
	return classicRequest.Finish(vs...)
}

type api_token_request request

func (me *api_token_request) Expect(codes ...int) Request {
	me.expect = statuscodes(codes)
	return me
}

func (me *api_token_request) OnResponse(onResponse func(resp *http.Response)) Request {
	me.onResponse = onResponse
	return me
}

func (me *api_token_request) SetHeader(name string, value string) {
	if me.headers == nil {
		me.headers = map[string]string{}
	}
	me.headers[name] = value
}

type classic_request request

var classicClientCache = map[string]*rest.Client{}

var classicClientCacheMutex sync.Mutex

func CreateClassicOAuthBasedClient(ctx context.Context, credentials *Credentials) *rest.Client {
	var parsedURL *url.URL
	parsedURL, err := url.Parse(credentials.URL)
	if err != nil {
		return nil
	}

	if credentials.OAuth.ClientID == "" || credentials.OAuth.ClientSecret == "" {
		return nil
	}

	logging.InstallRoundTripper()

	oauthClient := rest.NewClient(
		parsedURL,
		auth.NewOAuthClient(
			ctx,
			&clientcredentials.Config{
				ClientID:     credentials.OAuth.ClientID,
				ClientSecret: credentials.OAuth.ClientSecret,
				TokenURL:     credentials.OAuth.TokenURL,
				AuthStyle:    oauth2.AuthStyleInParams}),
		crest.WithHTTPListener(logging.HTTPListener("classic ")),
	)

	oauthClient.SetHeader("User-Agent", version.UserAgent())
	return oauthClient
}

func CreateClassicClient(classicURL string, apiToken string) (*rest.Client, error) {
	classicClientCacheMutex.Lock()
	defer classicClientCacheMutex.Unlock()

	if client, found := classicClientCache[classicURL]; found {
		return client, nil
	}

	factory := clients.Factory()
	factory = factory.WithUserAgent(version.UserAgent())
	factory = factory.WithClassicURL(classicURL)
	factory = factory.WithAccessToken(apiToken)
	factory = factory.WithHTTPListener(logging.HTTPListener("classic "))

	client, err := factory.CreateClassicClient()

	classicClientCache[classicURL] = client

	return client, err
}

func (me *classic_request) Finish(optionalTarget ...any) error {
	var target any
	if len(optionalTarget) > 0 {
		target = optionalTarget[0]
	}

	PreRequest()

	classicURL := request(*me).evalClassicURL()
	if len(classicURL) == 0 {
		// sanity check - this should not be empty anymore at this point
		return NoClassicURLDefinedErr
	}

	client, err := CreateClassicClient(classicURL, me.client.Credentials().Token)
	if err != nil {
		return err
	}
	for headername, headervalue := range me.headers {
		client.SetHeader(headername, headervalue)
	}

	pathURL, err := url.Parse(me.url)
	if err != nil {
		return err
	}

	return request(*me).HandleResponse(client, pathURL, target)
}
