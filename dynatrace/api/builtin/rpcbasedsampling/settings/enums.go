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

package rpcbasedsampling

type ComparisonType string

var ComparisonTypes = struct {
	Contains         ComparisonType
	DoesNotContain   ComparisonType
	DoesNotEndWith   ComparisonType
	DoesNotEqual     ComparisonType
	DoesNotStartWith ComparisonType
	EndsWith         ComparisonType
	Equals           ComparisonType
	StartsWith       ComparisonType
}{
	"CONTAINS",
	"DOES_NOT_CONTAIN",
	"DOES_NOT_END_WITH",
	"DOES_NOT_EQUAL",
	"DOES_NOT_START_WITH",
	"ENDS_WITH",
	"EQUALS",
	"STARTS_WITH",
}

type SamplingScaleFactor string

var SamplingScaleFactors = struct {
	SamplingScaleFactor0  SamplingScaleFactor
	SamplingScaleFactor1  SamplingScaleFactor
	SamplingScaleFactor10 SamplingScaleFactor
	SamplingScaleFactor11 SamplingScaleFactor
	SamplingScaleFactor12 SamplingScaleFactor
	SamplingScaleFactor13 SamplingScaleFactor
	SamplingScaleFactor14 SamplingScaleFactor
	SamplingScaleFactor2  SamplingScaleFactor
	SamplingScaleFactor3  SamplingScaleFactor
	SamplingScaleFactor4  SamplingScaleFactor
	SamplingScaleFactor5  SamplingScaleFactor
	SamplingScaleFactor6  SamplingScaleFactor
	SamplingScaleFactor8  SamplingScaleFactor
	SamplingScaleFactor9  SamplingScaleFactor
}{
	"0",
	"1",
	"10",
	"11",
	"12",
	"13",
	"14",
	"2",
	"3",
	"4",
	"5",
	"6",
	"8",
	"9",
}

type WireProtocolType string

var WireProtocolTypes = struct {
	WireProtocolType1  WireProtocolType
	WireProtocolType10 WireProtocolType
	WireProtocolType2  WireProtocolType
	WireProtocolType3  WireProtocolType
	WireProtocolType4  WireProtocolType
	WireProtocolType5  WireProtocolType
	WireProtocolType6  WireProtocolType
	WireProtocolType7  WireProtocolType
	WireProtocolType8  WireProtocolType
	WireProtocolType9  WireProtocolType
}{
	"1",
	"10",
	"2",
	"3",
	"4",
	"5",
	"6",
	"7",
	"8",
	"9",
}
