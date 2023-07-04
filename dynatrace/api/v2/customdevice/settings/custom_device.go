package customdevice

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CustomDevice struct {
	EntityId *string `json:"entityId,omitempty"` // The ID of the custom device.
	// Type        *string        `json:"type,omitempty"`        // The type of the custom device.
	DisplayName *string `json:"displayName,omitempty"` // The name of the custom device, displayed in the UI.
	// Tags        Tags           `json:"tags,omitempty"`        // A set of tags assigned to the custom device.
	// Properties  map[string]any `json:"properties"`
	CustomDeviceID string `json:"customDeviceId,omitempty"`
}

type CustomDeviceList struct {
	Entities []*CustomDevice `json:"entities"` // An unordered list of custom devices
}

func (me *CustomDevice) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"entity_id": {
			Type:        schema.TypeString,
			Description: "The ID of the custom device.",
			Computed:    true,
		},
		"display_name": {
			Type:        schema.TypeString,
			Description: "The name of the custom device, displayed in the UI.",
			Required:    true,
		},
		"custom_device_id": {
			Type:        schema.TypeString,
			Description: "The unique name of the custom device.",
			Optional:    true,
			Computed:    true,
		},
	}
}

func (me *CustomDevice) MarshalHCL(properties hcl.Properties) error {
	if err := properties.EncodeAll(map[string]any{
		"entity_id":        me.EntityId,
		"display_name":     me.DisplayName,
		"custom_device_id": me.CustomDeviceID,
	}); err != nil {
		return err
	}
	return nil
}

func (me *CustomDevice) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"entity_id":        &me.EntityId,
		"display_name":     &me.DisplayName,
		"custom_device_id": &me.CustomDeviceID,
	})
}

func (me *CustomDevice) Name() string {
	return *me.DisplayName
}
