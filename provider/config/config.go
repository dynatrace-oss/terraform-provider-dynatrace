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

// ProviderConfiguration contains credentials to communicate with the Dynatrace API
type ProviderConfiguration struct {
	EnvironmentURL  string
	ClusterAPIV2URL string
	ClusterAPIToken string
	APIToken        string
	IAM             IAM
	Platform        rest.PlatformCredentials

	// iamClient is built once in ProviderConfigureGeneric (when IAM credentials are available) and
	// reused by all IAM services. It is created there because that context outlives individual
	// requests, which the IAM client's OAuth token refresh relies on.
	iamClient rest.IAMClient
}

func (c *ProviderConfiguration) Credentials() *rest.Credentials {
	credentials := &rest.Credentials{
		Token:    c.APIToken,
		URL:      c.EnvironmentURL,
		IAM:      c.IAM,
		Platform: c.Platform,
	}
	credentials.Cluster.URL = c.ClusterAPIV2URL
	credentials.Cluster.Token = c.ClusterAPIToken
	return credentials
}

func (c *ProviderConfiguration) IAMClient() rest.IAMClient {
	return c.iamClient
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

var regexpProdTenant = regexp.MustCompile(`https:\/\/(.*).(live|apps).dynatrace.com`)

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
		ClusterAPIV2URL: getURLString(d, "dt_cluster_url"),
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

	// Build the shared IAM client once, using this long-lived context, when IAM credentials are
	// available. IAM services reuse it via pc.IAMClient() rather than constructing their own.
	if validateCredentials(pc, CredValIAM) == nil {
		pc.iamClient = rest.NewIAMClient(context.Background(), pc.Credentials())
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

// ClientSet validates the provider meta and returns it as a rest.ClientSet. *ProviderConfiguration
// implements rest.ClientSet (see client_set.go), so services/resources/datasources depend on the
// interface rather than the concrete credentials.
func ClientSet(m any, CredentialValidation int) (rest.ClientSet, error) {
	conf := m.(*ProviderConfiguration)
	if err := validateCredentials(conf, CredentialValidation); err != nil {
		return nil, err
	}
	return conf, nil
}

// getClassicEnvironmentURL retrieves the classic environment URL from the "dt_env_url" key in the provided configuration.
// It ensures that the URL is in the correct format for classic usage.
func getClassicEnvironmentURL(d Getter) string {
	return ensureClassicEnvironmentURL(getURLString(d, "dt_env_url"))
}

// getPlatformEnvironmentURL retrieves the platform environment URL from the provided configuration.
// It first checks for the "automation_env_url" key, and if not found, it derives it based on the classic environment URL.
func getPlatformEnvironmentURL(d Getter) string {
	if platformEnvironmentURL := getURLString(d, "automation_env_url"); platformEnvironmentURL != "" {
		return platformEnvironmentURL
	}
	return ensurePlatformEnvironmentURL(getClassicEnvironmentURL(d))
}

// getIAMEndpointURL retrieves the IAM API endpoint URL from the provided configuration.
// It first checks for the "iam_endpoint_url" key, and if not found, it derives it based on the environment URL.
func getIAMEndpointURL(d Getter) string {
	if iamEndpointURL := getURLString(d, "iam_endpoint_url"); iamEndpointURL != "" {
		return iamEndpointURL
	}
	return getIAMEndpointURLForEnvironment(getClassicEnvironmentURL(d))
}

// getIAMTokenURL retrieves the SSO token URL for IAM from the provided configuration.
// It first checks for the "iam_token_url" key, and if not found, it checks for the "token_url" key and then "automation_token_url" key.
// If neither is found, it derives the token URL based on the environment URL.
func getIAMTokenURL(d Getter) string {
	if iamTokenURL := getURLString(d, "iam_token_url"); iamTokenURL != "" {
		return iamTokenURL
	}

	if tokenURL := getURLString(d, "token_url"); tokenURL != "" {
		return tokenURL
	}

	if platformTokenURL := getURLString(d, "automation_token_url"); platformTokenURL != "" {
		return platformTokenURL
	}
	return getTokenURLForEnvironment(getClassicEnvironmentURL(d))
}

// getIAMClientID retrieves the OAuth client ID for IAM from the provided configuration.
// It first checks for the "iam_client_id" key, and if not found, it checks for the "client_id" and "automation_client_id" keys.
func getIAMClientID(d Getter) string {
	if iamClientID := getString(d, "iam_client_id"); iamClientID != "" {
		return iamClientID
	}
	if clientID := getString(d, "client_id"); clientID != "" {
		return clientID
	}
	return getString(d, "automation_client_id")
}

// getIAMClientSecret retrieves the OAuth client secret for IAM from the provided configuration.
// It first checks for the "iam_client_secret" key, and if not found, it checks for the "client_secret" and "automation_client_secret" keys.
func getIAMClientSecret(d Getter) string {
	if iamClientSecret := getString(d, "iam_client_secret"); iamClientSecret != "" {
		return iamClientSecret
	}
	if clientSecret := getString(d, "client_secret"); clientSecret != "" {
		return clientSecret
	}
	return getString(d, "automation_client_secret")
}

// getAccountID retrieves the account ID from the provided configuration.
// It first checks for the "iam_account_id" key, and if not found, it checks for the "account_id" key.
// The function returns the account ID with the "urn:dtaccount:" prefix removed, if present.
func getAccountID(d Getter) string {
	if iamAccountID := getString(d, "iam_account_id"); iamAccountID != "" {
		return strings.TrimPrefix(iamAccountID, "urn:dtaccount:")
	}
	if accountID := getString(d, "account_id"); accountID != "" {
		return strings.TrimPrefix(accountID, "urn:dtaccount:")
	}
	return ""
}

// getPlatformClientID retrieves the OAuth client ID for platform from the provided configuration.
// It first checks for the "automation_client_id" key, then "client_id", and finally "iam_client_id".
func getPlatformClientID(d Getter) string {
	if automationClientID := getString(d, "automation_client_id"); automationClientID != "" {
		return automationClientID
	}
	if clientID := getString(d, "client_id"); clientID != "" {
		return clientID
	}
	return getString(d, "iam_client_id")
}

// getPlatformClientSecret retrieves the OAuth client secret for platform from the provided configuration.
// It first checks for the "automation_client_secret" key, then "client_secret", and finally "iam_client
func getPlatformClientSecret(d Getter) string {
	if automation_client_secret := getString(d, "automation_client_secret"); automation_client_secret != "" {
		return automation_client_secret
	}
	if clientSecret := getString(d, "client_secret"); clientSecret != "" {
		return clientSecret
	}
	return getString(d, "iam_client_secret")
}

// getPlatformTokenURL returns the SSO token URL for platform based on the provided configuration.
// It first checks if a custom token URL is specified in the configuration using the "automation_token_url" key.
// If not, it checks the "token_url" and "iam_token_url" keys.
// If none are found, it derives the token URL based on the environment URL.
func getPlatformTokenURL(d Getter) string {
	if platformTokenURL := getURLString(d, "automation_token_url"); platformTokenURL != "" {
		return platformTokenURL
	}

	if tokenURL := getURLString(d, "token_url"); tokenURL != "" {
		return tokenURL
	}

	if iamTokenURL := getURLString(d, "iam_token_url"); iamTokenURL != "" {
		return iamTokenURL
	}

	return getTokenURLForEnvironment(getClassicEnvironmentURL(d))
}

// ensureClassicEnvironmentURL ensures that the provided Dynatrace environment URL is in the correct format for classic usage.
// It checks if the environment URL corresponds to a SaaS, Sprint, or Dev environment and returns the appropriate classic environment URL.
// If the environment URL does not match any of these, it returns the original environment URL.
func ensureClassicEnvironmentURL(dtEnvURL string) string {
	if envID, ok := extractProdEnvironmentID(dtEnvURL); ok {
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

// ensurePlatformEnvironmentURL ensures that the provided Dynatrace environment URL is in the correct format for platform usage.
// It checks if the environment URL corresponds to a SaaS, Sprint, or Dev environment and returns the appropriate platform environment URL.
// If the environment URL does not match any of these, it returns the original environment URL.
func ensurePlatformEnvironmentURL(dtEnvURL string) string {
	if envID, ok := extractProdEnvironmentID(dtEnvURL); ok {
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

// getTokenURLForEnvironment returns the token URL based on the provided Dynatrace environment URL.
// It checks if the environment URL corresponds to a SaaS, Sprint, or Dev environment and returns the appropriate token URL.
// If the environment URL does not match any of these, it defaults to the production token URL.
func getTokenURLForEnvironment(dtEnvURL string) string {
	if _, ok := extractProdEnvironmentID(dtEnvURL); ok {
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

// getIAMEndpointURLForEnvironment returns the IAM endpoint URL based on the provided Dynatrace environment URL.
// It checks if the environment URL corresponds to a SaaS, Sprint, or Dev environment and returns the appropriate IAM endpoint URL.
// If the environment URL does not match any of these, it defaults to the production IAM endpoint URL.
func getIAMEndpointURLForEnvironment(dtEnvURL string) string {
	if _, ok := extractProdEnvironmentID(dtEnvURL); ok {
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

// extractProdEnvironmentID extracts the environment ID from a Dynatrace SaaS environment URL.
// It returns the environment ID and a boolean indicating whether the extraction was successful.
func extractProdEnvironmentID(envURL string) (string, bool) {
	if match := regexpProdTenant.FindStringSubmatch(envURL); len(match) > 0 {
		return match[1], true
	}
	return "", false
}

// extractSprintEnvironmentID extracts the environment ID from a Dynatrace Sprint environment URL.
// It returns the environment ID and a boolean indicating whether the extraction was successful.
func extractSprintEnvironmentID(envURL string) (string, bool) {
	if match := regexpSprintTenant.FindStringSubmatch(envURL); len(match) > 0 {
		return match[1], true
	}
	return "", false
}

// extractDevEnvironmentID extracts the environment ID from a Dynatrace Dev environment URL.
// It returns the environment ID and a boolean indicating whether the extraction was successful.
func extractDevEnvironmentID(envURL string) (string, bool) {
	if match := regexpDevTenant.FindStringSubmatch(envURL); len(match) > 0 {
		return match[1], true
	}
	return "", false
}

// getURLString retrieves the URL string for the given key from the configuration,
// It removes trailing slashes and whitespace.
func getURLString(d Getter, key string) string {
	return strings.TrimSuffix(strings.TrimSpace(getString(d, key)), "/")
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
