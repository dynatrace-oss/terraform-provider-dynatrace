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

package ingestsources

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type BucketAssignmentAttributes struct {
	BucketName string `json:"bucketName"` // Bucket name
}

func (me *BucketAssignmentAttributes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"bucket_name": {
			Type:        schema.TypeString,
			Description: "Bucket name",
			Required:    true,
		},
	}
}

func (me *BucketAssignmentAttributes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"bucket_name": me.BucketName,
	})
}

func (me *BucketAssignmentAttributes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"bucket_name": &me.BucketName,
	})
}
