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

package match

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SpanMatcher struct {
	Source        Source     `json:"source"`
	SpanKindValue *SpanKind  `json:"spanKindValue"`
	Type          Comparison `json:"type"`
	SourceKey     *string    `json:"sourceKey"`
	Value         *string    `json:"value"`
	CaseSensitive bool       `json:"caseSensitive"`
}

func (me *SpanMatcher) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"comparison": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Possible values are `EQUALS`, `CONTAINS`, `STARTS_WITH`, `ENDS_WITH`, `DOES_NOT_EQUAL`, `DOES_NOT_CONTAIN`, `DOES_NOT_START_WITH` and `DOES_NOT_END_WITH`.",
		},
		"source": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "What to match against. Possible values are `SPAN_NAME`, `SPAN_KIND`, `ATTRIBUTE`, `INSTRUMENTATION_LIBRARY_NAME` and `INSTRUMENTATION_LIBRARY_VERSION`",
		},
		"case_sensitive": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether to match strings case sensitively or not",
		},
		"value": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The value to compare against. When `source` is `SPAN_KIND` the only allowed values are `INTERNAL`, `SERVER`, `CLIENT`, `PRODUCER` and `CONSUMER`",
		},
		"key": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the attribute if `source` is `ATTRIBUTE`",
		},
	}
}

func (me *SpanMatcher) MarshalHCL(properties hcl.Properties) error {

	m := map[string]any{
		"source":         me.Source,
		"comparison":     me.Type,
		"key":            me.SourceKey,
		"case_sensitive": me.CaseSensitive,
	}

	if me.Source == Sources.SpanKind && me.SpanKindValue != nil {
		m["value"] = string(*me.SpanKindValue)
	} else {
		m["value"] = *me.Value
	}
	if !me.CaseSensitive {
		delete(m, "case_sensitive")
	}
	return properties.EncodeAll(m)
}

func (me *SpanMatcher) UnmarshalHCL(decoder hcl.Decoder) error {
	m := map[string]any{
		"source":         &me.Source,
		"comparison":     &me.Type,
		"key":            &me.SourceKey,
		"case_sensitive": &me.CaseSensitive,
		"value":          &me.Value,
	}
	if err := decoder.DecodeAll(m); err != nil {
		return err
	}

	if me.Source == Sources.SpanKind {
		if me.Value != nil {
			me.SpanKindValue = SpanKind(*me.Value).Ref()
		}
		me.Value = nil
	}

	if me.Source == Sources.InstrumentationLibraryName {
		me.Source = Sources.InstrumentationScopeName
	} else if me.Source == Sources.InstrumentationLibraryVersion {
		me.Source = Sources.InstrumentationScopeVersion
	}

	return nil
}

type SpanMatchers []*SpanMatcher

func (me *SpanMatchers) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"match": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Matching strategies for the Span",
			Elem:        &schema.Resource{Schema: new(SpanMatcher).Schema()},
		},
	}
}

func (me SpanMatchers) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("match", me)
}

func (me *SpanMatchers) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("match", me)
}
