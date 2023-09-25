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

package synthetic

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CalculatedSyntheticMetric Descriptor of a calculated synthetic metric.
type CalculatedSyntheticMetric struct {
	Description       *string    `json:"description,omitempty"` // Descriptor of a calculated synthetic metric.
	MonitorIdentifier string     `json:"monitorIdentifier"`     // The Dynatrace entity ID of the application to which the metric belongs.
	Name              string     `json:"name"`                  // The displayed name of the metric.
	MetricKey         string     `json:"metricKey"`             // The unique key of the calculated synthetic metric.
	Enabled           bool       `json:"enabled"`               // The metric is enabled (`true`) or disabled (`false`).
	Metric            Metric     `json:"metric"`                // The type of the synthetic metric. Possible values: [ ApplicationCache, Callback, CumulativeLayoutShift, DNSLookup, DOMComplete, DOMContentLoaded, DOMInteractive, FailedRequestsResources, FirstContentfulPaint, FirstInputDelay, FirstInputStart, FirstPaint, HTMLDownloaded, HttpErrors, JavaScriptErrors, LargestContentfulPaint, LoadEventEnd, LoadEventStart, LongTasks, NavigationStart, OnDOMContentLoaded, OnLoad, Processing, RedirectTime, Request, RequestStart, ResourceCount, Response, SecureConnect, SpeedIndex, TCPConnect, TimeToFirstByte, TotalDuration, TransferSize, UserActionDuration, VisuallyComplete ]
	Dimensions        Dimensions `json:"dimensions,omitempty"`  // Dimension of the calculated synthetic metric.
	Filter            *Filter    `json:"filter,omitempty"`      // Filter of the calculated synthetic metric.
}

func (me *CalculatedSyntheticMetric) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Descriptor of a calculated synthetic metric.",
		},
		"monitor_identifier": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The Dynatrace entity ID of the monitor to which the metric belongs.",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The displayed name of the metric.",
		},
		"metric_key": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The unique key of the calculated synthetic metric.",
			ForceNew:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "The metric is enabled (`true`) or disabled (`false`)",
		},
		"metric": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The type of the synthetic metric. Possible values: [ ApplicationCache, Callback, CumulativeLayoutShift, DNSLookup, DOMComplete, DOMContentLoaded, DOMInteractive, FailedRequestsResources, FirstContentfulPaint, FirstInputDelay, FirstInputStart, FirstPaint, HTMLDownloaded, HttpErrors, JavaScriptErrors, LargestContentfulPaint, LoadEventEnd, LoadEventStart, LongTasks, NavigationStart, OnDOMContentLoaded, OnLoad, Processing, RedirectTime, Request, RequestStart, ResourceCount, Response, SecureConnect, SpeedIndex, TCPConnect, TimeToFirstByte, TotalDuration, TransferSize, UserActionDuration, VisuallyComplete ]",
		},
		"dimensions": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Dimension of the calculated synthetic metric.",
			Elem:        &schema.Resource{Schema: new(Dimensions).Schema()},
		},
		"filter": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Filter of the calculated synthetic metric.",
			Elem:        &schema.Resource{Schema: new(Filter).Schema()},
		},
	}
}

func (me *CalculatedSyntheticMetric) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"description":        me.Description,
		"monitor_identifier": me.MonitorIdentifier,
		"name":               me.Name,
		"metric_key":         me.MetricKey,
		"enabled":            me.Enabled,
		"metric":             me.Metric,
		"dimensions":         me.Dimensions,
		"filter":             me.Filter,
	})
}

func (me *CalculatedSyntheticMetric) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"description":        &me.Description,
		"monitor_identifier": &me.MonitorIdentifier,
		"name":               &me.Name,
		"metric_key":         &me.MetricKey,
		"enabled":            &me.Enabled,
		"metric":             &me.Metric,
		"dimensions":         &me.Dimensions,
		"filter":             &me.Filter,
	})
}

type Metric string

var Metrics = struct {
	ApplicationCache        Metric
	Callback                Metric
	CumulativeLayoutShift   Metric
	DNSLookup               Metric
	DOMComplete             Metric
	DOMContentLoaded        Metric
	DOMInteractive          Metric
	FailedRequestsResources Metric
	FirstContentfulPaint    Metric
	FirstInputDelay         Metric
	FirstInputStart         Metric
	FirstPaint              Metric
	HTMLDownloaded          Metric
	HttpErrors              Metric
	JavaScriptErrors        Metric
	LargestContentfulPaint  Metric
	LoadEventEnd            Metric
	LoadEventStart          Metric
	LongTasks               Metric
	NavigationStart         Metric
	OnDOMContentLoaded      Metric
	OnLoad                  Metric
	Processing              Metric
	RedirectTime            Metric
	Request                 Metric
	RequestStart            Metric
	ResourceCount           Metric
	Response                Metric
	SecureConnect           Metric
	SpeedIndex              Metric
	TCPConnect              Metric
	TimeToFirstByte         Metric
	TotalDuration           Metric
	TransferSize            Metric
	UserActionDuration      Metric
	VisuallyComplete        Metric
}{
	"ApplicationCache",
	"Callback",
	"CumulativeLayoutShift",
	"DNSLookup",
	"DOMComplete",
	"DOMContentLoaded",
	"DOMInteractive",
	"FailedRequestsResources",
	"FirstContentfulPaint",
	"FirstInputDelay",
	"FirstInputStart",
	"FirstPaint",
	"HTMLDownloaded",
	"HttpErrors",
	"JavaScriptErrors",
	"LargestContentfulPaint",
	"LoadEventEnd",
	"LoadEventStart",
	"LongTasks",
	"NavigationStart",
	"OnDOMContentLoaded",
	"OnLoad",
	"Processing",
	"RedirectTime",
	"Request",
	"RequestStart",
	"ResourceCount",
	"Response",
	"SecureConnect",
	"SpeedIndex",
	"TCPConnect",
	"TimeToFirstByte",
	"TotalDuration",
	"TransferSize",
	"UserActionDuration",
	"VisuallyComplete",
}
