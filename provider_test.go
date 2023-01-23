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

package main_test

import (
	"context"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider
var testAccExternalProviders map[string]resource.ExternalProvider
var testAccProviderFactories = map[string]func() (*schema.Provider, error){
	"dynatrace": func() (*schema.Provider, error) {
		return provider.Provider(), nil
	},
}

func init() {
	testAccProvider = provider.Provider()
	testAccProviders = map[string]*schema.Provider{
		"dynatrace": testAccProvider,
	}
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"dynatrace": func() (*schema.Provider, error) {
			return provider.Provider(), nil
		},
	}
	testAccExternalProviders = map[string]resource.ExternalProvider{
		"dynatrace": {
			VersionConstraint: "1.0.2",
			Source:            "dynatrace.com/com/dynatrace",
		},
	}
}

func TestProvider(t *testing.T) {
	t.Skip()
	provider := provider.Provider()
	if err := provider.InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	t.Skip()
	var _ schema.Provider = *provider.Provider()
}

func TestProvider_configure(t *testing.T) {
	t.Skip()
	ctx := context.TODO()

	rc := terraform.NewResourceConfigRaw(map[string]any{})
	p := provider.Provider()
	diags := p.Configure(ctx, rc)
	if diags.HasError() {
		t.Fatal(diags)
	}
}
