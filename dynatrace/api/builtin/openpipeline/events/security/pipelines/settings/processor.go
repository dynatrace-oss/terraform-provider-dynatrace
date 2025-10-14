/**
* @license
* Copyright 2020 Dynatrace LLC
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
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Processors []*Processor

func (me *Processors) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"processor": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(Processor).Schema()},
		},
	}
}

func (me Processors) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("processor", me)
}

func (me *Processors) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("processor", me)
}

// Processor. Processor definition
type Processor struct {
	AzureLogForwarding         *AzureLogForwardingAttributes         `json:"azureLogForwarding,omitempty"` // Azure log forwarding processor attributes
	Bizevent                   *BizeventAttributes                   `json:"bizevent,omitempty"`           // Bizevent extraction processor attributes
	BucketAssignment           *BucketAssignmentAttributes           `json:"bucketAssignment,omitempty"`   // Bucket assignment processor attributes
	CostAllocation             *CostAllocationAttributes             `json:"costAllocation,omitempty"`     // Cost allocation processor attributes
	CounterMetric              *CounterMetricAttributes              `json:"counterMetric,omitempty"`      // Counter metric processor attributes
	Davis                      *DavisAttributes                      `json:"davis,omitempty"`              // Davis event extraction processor attributes
	Description                string                                `json:"description"`
	Dql                        *DqlAttributes                        `json:"dql,omitempty"`                        // DQL processor attributes
	Enabled                    bool                                  `json:"enabled"`                              // This setting is enabled (`true`) or disabled (`false`)
	FieldsAdd                  *FieldsAddAttributes                  `json:"fieldsAdd,omitempty"`                  // Fields add processor attributes
	FieldsRemove               *FieldsRemoveAttributes               `json:"fieldsRemove,omitempty"`               // Fields remove processor attributes
	FieldsRename               *FieldsRenameAttributes               `json:"fieldsRename,omitempty"`               // Fields rename processor attributes
	HistogramMetric            *HistogramMetricAttributes            `json:"histogramMetric,omitempty"`            // Histogram metric processor attributes
	ID                         string                                `json:"id"`                                   // Processor identifier
	Matcher                    *string                               `json:"matcher,omitempty"`                    // [See our documentation](https://dt-url.net/bp234rv)
	ProductAllocation          *ProductAllocationAttributes          `json:"productAllocation,omitempty"`          // Product allocation processor attributes
	SampleData                 *string                               `json:"sampleData,omitempty"`                 // Sample data
	SamplingAwareCounterMetric *SamplingAwareCounterMetricAttributes `json:"samplingAwareCounterMetric,omitempty"` // Sampling-aware counter metric processor attributes
	SamplingAwareValueMetric   *SamplingAwareValueMetricAttributes   `json:"samplingAwareValueMetric,omitempty"`   // Sampling aware value metric processor attributes
	SecurityContext            *SecurityContextAttributes            `json:"securityContext,omitempty"`            // Security context processor attributes
	SecurityEvent              *SecurityEventAttributes              `json:"securityEvent,omitempty"`              // Security event extraction processor attributes
	Technology                 *TechnologyAttributes                 `json:"technology,omitempty"`                 // Technology processor attributes
	Type                       ProcessorType                         `json:"type"`                                 // Processor type. Possible Values: `azureLogForwarding`, `bizevent`, `bucketAssignment`, `costAllocation`, `counterMetric`, `davis`, `dql`, `drop`, `fieldsAdd`, `fieldsRemove`, `fieldsRename`, `histogramMetric`, `noStorage`, `productAllocation`, `samplingAwareCounterMetric`, `samplingAwareValueMetric`, `securityContext`, `securityEvent`, `technology`, `valueMetric`.
	ValueMetric                *ValueMetricAttributes                `json:"valueMetric,omitempty"`                // Value metric processor attributes
}

func (me *Processor) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"azure_log_forwarding": {
			Type:        schema.TypeList,
			Description: "Azure log forwarding processor attributes",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(AzureLogForwardingAttributes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"bizevent": {
			Type:        schema.TypeList,
			Description: "Bizevent extraction processor attributes",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(BizeventAttributes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"bucket_assignment": {
			Type:        schema.TypeList,
			Description: "Bucket assignment processor attributes",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(BucketAssignmentAttributes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"cost_allocation": {
			Type:        schema.TypeList,
			Description: "Cost allocation processor attributes",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(CostAllocationAttributes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"counter_metric": {
			Type:        schema.TypeList,
			Description: "Counter metric processor attributes",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(CounterMetricAttributes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"davis": {
			Type:        schema.TypeList,
			Description: "Davis event extraction processor attributes",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(DavisAttributes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Required:    true,
		},
		"dql": {
			Type:        schema.TypeList,
			Description: "DQL processor attributes",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(DqlAttributes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"fields_add": {
			Type:        schema.TypeList,
			Description: "Fields add processor attributes",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(FieldsAddAttributes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"fields_remove": {
			Type:        schema.TypeList,
			Description: "Fields remove processor attributes",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(FieldsRemoveAttributes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"fields_rename": {
			Type:        schema.TypeList,
			Description: "Fields rename processor attributes",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(FieldsRenameAttributes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"histogram_metric": {
			Type:        schema.TypeList,
			Description: "Histogram metric processor attributes",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(HistogramMetricAttributes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"id": {
			Type:        schema.TypeString,
			Description: "Processor identifier",
			Required:    true,
		},
		"matcher": {
			Type:        schema.TypeString,
			Description: "[See our documentation](https://dt-url.net/bp234rv)",
			Optional:    true, // precondition
		},
		"product_allocation": {
			Type:        schema.TypeList,
			Description: "Product allocation processor attributes",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(ProductAllocationAttributes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"sample_data": {
			Type:        schema.TypeString,
			Description: "Sample data",
			Optional:    true, // nullable
		},
		"sampling_aware_counter_metric": {
			Type:        schema.TypeList,
			Description: "Sampling-aware counter metric processor attributes",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(SamplingAwareCounterMetricAttributes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"sampling_aware_value_metric": {
			Type:        schema.TypeList,
			Description: "Sampling aware value metric processor attributes",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(SamplingAwareValueMetricAttributes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"security_context": {
			Type:        schema.TypeList,
			Description: "Security context processor attributes",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(SecurityContextAttributes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"security_event": {
			Type:        schema.TypeList,
			Description: "Security event extraction processor attributes",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(SecurityEventAttributes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"technology": {
			Type:        schema.TypeList,
			Description: "Technology processor attributes",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(TechnologyAttributes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Processor type. Possible Values: `azureLogForwarding`, `bizevent`, `bucketAssignment`, `costAllocation`, `counterMetric`, `davis`, `dql`, `drop`, `fieldsAdd`, `fieldsRemove`, `fieldsRename`, `histogramMetric`, `noStorage`, `productAllocation`, `samplingAwareCounterMetric`, `samplingAwareValueMetric`, `securityContext`, `securityEvent`, `technology`, `valueMetric`.",
			Required:    true,
		},
		"value_metric": {
			Type:        schema.TypeList,
			Description: "Value metric processor attributes",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(ValueMetricAttributes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Processor) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"azure_log_forwarding":          me.AzureLogForwarding,
		"bizevent":                      me.Bizevent,
		"bucket_assignment":             me.BucketAssignment,
		"cost_allocation":               me.CostAllocation,
		"counter_metric":                me.CounterMetric,
		"davis":                         me.Davis,
		"description":                   me.Description,
		"dql":                           me.Dql,
		"enabled":                       me.Enabled,
		"fields_add":                    me.FieldsAdd,
		"fields_remove":                 me.FieldsRemove,
		"fields_rename":                 me.FieldsRename,
		"histogram_metric":              me.HistogramMetric,
		"id":                            me.ID,
		"matcher":                       me.Matcher,
		"product_allocation":            me.ProductAllocation,
		"sample_data":                   me.SampleData,
		"sampling_aware_counter_metric": me.SamplingAwareCounterMetric,
		"sampling_aware_value_metric":   me.SamplingAwareValueMetric,
		"security_context":              me.SecurityContext,
		"security_event":                me.SecurityEvent,
		"technology":                    me.Technology,
		"type":                          me.Type,
		"value_metric":                  me.ValueMetric,
	})
}

func (me *Processor) HandlePreconditions() error {
	if (me.Matcher == nil) && (string(me.Type) != "technology") {
		me.Matcher = opt.NewString("")
	}
	if (me.AzureLogForwarding == nil) && (string(me.Type) == "azureLogForwarding") {
		return fmt.Errorf("'azure_log_forwarding' must be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.AzureLogForwarding != nil) && (string(me.Type) != "azureLogForwarding") {
		return fmt.Errorf("'azure_log_forwarding' must not be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.Bizevent == nil) && (string(me.Type) == "bizevent") {
		return fmt.Errorf("'bizevent' must be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.Bizevent != nil) && (string(me.Type) != "bizevent") {
		return fmt.Errorf("'bizevent' must not be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.BucketAssignment == nil) && (string(me.Type) == "bucketAssignment") {
		return fmt.Errorf("'bucket_assignment' must be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.BucketAssignment != nil) && (string(me.Type) != "bucketAssignment") {
		return fmt.Errorf("'bucket_assignment' must not be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.CostAllocation == nil) && (string(me.Type) == "costAllocation") {
		return fmt.Errorf("'cost_allocation' must be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.CostAllocation != nil) && (string(me.Type) != "costAllocation") {
		return fmt.Errorf("'cost_allocation' must not be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.CounterMetric == nil) && (string(me.Type) == "counterMetric") {
		return fmt.Errorf("'counter_metric' must be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.CounterMetric != nil) && (string(me.Type) != "counterMetric") {
		return fmt.Errorf("'counter_metric' must not be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.Davis == nil) && (string(me.Type) == "davis") {
		return fmt.Errorf("'davis' must be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.Davis != nil) && (string(me.Type) != "davis") {
		return fmt.Errorf("'davis' must not be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.Dql == nil) && (string(me.Type) == "dql") {
		return fmt.Errorf("'dql' must be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.Dql != nil) && (string(me.Type) != "dql") {
		return fmt.Errorf("'dql' must not be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.FieldsAdd == nil) && (string(me.Type) == "fieldsAdd") {
		return fmt.Errorf("'fields_add' must be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.FieldsAdd != nil) && (string(me.Type) != "fieldsAdd") {
		return fmt.Errorf("'fields_add' must not be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.FieldsRemove == nil) && (string(me.Type) == "fieldsRemove") {
		return fmt.Errorf("'fields_remove' must be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.FieldsRemove != nil) && (string(me.Type) != "fieldsRemove") {
		return fmt.Errorf("'fields_remove' must not be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.FieldsRename == nil) && (string(me.Type) == "fieldsRename") {
		return fmt.Errorf("'fields_rename' must be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.FieldsRename != nil) && (string(me.Type) != "fieldsRename") {
		return fmt.Errorf("'fields_rename' must not be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.HistogramMetric == nil) && (string(me.Type) == "histogramMetric") {
		return fmt.Errorf("'histogram_metric' must be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.HistogramMetric != nil) && (string(me.Type) != "histogramMetric") {
		return fmt.Errorf("'histogram_metric' must not be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.ProductAllocation == nil) && (string(me.Type) == "productAllocation") {
		return fmt.Errorf("'product_allocation' must be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.ProductAllocation != nil) && (string(me.Type) != "productAllocation") {
		return fmt.Errorf("'product_allocation' must not be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.SamplingAwareCounterMetric == nil) && (string(me.Type) == "samplingAwareCounterMetric") {
		return fmt.Errorf("'sampling_aware_counter_metric' must be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.SamplingAwareCounterMetric != nil) && (string(me.Type) != "samplingAwareCounterMetric") {
		return fmt.Errorf("'sampling_aware_counter_metric' must not be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.SamplingAwareValueMetric == nil) && (string(me.Type) == "samplingAwareValueMetric") {
		return fmt.Errorf("'sampling_aware_value_metric' must be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.SamplingAwareValueMetric != nil) && (string(me.Type) != "samplingAwareValueMetric") {
		return fmt.Errorf("'sampling_aware_value_metric' must not be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.SecurityContext == nil) && (string(me.Type) == "securityContext") {
		return fmt.Errorf("'security_context' must be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.SecurityContext != nil) && (string(me.Type) != "securityContext") {
		return fmt.Errorf("'security_context' must not be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.SecurityEvent == nil) && (string(me.Type) == "securityEvent") {
		return fmt.Errorf("'security_event' must be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.SecurityEvent != nil) && (string(me.Type) != "securityEvent") {
		return fmt.Errorf("'security_event' must not be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.Technology == nil) && (string(me.Type) == "technology") {
		return fmt.Errorf("'technology' must be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.Technology != nil) && (string(me.Type) != "technology") {
		return fmt.Errorf("'technology' must not be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.ValueMetric == nil) && (string(me.Type) == "valueMetric") {
		return fmt.Errorf("'value_metric' must be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.ValueMetric != nil) && (string(me.Type) != "valueMetric") {
		return fmt.Errorf("'value_metric' must not be specified if 'type' is set to '%v'", me.Type)
	}
	return nil
}

func (me *Processor) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"azure_log_forwarding":          &me.AzureLogForwarding,
		"bizevent":                      &me.Bizevent,
		"bucket_assignment":             &me.BucketAssignment,
		"cost_allocation":               &me.CostAllocation,
		"counter_metric":                &me.CounterMetric,
		"davis":                         &me.Davis,
		"description":                   &me.Description,
		"dql":                           &me.Dql,
		"enabled":                       &me.Enabled,
		"fields_add":                    &me.FieldsAdd,
		"fields_remove":                 &me.FieldsRemove,
		"fields_rename":                 &me.FieldsRename,
		"histogram_metric":              &me.HistogramMetric,
		"id":                            &me.ID,
		"matcher":                       &me.Matcher,
		"product_allocation":            &me.ProductAllocation,
		"sample_data":                   &me.SampleData,
		"sampling_aware_counter_metric": &me.SamplingAwareCounterMetric,
		"sampling_aware_value_metric":   &me.SamplingAwareValueMetric,
		"security_context":              &me.SecurityContext,
		"security_event":                &me.SecurityEvent,
		"technology":                    &me.Technology,
		"type":                          &me.Type,
		"value_metric":                  &me.ValueMetric,
	})
}
