package monitoring_config

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Scope     string `json:"-"`
	Extension string `json:"-"`
	Value     string `json:"value"`
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return nil
}
