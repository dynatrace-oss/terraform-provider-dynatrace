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

package connection

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Name        string `json:"name"`      // The name of the EdgeConnect deployment
	Namespace   string `json:"namespace"` // The namespace where EdgeConnect is deployed
	Token       string `json:"token"`     // Token
	Uid         string `json:"uid"`       // A pseudo-ID for the cluster, set to the UID of the kube-system namespace
	InsertAfter string `json:"-"`
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the EdgeConnect deployment",
			Required:    true,
		},
		"namespace": {
			Type:        schema.TypeString,
			Description: "The namespace where EdgeConnect is deployed",
			Required:    true,
		},
		"token": {
			Type:        schema.TypeString,
			Description: "Token",
			Required:    true,
			Sensitive:   true,
		},
		"uid": {
			Type:        schema.TypeString,
			Description: "A pseudo-ID for the cluster, set to the UID of the kube-system namespace",
			Required:    true,
		},
		"insert_after": {
			Type:        schema.TypeString,
			Description: "Because this resource allows for ordering you may specify the ID of the resource instance that comes before this instance regarding order. If not specified when creating the setting will be added to the end of the list. If not specified during update the order will remain untouched",
			Optional:    true,
			Computed:    true,
			Deprecated:  "This resource is no longer ordered, please remove this attribute from the configuration",
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":         me.Name,
		"namespace":    me.Namespace,
		"token":        "${state.secret_value}",
		"uid":          me.Uid,
		"insert_after": me.InsertAfter,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":         &me.Name,
		"namespace":    &me.Namespace,
		"token":        &me.Token,
		"uid":          &me.Uid,
		"insert_after": &me.InsertAfter,
	})
}

const credsNotProvided = "REST API didn't provide token data"

func (me *Settings) FillDemoValues() []string {
	me.Token = "################"
	return []string{credsNotProvided}
}
