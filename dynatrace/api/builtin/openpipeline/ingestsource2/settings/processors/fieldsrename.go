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

type FieldsRenameProcessor struct {
	ProcessorBasic
	FieldsRename processors.FieldsRenameAttributes `json:"fieldsRename,omitempty"`
}

func (da *FieldsRenameProcessor) Schema() map[string]*schema.Schema {
	scm := new(ProcessorBasic).Schema()
	scm["fields_rename"] = &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem:     &schema.Resource{Schema: new(processors.FieldsRenameAttributes).Schema()},
	}
	return scm
}

func (da *FieldsRenameProcessor) MarshalHCL(properties hcl.Properties) error {
	err := da.ProcessorBasic.MarshalHCL(properties)
	if err != nil {
		return err
	}
	return properties.Encode("fields_rename", da.FieldsRename)
}

func (da *FieldsRenameProcessor) UnmarshalHCL(decoder hcl.Decoder) error {
	err := da.ProcessorBasic.UnmarshalHCL(decoder)
	if err != nil {
		return err
	}
	return decoder.Decode("fields_rename", &da.FieldsRename)
}

func (ep FieldsRenameProcessor) MarshalJSON() ([]byte, error) {
	type processor FieldsRenameProcessor
	return openpipeline.MarshalAsJSONWithType((processor)(ep), FieldsRenameProcessorType)
}
