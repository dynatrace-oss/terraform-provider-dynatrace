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

package v2bindings

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Policies []*Policy

type Policy struct {
	ID         string
	Parameters map[string]string
	Metadata   map[string]string
	Boundaries []string
}

func (me *Policy) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Either the attribute `id` or the attribute `uuid` of a `dynatrace_iam_policy`. Initially just the `id` attribute was supported (which is a concatenation of several configuration settings) - and is still supported for backwards compatibility",
		},
		"parameters": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"metadata": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"boundaries": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *Policy) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("id", me.ID); err != nil {
		return err
	}
	if err := properties.Encode("parameters", me.Parameters); err != nil {
		return err
	}
	if err := properties.Encode("metadata", me.Metadata); err != nil {
		return err
	}
	if err := properties.Encode("boundaries", me.Boundaries); err != nil {
		return err
	}
	return nil
}

func (me *Policy) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("id", &me.ID); err != nil {
		return err
	}
	if err := decoder.Decode("parameters", &me.Parameters); err != nil {
		return err
	}
	if err := decoder.Decode("metadata", &me.Metadata); err != nil {
		return err
	}
	if err := decoder.Decode("boundaries", &me.Boundaries); err != nil {
		return err
	}
	return nil
}
