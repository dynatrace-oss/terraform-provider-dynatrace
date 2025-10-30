/**
* @license
* Copyright 2025 Dynatrace LLC
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

package pipelines

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	CostAllocation    *Stage `json:"costAllocation"`    // Cost allocation stage
	CustomID          string `json:"customId"`          // Custom pipeline id
	DataExtraction    *Stage `json:"dataExtraction"`    // Data extraction stage
	Davis             *Stage `json:"davis"`             // Davis event extraction stage
	DisplayName       string `json:"displayName"`       // Display name
	MetricExtraction  *Stage `json:"metricExtraction"`  // Metrics extraction stage
	Processing        *Stage `json:"processing"`        // Processing stage
	ProductAllocation *Stage `json:"productAllocation"` // Product allocation stage
	SecurityContext   *Stage `json:"securityContext"`   // Security context stage
	Storage           *Stage `json:"storage"`           // Storage stage
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cost_allocation": {
			Type:        schema.TypeList,
			Description: "Cost allocation stage",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Stage).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"custom_id": {
			Type:        schema.TypeString,
			Description: "Custom pipeline id",
			Required:    true,
		},
		"data_extraction": {
			Type:        schema.TypeList,
			Description: "Data extraction stage",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Stage).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"davis": {
			Type:        schema.TypeList,
			Description: "Davis event extraction stage",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Stage).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"display_name": {
			Type:        schema.TypeString,
			Description: "Display name",
			Required:    true,
		},
		"metric_extraction": {
			Type:        schema.TypeList,
			Description: "Metrics extraction stage",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Stage).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"processing": {
			Type:        schema.TypeList,
			Description: "Processing stage",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Stage).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"product_allocation": {
			Type:        schema.TypeList,
			Description: "Product allocation stage",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Stage).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"security_context": {
			Type:        schema.TypeList,
			Description: "Security context stage",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Stage).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"storage": {
			Type:        schema.TypeList,
			Description: "Storage stage",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Stage).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"cost_allocation":    me.CostAllocation,
		"custom_id":          me.CustomID,
		"data_extraction":    me.DataExtraction,
		"davis":              me.Davis,
		"display_name":       me.DisplayName,
		"metric_extraction":  me.MetricExtraction,
		"processing":         me.Processing,
		"product_allocation": me.ProductAllocation,
		"security_context":   me.SecurityContext,
		"storage":            me.Storage,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"cost_allocation":    &me.CostAllocation,
		"custom_id":          &me.CustomID,
		"data_extraction":    &me.DataExtraction,
		"davis":              &me.Davis,
		"display_name":       &me.DisplayName,
		"metric_extraction":  &me.MetricExtraction,
		"processing":         &me.Processing,
		"product_allocation": &me.ProductAllocation,
		"security_context":   &me.SecurityContext,
		"storage":            &me.Storage,
	})
}
