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

package urlbasedsampling

type HttpMethod string

var HttpMethods = struct {
	Get     HttpMethod
	Post    HttpMethod
	Put     HttpMethod
	Delete  HttpMethod
	Head    HttpMethod
	Connect HttpMethod
	Options HttpMethod
	Trace   HttpMethod
	Patch   HttpMethod
}{
	"GET",
	"POST",
	"PUT",
	"DELETE",
	"HEAD",
	"CONNECT",
	"OPTIONS",
	"TRACE",
	"PATCH",
}

type PathComparisonType string

var PathComparisonTypes = struct {
	Equals           PathComparisonType
	DoesNotEqual     PathComparisonType
	Contains         PathComparisonType
	DoesNotContain   PathComparisonType
	StartsWith       PathComparisonType
	DoesNotStartWith PathComparisonType
	EndsWith         PathComparisonType
	DoesNotEndWith   PathComparisonType
}{
	"EQUALS",
	"DOES_NOT_EQUAL",
	"CONTAINS",
	"DOES_NOT_CONTAIN",
	"STARTS_WITH",
	"DOES_NOT_START_WITH",
	"ENDS_WITH",
	"DOES_NOT_END_WITH",
}

type SamplingScaleFactor string

var SamplingScaleFactors = struct {
	IncreaseCapturing128Times  SamplingScaleFactor
	IncreaseCapturing64Times   SamplingScaleFactor
	IncreaseCapturing32Times   SamplingScaleFactor
	IncreaseCapturing16Times   SamplingScaleFactor
	IncreaseCapturing8Times    SamplingScaleFactor
	IncreaseCapturing4Times    SamplingScaleFactor
	IncreaseCapturing2Times    SamplingScaleFactor
	ReduceCapturingByFactor2   SamplingScaleFactor
	ReduceCapturingByFactor4   SamplingScaleFactor
	ReduceCapturingByFactor8   SamplingScaleFactor
	ReduceCapturingByFactor16  SamplingScaleFactor
	ReduceCapturingByFactor32  SamplingScaleFactor
	ReduceCapturingByFactor64  SamplingScaleFactor
	ReduceCapturingByFactor128 SamplingScaleFactor
}{
	"0",
	"1",
	"2",
	"3",
	"4",
	"5",
	"6",
	"8",
	"9",
	"10",
	"11",
	"12",
	"13",
	"14",
}

var SamplingScaleFactorLookup = map[string]SamplingScaleFactor{
	"IncreaseCapturing128Times":  SamplingScaleFactors.IncreaseCapturing128Times,
	"IncreaseCapturing64Times":   SamplingScaleFactors.IncreaseCapturing64Times,
	"IncreaseCapturing32Times":   SamplingScaleFactors.IncreaseCapturing32Times,
	"IncreaseCapturing16Times":   SamplingScaleFactors.IncreaseCapturing16Times,
	"IncreaseCapturing8Times":    SamplingScaleFactors.IncreaseCapturing8Times,
	"IncreaseCapturing4Times":    SamplingScaleFactors.IncreaseCapturing4Times,
	"IncreaseCapturing2Times":    SamplingScaleFactors.IncreaseCapturing2Times,
	"ReduceCapturingByFactor2":   SamplingScaleFactors.ReduceCapturingByFactor2,
	"ReduceCapturingByFactor4":   SamplingScaleFactors.ReduceCapturingByFactor4,
	"ReduceCapturingByFactor8":   SamplingScaleFactors.ReduceCapturingByFactor8,
	"ReduceCapturingByFactor16":  SamplingScaleFactors.ReduceCapturingByFactor16,
	"ReduceCapturingByFactor32":  SamplingScaleFactors.ReduceCapturingByFactor32,
	"ReduceCapturingByFactor64":  SamplingScaleFactors.ReduceCapturingByFactor64,
	"ReduceCapturingByFactor128": SamplingScaleFactors.ReduceCapturingByFactor128,
}
