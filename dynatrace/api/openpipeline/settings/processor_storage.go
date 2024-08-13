package openpipeline

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type StorageStageProcessors struct {
	Processors []StorageStageProcessor
}

func (ep *StorageStageProcessors) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"processor": {
			Type:        schema.TypeSet,
			Description: "todo",
			Elem:        &schema.Resource{Schema: new(StorageStageProcessor).Schema()},
			Optional:    true,
		},
	}
}

func (ep *StorageStageProcessors) MarshalHCL(properties hcl.Properties) error {
	return properties.Encode("processors", ep.Processors)
}

func (ep *StorageStageProcessors) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.Decode("processors", &ep.Processors)
}

type StorageStageProcessor struct {
	bucketAssignmentProcessor *BucketAssignmentProcessor
	noStorageProcessor        *NoStorageProcessor
}

func (ep *StorageStageProcessor) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"bucket_assignment_processor": {
			Type:        schema.TypeList,
			Description: "Processor to assign a bucket.",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(BucketAssignmentProcessor).Schema()},
			Optional:    true,
		},
		"no_storage_processor": {
			Type:        schema.TypeList,
			Description: "Processor to skip storage assignment.",
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
		"bucket_assignment_processor": ep.bucketAssignmentProcessor,
		"no_storage_processor":        ep.noStorageProcessor,
	})
}
