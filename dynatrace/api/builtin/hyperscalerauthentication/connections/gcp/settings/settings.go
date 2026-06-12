/**
* @license
* Copyright 2026 Dynatrace LLC
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

package gcp

import (
	"fmt"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const DefaultCreateTimeout = 2 * time.Minute

type Settings struct {
	Name                        string                       `json:"name"` // The name of the connection
	ServiceAccountImpersonation *ServiceAccountImpersonation `json:"serviceAccountImpersonation,omitempty"`
	Type                        Type                         `json:"type"` // GCP Authentication mechanism to be used by the connection. Possible values: `serviceAccountImpersonation`
}

func (me *Settings) Timeouts() *schema.ResourceTimeout {
	return &schema.ResourceTimeout{
		Create: schema.DefaultTimeout(DefaultCreateTimeout),
	}
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the connection",
			Required:    true,
		},
		"service_account_impersonation": {
			Type:        schema.TypeList,
			Description: "No documentation available",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(ServiceAccountImpersonation).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "GCP Authentication mechanism to be used by the connection. Possible values: `serviceAccountImpersonation`",
			ForceNew:    true,
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":                          me.Name,
		"service_account_impersonation": me.ServiceAccountImpersonation,
		"type":                          me.Type,
	})
}

func (me *Settings) HandlePreconditions() error {
	if (me.ServiceAccountImpersonation != nil) && (string(me.Type) != "serviceAccountImpersonation") {
		return fmt.Errorf("'service_account_impersonation' must not be specified unless 'type' is set to 'serviceAccountImpersonation'; got 'type'='%v'", me.Type)
	}
	if (me.ServiceAccountImpersonation == nil) && (string(me.Type) == "serviceAccountImpersonation") {
		return fmt.Errorf("'service_account_impersonation' must be specified when 'type' is set to 'serviceAccountImpersonation'; got 'type'='%v'", me.Type)
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":                          &me.Name,
		"service_account_impersonation": &me.ServiceAccountImpersonation,
		"type":                          &me.Type,
	})
}
