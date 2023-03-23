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

package externalwebrequest

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type IdContributorsType struct {
	ApplicationID    *ServiceIdContributor      `json:"applicationId"`    // Application identifier
	ContextRoot      *ContextIdContributor      `json:"contextRoot"`      // URL context root
	PortForServiceID bool                       `json:"portForServiceId"` // Let the Port contribute to the Service Id
	PublicDomainName *PublicDomainIdContributor `json:"publicDomainName"` // Public domain name
}

func (me *IdContributorsType) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"application_id": {
			Type:        schema.TypeList,
			Description: "Application identifier",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(ServiceIdContributor).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"context_root": {
			Type:        schema.TypeList,
			Description: "URL context root",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(ContextIdContributor).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"port_for_service_id": {
			Type:        schema.TypeBool,
			Description: "Let the Port contribute to the Service Id",
			Required:    true,
		},
		"public_domain_name": {
			Type:        schema.TypeList,
			Description: "Public domain name",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(PublicDomainIdContributor).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *IdContributorsType) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"application_id":      me.ApplicationID,
		"context_root":        me.ContextRoot,
		"port_for_service_id": me.PortForServiceID,
		"public_domain_name":  me.PublicDomainName,
	})
}

func (me *IdContributorsType) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"application_id":      &me.ApplicationID,
		"context_root":        &me.ContextRoot,
		"port_for_service_id": &me.PortForServiceID,
		"public_domain_name":  &me.PublicDomainName,
	})
}
