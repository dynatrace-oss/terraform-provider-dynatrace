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

package fullwebrequest

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ContextIdContributor struct {
	EnableIdContributor  bool         `json:"enableIdContributor"` // Transform this value before letting it contribute to the Service Id
	ServiceIdContributor *ContextRoot `json:"serviceIdContributor,omitempty"`
}

func (me *ContextIdContributor) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enable_id_contributor": {
			Type:        schema.TypeBool,
			Description: "Transform this value before letting it contribute to the Service Id",
			Required:    true,
		},
		"service_id_contributor": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(ContextRoot).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *ContextIdContributor) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enable_id_contributor":  me.EnableIdContributor,
		"service_id_contributor": me.ServiceIdContributor,
	})
}

func (me *ContextIdContributor) HandlePreconditions() error {
	if me.ServiceIdContributor == nil && me.EnableIdContributor {
		return fmt.Errorf("'service_id_contributor' must be specified if 'enable_id_contributor' is set to '%v'", me.EnableIdContributor)
	}
	if me.ServiceIdContributor != nil && !me.EnableIdContributor {
		return fmt.Errorf("'service_id_contributor' must not be specified if 'enable_id_contributor' is set to '%v'", me.EnableIdContributor)
	}
	return nil
}

func (me *ContextIdContributor) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enable_id_contributor":  &me.EnableIdContributor,
		"service_id_contributor": &me.ServiceIdContributor,
	})
}
