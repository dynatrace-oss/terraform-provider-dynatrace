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

// DynamicKey The key of the attribute, which need dynamic keys.
// Not applicable otherwise, as the attibute itself acts as a key.
type DynamicKey string

func (v DynamicKey) Ref() *DynamicKey {
	return &v
}

// DynamicKeys offers the known enum values
var DynamicKeys = struct {
	AmazonECRImageAccountID                  DynamicKey
	AmazonECRImageRegion                     DynamicKey
	AmazonLambdaFunctionName                 DynamicKey
	AmazonRegion                             DynamicKey
	ApacheConfigPath                         DynamicKey
	ApacheSparkMasterIPAddress               DynamicKey
	ASPDotNetCoreApplicationPath             DynamicKey
	AWSECSCluster                            DynamicKey
	AWSECSContainername                      DynamicKey
	AWSECSFamily                             DynamicKey
	AWSECSRevision                           DynamicKey
	CassandraClusterName                     DynamicKey
	CatalinaBase                             DynamicKey
	CatalinaHome                             DynamicKey
	CloudFoundryAppID                        DynamicKey
	CloudFoundryAppName                      DynamicKey
	CloudFoundryInstanceIndex                DynamicKey
	CloudFoundrySpaceID                      DynamicKey
	CloudFoundrySpaceName                    DynamicKey
	ColdfusionJVMConfigFile                  DynamicKey
	ColdfusionServiceName                    DynamicKey
	CommandLineArgs                          DynamicKey
	DotNetCommand                            DynamicKey
	DotNetCommandPath                        DynamicKey
	DynatraceClusterID                       DynamicKey
	DynatraceNodeID                          DynamicKey
	ElasticSearchClusterName                 DynamicKey
	ElasticSearchNodeName                    DynamicKey
	EquinoxConfigPath                        DynamicKey
	ExeName                                  DynamicKey
	ExePath                                  DynamicKey
	GlassFishDomainName                      DynamicKey
	GlassFishInstanceName                    DynamicKey
	GoogleAppEngineInstance                  DynamicKey
	GoogleAppEngineService                   DynamicKey
	GoogleCloudProject                       DynamicKey
	HybrisBinDirectory                       DynamicKey
	HybrisConfigDirectory                    DynamicKey
	HybrisDataDirectory                      DynamicKey
	IBMCICSRegion                            DynamicKey
	IBMCTGName                               DynamicKey
	IBMIMSConnectRegion                      DynamicKey
	IBMIMSControlRegion                      DynamicKey
	IBMIMSMessageProcessingRegion            DynamicKey
	IBMIMSSoapGwName                         DynamicKey
	IBMIntegrationNodeName                   DynamicKey
	IBMIntegrationServerName                 DynamicKey
	IISAppPool                               DynamicKey
	IISRoleName                              DynamicKey
	JavaJarFile                              DynamicKey
	JavaJarPath                              DynamicKey
	JavaMainClass                            DynamicKey
	JavaMainModule                           DynamicKey
	JBossHome                                DynamicKey
	JBossMode                                DynamicKey
	JBossServerName                          DynamicKey
	KubernetesBasePodName                    DynamicKey
	KubernetesContainerName                  DynamicKey
	KubernetesFullPodName                    DynamicKey
	KubernetesNamespace                      DynamicKey
	KubernetesPodUID                         DynamicKey
	MSSQLInstanceName                        DynamicKey
	NodeJsAppBaseDirectory                   DynamicKey
	NodeJsAppName                            DynamicKey
	NodeJsScriptName                         DynamicKey
	OracleSid                                DynamicKey
	PgIDCalcInputKeyLinkage                  DynamicKey
	PHPScriptPath                            DynamicKey
	PHPWorkingDirectory                      DynamicKey
	RubyAppRootPath                          DynamicKey
	RubyScriptPath                           DynamicKey
	RuleResult                               DynamicKey
	SoftwareAGInstallRoot                    DynamicKey
	SoftwareAGProductpropname                DynamicKey
	SpringBootAppName                        DynamicKey
	SpringBootProfileName                    DynamicKey
	SpringBootStartupClass                   DynamicKey
	TibcoBusinessworksCeAppName              DynamicKey
	TibcoBusinessworksCeVersion              DynamicKey
	TibcoBusinessWorksAppNodeName            DynamicKey
	TibcoBusinessWorksAppSpaceName           DynamicKey
	TibcoBusinessWorksDomainName             DynamicKey
	TibcoBusinessWorksEnginePropertyFile     DynamicKey
	TibcoBusinessWorksEnginePropertyFilePath DynamicKey
	TibcoBusinessWorksHome                   DynamicKey
	VarnishInstanceName                      DynamicKey
	WebLogicClusterName                      DynamicKey
	WebLogicDomainName                       DynamicKey
	WebLogicHome                             DynamicKey
	WebLogicName                             DynamicKey
	WebSphereCellName                        DynamicKey
	WebSphereClusterName                     DynamicKey
	WebSphereNodeName                        DynamicKey
	WebSphereServerName                      DynamicKey
}{
	"AMAZON_ECR_IMAGE_ACCOUNT_ID",
	"AMAZON_ECR_IMAGE_REGION",
	"AMAZON_LAMBDA_FUNCTION_NAME",
	"AMAZON_REGION",
	"APACHE_CONFIG_PATH",
	"APACHE_SPARK_MASTER_IP_ADDRESS",
	"ASP_DOT_NET_CORE_APPLICATION_PATH",
	"AWS_ECS_CLUSTER",
	"AWS_ECS_CONTAINERNAME",
	"AWS_ECS_FAMILY",
	"AWS_ECS_REVISION",
	"CASSANDRA_CLUSTER_NAME",
	"CATALINA_BASE",
	"CATALINA_HOME",
	"CLOUD_FOUNDRY_APP_ID",
	"CLOUD_FOUNDRY_APP_NAME",
	"CLOUD_FOUNDRY_INSTANCE_INDEX",
	"CLOUD_FOUNDRY_SPACE_ID",
	"CLOUD_FOUNDRY_SPACE_NAME",
	"COLDFUSION_JVM_CONFIG_FILE",
	"COLDFUSION_SERVICE_NAME",
	"COMMAND_LINE_ARGS",
	"DOTNET_COMMAND",
	"DOTNET_COMMAND_PATH",
	"DYNATRACE_CLUSTER_ID",
	"DYNATRACE_NODE_ID",
	"ELASTICSEARCH_CLUSTER_NAME",
	"ELASTICSEARCH_NODE_NAME",
	"EQUINOX_CONFIG_PATH",
	"EXE_NAME",
	"EXE_PATH",
	"GLASS_FISH_DOMAIN_NAME",
	"GLASS_FISH_INSTANCE_NAME",
	"GOOGLE_APP_ENGINE_INSTANCE",
	"GOOGLE_APP_ENGINE_SERVICE",
	"GOOGLE_CLOUD_PROJECT",
	"HYBRIS_BIN_DIRECTORY",
	"HYBRIS_CONFIG_DIRECTORY",
	"HYBRIS_DATA_DIRECTORY",
	"IBM_CICS_REGION",
	"IBM_CTG_NAME",
	"IBM_IMS_CONNECT_REGION",
	"IBM_IMS_CONTROL_REGION",
	"IBM_IMS_MESSAGE_PROCESSING_REGION",
	"IBM_IMS_SOAP_GW_NAME",
	"IBM_INTEGRATION_NODE_NAME",
	"IBM_INTEGRATION_SERVER_NAME",
	"IIS_APP_POOL",
	"IIS_ROLE_NAME",
	"JAVA_JAR_FILE",
	"JAVA_JAR_PATH",
	"JAVA_MAIN_CLASS",
	"JAVA_MAIN_MODULE",
	"JBOSS_HOME",
	"JBOSS_MODE",
	"JBOSS_SERVER_NAME",
	"KUBERNETES_BASE_POD_NAME",
	"KUBERNETES_CONTAINER_NAME",
	"KUBERNETES_FULL_POD_NAME",
	"KUBERNETES_NAMESPACE",
	"KUBERNETES_POD_UID",
	"MSSQL_INSTANCE_NAME",
	"NODE_JS_APP_BASE_DIRECTORY",
	"NODE_JS_APP_NAME",
	"NODE_JS_SCRIPT_NAME",
	"ORACLE_SID",
	"PG_ID_CALC_INPUT_KEY_LINKAGE",
	"PHP_SCRIPT_PATH",
	"PHP_WORKING_DIRECTORY",
	"RUBY_APP_ROOT_PATH",
	"RUBY_SCRIPT_PATH",
	"RULE_RESULT",
	"SOFTWAREAG_INSTALL_ROOT",
	"SOFTWAREAG_PRODUCTPROPNAME",
	"SPRINGBOOT_APP_NAME",
	"SPRINGBOOT_PROFILE_NAME",
	"SPRINGBOOT_STARTUP_CLASS",
	"TIBCO_BUSINESSWORKS_CE_APP_NAME",
	"TIBCO_BUSINESSWORKS_CE_VERSION",
	"TIBCO_BUSINESS_WORKS_APP_NODE_NAME",
	"TIBCO_BUSINESS_WORKS_APP_SPACE_NAME",
	"TIBCO_BUSINESS_WORKS_DOMAIN_NAME",
	"TIBCO_BUSINESS_WORKS_ENGINE_PROPERTY_FILE",
	"TIBCO_BUSINESS_WORKS_ENGINE_PROPERTY_FILE_PATH",
	"TIBCO_BUSINESS_WORKS_HOME",
	"VARNISH_INSTANCE_NAME",
	"WEB_LOGIC_CLUSTER_NAME",
	"WEB_LOGIC_DOMAIN_NAME",
	"WEB_LOGIC_HOME",
	"WEB_LOGIC_NAME",
	"WEB_SPHERE_CELL_NAME",
	"WEB_SPHERE_CLUSTER_NAME",
	"WEB_SPHERE_NODE_NAME",
	"WEB_SPHERE_SERVER_NAME",
}
