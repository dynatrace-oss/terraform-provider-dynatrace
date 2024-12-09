package cities

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Regions []*Region

func (me *Regions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"region": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Resource{Schema: new(Region).Schema()},
		},
	}
}

func (me Regions) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("region", me)
}

func (me *Regions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("region", me)
}

type Region struct {
	Cities Cities `json:"cities"`
}

func (me *Region) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cities": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Resource{Schema: new(Cities).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
	}
}

func (me *Region) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"cities": me.Cities,
	})
}

func (me *Region) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"cities": &me.Cities,
	})
}
