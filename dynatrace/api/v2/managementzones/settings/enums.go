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

package managementzones

type DimensionType string

var DimensionTypes = struct {
	Metric DimensionType
	Any    DimensionType
	Log    DimensionType
}{
	DimensionType("METRIC"),
	DimensionType("ANY"),
	DimensionType("LOG"),
}

type ManagementZoneMeType string

var ManagementZoneMeTypes = struct {
	Azure                        ManagementZoneMeType
	DataCenterService            ManagementZoneMeType
	Service                      ManagementZoneMeType
	MobileApplication            ManagementZoneMeType
	CloudFoundryFoundation       ManagementZoneMeType
	OpenstackAccount             ManagementZoneMeType
	AppmonSystemProfile          ManagementZoneMeType
	EsxiHost                     ManagementZoneMeType
	CustomDevice                 ManagementZoneMeType
	CloudApplication             ManagementZoneMeType
	AwsRelationalDatabaseService ManagementZoneMeType
	ExternalMonitor              ManagementZoneMeType
	EnterpriseApplication        ManagementZoneMeType
	HttpMonitor                  ManagementZoneMeType
	ProcessGroup                 ManagementZoneMeType
	Host                         ManagementZoneMeType
	Queue                        ManagementZoneMeType
	AwsClassicLoadBalancer       ManagementZoneMeType
	AwsApplicationLoadBalancer   ManagementZoneMeType
	CustomApplication            ManagementZoneMeType
	AwsNetworkLoadBalancer       ManagementZoneMeType
	BrowserMonitor               ManagementZoneMeType
	AppmonServer                 ManagementZoneMeType
	CustomDeviceGroup            ManagementZoneMeType
	KubernetesCluster            ManagementZoneMeType
	HostGroup                    ManagementZoneMeType
	WebApplication               ManagementZoneMeType
	AwsAccount                   ManagementZoneMeType
	KubernetesService            ManagementZoneMeType
	AwsAutoScalingGroup          ManagementZoneMeType
	CloudApplicationNamespace    ManagementZoneMeType
}{
	ManagementZoneMeType("AZURE"),
	ManagementZoneMeType("DATA_CENTER_SERVICE"),
	ManagementZoneMeType("SERVICE"),
	ManagementZoneMeType("MOBILE_APPLICATION"),
	ManagementZoneMeType("CLOUD_FOUNDRY_FOUNDATION"),
	ManagementZoneMeType("OPENSTACK_ACCOUNT"),
	ManagementZoneMeType("APPMON_SYSTEM_PROFILE"),
	ManagementZoneMeType("ESXI_HOST"),
	ManagementZoneMeType("CUSTOM_DEVICE"),
	ManagementZoneMeType("CLOUD_APPLICATION"),
	ManagementZoneMeType("AWS_RELATIONAL_DATABASE_SERVICE"),
	ManagementZoneMeType("EXTERNAL_MONITOR"),
	ManagementZoneMeType("ENTERPRISE_APPLICATION"),
	ManagementZoneMeType("HTTP_MONITOR"),
	ManagementZoneMeType("PROCESS_GROUP"),
	ManagementZoneMeType("HOST"),
	ManagementZoneMeType("QUEUE"),
	ManagementZoneMeType("AWS_CLASSIC_LOAD_BALANCER"),
	ManagementZoneMeType("AWS_APPLICATION_LOAD_BALANCER"),
	ManagementZoneMeType("CUSTOM_APPLICATION"),
	ManagementZoneMeType("AWS_NETWORK_LOAD_BALANCER"),
	ManagementZoneMeType("BROWSER_MONITOR"),
	ManagementZoneMeType("APPMON_SERVER"),
	ManagementZoneMeType("CUSTOM_DEVICE_GROUP"),
	ManagementZoneMeType("KUBERNETES_CLUSTER"),
	ManagementZoneMeType("HOST_GROUP"),
	ManagementZoneMeType("WEB_APPLICATION"),
	ManagementZoneMeType("AWS_ACCOUNT"),
	ManagementZoneMeType("KUBERNETES_SERVICE"),
	ManagementZoneMeType("AWS_AUTO_SCALING_GROUP"),
	ManagementZoneMeType("CLOUD_APPLICATION_NAMESPACE"),
}

type RuleType string

var RuleTypes = struct {
	Me        RuleType
	Dimension RuleType
	Selector  RuleType
}{
	RuleType("ME"),
	RuleType("DIMENSION"),
	RuleType("SELECTOR"),
}

type DimensionConditionType string

var DimensionConditionTypes = struct {
	Dimension   DimensionConditionType
	LogFileName DimensionConditionType
	MetricKey   DimensionConditionType
}{
	DimensionConditionType("DIMENSION"),
	DimensionConditionType("LOG_FILE_NAME"),
	DimensionConditionType("METRIC_KEY"),
}

type Attribute string

var Attributes = struct {
	GeolocationSiteName                       Attribute
	HostGroupName                             Attribute
	DockerFullImageName                       Attribute
	OpenstackAvailabilityZoneName             Attribute
	HostAixLogicalCpuCount                    Attribute
	ServiceAkkaActorSystem                    Attribute
	ServiceWebServiceNamespace                Attribute
	OpenstackVmSecurityGroup                  Attribute
	DataCenterServiceName                     Attribute
	HostOneagentCustomHostName                Attribute
	HostName                                  Attribute
	GoogleComputeInstanceProject              Attribute
	DataCenterServiceTags                     Attribute
	OpenstackProjectName                      Attribute
	ProcessGroupTechnologyVersion             Attribute
	GoogleComputeInstanceMachineType          Attribute
	DockerImageVersion                        Attribute
	ProcessGroupDetectedName                  Attribute
	AzureSubscriptionUuid                     Attribute
	DockerStrippedImageName                   Attribute
	ProcessGroupTechnology                    Attribute
	WebApplicationTags                        Attribute
	ProcessGroupAzureSiteName                 Attribute
	ServiceType                               Attribute
	VmwareVmName                              Attribute
	ServiceRemoteServiceName                  Attribute
	AzureEntityName                           Attribute
	VmwareDatacenterName                      Attribute
	AzureSubscriptionName                     Attribute
	AwsRelationalDatabaseServiceName          Attribute
	HostBoshAvailabilityZone                  Attribute
	HostAixSimultaneousThreads                Attribute
	MobileApplicationPlatform                 Attribute
	KubernetesNodeName                        Attribute
	AwsRelationalDatabaseServiceTags          Attribute
	ProcessGroupName                          Attribute
	OpenstackAccountProjectName               Attribute
	HttpMonitorTags                           Attribute
	HostGroupId                               Attribute
	AzureTenantName                           Attribute
	Ec2InstanceAwsInstanceType                Attribute
	EnterpriseApplicationMetadata             Attribute
	DataCenterServiceDecoderType              Attribute
	ServiceWebServiceName                     Attribute
	AppmonSystemProfileName                   Attribute
	ExternalMonitorTags                       Attribute
	AzureMgmtGroupName                        Attribute
	ServicePublicDomainName                   Attribute
	Ec2InstanceAmiId                          Attribute
	EnterpriseApplicationTags                 Attribute
	GoogleCloudPlatformZoneName               Attribute
	HostKubernetesLabels                      Attribute
	HostAwsNameTag                            Attribute
	WebApplicationType                        Attribute
	HostArchitecture                          Attribute
	AwsAutoScalingGroupTags                   Attribute
	ServiceWebServerEndpoint                  Attribute
	HostBitness                               Attribute
	ServicePort                               Attribute
	EnterpriseApplicationPort                 Attribute
	ServiceWebContextRoot                     Attribute
	ExternalMonitorEngineDescription          Attribute
	HostAzureWebApplicationSiteNames          Attribute
	EsxiHostHardwareModel                     Attribute
	HostBoshStemcellVersion                   Attribute
	AzureVmName                               Attribute
	ServiceWebServerName                      Attribute
	QueueVendor                               Attribute
	CustomApplicationPlatform                 Attribute
	HostDetectedName                          Attribute
	HttpMonitorName                           Attribute
	OpenstackAccountName                      Attribute
	ProcessGroupTechnologyEdition             Attribute
	HostAzureSku                              Attribute
	CustomDeviceDnsAddress                    Attribute
	HostOsVersion                             Attribute
	AwsRelationalDatabaseServicePort          Attribute
	CloudFoundryFoundationName                Attribute
	EsxiHostClusterName                       Attribute
	CloudApplicationNamespaceName             Attribute
	CustomDevicePort                          Attribute
	ServiceCtgServiceName                     Attribute
	GoogleComputeInstancePublicIpAddresses    Attribute
	EsxiHostProductVersion                    Attribute
	AwsRelationalDatabaseServiceDbName        Attribute
	HostCloudType                             Attribute
	GoogleComputeInstanceName                 Attribute
	DockerContainerName                       Attribute
	AzureTenantUuid                           Attribute
	HostTags                                  Attribute
	CustomApplicationType                     Attribute
	HostTechnology                            Attribute
	AwsAccountName                            Attribute
	CloudApplicationNamespaceLabels           Attribute
	ProcessGroupCustomMetadata                Attribute
	AzureMgmtGroupUuid                        Attribute
	ProcessGroupListenPort                    Attribute
	AzureEntityTags                           Attribute
	EsxiHostHardwareVendor                    Attribute
	OpenstackVmInstanceType                   Attribute
	HostBoshInstanceName                      Attribute
	ServiceDatabaseVendor                     Attribute
	HostIpAddress                             Attribute
	HostOsType                                Attribute
	HostLogicalCpuCores                       Attribute
	EnterpriseApplicationName                 Attribute
	OpenstackVmName                           Attribute
	Ec2InstancePrivateHostName                Attribute
	ProcessGroupId                            Attribute
	BrowserMonitorName                        Attribute
	QueueTechnology                           Attribute
	ServiceTechnologyVersion                  Attribute
	HostAzureComputeMode                      Attribute
	EsxiHostName                              Attribute
	ServiceRemoteEndpoint                     Attribute
	MobileApplicationName                     Attribute
	GoogleComputeInstanceId                   Attribute
	AwsClassicLoadBalancerName                Attribute
	WebApplicationNamePattern                 Attribute
	CloudFoundryOrgName                       Attribute
	EsxiHostTags                              Attribute
	MobileApplicationTags                     Attribute
	Ec2InstanceId                             Attribute
	ServiceWebApplicationId                   Attribute
	ServiceTechnology                         Attribute
	QueueName                                 Attribute
	AwsNetworkLoadBalancerTags                Attribute
	CustomDeviceGroupTags                     Attribute
	HostPaasType                              Attribute
	ServiceName                               Attribute
	ServiceMessagingListenerClassName         Attribute
	HostCustomMetadata                        Attribute
	AwsApplicationLoadBalancerName            Attribute
	AwsAccountId                              Attribute
	ServiceIbmCtgGatewayUrl                   Attribute
	ServiceTechnologyEdition                  Attribute
	AwsRelationalDatabaseServiceEndpoint      Attribute
	ServiceDatabaseTopology                   Attribute
	Ec2InstanceTags                           Attribute
	AwsAvailabilityZoneName                   Attribute
	ServiceDatabaseName                       Attribute
	AwsAutoScalingGroupName                   Attribute
	ExternalMonitorEngineName                 Attribute
	KubernetesClusterName                     Attribute
	CustomDeviceGroupName                     Attribute
	EnterpriseApplicationDecoderType          Attribute
	HostBoshName                              Attribute
	ServiceDatabaseHostName                   Attribute
	ProcessGroupPredefinedMetadata            Attribute
	BrowserMonitorTags                        Attribute
	EsxiHostProductName                       Attribute
	CustomApplicationTags                     Attribute
	ServiceTags                               Attribute
	HostBoshInstanceId                        Attribute
	AwsClassicLoadBalancerFrontendPorts       Attribute
	HostAixVirtualCpuCount                    Attribute
	DataCenterServiceIpAddress                Attribute
	AwsNetworkLoadBalancerName                Attribute
	HostAzureWebApplicationHostNames          Attribute
	AppmonServerName                          Attribute
	CustomDeviceMetadata                      Attribute
	AwsApplicationLoadBalancerTags            Attribute
	CustomDeviceName                          Attribute
	HostCpuCores                              Attribute
	CustomDeviceIpAddress                     Attribute
	ServiceDetectedName                       Attribute
	HostHypervisorType                        Attribute
	CloudApplicationName                      Attribute
	KubernetesServiceName                     Attribute
	ProcessGroupAzureHostName                 Attribute
	NameOfComputeNode                         Attribute
	AzureRegionName                           Attribute
	OpenstackRegionName                       Attribute
	CustomDeviceTags                          Attribute
	Ec2InstanceBeanstalkEnvName               Attribute
	DataCenterServicePort                     Attribute
	EnterpriseApplicationIpAddress            Attribute
	Ec2InstanceName                           Attribute
	AwsRelationalDatabaseServiceEngine        Attribute
	HostBoshDeploymentId                      Attribute
	AwsClassicLoadBalancerTags                Attribute
	Ec2InstanceAwsSecurityGroup               Attribute
	HostPaasMemoryLimit                       Attribute
	ServiceEsbApplicationName                 Attribute
	GoogleComputeInstanceProjectId            Attribute
	ExternalMonitorEngineType                 Attribute
	WebApplicationName                        Attribute
	CustomApplicationName                     Attribute
	AzureScaleSetName                         Attribute
	CustomDeviceTechnology                    Attribute
	DataCenterServiceMetadata                 Attribute
	AwsRelationalDatabaseServiceInstanceClass Attribute
	ProcessGroupTags                          Attribute
	Ec2InstancePublicHostName                 Attribute
	ServiceTopology                           Attribute
	CloudApplicationLabels                    Attribute
	ExternalMonitorName                       Attribute
}{
	Attribute("GEOLOCATION_SITE_NAME"),
	Attribute("HOST_GROUP_NAME"),
	Attribute("DOCKER_FULL_IMAGE_NAME"),
	Attribute("OPENSTACK_AVAILABILITY_ZONE_NAME"),
	Attribute("HOST_AIX_LOGICAL_CPU_COUNT"),
	Attribute("SERVICE_AKKA_ACTOR_SYSTEM"),
	Attribute("SERVICE_WEB_SERVICE_NAMESPACE"),
	Attribute("OPENSTACK_VM_SECURITY_GROUP"),
	Attribute("DATA_CENTER_SERVICE_NAME"),
	Attribute("HOST_ONEAGENT_CUSTOM_HOST_NAME"),
	Attribute("HOST_NAME"),
	Attribute("GOOGLE_COMPUTE_INSTANCE_PROJECT"),
	Attribute("DATA_CENTER_SERVICE_TAGS"),
	Attribute("OPENSTACK_PROJECT_NAME"),
	Attribute("PROCESS_GROUP_TECHNOLOGY_VERSION"),
	Attribute("GOOGLE_COMPUTE_INSTANCE_MACHINE_TYPE"),
	Attribute("DOCKER_IMAGE_VERSION"),
	Attribute("PROCESS_GROUP_DETECTED_NAME"),
	Attribute("AZURE_SUBSCRIPTION_UUID"),
	Attribute("DOCKER_STRIPPED_IMAGE_NAME"),
	Attribute("PROCESS_GROUP_TECHNOLOGY"),
	Attribute("WEB_APPLICATION_TAGS"),
	Attribute("PROCESS_GROUP_AZURE_SITE_NAME"),
	Attribute("SERVICE_TYPE"),
	Attribute("VMWARE_VM_NAME"),
	Attribute("SERVICE_REMOTE_SERVICE_NAME"),
	Attribute("AZURE_ENTITY_NAME"),
	Attribute("VMWARE_DATACENTER_NAME"),
	Attribute("AZURE_SUBSCRIPTION_NAME"),
	Attribute("AWS_RELATIONAL_DATABASE_SERVICE_NAME"),
	Attribute("HOST_BOSH_AVAILABILITY_ZONE"),
	Attribute("HOST_AIX_SIMULTANEOUS_THREADS"),
	Attribute("MOBILE_APPLICATION_PLATFORM"),
	Attribute("KUBERNETES_NODE_NAME"),
	Attribute("AWS_RELATIONAL_DATABASE_SERVICE_TAGS"),
	Attribute("PROCESS_GROUP_NAME"),
	Attribute("OPENSTACK_ACCOUNT_PROJECT_NAME"),
	Attribute("HTTP_MONITOR_TAGS"),
	Attribute("HOST_GROUP_ID"),
	Attribute("AZURE_TENANT_NAME"),
	Attribute("EC2_INSTANCE_AWS_INSTANCE_TYPE"),
	Attribute("ENTERPRISE_APPLICATION_METADATA"),
	Attribute("DATA_CENTER_SERVICE_DECODER_TYPE"),
	Attribute("SERVICE_WEB_SERVICE_NAME"),
	Attribute("APPMON_SYSTEM_PROFILE_NAME"),
	Attribute("EXTERNAL_MONITOR_TAGS"),
	Attribute("AZURE_MGMT_GROUP_NAME"),
	Attribute("SERVICE_PUBLIC_DOMAIN_NAME"),
	Attribute("EC2_INSTANCE_AMI_ID"),
	Attribute("ENTERPRISE_APPLICATION_TAGS"),
	Attribute("GOOGLE_CLOUD_PLATFORM_ZONE_NAME"),
	Attribute("HOST_KUBERNETES_LABELS"),
	Attribute("HOST_AWS_NAME_TAG"),
	Attribute("WEB_APPLICATION_TYPE"),
	Attribute("HOST_ARCHITECTURE"),
	Attribute("AWS_AUTO_SCALING_GROUP_TAGS"),
	Attribute("SERVICE_WEB_SERVER_ENDPOINT"),
	Attribute("HOST_BITNESS"),
	Attribute("SERVICE_PORT"),
	Attribute("ENTERPRISE_APPLICATION_PORT"),
	Attribute("SERVICE_WEB_CONTEXT_ROOT"),
	Attribute("EXTERNAL_MONITOR_ENGINE_DESCRIPTION"),
	Attribute("HOST_AZURE_WEB_APPLICATION_SITE_NAMES"),
	Attribute("ESXI_HOST_HARDWARE_MODEL"),
	Attribute("HOST_BOSH_STEMCELL_VERSION"),
	Attribute("AZURE_VM_NAME"),
	Attribute("SERVICE_WEB_SERVER_NAME"),
	Attribute("QUEUE_VENDOR"),
	Attribute("CUSTOM_APPLICATION_PLATFORM"),
	Attribute("HOST_DETECTED_NAME"),
	Attribute("HTTP_MONITOR_NAME"),
	Attribute("OPENSTACK_ACCOUNT_NAME"),
	Attribute("PROCESS_GROUP_TECHNOLOGY_EDITION"),
	Attribute("HOST_AZURE_SKU"),
	Attribute("CUSTOM_DEVICE_DNS_ADDRESS"),
	Attribute("HOST_OS_VERSION"),
	Attribute("AWS_RELATIONAL_DATABASE_SERVICE_PORT"),
	Attribute("CLOUD_FOUNDRY_FOUNDATION_NAME"),
	Attribute("ESXI_HOST_CLUSTER_NAME"),
	Attribute("CLOUD_APPLICATION_NAMESPACE_NAME"),
	Attribute("CUSTOM_DEVICE_PORT"),
	Attribute("SERVICE_CTG_SERVICE_NAME"),
	Attribute("GOOGLE_COMPUTE_INSTANCE_PUBLIC_IP_ADDRESSES"),
	Attribute("ESXI_HOST_PRODUCT_VERSION"),
	Attribute("AWS_RELATIONAL_DATABASE_SERVICE_DB_NAME"),
	Attribute("HOST_CLOUD_TYPE"),
	Attribute("GOOGLE_COMPUTE_INSTANCE_NAME"),
	Attribute("DOCKER_CONTAINER_NAME"),
	Attribute("AZURE_TENANT_UUID"),
	Attribute("HOST_TAGS"),
	Attribute("CUSTOM_APPLICATION_TYPE"),
	Attribute("HOST_TECHNOLOGY"),
	Attribute("AWS_ACCOUNT_NAME"),
	Attribute("CLOUD_APPLICATION_NAMESPACE_LABELS"),
	Attribute("PROCESS_GROUP_CUSTOM_METADATA"),
	Attribute("AZURE_MGMT_GROUP_UUID"),
	Attribute("PROCESS_GROUP_LISTEN_PORT"),
	Attribute("AZURE_ENTITY_TAGS"),
	Attribute("ESXI_HOST_HARDWARE_VENDOR"),
	Attribute("OPENSTACK_VM_INSTANCE_TYPE"),
	Attribute("HOST_BOSH_INSTANCE_NAME"),
	Attribute("SERVICE_DATABASE_VENDOR"),
	Attribute("HOST_IP_ADDRESS"),
	Attribute("HOST_OS_TYPE"),
	Attribute("HOST_LOGICAL_CPU_CORES"),
	Attribute("ENTERPRISE_APPLICATION_NAME"),
	Attribute("OPENSTACK_VM_NAME"),
	Attribute("EC2_INSTANCE_PRIVATE_HOST_NAME"),
	Attribute("PROCESS_GROUP_ID"),
	Attribute("BROWSER_MONITOR_NAME"),
	Attribute("QUEUE_TECHNOLOGY"),
	Attribute("SERVICE_TECHNOLOGY_VERSION"),
	Attribute("HOST_AZURE_COMPUTE_MODE"),
	Attribute("ESXI_HOST_NAME"),
	Attribute("SERVICE_REMOTE_ENDPOINT"),
	Attribute("MOBILE_APPLICATION_NAME"),
	Attribute("GOOGLE_COMPUTE_INSTANCE_ID"),
	Attribute("AWS_CLASSIC_LOAD_BALANCER_NAME"),
	Attribute("WEB_APPLICATION_NAME_PATTERN"),
	Attribute("CLOUD_FOUNDRY_ORG_NAME"),
	Attribute("ESXI_HOST_TAGS"),
	Attribute("MOBILE_APPLICATION_TAGS"),
	Attribute("EC2_INSTANCE_ID"),
	Attribute("SERVICE_WEB_APPLICATION_ID"),
	Attribute("SERVICE_TECHNOLOGY"),
	Attribute("QUEUE_NAME"),
	Attribute("AWS_NETWORK_LOAD_BALANCER_TAGS"),
	Attribute("CUSTOM_DEVICE_GROUP_TAGS"),
	Attribute("HOST_PAAS_TYPE"),
	Attribute("SERVICE_NAME"),
	Attribute("SERVICE_MESSAGING_LISTENER_CLASS_NAME"),
	Attribute("HOST_CUSTOM_METADATA"),
	Attribute("AWS_APPLICATION_LOAD_BALANCER_NAME"),
	Attribute("AWS_ACCOUNT_ID"),
	Attribute("SERVICE_IBM_CTG_GATEWAY_URL"),
	Attribute("SERVICE_TECHNOLOGY_EDITION"),
	Attribute("AWS_RELATIONAL_DATABASE_SERVICE_ENDPOINT"),
	Attribute("SERVICE_DATABASE_TOPOLOGY"),
	Attribute("EC2_INSTANCE_TAGS"),
	Attribute("AWS_AVAILABILITY_ZONE_NAME"),
	Attribute("SERVICE_DATABASE_NAME"),
	Attribute("AWS_AUTO_SCALING_GROUP_NAME"),
	Attribute("EXTERNAL_MONITOR_ENGINE_NAME"),
	Attribute("KUBERNETES_CLUSTER_NAME"),
	Attribute("CUSTOM_DEVICE_GROUP_NAME"),
	Attribute("ENTERPRISE_APPLICATION_DECODER_TYPE"),
	Attribute("HOST_BOSH_NAME"),
	Attribute("SERVICE_DATABASE_HOST_NAME"),
	Attribute("PROCESS_GROUP_PREDEFINED_METADATA"),
	Attribute("BROWSER_MONITOR_TAGS"),
	Attribute("ESXI_HOST_PRODUCT_NAME"),
	Attribute("CUSTOM_APPLICATION_TAGS"),
	Attribute("SERVICE_TAGS"),
	Attribute("HOST_BOSH_INSTANCE_ID"),
	Attribute("AWS_CLASSIC_LOAD_BALANCER_FRONTEND_PORTS"),
	Attribute("HOST_AIX_VIRTUAL_CPU_COUNT"),
	Attribute("DATA_CENTER_SERVICE_IP_ADDRESS"),
	Attribute("AWS_NETWORK_LOAD_BALANCER_NAME"),
	Attribute("HOST_AZURE_WEB_APPLICATION_HOST_NAMES"),
	Attribute("APPMON_SERVER_NAME"),
	Attribute("CUSTOM_DEVICE_METADATA"),
	Attribute("AWS_APPLICATION_LOAD_BALANCER_TAGS"),
	Attribute("CUSTOM_DEVICE_NAME"),
	Attribute("HOST_CPU_CORES"),
	Attribute("CUSTOM_DEVICE_IP_ADDRESS"),
	Attribute("SERVICE_DETECTED_NAME"),
	Attribute("HOST_HYPERVISOR_TYPE"),
	Attribute("CLOUD_APPLICATION_NAME"),
	Attribute("KUBERNETES_SERVICE_NAME"),
	Attribute("PROCESS_GROUP_AZURE_HOST_NAME"),
	Attribute("NAME_OF_COMPUTE_NODE"),
	Attribute("AZURE_REGION_NAME"),
	Attribute("OPENSTACK_REGION_NAME"),
	Attribute("CUSTOM_DEVICE_TAGS"),
	Attribute("EC2_INSTANCE_BEANSTALK_ENV_NAME"),
	Attribute("DATA_CENTER_SERVICE_PORT"),
	Attribute("ENTERPRISE_APPLICATION_IP_ADDRESS"),
	Attribute("EC2_INSTANCE_NAME"),
	Attribute("AWS_RELATIONAL_DATABASE_SERVICE_ENGINE"),
	Attribute("HOST_BOSH_DEPLOYMENT_ID"),
	Attribute("AWS_CLASSIC_LOAD_BALANCER_TAGS"),
	Attribute("EC2_INSTANCE_AWS_SECURITY_GROUP"),
	Attribute("HOST_PAAS_MEMORY_LIMIT"),
	Attribute("SERVICE_ESB_APPLICATION_NAME"),
	Attribute("GOOGLE_COMPUTE_INSTANCE_PROJECT_ID"),
	Attribute("EXTERNAL_MONITOR_ENGINE_TYPE"),
	Attribute("WEB_APPLICATION_NAME"),
	Attribute("CUSTOM_APPLICATION_NAME"),
	Attribute("AZURE_SCALE_SET_NAME"),
	Attribute("CUSTOM_DEVICE_TECHNOLOGY"),
	Attribute("DATA_CENTER_SERVICE_METADATA"),
	Attribute("AWS_RELATIONAL_DATABASE_SERVICE_INSTANCE_CLASS"),
	Attribute("PROCESS_GROUP_TAGS"),
	Attribute("EC2_INSTANCE_PUBLIC_HOST_NAME"),
	Attribute("SERVICE_TOPOLOGY"),
	Attribute("CLOUD_APPLICATION_LABELS"),
	Attribute("EXTERNAL_MONITOR_NAME"),
}

type DimensionOperator string

var DimensionOperators = struct {
	Equals     DimensionOperator
	BeginsWith DimensionOperator
}{
	DimensionOperator("EQUALS"),
	DimensionOperator("BEGINS_WITH"),
}

type Operator string

var Operators = struct {
	NotGreaterThan        Operator
	NotRegexMatches       Operator
	IsIpInRange           Operator
	NotIsIpInRange        Operator
	LowerThan             Operator
	TagKeyEquals          Operator
	NotEndsWith           Operator
	NotTagKeyEquals       Operator
	NotGreaterThanOrEqual Operator
	LowerThanOrEqual      Operator
	Exists                Operator
	Contains              Operator
	RegexMatches          Operator
	EndsWith              Operator
	GreaterThanOrEqual    Operator
	NotExists             Operator
	NotContains           Operator
	NotLowerThanOrEqual   Operator
	NotEquals             Operator
	Equals                Operator
	BeginsWith            Operator
	NotBeginsWith         Operator
	NotLowerThan          Operator
	GreaterThan           Operator
}{
	Operator("NOT_GREATER_THAN"),
	Operator("NOT_REGEX_MATCHES"),
	Operator("IS_IP_IN_RANGE"),
	Operator("NOT_IS_IP_IN_RANGE"),
	Operator("LOWER_THAN"),
	Operator("TAG_KEY_EQUALS"),
	Operator("NOT_ENDS_WITH"),
	Operator("NOT_TAG_KEY_EQUALS"),
	Operator("NOT_GREATER_THAN_OR_EQUAL"),
	Operator("LOWER_THAN_OR_EQUAL"),
	Operator("EXISTS"),
	Operator("CONTAINS"),
	Operator("REGEX_MATCHES"),
	Operator("ENDS_WITH"),
	Operator("GREATER_THAN_OR_EQUAL"),
	Operator("NOT_EXISTS"),
	Operator("NOT_CONTAINS"),
	Operator("NOT_LOWER_THAN_OR_EQUAL"),
	Operator("NOT_EQUALS"),
	Operator("EQUALS"),
	Operator("BEGINS_WITH"),
	Operator("NOT_BEGINS_WITH"),
	Operator("NOT_LOWER_THAN"),
	Operator("GREATER_THAN"),
}
