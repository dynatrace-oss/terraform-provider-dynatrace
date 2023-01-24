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
}

func (me *SpanEntryPoint) Name() string {
	return me.EntryPointRule.Name
}

func (me *SpanEntryPoint) Schema() map[string]*schema.Schema {
	return new(SpanEntrypointRule).Schema()
}

func (me *SpanEntryPoint) MarshalHCL(properties hcl.Properties) error {
	return me.EntryPointRule.MarshalHCL(properties)
}

func (me *SpanEntryPoint) UnmarshalHCL(decoder hcl.Decoder) error {
	me.EntryPointRule = new(SpanEntrypointRule)
	return me.EntryPointRule.UnmarshalHCL(decoder)
}
