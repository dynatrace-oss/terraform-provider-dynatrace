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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/processors"
	openpipeline "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/openpipeline/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DqlProcessor struct {
	ProcessorBasic
	Dql processors.DqlAttributes `json:"dql,omitempty"`
}

func (da *DqlProcessor) Schema() map[string]*schema.Schema {
	scm := new(ProcessorBasic).Schema()
	scm["dql"] = &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		MinItems: 0,
		Elem:     &schema.Resource{Schema: new(processors.DqlAttributes).Schema()},
	}
	return scm
}

func (da *DqlProcessor) MarshalHCL(properties hcl.Properties) error {
	err := da.ProcessorBasic.MarshalHCL(properties)
	if err != nil {
		return err
	}
	return properties.Encode("dql", da.Dql)
}

func (da *DqlProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	err := da.ProcessorBasic.UnmarshalHCL(decoder)
	if err != nil {
		return err
	}
	return decoder.Decode("dql", &da.Dql)
}

func (ep DqlProcessor) MarshalJSON() ([]byte, error) {
	type dqlProcessor DqlProcessor
	return openpipeline.MarshalAsJSONWithType((dqlProcessor)(ep), DqlProcessorType)
}
