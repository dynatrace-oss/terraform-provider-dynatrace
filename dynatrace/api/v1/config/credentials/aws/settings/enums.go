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

package aws

// Statistic The statistic (aggregation) to be used for the metric. AVG_MIN_MAX value is 3 statistics at once: AVERAGE, MINIMUM and MAXIMUM
type Statistic string

// Statistics offers the known enum values
var Statistics = struct {
	Average     Statistic
	AvgMinMax   Statistic
	Maximum     Statistic
	Minimum     Statistic
	SampleCount Statistic
	Sum         Statistic
}{
	"AVERAGE",
	"AVG_MIN_MAX",
	"MAXIMUM",
	"MINIMUM",
	"SAMPLE_COUNT",
	"SUM",
}

// Type The type of the authentication: role-based or key-based.
type Type string

// Types offers the known enum values
var Types = struct {
	Keys Type
	Role Type
}{
	"KEYS",
	"ROLE",
}

// ConnectionStatus The status of the connection to the AWS environment.
//  * `CONNECTED`: There was a connection within last 10 minutes.
// * `DISCONNECTED`: A problem occurred with establishing connection using these credentials. Check whether the data is correct.
// * `UNINITIALIZED`: The successful connection has never been established for these credentials.
type ConnectionStatus string

// ConnectionStati offers the known enum values
var ConnectionStati = struct {
	Connected     ConnectionStatus
	Disconnected  ConnectionStatus
	Uninitialized ConnectionStatus
}{
	"CONNECTED",
	"DISCONNECTED",
	"UNINITIALIZED",
}

// PartitionType The type of the AWS partition.
type PartitionType string

// PartitionTypes offers the known enum values
var PartitionTypes = struct {
	AWSCn      PartitionType
	AWSDefault PartitionType
	AWSUsGov   PartitionType
}{
	"AWS_CN",
	"AWS_DEFAULT",
	"AWS_US_GOV",
}
