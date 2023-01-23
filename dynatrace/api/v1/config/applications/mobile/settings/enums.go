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

package mobile

type ApplicationType string

var ApplicationTypes = struct {
	CustomApplication ApplicationType
	MobileApplication ApplicationType
}{
	CustomApplication: ApplicationType("CUSTOM_APPLICATION"),
	MobileApplication: ApplicationType("MOBILE_APPLICATION"),
}

type BeaconEndpointType string

var BeaconEndpointTypes = struct {
	ClusterActiveGate     BeaconEndpointType
	EnvironmentActiveGate BeaconEndpointType
	InstrumentedWebServer BeaconEndpointType
}{
	ClusterActiveGate:     BeaconEndpointType("CLUSTER_ACTIVE_GATE"),
	EnvironmentActiveGate: BeaconEndpointType("ENVIRONMENT_ACTIVE_GATE"),
	InstrumentedWebServer: BeaconEndpointType("INSTRUMENTED_WEB_SERVER"),
}

type PropertyType string

var PropertyTypes = struct {
	Double PropertyType
	Long   PropertyType
	String PropertyType
}{
	Double: PropertyType("DOUBLE"),
	Long:   PropertyType("LONG"),
	String: PropertyType("STRING"),
}

type Origin string

var Origins = struct {
	API                        Origin
	ServerSideRequestAttribute Origin
}{
	API:                        Origin("API"),
	ServerSideRequestAttribute: Origin("SERVER_SIDE_REQUEST_ATTRIBUTE"),
}

type Aggregation string

var Aggregations = struct {
	First Aggregation
	Last  Aggregation
	Max   Aggregation
	Min   Aggregation
	Sum   Aggregation
}{
	First: Aggregation("FIRST"),
	Last:  Aggregation("LAST"),
	Max:   Aggregation("MAX"),
	Min:   Aggregation("MIN"),
	Sum:   Aggregation("SUM"),
}
