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

package entrypoints

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SpanEntryPoint OpenTelemetry/OpenTracing spans can start new PurePaths. Define rules that define which spans should not be considered as entry points.\n\nNote: This config does not apply to Trace ingest
type SpanEntryPoint struct {
	EntryPointRule *SpanEntrypointRule `json:"entryPointRule"`
	InsertAfter    string              `json:"-"`
}

func (me *SpanEntryPoint) Name() string {
	return me.EntryPointRule.Name
}

func (me *SpanEntryPoint) Schema() map[string]*schema.Schema {
	var sch = new(SpanEntrypointRule).Schema()
	sch["insert_after"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Because this resource allows for ordering you may specify the ID of the resource instance that comes before this instance regarding order. If not specified when creating the setting will be added to the end of the list. If not specified during update the order will remain untouched",
		Optional:    true,
		Computed:    true,
	}
	return sch
}

func (me *SpanEntryPoint) MarshalHCL(properties hcl.Properties) error {
	if err := me.EntryPointRule.MarshalHCL(properties); err != nil {
		return err
	}
	if len(me.InsertAfter) > 0 {
		properties["insert_after"] = me.InsertAfter
	}
	return nil

}

func (me *SpanEntryPoint) UnmarshalHCL(decoder hcl.Decoder) error {
	me.EntryPointRule = new(SpanEntrypointRule)
	if err := me.EntryPointRule.UnmarshalHCL(decoder); err != nil {
		return err
	}
	if v, ok := decoder.GetOk("insert_after"); ok {
		if sv, ok := v.(string); ok {
			me.InsertAfter = sv
		}
	}
	return nil
}
