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

package filters

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Filters TODO: documentation
type Filters struct {
	CICSMQQueueIdIncludes []string `json:"cicsMqQueueIdIncludes,omitempty"`
	CICSMQQueueIdExcludes []string `json:"cicsMqQueueIdExcludes,omitempty"`
	IMSMQQueueIdIncludes  []string `json:"imsMqQueueIdIncludes,omitempty"`
	IMSMQQueueIdExcludes  []string `json:"imsMqQueueIdExcludes,omitempty"`
	IMSCrTrnIdIncludes    []string `json:"imsCrTrnIdIncludes,omitempty"`
	IMSCrTrnIdExcludes    []string `json:"imsCrTrnIdExcludes,omitempty"`
}

func (me *Filters) Name() string {
	return "ibm_mq_filters"
}

func (me *Filters) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cics_mq_queue_id_includes": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "CICS: Included MQ queues",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"cics_mq_queue_id_excludes": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "CICS: Excluded MQ queues",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"ims_mq_queue_id_includes": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "IMS: Included MQ queues",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"ims_mq_queue_id_excludes": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "IMS: Excluded MQ queues",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"ims_cr_trn_id_includes": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "IMS bridge: Included transaction IDs",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"ims_cr_trn_id_excludes": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "IMS bridge: Excluded transaction IDs",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *Filters) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"cics_mq_queue_id_includes": me.CICSMQQueueIdIncludes,
		"cics_mq_queue_id_excludes": me.CICSMQQueueIdExcludes,
		"ims_mq_queue_id_includes":  me.IMSMQQueueIdIncludes,
		"ims_mq_queue_id_excludes":  me.IMSMQQueueIdExcludes,
		"ims_cr_trn_id_includes":    me.IMSCrTrnIdIncludes,
		"ims_cr_trn_id_excludes":    me.IMSCrTrnIdExcludes,
	})
}

func (me *Filters) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"cics_mq_queue_id_includes": &me.CICSMQQueueIdIncludes,
		"cics_mq_queue_id_excludes": &me.CICSMQQueueIdExcludes,
		"ims_mq_queue_id_includes":  &me.IMSMQQueueIdIncludes,
		"ims_mq_queue_id_excludes":  &me.IMSMQQueueIdExcludes,
		"ims_cr_trn_id_includes":    &me.IMSCrTrnIdIncludes,
		"ims_cr_trn_id_excludes":    &me.IMSCrTrnIdExcludes,
	})
}
