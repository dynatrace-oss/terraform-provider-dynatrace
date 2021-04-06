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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// HTTPVerbose if set to `true` terraform-provider-dynatrace.log will contain request and response payload
const HTTPVerbose = false

// ProviderConfiguration contains the initialized API clients to communicate with the Dynatrace API
type ProviderConfiguration struct {
	DTenvURL string
	APIToken string
}

// ProviderConfigure has no documentation
func ProviderConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	dtEnvURL := d.Get("dt_env_url").(string)
	apiToken := d.Get("dt_api_token").(string)

	fullURL := dtEnvURL + "/api/config/v1"
	var diags diag.Diagnostics

	return &ProviderConfiguration{
		DTenvURL: fullURL,
		APIToken: apiToken,
	}, diags

}
