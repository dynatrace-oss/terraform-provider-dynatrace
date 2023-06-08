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

package externalwebservice

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type IdContributorsType struct {
	DetectAsWebRequestService bool                  `json:"detectAsWebRequestService"`  // Detect the matching requests as web request services instead of web services.\n\nThis prevents detecting of matching requests as opaque web services. An opaque web request service is created instead. If you need to further modify the resulting web request service, you need to create a separate [Opaque/external web request rule](builtin:service-detection.full-web-request).
	PortForServiceID          *bool                 `json:"portForServiceId,omitempty"` // Let the port contribute to the Service Id
	UrlPath                   *ServiceIdContributor `json:"urlPath,omitempty"`          // URL path
}

func (me *IdContributorsType) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"detect_as_web_request_service": {
			Type:        schema.TypeBool,
			Description: "Detect the matching requests as web request services instead of web services.\n\nThis prevents detecting of matching requests as opaque web services. An opaque web request service is created instead. If you need to further modify the resulting web request service, you need to create a separate [Opaque/external web request rule](builtin:service-detection.full-web-request).",
			Required:    true,
		},
		"port_for_service_id": {
			Type:        schema.TypeBool,
			Description: "Let the port contribute to the Service Id",
			Optional:    true, // precondition
		},
		"url_path": {
			Type:        schema.TypeList,
			Description: "URL path",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(ServiceIdContributor).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *IdContributorsType) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"detect_as_web_request_service": me.DetectAsWebRequestService,
		"port_for_service_id":           me.PortForServiceID,
		"url_path":                      me.UrlPath,
	})
}

func (me *IdContributorsType) HandlePreconditions() error {
	if me.PortForServiceID == nil && !me.DetectAsWebRequestService {
		me.PortForServiceID = opt.NewBool(false)
	}
	if me.UrlPath == nil && !me.DetectAsWebRequestService {
		return fmt.Errorf("'url_path' must be specified if 'detect_as_web_request_service' is set to '%v'", me.DetectAsWebRequestService)
	}
	if me.UrlPath != nil && me.DetectAsWebRequestService {
		return fmt.Errorf("'url_path' must not be specified if 'detect_as_web_request_service' is set to '%v'", me.DetectAsWebRequestService)
	}
	return nil
}

func (me *IdContributorsType) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"detect_as_web_request_service": &me.DetectAsWebRequestService,
		"port_for_service_id":           &me.PortForServiceID,
		"url_path":                      &me.UrlPath,
	})
}
