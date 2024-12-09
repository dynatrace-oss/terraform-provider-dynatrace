package countries

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Countries []*Country

func (me *Countries) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"country": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Resource{Schema: new(Country).Schema()},
		},
	}
}

func (me Countries) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("country", me)
}

func (me *Countries) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("country", me)
}

type Country struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func (me *Country) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"code": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

func (me *Country) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name": me.Name,
		"code": me.Code,
	})
}

func (me *Country) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name": &me.Name,
		"code": &me.Code,
	})
}
