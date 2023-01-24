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

package request

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Config contains the setup of the monitor
type Config struct {
	UserAgent            *string `json:"userAgent,omitempty"`            // The User agent of the request
	AcceptAnyCertificate *bool   `json:"acceptAnyCertificate,omitempty"` // If set to `false`, then the monitor fails with invalid SSL certificates.\n\nIf not set, the `false` option is used
	FollowRedirects      *bool   `json:"followRedirects,omitempty"`      // If set to `false`, redirects are reported as successful requests with response code 3xx.\n\nIf not set, the `false` option is used.
	RequestHeaders       Headers `json:"requestHeaders,omitempty"`       // By default, only the `User-Agent` header is set.\n\nYou can't set or modify this header here. Use the `userAgent` field for that.
}

func (me *Config) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"user_agent": {
			Type:        schema.TypeString,
			Description: "The User agent of the request",
			Optional:    true,
		},
		"accept_any_certificate": {
			Type:        schema.TypeBool,
			Description: "If set to `false`, then the monitor fails with invalid SSL certificates.\n\nIf not set, the `false` option is used",
			Optional:    true,
		},
		"follow_redirects": {
			Type:        schema.TypeBool,
			Description: "If set to `false`, redirects are reported as successful requests with response code 3xx.\n\nIf not set, the `false` option is used.",
			Optional:    true,
		},
		"headers": {
			Type:        schema.TypeList,
			Description: "The setup of the monitor",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Headers).Schema()},
		},
	}
}

func (me *Config) MarshalHCL(properties hcl.Properties) error {
	if me.UserAgent != nil && len(*me.UserAgent) > 0 {
		if err := properties.Encode("user_agent", me.UserAgent); err != nil {
			return err
		}
	}
	if me.AcceptAnyCertificate != nil && *me.AcceptAnyCertificate {
		if err := properties.Encode("accept_any_certificate", me.AcceptAnyCertificate); err != nil {
			return err
		}
	}
	if me.FollowRedirects != nil && *me.FollowRedirects {
		if err := properties.Encode("follow_redirects", me.FollowRedirects); err != nil {
			return err
		}
	} else {
		if err := properties.Encode("follow_redirects", false); err != nil {
			return err
		}
	}
	if err := properties.Encode("headers", me.RequestHeaders); err != nil {
		return err
	}
	return nil
}

func (me *Config) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("user_agent", &me.UserAgent); err != nil {
		return err
	}
	if err := decoder.Decode("accept_any_certificate", &me.AcceptAnyCertificate); err != nil {
		return err
	}
	if err := decoder.Decode("follow_redirects", &me.FollowRedirects); err != nil {
		return err
	}
	if _, ok := decoder.GetOk("headers.#"); ok {
		me.RequestHeaders = Headers{}
		if err := me.RequestHeaders.UnmarshalHCL(hcl.NewDecoder(decoder, "headers", 0)); err != nil {
			return err
		}
	}
	return nil
}
