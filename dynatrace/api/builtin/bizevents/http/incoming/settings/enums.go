/**
* @license
* Copyright 2025 Dynatrace LLC
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

package incoming

type ComparisonEnum string

var ComparisonEnums = struct {
	Contains    ComparisonEnum
	EndsWith    ComparisonEnum
	Equals      ComparisonEnum
	Exists      ComparisonEnum
	NContains   ComparisonEnum
	NEndsWith   ComparisonEnum
	NEquals     ComparisonEnum
	NExists     ComparisonEnum
	NStartsWith ComparisonEnum
	StartsWith  ComparisonEnum
}{
	"CONTAINS",
	"ENDS_WITH",
	"EQUALS",
	"EXISTS",
	"N_CONTAINS",
	"N_ENDS_WITH",
	"N_EQUALS",
	"N_EXISTS",
	"N_STARTS_WITH",
	"STARTS_WITH",
}

type DataSourceEnum string

var DataSourceEnums = struct {
	RequestBody        DataSourceEnum
	RequestHeaders     DataSourceEnum
	RequestMethod      DataSourceEnum
	RequestParameters  DataSourceEnum
	RequestPath        DataSourceEnum
	RequestURL         DataSourceEnum
	ResponseBody       DataSourceEnum
	ResponseHeaders    DataSourceEnum
	ResponseStatuscode DataSourceEnum
}{
	"request.body",
	"request.headers",
	"request.method",
	"request.parameters",
	"request.path",
	"request.url",
	"response.body",
	"response.headers",
	"response.statusCode",
}

type DataSourceWithStaticStringEnum string

var DataSourceWithStaticStringEnums = struct {
	ConstantString     DataSourceWithStaticStringEnum
	RequestBody        DataSourceWithStaticStringEnum
	RequestHeaders     DataSourceWithStaticStringEnum
	RequestMethod      DataSourceWithStaticStringEnum
	RequestParameters  DataSourceWithStaticStringEnum
	RequestPath        DataSourceWithStaticStringEnum
	RequestURL         DataSourceWithStaticStringEnum
	ResponseBody       DataSourceWithStaticStringEnum
	ResponseHeaders    DataSourceWithStaticStringEnum
	ResponseStatuscode DataSourceWithStaticStringEnum
}{
	"constant.string",
	"request.body",
	"request.headers",
	"request.method",
	"request.parameters",
	"request.path",
	"request.url",
	"response.body",
	"response.headers",
	"response.statusCode",
}
