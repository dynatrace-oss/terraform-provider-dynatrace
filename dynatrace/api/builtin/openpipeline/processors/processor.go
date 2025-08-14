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

package processors

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	DqlProcessorType          = "dql"
	FieldsAddProcessorType    = "fieldsAdd"
	FieldsRemoveProcessorType = "fieldsRemove"
	FieldsRenameProcessorType = "fieldsRename"
	DropProcessorType         = "drop"

	CounterMetricProcessorType = "counterMetric"
	ValueMetricProcessorType   = "valueMetric"

	DavisEventExtractionProcessorType = "davis"
	BizEventExtractionProcessorType   = "bizevent"

	SecurityContextProcessorType = "securityContext"

	NoStorageStageProcessorType        = "noStorage"
	BucketAssignmentStageProcessorType = "bucketAssignment"

	TechnologyProcessorType = "technology"

	AzureLogForwardingProcessorType = "azureLogForwarding"
	SecurityEventProcessorType      = "securityEvent"
	CostAllocationProcessorType     = "costAllocation"
	ProductAllocationProcessorType  = "productAllocation"
)

type Processor struct {
	Enabled      bool                    `json:"enabled"`
	Id           string                  `json:"id"`
	Type         string                  `json:"type"`
	Description  string                  `json:"description"`
	SampleData   *string                 `json:"sampleData,omitempty"`
	Matcher      *string                 `json:"matcher,omitempty"`
	Dql          *DqlAttributes          `json:"dql,omitempty"`
	FieldsAdd    *FieldsAddAttributes    `json:"fieldsAdd,omitempty"`
	FieldsRename *FieldsRenameAttributes `json:"fieldsRename,omitempty"`
	FieldsRemove *FieldsRemoveAttributes `json:"fieldsRemove,omitempty"`
}

func (p *Processor) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Indicates if the object is active",
			Default:     true,
		},
		"id": {
			Type:        schema.TypeString,
			Description: "Identifier of the processor. Must be unique within a stage. Must not start with 'dt.' or 'dynatrace.'",
			Required:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Type of the processor.",
			Required:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "Name or description of the processor",
			Required:    true,
		},
		"sample_data": {
			Type:        schema.TypeString,
			Description: "Sample data related to the processor for documentation or testing",
			Optional:    true,
		},
		"matcher": {
			Type:        schema.TypeString,
			Description: "Matching condition to apply on incoming records",
			Optional:    true,
		},
		"dql": {
			Type:        schema.TypeList,
			MaxItems:    1,
			MinItems:    0,
			Description: "Properties of the DQL processor",
			Elem:        &schema.Resource{Schema: new(DqlAttributes).Schema()},
		},
		"fieldsAdd": {
			Type:        schema.TypeList,
			MaxItems:    1,
			MinItems:    0,
			Description: "Properties of the fieldsAdd processor",
			Elem:        &schema.Resource{Schema: new(FieldsAddAttributes).Schema()},
		},
		"fieldsRename": {
			Type:        schema.TypeList,
			MaxItems:    1,
			MinItems:    0,
			Description: "Properties of the fieldsRename processor",
			Elem:        &schema.Resource{Schema: new(FieldsRenameAttributes).Schema()},
		},
		"fieldsRemove": {
			Type:        schema.TypeList,
			MaxItems:    1,
			MinItems:    0,
			Description: "Properties of the fieldsRemove processor",
			Elem:        &schema.Resource{Schema: new(FieldsRemoveAttributes).Schema()},
		},
	}
}

func (p *Processor) MarshalHCL(properties hcl.Properties) error {
	err := properties.EncodeAll(map[string]any{
		"enabled":     p.Enabled,
		"id":          p.Id,
		"type":        p.Type,
		"description": p.Description,
	})

	if err != nil {
		return err
	}

	if p.SampleData != nil {
		err = properties.Encode("sample_data", p.SampleData)
		if err != nil {
			return err
		}
	}
	if p.Matcher != nil {
		err = properties.Encode("matcher", p.Matcher)
		if err != nil {
			return err
		}
	}
	if p.Dql != nil {
		err = properties.Encode("dql", p.Dql)
		if err != nil {
			return err
		}
	}
	if p.FieldsAdd != nil {
		err = properties.Encode("fieldsAdd", p.FieldsAdd)
		if err != nil {
			return err
		}
	}
	if p.FieldsRename != nil {
		err = properties.Encode("fieldsRename", p.FieldsRename)
		if err != nil {
			return err
		}
	}
	if p.FieldsRemove != nil {
		err = properties.Encode("fieldsRemove", p.FieldsRemove)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Processor) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":      &p.Enabled,
		"id":           &p.Id,
		"type":         &p.Type,
		"description":  &p.Description,
		"sample_data":  &p.SampleData,
		"matcher":      &p.Matcher,
		"dql":          &p.Dql,
		"fieldsAdd":    &p.FieldsAdd,
		"fieldsRename": &p.FieldsRename,
		"fieldsRemove": &p.FieldsRemove,
	})
}
