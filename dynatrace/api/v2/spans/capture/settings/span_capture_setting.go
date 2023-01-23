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

package capture

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SpanCaptureSetting OpenTelemetry/OpenTracing spans can start new PurePaths. Define rules that define which spans should not be considered as entry points.\n\nNote: This config does not apply to Trace ingest
type SpanCaptureSetting struct {
	SpanCaptureRule *SpanCaptureRule `json:"spanCaptureRule"`
}

func (me *SpanCaptureSetting) Name() string {
	return me.SpanCaptureRule.Name
}

func (me *SpanCaptureSetting) Schema() map[string]*schema.Schema {
	return new(SpanCaptureRule).Schema()
}

func (me *SpanCaptureSetting) MarshalHCL(properties hcl.Properties) error {
	return me.SpanCaptureRule.MarshalHCL(properties)
}

func (me *SpanCaptureSetting) UnmarshalHCL(decoder hcl.Decoder) error {
	me.SpanCaptureRule = new(SpanCaptureRule)
	return me.SpanCaptureRule.UnmarshalHCL(decoder)
}
