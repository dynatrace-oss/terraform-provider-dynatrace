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

type ApplicationType string

var ApplicationTypes = struct {
	AutoInjected             ApplicationType
	BrowserExtensionInjected ApplicationType
	ManuallyInjected         ApplicationType
}{
	"AUTO_INJECTED",
	"BROWSER_EXTENSION_INJECTED",
	"MANUALLY_INJECTED",
}

type ConversionGoalType string

var ConversionGoalTypes = struct {
	Destination     ConversionGoalType
	UserAction      ConversionGoalType
	VisitDuration   ConversionGoalType
	VisitNumActions ConversionGoalType
}{
	"Destination",
	"UserAction",
	"VisitDuration",
	"VisitNumActions",
}

type InjectionMode string

var InjectionModes = struct {
	CodeSnippet      InjectionMode
	CodeSnippetAsync InjectionMode
	InlineCode       InjectionMode
	JavaScriptTag    InjectionMode
}{
	"CODE_SNIPPET",
	"CODE_SNIPPET_ASYNC",
	"INLINE_CODE",
	"JAVASCRIPT_TAG",
}

type InjectionTarget string

var InjectionTargets = struct {
	PageQuery InjectionTarget
	URL       InjectionTarget
}{
	"PAGE_QUERY",
	"URL",
}

type URLOperator string

var URLOperators = struct {
	AllPages   URLOperator
	Contains   URLOperator
	EndsWith   URLOperator
	Equals     URLOperator
	StartsWith URLOperator
}{
	"ALL_PAGES",
	"CONTAINS",
	"ENDS_WITH",
	"EQUALS",
	"STARTS_WITH",
}

type JSInjectionRule string

var JSInjectionRules = struct {
	AfterSpecificHTML  JSInjectionRule
	AutomaticInjection JSInjectionRule
	BeforeSpecificHTML JSInjectionRule
	DoNotInject        JSInjectionRule
}{
	"AFTER_SPECIFIC_HTML",
	"AUTOMATIC_INJECTION",
	"BEFORE_SPECIFIC_HTML",
	"DO_NOT_INJECT",
}

type LoadActionKeyPerformanceMetric string

var LoadActionKeyPerformanceMetrics = struct {
	ActionDuration         LoadActionKeyPerformanceMetric
	CumulativeLayoutShift  LoadActionKeyPerformanceMetric
	DomInteractive         LoadActionKeyPerformanceMetric
	FirstInputDelay        LoadActionKeyPerformanceMetric
	LargestContentfulPaint LoadActionKeyPerformanceMetric
	LoadEventEnd           LoadActionKeyPerformanceMetric
	LoadEventStart         LoadActionKeyPerformanceMetric
	ResponseEnd            LoadActionKeyPerformanceMetric
	ResponseStart          LoadActionKeyPerformanceMetric
	SpeedIndex             LoadActionKeyPerformanceMetric
	VisuallyComplete       LoadActionKeyPerformanceMetric
}{
	"ACTION_DURATION",
	"CUMULATIVE_LAYOUT_SHIFT",
	"DOM_INTERACTIVE",
	"FIRST_INPUT_DELAY",
	"LARGEST_CONTENTFUL_PAINT",
	"LOAD_EVENT_END",
	"LOAD_EVENT_START",
	"RESPONSE_END",
	"RESPONSE_START",
	"SPEED_INDEX",
	"VISUALLY_COMPLETE",
}

type ResourceTimingCaptureType string

var ResourceTimingCaptureTypes = struct {
	CaptureAllSummaries     ResourceTimingCaptureType
	CaptureFullDetails      ResourceTimingCaptureType
	CaptureLimitedSummaries ResourceTimingCaptureType
}{
	"CAPTURE_ALL_SUMMARIES",
	"CAPTURE_FULL_DETAILS",
	"CAPTURE_LIMITED_SUMMARIES",
}

type XHRActionKeyPerformanceMetric string

var XHRActionKeyPerformanceMetrics = struct {
	ActionDuration   XHRActionKeyPerformanceMetric
	ResponseEnd      XHRActionKeyPerformanceMetric
	ResponseStart    XHRActionKeyPerformanceMetric
	VisuallyComplete XHRActionKeyPerformanceMetric
}{
	"ACTION_DURATION",
	"RESPONSE_END",
	"RESPONSE_START",
	"VISUALLY_COMPLETE",
}

type MetaDataCapturingType string

var MetaDataCapturingTypes = struct {
	Cookie             MetaDataCapturingType
	CSSSelector        MetaDataCapturingType
	JavaScriptFunction MetaDataCapturingType
	JavaScriptVariable MetaDataCapturingType
	MetaTag            MetaDataCapturingType
	QueryString        MetaDataCapturingType
}{
	"COOKIE",
	"CSS_SELECTOR",
	"JAVA_SCRIPT_FUNCTION",
	"JAVA_SCRIPT_VARIABLE",
	"META_TAG",
	"QUERY_STRING",
}

type PropertyType string

var PropertyTypes = struct {
	Date       PropertyType
	Double     PropertyType
	Long       PropertyType
	LongString PropertyType
	String     PropertyType
}{
	"DATE",
	"DOUBLE",
	"LONG",
	"LONG_STRING",
	"STRING",
}

type PropertyOrigin string

var PropertyOrigins = struct {
	JavaScriptAPI              PropertyOrigin
	MetaData                   PropertyOrigin
	ServerSideRequestAttribute PropertyOrigin
}{
	"JAVASCRIPT_API",
	"META_DATA",
	"SERVER_SIDE_REQUEST_ATTRIBUTE",
}

type Aggregation string

var Aggregations = struct {
	Average Aggregation
	First   Aggregation
	Last    Aggregation
	Maximum Aggregation
	Minimum Aggregation
	Sum     Aggregation
}{
	"AVERAGE",
	"FIRST",
	"LAST",
	"MAXIMUM",
	"MINIMUM",
	"SUM",
}

type KeyUserActionType string

var KeyUserActionTypes = struct {
	Custom KeyUserActionType
	Load   KeyUserActionType
	XHR    KeyUserActionType
}{
	"Custom",
	"Load",
	"Xhr",
}
