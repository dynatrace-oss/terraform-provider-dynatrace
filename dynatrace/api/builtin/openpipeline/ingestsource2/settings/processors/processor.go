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
	"encoding/json"
	"errors"
	"fmt"

	openpipeline "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/openpipeline/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//const (
//	DqlProcessorType          = "dql"
//	FieldsAddProcessorType    = "fieldsAdd"
//	FieldsRemoveProcessorType = "fieldsRemove"
//	FieldsRenameProcessorType = "fieldsRename"
//	DropProcessorType         = "drop"
//
//	CounterMetricProcessorType = "counterMetric"
//	ValueMetricProcessorType   = "valueMetric"
//
//	DavisEventExtractionProcessorType = "davis"
//	BizEventExtractionProcessorType   = "bizevent"
//
//	SecurityContextProcessorType = "securityContext"
//
//	NoStorageStageProcessorType        = "noStorage"
//	BucketAssignmentStageProcessorType = "bucketAssignment"
//
//	TechnologyProcessorType = "technology"
//
//	AzureLogForwardingProcessorType = "azureLogForwarding"
//	SecurityEventProcessorType      = "securityEvent"
//	CostAllocationProcessorType     = "costAllocation"
//	ProductAllocationProcessorType  = "productAllocation"
//)

//var AvailableProcessorTypes = []string{DqlProcessorType, FieldsAddProcessorType, FieldsRemoveProcessorType,
//	FieldsRenameProcessorType, TechnologyProcessorType, DropProcessorType, BucketAssignmentStageProcessorType,
//	NoStorageStageProcessorType, SecurityContextProcessorType, CounterMetricProcessorType, ValueMetricProcessorType,
//	DavisEventExtractionProcessorType, BizEventExtractionProcessorType, AzureLogForwardingProcessorType,
//	SecurityEventProcessorType, CostAllocationProcessorType, ProductAllocationProcessorType}

type Processor struct {
	DqlProcessor          *DqlProcessor
	FieldsAddProcessor    *FieldsAddProcessor
	FieldsRenameProcessor *FieldsRenameProcessor
	FieldsRemoveProcessor *FieldsRemoveProcessor
	DropProcessor         *DropProcessor
}

func (p *Processor) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"dql_processor": {
			Type:        schema.TypeList,
			MaxItems:    1,
			MinItems:    0,
			Elem:        &schema.Resource{Schema: new(DqlProcessor).Schema()},
			Description: "DQL processor",
			Optional:    true,
		},
		"fields_add_processor": {
			Type:        schema.TypeList,
			MaxItems:    1,
			MinItems:    0,
			Elem:        &schema.Resource{Schema: new(FieldsAddProcessor).Schema()},
			Description: "fields_add processor",
		},
		"fields_rename_processor": {
			Type:        schema.TypeList,
			MaxItems:    1,
			MinItems:    0,
			Elem:        &schema.Resource{Schema: new(FieldsRenameProcessor).Schema()},
			Description: "fields_rename processor",
		},
		"fields_remove_processor": {
			Type:        schema.TypeList,
			MaxItems:    1,
			MinItems:    0,
			Elem:        &schema.Resource{Schema: new(FieldsRemoveProcessor).Schema()},
			Description: "fields_remove processor",
		},
		"drop_processor": {
			Type:     schema.TypeList,
			MaxItems: 1,
			MinItems: 0,
			Elem:     &schema.Resource{Schema: new(DropProcessor).Schema()},
		},
	}
}

func (p *Processor) MarshalHCL(properties hcl.Properties) error {
	err := properties.Encode("dql_processor", p.DqlProcessor)
	if err != nil {
		return err
	}

	err = properties.Encode("fields_add_processor", p.FieldsAddProcessor)
	if err != nil {
		return err
	}

	err = properties.Encode("fields_remove_processor", p.FieldsRemoveProcessor)
	if err != nil {
		return err
	}

	err = properties.Encode("fields_rename_processor", p.FieldsRenameProcessor)
	if err != nil {
		return err
	}

	err = properties.Encode("drop_processor", p.DropProcessor)
	if err != nil {
		return err
	}

	return nil
}

func (p *Processor) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"dql_processor":           &p.DqlProcessor,
		"fields_add_processor":    &p.FieldsAddProcessor,
		"fields_remove_processor": &p.FieldsRemoveProcessor,
		"fields_rename_processor": &p.FieldsRenameProcessor,
		"drop_processor":          &p.FieldsRenameProcessor,
	})
}

func (ep Processor) MarshalJSON() ([]byte, error) {
	if ep.DqlProcessor != nil {
		return json.Marshal(ep.DqlProcessor)
	}
	if ep.FieldsAddProcessor != nil {
		return json.Marshal(ep.FieldsAddProcessor)
	}
	if ep.FieldsRemoveProcessor != nil {
		return json.Marshal(ep.FieldsRemoveProcessor)
	}
	if ep.FieldsRenameProcessor != nil {
		return json.Marshal(ep.FieldsRenameProcessor)
	}
	if ep.DropProcessor != nil {
		return json.Marshal(ep.DropProcessor)
	}

	return nil, errors.New("missing EndpointProcessor value")
}

func (ep *Processor) UnmarshalJSON(b []byte) error {
	ttype, err := openpipeline.ExtractType(b)
	if err != nil {
		return err
	}

	switch ttype {
	case DqlProcessorType:
		dqlProcessor := DqlProcessor{}
		if err := json.Unmarshal(b, &dqlProcessor); err != nil {
			return err
		}
		ep.DqlProcessor = &dqlProcessor

	case FieldsAddProcessorType:
		fieldsAddProcessor := FieldsAddProcessor{}
		if err := json.Unmarshal(b, &fieldsAddProcessor); err != nil {
			return err
		}
		ep.FieldsAddProcessor = &fieldsAddProcessor

	case FieldsRemoveProcessorType:
		fieldsRemoveProcessor := FieldsRemoveProcessor{}
		if err := json.Unmarshal(b, &fieldsRemoveProcessor); err != nil {
			return err
		}
		ep.FieldsRemoveProcessor = &fieldsRemoveProcessor

	case FieldsRenameProcessorType:
		fieldsRenameProcessor := FieldsRenameProcessor{}
		if err := json.Unmarshal(b, &fieldsRenameProcessor); err != nil {
			return err
		}
		ep.FieldsRenameProcessor = &fieldsRenameProcessor

	case DropProcessorType:
		dropProcessor := DropProcessor{}
		if err := json.Unmarshal(b, &dropProcessor); err != nil {
			return err
		}
		ep.DropProcessor = &dropProcessor

	default:
		return fmt.Errorf("unknown EndpointProcessor type: %s", ttype)
	}

	return nil
}
