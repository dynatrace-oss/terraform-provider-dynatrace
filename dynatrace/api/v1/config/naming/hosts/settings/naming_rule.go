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

package hosts

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// NamingRule The rule for the conditional naming.
type NamingRule struct {
	Name       string                     `json:"displayName"` // The name of the rule
	Enabled    bool                       `json:"enabled"`     // The rule is enabled (`true`) or disabled (`false`).
	Format     string                     `json:"nameFormat"`  // The name to be assigned to matching entities. You can use the following placeholders here:  * `{AwsAutoScalingGroup:Name}`  * `{AwsAvailabilityZone:Name}`  * `{AwsElasticLoadBalancer:Name}`  * `{AwsRelationalDatabaseService:DBName}`  * `{AwsRelationalDatabaseService:Endpoint}`  * `{AwsRelationalDatabaseService:Engine}`  * `{AwsRelationalDatabaseService:InstanceClass}`  * `{AwsRelationalDatabaseService:Name}`  * `{AwsRelationalDatabaseService:Port}`  * `{AzureRegion:Name}`  * `{AzureScaleSet:Name}`  * `{AzureVm:Name}`  * `{CloudFoundryOrganization:Name}`  * `{CustomDevice:DetectedName}`  * `{CustomDevice:DnsName}`  * `{CustomDevice:IpAddress}`  * `{CustomDevice:Port}`  * `{DockerContainerGroupInstance:ContainerName}`  * `{DockerContainerGroupInstance:FullImageName}`  * `{DockerContainerGroupInstance:ImageVersion}`  * `{DockerContainerGroupInstance:StrippedImageName}`  * `{ESXIHost:HardwareModel}`  * `{ESXIHost:HardwareVendor}`  * `{ESXIHost:Name}`  * `{ESXIHost:ProductName}`  * `{ESXIHost:ProductVersion}`  * `{Ec2Instance:AmiId}`  * `{Ec2Instance:BeanstalkEnvironmentName}`  * `{Ec2Instance:InstanceId}`  * `{Ec2Instance:InstanceType}`  * `{Ec2Instance:LocalHostName}`  * `{Ec2Instance:Name}`  * `{Ec2Instance:PublicHostName}`  * `{Ec2Instance:SecurityGroup}`  * `{GoogleComputeInstance:Id}`  * `{GoogleComputeInstance:IpAddresses}`  * `{GoogleComputeInstance:MachineType}`  * `{GoogleComputeInstance:Name}`  * `{GoogleComputeInstance:ProjectId}`  * `{GoogleComputeInstance:Project}`  * `{Host:AWSNameTag}`  * `{Host:AixLogicalCpuCount}`  * `{Host:AzureHostName}`  * `{Host:AzureSiteName}`  * `{Host:BoshDeploymentId}`  * `{Host:BoshInstanceId}`  * `{Host:BoshInstanceName}`  * `{Host:BoshName}`  * `{Host:BoshStemcellVersion}`  * `{Host:CpuCores}`  * `{Host:DetectedName}`  * `{Host:Environment:AppName}`  * `{Host:Environment:BoshReleaseVersion}`  * `{Host:Environment:Environment}`  * `{Host:Environment:Link}`  * `{Host:Environment:Organization}`  * `{Host:Environment:Owner}`  * `{Host:Environment:Support}`  * `{Host:IpAddress}`  * `{Host:LogicalCpuCores}`  * `{Host:OneAgentCustomHostName}`  * `{Host:OperatingSystemVersion}`  * `{Host:PaasMemoryLimit}`  * `{HostGroup:Name}`  * `{KubernetesCluster:Name}`  * `{KubernetesNode:DetectedName}`  * `{OpenstackAvailabilityZone:Name}`  * `{OpenstackZone:Name}`  * `{OpenstackComputeNode:Name}`  * `{OpenstackProject:Name}`  * `{OpenstackVm:InstanceType}`  * `{OpenstackVm:Name}`  * `{OpenstackVm:SecurityGroup}`  * `{ProcessGroup:AmazonECRImageAccountId}`  * `{ProcessGroup:AmazonECRImageRegion}`  * `{ProcessGroup:AmazonECSCluster}`  * `{ProcessGroup:AmazonECSContainerName}`  * `{ProcessGroup:AmazonECSFamily}`  * `{ProcessGroup:AmazonECSRevision}`  * `{ProcessGroup:AmazonLambdaFunctionName}`  * `{ProcessGroup:AmazonRegion}`  * `{ProcessGroup:ApacheConfigPath}`  * `{ProcessGroup:ApacheSparkMasterIpAddress}`  * `{ProcessGroup:AspDotNetCoreApplicationPath}`  * `{ProcessGroup:AspDotNetCoreApplicationPath}`  * `{ProcessGroup:AzureHostName}`  * `{ProcessGroup:AzureSiteName}`  * `{ProcessGroup:CassandraClusterName}`  * `{ProcessGroup:CatalinaBase}`  * `{ProcessGroup:CatalinaHome}`  * `{ProcessGroup:CloudFoundryAppId}`  * `{ProcessGroup:CloudFoundryAppName}`  * `{ProcessGroup:CloudFoundryInstanceIndex}`  * `{ProcessGroup:CloudFoundrySpaceId}`  * `{ProcessGroup:CloudFoundrySpaceName}`  * `{ProcessGroup:ColdFusionJvmConfigFile}`  * `{ProcessGroup:ColdFusionServiceName}`  * `{ProcessGroup:CommandLineArgs}`  * `{ProcessGroup:DetectedName}`  * `{ProcessGroup:DotNetCommandPath}`  * `{ProcessGroup:DotNetCommand}`  * `{ProcessGroup:DotNetClusterId}`  * `{ProcessGroup:DotNetNodeId}`  * `{ProcessGroup:ElasticsearchClusterName}`  * `{ProcessGroup:ElasticsearchNodeName}`  * `{ProcessGroup:EquinoxConfigPath}`  * `{ProcessGroup:ExeName}`  * `{ProcessGroup:ExePath}`  * `{ProcessGroup:GlassFishDomainName}`  * `{ProcessGroup:GlassFishInstanceName}`  * `{ProcessGroup:GoogleAppEngineInstance}`  * `{ProcessGroup:GoogleAppEngineService}`  * `{ProcessGroup:GoogleCloudProject}`  * `{ProcessGroup:HybrisBinDirectory}`  * `{ProcessGroup:HybrisConfigDirectory}`  * `{ProcessGroup:HybrisConfigDirectory}`  * `{ProcessGroup:HybrisDataDirectory}`  * `{ProcessGroup:IBMCicsRegion}`  * `{ProcessGroup:IBMCtgName}`  * `{ProcessGroup:IBMImsConnectRegion}`  * `{ProcessGroup:IBMImsControlRegion}`  * `{ProcessGroup:IBMImsMessageProcessingRegion}`  * `{ProcessGroup:IBMImsSoapGwName}`  * `{ProcessGroup:IBMIntegrationNodeName}`  * `{ProcessGroup:IBMIntegrationServerName}`  * `{ProcessGroup:IISAppPool}`  * `{ProcessGroup:IISRoleName}`  * `{ProcessGroup:JbossHome}`  * `{ProcessGroup:JbossMode}`  * `{ProcessGroup:JbossServerName}`  * `{ProcessGroup:JavaJarFile}`  * `{ProcessGroup:JavaJarPath}`  * `{ProcessGroup:JavaMainCLass}`  * `{ProcessGroup:KubernetesBasePodName}`  * `{ProcessGroup:KubernetesContainerName}`  * `{ProcessGroup:KubernetesFullPodName}`  * `{ProcessGroup:KubernetesNamespace}`  * `{ProcessGroup:KubernetesPodUid}`  * `{ProcessGroup:MssqlInstanceName}`  * `{ProcessGroup:NodeJsAppBaseDirectory}`  * `{ProcessGroup:NodeJsAppName}`  * `{ProcessGroup:NodeJsScriptName}`  * `{ProcessGroup:OracleSid}`  * `{ProcessGroup:PHPScriptPath}`  * `{ProcessGroup:PHPWorkingDirectory}`  * `{ProcessGroup:Ports}`  * `{ProcessGroup:RubyAppRootPath}`  * `{ProcessGroup:RubyScriptPath}`  * `{ProcessGroup:SoftwareAGInstallRoot}`  * `{ProcessGroup:SoftwareAGProductPropertyName}`  * `{ProcessGroup:SpringBootAppName}`  * `{ProcessGroup:SpringBootProfileName}`  * `{ProcessGroup:SpringBootStartupClass}`  * `{ProcessGroup:TIBCOBusinessWorksAppNodeName}`  * `{ProcessGroup:TIBCOBusinessWorksAppSpaceName}`  * `{ProcessGroup:TIBCOBusinessWorksCeAppName}`  * `{ProcessGroup:TIBCOBusinessWorksCeVersion}`  * `{ProcessGroup:TIBCOBusinessWorksDomainName}`  * `{ProcessGroup:TIBCOBusinessWorksEnginePropertyFilePath}`  * `{ProcessGroup:TIBCOBusinessWorksEnginePropertyFile}`  * `{ProcessGroup:TIBCOBusinessWorksHome}`  * `{ProcessGroup:VarnishInstanceName}`  * `{ProcessGroup:WebLogicClusterName}`  * `{ProcessGroup:WebLogicDomainName}`  * `{ProcessGroup:WebLogicHome}`  * `{ProcessGroup:WebLogicName}`  * `{ProcessGroup:WebSphereCellName}`  * `{ProcessGroup:WebSphereClusterName}`  * `{ProcessGroup:WebSphereNodeName}`  * `{ProcessGroup:WebSphereServerName}`  * `{ProcessGroup:ActorSystem}`  * `{Service:STGServerName}`  * `{Service:DatabaseHostName}`  * `{Service:DatabaseName}`  * `{Service:DatabaseVendor}`  * `{Service:DetectedName}`  * `{Service:EndpointPath}`  * `{Service:EndpointPathGatewayUrl}`  * `{Service:IIBApplicationName}`  * `{Service:MessageListenerClassName}`  * `{Service:Port}`  * `{Service:PublicDomainName}`  * `{Service:RemoteEndpoint}`  * `{Service:RemoteName}`  * `{Service:WebApplicationId}`  * `{Service:WebContextRoot}`  * `{Service:WebServerName}`  * `{Service:WebServiceNamespace}`  * `{Service:WebServiceName}`  * `{VmwareDatacenter:Name}`  * `{VmwareVm:Name}`
	Conditions Conditions                 `json:"rules"`       // A list of matching conditions of the rule.  The rule applies only if **all** conditions are fulfilled.
	Unknowns   map[string]json.RawMessage `json:"-"`
}

func (me *NamingRule) Validate() []string {
	result := []string{}
	if len(me.Conditions) > 0 {
		for _, cond := range me.Conditions {
			result = append(result, cond.Validate()...)
		}
	}
	return result
}

func (me *NamingRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the rule",
		},
		"enabled": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "The rule is enabled (`true`) or disabled (`false`)",
		},
		"format": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name to be assigned to matching entities. You can use the following placeholders here:  * `{AwsAutoScalingGroup:Name}`  * `{AwsAvailabilityZone:Name}`  * `{AwsElasticLoadBalancer:Name}`  * `{AwsRelationalDatabaseService:DBName}`  * `{AwsRelationalDatabaseService:Endpoint}`  * `{AwsRelationalDatabaseService:Engine}`  * `{AwsRelationalDatabaseService:InstanceClass}`  * `{AwsRelationalDatabaseService:Name}`  * `{AwsRelationalDatabaseService:Port}`  * `{AzureRegion:Name}`  * `{AzureScaleSet:Name}`  * `{AzureVm:Name}`  * `{CloudFoundryOrganization:Name}`  * `{CustomDevice:DetectedName}`  * `{CustomDevice:DnsName}`  * `{CustomDevice:IpAddress}`  * `{CustomDevice:Port}`  * `{DockerContainerGroupInstance:ContainerName}`  * `{DockerContainerGroupInstance:FullImageName}`  * `{DockerContainerGroupInstance:ImageVersion}`  * `{DockerContainerGroupInstance:StrippedImageName}`  * `{ESXIHost:HardwareModel}`  * `{ESXIHost:HardwareVendor}`  * `{ESXIHost:Name}`  * `{ESXIHost:ProductName}`  * `{ESXIHost:ProductVersion}`  * `{Ec2Instance:AmiId}`  * `{Ec2Instance:BeanstalkEnvironmentName}`  * `{Ec2Instance:InstanceId}`  * `{Ec2Instance:InstanceType}`  * `{Ec2Instance:LocalHostName}`  * `{Ec2Instance:Name}`  * `{Ec2Instance:PublicHostName}`  * `{Ec2Instance:SecurityGroup}`  * `{GoogleComputeInstance:Id}`  * `{GoogleComputeInstance:IpAddresses}`  * `{GoogleComputeInstance:MachineType}`  * `{GoogleComputeInstance:Name}`  * `{GoogleComputeInstance:ProjectId}`  * `{GoogleComputeInstance:Project}`  * `{Host:AWSNameTag}`  * `{Host:AixLogicalCpuCount}`  * `{Host:AzureHostName}`  * `{Host:AzureSiteName}`  * `{Host:BoshDeploymentId}`  * `{Host:BoshInstanceId}`  * `{Host:BoshInstanceName}`  * `{Host:BoshName}`  * `{Host:BoshStemcellVersion}`  * `{Host:CpuCores}`  * `{Host:DetectedName}`  * `{Host:Environment:AppName}`  * `{Host:Environment:BoshReleaseVersion}`  * `{Host:Environment:Environment}`  * `{Host:Environment:Link}`  * `{Host:Environment:Organization}`  * `{Host:Environment:Owner}`  * `{Host:Environment:Support}`  * `{Host:IpAddress}`  * `{Host:LogicalCpuCores}`  * `{Host:OneAgentCustomHostName}`  * `{Host:OperatingSystemVersion}`  * `{Host:PaasMemoryLimit}`  * `{HostGroup:Name}`  * `{KubernetesCluster:Name}`  * `{KubernetesNode:DetectedName}`  * `{OpenstackAvailabilityZone:Name}`  * `{OpenstackZone:Name}`  * `{OpenstackComputeNode:Name}`  * `{OpenstackProject:Name}`  * `{OpenstackVm:InstanceType}`  * `{OpenstackVm:Name}`  * `{OpenstackVm:SecurityGroup}`  * `{ProcessGroup:AmazonECRImageAccountId}`  * `{ProcessGroup:AmazonECRImageRegion}`  * `{ProcessGroup:AmazonECSCluster}`  * `{ProcessGroup:AmazonECSContainerName}`  * `{ProcessGroup:AmazonECSFamily}`  * `{ProcessGroup:AmazonECSRevision}`  * `{ProcessGroup:AmazonLambdaFunctionName}`  * `{ProcessGroup:AmazonRegion}`  * `{ProcessGroup:ApacheConfigPath}`  * `{ProcessGroup:ApacheSparkMasterIpAddress}`  * `{ProcessGroup:AspDotNetCoreApplicationPath}`  * `{ProcessGroup:AspDotNetCoreApplicationPath}`  * `{ProcessGroup:AzureHostName}`  * `{ProcessGroup:AzureSiteName}`  * `{ProcessGroup:CassandraClusterName}`  * `{ProcessGroup:CatalinaBase}`  * `{ProcessGroup:CatalinaHome}`  * `{ProcessGroup:CloudFoundryAppId}`  * `{ProcessGroup:CloudFoundryAppName}`  * `{ProcessGroup:CloudFoundryInstanceIndex}`  * `{ProcessGroup:CloudFoundrySpaceId}`  * `{ProcessGroup:CloudFoundrySpaceName}`  * `{ProcessGroup:ColdFusionJvmConfigFile}`  * `{ProcessGroup:ColdFusionServiceName}`  * `{ProcessGroup:CommandLineArgs}`  * `{ProcessGroup:DetectedName}`  * `{ProcessGroup:DotNetCommandPath}`  * `{ProcessGroup:DotNetCommand}`  * `{ProcessGroup:DotNetClusterId}`  * `{ProcessGroup:DotNetNodeId}`  * `{ProcessGroup:ElasticsearchClusterName}`  * `{ProcessGroup:ElasticsearchNodeName}`  * `{ProcessGroup:EquinoxConfigPath}`  * `{ProcessGroup:ExeName}`  * `{ProcessGroup:ExePath}`  * `{ProcessGroup:GlassFishDomainName}`  * `{ProcessGroup:GlassFishInstanceName}`  * `{ProcessGroup:GoogleAppEngineInstance}`  * `{ProcessGroup:GoogleAppEngineService}`  * `{ProcessGroup:GoogleCloudProject}`  * `{ProcessGroup:HybrisBinDirectory}`  * `{ProcessGroup:HybrisConfigDirectory}`  * `{ProcessGroup:HybrisConfigDirectory}`  * `{ProcessGroup:HybrisDataDirectory}`  * `{ProcessGroup:IBMCicsRegion}`  * `{ProcessGroup:IBMCtgName}`  * `{ProcessGroup:IBMImsConnectRegion}`  * `{ProcessGroup:IBMImsControlRegion}`  * `{ProcessGroup:IBMImsMessageProcessingRegion}`  * `{ProcessGroup:IBMImsSoapGwName}`  * `{ProcessGroup:IBMIntegrationNodeName}`  * `{ProcessGroup:IBMIntegrationServerName}`  * `{ProcessGroup:IISAppPool}`  * `{ProcessGroup:IISRoleName}`  * `{ProcessGroup:JbossHome}`  * `{ProcessGroup:JbossMode}`  * `{ProcessGroup:JbossServerName}`  * `{ProcessGroup:JavaJarFile}`  * `{ProcessGroup:JavaJarPath}`  * `{ProcessGroup:JavaMainCLass}`  * `{ProcessGroup:KubernetesBasePodName}`  * `{ProcessGroup:KubernetesContainerName}`  * `{ProcessGroup:KubernetesFullPodName}`  * `{ProcessGroup:KubernetesNamespace}`  * `{ProcessGroup:KubernetesPodUid}`  * `{ProcessGroup:MssqlInstanceName}`  * `{ProcessGroup:NodeJsAppBaseDirectory}`  * `{ProcessGroup:NodeJsAppName}`  * `{ProcessGroup:NodeJsScriptName}`  * `{ProcessGroup:OracleSid}`  * `{ProcessGroup:PHPScriptPath}`  * `{ProcessGroup:PHPWorkingDirectory}`  * `{ProcessGroup:Ports}`  * `{ProcessGroup:RubyAppRootPath}`  * `{ProcessGroup:RubyScriptPath}`  * `{ProcessGroup:SoftwareAGInstallRoot}`  * `{ProcessGroup:SoftwareAGProductPropertyName}`  * `{ProcessGroup:SpringBootAppName}`  * `{ProcessGroup:SpringBootProfileName}`  * `{ProcessGroup:SpringBootStartupClass}`  * `{ProcessGroup:TIBCOBusinessWorksAppNodeName}`  * `{ProcessGroup:TIBCOBusinessWorksAppSpaceName}`  * `{ProcessGroup:TIBCOBusinessWorksCeAppName}`  * `{ProcessGroup:TIBCOBusinessWorksCeVersion}`  * `{ProcessGroup:TIBCOBusinessWorksDomainName}`  * `{ProcessGroup:TIBCOBusinessWorksEnginePropertyFilePath}`  * `{ProcessGroup:TIBCOBusinessWorksEnginePropertyFile}`  * `{ProcessGroup:TIBCOBusinessWorksHome}`  * `{ProcessGroup:VarnishInstanceName}`  * `{ProcessGroup:WebLogicClusterName}`  * `{ProcessGroup:WebLogicDomainName}`  * `{ProcessGroup:WebLogicHome}`  * `{ProcessGroup:WebLogicName}`  * `{ProcessGroup:WebSphereCellName}`  * `{ProcessGroup:WebSphereClusterName}`  * `{ProcessGroup:WebSphereNodeName}`  * `{ProcessGroup:WebSphereServerName}`  * `{ProcessGroup:ActorSystem}`  * `{Service:STGServerName}`  * `{Service:DatabaseHostName}`  * `{Service:DatabaseName}`  * `{Service:DatabaseVendor}`  * `{Service:DetectedName}`  * `{Service:EndpointPath}`  * `{Service:EndpointPathGatewayUrl}`  * `{Service:IIBApplicationName}`  * `{Service:MessageListenerClassName}`  * `{Service:Port}`  * `{Service:PublicDomainName}`  * `{Service:RemoteEndpoint}`  * `{Service:RemoteName}`  * `{Service:WebApplicationId}`  * `{Service:WebContextRoot}`  * `{Service:WebServerName}`  * `{Service:WebServiceNamespace}`  * `{Service:WebServiceName}`  * `{VmwareDatacenter:Name}`  * `{VmwareVm:Name}",
		},
		"conditions": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "A list of matching conditions of the rule.  The rule applies only if **all** conditions are fulfilled",
			Elem:        &schema.Resource{Schema: new(Conditions).Schema()},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *NamingRule) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	delete(me.Unknowns, "metadata")
	return properties.EncodeAll(map[string]any{
		"name":       me.Name,
		"enabled":    me.Enabled,
		"format":     me.Format,
		"conditions": me.Conditions,
		"unknowns":   me.Unknowns,
	})
}

func (me *NamingRule) UnmarshalHCL(decoder hcl.Decoder) error {
	err := decoder.DecodeAll(map[string]any{
		"name":       &me.Name,
		"enabled":    &me.Enabled,
		"format":     &me.Format,
		"conditions": &me.Conditions,
		"unknowns":   &me.Unknowns,
	})
	delete(me.Unknowns, "metadata")
	return err
}

func (me *NamingRule) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	delete(properties, "id")
	delete(properties, "metadata")
	if err := properties.MarshalAll(map[string]any{
		"type":        "HOST",
		"displayName": me.Name,
		"enabled":     me.Enabled,
		"nameFormat":  me.Format,
		"rules":       me.Conditions,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *NamingRule) UnmarshalJSON(data []byte) error {
	properties := xjson.NewProperties(me.Unknowns)
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	delete(me.Unknowns, "id")
	delete(me.Unknowns, "metadata")
	var typ string
	err := properties.UnmarshalAll(map[string]any{
		"type":        &typ,
		"displayName": &me.Name,
		"enabled":     &me.Enabled,
		"nameFormat":  &me.Format,
		"rules":       &me.Conditions,
	})
	return err
}
