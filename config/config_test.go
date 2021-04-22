package config_test

import (
	"context"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/config"
)

type mockResourceData map[string]interface{}

func (mrd mockResourceData) Get(k string) interface{} {
	return mrd[k]
}

func TestProviderConfigure(t *testing.T) {
	ctx := context.TODO()
	d := mockResourceData{
		"dt_env_url":   "https://something.live.dynatrace.com",
		"dt_api_token": "faketoken",
	}

	result, _ := config.ProviderConfigureGeneric(ctx, d)
	configuration := result.(*config.ProviderConfiguration)
	if configuration.DTenvURL != "https://something.live.dynatrace.com/api/config/v1" {
		t.Fail()
	}
}
