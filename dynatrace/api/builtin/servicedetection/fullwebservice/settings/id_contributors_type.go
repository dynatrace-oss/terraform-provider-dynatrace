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

package fullwebservice

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type IdContributorsType struct {
	ApplicationID             *ServiceIdContributor `json:"applicationId,omitempty"`       // Application identifier
	ContextRoot               *ContextIdContributor `json:"contextRoot,omitempty"`         // The context root is the first segment of the request URL after the Server name. For example, in the `www.dynatrace.com/support/help/dynatrace-api/` URL the context root is `/support`. The context root value can be found on the **Service overview page** under **Properties and tags**.
	DetectAsWebRequestService bool                  `json:"detectAsWebRequestService"`     // Detect the matching requests as full web services (false) or web request services (true).\n\nSetting this field to true prevents detecting of matching requests as full web services. A web request service is created instead. If you need to further modify the resulting web request service, you need to create a separate [Full web request rule](builtin:service-detection.full-web-request).
	ServerName                *ServiceIdContributor `json:"serverName,omitempty"`          // Server name
	WebServiceName            *ServiceIdContributor `json:"webServiceName,omitempty"`      // Web service name
	WebServiceNamespace       *ServiceIdContributor `json:"webServiceNamespace,omitempty"` // Web service namespace
}

func (me *IdContributorsType) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"application_id": {
			Type:        schema.TypeList,
			Description: "Application identifier",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(ServiceIdContributor).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"context_root": {
			Type:        schema.TypeList,
			Description: "The context root is the first segment of the request URL after the Server name. For example, in the `www.dynatrace.com/support/help/dynatrace-api/` URL the context root is `/support`. The context root value can be found on the **Service overview page** under **Properties and tags**.",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(ContextIdContributor).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"detect_as_web_request_service": {
			Type:        schema.TypeBool,
			Description: "Detect the matching requests as full web services (false) or web request services (true).\n\nSetting this field to true prevents detecting of matching requests as full web services. A web request service is created instead. If you need to further modify the resulting web request service, you need to create a separate [Full web request rule](builtin:service-detection.full-web-request).",
			Required:    true,
		},
		"server_name": {
			Type:        schema.TypeList,
			Description: "Server name",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(ServiceIdContributor).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"web_service_name": {
			Type:        schema.TypeList,
			Description: "Web service name",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(ServiceIdContributor).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"web_service_namespace": {
			Type:        schema.TypeList,
			Description: "Web service namespace",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(ServiceIdContributor).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *IdContributorsType) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"application_id":                me.ApplicationID,
		"context_root":                  me.ContextRoot,
		"detect_as_web_request_service": me.DetectAsWebRequestService,
		"server_name":                   me.ServerName,
		"web_service_name":              me.WebServiceName,
		"web_service_namespace":         me.WebServiceNamespace,
	})
}

func (me *IdContributorsType) HandlePreconditions() error {
	if me.ApplicationID == nil && !me.DetectAsWebRequestService {
		return fmt.Errorf("'application_id' must be specified if 'detect_as_web_request_service' is set to '%v'", me.DetectAsWebRequestService)
	}
	if me.ApplicationID != nil && me.DetectAsWebRequestService {
		return fmt.Errorf("'application_id' must not be specified if 'detect_as_web_request_service' is set to '%v'", me.DetectAsWebRequestService)
	}
	if me.ContextRoot == nil && !me.DetectAsWebRequestService {
		return fmt.Errorf("'context_root' must be specified if 'detect_as_web_request_service' is set to '%v'", me.DetectAsWebRequestService)
	}
	if me.ContextRoot != nil && me.DetectAsWebRequestService {
		return fmt.Errorf("'context_root' must not be specified if 'detect_as_web_request_service' is set to '%v'", me.DetectAsWebRequestService)
	}
	if me.ServerName == nil && !me.DetectAsWebRequestService {
		return fmt.Errorf("'server_name' must be specified if 'detect_as_web_request_service' is set to '%v'", me.DetectAsWebRequestService)
	}
	if me.ServerName != nil && me.DetectAsWebRequestService {
		return fmt.Errorf("'server_name' must not be specified if 'detect_as_web_request_service' is set to '%v'", me.DetectAsWebRequestService)
	}
	if me.WebServiceName == nil && !me.DetectAsWebRequestService {
		return fmt.Errorf("'web_service_name' must be specified if 'detect_as_web_request_service' is set to '%v'", me.DetectAsWebRequestService)
	}
	if me.WebServiceName != nil && me.DetectAsWebRequestService {
		return fmt.Errorf("'web_service_name' must not be specified if 'detect_as_web_request_service' is set to '%v'", me.DetectAsWebRequestService)
	}
	if me.WebServiceNamespace == nil && !me.DetectAsWebRequestService {
		return fmt.Errorf("'web_service_namespace' must be specified if 'detect_as_web_request_service' is set to '%v'", me.DetectAsWebRequestService)
	}
	if me.WebServiceNamespace != nil && me.DetectAsWebRequestService {
		return fmt.Errorf("'web_service_namespace' must not be specified if 'detect_as_web_request_service' is set to '%v'", me.DetectAsWebRequestService)
	}
	return nil
}

func (me *IdContributorsType) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"application_id":                &me.ApplicationID,
		"context_root":                  &me.ContextRoot,
		"detect_as_web_request_service": &me.DetectAsWebRequestService,
		"server_name":                   &me.ServerName,
		"web_service_name":              &me.WebServiceName,
		"web_service_namespace":         &me.WebServiceNamespace,
	})
}
