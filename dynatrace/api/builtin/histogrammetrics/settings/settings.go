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

package histogrammetrics

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	EnableHistogramBucketIngest bool `json:"enableHistogramBucketIngest"` // When enabled, you can ingest the `le` dimension, representing explicit histogram buckets.\\\n Enable this if you are using OpenTelemetry histograms or Prometheus histogram metrics.\\\nWhen disabled, only your histograms' sum and count metrics will be ingested.
}

func (me *Settings) Name() string {
	return "histogram_metrics"
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enable_histogram_bucket_ingest": {
			Type:        schema.TypeBool,
			Description: "When enabled, you can ingest the `le` dimension, representing explicit histogram buckets.\\\n Enable this if you are using OpenTelemetry histograms or Prometheus histogram metrics.\\\nWhen disabled, only your histograms' sum and count metrics will be ingested.",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enable_histogram_bucket_ingest": me.EnableHistogramBucketIngest,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enable_histogram_bucket_ingest": &me.EnableHistogramBucketIngest,
	})
}
