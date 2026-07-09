/**
* @license
* Copyright 2026 Dynatrace LLC
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
	"regexp"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/envutils"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type IAM struct {
	ClientID     string
	AccountID    string
	ClientSecret string
	TokenURL     string
	EndpointURL  string
}

// ProviderConfiguration contains the initialized API clients to communicate with the Dynatrace API
type ProviderConfiguration struct {
	EnvironmentURL  string
	ClusterAPIV2URL string
	ClusterAPIToken string
	APIToken        string
	IAM             IAM
	Platform        rest.PlatformCredentials
}

type Getter interface {
	Get(key string) any
}

type ConfigGetter struct {
	Provider *schema.Provider
}

const (
	CredValDefault = iota
	CredValIAM
	CredValCluster
	CredValPlatform
	CredValNone
	CredValExport
	CredValExportIAM
)

const (
	ProdTokenURL   = "https://sso.dynatrace.com/sso/oauth2/token"
	SprintTokenURL = "https://sso-sprint.dynatracelabs.com/sso/oauth2/token"
	DevTokenURL    = "https://sso-dev.dynatracelabs.com/sso/oauth2/token"

	ProdIAMEndpointURL   = "https://api.dynatrace.com"
	SprintIAMEndpointURL = "https://api-hardening.internal.dynatracelabs.com"
	DevIAMEndpointURL    = "https://api-dev.internal.dynatracelabs.com"
)

var regexpSaasTenant = regexp.MustCompile(`https:\/\/(.*).(live|apps).dynatrace.com`)

var regexpSprintTenant = regexp.MustCompile(`https:\/\/(.*).sprint(?:\.apps)?.dynatracelabs.com`)

var regexpDevTenant = regexp.MustCompile(`https:\/\/(.*).dev(?:\.apps)?.dynatracelabs.com`)

func ProviderConfigure(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
	return ProviderConfigureGeneric(ctx, d)
}

func ProviderConfigureGeneric(ctx context.Context, d Getter) (any, diag.Diagnostics) {
	pc := &ProviderConfiguration{
		EnvironmentURL:  getClassicEnvironmentURL(d),
		APIToken:        getString(d, "dt_api_token"),
		ClusterAPIToken: getString(d, "dt_cluster_api_token"),
		ClusterAPIV2URL: cleanURL(getString(d, "dt_cluster_url")),
		IAM: IAM{
			ClientID:     getIAMClientID(d),
			AccountID:    getAccountID(d),
			ClientSecret: getIAMClientSecret(d),
			TokenURL:     getIAMTokenURL(d),
			EndpointURL:  getIAMEndpointURL(d),
		},
		Platform: rest.PlatformCredentials{
			PlatformToken:  getString(d, "platform_token"),
			ClientID:       getPlatformClientID(d),
			ClientSecret:   getPlatformClientSecret(d),
			TokenURL:       getPlatformTokenURL(d),
			EnvironmentURL: getPlatformEnvironmentURL(d),
		},
	}
	return pc, diag.Diagnostics{}
}

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
	case CredValPlatform:
		if len(conf.Platform.ClientID) == 0 {
			return fmt.Errorf(" No OAuth Client ID for the Automation API has been specified. Use either the environment variable `DT_AUTOMATION_CLIENT_ID` or the configuration attribute `automation_client_id` of the provider for that")
		}
		if len(conf.Platform.ClientSecret) == 0 {
			return fmt.Errorf(" No OAuth Client Secret for the Automation API has been specified. Use either the environment variable `DT_AUTOMATION_CLIENT_SECRET` or the configuration attribute `automation_client_secret` of the provider for that")
		}
		if len(conf.Platform.TokenURL) == 0 {
			return fmt.Errorf(" No Token URL for the Automation API has been specified. Use either the environment variable `DT_AUTOMATION_TOKEN_URL` or the configuration attribute `automation_token_url` of the provider for that")
		}
		if len(conf.Platform.EnvironmentURL) == 0 {
			return fmt.Errorf(" No Environment URL for the Automation API has been specified. Use either the environment variable `DT_AUTOMATION_ENVIRONMENT_URL` or the configuration attribute `automation_env_url` of the provider for that")
		}
	case CredValExport:
		if len(conf.EnvironmentURL) == 0 {
			return fmt.Errorf(" No Environment URL has been specified. Use either the environment variable `DYNATRACE_ENV_URL` or the configuration attribute `dt_env_url` of the provider for that")
		}
		if len(conf.APIToken) == 0 && len(conf.Platform.PlatformToken) == 0 && validateCredentials(conf, CredValPlatform) != nil {
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

// Validate asserts that m is a *ProviderConfiguration, runs the requested credential validation and
// returns the configuration so callers can use it directly to construct clients.
func Validate(m any, CredentialValidation int) (*ProviderConfiguration, error) {
	conf := m.(*ProviderConfiguration)
	if err := validateCredentials(conf, CredentialValidation); err != nil {
		return nil, err
	}
	return conf, nil
}

// ContainsOAuth reports whether Platform OAuth client credentials have been provided.
func (c *ProviderConfiguration) ContainsOAuth() bool {
	return len(c.Platform.ClientID) > 0 && len(c.Platform.ClientSecret) > 0
}

// ContainsPlatformToken reports whether a platform token has been provided.
func (c *ProviderConfiguration) ContainsPlatformToken() bool {
	return len(c.Platform.PlatformToken) > 0
}

// ContainsOAuthOrPlatformToken reports whether either Platform OAuth credentials or a platform token
// have been provided.
func (c *ProviderConfiguration) ContainsOAuthOrPlatformToken() bool {
	return c.ContainsOAuth() || c.ContainsPlatformToken()
}

func getClassicEnvironmentURL(d Getter) string {
	return ensureClassicEnvironmentURL(cleanURL(getString(d, "dt_env_url")))
}

func getPlatformEnvironmentURL(d Getter) string {
	if platformEnvironmentURL := cleanURL(getString(d, "automation_env_url")); platformEnvironmentURL != "" {
		return platformEnvironmentURL
	}
	return ensurePlatformEnvironmentURL(getClassicEnvironmentURL(d))
}

func getIAMEndpointURL(d Getter) string {
	if iamEndpointURL := cleanURL(getString(d, "iam_endpoint_url")); iamEndpointURL != "" {
		return iamEndpointURL
	}
	return getIAMEndpointURLForEnvironment(getClassicEnvironmentURL(d))
}

func getIAMTokenURL(d Getter) string {
	if iamTokenURL := cleanURL(getString(d, "iam_token_url")); iamTokenURL != "" {
		return iamTokenURL
	}

	if tokenURL := cleanURL(getString(d, "token_url")); tokenURL != "" {
		return tokenURL
	}
	return getPlatformTokenURL(d)
}

func getIAMClientID(d Getter) string {
	iamClientID := getString(d, "iam_client_id")
	if iamClientID == "" {
		iamClientID = getString(d, "client_id")
	}
	if iamClientID == "" {
		iamClientID = getString(d, "automation_client_id")
	}
	return iamClientID
}

func getIAMClientSecret(d Getter) string {
	iamClientSecret := getString(d, "iam_client_secret")
	if iamClientSecret == "" {
		iamClientSecret = getString(d, "client_secret")
	}
	if iamClientSecret == "" {
		iamClientSecret = getString(d, "automation_client_secret")
	}
	return iamClientSecret
}

func getAccountID(d Getter) string {
	iamAccountID := getString(d, "iam_account_id")
	if iamAccountID == "" {
		iamAccountID = getString(d, "account_id")
	}

	return strings.TrimPrefix(iamAccountID, "urn:dtaccount:")
}

func getPlatformClientID(d Getter) string {
	platformClientID := getString(d, "automation_client_id")
	if platformClientID == "" {
		platformClientID = getString(d, "client_id")
	}
	if platformClientID == "" {
		platformClientID = getString(d, "iam_client_id")
	}
	return platformClientID
}

func getPlatformClientSecret(d Getter) string {
	platformClientSecret := getString(d, "automation_client_secret")
	if platformClientSecret == "" {
		platformClientSecret = getString(d, "client_secret")
	}
	if platformClientSecret == "" {
		platformClientSecret = getString(d, "iam_client_secret")
	}
	return platformClientSecret
}

func getPlatformTokenURL(d Getter) string {
	if platformTokenURL := getString(d, "automation_token_url"); platformTokenURL != "" {
		return platformTokenURL
	}
	return getTokenURLForEnvironment(getClassicEnvironmentURL(d))
}

func ensureClassicEnvironmentURL(dtEnvURL string) string {
	dtEnvURL = strings.TrimSuffix(strings.TrimSuffix(dtEnvURL, " "), "/")
	if dtEnvURL == "" {
		return ""
	}

	if envID, ok := extractSaasEnvironmentID(dtEnvURL); ok {
		return fmt.Sprintf("https://%s.live.dynatrace.com", envID)
	}

	if envID, ok := extractSprintEnvironmentID(dtEnvURL); ok {
		return fmt.Sprintf("https://%s.sprint.dynatracelabs.com", envID)
	}

	if envID, ok := extractDevEnvironmentID(dtEnvURL); ok {
		return fmt.Sprintf("https://%s.dev.dynatracelabs.com", envID)
	}

	return dtEnvURL
}

func ensurePlatformEnvironmentURL(dtEnvURL string) string {
	if envID, ok := extractSaasEnvironmentID(dtEnvURL); ok {
		return fmt.Sprintf("https://%s.apps.dynatrace.com", envID)
	}
	if envID, ok := extractSprintEnvironmentID(dtEnvURL); ok {
		return fmt.Sprintf("https://%s.sprint.apps.dynatracelabs.com", envID)
	}
	if envID, ok := extractDevEnvironmentID(dtEnvURL); ok {
		return fmt.Sprintf("https://%s.dev.apps.dynatracelabs.com", envID)
	}

	return dtEnvURL
}

func getTokenURLForEnvironment(dtEnvURL string) string {
	if _, ok := extractSaasEnvironmentID(dtEnvURL); ok {
		return ProdTokenURL
	}
	if _, ok := extractSprintEnvironmentID(dtEnvURL); ok {
		return SprintTokenURL
	}
	if _, ok := extractDevEnvironmentID(dtEnvURL); ok {
		return DevTokenURL
	}
	return ProdTokenURL
}

func getIAMEndpointURLForEnvironment(dtEnvURL string) string {
	if _, ok := extractSaasEnvironmentID(dtEnvURL); ok {
		return ProdIAMEndpointURL
	}
	if _, ok := extractSprintEnvironmentID(dtEnvURL); ok {
		return SprintIAMEndpointURL
	}
	if _, ok := extractDevEnvironmentID(dtEnvURL); ok {
		return DevIAMEndpointURL
	}
	return ProdIAMEndpointURL
}

func extractSaasEnvironmentID(envURL string) (string, bool) {
	if match := regexpSaasTenant.FindStringSubmatch(envURL); len(match) > 0 {
		return match[1], true
	}
	return "", false
}

func extractSprintEnvironmentID(envURL string) (string, bool) {
	if match := regexpSprintTenant.FindStringSubmatch(envURL); len(match) > 0 {
		return match[1], true
	}
	return "", false
}

func extractDevEnvironmentID(envURL string) (string, bool) {
	if match := regexpDevTenant.FindStringSubmatch(envURL); len(match) > 0 {
		return match[1], true
	}
	return "", false
}

func cleanURL(url string) string {
	return strings.TrimSuffix(strings.TrimSpace(url), "/")
}

func getString(d Getter, key string) string {
	if value := d.Get(key); value != nil {
		return value.(string)
	}
	return ""
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
	var sourceEnvVar *envutils.MultiStringEnvVar
	switch key {
	case "dt_env_url", "automation_env_url":
		sourceEnvVar = &envutils.DynatraceSourceEnvURL
	case "dt_api_token":
		sourceEnvVar = &envutils.DynatraceSourceAPIToken
	case "client_id", "automation_client_id":
		sourceEnvVar = &envutils.DynatraceSourceClientID
	case "account_id":
		sourceEnvVar = &envutils.DynatraceSourceAccountID
	case "client_secret", "automation_client_secret":
		sourceEnvVar = &envutils.DynatraceSourceClientSecret
	case "platform_token":
		sourceEnvVar = &envutils.DynatraceSourcePlatformToken
	}
	if sourceEnvVar != nil {
		if sourceValue := sourceEnvVar.Get(); len(sourceValue) > 0 {
			return sourceValue
		}
	}

	return result
}
