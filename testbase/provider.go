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
