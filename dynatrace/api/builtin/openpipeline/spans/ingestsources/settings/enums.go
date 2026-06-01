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

package ingestsources

type Aggregation string

var Aggregations = struct {
	Disabled Aggregation
	Enabled  Aggregation
}{
	"disabled",
	"enabled",
}

type AssignmentType string

var AssignmentTypes = struct {
	Constant           AssignmentType
	Field              AssignmentType
	Multivalueconstant AssignmentType
}{
	"constant",
	"field",
	"multiValueConstant",
}

type FieldExtractionStrategy string

var FieldExtractionStrategies = struct {
	Equals     FieldExtractionStrategy
	Startswith FieldExtractionStrategy
}{
	"equals",
	"startsWith",
}

type FieldExtractionType string

var FieldExtractionTypes = struct {
	Exclude    FieldExtractionType
	Include    FieldExtractionType
	Includeall FieldExtractionType
}{
	"exclude",
	"include",
	"includeAll",
}

type FieldValueExtractionType string

var FieldValueExtractionTypes = struct {
	Constant FieldValueExtractionType
	Field    FieldValueExtractionType
}{
	"constant",
	"field",
}

type GeoOutputField string

var GeoOutputFields = struct {
	Cityname            GeoOutputField
	Continentisocode    GeoOutputField
	Continentname       GeoOutputField
	Countryisocode      GeoOutputField
	Countryname         GeoOutputField
	Location            GeoOutputField
	Postalcode          GeoOutputField
	Regionisocode       GeoOutputField
	Regionname          GeoOutputField
	Subdivisionisocodes GeoOutputField
}{
	"cityName",
	"continentIsoCode",
	"continentName",
	"countryIsoCode",
	"countryName",
	"location",
	"postalCode",
	"regionIsoCode",
	"regionName",
	"subdivisionIsoCodes",
}

type IngestSourceType string

var IngestSourceTypes = struct {
	Extension IngestSourceType
	Http      IngestSourceType
}{
	"extension",
	"http",
}

type Measurement string

var Measurements = struct {
	Duration Measurement
	Field    Measurement
}{
	"duration",
	"field",
}

type PipelineType string

var PipelineTypes = struct {
	Builtin PipelineType
	Custom  PipelineType
}{
	"builtin",
	"custom",
}

type ProcessorType string

var ProcessorTypes = struct {
	Azurelogforwarding           ProcessorType
	Bizevent                     ProcessorType
	Bucketassignment             ProcessorType
	Costallocation               ProcessorType
	Countermetric                ProcessorType
	Davis                        ProcessorType
	Dql                          ProcessorType
	Drop                         ProcessorType
	Fieldsadd                    ProcessorType
	Fieldsremove                 ProcessorType
	Fieldsrename                 ProcessorType
	Geolookup                    ProcessorType
	Histogrammetric              ProcessorType
	Nostorage                    ProcessorType
	Productallocation            ProcessorType
	Samplingawarecountermetric   ProcessorType
	Samplingawarehistogrammetric ProcessorType
	Samplingawarevaluemetric     ProcessorType
	Sdlcevent                    ProcessorType
	Securitycontext              ProcessorType
	Securityevent                ProcessorType
	Smartscapeedge               ProcessorType
	Smartscapenode               ProcessorType
	Technology                   ProcessorType
	Valuemetric                  ProcessorType
}{
	"azureLogForwarding",
	"bizevent",
	"bucketAssignment",
	"costAllocation",
	"counterMetric",
	"davis",
	"dql",
	"drop",
	"fieldsAdd",
	"fieldsRemove",
	"fieldsRename",
	"geoLookup",
	"histogramMetric",
	"noStorage",
	"productAllocation",
	"samplingAwareCounterMetric",
	"samplingAwareHistogramMetric",
	"samplingAwareValueMetric",
	"sdlcEvent",
	"securityContext",
	"securityEvent",
	"smartscapeEdge",
	"smartscapeNode",
	"technology",
	"valueMetric",
}

type Sampling string

var Samplings = struct {
	Disabled Sampling
	Enabled  Sampling
}{
	"disabled",
	"enabled",
}
