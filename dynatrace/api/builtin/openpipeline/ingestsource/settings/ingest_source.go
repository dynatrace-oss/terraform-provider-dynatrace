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

package settings

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/processors"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var AllowedKinds = []string{
	"logs", "events", "events.security", "security.events", "bizevents", "spans",
	"events.sdlc", "metrics", "usersessions", "davis.problems", "davis.events",
	"system.events", "azure.logs.forwarding", "user.events",
}

type IngestSource struct {
	Kind          string             `json:"-"`
	DefaultBucket *string            `json:"defaultBucket,omitempty"`
	DisplayName   string             `json:"displayName"`
	Enabled       bool               `json:"enabled"`
	PathSegment   string             `json:"pathSegment"`
	StaticRouting *PipelineReference `json:"staticRouting,omitempty"`
	Processing    *Processing        `json:"processing,omitempty"`
}

const DisplayNameMaxLength = 512
const PathSegmentMaxLength = 100
const PathSegmentRegex = "^[a-zA-Z](\\.?[a-zA-Z0-9])*$"

var PathSegmentErrorMessage = fmt.Sprintf("expected pattern: %s", PathSegmentRegex)

const DefaultBucketMaxLength = 500

func (is *IngestSource) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"kind": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "Indicates OpenPipeline data source",
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(AllowedKinds, true),
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Indicates if the object is active",
			Default:     true,
			Optional:    true,
		},
		"default_bucket": {
			Type:         schema.TypeString,
			Description:  "The default bucket assigned to records for the ingest source",
			ValidateFunc: validation.StringLenBetween(1, DefaultBucketMaxLength),
			Optional:     true,
		},
		"display_name": {
			Type:         schema.TypeString,
			Description:  "Display name of the ingest source",
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, DisplayNameMaxLength),
		},
		"path_segment": {
			Type:        schema.TypeString,
			Description: "The segment of the ingest source, which is applied to the base path",
			Required:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, PathSegmentMaxLength),
				validation.StringMatch(regexp.MustCompile(PathSegmentRegex), PathSegmentErrorMessage)),
		},
		"static_routing": {
			Type:        schema.TypeList,
			Description: "References a pipeline, if present",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(PipelineReference).Schema("static_routing.0.")},
			Optional:    true,
		},
		"processing": {
			Type:        schema.TypeList,
			Description: "The pre-processing done in the ingest source",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Processing).Schema()},
			Optional:    true,
		},
	}
}

func (is *IngestSource) MarshalHCL(properties hcl.Properties) error {
	err := properties.EncodeAll(map[string]any{
		"kind":           is.Kind,
		"enabled":        is.Enabled,
		"display_name":   is.DisplayName,
		"path_segment":   is.PathSegment,
		"default_bucket": is.DefaultBucket,
		"static_routing": is.StaticRouting,
		"processing":     is.Processing,
	})
	openpipeline.RemoveNils(properties)

	return err
}

func (is *IngestSource) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"kind":           &is.Kind,
		"enabled":        &is.Enabled,
		"default_bucket": &is.DefaultBucket,
		"display_name":   &is.DisplayName,
		"path_segment":   &is.PathSegment,
		"static_routing": &is.StaticRouting,
		"processing":     &is.Processing,
	})
}

func (is *IngestSource) MarshalJSON() ([]byte, error) {
	var temp = *is
	if temp.Processing == nil {
		// The API expected "processing" to be not null.
		temp.Processing = &Processing{}
	}
	return json.Marshal(temp)
}

type Processing struct {
	Processors []*processors.Processor `json:"processors,omitempty"`
}

func (p *Processing) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"processor": {
			Type:        schema.TypeList,
			Description: "Groups all processors applicable for processing in the ingest-source.\nApplicable processors types are dql, fieldsAdd, fieldsRemove, fieldsRename, and drop",
			Elem:        &schema.Resource{Schema: new(processors.Processor).Schema()},
			Required:    true,
		},
	}
}

func (p *Processing) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("processor", p.Processors)
}

func (p *Processing) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("processor", &p.Processors)
}
