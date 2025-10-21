/**
* @license
* Copyright 2020 Dynatrace LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package config

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// HTTPVerbose if set to `true` terraform-provider-dynatrace.log will contain request and response payload
var HTTPVerbose = (strings.TrimSpace(os.Getenv("DYNATRACE_DEBUG")) == "true")

type IAM struct {
	ClientID     string
	AccountID    string
	ClientSecret string
	TokenURL     string
	EndpointURL  string
}

const (
	CredValDefault = iota
	CredValIAM
	CredValCluster
	CredValAutomation
	CredValNone
	CredValExport
	CredValExportIAM
)

func validateCredentials(conf *ProviderConfiguration, CredentialValidation int) error {
	switch CredentialValidation {
	case CredValDefault:
		if len(conf.EnvironmentURL) == 0 {
			return fmt.Errorf(" No Environment URL has been specified. Use either the environment variable `DYNATRACE_ENV_URL` or the configuration attribute `dt_env_url` of the provider for that")
		}
		if !strings.HasPrefix(conf.EnvironmentURL, "https://") && !strings.HasPrefix(conf.EnvironmentURL, "http://") {
			return fmt.Errorf(" The Environment URL `%s` neither starts with `https://` nor with `http://`. Please check your configuration.\nFor SaaS environments: `https://######.live.dynatrace.com`.\nFor Managed environments: `https://############/e/########-####-####-####-############`", conf.EnvironmentURL)
		}
		// We cannot check for the existence of an API Token
		// OAuth and Platform Tokens are quite possible
	case CredValIAM:
		if len(conf.IAM.AccountID) == 0 {
			return fmt.Errorf(" No OAuth Account ID has been specified. Use either the environment variable `DT_ACCOUNT_ID` or the configuration attribute `iam_account_id` of the provider for that")
		}
		if len(conf.IAM.ClientID) == 0 {
			return fmt.Errorf(" No OAuth Client ID has been specified. Use either the environment variable `DT_CLIENT_ID` or the configuration attribute `iam_client_id` of the provider for that")
		}
		if len(conf.IAM.ClientSecret) == 0 {
			return fmt.Errorf(" No OAuth Client Secret has been specified. Use either the environment variable `DT_CLIENT_SECRET` or the configuration attribute `iam_client_secret` of the provider for that")
		}
		// We don't complain about a missing Token URL anymore
		// It is either getting deducted from the Environment URL or assumed to be the default for a SaaS Production Tenant
		//
		// if len(conf.IAM.TokenURL) == 0 {
		// 	return fmt.Errorf(" No OAuth TokenURL has been specified. Use either the environment variable `DT_TOKEN_URL` or the configuration attribute `iam_token_url` of the provider for that")
		// }
	case CredValCluster:
		if len(conf.ClusterAPIToken) == 0 {
			return fmt.Errorf(" No Cluster API Token has been specified. Use either the environment variable `DT_CLUSTER_API_TOKEN` or the configuration attribute `dt_cluster_api_token` of the provider for that")
		}
		if len(conf.ClusterAPIV2URL) == 0 {
			return fmt.Errorf(" No Cluster URL has been specified. Use either the environment variable `DT_CLUSTER_URL` or the configuration attribute `dt_cluster_url` of the provider for that")
		}
	case CredValAutomation:
		if len(conf.Automation.ClientID) == 0 {
			return fmt.Errorf(" No OAuth Client ID for the Automation API has been specified. Use either the environment variable `DT_AUTOMATION_CLIENT_ID` or the configuration attribute `automation_client_id` of the provider for that")
		}
		if len(conf.Automation.ClientSecret) == 0 {
			return fmt.Errorf(" No OAuth Client Secret for the Automation API has been specified. Use either the environment variable `DT_AUTOMATION_CLIENT_SECRET` or the configuration attribute `automation_client_secret` of the provider for that")
		}
		if len(conf.Automation.TokenURL) == 0 {
			return fmt.Errorf(" No Token URL for the Automation API has been specified. Use either the environment variable `DT_AUTOMATION_TOKEN_URL` or the configuration attribute `automation_token_url` of the provider for that")
		}
		if len(conf.Automation.EnvironmentURL) == 0 {
			return fmt.Errorf(" No Environment URL for the Automation API has been specified. Use either the environment variable `DT_AUTOMATION_ENVIRONMENT_URL` or the configuration attribute `automation_env_url` of the provider for that")
		}
	case CredValExport:
		if len(conf.EnvironmentURL) == 0 {
			return fmt.Errorf(" No Environment URL has been specified. Use either the environment variable `DYNATRACE_ENV_URL` or the configuration attribute `dt_env_url` of the provider for that")
		}
		if len(conf.APIToken) == 0 && len(conf.Automation.PlatformToken) == 0 && validateCredentials(conf, CredValAutomation) != nil {
			return fmt.Errorf(" No API Token, Platform Token, or OAuth has been specified for export. More detailed information can be found in the documentation at https://registry.terraform.io/providers/dynatrace-oss/dynatrace/latest/docs#configure-the-dynatrace-provider")
		}
	case CredValExportIAM:
		if conf.IAM.AccountID == "" {
			return fmt.Errorf(" No Account ID has been specified for export IAM validation")
		}

		return validateCredentials(conf, CredValExport)
	}
	return nil
}

func Credentials(m any, CredentialValidation int) (*rest.Credentials, error) {
	conf := m.(*ProviderConfiguration)
	if err := validateCredentials(conf, CredentialValidation); err != nil {
		return nil, err
	}
	return &rest.Credentials{
		Token: conf.APIToken,
		URL:   conf.EnvironmentURL,
		IAM:   conf.IAM,
		OAuth: conf.Automation,
		Cluster: struct {
			URL   string
			Token string
		}{
			URL:   conf.ClusterAPIV2URL,
			Token: conf.ClusterAPIToken,
		},
	}, nil
}

// ProviderConfiguration contains the initialized API clients to communicate with the Dynatrace API
type ProviderConfiguration struct {
	EnvironmentURL    string
	DTenvURL          string
	DTNonConfigEnvURL string
	DTApiV2URL        string
	ClusterAPIV2URL   string
	ClusterAPIToken   string
	APIToken          string
	IAM               IAM
	Automation        rest.OAuthCredentials
}

type Getter interface {
	Get(key string) any
}

func ProviderConfigure(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
	return ProviderConfigureGeneric(ctx, d)
}

var regexpSaasTenant = regexp.MustCompile(`https:\/\/(.*).(live|apps).dynatrace.com`)
var regexpSprintTenant = regexp.MustCompile(`https:\/\/(.*).sprint(?:\.apps)?.dynatracelabs.com`)
var regexpDevTenant = regexp.MustCompile(`https:\/\/(.*).dev(?:\.apps)?.dynatracelabs.com`)

func ProviderConfigureGeneric(ctx context.Context, d Getter) (any, diag.Diagnostics) {
	dtEnvURL := d.Get("dt_env_url").(string)
	apiToken := d.Get("dt_api_token").(string)
	clusterAPIToken := getString(d, "dt_cluster_api_token")
	clusterURL := getString(d, "dt_cluster_url")

	dtEnvURL = strings.TrimSuffix(strings.TrimSuffix(dtEnvURL, " "), "/")
	if len(dtEnvURL) != 0 {
		if match := regexpSaasTenant.FindStringSubmatch(dtEnvURL); len(match) > 0 {
			dtEnvURL = fmt.Sprintf("https://%s.live.dynatrace.com", match[1])
		}
		if match := regexpSprintTenant.FindStringSubmatch(dtEnvURL); len(match) > 0 {
			dtEnvURL = fmt.Sprintf("https://%s.sprint.dynatracelabs.com", match[1])
		}
		if match := regexpDevTenant.FindStringSubmatch(dtEnvURL); len(match) > 0 {
			dtEnvURL = fmt.Sprintf("https://%s.dev.dynatracelabs.com", match[1])
		}
	}

	clusterURL = strings.TrimSuffix(strings.TrimSuffix(clusterURL, " "), "/")

	fullURL := dtEnvURL + "/api/config/v1"
	fullNonConfigURL := dtEnvURL + "/api/v1"
	fullApiV2URL := dtEnvURL + "/api/v2"

	automation_environment_url := getString(d, "automation_env_url")
	automation_token_url := getString(d, "automation_token_url")
	if len(automation_environment_url) == 0 {
		if match := regexpSaasTenant.FindStringSubmatch(dtEnvURL); len(match) > 0 {
			automation_environment_url = fmt.Sprintf("https://%s.apps.dynatrace.com", match[1])
			automation_token_url = rest.ProdTokenURL
		}
		if match := regexpSprintTenant.FindStringSubmatch(dtEnvURL); len(match) > 0 {
			automation_environment_url = fmt.Sprintf("https://%s.sprint.apps.dynatracelabs.com", match[1])
			automation_token_url = rest.SprintTokenURL
		}
		if match := regexpDevTenant.FindStringSubmatch(dtEnvURL); len(match) > 0 {
			automation_environment_url = fmt.Sprintf("https://%s.dev.apps.dynatracelabs.com", match[1])
			automation_token_url = rest.DevTokenURL
		}
	}

	client_id := getString(d, "client_id")
	client_secret := getString(d, "client_secret")
	account_id := getString(d, "account_id")
	token_url := getString(d, "token_url")
	platform_token := getString(d, "platform_token")

	oauth_endpoint_url := "https://api.dynatrace.com"
	if strings.Contains(dtEnvURL, ".live.dynatrace.com") || strings.Contains(dtEnvURL, ".apps.dynatrace.com") {
		oauth_endpoint_url = rest.ProdIAMEndpointURL
	} else if strings.Contains(dtEnvURL, ".sprint.dynatracelabs.com") || strings.Contains(dtEnvURL, ".sprint.apps.dynatracelabs.com") {
		oauth_endpoint_url = rest.SprintIAMEndpointURL
	} else if strings.Contains(dtEnvURL, ".dev.dynatracelabs.com") || strings.Contains(dtEnvURL, ".dev.apps.dynatracelabs.com") {
		oauth_endpoint_url = rest.DevIAMEndpointURL
	}

	iam_client_id := getString(d, "iam_client_id")
	iam_account_id := getString(d, "iam_account_id")
	iam_client_secret := getString(d, "iam_client_secret")
	iam_token_url := strings.TrimSuffix(strings.TrimSpace(getString(d, "iam_token_url")), "/")
	iam_endpoint_url := strings.TrimSuffix(strings.TrimSpace(getString(d, "iam_endpoint_url")), "/")

	automation_client_id := getString(d, "automation_client_id")
	if len(automation_client_id) == 0 {
		automation_client_id = client_id
	}
	automation_client_secret := getString(d, "automation_client_secret")
	if len(automation_client_secret) == 0 {
		automation_client_secret = client_secret
	}

	automation_client_id = streamlineOAuthCreds(automation_client_id, client_id, iam_client_id)
	automation_client_secret = streamlineOAuthCreds(automation_client_secret, client_secret, iam_client_secret)
	automation_token_url = streamlineOAuthCreds(automation_token_url, token_url, iam_token_url, rest.ProdTokenURL)

	iam_client_id = streamlineOAuthCreds(iam_client_id, client_id, automation_client_id)
	iam_client_secret = streamlineOAuthCreds(iam_client_secret, client_secret, automation_client_secret)
	iam_token_url = streamlineOAuthCreds(iam_token_url, token_url, automation_token_url, rest.ProdTokenURL)
	iam_account_id = streamlineOAuthCreds(iam_account_id, account_id)
	iam_endpoint_url = streamlineOAuthCreds(iam_endpoint_url, oauth_endpoint_url, rest.ProdIAMEndpointURL)

	var diags diag.Diagnostics

	pc := &ProviderConfiguration{
		EnvironmentURL:    dtEnvURL,
		DTenvURL:          fullURL,
		DTApiV2URL:        fullApiV2URL,
		DTNonConfigEnvURL: fullNonConfigURL,
		APIToken:          apiToken,
		ClusterAPIToken:   clusterAPIToken,
		ClusterAPIV2URL:   clusterURL,
		IAM: IAM{
			ClientID:     iam_client_id,
			AccountID:    iam_account_id,
			ClientSecret: iam_client_secret,
			TokenURL:     iam_token_url,
			EndpointURL:  iam_endpoint_url,
		},
		Automation: rest.OAuthCredentials{
			PlatformToken:  platform_token,
			ClientID:       automation_client_id,
			ClientSecret:   automation_client_secret,
			TokenURL:       automation_token_url,
			EnvironmentURL: automation_environment_url,
		},
	}
	return pc, diags
}

func streamlineOAuthCreds(values ...string) string {
	if len(values) == 0 {
		return ""
	}
	for _, value := range values {
		if len(value) != 0 {
			return value
		}
	}
	return ""
}

func getString(d Getter, key string) string {
	if value := d.Get(key); value != nil {
		return value.(string)
	}
	return ""
}

type ConfigGetter struct {
	Provider *schema.Provider
}

func (me ConfigGetter) Get(key string) any {
	schema, found := me.Provider.Schema[key]
	if !found {
		return ""
	}
	if schema.DefaultFunc == nil {
		return ""
	}
	result, _ := schema.DefaultFunc()
	if result == nil {
		return ""
	}
	srcEnvVars := []string{}
	switch key {
	case "dt_env_url", "automation_env_url":
		srcEnvVars = sourceEnvURLEnvVars
	case "dt_api_token":
		srcEnvVars = sourceAPITokenEnvVars
	case "client_id", "automation_client_id":
		srcEnvVars = sourceClientIDEnvVars
	case "account_id":
		srcEnvVars = sourceAccountIDEnvVars
	case "client_secret", "automation_client_secret":
		srcEnvVars = sourceClientSecretEnvVars
	case "platform_token":
		srcEnvVars = sourcePlatformTokenEnvVars
	}
	if len(srcEnvVars) > 0 {
		sourceValue := evalSourceEnv(sourceEnvURLEnvVars)
		if len(sourceValue) > 0 {
			return sourceValue
		}
	}

	return result
}

var sourceEnvURLEnvVars = []string{
	"DYNATRACE_SOURCE_ENV_URL",
	"DT_SOURCE_ENV_URL",
	"DYNATRACE_SOURCE_ENVIRONMENT_URL",
	"DT_SOURCE_ENVIRONMENT_URL",
}

var sourceAPITokenEnvVars = []string{
	"DYNATRACE_SOURCE_API_TOKEN",
	"DT_SOURCE_API_TOKEN",
}

var sourceClientIDEnvVars = []string{
	"DT_SOURCE_CLIENT_ID",
	"DYNATRACE_SOURCE_CLIENT_ID",
}

var sourceAccountIDEnvVars = []string{
	"DT_SOURCE_ACCOUNT_ID",
	"DYNATRACE_SOURCE_ACCOUNT_ID",
}

var sourceClientSecretEnvVars = []string{
	"DT_SOURCE_CLIENT_SECRET",
	"DYNATRACE_SOURCE_CLIENT_SECRET",
}

var sourcePlatformTokenEnvVars = []string{
	"DYNATRACE_SOURCE_PLATFORM_TOKEN",
	"DT_SOURCE_PLATFORM_TOKEN",
}

func evalSourceEnv(envVars []string) string {
	for _, envVar := range envVars {
		value := os.Getenv(envVar)
		if len(value) > 0 {
			return value
		}
	}
	return ""
}
