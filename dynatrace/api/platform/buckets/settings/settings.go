package buckets

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Bucket struct {
	Name          string  `json:"bucketName"`       // TODO
	Table         Table   `json:"table"`            // TODO
	DisplayName   string  `json:"displayName"`      // TODO
	Status        *Status `json:"status,omitempty"` // TODO
	RetentionDays int     `json:"retentionDays"`    // TODO
	Version       int     `json:"version"`          // TODO
}

func (me *Bucket) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "TODO",
			Required:    true,
			ForceNew:    true,
		},
		"table": {
			Type:        schema.TypeString,
			Description: "TODO",
			Required:    true,
			ForceNew:    true,
		},
		"display_name": {
			Type:        schema.TypeString,
			Description: "TODO",
			Optional:    true,
		},
		"retention": {
			Type:        schema.TypeInt,
			Description: "The retention in days",
			Required:    true,
		},
		"status": {
			Type:        schema.TypeString,
			Description: "TODO",
			Computed:    true,
		},
		"version": {
			Type:        schema.TypeInt,
			Description: "TODO",
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
		"version":      me.Version,
	})
}

func (me *Bucket) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":         &me.Name,
		"table":        &me.Table,
		"display_name": &me.DisplayName,
		"retention":    &me.RetentionDays,
		"status":       &me.Status,
		"version":      &me.Version,
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
