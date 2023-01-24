package ddupool

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DDUPoolConfig struct {
	LimitEnabled bool    `json:"limitEnabled"`
	LimitType    *string `json:"limitType,omitempty"`
	LimitValue   *int    `json:"limitValue,omitempty"`
}

func (me *DDUPoolConfig) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Is the limit configuration enabled",
		},
		"type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Type of the limit applied: MONTHLY or ANNUAL",
		},
		"value": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Value of the DDU limit applied for provided timerange",
		},
	}
}

func (me *DDUPoolConfig) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]interface{}{
		"enabled": me.LimitEnabled,
		"type":    me.LimitType,
		"value":   me.LimitValue,
	})
}

func (me *DDUPoolConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	err := decoder.DecodeAll(map[string]interface{}{
		"enabled": &me.LimitEnabled,
		"type":    &me.LimitType,
		"value":   &me.LimitValue,
	})

	if err != nil {
		return err
	}

	//  Sanity check -> if limit_enabled is false, the type and value must not be sent
	if !me.LimitEnabled {
		me.LimitType = nil
		me.LimitValue = nil
	}

	return nil

}
