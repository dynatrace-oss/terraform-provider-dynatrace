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

package dataforwarding

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type GcpConnection struct {
	BucketName   string `json:"bucketName"`   // GCS Bucket Name
	ConnectionID string `json:"connectionId"` // GCP connection
}

func (me *GcpConnection) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"bucket_name": {
			Type:        schema.TypeString,
			Description: "GCS Bucket Name",
			Required:    true,
		},
		"connection_id": {
			Type:        schema.TypeString,
			Description: "GCP connection",
			Required:    true,
		},
	}
}

func (me *GcpConnection) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"bucket_name":   me.BucketName,
		"connection_id": me.ConnectionID,
	})
}

func (me *GcpConnection) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"bucket_name":   &me.BucketName,
		"connection_id": &me.ConnectionID,
	})
}
