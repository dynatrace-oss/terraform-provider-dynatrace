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

package dashboards

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CustomFilterConfig Configuration of the custom filter of a tile
type CustomFilterConfig struct {
	Type                 CustomFilterConfigType         `json:"type"`                 // The type of the filter
	CustomName           string                         `json:"customName"`           // The name of the tile, set by user
	DefaultName          string                         `json:"defaultName"`          // The default name of the tile
	ChartConfig          *CustomFilterChartConfig       `json:"chartConfig"`          // Config Configuration of a custom chart
	FiltersPerEntityType map[string]map[string][]string `json:"filtersPerEntityType"` // A list of filters, applied to specific entity types
	Unknowns             map[string]json.RawMessage     `json:"-"`
}

func (me *CustomFilterConfig) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "The type of the filter. Possible values are `ALB`, `APPLICATION`, `APPLICATION_METHOD`, `APPMON`, `ASG`, `AWS_CREDENTIALS`, `AWS_CUSTOM_SERVICE`, `AWS_LAMBDA_FUNCTION`, `CLOUD_APPLICATION`, `CLOUD_APPLICATION_INSTANCE`, `CLOUD_APPLICATION_NAMESPACE`, `CONTAINER_GROUP_INSTANCE`, `CUSTOM_APPLICATION`, `CUSTOM_DEVICES`, `CUSTOM_SERVICES`, `DATABASE`, `DATABASE_KEY_REQUEST`, `DCRUM_APPLICATION`, `DCRUM_ENTITY`, `DYNAMO_DB`, `EBS`, `EC2`, `ELB`, `ENVIRONMENT`, `ESXI`, `EXTERNAL_SYNTHETIC_TEST`, `GLOBAL_BACKGROUND_ACTIVITY`, `HOST`, `IOT`, `KUBERNETES_CLUSTER`, `KUBERNETES_NODE`, `MDA_SERVICE`, `MIXED`, `MOBILE_APPLICATION`, `MONITORED_ENTITY`, `NLB`, `PG_BACKGROUND_ACTIVITY`, `PROBLEM`, `PROCESS_GROUP_INSTANCE`, `RDS`, `REMOTE_PLUGIN`, `SERVICE`, `SERVICE_KEY_REQUEST`, `SYNTHETIC_BROWSER_MONITOR`, `SYNTHETIC_HTTPCHECK`, `SYNTHETIC_HTTPCHECK_STEP`, `SYNTHETIC_LOCATION`, `SYNTHETIC_TEST`, `SYNTHETIC_TEST_STEP`, `UI_ENTITY`, `VIRTUAL_MACHINE`, `WEB_CHECK`.",
			Required:    true,
		},
		"custom_name": {
			Type:        schema.TypeString,
			Description: "The name of the tile, set by user",
			Required:    true,
		},
		"default_name": {
			Type:        schema.TypeString,
			Description: "The default name of the tile",
			Required:    true,
		},
		"chart_config": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration of a custom chart",
			Elem: &schema.Resource{
				Schema: new(CustomFilterChartConfig).Schema(),
			},
		},
		"filters": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration of a custom chart",
			Elem: &schema.Resource{
				Schema: new(FiltersPerEntityType).Schema(),
			},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *CustomFilterConfig) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("type", string(me.Type)); err != nil {
		return err
	}
	if err := properties.Encode("custom_name", me.CustomName); err != nil {
		return err
	}
	if err := properties.Encode("default_name", me.DefaultName); err != nil {
		return err
	}
	if err := properties.Encode("chart_config", me.ChartConfig); err != nil {
		return err
	}
	if len(me.FiltersPerEntityType) > 0 {
		filtersPerEntityType := &FiltersPerEntityType{Filters: []*FilterForEntityType{}}
		for k, v := range me.FiltersPerEntityType {
			filterForEntityType := &FilterForEntityType{
				EntityType: "",
				Filters:    []*FilterMatch{},
			}
			filterForEntityType.EntityType = k
			for mk, mv := range v {
				filterForEntityType.Filters = append(filterForEntityType.Filters, &FilterMatch{
					Key:    mk,
					Values: mv,
				})
			}
			filtersPerEntityType.Filters = append(filtersPerEntityType.Filters, filterForEntityType)
		}
		if err := properties.Encode("filters", filtersPerEntityType); err != nil {
			return err
		}
	}
	return nil
}

func (me *CustomFilterConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "type")
		delete(me.Unknowns, "custom_name")
		delete(me.Unknowns, "default_name")
		delete(me.Unknowns, "chart_config")
		delete(me.Unknowns, "filters")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("type"); ok {
		me.Type = CustomFilterConfigType(value.(string))
	}
	if value, ok := decoder.GetOk("custom_name"); ok {
		me.CustomName = value.(string)
	}
	if value, ok := decoder.GetOk("default_name"); ok {
		me.DefaultName = value.(string)
	}
	if _, ok := decoder.GetOk("chart_config.#"); ok {
		me.ChartConfig = new(CustomFilterChartConfig)
		if err := me.ChartConfig.UnmarshalHCL(hcl.NewDecoder(decoder, "chart_config", 0)); err != nil {
			return err
		}
	}
	if _, ok := decoder.GetOk("filters.#"); ok {
		filtersPerEntityType := new(FiltersPerEntityType)
		if err := filtersPerEntityType.UnmarshalHCL(hcl.NewDecoder(decoder, "filters", 0)); err != nil {
			return err
		}
		me.FiltersPerEntityType = map[string]map[string][]string{}
		if len(filtersPerEntityType.Filters) > 0 {
			for _, filterForEntityType := range filtersPerEntityType.Filters {
				filterMatches := map[string][]string{}
				for _, filterMatch := range filterForEntityType.Filters {
					filterMatches[filterMatch.Key] = filterMatch.Values
				}
				if len(filterMatches) > 0 {
					me.FiltersPerEntityType[filterForEntityType.EntityType] = filterMatches
				}
			}
		}
	}
	return nil
}

func (me *CustomFilterConfig) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(me.Unknowns) > 0 {
		for k, v := range me.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(me.Type)
		if err != nil {
			return nil, err
		}
		m["type"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(me.CustomName)
		if err != nil {
			return nil, err
		}
		m["customName"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(me.DefaultName)
		if err != nil {
			return nil, err
		}
		m["defaultName"] = rawMessage
	}
	if me.ChartConfig != nil {
		rawMessage, err := json.Marshal(me.ChartConfig)
		if err != nil {
			return nil, err
		}
		m["chartConfig"] = rawMessage
	}
	filtersPerEntityType := me.FiltersPerEntityType
	if filtersPerEntityType == nil {
		filtersPerEntityType = map[string]map[string][]string{}
	}
	{
		rawMessage, err := json.Marshal(filtersPerEntityType)
		if err != nil {
			return nil, err
		}
		m["filtersPerEntityType"] = rawMessage
	}
	return json.Marshal(m)
}

func (me *CustomFilterConfig) UnmarshalJSON(data []byte) error {
	m := xjson.Properties{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if err := m.Unmarshal("type", &me.Type); err != nil {
		return err
	}
	if err := m.Unmarshal("customName", &me.CustomName); err != nil {
		return err
	}
	if err := m.Unmarshal("defaultName", &me.DefaultName); err != nil {
		return err
	}
	if err := m.Unmarshal("chartConfig", &me.ChartConfig); err != nil {
		return err
	}
	if err := m.Unmarshal("filtersPerEntityType", &me.FiltersPerEntityType); err != nil {
		return err
	}
	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}

// CustomFilterConfigType has no documentation
type CustomFilterConfigType string

// CustomFilterConfigTypes offers the known enum values
var CustomFilterConfigTypes = struct {
	Alb                       CustomFilterConfigType
	Application               CustomFilterConfigType
	ApplicationMethod         CustomFilterConfigType
	Appmon                    CustomFilterConfigType
	Asg                       CustomFilterConfigType
	AwsCredentials            CustomFilterConfigType
	AwsCustomService          CustomFilterConfigType
	AwsLambdaFunction         CustomFilterConfigType
	CloudApplication          CustomFilterConfigType
	CloudApplicationInstance  CustomFilterConfigType
	CloudApplicationNamespace CustomFilterConfigType
	ContainerGroupInstance    CustomFilterConfigType
	CustomApplication         CustomFilterConfigType
	CustomDevices             CustomFilterConfigType
	CustomServices            CustomFilterConfigType
	Database                  CustomFilterConfigType
	DatabaseKeyRequest        CustomFilterConfigType
	DcrumApplication          CustomFilterConfigType
	DcrumEntity               CustomFilterConfigType
	DynamoDb                  CustomFilterConfigType
	Ebs                       CustomFilterConfigType
	Ec2                       CustomFilterConfigType
	Elb                       CustomFilterConfigType
	Environment               CustomFilterConfigType
	Esxi                      CustomFilterConfigType
	ExternalSyntheticTest     CustomFilterConfigType
	GlobalBackgroundActivity  CustomFilterConfigType
	Host                      CustomFilterConfigType
	Iot                       CustomFilterConfigType
	KubernetesCluster         CustomFilterConfigType
	KubernetesNode            CustomFilterConfigType
	MdaService                CustomFilterConfigType
	Mixed                     CustomFilterConfigType
	MobileApplication         CustomFilterConfigType
	MonitoredEntity           CustomFilterConfigType
	Nlb                       CustomFilterConfigType
	PgBackgroundActivity      CustomFilterConfigType
	Problem                   CustomFilterConfigType
	ProcessGroupInstance      CustomFilterConfigType
	Rds                       CustomFilterConfigType
	RemotePlugin              CustomFilterConfigType
	Service                   CustomFilterConfigType
	ServiceKeyRequest         CustomFilterConfigType
	SyntheticBrowserMonitor   CustomFilterConfigType
	SyntheticHTTPcheck        CustomFilterConfigType
	SyntheticHTTPcheckStep    CustomFilterConfigType
	SyntheticLocation         CustomFilterConfigType
	SyntheticTest             CustomFilterConfigType
	SyntheticTestStep         CustomFilterConfigType
	UUEntity                  CustomFilterConfigType
	VirtualMachine            CustomFilterConfigType
	WebCheck                  CustomFilterConfigType
}{
	"ALB",
	"APPLICATION",
	"APPLICATION_METHOD",
	"APPMON",
	"ASG",
	"AWS_CREDENTIALS",
	"AWS_CUSTOM_SERVICE",
	"AWS_LAMBDA_FUNCTION",
	"CLOUD_APPLICATION",
	"CLOUD_APPLICATION_INSTANCE",
	"CLOUD_APPLICATION_NAMESPACE",
	"CONTAINER_GROUP_INSTANCE",
	"CUSTOM_APPLICATION",
	"CUSTOM_DEVICES",
	"CUSTOM_SERVICES",
	"DATABASE",
	"DATABASE_KEY_REQUEST",
	"DCRUM_APPLICATION",
	"DCRUM_ENTITY",
	"DYNAMO_DB",
	"EBS",
	"EC2",
	"ELB",
	"ENVIRONMENT",
	"ESXI",
	"EXTERNAL_SYNTHETIC_TEST",
	"GLOBAL_BACKGROUND_ACTIVITY",
	"HOST",
	"IOT",
	"KUBERNETES_CLUSTER",
	"KUBERNETES_NODE",
	"MDA_SERVICE",
	"MIXED",
	"MOBILE_APPLICATION",
	"MONITORED_ENTITY",
	"NLB",
	"PG_BACKGROUND_ACTIVITY",
	"PROBLEM",
	"PROCESS_GROUP_INSTANCE",
	"RDS",
	"REMOTE_PLUGIN",
	"SERVICE",
	"SERVICE_KEY_REQUEST",
	"SYNTHETIC_BROWSER_MONITOR",
	"SYNTHETIC_HTTPCHECK",
	"SYNTHETIC_HTTPCHECK_STEP",
	"SYNTHETIC_LOCATION",
	"SYNTHETIC_TEST",
	"SYNTHETIC_TEST_STEP",
	"UI_ENTITY",
	"VIRTUAL_MACHINE",
	"WEB_CHECK",
}
