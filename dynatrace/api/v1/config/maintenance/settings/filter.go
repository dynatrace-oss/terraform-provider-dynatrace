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

package maintenance

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Filter A matching rule for Dynatrace entities.
type Filter struct {
	Type           *FilterType                `json:"type,omitempty"`           // The type of the Dynatrace entities (for example, hosts or services) you want to pick up by matching.
	MzID           *string                    `json:"mzId,omitempty"`           // The ID of a management zone to which the matched entities must belong.
	TagCombination *TagCombination            `json:"tagCombination,omitempty"` // The logic that applies when several tags are specified: AND/OR.  If not set, the OR logic is used.
	Tags           TagInfos                   `json:"tags"`                     // The tag you want to use for matching.  You can use custom tags from the UI, AWS tags, Cloud Foundry tags, OpenShift/Kubernetes, and tags based on environment variables.
	Unknowns       map[string]json.RawMessage `json:"-"`
}

func (me *Filter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "The type of the Dynatrace entities (for example, hosts or services) you want to pick up by matching",
			Optional:    true,
		},
		"mz_id": {
			Type:        schema.TypeString,
			Description: "The ID of a management zone to which the matched entities must belong",
			Optional:    true,
		},
		"tag_combination": {
			Type:        schema.TypeString,
			Description: "The logic that applies when several tags are specified: AND/OR.  If not set, the OR logic is used",
			Optional:    true,
		},
		"tags": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "The tag you want to use for matching.  You can use custom tags from the UI, AWS tags, Cloud Foundry tags, OpenShift/Kubernetes, and tags based on environment variables",
			Elem:        &schema.Resource{Schema: new(TagInfo).Schema()},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *Filter) MarshalHCL(properties hcl.Properties) error {
	if len(me.Unknowns) > 0 {
		data, err := json.Marshal(me.Unknowns)
		if err != nil {
			return err
		}
		if dmgmzid, ok := me.Unknowns["managementZoneId"]; ok {
			json.Unmarshal(dmgmzid, &me.MzID)
		}
		delete(me.Unknowns, "managementZoneId")
		if err := properties.Encode("unknowns", string(data)); err != nil {
			return err
		}
	}
	if err := properties.Encode("type", me.Type); err != nil {
		return err
	}
	if err := properties.Encode("mz_id", me.MzID); err != nil {
		return err
	}
	if err := properties.Encode("tag_combination", me.TagCombination); err != nil {
		return err
	}
	if err := properties.Encode("tags", me.Tags); err != nil {
		return err
	}
	return nil
}

func (me *Filter) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "type")
		delete(me.Unknowns, "mz_id")
		delete(me.Unknowns, "tag_combination")
		delete(me.Unknowns, "tags")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("type"); ok {
		me.Type = FilterType(value.(string)).Ref()
	}
	if value, ok := decoder.GetOk("mz_id"); ok {
		me.MzID = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("tag_combination"); ok {
		me.TagCombination = TagCombination(value.(string)).Ref()
	}
	if err := decoder.Decode("tags", &me.Tags); err != nil {
		return err
	}
	return nil
}

func (me *Filter) MarshalJSON() ([]byte, error) {
	m := xjson.NewProperties(me.Unknowns)
	if err := m.Marshal("type", me.Type); err != nil {
		return nil, err
	}
	if err := m.Marshal("mzId", me.MzID); err != nil {
		return nil, err
	}
	if err := m.Marshal("managementZoneId", me.MzID); err != nil {
		return nil, err
	}
	if err := m.Marshal("tagCombination", me.TagCombination); err != nil {
		return nil, err
	}
	if err := m.Marshal("tags", me.Tags); err != nil {
		return nil, err
	}
	// REST API doesn't accept omitting tags or specifying nil
	if me.Tags == nil {
		m["tags"] = []byte("[]")
	}
	return json.Marshal(m)
}

func (me *Filter) UnmarshalJSON(data []byte) error {
	m := xjson.Properties{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if err := m.Unmarshal("type", &me.Type); err != nil {
		return err
	}
	if err := m.Unmarshal("mzId", &me.MzID); err != nil {
		return err
	}
	if err := m.Unmarshal("managementZoneId", &me.MzID); err != nil {
		return err
	}
	if err := m.Unmarshal("tagCombination", &me.TagCombination); err != nil {
		return err
	}
	if err := m.Unmarshal("tags", &me.Tags); err != nil {
		return err
	}

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}

// TagCombination The logic that applies when several tags are specified: AND/OR.
// If not set, the OR logic is used.
type TagCombination string

func (me TagCombination) Ref() *TagCombination {
	return &me
}

// TagCombinations offers the known enum values
var TagCombinations = struct {
	And TagCombination
	Or  TagCombination
}{
	"AND",
	"OR",
}

// FilterType The type of the Dynatrace entities (for example, hosts or services) you want to pick up by matching.
type FilterType string

func (me FilterType) Ref() *FilterType {
	return &me
}

// FilterTypes offers the known enum values
var FilterTypes = struct {
	ApmSecurityGateway           FilterType
	Application                  FilterType
	ApplicationMethod            FilterType
	ApplicationMethodGroup       FilterType
	AppMonServer                 FilterType
	AppMonSystemProfile          FilterType
	AutoScalingGroup             FilterType
	AuxiliarySyntheticTest       FilterType
	AWSApplicationLoadBalancer   FilterType
	AWSAvailabilityZone          FilterType
	AWSCredentials               FilterType
	AWSLambdaFunction            FilterType
	AWSNetworkLoadBalancer       FilterType
	AzureAPIManagementService    FilterType
	AzureApplicationGateway      FilterType
	AzureCosmosDB                FilterType
	AzureCredentials             FilterType
	AzureEventHub                FilterType
	AzureEventHubNamespace       FilterType
	AzureFunctionApp             FilterType
	AzureIotHub                  FilterType
	AzureLoadBalancer            FilterType
	AzureMgmtGroup               FilterType
	AzureRedisCache              FilterType
	AzureRegion                  FilterType
	AzureServiceBusNamespace     FilterType
	AzureServiceBusQueue         FilterType
	AzureServiceBusTopic         FilterType
	AzureSQLDatabase             FilterType
	AzureSQLElasticPool          FilterType
	AzureSQLServer               FilterType
	AzureStorageAccount          FilterType
	AzureSubscription            FilterType
	AzureTenant                  FilterType
	AzureVM                      FilterType
	AzureVMScaleSet              FilterType
	AzureWebApp                  FilterType
	CfApplication                FilterType
	CfFoundation                 FilterType
	CinderVolume                 FilterType
	CloudApplication             FilterType
	CloudApplicationInstance     FilterType
	CloudApplicationNamespace    FilterType
	ContainerGroup               FilterType
	ContainerGroupInstance       FilterType
	CustomApplication            FilterType
	CustomDevice                 FilterType
	CustomDeviceGroup            FilterType
	DCRumApplication             FilterType
	DCRumService                 FilterType
	DCRumServiceInstance         FilterType
	DeviceApplicationMethod      FilterType
	Disk                         FilterType
	DockerContainerGroup         FilterType
	DockerContainerGroupInstance FilterType
	DynamoDBTable                FilterType
	EbsVolume                    FilterType
	EC2Instance                  FilterType
	ElasticLoadBalancer          FilterType
	Environment                  FilterType
	ExternalSyntheticTestStep    FilterType
	GcpZone                      FilterType
	Geolocation                  FilterType
	GeolocSite                   FilterType
	GoogleComputeEngine          FilterType
	Host                         FilterType
	HostGroup                    FilterType
	HTTPCheck                    FilterType
	HTTPCheckStep                FilterType
	Hypervisor                   FilterType
	KubernetesCluster            FilterType
	KubernetesNode               FilterType
	MobileApplication            FilterType
	NetworkInterface             FilterType
	NeutronSubnet                FilterType
	OpenStackProject             FilterType
	OpenStackRegion              FilterType
	OpenStackVM                  FilterType
	OS                           FilterType
	ProcessGroup                 FilterType
	ProcessGroupInstance         FilterType
	RelationalDatabaseService    FilterType
	Service                      FilterType
	ServiceInstance              FilterType
	ServiceMethod                FilterType
	ServiceMethodGroup           FilterType
	SwiftContainer               FilterType
	SyntheticLocation            FilterType
	SyntheticTest                FilterType
	SyntheticTestStep            FilterType
	Virtualmachine               FilterType
	VmwareDatacenter             FilterType
}{
	"APM_SECURITY_GATEWAY",
	"APPLICATION",
	"APPLICATION_METHOD",
	"APPLICATION_METHOD_GROUP",
	"APPMON_SERVER",
	"APPMON_SYSTEM_PROFILE",
	"AUTO_SCALING_GROUP",
	"AUXILIARY_SYNTHETIC_TEST",
	"AWS_APPLICATION_LOAD_BALANCER",
	"AWS_AVAILABILITY_ZONE",
	"AWS_CREDENTIALS",
	"AWS_LAMBDA_FUNCTION",
	"AWS_NETWORK_LOAD_BALANCER",
	"AZURE_API_MANAGEMENT_SERVICE",
	"AZURE_APPLICATION_GATEWAY",
	"AZURE_COSMOS_DB",
	"AZURE_CREDENTIALS",
	"AZURE_EVENT_HUB",
	"AZURE_EVENT_HUB_NAMESPACE",
	"AZURE_FUNCTION_APP",
	"AZURE_IOT_HUB",
	"AZURE_LOAD_BALANCER",
	"AZURE_MGMT_GROUP",
	"AZURE_REDIS_CACHE",
	"AZURE_REGION",
	"AZURE_SERVICE_BUS_NAMESPACE",
	"AZURE_SERVICE_BUS_QUEUE",
	"AZURE_SERVICE_BUS_TOPIC",
	"AZURE_SQL_DATABASE",
	"AZURE_SQL_ELASTIC_POOL",
	"AZURE_SQL_SERVER",
	"AZURE_STORAGE_ACCOUNT",
	"AZURE_SUBSCRIPTION",
	"AZURE_TENANT",
	"AZURE_VM",
	"AZURE_VM_SCALE_SET",
	"AZURE_WEB_APP",
	"CF_APPLICATION",
	"CF_FOUNDATION",
	"CINDER_VOLUME",
	"CLOUD_APPLICATION",
	"CLOUD_APPLICATION_INSTANCE",
	"CLOUD_APPLICATION_NAMESPACE",
	"CONTAINER_GROUP",
	"CONTAINER_GROUP_INSTANCE",
	"CUSTOM_APPLICATION",
	"CUSTOM_DEVICE",
	"CUSTOM_DEVICE_GROUP",
	"DCRUM_APPLICATION",
	"DCRUM_SERVICE",
	"DCRUM_SERVICE_INSTANCE",
	"DEVICE_APPLICATION_METHOD",
	"DISK",
	"DOCKER_CONTAINER_GROUP",
	"DOCKER_CONTAINER_GROUP_INSTANCE",
	"DYNAMO_DB_TABLE",
	"EBS_VOLUME",
	"EC2_INSTANCE",
	"ELASTIC_LOAD_BALANCER",
	"ENVIRONMENT",
	"EXTERNAL_SYNTHETIC_TEST_STEP",
	"GCP_ZONE",
	"GEOLOCATION",
	"GEOLOC_SITE",
	"GOOGLE_COMPUTE_ENGINE",
	"HOST",
	"HOST_GROUP",
	"HTTP_CHECK",
	"HTTP_CHECK_STEP",
	"HYPERVISOR",
	"KUBERNETES_CLUSTER",
	"KUBERNETES_NODE",
	"MOBILE_APPLICATION",
	"NETWORK_INTERFACE",
	"NEUTRON_SUBNET",
	"OPENSTACK_PROJECT",
	"OPENSTACK_REGION",
	"OPENSTACK_VM",
	"OS",
	"PROCESS_GROUP",
	"PROCESS_GROUP_INSTANCE",
	"RELATIONAL_DATABASE_SERVICE",
	"SERVICE",
	"SERVICE_INSTANCE",
	"SERVICE_METHOD",
	"SERVICE_METHOD_GROUP",
	"SWIFT_CONTAINER",
	"SYNTHETIC_LOCATION",
	"SYNTHETIC_TEST",
	"SYNTHETIC_TEST_STEP",
	"VIRTUALMACHINE",
	"VMWARE_DATACENTER",
}
