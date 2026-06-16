/**
* @license
* Copyright 2026 Dynatrace LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package dataforwarding

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	AwsConnection        *AwsConnection     `json:"awsConnection,omitempty"`        // AWS Connection
	AzureConnection      *AzureConnection   `json:"azureConnection,omitempty"`      // Azure Connection
	BuiltinIngestSources []string           `json:"builtinIngestSources,omitempty"` // List of built-in ingest sources
	BuiltinPipelines     []string           `json:"builtinPipelines,omitempty"`     // Built-in pipelines
	BulkPattern          string             `json:"bulkPattern"`                    // Segmentation and prefix of the data
	BulkSize             *int               `json:"bulkSize,omitempty"`             // Bulk size for transmission
	CloudVendorType      CloudVendorType    `json:"cloudVendorType"`                // Cloud Vendor Type. Possible values: `aws`, `azure`, `gcp`
	DataForwardingType   DataForwardingType `json:"dataForwardingType"`             // Pipeline Type. Possible values: `processed`, `raw`
	Enabled              bool               `json:"enabled"`                        // This setting is enabled (`true`) or disabled (`false`)
	ForwardingName       string             `json:"forwardingName"`                 // Forwarding name
	GcpConnection        *GcpConnection     `json:"gcpConnection,omitempty"`        // GCP Connection
	IngestSources        []string           `json:"ingestSources,omitempty"`        // List of ingest sources
	Matcher              string             `json:"matcher"`                        // Query which determines whether the record should be routed to the target pipeline of this rule.
	Pipelines            []string           `json:"pipelines,omitempty"`            // Pipelines
	Processing           *Stage             `json:"processing,omitempty"`           // Processing
}

func (me *Settings) Name() string {
	return me.ForwardingName
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"aws_connection": {
			Type:        schema.TypeList,
			Description: "AWS Connection",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(AwsConnection).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"azure_connection": {
			Type:        schema.TypeList,
			Description: "Azure Connection",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(AzureConnection).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"builtin_ingest_sources": {
			Type:        schema.TypeSet,
			Description: "List of built-in ingest sources",
			Optional:    true, // precondition & minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"builtin_pipelines": {
			Type:        schema.TypeSet,
			Description: "Built-in pipelines",
			Optional:    true, // precondition & minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"bulk_pattern": {
			Type:        schema.TypeString,
			Description: "Segmentation and prefix of the data",
			Required:    true,
		},
		"bulk_size": {
			Type:        schema.TypeInt,
			Description: "Bulk size for transmission",
			Optional:    true, // nullable
		},
		"cloud_vendor_type": {
			Type:        schema.TypeString,
			Description: "Cloud Vendor Type. Possible values: `aws`, `azure`, `gcp`",
			Required:    true,
		},
		"data_forwarding_type": {
			Type:        schema.TypeString,
			Description: "Pipeline Type. Possible values: `processed`, `raw`",
			Required:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"forwarding_name": {
			Type:        schema.TypeString,
			Description: "Forwarding name",
			Required:    true,
		},
		"gcp_connection": {
			Type:        schema.TypeList,
			Description: "GCP Connection",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(GcpConnection).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"ingest_sources": {
			Type:        schema.TypeSet,
			Description: "List of ingest sources",
			Optional:    true, // precondition & minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"matcher": {
			Type:        schema.TypeString,
			Description: "Query which determines whether the record should be routed to the target pipeline of this rule.",
			Required:    true,
		},
		"pipelines": {
			Type:        schema.TypeSet,
			Description: "Pipelines",
			Optional:    true, // precondition & minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"processing": {
			Type:        schema.TypeList,
			Description: "Processing",
			Optional:    true, // nullable
			Elem:        &schema.Resource{Schema: new(Stage).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"aws_connection":         me.AwsConnection,
		"azure_connection":       me.AzureConnection,
		"builtin_ingest_sources": me.BuiltinIngestSources,
		"builtin_pipelines":      me.BuiltinPipelines,
		"bulk_pattern":           me.BulkPattern,
		"bulk_size":              me.BulkSize,
		"cloud_vendor_type":      me.CloudVendorType,
		"data_forwarding_type":   me.DataForwardingType,
		"enabled":                me.Enabled,
		"forwarding_name":        me.ForwardingName,
		"gcp_connection":         me.GcpConnection,
		"ingest_sources":         me.IngestSources,
		"matcher":                me.Matcher,
		"pipelines":              me.Pipelines,
		"processing":             me.Processing,
	})
}

func (me *Settings) HandlePreconditions() error {
	if (me.AwsConnection != nil) && (string(me.CloudVendorType) != "aws") {
		return fmt.Errorf("'aws_connection' must not be specified unless 'cloud_vendor_type' is set to 'aws'; got 'cloud_vendor_type'='%v'", me.CloudVendorType)
	}
	if (me.AwsConnection == nil) && (string(me.CloudVendorType) == "aws") {
		return fmt.Errorf("'aws_connection' must be specified when 'cloud_vendor_type' is set to 'aws'; got 'cloud_vendor_type'='%v'", me.CloudVendorType)
	}
	if (me.AzureConnection != nil) && (string(me.CloudVendorType) != "azure") {
		return fmt.Errorf("'azure_connection' must not be specified unless 'cloud_vendor_type' is set to 'azure'; got 'cloud_vendor_type'='%v'", me.CloudVendorType)
	}
	if (me.AzureConnection == nil) && (string(me.CloudVendorType) == "azure") {
		return fmt.Errorf("'azure_connection' must be specified when 'cloud_vendor_type' is set to 'azure'; got 'cloud_vendor_type'='%v'", me.CloudVendorType)
	}
	if (me.GcpConnection != nil) && (string(me.CloudVendorType) != "gcp") {
		return fmt.Errorf("'gcp_connection' must not be specified unless 'cloud_vendor_type' is set to 'gcp'; got 'cloud_vendor_type'='%v'", me.CloudVendorType)
	}
	if (me.GcpConnection == nil) && (string(me.CloudVendorType) == "gcp") {
		return fmt.Errorf("'gcp_connection' must be specified when 'cloud_vendor_type' is set to 'gcp'; got 'cloud_vendor_type'='%v'", me.CloudVendorType)
	}
	// ---- BuiltinIngestSources []string -> {"preconditions":[{"expectedValue":"raw","property":"dataForwardingType","type":"EQUALS"}],"type":"AND"}
	// ---- BuiltinPipelines []string -> {"preconditions":[{"expectedValue":"processed","property":"dataForwardingType","type":"EQUALS"}],"type":"AND"}
	// ---- IngestSources []string -> {"preconditions":[{"expectedValue":"raw","property":"dataForwardingType","type":"EQUALS"}],"type":"AND"}
	// ---- Pipelines []string -> {"preconditions":[{"expectedValue":"processed","property":"dataForwardingType","type":"EQUALS"}],"type":"AND"}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"aws_connection":         &me.AwsConnection,
		"azure_connection":       &me.AzureConnection,
		"builtin_ingest_sources": &me.BuiltinIngestSources,
		"builtin_pipelines":      &me.BuiltinPipelines,
		"bulk_pattern":           &me.BulkPattern,
		"bulk_size":              &me.BulkSize,
		"cloud_vendor_type":      &me.CloudVendorType,
		"data_forwarding_type":   &me.DataForwardingType,
		"enabled":                &me.Enabled,
		"forwarding_name":        &me.ForwardingName,
		"gcp_connection":         &me.GcpConnection,
		"ingest_sources":         &me.IngestSources,
		"matcher":                &me.Matcher,
		"pipelines":              &me.Pipelines,
		"processing":             &me.Processing,
	})
}
