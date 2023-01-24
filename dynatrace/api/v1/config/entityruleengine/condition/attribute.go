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

package condition

// Attribute The attribute to be used for comparision.
type Attribute string

// Attributes offers the known enum values
var Attributes = struct {
	AppMonServerName                          Attribute
	AppMonSystemProfileName                   Attribute
	AWSAccountID                              Attribute
	AWSAccountName                            Attribute
	AWSApplicationLoadBalancerName            Attribute
	AWSApplicationLoadBalancerTags            Attribute
	AWSAutoScalingGroupName                   Attribute
	AWSAutoScalingGroupTags                   Attribute
	AWSAvailabilityZoneName                   Attribute
	AWSClassicLoadBalancerFrontendPorts       Attribute
	AWSClassicLoadBalancerName                Attribute
	AWSClassicLoadBalancerTags                Attribute
	AWSNetworkLoadBalancerName                Attribute
	AWSNetworkLoadBalancerTags                Attribute
	AWSRelationalDatabaseServiceDBName        Attribute
	AWSRelationalDatabaseServiceEndpoint      Attribute
	AWSRelationalDatabaseServiceEngine        Attribute
	AWSRelationalDatabaseServiceInstanceClass Attribute
	AWSRelationalDatabaseServiceName          Attribute
	AWSRelationalDatabaseServicePort          Attribute
	AWSRelationalDatabaseServiceTags          Attribute
	AzureEntityName                           Attribute
	AzureEntityTags                           Attribute
	AzureMgmtGroupName                        Attribute
	AzureMgmtGroupUUID                        Attribute
	AzureRegionName                           Attribute
	AzureScaleSetName                         Attribute
	AzureSubscriptionName                     Attribute
	AzureSubscriptionUUID                     Attribute
	AzureTenantName                           Attribute
	AzureTenantUUID                           Attribute
	AzureVMName                               Attribute
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
	CustomDeviceDetectedName                  Attribute
	CustomDeviceDNSAddress                    Attribute
	CustomDeviceGroupName                     Attribute
	CustomDeviceGroupTags                     Attribute
	CustomDeviceIPAddress                     Attribute
	CustomDeviceMetadata                      Attribute
	CustomDeviceName                          Attribute
	CustomDevicePort                          Attribute
	CustomDeviceTags                          Attribute
	CustomDeviceTechnology                    Attribute
	DataCenterServiceDecoderType              Attribute
	DataCenterServiceIPAddress                Attribute
	DataCenterServiceMetadata                 Attribute
	DataCenterServiceName                     Attribute
	DataCenterServicePort                     Attribute
	DataCenterServiceTags                     Attribute
	DockerContainerName                       Attribute
	DockerFullImageName                       Attribute
	DockerImageVersion                        Attribute
	DockerStrippedImageName                   Attribute
	EC2InstanceAmiID                          Attribute
	EC2InstanceAWSInstanceType                Attribute
	EC2InstanceAWSSecurityGroup               Attribute
	EC2InstanceBeanstalkEnvName               Attribute
	EC2InstanceID                             Attribute
	EC2InstanceName                           Attribute
	EC2InstancePrivateHostName                Attribute
	EC2InstancePublicHostName                 Attribute
	EC2InstanceTags                           Attribute
	EnterpriseApplicationDecoderType          Attribute
	EnterpriseApplicationIPAddress            Attribute
	EnterpriseApplicationMetadata             Attribute
	EnterpriseApplicationName                 Attribute
	EnterpriseApplicationPort                 Attribute
	EnterpriseApplicationTags                 Attribute
	ESXIHostClusterName                       Attribute
	ESXIHostHardwareModel                     Attribute
	ESXIHostHardwareVendor                    Attribute
	ESXIHostName                              Attribute
	ESXIHostProductName                       Attribute
	ESXIHostProductVersion                    Attribute
	ESXIHostTags                              Attribute
	ExternalMonitorEngineDescription          Attribute
	ExternalMonitorEngineName                 Attribute
	ExternalMonitorEngineType                 Attribute
	ExternalMonitorName                       Attribute
	ExternalMonitorTags                       Attribute
	GeolocationSiteName                       Attribute
	GoogleCloudPlatformZoneName               Attribute
	GoogleComputeInstanceID                   Attribute
	GoogleComputeInstanceMachineType          Attribute
	GoogleComputeInstanceName                 Attribute
	GoogleComputeInstanceProject              Attribute
	GoogleComputeInstanceProjectID            Attribute
	GoogleComputeInstancePublicIPAddresses    Attribute
	HostAIXLogicalCPUCount                    Attribute
	HostAIXSimultaneousThreads                Attribute
	HostAIXVirtualCPUCount                    Attribute
	HostArchitecture                          Attribute
	HostAWSNameTag                            Attribute
	HostAzureComputeMode                      Attribute
	HostAzureSku                              Attribute
	HostAzureWebApplicationHostNames          Attribute
	HostAzureWebApplicationSiteNames          Attribute
	HostBitness                               Attribute
	HostBoshAvailabilityZone                  Attribute
	HostBoshDeploymentID                      Attribute
	HostBoshInstanceID                        Attribute
	HostBoshInstanceName                      Attribute
	HostBoshName                              Attribute
	HostBoshStemcellVersion                   Attribute
	HostCloudType                             Attribute
	HostCPUCores                              Attribute
	HostCustomMetadata                        Attribute
	HostDetectedName                          Attribute
	HostGroupID                               Attribute
	HostGroupName                             Attribute
	HostHypervisorType                        Attribute
	HostIPAddress                             Attribute
	HostKubernetesLabels                      Attribute
	HostLogicalCPUCores                       Attribute
	HostName                                  Attribute
	HostOneAgentCustomHostName                Attribute
	HostOSType                                Attribute
	HostOSVersion                             Attribute
	HostPaasMemoryLimit                       Attribute
	HostPaasType                              Attribute
	HostTags                                  Attribute
	HostTechnology                            Attribute
	HTTPMonitorName                           Attribute
	HTTPMonitorTags                           Attribute
	KubernetesClusterName                     Attribute
	KubernetesNodeName                        Attribute
	MobileApplicationName                     Attribute
	MobileApplicationPlatform                 Attribute
	MobileApplicationTags                     Attribute
	NameOfComputeNode                         Attribute
	OpenStackAccountName                      Attribute
	OpenStackAccountProjectName               Attribute
	OpenStackAvailabilityZoneName             Attribute
	OpenStackProjectName                      Attribute
	OpenStackRegionName                       Attribute
	OpenStackVMInstanceType                   Attribute
	OpenStackVMName                           Attribute
	OpenStackVMSecurityGroup                  Attribute
	ProcessGroupAzureHostName                 Attribute
	ProcessGroupAzureSiteName                 Attribute
	ProcessGroupCustomMetadata                Attribute
	ProcessGroupDetectedName                  Attribute
	ProcessGroupID                            Attribute
	ProcessGroupListenPort                    Attribute
	ProcessGroupName                          Attribute
	ProcessGroupPredefinedMetadata            Attribute
	ProcessGroupTags                          Attribute
	ProcessGroupTechnology                    Attribute
	ProcessGroupTechnologyEdition             Attribute
	ProcessGroupTechnologyVersion             Attribute
	ServiceAkkaActorSystem                    Attribute
	ServiceCTGServiceName                     Attribute
	ServiceDatabaseHostName                   Attribute
	ServiceDatabaseName                       Attribute
	ServiceDatabaseTopology                   Attribute
	ServiceDatabaseVendor                     Attribute
	ServiceDetectedName                       Attribute
	ServiceEsbApplicationName                 Attribute
	ServiceIBMCTGGatewayURL                   Attribute
	ServiceIibApplicationName                 Attribute
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
	ServiceWebApplicationID                   Attribute
	ServiceWebContextRoot                     Attribute
	ServiceWebServerEndpoint                  Attribute
	ServiceWebServerName                      Attribute
	ServiceWebServiceName                     Attribute
	ServiceWebServiceNamespace                Attribute
	VmwareDatacenterName                      Attribute
	VmwareVMName                              Attribute
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
	"CUSTOM_DEVICE_DETECTED_NAME",
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
	"DOCKER_STRIPPED_IMAGE_NAME",
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
	"MOBILE_APPLICATION_NAME",
	"MOBILE_APPLICATION_PLATFORM",
	"MOBILE_APPLICATION_TAGS",
	"NAME_OF_COMPUTE_NODE",
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
	"SERVICE_AKKA_ACTOR_SYSTEM",
	"SERVICE_CTG_SERVICE_NAME",
	"SERVICE_DATABASE_HOST_NAME",
	"SERVICE_DATABASE_NAME",
	"SERVICE_DATABASE_TOPOLOGY",
	"SERVICE_DATABASE_VENDOR",
	"SERVICE_DETECTED_NAME",
	"SERVICE_ESB_APPLICATION_NAME",
	"SERVICE_IBM_CTG_GATEWAY_URL",
	"SERVICE_IIB_APPLICATION_NAME",
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
