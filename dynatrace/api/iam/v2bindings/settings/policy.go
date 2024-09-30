package v2bindings

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Policies []*Policy

type Policy struct {
	ID         string
	Parameters map[string]string
	Metadata   map[string]string
}

func (me *Policy) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Either the attribute `id` or the attribute `uuid` of a `dynatrace_iam_policy`. Initially just the `id` attribute was supported (which is a concatenation of several configuration settings) - and is still supported for backwards compatibility",
		},
		"parameters": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"metadata": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *Policy) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("id", me.ID); err != nil {
		return err
	}
	if err := properties.Encode("parameters", me.Parameters); err != nil {
		return err
	}
	if err := properties.Encode("metadata", me.Metadata); err != nil {
		return err
	}
	return nil
}

func (me *Policy) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("id", &me.ID); err != nil {
		return err
	}
	if err := decoder.Decode("parameters", &me.Parameters); err != nil {
		return err
	}
	if err := decoder.Decode("metadata", &me.Metadata); err != nil {
		return err
	}
	return nil
}
