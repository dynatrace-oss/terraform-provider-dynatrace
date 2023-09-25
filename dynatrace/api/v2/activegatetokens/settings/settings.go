package activegatetokens

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type Settings struct {
	Name           string  `json:"name,omitempty"`
	Type           string  `json:"activeGateType"`
	ExpirationDate *string `json:"expirationDate,omitempty"`
	SeedToken      bool    `json:"seedToken,omitempty"`
	Token          *string `json:"token,omitempty"`
	TenantToken    *string `json:"-"`
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the token.",
			Required:    true,
		},
		"type": {
			Type:         schema.TypeString,
			Description:  "The type of the ActiveGate for which the token is valid. Possible values are `ENVIRONMENT` or `CLUSTER`",
			ValidateFunc: validation.StringInSlice([]string{"ENVIRONMENT", "CLUSTER"}, false),
			Required:     true,
		},
		"expiration_date": {
			Type:             schema.TypeString,
			Description:      "The expiration date of the token.\n\n    You can use one of the following formats:\n\n    * Timestamp in UTC milliseconds.\n    * Human-readable format of 2021-01-25T05:57:01.123+01:00. If no time zone is specified, UTC is used. You can use a space character instead of the T. Seconds and fractions of a second are optional.\n    * Relative timeframe, back from now. The format is now-NU/A, where N is the amount of time, U is the unit of time, and A is an alignment. The alignment rounds all the smaller values to the nearest zero in the past. For example, now-1y/w is one year back, aligned by a week. You can also specify relative timeframe without an alignment: now-NU. Supported time units for the relative timeframe are:\n      - m: minutes\n      - h: hours\n      - d: days\n      - w: weeks\n      - M: months\n      - y: years",
			Optional:         true,
			DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool { return len(oldValue) > 0 },
		},
		"token": {
			Type:        schema.TypeString,
			Description: "The secret of the token.",
			Sensitive:   true,
			Computed:    true,
		},
		"tenant_token": {
			Type:        schema.TypeString,
			Description: "The tenant token. This information isn't directly related to the Active Gate Token. It's included for convenience. You require the permission `InstallerDownload` for that attribute to get populated",
			Sensitive:   true,
			Computed:    true,
		},
		"seed": {
			Type:        schema.TypeBool,
			Description: "The token is a seed token (true) or an individual token (false). We recommend the individual token option (false)",
			Optional:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	if err := properties.EncodeAll(map[string]any{
		"name":            me.Name,
		"type":            me.Type,
		"expiration_date": me.ExpirationDate,
		"seed":            me.SeedToken,
		"tenant_token":    me.TenantToken,
	}); err != nil {
		return err
	}
	if me.Token != nil && len(*me.Token) > 0 {
		if err := properties.Encode("token", me.Token); err != nil {
			return err
		}
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":            &me.Name,
		"type":            &me.Type,
		"expiration_date": &me.ExpirationDate,
		"seed":            &me.SeedToken,
	})
}
