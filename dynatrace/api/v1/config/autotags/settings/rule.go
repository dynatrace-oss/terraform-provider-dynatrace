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

package autotags

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/entityruleengine"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Rule A rule for the auto-tag.
//
//	Defines the conditions of tag usage.
type Rule struct {
	Type             RuleType                      `json:"type"`                       // Type of entities to which the rule applies.
	Enabled          bool                          `json:"enabled"`                    // Tag rule is enabled (`true`) or disabled (`false`).
	PropagationTypes []PropagationType             `json:"propagationTypes,omitempty"` // How to apply the tag to underlying entities:  * `SERVICE_TO_PROCESS_GROUP_LIKE`: Apply to underlying process groups of matching services.  * `SERVICE_TO_HOST_LIKE`: Apply to underlying hosts of matching services.  * `PROCESS_GROUP_TO_HOST`: Apply to underlying hosts of matching process groups.  * `PROCESS_GROUP_TO_SERVICE`: Apply to all services provided by the process groups.  * `HOST_TO_PROCESS_GROUP_INSTANCE`: Apply to processes running on matching hosts.  *  `AZURE_TO_PG`: Apply to process groups connected to matching Azure entities.  * `AZURE_TO_SERVICE`: Apply to services provided by matching Azure entities.
	Conditions       []*entityruleengine.Condition `json:"conditions"`                 // A list of matching rules for the auto-tag.  The tag applies only when **all** conditions are fulfilled.
	ValueFormat      *string                       `json:"valueFormat,omitempty"`      // The value of the auto-tag. If specified, the tag is used in the `name:valueFormat` format.  For example, you can extend the `Infrastructure` tag to `Infrastructure:Windows` and `Infrastructure:Linux`.  You can use the following placeholders here:  * `{AwsAutoScalingGroup:Name}`  * `{AwsAvailabilityZone:Name}`  * `{AwsElasticLoadBalancer:Name}`  * `{AwsRelationalDatabaseService:DBName}`  * `{AwsRelationalDatabaseService:Endpoint}`  * `{AwsRelationalDatabaseService:Engine}`  * `{AwsRelationalDatabaseService:InstanceClass}`  * `{AwsRelationalDatabaseService:Name}`  * `{AwsRelationalDatabaseService:Port}`  * `{AzureRegion:Name}`  * `{AzureScaleSet:Name}`  * `{AzureVm:Name}`  * `{CloudFoundryOrganization:Name}`  * `{CustomDevice:DetectedName}`  * `{CustomDevice:DnsName}`  * `{CustomDevice:IpAddress}`  * `{CustomDevice:Port}`  * `{DockerContainerGroupInstance:ContainerName}`  * `{DockerContainerGroupInstance:FullImageName}`  * `{DockerContainerGroupInstance:ImageVersion}`  * `{DockerContainerGroupInstance:StrippedImageName}`  * `{ESXIHost:HardwareModel}`  * `{ESXIHost:HardwareVendor}`  * `{ESXIHost:Name}`  * `{ESXIHost:ProductName}`  * `{ESXIHost:ProductVersion}`  * `{Ec2Instance:AmiId}`  * `{Ec2Instance:BeanstalkEnvironmentName}`  * `{Ec2Instance:InstanceId}`  * `{Ec2Instance:InstanceType}`  * `{Ec2Instance:LocalHostName}`  * `{Ec2Instance:Name}`  * `{Ec2Instance:PublicHostName}`  * `{Ec2Instance:SecurityGroup}`  * `{GoogleComputeInstance:Id}`  * `{GoogleComputeInstance:IpAddresses}`  * `{GoogleComputeInstance:MachineType}`  * `{GoogleComputeInstance:Name}`  * `{GoogleComputeInstance:ProjectId}`  * `{GoogleComputeInstance:Project}`  * `{Host:AWSNameTag}`  * `{Host:AixLogicalCpuCount}`  * `{Host:AzureHostName}`  * `{Host:AzureSiteName}`  * `{Host:BoshDeploymentId}`  * `{Host:BoshInstanceId}`  * `{Host:BoshInstanceName}`  * `{Host:BoshName}`  * `{Host:BoshStemcellVersion}`  * `{Host:CpuCores}`  * `{Host:DetectedName}`  * `{Host:Environment:AppName}`  * `{Host:Environment:BoshReleaseVersion}`  * `{Host:Environment:Environment}`  * `{Host:Environment:Link}`  * `{Host:Environment:Organization}`  * `{Host:Environment:Owner}`  * `{Host:Environment:Support}`  * `{Host:IpAddress}`  * `{Host:LogicalCpuCores}`  * `{Host:OneAgentCustomHostName}`  * `{Host:OperatingSystemVersion}`  * `{Host:PaasMemoryLimit}`  * `{HostGroup:Name}`  * `{KubernetesCluster:Name}`  * `{KubernetesNode:DetectedName}`  * `{OpenstackAvailabilityZone:Name}`  * `{OpenstackZone:Name}`  * `{OpenstackComputeNode:Name}`  * `{OpenstackProject:Name}`  * `{OpenstackVm:UnstanceType}`  * `{OpenstackVm:Name}`  * `{OpenstackVm:SecurityGroup}`  * `{ProcessGroup:AmazonECRImageAccountId}`  * `{ProcessGroup:AmazonECRImageRegion}`  * `{ProcessGroup:AmazonECSCluster}`  * `{ProcessGroup:AmazonECSContainerName}`  * `{ProcessGroup:AmazonECSFamily}`  * `{ProcessGroup:AmazonECSRevision}`  * `{ProcessGroup:AmazonLambdaFunctionName}`  * `{ProcessGroup:AmazonRegion}`  * `{ProcessGroup:ApacheConfigPath}`  * `{ProcessGroup:ApacheSparkMasterIpAddress}`  * `{ProcessGroup:AspDotNetCoreApplicationPath}`  * `{ProcessGroup:AspDotNetCoreApplicationPath}`  * `{ProcessGroup:AzureHostName}`  * `{ProcessGroup:AzureSiteName}`  * `{ProcessGroup:CassandraClusterName}`  * `{ProcessGroup:CatalinaBase}`  * `{ProcessGroup:CatalinaHome}`  * `{ProcessGroup:CloudFoundryAppId}`  * `{ProcessGroup:CloudFoundryAppName}`  * `{ProcessGroup:CloudFoundryInstanceIndex}`  * `{ProcessGroup:CloudFoundrySpaceId}`  * `{ProcessGroup:CloudFoundrySpaceName}`  * `{ProcessGroup:ColdFusionJvmConfigFile}`  * `{ProcessGroup:ColdFusionServiceName}`  * `{ProcessGroup:CommandLineArgs}`  * `{ProcessGroup:DetectedName}`  * `{ProcessGroup:DotNetCommandPath}`  * `{ProcessGroup:DotNetCommand}`  * `{ProcessGroup:DotNetClusterId}`  * `{ProcessGroup:DotNetNodeId}`  * `{ProcessGroup:ElasticsearchClusterName}`  * `{ProcessGroup:ElasticsearchNodeName}`  * `{ProcessGroup:EquinoxConfigPath}`  * `{ProcessGroup:ExeName}`  * `{ProcessGroup:ExePath}`  * `{ProcessGroup:GlassFishDomainName}`  * `{ProcessGroup:GlassFishInstanceName}`  * `{ProcessGroup:GoogleAppEngineInstance}`  * `{ProcessGroup:GoogleAppEngineService}`  * `{ProcessGroup:GoogleCloudProject}`  * `{ProcessGroup:HybrisBinDirectory}`  * `{ProcessGroup:HybrisConfigDirectory}`  * `{ProcessGroup:HybrisConfigDirectory}`  * `{ProcessGroup:HybrisDataDirectory}`  * `{ProcessGroup:IBMCicsRegion}`  * `{ProcessGroup:IBMCtgName}`  * `{ProcessGroup:IBMImsConnectRegion}`  * `{ProcessGroup:IBMImsControlRegion}`  * `{ProcessGroup:IBMImsMessageProcessingRegion}`  * `{ProcessGroup:IBMImsSoapGwName}`  * `{ProcessGroup:IBMIntegrationNodeName}`  * `{ProcessGroup:IBMIntegrationServerName}`  * `{ProcessGroup:IISAppPool}`  * `{ProcessGroup:IISRoleName}`  * `{ProcessGroup:JbossHome}`  * `{ProcessGroup:JbossMode}`  * `{ProcessGroup:JbossServerName}`  * `{ProcessGroup:JavaJarFile}`  * `{ProcessGroup:JavaJarPath}`  * `{ProcessGroup:JavaMainCLass}`  * `{ProcessGroup:KubernetesBasePodName}`  * `{ProcessGroup:KubernetesContainerName}`  * `{ProcessGroup:KubernetesFullPodName}`  * `{ProcessGroup:KubernetesNamespace}`  * `{ProcessGroup:KubernetesPodUid}`  * `{ProcessGroup:MssqlInstanceName}`  * `{ProcessGroup:NodeJsAppBaseDirectory}`  * `{ProcessGroup:NodeJsAppName}`  * `{ProcessGroup:NodeJsScriptName}`  * `{ProcessGroup:OracleSid}`  * `{ProcessGroup:PHPScriptPath}`  * `{ProcessGroup:PHPWorkingDirectory}`  * `{ProcessGroup:Ports}`  * `{ProcessGroup:RubyAppRootPath}`  * `{ProcessGroup:RubyScriptPath}`  * `{ProcessGroup:SoftwareAGInstallRoot}`  * `{ProcessGroup:SoftwareAGProductPropertyName}`  * `{ProcessGroup:SpringBootAppName}`  * `{ProcessGroup:SpringBootProfileName}`  * `{ProcessGroup:SpringBootStartupClass}`  * `{ProcessGroup:TIBCOBusinessWorksAppNodeName}`  * `{ProcessGroup:TIBCOBusinessWorksAppSpaceName}`  * `{ProcessGroup:TIBCOBusinessWorksCeAppName}`  * `{ProcessGroup:TIBCOBusinessWorksCeVersion}`  * `{ProcessGroup:TIBCOBusinessWorksDomainName}`  * `{ProcessGroup:TIBCOBusinessWorksEnginePropertyFilePath}`  * `{ProcessGroup:TIBCOBusinessWorksEnginePropertyFile}`  * `{ProcessGroup:TIBCOBusinessWorksHome}`  * `{ProcessGroup:VarnishInstanceName}`  * `{ProcessGroup:WebLogicClusterName}`  * `{ProcessGroup:WebLogicDomainName}`  * `{ProcessGroup:WebLogicHome}`  * `{ProcessGroup:WebLogicName}`  * `{ProcessGroup:WebSphereCellName}`  * `{ProcessGroup:WebSphereClusterName}`  * `{ProcessGroup:WebSphereNodeName}`  * `{ProcessGroup:WebSphereServerName}`  * `{ProcessGroup:ActorSystem}`  * `{Service:STGServerName}`  * `{Service:DatabaseHostName}`  * `{Service:DatabaseName}`  * `{Service:DatabaseVendor}`  * `{Service:DetectedName}`  * `{Service:EndpointPath}`  * `{Service:EndpointPathGatewayUrl}`  * `{Service:IIBApplicationName}`  * `{Service:MessageListenerClassName}`  * `{Service:Port}`  * `{Service:PublicDomainName}`  * `{Service:RemoteEndpoint}`  * `{Service:RemoteName}`  * `{Service:WebApplicationId}`  * `{Service:WebContextRoot}`  * `{Service:WebServerName}`  * `{Service:WebServiceNamespace}`  * `{Service:WebServiceName}`  * `{VmwareDatacenter:Name}`  * `{VmwareVm:Name}`
	Normalization    *string                       `json:"normalization,omitempty"`    // Changes applied to the value after applying the value format. Default is LEAVE_TEXT_AS_IS
	Unknowns         map[string]json.RawMessage    `json:"-"`
}

func (mzr *Rule) Validate() []string {
	result := []string{}
	if len(mzr.Conditions) > 0 {
		for _, cond := range mzr.Conditions {
			result = append(result, cond.Validate()...)
		}
	}
	return result
}

func (me *Rule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "The type of Dynatrace entities the management zone can be applied to",
			Required:    true,
		},
		"normalization": {
			Type:        schema.TypeString,
			Description: "Changes applied to the value after applying the value format. Possible values are `LEAVE_TEXT_AS_IS`, `TO_LOWER_CASE` and `TO_UPPER_CASE`. Default is `LEAVE_TEXT_AS_IS`",
			Optional:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "The rule is enabled (`true`) or disabled (`false`)",
			Optional:    true,
		},
		"propagation_types": {
			Type:        schema.TypeList,
			Description: "How to apply the management zone to underlying entities:\n   - `SERVICE_TO_HOST_LIKE`: Apply to underlying hosts of matching services\n   - `SERVICE_TO_PROCESS_GROUP_LIKE`: Apply to underlying process groups of matching services\n   - `PROCESS_GROUP_TO_HOST`: Apply to underlying hosts of matching process groups\n   - `PROCESS_GROUP_TO_SERVICE`: Apply to all services provided by matching process groups\n   - `HOST_TO_PROCESS_GROUP_INSTANCE`: Apply to processes running on matching hosts\n   - `CUSTOM_DEVICE_GROUP_TO_CUSTOM_DEVICE`: Apply to custom devices in matching custom device groups\n   - `AZURE_TO_PG`: Apply to process groups connected to matching Azure entities\n   - `AZURE_TO_SERVICE`: Apply to services provided by matching Azure entities",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"conditions": {
			Type:        schema.TypeSet,
			Description: "A list of matching rules for the management zone. The management zone applies only if **all** conditions are fulfilled",
			MinItems:    1,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(entityruleengine.Condition).Schema(),
			},
		},
		"value_format": {
			Type:        schema.TypeString,
			Description: "The value of the auto-tag. If specified, the tag is used in the `name:valueFormat` format.  For example, you can extend the `Infrastructure` tag to `Infrastructure:Windows` and `Infrastructure:Linux`.  You can use the following placeholders here:  * `{AwsAutoScalingGroup:Name}`  * `{AwsAvailabilityZone:Name}`  * `{AwsElasticLoadBalancer:Name}`  * `{AwsRelationalDatabaseService:DBName}`  * `{AwsRelationalDatabaseService:Endpoint}`  * `{AwsRelationalDatabaseService:Engine}`  * `{AwsRelationalDatabaseService:InstanceClass}`  * `{AwsRelationalDatabaseService:Name}`  * `{AwsRelationalDatabaseService:Port}`  * `{AzureRegion:Name}`  * `{AzureScaleSet:Name}`  * `{AzureVm:Name}`  * `{CloudFoundryOrganization:Name}`  * `{CustomDevice:DetectedName}`  * `{CustomDevice:DnsName}`  * `{CustomDevice:IpAddress}`  * `{CustomDevice:Port}`  * `{DockerContainerGroupInstance:ContainerName}`  * `{DockerContainerGroupInstance:FullImageName}`  * `{DockerContainerGroupInstance:ImageVersion}`  * `{DockerContainerGroupInstance:StrippedImageName}`  * `{ESXIHost:HardwareModel}`  * `{ESXIHost:HardwareVendor}`  * `{ESXIHost:Name}`  * `{ESXIHost:ProductName}`  * `{ESXIHost:ProductVersion}`  * `{Ec2Instance:AmiId}`  * `{Ec2Instance:BeanstalkEnvironmentName}`  * `{Ec2Instance:InstanceId}`  * `{Ec2Instance:InstanceType}`  * `{Ec2Instance:LocalHostName}`  * `{Ec2Instance:Name}`  * `{Ec2Instance:PublicHostName}`  * `{Ec2Instance:SecurityGroup}`  * `{GoogleComputeInstance:Id}`  * `{GoogleComputeInstance:IpAddresses}`  * `{GoogleComputeInstance:MachineType}`  * `{GoogleComputeInstance:Name}`  * `{GoogleComputeInstance:ProjectId}`  * `{GoogleComputeInstance:Project}`  * `{Host:AWSNameTag}`  * `{Host:AixLogicalCpuCount}`  * `{Host:AzureHostName}`  * `{Host:AzureSiteName}`  * `{Host:BoshDeploymentId}`  * `{Host:BoshInstanceId}`  * `{Host:BoshInstanceName}`  * `{Host:BoshName}`  * `{Host:BoshStemcellVersion}`  * `{Host:CpuCores}`  * `{Host:DetectedName}`  * `{Host:Environment:AppName}`  * `{Host:Environment:BoshReleaseVersion}`  * `{Host:Environment:Environment}`  * `{Host:Environment:Link}`  * `{Host:Environment:Organization}`  * `{Host:Environment:Owner}`  * `{Host:Environment:Support}`  * `{Host:IpAddress}`  * `{Host:LogicalCpuCores}`  * `{Host:OneAgentCustomHostName}`  * `{Host:OperatingSystemVersion}`  * `{Host:PaasMemoryLimit}`  * `{HostGroup:Name}`  * `{KubernetesCluster:Name}`  * `{KubernetesNode:DetectedName}`  * `{OpenstackAvailabilityZone:Name}`  * `{OpenstackZone:Name}`  * `{OpenstackComputeNode:Name}`  * `{OpenstackProject:Name}`  * `{OpenstackVm:UnstanceType}`  * `{OpenstackVm:Name}`  * `{OpenstackVm:SecurityGroup}`  * `{ProcessGroup:AmazonECRImageAccountId}`  * `{ProcessGroup:AmazonECRImageRegion}`  * `{ProcessGroup:AmazonECSCluster}`  * `{ProcessGroup:AmazonECSContainerName}`  * `{ProcessGroup:AmazonECSFamily}`  * `{ProcessGroup:AmazonECSRevision}`  * `{ProcessGroup:AmazonLambdaFunctionName}`  * `{ProcessGroup:AmazonRegion}`  * `{ProcessGroup:ApacheConfigPath}`  * `{ProcessGroup:ApacheSparkMasterIpAddress}`  * `{ProcessGroup:AspDotNetCoreApplicationPath}`  * `{ProcessGroup:AspDotNetCoreApplicationPath}`  * `{ProcessGroup:AzureHostName}`  * `{ProcessGroup:AzureSiteName}`  * `{ProcessGroup:CassandraClusterName}`  * `{ProcessGroup:CatalinaBase}`  * `{ProcessGroup:CatalinaHome}`  * `{ProcessGroup:CloudFoundryAppId}`  * `{ProcessGroup:CloudFoundryAppName}`  * `{ProcessGroup:CloudFoundryInstanceIndex}`  * `{ProcessGroup:CloudFoundrySpaceId}`  * `{ProcessGroup:CloudFoundrySpaceName}`  * `{ProcessGroup:ColdFusionJvmConfigFile}`  * `{ProcessGroup:ColdFusionServiceName}`  * `{ProcessGroup:CommandLineArgs}`  * `{ProcessGroup:DetectedName}`  * `{ProcessGroup:DotNetCommandPath}`  * `{ProcessGroup:DotNetCommand}`  * `{ProcessGroup:DotNetClusterId}`  * `{ProcessGroup:DotNetNodeId}`  * `{ProcessGroup:ElasticsearchClusterName}`  * `{ProcessGroup:ElasticsearchNodeName}`  * `{ProcessGroup:EquinoxConfigPath}`  * `{ProcessGroup:ExeName}`  * `{ProcessGroup:ExePath}`  * `{ProcessGroup:GlassFishDomainName}`  * `{ProcessGroup:GlassFishInstanceName}`  * `{ProcessGroup:GoogleAppEngineInstance}`  * `{ProcessGroup:GoogleAppEngineService}`  * `{ProcessGroup:GoogleCloudProject}`  * `{ProcessGroup:HybrisBinDirectory}`  * `{ProcessGroup:HybrisConfigDirectory}`  * `{ProcessGroup:HybrisConfigDirectory}`  * `{ProcessGroup:HybrisDataDirectory}`  * `{ProcessGroup:IBMCicsRegion}`  * `{ProcessGroup:IBMCtgName}`  * `{ProcessGroup:IBMImsConnectRegion}`  * `{ProcessGroup:IBMImsControlRegion}`  * `{ProcessGroup:IBMImsMessageProcessingRegion}`  * `{ProcessGroup:IBMImsSoapGwName}`  * `{ProcessGroup:IBMIntegrationNodeName}`  * `{ProcessGroup:IBMIntegrationServerName}`  * `{ProcessGroup:IISAppPool}`  * `{ProcessGroup:IISRoleName}`  * `{ProcessGroup:JbossHome}`  * `{ProcessGroup:JbossMode}`  * `{ProcessGroup:JbossServerName}`  * `{ProcessGroup:JavaJarFile}`  * `{ProcessGroup:JavaJarPath}`  * `{ProcessGroup:JavaMainCLass}`  * `{ProcessGroup:KubernetesBasePodName}`  * `{ProcessGroup:KubernetesContainerName}`  * `{ProcessGroup:KubernetesFullPodName}`  * `{ProcessGroup:KubernetesNamespace}`  * `{ProcessGroup:KubernetesPodUid}`  * `{ProcessGroup:MssqlInstanceName}`  * `{ProcessGroup:NodeJsAppBaseDirectory}`  * `{ProcessGroup:NodeJsAppName}`  * `{ProcessGroup:NodeJsScriptName}`  * `{ProcessGroup:OracleSid}`  * `{ProcessGroup:PHPScriptPath}`  * `{ProcessGroup:PHPWorkingDirectory}`  * `{ProcessGroup:Ports}`  * `{ProcessGroup:RubyAppRootPath}`  * `{ProcessGroup:RubyScriptPath}`  * `{ProcessGroup:SoftwareAGInstallRoot}`  * `{ProcessGroup:SoftwareAGProductPropertyName}`  * `{ProcessGroup:SpringBootAppName}`  * `{ProcessGroup:SpringBootProfileName}`  * `{ProcessGroup:SpringBootStartupClass}`  * `{ProcessGroup:TIBCOBusinessWorksAppNodeName}`  * `{ProcessGroup:TIBCOBusinessWorksAppSpaceName}`  * `{ProcessGroup:TIBCOBusinessWorksCeAppName}`  * `{ProcessGroup:TIBCOBusinessWorksCeVersion}`  * `{ProcessGroup:TIBCOBusinessWorksDomainName}`  * `{ProcessGroup:TIBCOBusinessWorksEnginePropertyFilePath}`  * `{ProcessGroup:TIBCOBusinessWorksEnginePropertyFile}`  * `{ProcessGroup:TIBCOBusinessWorksHome}`  * `{ProcessGroup:VarnishInstanceName}`  * `{ProcessGroup:WebLogicClusterName}`  * `{ProcessGroup:WebLogicDomainName}`  * `{ProcessGroup:WebLogicHome}`  * `{ProcessGroup:WebLogicName}`  * `{ProcessGroup:WebSphereCellName}`  * `{ProcessGroup:WebSphereClusterName}`  * `{ProcessGroup:WebSphereNodeName}`  * `{ProcessGroup:WebSphereServerName}`  * `{ProcessGroup:ActorSystem}`  * `{Service:STGServerName}`  * `{Service:DatabaseHostName}`  * `{Service:DatabaseName}`  * `{Service:DatabaseVendor}`  * `{Service:DetectedName}`  * `{Service:EndpointPath}`  * `{Service:EndpointPathGatewayUrl}`  * `{Service:IIBApplicationName}`  * `{Service:MessageListenerClassName}`  * `{Service:Port}`  * `{Service:PublicDomainName}`  * `{Service:RemoteEndpoint}`  * `{Service:RemoteName}`  * `{Service:WebApplicationId}`  * `{Service:WebContextRoot}`  * `{Service:WebServerName}`  * `{Service:WebServiceNamespace}`  * `{Service:WebServiceName}`  * `{VmwareDatacenter:Name}`  * `{VmwareVm:Name}`",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *Rule) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("enabled", me.Enabled); err != nil {
		return err
	}
	if err := properties.Encode("type", me.Type); err != nil {
		return err
	}
	if err := properties.Encode("propagation_types", me.PropagationTypes); err != nil {
		return err
	}
	if err := properties.Encode("conditions", me.Conditions); err != nil {
		return err
	}
	if err := properties.Encode("value_format", me.ValueFormat); err != nil {
		return err
	}
	if err := properties.Encode("normalization", me.Normalization); err != nil {
		return err
	}
	return nil
}

func (me *Rule) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "type")
		delete(me.Unknowns, "enabled")
		delete(me.Unknowns, "propagation_types")
		delete(me.Unknowns, "conditions")
		delete(me.Unknowns, "value_format")
		delete(me.Unknowns, "normalization")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("type"); ok {
		me.Type = RuleType(value.(string))
	}
	if value, ok := decoder.GetOk("enabled"); ok {
		me.Enabled = value.(bool)
	}
	if _, ok := decoder.GetOk("propagation_types.#"); ok {
		me.PropagationTypes = []PropagationType{}
		if entries, ok := decoder.GetOk("propagation_types"); ok {
			for _, entry := range entries.([]any) {
				me.PropagationTypes = append(me.PropagationTypes, PropagationType(entry.(string)))
			}
		}
	}
	if value, ok := decoder.GetOk("conditions"); ok {
		conditionSet := value.(*schema.Set)

		for _, conditionMap := range conditionSet.List() {
			hash := conditionSet.F(conditionMap)
			entry := new(entityruleengine.Condition)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "conditions", hash)); err != nil {
				return err
			}
			me.Conditions = append(me.Conditions, entry)
		}
	}
	if value, ok := decoder.GetOk("value_format"); ok {
		me.ValueFormat = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("normalization"); ok {
		me.Normalization = opt.NewString(value.(string))
	}
	return nil
}

func (me *Rule) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(me.Unknowns) > 0 {
		for k, v := range me.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(me.Enabled)
		if err != nil {
			return nil, err
		}
		m["enabled"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(me.Type)
		if err != nil {
			return nil, err
		}
		m["type"] = rawMessage
	}
	if len(me.PropagationTypes) > 0 {
		rawMessage, err := json.Marshal(me.PropagationTypes)
		if err != nil {
			return nil, err
		}
		m["propagationTypes"] = rawMessage
	}
	if len(me.Conditions) > 0 {
		rawMessage, err := json.Marshal(me.Conditions)
		if err != nil {
			return nil, err
		}
		m["conditions"] = rawMessage
	}
	if me.ValueFormat != nil {
		rawMessage, err := json.Marshal(me.ValueFormat)
		if err != nil {
			return nil, err
		}
		m["valueFormat"] = rawMessage
	}
	if me.Normalization != nil {
		rawMessage, err := json.Marshal(me.Normalization)
		if err != nil {
			return nil, err
		}
		m["normalization"] = rawMessage
	}

	return json.Marshal(m)
}

func (me *Rule) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["enabled"]; found {
		if err := json.Unmarshal(v, &me.Enabled); err != nil {
			return err
		}
	}
	if v, found := m["type"]; found {
		if err := json.Unmarshal(v, &me.Type); err != nil {
			return err
		}
	}
	if v, found := m["propagationTypes"]; found {
		if err := json.Unmarshal(v, &me.PropagationTypes); err != nil {
			return err
		}
	}
	if v, found := m["conditions"]; found {
		if err := json.Unmarshal(v, &me.Conditions); err != nil {
			return err
		}
	}
	if v, found := m["valueFormat"]; found {
		if err := json.Unmarshal(v, &me.ValueFormat); err != nil {
			return err
		}
	}
	if v, found := m["normalization"]; found {
		if err := json.Unmarshal(v, &me.Normalization); err != nil {
			return err
		}
	}
	delete(m, "enabled")
	delete(m, "type")
	delete(m, "propagationTypes")
	delete(m, "conditions")
	delete(m, "valueFormat")
	delete(m, "normalization")
	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}

// RuleType Type of entities to which the rule applies.
type RuleType string

// RuleTypes offers the known enum values
var RuleTypes = struct {
	Application                  RuleType
	AWSApplicationLoadBalancer   RuleType
	AWSClassicLoadBalancer       RuleType
	AWSNetworkLoadBalancer       RuleType
	AWSRelationalDatabaseService RuleType
	Azure                        RuleType
	CustomApplication            RuleType
	CustomDevice                 RuleType
	DCRumApplication             RuleType
	ESXIHost                     RuleType
	ExternalSyntheticTest        RuleType
	Host                         RuleType
	HTTPCheck                    RuleType
	MobileApplication            RuleType
	ProcessGroup                 RuleType
	Service                      RuleType
	SyntheticTest                RuleType
}{
	"APPLICATION",
	"AWS_APPLICATION_LOAD_BALANCER",
	"AWS_CLASSIC_LOAD_BALANCER",
	"AWS_NETWORK_LOAD_BALANCER",
	"AWS_RELATIONAL_DATABASE_SERVICE",
	"AZURE",
	"CUSTOM_APPLICATION",
	"CUSTOM_DEVICE",
	"DCRUM_APPLICATION",
	"ESXI_HOST",
	"EXTERNAL_SYNTHETIC_TEST",
	"HOST",
	"HTTP_CHECK",
	"MOBILE_APPLICATION",
	"PROCESS_GROUP",
	"SERVICE",
	"SYNTHETIC_TEST",
}
