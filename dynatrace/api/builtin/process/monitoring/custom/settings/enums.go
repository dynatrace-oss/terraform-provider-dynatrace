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
}{
	AgentItemName("APACHE_CONFIG_PATH"),
	AgentItemName("APACHE_SPARK_MASTER_IP_ADDRESS"),
	AgentItemName("ASP_NET_CORE_APPLICATION_PATH"),
	AgentItemName("AWS_ECR_ACCOUNT_ID"),
	AgentItemName("AWS_ECR_REGION"),
	AgentItemName("AWS_ECS_CLUSTER"),
	AgentItemName("AWS_ECS_CONTAINERNAME"),
	AgentItemName("AWS_ECS_FAMILY"),
	AgentItemName("AWS_ECS_REVISION"),
	AgentItemName("AWS_LAMBDA_FUNCTION_NAME"),
	AgentItemName("AWS_REGION"),
	AgentItemName("CATALINA_BASE"),
	AgentItemName("CATALINA_HOME"),
	AgentItemName("CLOUD_FOUNDRY_APP_NAME"),
	AgentItemName("CLOUD_FOUNDRY_APPLICATION_ID"),
	AgentItemName("CLOUD_FOUNDRY_INSTANCE_INDEX"),
	AgentItemName("CLOUD_FOUNDRY_SPACE_ID"),
	AgentItemName("CLOUD_FOUNDRY_SPACE_NAME"),
	AgentItemName("COLDFUSION_JVM_CONFIG_FILE"),
	AgentItemName("COMMAND_LINE_ARGS"),
	AgentItemName("CONTAINER_ID"),
	AgentItemName("CONTAINER_IMAGE_NAME"),
	AgentItemName("CONTAINER_IMAGE_VERSION"),
	AgentItemName("CONTAINER_NAME"),
	AgentItemName("DECLARATIVE_ID"),
	AgentItemName("DOTNET_COMMAND"),
	AgentItemName("DOTNET_COMMAND_PATH"),
	AgentItemName("ELASTIC_SEARCH_CLUSTER_NAME"),
	AgentItemName("ELASTIC_SEARCH_NODE_NAME"),
	AgentItemName("EQUINOX_CONFIG_PATH"),
	AgentItemName("EXE_NAME"),
	AgentItemName("EXE_PATH"),
	AgentItemName("GAE_INSTANCE"),
	AgentItemName("GAE_SERVICE"),
	AgentItemName("GLASSFISH_DOMAIN_NAME"),
	AgentItemName("GLASSFISH_INSTANCE_NAME"),
	AgentItemName("GOOGLE_CLOUD_PROJECT"),
	AgentItemName("HYBRIS_BIN_DIR"),
	AgentItemName("HYBRIS_CONFIG_DIR"),
	AgentItemName("HYBRIS_DATA_DIR"),
	AgentItemName("IBM_CICS_IMS_APPLID"),
	AgentItemName("IBM_CICS_IMS_JOBNAME"),
	AgentItemName("IBM_CICS_REGION"),
	AgentItemName("IBM_CTG_NAME"),
	AgentItemName("IBM_IMS_CONNECT"),
	AgentItemName("IBM_IMS_CONTROL"),
	AgentItemName("IBM_IMS_MPR"),
	AgentItemName("IBM_IMS_SOAP_GW_NAME"),
	AgentItemName("IIB_BROKER_NAME"),
	AgentItemName("IIB_EXECUTION_GROUP_NAME"),
	AgentItemName("IIS_APP_POOL"),
	AgentItemName("IIS_ROLE_NAME"),
	AgentItemName("JAVA_JAR_FILE"),
	AgentItemName("JAVA_JAR_PATH"),
	AgentItemName("JAVA_MAIN_CLASS"),
	AgentItemName("JBOSS_HOME"),
	AgentItemName("JBOSS_MODE"),
	AgentItemName("JBOSS_SERVER_NAME"),
	AgentItemName("KUBERNETES_BASEPODNAME"),
	AgentItemName("KUBERNETES_CONTAINERNAME"),
	AgentItemName("KUBERNETES_FULLPODNAME"),
	AgentItemName("KUBERNETES_NAMESPACE"),
	AgentItemName("KUBERNETES_PODUID"),
	AgentItemName("MSSQL_INSTANCE_NAME"),
	AgentItemName("NODEJS_APP_BASE_DIR"),
	AgentItemName("NODEJS_APP_NAME"),
	AgentItemName("NODEJS_SCRIPT_NAME"),
	AgentItemName("ORACLE_SID"),
	AgentItemName("PG_ID_CALC_INPUT_KEY_LINKAGE"),
	AgentItemName("PHP_CLI_SCRIPT_PATH"),
	AgentItemName("PHP_CLI_WORKING_DIR"),
	AgentItemName("RUXIT_CLUSTER_ID"),
	AgentItemName("RUXIT_NODE_ID"),
	AgentItemName("SERVICE_NAME"),
	AgentItemName("SOFTWAREAG_INSTALL_ROOT"),
	AgentItemName("SOFTWAREAG_PRODUCTPROPNAME"),
	AgentItemName("SPRINGBOOT_APP_NAME"),
	AgentItemName("SPRINGBOOT_PROFILE_NAME"),
	AgentItemName("SPRINGBOOT_STARTUP_CLASS"),
	AgentItemName("TIBCO_BUSINESSWORKS_APP_NODE_NAME"),
	AgentItemName("TIBCO_BUSINESSWORKS_APP_SPACE_NAME"),
	AgentItemName("TIBCO_BUSINESSWORKS_CE_APP_NAME"),
	AgentItemName("TIBCO_BUSINESSWORKS_CE_VERSION"),
	AgentItemName("TIBCO_BUSINESSWORKS_DOMAIN_NAME"),
	AgentItemName("TIBCO_BUSINESSWORKS_HOME"),
	AgentItemName("TIPCO_BUSINESSWORKS_PROPERTY_FILE"),
	AgentItemName("TIPCO_BUSINESSWORKS_PROPERTY_FILE_PATH"),
	AgentItemName("UNKNOWN"),
	AgentItemName("VARNISH_INSTANCE_NAME"),
	AgentItemName("WEBLOGIC_CLUSTER_NAME"),
	AgentItemName("WEBLOGIC_DOMAIN_NAME"),
	AgentItemName("WEBLOGIC_HOME"),
	AgentItemName("WEBLOGIC_NAME"),
	AgentItemName("WEBSPHERE_CELL_NAME"),
	AgentItemName("WEBSPHERE_CLUSTER_NAME"),
	AgentItemName("WEBSPHERE_LIBERTY_SERVER_NAME"),
	AgentItemName("WEBSPHERE_NODE_NAME"),
	AgentItemName("WEBSPHERE_SERVER_NAME"),
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
	ConditionOperator("CONTAINS"),
	ConditionOperator("ENDS"),
	ConditionOperator("EQUALS"),
	ConditionOperator("EXISTS"),
	ConditionOperator("NOT_CONTAINS"),
	ConditionOperator("NOT_ENDS"),
	ConditionOperator("NOT_EQUALS"),
	ConditionOperator("NOT_EXISTS"),
	ConditionOperator("NOT_STARTS"),
	ConditionOperator("STARTS"),
}

type MonitoringMode string

var MonitoringModes = struct {
	MonitoringOff MonitoringMode
	MonitoringOn  MonitoringMode
}{
	MonitoringMode("MONITORING_OFF"),
	MonitoringMode("MONITORING_ON"),
}
