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

package web

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dMetricDefinition The definition of a calculated web metric.
type MetricDefinition struct {
	Metric      Metric  `json:"metric"`                // The type of the web application metric
	PropertyKey *string `json:"propertyKey,omitempty"` // The key of the user action property. Only applicable for DoubleProperty and LongProperty metrics.
}

func (me *MetricDefinition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metric": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The metric to be captured. Possible values are `Apdex`, `ApplicationCache`, `Callback`, `CumulativeLayoutShift`, `DNSLookup`, `DOMComplete`, `DOMContentLoaded`, `DOMInteractive`, `DoubleProperty`, `ErrorCount`, `FirstContentfulPaint`, `FirstInputDelay`, `FirstInputStart`, `FirstPaint`, `HTMLDownloaded`, `LargestContentfulPaint`, `LoadEventEnd`, `LoadEventStart`, `LongProperty`, `LongTasksTime`, `NavigationStart`, `OnDOMContentLoaded`, `OnLoad`, `Processing`, `RedirectTime`, `Request`, `RequestStart`, `Response`, `SecureConnect`, `SpeedIndex`, `TCPConnect`, `TimeToFirstByte`, `UserActionDuration`, `VisuallyComplete`",
		},
		"property_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The key of the user action property. Only applicable for DoubleProperty and LongProperty metrics.",
		},
	}
}

func (me *MetricDefinition) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"metric":       me.Metric,
		"property_key": me.PropertyKey,
	})
}

func (me *MetricDefinition) UnmarshalHCL(decoder hcl.Decoder) error {
	err := decoder.DecodeAll(map[string]any{
		"metric":       &me.Metric,
		"property_key": &me.PropertyKey,
	})
	return err
}

// Metric The metric to be captured.
type Metric string

// Metrics offers the known enum values
var Metrics = struct {
	Apdex                  Metric
	ApplicationCache       Metric
	Callback               Metric
	CumulativeLayoutShift  Metric
	DNSLookup              Metric
	DOMComplete            Metric
	DOMContentLoaded       Metric
	DOMInteractive         Metric
	DoubleProperty         Metric
	ErrorCount             Metric
	FirstContentfulPaint   Metric
	FirstInputDelay        Metric
	FirstInputStart        Metric
	FirstPaint             Metric
	HTMLDownloaded         Metric
	LargestContentfulPaint Metric
	LoadEventEnd           Metric
	LoadEventStart         Metric
	LongProperty           Metric
	LongTasksTime          Metric
	NavigationStart        Metric
	OnDOMContentLoaded     Metric
	OnLoad                 Metric
	Processing             Metric
	RedirectTime           Metric
	Request                Metric
	RequestStart           Metric
	Response               Metric
	SecureConnect          Metric
	SpeedIndex             Metric
	TCPConnect             Metric
	TimeToFirstByte        Metric
	UserActionDuration     Metric
	VisuallyComplete       Metric
}{
	"Apdex",
	"ApplicationCache",
	"Callback",
	"CumulativeLayoutShift",
	"DNSLookup",
	"DOMComplete",
	"DOMContentLoaded",
	"DOMInteractive",
	"DoubleProperty",
	"ErrorCount",
	"FirstContentfulPaint",
	"FirstInputDelay",
	"FirstInputStart",
	"FirstPaint",
	"HTMLDownloaded",
	"LargestContentfulPaint",
	"LoadEventEnd",
	"LoadEventStart",
	"LongProperty",
	"LongTasksTime",
	"NavigationStart",
	"OnDOMContentLoaded",
	"OnLoad",
	"Processing",
	"RedirectTime",
	"Request",
	"RequestStart",
	"Response",
	"SecureConnect",
	"SpeedIndex",
	"TCPConnect",
	"TimeToFirstByte",
	"UserActionDuration",
	"VisuallyComplete",
}
