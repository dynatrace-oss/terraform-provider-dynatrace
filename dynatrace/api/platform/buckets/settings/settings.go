package buckets

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type Bucket struct {
	Name          string `json:"bucketName"`       // The name / id of the bucket definition
	Table         Table  `json:"table"`            // The table the bucket definition applies to. Possible values are `logs`, `spans`,	`events` and `bizevents`
	DisplayName   string `json:"displayName"`      // The name of the bucket definition when visualized within the UI
	Status        Status `json:"status,omitempty"` // The current status of the bucket definition. Possible values are `creating`, `active`, `updating` and `deleting`
	RetentionDays int    `json:"retentionDays"`    // The retention in days
	Version       int    `json:"version"`          // The REST API keeps track of changes by increasing that value with every update
}

func (me *Bucket) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name / id of the bucket definition",
			Required:    true,
			ForceNew:    true,
		},
		"table": {
			Type:         schema.TypeString,
			Description:  "The table the bucket definition applies to. Possible values are `logs`, `spans`,	`events` and `bizevents`. Changing this attribute will result in deleting and re-creating the bucket definition",
			ValidateFunc: validation.StringInSlice([]string{`logs`, `spans`, `events`, `bizevents`}, false),
			Required:     true,
			ForceNew:     true,
		},
		"display_name": {
			Type:        schema.TypeString,
			Description: "The name of the bucket definition when visualized within the UI",
			Optional:    true,
		},
		"retention": {
			Type:        schema.TypeInt,
			Description: "The retention of stored data in days",
			Required:    true,
		},
		"status": {
			Type:        schema.TypeString,
			Description: "The status of the bucket definition. Usually has the value `active` unless an update or delete is currently happening",
			Computed:    true,
		},
	}
}

func (me *Bucket) ForUpdate() *BucketUpdate {
	return &BucketUpdate{
		DisplayName:   me.DisplayName,
		RetentionDays: me.RetentionDays,
	}
}

func (me *Bucket) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":         me.Name,
		"table":        me.Table,
		"display_name": me.DisplayName,
		"retention":    me.RetentionDays,
		"status":       me.Status,
	})
}

func (me *Bucket) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":         &me.Name,
		"table":        &me.Table,
		"display_name": &me.DisplayName,
		"retention":    &me.RetentionDays,
		"status":       &me.Status,
	})
}

type BucketUpdate struct {
	DisplayName   string `json:"displayName"`   // TODO
	RetentionDays int    `json:"retentionDays"` // TODO
}

type Status string

var Statuses = struct {
	Creating Status
	Active   Status
	Updating Status
	Deleting Status
}{
	"creating",
	"active",
	"updating",
	"deleting",
}

type Table string

var Tables = struct {
	Logs      Table
	Spans     Table
	Events    Table
	BizEvents Table
}{
	"logs",
	"spans",
	"events",
	"bizevents",
}
