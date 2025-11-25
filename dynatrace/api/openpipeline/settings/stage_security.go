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

type SecurityContextStage struct {
	Editable   *bool                  `json:"editable"`
	Processors []*SecContextProcessor `json:"processors"`
}

func (f *SecurityContextStage) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"processor": {
			Type:        schema.TypeList,
			Description: "Groups all processors applicable for the SecurityContextStage.\nApplicable processor is SecurityContextProcessor.",
			Elem:        &schema.Resource{Schema: new(SecContextProcessor).Schema()},
			Optional:    true,
		},
	}
}

func (f *SecurityContextStage) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("processor", f.Processors)
}

func (f *SecurityContextStage) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("processor", &f.Processors)
}

type SecContextProcessor struct {
	securityContextProcessor *SecurityContextProcessor
}

func (ep *SecContextProcessor) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"security_context_processor": {
			Type:        schema.TypeList,
			Description: "Processor to set the security context field",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(SecurityContextProcessor).Schema()},
			Optional:    true,
		},
	}
}

func (ep *SecContextProcessor) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"security_context_processor": ep.securityContextProcessor,
	})
}

func (ep *SecContextProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"security_context_processor": &ep.securityContextProcessor,
	})
}

func (ep SecContextProcessor) MarshalJSON() ([]byte, error) {
	if ep.securityContextProcessor != nil {
		return json.Marshal(ep.securityContextProcessor)
	}

	return nil, errors.New("missing SecurityContextProcessor value")
}

func (ep *SecContextProcessor) UnmarshalJSON(b []byte) error {
	ttype, err := ExtractType(b)
	if err != nil {
		return err
	}

	switch ttype {
	case SecurityContextProcessorType:
		securityContextProcessor := SecurityContextProcessor{}
		if err := json.Unmarshal(b, &securityContextProcessor); err != nil {
			return err
		}
		ep.securityContextProcessor = &securityContextProcessor

	default:
		return fmt.Errorf("unknown SecurityContextProcessor type: %s", ttype)
	}

	return nil
}
