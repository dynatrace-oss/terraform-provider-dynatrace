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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// HTTPVerbose if set to `true` terraform-provider-dynatrace.log will contain request and response payload
var HTTPVerbose = (strings.TrimSpace(os.Getenv("DYNATRACE_DEBUG")) == "true")

type IAM struct {
	ClientID     string
	AccountID    string
	ClientSecret string
}

type Automation struct {
	ClientID       string
	ClientSecret   string
	TokenURL       string
	EnvironmentURL string
}

const (
	CredValDefault = iota
	CredValIAM
	CredValCluster
	CredValAutomation
	CredValNone
)

func validateCredentials(conf *ProviderConfiguration, CredentialValidation int) error {
	switch CredentialValidation {
	case CredValDefault:
		if len(conf.EnvironmentURL) == 0 {
			return fmt.Errorf("No Environment URL has been specified. Use either the environment variable `DYNATRACE_ENV_URL` or the configuration attribute `dt_env_url` of the provider for that.")
		}
		if !strings.HasPrefix(conf.EnvironmentURL, "https://") && !strings.HasPrefix(conf.EnvironmentURL, "http://") {
			return fmt.Errorf("The Environment URL `%s` neither starts with `https://` nor with `http://`. Please check your configuration.\nFor SaaS environments: `https://######.live.dynatrace.com`.\nFor Managed environments: `https://############/e/########-####-####-####-############`", conf.EnvironmentURL)
		}
		if len(conf.APIToken) == 0 {
			return fmt.Errorf("No API Token has been specified. Use either the environment variable `DYNATRACE_API_TOKEN` or the configuration attribute `dt_api_token` of the provider for that.")
		}
	case CredValIAM:
		if len(conf.IAM.AccountID) == 0 {
			return fmt.Errorf("No OAuth Account ID has been specified. Use either the environment variable `DT_ACCOUNT_ID` or the configuration attribute `iam_account_id` of the provider for that.")
		}
		if len(conf.IAM.ClientID) == 0 {
			return fmt.Errorf("No OAuth Client ID has been specified. Use either the environment variable `DT_CLIENT_ID` or the configuration attribute `iam_client_id` of the provider for that.")
		}
		if len(conf.IAM.ClientSecret) == 0 {
			return fmt.Errorf("No OAuth Client Secret has been specified. Use either the environment variable `DT_CLIENT_SECRET` or the configuration attribute `iam_client_secret` of the provider for that.")
		}
	case CredValCluster:
		if len(conf.ClusterAPIToken) == 0 {
			return fmt.Errorf("No Cluster API Token has been specified. Use either the environment variable `DT_CLUSTER_API_TOKEN` or the configuration attribute `dt_cluster_api_token` of the provider for that.")
		}
		if len(conf.ClusterAPIV2URL) == 0 {
			return fmt.Errorf("No Cluster URL has been specified. Use either the environment variable `DT_CLUSTER_URL` or the configuration attribute `dt_cluster_url` of the provider for that.")
		}
	case CredValAutomation:
		if len(conf.Automation.ClientID) == 0 {
			return fmt.Errorf("No OAuth Client ID for the Automation API has been specified. Use either the environment variable `DT_AUTOMATION_CLIENT_ID` or the configuration attribute `automation_client_id` of the provider for that.")
		}
		if len(conf.Automation.ClientSecret) == 0 {
			return fmt.Errorf("No OAuth Client Secret for the Automation API has been specified. Use either the environment variable `DT_AUTOMATION_CLIENT_SECRET` or the configuration attribute `automation_client_secret` of the provider for that.")
		}
		if len(conf.Automation.TokenURL) == 0 {
			return fmt.Errorf("No Token URL for the Automation API has been specified. Use either the environment variable `DT_AUTOMATION_TOKEN_URL` or the configuration attribute `automation_token_url` of the provider for that.")
		}
		if len(conf.Automation.EnvironmentURL) == 0 {
			return fmt.Errorf("No Environment URL for the Automation API has been specified. Use either the environment variable `DT_AUTOMATION_ENVIRONMENT_URL` or the configuration attribute `automation_env_url` of the provider for that.")
		}
	}
	return nil
}

func Credentials(m any, CredentialValidation int) (*settings.Credentials, error) {
	conf := m.(*ProviderConfiguration)
	if err := validateCredentials(conf, CredentialValidation); err != nil {
		return nil, err
	}
	return &settings.Credentials{
		Token:      conf.APIToken,
		URL:        conf.EnvironmentURL,
		IAM:        conf.IAM,
		Automation: conf.Automation,
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
	Automation        Automation
}

type Getter interface {
	Get(key string) any
}

func ProviderConfigure(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
	return ProviderConfigureGeneric(ctx, d)
}

func ProviderConfigureGeneric(ctx context.Context, d Getter) (any, diag.Diagnostics) {
	dtEnvURL := d.Get("dt_env_url").(string)
	apiToken := d.Get("dt_api_token").(string)
	clusterAPIToken := getString(d, "dt_cluster_api_token")
	clusterURL := getString(d, "dt_cluster_url")

	dtEnvURL = strings.TrimSuffix(strings.TrimSuffix(dtEnvURL, " "), "/")
	clusterURL = strings.TrimSuffix(strings.TrimSuffix(clusterURL, " "), "/")

	fullURL := dtEnvURL + "/api/config/v1"
	fullNonConfigURL := dtEnvURL + "/api/v1"
	fullApiV2URL := dtEnvURL + "/api/v2"

	automationEnvironmentURL := getString(d, "automation_env_url")
	automationTokenURL := getString(d, "automation_token_url")
	if len(automationEnvironmentURL) == 0 {
		re := regexp.MustCompile(`https:\/\/(.*).(live|apps).dynatrace.com`)
		if match := re.FindStringSubmatch(dtEnvURL); match != nil && len(match) > 0 {
			automationEnvironmentURL = fmt.Sprintf("https://%s.apps.dynatrace.com", match[1])
			automationTokenURL = "https://sso.dynatrace.com/sso/oauth2/token"
		}
	}

	var diags diag.Diagnostics

	return &ProviderConfiguration{
		EnvironmentURL:    dtEnvURL,
		DTenvURL:          fullURL,
		DTApiV2URL:        fullApiV2URL,
		DTNonConfigEnvURL: fullNonConfigURL,
		APIToken:          apiToken,
		ClusterAPIToken:   clusterAPIToken,
		ClusterAPIV2URL:   clusterURL,
		IAM: IAM{
			ClientID:     getString(d, "iam_client_id"),
			AccountID:    getString(d, "iam_account_id"),
			ClientSecret: getString(d, "iam_client_secret"),
		},
		Automation: Automation{
			ClientID:       getString(d, "automation_client_id"),
			ClientSecret:   getString(d, "automation_client_secret"),
			TokenURL:       automationTokenURL,
			EnvironmentURL: automationEnvironmentURL,
		},
	}, diags
}

func getString(d Getter, key string) string {
	if value := d.Get(key); value != nil {
		return value.(string)
	}
	return ""
}
