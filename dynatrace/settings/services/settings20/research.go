package settings20

import (
	"context"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/automation/httplog"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/clients"
	"golang.org/x/oauth2/clientcredentials"
)

var regexpSaasTenant = regexp.MustCompile(`https:\/\/(.*).(live|apps).dynatrace.com`)
var regexpSprintTenant = regexp.MustCompile(`https:\/\/(.*).sprint(?:\.apps)?.dynatracelabs.com`)
var regexpDevTenant = regexp.MustCompile(`https:\/\/(.*).dev(?:\.apps)?.dynatracelabs.com`)

const (
	ProdTokenURL   = "https://sso.dynatrace.com/sso/oauth2/token"
	SprintTokenURL = "https://sso-sprint.dynatracelabs.com/sso/oauth2/token"
	DevTokenURL    = "https://sso-dev.dynatracelabs.com/sso/oauth2/token"
)

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

func evalEnvURL() string {
	return strings.TrimSuffix(os.Getenv("DYNATRACE_ENV_URL"), "/")
}

func evalClassicURL() string {
	envURL := evalEnvURL()
	if len(envURL) == 0 {
		return envURL
	}
	envURL = strings.ReplaceAll(envURL, ".dev.apps.dynatracelabs.", ".dev.dynatracelabs.")
	envURL = strings.ReplaceAll(envURL, ".sprint.apps.dynatracelabs.", ".sprint.dynatracelabs.")
	envURL = strings.ReplaceAll(envURL, ".apps.dynatrace.", ".live.dynatrace.")
	return envURL
}

func evalPlatformURL() string {
	envURL := evalEnvURL()
	if len(envURL) == 0 {
		return envURL
	}
	envURL = strings.ReplaceAll(envURL, ".dev.dynatracelabs.", ".dev.apps.dynatracelabs.")
	envURL = strings.ReplaceAll(envURL, ".live.dynatrace.", ".apps.dynatrace.")
	envURL = strings.ReplaceAll(envURL, ".sprint.dynatracelabs.", ".sprint.apps.dynatracelabs.")
	return envURL

}

func SendClassicClientRequest() error {
	ctx := context.Background()
	httplog.InstallRoundTripper()
	clientsFactory := clients.Factory()
	clientsFactory = clientsFactory.WithUserAgent("Dynatrace Terraform Provider")
	clientsFactory = clientsFactory.WithClassicURL(evalClassicURL())
	clientsFactory = clientsFactory.WithAccessToken(os.Getenv("DYNATRACE_API_TOKEN"))
	clientsFactory = clientsFactory.WithHTTPListener(httplog.HTTPListener)

	restClient, err := clientsFactory.CreateClassicClient()
	if err != nil {
		return err
	}

	response, err := restClient.GET(
		ctx,
		"/api/v2/settings/objects",
		rest.RequestOptions{
			QueryParams: url.Values{
				"schemaIds":   []string{"builtin:management-zones"},
				"fields":      []string{"objectId"},
				"adminAccess": []string{"false"},
			},
		},
	)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
}

func SendPlatformClientRequest() error {
	ctx := context.Background()
	httplog.InstallRoundTripper()
	clientsFactory := clients.Factory()
	clientsFactory = clientsFactory.WithUserAgent("Dynatrace Terraform Provider")
	clientsFactory = clientsFactory.WithPlatformURL(evalPlatformURL())
	clientsFactory = clientsFactory.WithHTTPListener(httplog.HTTPListener)
	clientsFactory = clientsFactory.WithOAuthCredentials(clientcredentials.Config{
		ClientID:     os.Getenv("DT_CLIENT_ID"),
		ClientSecret: os.Getenv("DT_CLIENT_SECRET"),
		TokenURL:     evalTokenURL(evalPlatformURL()),
	})

	platformClient, err := clientsFactory.CreatePlatformClient(ctx)
	if err != nil {
		return err
	}
	response, err := platformClient.GET(
		ctx,
		"/platform/classic/environment-api/v2/settings/objects",
		rest.RequestOptions{
			QueryParams: url.Values{
				"schemaIds":   []string{"builtin:management-zones"},
				"fields":      []string{"objectId"},
				"adminAccess": []string{"false"},
			},
		},
	)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
}
