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
	"fmt"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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

var AvailableProcessorTypes = []string{DqlProcessorType, FieldsAddProcessorType, FieldsRemoveProcessorType,
	FieldsRenameProcessorType, TechnologyProcessorType, DropProcessorType, BucketAssignmentStageProcessorType,
	NoStorageStageProcessorType, SecurityContextProcessorType, CounterMetricProcessorType, ValueMetricProcessorType,
	DavisEventExtractionProcessorType, BizEventExtractionProcessorType, AzureLogForwardingProcessorType,
	SecurityEventProcessorType, CostAllocationProcessorType, ProductAllocationProcessorType}

type ProcessorBasic struct {
	Enabled     bool    `json:"enabled"`
	Id          string  `json:"id"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	SampleData  *string `json:"sampleData,omitempty"`
}

const IDMinLength = 4
const IDMaxLength = 100
const DescriptionMaxLength = 512
const SampleDataMaxLength = 8192
const MatcherMaxLength = 1500

func (p *ProcessorBasic) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Indicates if the object is active",
			Default:     true,
			Optional:    true,
		},
		"id": {
			Type:        schema.TypeString,
			Description: "Identifier of the processor. Must be unique within a stage. Must not start with 'dt.' or 'dynatrace.'",
			Required:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(IDMinLength, IDMaxLength),
				func(input interface{}, schema string) (warnings []string, errors []error) {
					id, ok := input.(string)
					if !ok {
						errors = append(errors, fmt.Errorf("expected type of %s to be string", schema))
						return warnings, errors
					}

					if strings.HasPrefix(id, "dt.") || strings.HasPrefix(id, "dynatrace.") {
						errors = append(errors,
							fmt.Errorf("%s must not start with 'dt.' or 'dynatrace.'", schema))
					}
					return warnings, errors
				}),
		},
		"description": {
			Type:         schema.TypeString,
			Description:  "Name or description of the processor",
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, DescriptionMaxLength),
		},
		"sample_data": {
			Type:         schema.TypeString,
			Description:  "Sample data related to the processor for documentation or testing",
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, SampleDataMaxLength),
		},
	}
}

func (p *ProcessorBasic) MarshalHCL(properties hcl.Properties) error {
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
	return nil
}

func (p *ProcessorBasic) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":     p.Enabled,
		"id":          p.Id,
		"type":        p.Type,
		"description": p.Description,
		"sample_data": p.SampleData,
	})
}
