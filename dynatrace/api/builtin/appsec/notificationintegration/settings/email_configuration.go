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

package notificationintegration

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type EmailConfiguration struct {
	BccRecipients []string `json:"bccRecipients,omitempty"` // BCC
	CcRecipients  []string `json:"ccRecipients,omitempty"`  // CC
	Recipients    []string `json:"recipients"`              // To
}

func (me *EmailConfiguration) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"bcc_recipients": {
			Type:        schema.TypeSet,
			Description: "BCC",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"cc_recipients": {
			Type:        schema.TypeSet,
			Description: "CC",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"recipients": {
			Type:        schema.TypeSet,
			Description: "To",
			Required:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *EmailConfiguration) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"bcc_recipients": me.BccRecipients,
		"cc_recipients":  me.CcRecipients,
		"recipients":     me.Recipients,
	})
}

func (me *EmailConfiguration) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"bcc_recipients": &me.BccRecipients,
		"cc_recipients":  &me.CcRecipients,
		"recipients":     &me.Recipients,
	})
}
