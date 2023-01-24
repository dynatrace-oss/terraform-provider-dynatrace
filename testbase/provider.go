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

package testbase

import (
	"context"
	"os"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var TestAccProvider *schema.Provider
var TestAccProviders map[string]*schema.Provider
var TestAccExternalProviders map[string]resource.ExternalProvider
var TestAccProviderFactories = map[string]func() (*schema.Provider, error){
	"dynatrace": func() (*schema.Provider, error) {
		return provider.Provider(), nil
	},
}

func init() {
	TestAccProvider = provider.Provider()
	TestAccProviders = map[string]*schema.Provider{
		"dynatrace": TestAccProvider,
	}
	TestAccProviderFactories = map[string]func() (*schema.Provider, error){
		"dynatrace": func() (*schema.Provider, error) {
			return provider.Provider(), nil
		},
	}
	TestAccExternalProviders = map[string]resource.ExternalProvider{
		"dynatrace": {
			VersionConstraint: "1.0.2",
			Source:            "dynatrace.com/com/dynatrace",
		},
	}
}

func TestAccPreCheck(t *testing.T) {
	ctx := context.TODO()

	if v := os.Getenv("DYNATRACE_ENV_URL"); v == "" {
		t.Fatalf("[WARN] DYNATRACE_ENV_URL has not been set for acceptance tests")
	}

	if v := os.Getenv("DYNATRACE_API_TOKEN"); v == "" {
		t.Fatalf("[WARN] DYNATRACE_API_TOKEN must be set for acceptance tests")
	}

	resourceConfig := terraform.NewResourceConfigRaw(nil)
	diags := TestAccProvider.Configure(ctx, resourceConfig)
	if diags.HasError() {
		t.Fatal(diags[0].Summary)
	}
}
