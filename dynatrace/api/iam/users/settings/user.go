package users

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type User struct {
	Email  string   `json:"email"`
	Groups []string `json:"-"`
}

func (me *User) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"email": {
			Type:     schema.TypeString,
			Required: true,
		},
		"groups": {
			Type:     schema.TypeSet,
			Optional: true,
			MinItems: 1,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *User) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"email":  me.Email,
		"groups": me.Groups,
	})
}

func (me *User) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"email":  &me.Email,
		"groups": &me.Groups,
	})
}
