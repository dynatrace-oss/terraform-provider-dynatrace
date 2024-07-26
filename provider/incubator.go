package provider

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/resources/goldenstate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func incubator(provider *schema.Provider) {
	provider.ResourcesMap["dynatrace_golden_state"] = goldenstate.Resource()
}
