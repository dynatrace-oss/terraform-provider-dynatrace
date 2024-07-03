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

type Attribute string

var Attributes = struct {
	AppmonServerName                          Attribute
	AppmonSystemProfileName                   Attribute
	AwsAccountId                              Attribute
	AwsAccountName                            Attribute
	AwsApplicationLoadBalancerName            Attribute
	AwsApplicationLoadBalancerTags            Attribute
	AwsAutoScalingGroupName                   Attribute
	AwsAutoScalingGroupTags                   Attribute
	AwsAvailabilityZoneName                   Attribute
	AwsClassicLoadBalancerFrontendPorts       Attribute
	AwsClassicLoadBalancerName                Attribute
	AwsClassicLoadBalancerTags                Attribute
	AwsNetworkLoadBalancerName                Attribute
	AwsNetworkLoadBalancerTags                Attribute
	AwsRelationalDatabaseServiceDbName        Attribute
	AwsRelationalDatabaseServiceEndpoint      Attribute
	AwsRelationalDatabaseServiceEngine        Attribute
	AwsRelationalDatabaseServiceInstanceClass Attribute
	AwsRelationalDatabaseServiceName          Attribute
	AwsRelationalDatabaseServicePort          Attribute
	AwsRelationalDatabaseServiceTags          Attribute
	AzureEntityName                           Attribute
	AzureEntityTags                           Attribute
	AzureMgmtGroupName                        Attribute
	AzureMgmtGroupUuid                        Attribute
	AzureRegionName                           Attribute
	AzureScaleSetName                         Attribute
	AzureSubscriptionName                     Attribute
	AzureSubscriptionUuid                     Attribute
	AzureTenantName                           Attribute
	AzureTenantUuid                           Attribute
	AzureVmName                               Attribute
	BrowserMonitorName                        Attribute
	BrowserMonitorTags                        Attribute
	CloudApplicationLabels                    Attribute
	CloudApplicationName                      Attribute
	CloudApplicationNamespaceLabels           Attribute
	CloudApplicationNamespaceName             Attribute
	CloudFoundryFoundationName                Attribute
	CloudFoundryOrgName                       Attribute
	CustomApplicationName                     Attribute
	CustomApplicationPlatform                 Attribute
	CustomApplicationTags                     Attribute
	CustomApplicationType                     Attribute
	CustomDeviceDnsAddress                    Attribute
	CustomDeviceGroupName                     Attribute
	CustomDeviceGroupTags                     Attribute
	CustomDeviceIpAddress                     Attribute
	CustomDeviceMetadata                      Attribute
	CustomDeviceName                          Attribute
	CustomDevicePort                          Attribute
	CustomDeviceTags                          Attribute
	CustomDeviceTechnology                    Attribute
	DataCenterServiceDecoderType              Attribute
	DataCenterServiceIpAddress                Attribute
	DataCenterServiceMetadata                 Attribute
	DataCenterServiceName                     Attribute
	DataCenterServicePort                     Attribute
	DataCenterServiceTags                     Attribute
	DockerContainerName                       Attribute
	DockerFullImageName                       Attribute
	DockerImageVersion                        Attribute
	Ec2InstanceAmiId                          Attribute
	Ec2InstanceAwsInstanceType                Attribute
	Ec2InstanceAwsSecurityGroup               Attribute
	Ec2InstanceBeanstalkEnvName               Attribute
	Ec2InstanceId                             Attribute
	Ec2InstanceName                           Attribute
	Ec2InstancePrivateHostName                Attribute
	Ec2InstancePublicHostName                 Attribute
	Ec2InstanceTags                           Attribute
	EnterpriseApplicationDecoderType          Attribute
	EnterpriseApplicationIpAddress            Attribute
	EnterpriseApplicationMetadata             Attribute
	EnterpriseApplicationName                 Attribute
	EnterpriseApplicationPort                 Attribute
	EnterpriseApplicationTags                 Attribute
	EsxiHostClusterName                       Attribute
	EsxiHostHardwareModel                     Attribute
	EsxiHostHardwareVendor                    Attribute
	EsxiHostName                              Attribute
	EsxiHostProductName                       Attribute
	EsxiHostProductVersion                    Attribute
	EsxiHostTags                              Attribute
	ExternalMonitorEngineDescription          Attribute
	ExternalMonitorEngineName                 Attribute
	ExternalMonitorEngineType                 Attribute
	ExternalMonitorName                       Attribute
	ExternalMonitorTags                       Attribute
	GeolocationSiteName                       Attribute
	GoogleCloudPlatformZoneName               Attribute
	GoogleComputeInstanceId                   Attribute
	GoogleComputeInstanceMachineType          Attribute
	GoogleComputeInstanceName                 Attribute
	GoogleComputeInstanceProject              Attribute
	GoogleComputeInstanceProjectId            Attribute
	GoogleComputeInstancePublicIpAddresses    Attribute
	HostAixLogicalCpuCount                    Attribute
	HostAixSimultaneousThreads                Attribute
	HostAixVirtualCpuCount                    Attribute
	HostArchitecture                          Attribute
	HostAwsNameTag                            Attribute
	HostAzureComputeMode                      Attribute
	HostAzureSku                              Attribute
	HostAzureWebApplicationHostNames          Attribute
	HostAzureWebApplicationSiteNames          Attribute
	HostBitness                               Attribute
	HostBoshAvailabilityZone                  Attribute
	HostBoshDeploymentId                      Attribute
	HostBoshInstanceId                        Attribute
	HostBoshInstanceName                      Attribute
	HostBoshName                              Attribute
	HostBoshStemcellVersion                   Attribute
	HostCloudType                             Attribute
	HostCpuCores                              Attribute
	HostCustomMetadata                        Attribute
	HostDetectedName                          Attribute
	HostGroupId                               Attribute
	HostGroupName                             Attribute
	HostHypervisorType                        Attribute
	HostIpAddress                             Attribute
	HostKubernetesLabels                      Attribute
	HostLogicalCpuCores                       Attribute
	HostName                                  Attribute
	HostOneagentCustomHostName                Attribute
	HostOsType                                Attribute
	HostOsVersion                             Attribute
	HostPaasMemoryLimit                       Attribute
	HostPaasType                              Attribute
	HostTags                                  Attribute
	HostTechnology                            Attribute
	HttpMonitorName                           Attribute
	HttpMonitorTags                           Attribute
	KubernetesClusterName                     Attribute
	KubernetesNodeName                        Attribute
	KubernetesServiceName                     Attribute
	MobileApplicationName                     Attribute
	MobileApplicationPlatform                 Attribute
	MobileApplicationTags                     Attribute
	NameOfComputeNode                         Attribute
	NetworkAvailabilityMonitorName            Attribute
	NetworkAvailabilityMonitorTags            Attribute
	OpenstackAccountName                      Attribute
	OpenstackAccountProjectName               Attribute
	OpenstackAvailabilityZoneName             Attribute
	OpenstackProjectName                      Attribute
	OpenstackRegionName                       Attribute
	OpenstackVmInstanceType                   Attribute
	OpenstackVmName                           Attribute
	OpenstackVmSecurityGroup                  Attribute
	ProcessGroupAzureHostName                 Attribute
	ProcessGroupAzureSiteName                 Attribute
	ProcessGroupCustomMetadata                Attribute
	ProcessGroupDetectedName                  Attribute
	ProcessGroupId                            Attribute
	ProcessGroupListenPort                    Attribute
	ProcessGroupName                          Attribute
	ProcessGroupPredefinedMetadata            Attribute
	ProcessGroupTags                          Attribute
	ProcessGroupTechnology                    Attribute
	ProcessGroupTechnologyEdition             Attribute
	ProcessGroupTechnologyVersion             Attribute
	QueueName                                 Attribute
	QueueTechnology                           Attribute
	QueueVendor                               Attribute
	ServiceAkkaActorSystem                    Attribute
	ServiceCtgServiceName                     Attribute
	ServiceDatabaseHostName                   Attribute
	ServiceDatabaseName                       Attribute
	ServiceDatabaseTopology                   Attribute
	ServiceDatabaseVendor                     Attribute
	ServiceDetectedName                       Attribute
	ServiceEsbApplicationName                 Attribute
	ServiceIbmCtgGatewayUrl                   Attribute
	ServiceMessagingListenerClassName         Attribute
	ServiceName                               Attribute
	ServicePort                               Attribute
	ServicePublicDomainName                   Attribute
	ServiceRemoteEndpoint                     Attribute
	ServiceRemoteServiceName                  Attribute
	ServiceTags                               Attribute
	ServiceTechnology                         Attribute
	ServiceTechnologyEdition                  Attribute
	ServiceTechnologyVersion                  Attribute
	ServiceTopology                           Attribute
	ServiceType                               Attribute
	ServiceWebApplicationId                   Attribute
	ServiceWebContextRoot                     Attribute
	ServiceWebServerEndpoint                  Attribute
	ServiceWebServerName                      Attribute
	ServiceWebServiceName                     Attribute
	ServiceWebServiceNamespace                Attribute
	VmwareDatacenterName                      Attribute
	VmwareVmName                              Attribute
	WebApplicationName                        Attribute
	WebApplicationNamePattern                 Attribute
	WebApplicationTags                        Attribute
	WebApplicationType                        Attribute
}{
	"APPMON_SERVER_NAME",
	"APPMON_SYSTEM_PROFILE_NAME",
	"AWS_ACCOUNT_ID",
	"AWS_ACCOUNT_NAME",
	"AWS_APPLICATION_LOAD_BALANCER_NAME",
	"AWS_APPLICATION_LOAD_BALANCER_TAGS",
	"AWS_AUTO_SCALING_GROUP_NAME",
	"AWS_AUTO_SCALING_GROUP_TAGS",
	"AWS_AVAILABILITY_ZONE_NAME",
	"AWS_CLASSIC_LOAD_BALANCER_FRONTEND_PORTS",
	"AWS_CLASSIC_LOAD_BALANCER_NAME",
	"AWS_CLASSIC_LOAD_BALANCER_TAGS",
	"AWS_NETWORK_LOAD_BALANCER_NAME",
	"AWS_NETWORK_LOAD_BALANCER_TAGS",
	"AWS_RELATIONAL_DATABASE_SERVICE_DB_NAME",
	"AWS_RELATIONAL_DATABASE_SERVICE_ENDPOINT",
	"AWS_RELATIONAL_DATABASE_SERVICE_ENGINE",
	"AWS_RELATIONAL_DATABASE_SERVICE_INSTANCE_CLASS",
	"AWS_RELATIONAL_DATABASE_SERVICE_NAME",
	"AWS_RELATIONAL_DATABASE_SERVICE_PORT",
	"AWS_RELATIONAL_DATABASE_SERVICE_TAGS",
	"AZURE_ENTITY_NAME",
	"AZURE_ENTITY_TAGS",
	"AZURE_MGMT_GROUP_NAME",
	"AZURE_MGMT_GROUP_UUID",
	"AZURE_REGION_NAME",
	"AZURE_SCALE_SET_NAME",
	"AZURE_SUBSCRIPTION_NAME",
	"AZURE_SUBSCRIPTION_UUID",
	"AZURE_TENANT_NAME",
	"AZURE_TENANT_UUID",
	"AZURE_VM_NAME",
	"BROWSER_MONITOR_NAME",
	"BROWSER_MONITOR_TAGS",
	"CLOUD_APPLICATION_LABELS",
	"CLOUD_APPLICATION_NAME",
	"CLOUD_APPLICATION_NAMESPACE_LABELS",
	"CLOUD_APPLICATION_NAMESPACE_NAME",
	"CLOUD_FOUNDRY_FOUNDATION_NAME",
	"CLOUD_FOUNDRY_ORG_NAME",
	"CUSTOM_APPLICATION_NAME",
	"CUSTOM_APPLICATION_PLATFORM",
	"CUSTOM_APPLICATION_TAGS",
	"CUSTOM_APPLICATION_TYPE",
	"CUSTOM_DEVICE_DNS_ADDRESS",
	"CUSTOM_DEVICE_GROUP_NAME",
	"CUSTOM_DEVICE_GROUP_TAGS",
	"CUSTOM_DEVICE_IP_ADDRESS",
	"CUSTOM_DEVICE_METADATA",
	"CUSTOM_DEVICE_NAME",
	"CUSTOM_DEVICE_PORT",
	"CUSTOM_DEVICE_TAGS",
	"CUSTOM_DEVICE_TECHNOLOGY",
	"DATA_CENTER_SERVICE_DECODER_TYPE",
	"DATA_CENTER_SERVICE_IP_ADDRESS",
	"DATA_CENTER_SERVICE_METADATA",
	"DATA_CENTER_SERVICE_NAME",
	"DATA_CENTER_SERVICE_PORT",
	"DATA_CENTER_SERVICE_TAGS",
	"DOCKER_CONTAINER_NAME",
	"DOCKER_FULL_IMAGE_NAME",
	"DOCKER_IMAGE_VERSION",
	"EC2_INSTANCE_AMI_ID",
	"EC2_INSTANCE_AWS_INSTANCE_TYPE",
	"EC2_INSTANCE_AWS_SECURITY_GROUP",
	"EC2_INSTANCE_BEANSTALK_ENV_NAME",
	"EC2_INSTANCE_ID",
	"EC2_INSTANCE_NAME",
	"EC2_INSTANCE_PRIVATE_HOST_NAME",
	"EC2_INSTANCE_PUBLIC_HOST_NAME",
	"EC2_INSTANCE_TAGS",
	"ENTERPRISE_APPLICATION_DECODER_TYPE",
	"ENTERPRISE_APPLICATION_IP_ADDRESS",
	"ENTERPRISE_APPLICATION_METADATA",
	"ENTERPRISE_APPLICATION_NAME",
	"ENTERPRISE_APPLICATION_PORT",
	"ENTERPRISE_APPLICATION_TAGS",
	"ESXI_HOST_CLUSTER_NAME",
	"ESXI_HOST_HARDWARE_MODEL",
	"ESXI_HOST_HARDWARE_VENDOR",
	"ESXI_HOST_NAME",
	"ESXI_HOST_PRODUCT_NAME",
	"ESXI_HOST_PRODUCT_VERSION",
	"ESXI_HOST_TAGS",
	"EXTERNAL_MONITOR_ENGINE_DESCRIPTION",
	"EXTERNAL_MONITOR_ENGINE_NAME",
	"EXTERNAL_MONITOR_ENGINE_TYPE",
	"EXTERNAL_MONITOR_NAME",
	"EXTERNAL_MONITOR_TAGS",
	"GEOLOCATION_SITE_NAME",
	"GOOGLE_CLOUD_PLATFORM_ZONE_NAME",
	"GOOGLE_COMPUTE_INSTANCE_ID",
	"GOOGLE_COMPUTE_INSTANCE_MACHINE_TYPE",
	"GOOGLE_COMPUTE_INSTANCE_NAME",
	"GOOGLE_COMPUTE_INSTANCE_PROJECT",
	"GOOGLE_COMPUTE_INSTANCE_PROJECT_ID",
	"GOOGLE_COMPUTE_INSTANCE_PUBLIC_IP_ADDRESSES",
	"HOST_AIX_LOGICAL_CPU_COUNT",
	"HOST_AIX_SIMULTANEOUS_THREADS",
	"HOST_AIX_VIRTUAL_CPU_COUNT",
	"HOST_ARCHITECTURE",
	"HOST_AWS_NAME_TAG",
	"HOST_AZURE_COMPUTE_MODE",
	"HOST_AZURE_SKU",
	"HOST_AZURE_WEB_APPLICATION_HOST_NAMES",
	"HOST_AZURE_WEB_APPLICATION_SITE_NAMES",
	"HOST_BITNESS",
	"HOST_BOSH_AVAILABILITY_ZONE",
	"HOST_BOSH_DEPLOYMENT_ID",
	"HOST_BOSH_INSTANCE_ID",
	"HOST_BOSH_INSTANCE_NAME",
	"HOST_BOSH_NAME",
	"HOST_BOSH_STEMCELL_VERSION",
	"HOST_CLOUD_TYPE",
	"HOST_CPU_CORES",
	"HOST_CUSTOM_METADATA",
	"HOST_DETECTED_NAME",
	"HOST_GROUP_ID",
	"HOST_GROUP_NAME",
	"HOST_HYPERVISOR_TYPE",
	"HOST_IP_ADDRESS",
	"HOST_KUBERNETES_LABELS",
	"HOST_LOGICAL_CPU_CORES",
	"HOST_NAME",
	"HOST_ONEAGENT_CUSTOM_HOST_NAME",
	"HOST_OS_TYPE",
	"HOST_OS_VERSION",
	"HOST_PAAS_MEMORY_LIMIT",
	"HOST_PAAS_TYPE",
	"HOST_TAGS",
	"HOST_TECHNOLOGY",
	"HTTP_MONITOR_NAME",
	"HTTP_MONITOR_TAGS",
	"KUBERNETES_CLUSTER_NAME",
	"KUBERNETES_NODE_NAME",
	"KUBERNETES_SERVICE_NAME",
	"MOBILE_APPLICATION_NAME",
	"MOBILE_APPLICATION_PLATFORM",
	"MOBILE_APPLICATION_TAGS",
	"NAME_OF_COMPUTE_NODE",
	"NETWORK_AVAILABILITY_MONITOR_NAME",
	"NETWORK_AVAILABILITY_MONITOR_TAGS",
	"OPENSTACK_ACCOUNT_NAME",
	"OPENSTACK_ACCOUNT_PROJECT_NAME",
	"OPENSTACK_AVAILABILITY_ZONE_NAME",
	"OPENSTACK_PROJECT_NAME",
	"OPENSTACK_REGION_NAME",
	"OPENSTACK_VM_INSTANCE_TYPE",
	"OPENSTACK_VM_NAME",
	"OPENSTACK_VM_SECURITY_GROUP",
	"PROCESS_GROUP_AZURE_HOST_NAME",
	"PROCESS_GROUP_AZURE_SITE_NAME",
	"PROCESS_GROUP_CUSTOM_METADATA",
	"PROCESS_GROUP_DETECTED_NAME",
	"PROCESS_GROUP_ID",
	"PROCESS_GROUP_LISTEN_PORT",
	"PROCESS_GROUP_NAME",
	"PROCESS_GROUP_PREDEFINED_METADATA",
	"PROCESS_GROUP_TAGS",
	"PROCESS_GROUP_TECHNOLOGY",
	"PROCESS_GROUP_TECHNOLOGY_EDITION",
	"PROCESS_GROUP_TECHNOLOGY_VERSION",
	"QUEUE_NAME",
	"QUEUE_TECHNOLOGY",
	"QUEUE_VENDOR",
	"SERVICE_AKKA_ACTOR_SYSTEM",
	"SERVICE_CTG_SERVICE_NAME",
	"SERVICE_DATABASE_HOST_NAME",
	"SERVICE_DATABASE_NAME",
	"SERVICE_DATABASE_TOPOLOGY",
	"SERVICE_DATABASE_VENDOR",
	"SERVICE_DETECTED_NAME",
	"SERVICE_ESB_APPLICATION_NAME",
	"SERVICE_IBM_CTG_GATEWAY_URL",
	"SERVICE_MESSAGING_LISTENER_CLASS_NAME",
	"SERVICE_NAME",
	"SERVICE_PORT",
	"SERVICE_PUBLIC_DOMAIN_NAME",
	"SERVICE_REMOTE_ENDPOINT",
	"SERVICE_REMOTE_SERVICE_NAME",
	"SERVICE_TAGS",
	"SERVICE_TECHNOLOGY",
	"SERVICE_TECHNOLOGY_EDITION",
	"SERVICE_TECHNOLOGY_VERSION",
	"SERVICE_TOPOLOGY",
	"SERVICE_TYPE",
	"SERVICE_WEB_APPLICATION_ID",
	"SERVICE_WEB_CONTEXT_ROOT",
	"SERVICE_WEB_SERVER_ENDPOINT",
	"SERVICE_WEB_SERVER_NAME",
	"SERVICE_WEB_SERVICE_NAME",
	"SERVICE_WEB_SERVICE_NAMESPACE",
	"VMWARE_DATACENTER_NAME",
	"VMWARE_VM_NAME",
	"WEB_APPLICATION_NAME",
	"WEB_APPLICATION_NAME_PATTERN",
	"WEB_APPLICATION_TAGS",
	"WEB_APPLICATION_TYPE",
}

type DimensionConditionType string

var DimensionConditionTypes = struct {
	Dimension   DimensionConditionType
	LogFileName DimensionConditionType
	MetricKey   DimensionConditionType
}{
	"DIMENSION",
	"LOG_FILE_NAME",
	"METRIC_KEY",
}

type DimensionOperator string

var DimensionOperators = struct {
	BeginsWith DimensionOperator
	Equals     DimensionOperator
}{
	"BEGINS_WITH",
	"EQUALS",
}

type DimensionType string

var DimensionTypes = struct {
	Any    DimensionType
	Log    DimensionType
	Metric DimensionType
}{
	"ANY",
	"LOG",
	"METRIC",
}

type ManagementZoneMeType string

var ManagementZoneMeTypes = struct {
	AppmonServer                 ManagementZoneMeType
	AppmonSystemProfile          ManagementZoneMeType
	AwsAccount                   ManagementZoneMeType
	AwsApplicationLoadBalancer   ManagementZoneMeType
	AwsAutoScalingGroup          ManagementZoneMeType
	AwsClassicLoadBalancer       ManagementZoneMeType
	AwsNetworkLoadBalancer       ManagementZoneMeType
	AwsRelationalDatabaseService ManagementZoneMeType
	Azure                        ManagementZoneMeType
	BrowserMonitor               ManagementZoneMeType
	CloudApplication             ManagementZoneMeType
	CloudApplicationNamespace    ManagementZoneMeType
	CloudFoundryFoundation       ManagementZoneMeType
	CustomApplication            ManagementZoneMeType
	CustomDevice                 ManagementZoneMeType
	CustomDeviceGroup            ManagementZoneMeType
	DataCenterService            ManagementZoneMeType
	EnterpriseApplication        ManagementZoneMeType
	EsxiHost                     ManagementZoneMeType
	ExternalMonitor              ManagementZoneMeType
	Host                         ManagementZoneMeType
	HostGroup                    ManagementZoneMeType
	HttpMonitor                  ManagementZoneMeType
	KubernetesCluster            ManagementZoneMeType
	KubernetesService            ManagementZoneMeType
	MobileApplication            ManagementZoneMeType
	NetworkAvailabilityMonitor   ManagementZoneMeType
	OpenstackAccount             ManagementZoneMeType
	ProcessGroup                 ManagementZoneMeType
	Queue                        ManagementZoneMeType
	Service                      ManagementZoneMeType
	WebApplication               ManagementZoneMeType
}{
	"APPMON_SERVER",
	"APPMON_SYSTEM_PROFILE",
	"AWS_ACCOUNT",
	"AWS_APPLICATION_LOAD_BALANCER",
	"AWS_AUTO_SCALING_GROUP",
	"AWS_CLASSIC_LOAD_BALANCER",
	"AWS_NETWORK_LOAD_BALANCER",
	"AWS_RELATIONAL_DATABASE_SERVICE",
	"AZURE",
	"BROWSER_MONITOR",
	"CLOUD_APPLICATION",
	"CLOUD_APPLICATION_NAMESPACE",
	"CLOUD_FOUNDRY_FOUNDATION",
	"CUSTOM_APPLICATION",
	"CUSTOM_DEVICE",
	"CUSTOM_DEVICE_GROUP",
	"DATA_CENTER_SERVICE",
	"ENTERPRISE_APPLICATION",
	"ESXI_HOST",
	"EXTERNAL_MONITOR",
	"HOST",
	"HOST_GROUP",
	"HTTP_MONITOR",
	"KUBERNETES_CLUSTER",
	"KUBERNETES_SERVICE",
	"MOBILE_APPLICATION",
	"NETWORK_AVAILABILITY_MONITOR",
	"OPENSTACK_ACCOUNT",
	"PROCESS_GROUP",
	"QUEUE",
	"SERVICE",
	"WEB_APPLICATION",
}

type Operator string

var Operators = struct {
	BeginsWith            Operator
	Contains              Operator
	EndsWith              Operator
	Equals                Operator
	Exists                Operator
	GreaterThan           Operator
	GreaterThanOrEqual    Operator
	IsIpInRange           Operator
	LowerThan             Operator
	LowerThanOrEqual      Operator
	NotBeginsWith         Operator
	NotContains           Operator
	NotEndsWith           Operator
	NotEquals             Operator
	NotExists             Operator
	NotGreaterThan        Operator
	NotGreaterThanOrEqual Operator
	NotIsIpInRange        Operator
	NotLowerThan          Operator
	NotLowerThanOrEqual   Operator
	NotRegexMatches       Operator
	NotTagKeyEquals       Operator
	RegexMatches          Operator
	TagKeyEquals          Operator
}{
	"BEGINS_WITH",
	"CONTAINS",
	"ENDS_WITH",
	"EQUALS",
	"EXISTS",
	"GREATER_THAN",
	"GREATER_THAN_OR_EQUAL",
	"IS_IP_IN_RANGE",
	"LOWER_THAN",
	"LOWER_THAN_OR_EQUAL",
	"NOT_BEGINS_WITH",
	"NOT_CONTAINS",
	"NOT_ENDS_WITH",
	"NOT_EQUALS",
	"NOT_EXISTS",
	"NOT_GREATER_THAN",
	"NOT_GREATER_THAN_OR_EQUAL",
	"NOT_IS_IP_IN_RANGE",
	"NOT_LOWER_THAN",
	"NOT_LOWER_THAN_OR_EQUAL",
	"NOT_REGEX_MATCHES",
	"NOT_TAG_KEY_EQUALS",
	"REGEX_MATCHES",
	"TAG_KEY_EQUALS",
}

type RuleType string

var RuleTypes = struct {
	Dimension RuleType
	Me        RuleType
	Selector  RuleType
}{
	"DIMENSION",
	"ME",
	"SELECTOR",
}
