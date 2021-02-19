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
	"fmt"
	"os"
	"testing"

	main "github.com/dynatrace-oss/terraform-provider-dynatrace"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type currentEnv struct {
	Config                  string
	DynatraceEnvironmentURL string
	DynatraceAPIToken       string
}

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider
var testAccExternalProviders map[string]resource.ExternalProvider
var testAccProviderFactories = map[string]func() (*schema.Provider, error){
	"dynatrace": func() (*schema.Provider, error) {
		return main.Provider(), nil
	},
}

func init() {
	testAccProvider = main.Provider()
	testAccProviders = map[string]*schema.Provider{
		"dynatrace": testAccProvider,
	}
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"dynatrace": func() (*schema.Provider, error) {
			return main.Provider(), nil
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
	provider := main.Provider()
	if err := provider.InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ schema.Provider = *main.Provider()
}

func TestProvider_configure(t *testing.T) {
	ctx := context.TODO()

	rc := terraform.NewResourceConfigRaw(map[string]interface{}{})
	p := main.Provider()
	diags := p.Configure(ctx, rc)
	if diags.HasError() {
		t.Fatal(diags)
	}
}

func unsetEnv(t *testing.T) func() {
	e := getEnv()

	if err := os.Unsetenv("DYNATRACE_ENV_URL"); err != nil {
		t.Fatalf("Error unsetting env var DYNATRACE_ENV_URL: %s", err)
	}
	if err := os.Unsetenv("DYNATRACE_API_TOKEN"); err != nil {
		t.Fatalf("Error unsetting env var DYNATRACE_API_TOKEN: %s", err)
	}

	return func() {
		if err := os.Setenv("DYNATRACE_ENV_URL", e.Config); err != nil {
			t.Fatalf("Error resetting env var DYNATRACE_ENV_URL: %s", err)
		}
		if err := os.Setenv("DYNATRACE_API_TOKEN", e.Config); err != nil {
			t.Fatalf("Error resetting env var DYNATRACE_API_TOKEN: %s", err)
		}
	}
}

func getEnv() *currentEnv {
	e := &currentEnv{
		DynatraceEnvironmentURL: os.Getenv("DYNATRACE_ENV_URL"),
		DynatraceAPIToken:       os.Getenv("DYNATRACE_API_TOKEN"),
	}

	return e
}

func testAccPreCheck(t *testing.T) {
	ctx := context.TODO()

	if v := os.Getenv("DYNATRACE_ENV_URL"); v == "" {
		t.Fatalf("[WARN] DYNATRACE_ENV_URL has not been set for acceptance tests")
	}

	if v := os.Getenv("DYNATRACE_API_TOKEN"); v == "" {
		t.Fatalf("[WARN] DYNATRACE_API_TOKEN must be set for acceptance tests")
	}

	diags := testAccProvider.Configure(ctx, terraform.NewResourceConfigRaw(nil))
	if diags.HasError() {
		t.Fatal(diags[0].Summary)
	}
	return
}

func requiredProviders() string {
	return fmt.Sprintf(`terraform {
  required_providers {
    dynatrace = {
      source  = "dynatrace.com/com/dynatrace"
      version = "1.0.2"
    }
  }
}
`)
}
