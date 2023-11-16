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

package customprocessmonitoring

type AgentItemName string

var AgentItemNames = struct {
	ApacheConfigPath                   AgentItemName
	ApacheSparkMasterIpAddress         AgentItemName
	AspNetCoreApplicationPath          AgentItemName
	AwsEcrAccountId                    AgentItemName
	AwsEcrRegion                       AgentItemName
	AwsEcsCluster                      AgentItemName
	AwsEcsContainername                AgentItemName
	AwsEcsFamily                       AgentItemName
	AwsEcsRevision                     AgentItemName
	AwsLambdaFunctionName              AgentItemName
	AwsRegion                          AgentItemName
	CatalinaBase                       AgentItemName
	CatalinaHome                       AgentItemName
	CloudFoundryAppName                AgentItemName
	CloudFoundryApplicationId          AgentItemName
	CloudFoundryInstanceIndex          AgentItemName
	CloudFoundrySpaceId                AgentItemName
	CloudFoundrySpaceName              AgentItemName
	ColdfusionJvmConfigFile            AgentItemName
	CommandLineArgs                    AgentItemName
	ContainerId                        AgentItemName
	ContainerImageName                 AgentItemName
	ContainerImageVersion              AgentItemName
	ContainerName                      AgentItemName
	DatasourceMonitoringConfigId       AgentItemName
	DeclarativeId                      AgentItemName
	DotnetCommand                      AgentItemName
	DotnetCommandPath                  AgentItemName
	ElasticSearchClusterName           AgentItemName
	ElasticSearchNodeName              AgentItemName
	EquinoxConfigPath                  AgentItemName
	ExeName                            AgentItemName
	ExePath                            AgentItemName
	GaeInstance                        AgentItemName
	GaeService                         AgentItemName
	GlassfishDomainName                AgentItemName
	GlassfishInstanceName              AgentItemName
	GoogleCloudProject                 AgentItemName
	HybrisBinDir                       AgentItemName
	HybrisConfigDir                    AgentItemName
	HybrisDataDir                      AgentItemName
	IbmCicsImsApplid                   AgentItemName
	IbmCicsImsJobname                  AgentItemName
	IbmCicsRegion                      AgentItemName
	IbmCtgName                         AgentItemName
	IbmImsConnect                      AgentItemName
	IbmImsControl                      AgentItemName
	IbmImsMpr                          AgentItemName
	IbmImsSoapGwName                   AgentItemName
	IibBrokerName                      AgentItemName
	IibExecutionGroupName              AgentItemName
	IisAppPool                         AgentItemName
	IisRoleName                        AgentItemName
	JavaJarFile                        AgentItemName
	JavaJarPath                        AgentItemName
	JavaMainClass                      AgentItemName
	JbossHome                          AgentItemName
	JbossMode                          AgentItemName
	JbossServerName                    AgentItemName
	KubernetesBasepodname              AgentItemName
	KubernetesContainername            AgentItemName
	KubernetesFullpodname              AgentItemName
	KubernetesNamespace                AgentItemName
	KubernetesPoduid                   AgentItemName
	MssqlInstanceName                  AgentItemName
	NodejsAppBaseDir                   AgentItemName
	NodejsAppName                      AgentItemName
	NodejsScriptName                   AgentItemName
	OracleSid                          AgentItemName
	PgIdCalcInputKeyLinkage            AgentItemName
	PhpCliScriptPath                   AgentItemName
	PhpCliWorkingDir                   AgentItemName
	Rke2Type                           AgentItemName
	RuxitClusterId                     AgentItemName
	RuxitNodeId                        AgentItemName
	ServiceName                        AgentItemName
	SoftwareagInstallRoot              AgentItemName
	SoftwareagProductpropname          AgentItemName
	SpringbootAppName                  AgentItemName
	SpringbootProfileName              AgentItemName
	SpringbootStartupClass             AgentItemName
	TibcoBusinessworksAppNodeName      AgentItemName
	TibcoBusinessworksAppSpaceName     AgentItemName
	TibcoBusinessworksCeAppName        AgentItemName
	TibcoBusinessworksCeVersion        AgentItemName
	TibcoBusinessworksDomainName       AgentItemName
	TibcoBusinessworksHome             AgentItemName
	TipcoBusinessworksPropertyFile     AgentItemName
	TipcoBusinessworksPropertyFilePath AgentItemName
	Unknown                            AgentItemName
	VarnishInstanceName                AgentItemName
	WeblogicClusterName                AgentItemName
	WeblogicDomainName                 AgentItemName
	WeblogicHome                       AgentItemName
	WeblogicName                       AgentItemName
	WebsphereCellName                  AgentItemName
	WebsphereClusterName               AgentItemName
	WebsphereLibertyServerName         AgentItemName
	WebsphereNodeName                  AgentItemName
	WebsphereServerName                AgentItemName
	ZCmVersion                         AgentItemName
}{
	"APACHE_CONFIG_PATH",
	"APACHE_SPARK_MASTER_IP_ADDRESS",
	"ASP_NET_CORE_APPLICATION_PATH",
	"AWS_ECR_ACCOUNT_ID",
	"AWS_ECR_REGION",
	"AWS_ECS_CLUSTER",
	"AWS_ECS_CONTAINERNAME",
	"AWS_ECS_FAMILY",
	"AWS_ECS_REVISION",
	"AWS_LAMBDA_FUNCTION_NAME",
	"AWS_REGION",
	"CATALINA_BASE",
	"CATALINA_HOME",
	"CLOUD_FOUNDRY_APP_NAME",
	"CLOUD_FOUNDRY_APPLICATION_ID",
	"CLOUD_FOUNDRY_INSTANCE_INDEX",
	"CLOUD_FOUNDRY_SPACE_ID",
	"CLOUD_FOUNDRY_SPACE_NAME",
	"COLDFUSION_JVM_CONFIG_FILE",
	"COMMAND_LINE_ARGS",
	"CONTAINER_ID",
	"CONTAINER_IMAGE_NAME",
	"CONTAINER_IMAGE_VERSION",
	"CONTAINER_NAME",
	"DATASOURCE_MONITORING_CONFIG_ID",
	"DECLARATIVE_ID",
	"DOTNET_COMMAND",
	"DOTNET_COMMAND_PATH",
	"ELASTIC_SEARCH_CLUSTER_NAME",
	"ELASTIC_SEARCH_NODE_NAME",
	"EQUINOX_CONFIG_PATH",
	"EXE_NAME",
	"EXE_PATH",
	"GAE_INSTANCE",
	"GAE_SERVICE",
	"GLASSFISH_DOMAIN_NAME",
	"GLASSFISH_INSTANCE_NAME",
	"GOOGLE_CLOUD_PROJECT",
	"HYBRIS_BIN_DIR",
	"HYBRIS_CONFIG_DIR",
	"HYBRIS_DATA_DIR",
	"IBM_CICS_IMS_APPLID",
	"IBM_CICS_IMS_JOBNAME",
	"IBM_CICS_REGION",
	"IBM_CTG_NAME",
	"IBM_IMS_CONNECT",
	"IBM_IMS_CONTROL",
	"IBM_IMS_MPR",
	"IBM_IMS_SOAP_GW_NAME",
	"IIB_BROKER_NAME",
	"IIB_EXECUTION_GROUP_NAME",
	"IIS_APP_POOL",
	"IIS_ROLE_NAME",
	"JAVA_JAR_FILE",
	"JAVA_JAR_PATH",
	"JAVA_MAIN_CLASS",
	"JBOSS_HOME",
	"JBOSS_MODE",
	"JBOSS_SERVER_NAME",
	"KUBERNETES_BASEPODNAME",
	"KUBERNETES_CONTAINERNAME",
	"KUBERNETES_FULLPODNAME",
	"KUBERNETES_NAMESPACE",
	"KUBERNETES_PODUID",
	"MSSQL_INSTANCE_NAME",
	"NODEJS_APP_BASE_DIR",
	"NODEJS_APP_NAME",
	"NODEJS_SCRIPT_NAME",
	"ORACLE_SID",
	"PG_ID_CALC_INPUT_KEY_LINKAGE",
	"PHP_CLI_SCRIPT_PATH",
	"PHP_CLI_WORKING_DIR",
	"RKE2_TYPE",
	"RUXIT_CLUSTER_ID",
	"RUXIT_NODE_ID",
	"SERVICE_NAME",
	"SOFTWAREAG_INSTALL_ROOT",
	"SOFTWAREAG_PRODUCTPROPNAME",
	"SPRINGBOOT_APP_NAME",
	"SPRINGBOOT_PROFILE_NAME",
	"SPRINGBOOT_STARTUP_CLASS",
	"TIBCO_BUSINESSWORKS_APP_NODE_NAME",
	"TIBCO_BUSINESSWORKS_APP_SPACE_NAME",
	"TIBCO_BUSINESSWORKS_CE_APP_NAME",
	"TIBCO_BUSINESSWORKS_CE_VERSION",
	"TIBCO_BUSINESSWORKS_DOMAIN_NAME",
	"TIBCO_BUSINESSWORKS_HOME",
	"TIPCO_BUSINESSWORKS_PROPERTY_FILE",
	"TIPCO_BUSINESSWORKS_PROPERTY_FILE_PATH",
	"UNKNOWN",
	"VARNISH_INSTANCE_NAME",
	"WEBLOGIC_CLUSTER_NAME",
	"WEBLOGIC_DOMAIN_NAME",
	"WEBLOGIC_HOME",
	"WEBLOGIC_NAME",
	"WEBSPHERE_CELL_NAME",
	"WEBSPHERE_CLUSTER_NAME",
	"WEBSPHERE_LIBERTY_SERVER_NAME",
	"WEBSPHERE_NODE_NAME",
	"WEBSPHERE_SERVER_NAME",
	"Z_CM_VERSION",
}

type ConditionOperator string

var ConditionOperators = struct {
	Contains    ConditionOperator
	Ends        ConditionOperator
	Equals      ConditionOperator
	Exists      ConditionOperator
	NotContains ConditionOperator
	NotEnds     ConditionOperator
	NotEquals   ConditionOperator
	NotExists   ConditionOperator
	NotStarts   ConditionOperator
	Starts      ConditionOperator
}{
	"CONTAINS",
	"ENDS",
	"EQUALS",
	"EXISTS",
	"NOT_CONTAINS",
	"NOT_ENDS",
	"NOT_EQUALS",
	"NOT_EXISTS",
	"NOT_STARTS",
	"STARTS",
}

type MonitoringMode string

var MonitoringModes = struct {
	MonitoringOff MonitoringMode
	MonitoringOn  MonitoringMode
}{
	"MONITORING_OFF",
	"MONITORING_ON",
}
