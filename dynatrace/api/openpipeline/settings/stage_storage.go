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

package openpipeline

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type StorageStage struct {
	Editable           *bool                    `json:"editable,omitempty"`
	CatchAllBucketName string                   `json:"catchAllBucketName"`
	Processors         []*StorageStageProcessor `json:"processors"`
}

func (f *StorageStage) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"catch_all_bucket_name": {
			Type:        schema.TypeString,
			Description: "Default bucket assigned to records which do not match any other storage processor",
			Optional:    true,
		},
		"processor": {
			Type:        schema.TypeList,
			Description: "Groups all processors applicable for the StorageStage.\nApplicable processors are BucketAssignmentProcessor and NoStorageProcessor.",
			Elem:        &schema.Resource{Schema: new(StorageStageProcessor).Schema()},
			Optional:    true,
		},
	}
}

func (f *StorageStage) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("catch_all_bucket_name", f.CatchAllBucketName); err != nil {
		return err
	}

	return properties.EncodeSlice("processor", f.Processors)
}

func (f *StorageStage) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("catch_all_bucket_name", &f.CatchAllBucketName); err != nil {
		return err
	}

	return decoder.DecodeSlice("processor", &f.Processors)
}

type StorageStageProcessor struct {
	bucketAssignmentProcessor *BucketAssignmentProcessor
	noStorageProcessor        *NoStorageProcessor
}

func (ep *StorageStageProcessor) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"bucket_assignment_processor": {
			Type:        schema.TypeList,
			Description: "Processor to assign a bucket",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(BucketAssignmentProcessor).Schema()},
			Optional:    true,
		},
		"no_storage_processor": {
			Type:        schema.TypeList,
			Description: "Processor to skip storage assignment",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(NoStorageProcessor).Schema()},
			Optional:    true,
		},
	}
}

func (ep *StorageStageProcessor) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"bucket_assignment_processor": ep.bucketAssignmentProcessor,
		"no_storage_processor":        ep.noStorageProcessor,
	})
}

func (ep *StorageStageProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"bucket_assignment_processor": &ep.bucketAssignmentProcessor,
		"no_storage_processor":        &ep.noStorageProcessor,
	})
}

func (ep StorageStageProcessor) MarshalJSON() ([]byte, error) {
	if ep.bucketAssignmentProcessor != nil {
		return json.Marshal(ep.bucketAssignmentProcessor)
	}
	if ep.noStorageProcessor != nil {
		return json.Marshal(ep.noStorageProcessor)
	}

	return nil, errors.New("missing StorageStageProcessor value")
}

func (ep *StorageStageProcessor) UnmarshalJSON(b []byte) error {
	ttype, err := ExtractType(b)
	if err != nil {
		return err
	}

	switch ttype {
	case NoStorageStageProcessorType:
		noStorageProcessor := NoStorageProcessor{}
		if err := json.Unmarshal(b, &noStorageProcessor); err != nil {
			return err
		}
		ep.noStorageProcessor = &noStorageProcessor

	case BucketAssignmentStageProcessorType:
		bucketAssignmentProcessor := BucketAssignmentProcessor{}
		if err := json.Unmarshal(b, &bucketAssignmentProcessor); err != nil {
			return err
		}
		ep.bucketAssignmentProcessor = &bucketAssignmentProcessor

	default:
		return fmt.Errorf("unknown StorageStageProcessor type: %s", ttype)
	}

	return nil
}
