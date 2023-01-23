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
	"os"
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

func Credentials(m any) *settings.Credentials {
	conf := m.(*ProviderConfiguration)
	return &settings.Credentials{
		Token: conf.APIToken,
		URL:   conf.EnvironmentURL,
	}
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
	}, diags
}

func getString(d Getter, key string) string {
	if value := d.Get(key); value != nil {
		return value.(string)
	}
	return ""
}
