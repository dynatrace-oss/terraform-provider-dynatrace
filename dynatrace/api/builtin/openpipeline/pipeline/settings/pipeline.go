package settings

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/processors"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var AllowedKinds = []string{
	"logs", "events", "events.security", "security.events", "bizevents", "spans",
	"events.sdlc", "metrics", "usersessions", "davis.problems", "davis.events",
	"system.events", "azure.logs.forwarding", "user.events",
}

type Pipeline struct {
	Kind              string `json:"-"`
	CustomID          string `json:"customId"`
	DisplayName       string `json:"displayName"`
	Processing        *Stage `json:"processing,omitempty"`
	SecurityContext   *Stage `json:"securityContext,omitempty"`
	CostAllocation    *Stage `json:"costAllocation,omitempty"`
	ProductAllocation *Stage `json:"productAllocation,omitempty"`
	Storage           *Stage `json:"storage,omitempty"`
	MetricExtraction  *Stage `json:"metricExtraction,omitempty"`
	Davis             *Stage `json:"davis,omitempty"`
	DataExtraction    *Stage `json:"dataExtraction,omitempty"`
}

const DisplayNameMaxLength = 500
const CustomIDMinLength = 4
const CustomIDMaxLength = 100

func (p *Pipeline) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"kind": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "Indicates OpenPipeline data source",
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(AllowedKinds, true),
		},
		"custom_id": {
			Type:        schema.TypeString,
			Description: "The ID of the pipeline. Must not start with 'dt.' or 'dynatrace.'",
			Required:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(CustomIDMinLength, CustomIDMaxLength),
				func(input interface{}, schema string) (warnings []string, errors []error) {
					id, ok := input.(string)

					// Schema has "NO_WHITESPACE" validation which is missing here.
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
		"display_name": {
			Type:         schema.TypeString,
			Description:  "Display name of the pipeline",
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, DisplayNameMaxLength),
		},
		"processing": {
			Type:        schema.TypeList,
			Description: "Processing stage",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Stage).Schema()},
			Optional:    true,
		},
		"security_context": {
			Type:        schema.TypeList,
			Description: "Security context stage",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Stage).Schema()},
			Optional:    true,
		},
		"cost_allocation": {
			Type:        schema.TypeList,
			Description: "Cost allocation stage",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Stage).Schema()},
			Optional:    true,
		},
		"product_allocation": {
			Type:        schema.TypeList,
			Description: "Product allocation stage",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Stage).Schema()},
			Optional:    true,
		},
		"storage": {
			Type:        schema.TypeList,
			Description: "Storage stage",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Stage).Schema()},
			Optional:    true,
		},
		"metric_extraction": {
			Type:        schema.TypeList,
			Description: "Metrics extraction stage",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Stage).Schema()},
			Optional:    true,
		},
		"davis": {
			Type:        schema.TypeList,
			Description: "Davis event extraction stage",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Stage).Schema()},
			Optional:    true,
		},
		"data_extraction": {
			Type:        schema.TypeList,
			Description: "Data extraction stage",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Stage).Schema()},
			Optional:    true,
		},
	}
}

func (p *Pipeline) MarshalHCL(properties hcl.Properties) error {
	err := properties.EncodeAll(map[string]any{
		"kind":               p.Kind,
		"custom_id":          p.CustomID,
		"display_name":       p.DisplayName,
		"processing":         p.Processing,
		"security_context":   p.SecurityContext,
		"cost_allocation":    p.CostAllocation,
		"product_allocation": p.ProductAllocation,
		"storage":            p.Storage,
		"metric_extraction":  p.MetricExtraction,
		"davis":              p.Davis,
		"data_extraction":    p.DataExtraction,
	})
	openpipeline.RemoveNils(properties)

	return err
}

func (p *Pipeline) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"kind":               &p.Kind,
		"custom_id":          &p.CustomID,
		"display_name":       &p.DisplayName,
		"processing":         &p.Processing,
		"security_context":   &p.SecurityContext,
		"cost_allocation":    &p.CostAllocation,
		"product_allocation": &p.ProductAllocation,
		"storage":            &p.Storage,
		"metric_extraction":  &p.MetricExtraction,
		"davis":              &p.Davis,
		"data_extraction":    &p.DataExtraction,
	})
}

func (p *Pipeline) MarshalJSON() ([]byte, error) {
	var temp = *p
	// The API expects the following fields to be non-nil
	if temp.Processing == nil {
		temp.Processing = &Stage{}
	}
	if temp.SecurityContext == nil {
		temp.SecurityContext = &Stage{}
	}
	if temp.CostAllocation == nil {
		temp.CostAllocation = &Stage{}
	}
	if temp.ProductAllocation == nil {
		temp.ProductAllocation = &Stage{}
	}
	if temp.Storage == nil {
		temp.Storage = &Stage{}
	}
	if temp.MetricExtraction == nil {
		temp.MetricExtraction = &Stage{}
	}
	if temp.Davis == nil {
		temp.Davis = &Stage{}
	}
	if temp.DataExtraction == nil {
		temp.DataExtraction = &Stage{}
	}

	return json.Marshal(temp)
}

type Stage struct {
	Processors []*processors.Processor `json:"processors,omitempty"`
}

func (s *Stage) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"processor": {
			Type:        schema.TypeList,
			Description: "Groups all processors of the stage",
			Elem:        &schema.Resource{Schema: new(processors.Processor).Schema()},
			Required:    true,
		},
	}
}

func (s *Stage) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("processor", s.Processors)
}

func (s *Stage) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("processor", &s.Processors)
}
