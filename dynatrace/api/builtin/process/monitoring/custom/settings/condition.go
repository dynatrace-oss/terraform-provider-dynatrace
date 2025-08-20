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

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Condition struct {
	EnvVar   *string           `json:"envVar,omitempty"` // supported only with OneAgent 1.167+
	Item     AgentItemName     `json:"item"`             // Condition target
	Operator ConditionOperator `json:"operator"`         // Condition operator
	Value    *string           `json:"value,omitempty"`  // Condition value
}

func (me *Condition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"env_var": {
			Type:        schema.TypeString,
			Description: "supported only with OneAgent 1.167+",
			Optional:    true,
		},
		"item": {
			Type:        schema.TypeString,
			Description: "Possible Values: `APACHE_CONFIG_PATH`, `APACHE_SPARK_MASTER_IP_ADDRESS`, `ASP_NET_CORE_APPLICATION_PATH`, `AWS_ECR_ACCOUNT_ID`, `AWS_ECR_REGION`, `AWS_ECS_CLUSTER`, `AWS_ECS_CONTAINERNAME`, `AWS_ECS_FAMILY`, `AWS_ECS_REVISION`, `AWS_LAMBDA_FUNCTION_NAME`, `AWS_REGION`, `AZURE_CONTAINER_APP_ENV_DNS_SUFFIX`, `AZURE_CONTAINER_APP_NAME`, `CATALINA_BASE`, `CATALINA_HOME`, `CLOUD_FOUNDRY_APPLICATION_ID`, `CLOUD_FOUNDRY_APP_NAME`, `CLOUD_FOUNDRY_INSTANCE_INDEX`, `CLOUD_FOUNDRY_SPACE_ID`, `CLOUD_FOUNDRY_SPACE_NAME`, `COLDFUSION_JVM_CONFIG_FILE`, `COMMAND_LINE_ARGS`, `CONTAINER_ID`, `CONTAINER_IMAGE_NAME`, `CONTAINER_IMAGE_VERSION`, `CONTAINER_NAME`, `DATASOURCE_MONITORING_CONFIG_ID`, `DECLARATIVE_ID`, `DOTNET_COMMAND`, `DOTNET_COMMAND_PATH`, `ELASTIC_SEARCH_CLUSTER_NAME`, `ELASTIC_SEARCH_NODE_NAME`, `EQUINOX_CONFIG_PATH`, `EXE_NAME`, `EXE_PATH`, `GAE_INSTANCE`, `GAE_SERVICE`, `GLASSFISH_DOMAIN_NAME`, `GLASSFISH_INSTANCE_NAME`, `GOOGLE_CLOUD_PROJECT`, `HYBRIS_BIN_DIR`, `HYBRIS_CONFIG_DIR`, `HYBRIS_DATA_DIR`, `IBM_APPLID`, `IBM_CICS_IMS_APPLID`, `IBM_CICS_IMS_JOBNAME`, `IBM_CICS_REGION`, `IBM_CTG_NAME`, `IBM_IMS_CONNECT`, `IBM_IMS_CONTROL`, `IBM_IMS_MPR`, `IBM_IMS_SOAP_GW_NAME`, `IBM_JOBNAME`, `IIB_BROKER_NAME`, `IIB_EXECUTION_GROUP_NAME`, `IIS_APP_POOL`, `IIS_ROLE_NAME`, `JAVA_JAR_FILE`, `JAVA_JAR_PATH`, `JAVA_MAIN_CLASS`, `JBOSS_HOME`, `JBOSS_MODE`, `JBOSS_SERVER_NAME`, `KUBERNETES_BASEPODNAME`, `KUBERNETES_CONTAINERNAME`, `KUBERNETES_FULLPODNAME`, `KUBERNETES_NAMESPACE`, `KUBERNETES_PODUID`, `MSSQL_INSTANCE_NAME`, `NODEJS_APP_BASE_DIR`, `NODEJS_APP_NAME`, `NODEJS_SCRIPT_NAME`, `ORACLE_SID`, `PG_ID_CALC_INPUT_KEY_LINKAGE`, `PHP_CLI_SCRIPT_PATH`, `PHP_CLI_WORKING_DIR`, `PYTHON_MODULE`, `PYTHON_SCRIPT`, `PYTHON_SCRIPT_PATH`, `RKE2_TYPE`, `RUXIT_CLUSTER_ID`, `RUXIT_NODE_ID`, `SERVICE_NAME`, `SOFTWAREAG_INSTALL_ROOT`, `SOFTWAREAG_PRODUCTPROPNAME`, `SPRINGBOOT_APP_NAME`, `SPRINGBOOT_PROFILE_NAME`, `SPRINGBOOT_STARTUP_CLASS`, `TIBCO_BUSINESSWORKS_APP_NODE_NAME`, `TIBCO_BUSINESSWORKS_APP_SPACE_NAME`, `TIBCO_BUSINESSWORKS_CE_APP_NAME`, `TIBCO_BUSINESSWORKS_CE_VERSION`, `TIBCO_BUSINESSWORKS_DOMAIN_NAME`, `TIBCO_BUSINESSWORKS_HOME`, `TIPCO_BUSINESSWORKS_PROPERTY_FILE`, `TIPCO_BUSINESSWORKS_PROPERTY_FILE_PATH`, `UNKNOWN`, `VARNISH_INSTANCE_NAME`, `WEBLOGIC_CLUSTER_NAME`, `WEBLOGIC_DOMAIN_NAME`, `WEBLOGIC_HOME`, `WEBLOGIC_NAME`, `WEBSPHERE_CELL_NAME`, `WEBSPHERE_CLUSTER_NAME`, `WEBSPHERE_LIBERTY_SERVER_NAME`, `WEBSPHERE_NODE_NAME`, `WEBSPHERE_SERVER_NAME`, `Z_CM_VERSION`",
			Required:    true,
		},
		"operator": {
			Type:        schema.TypeString,
			Description: "Possible Values: `CONTAINS`, `ENDS`, `EQUALS`, `EXISTS`, `NOT_CONTAINS`, `NOT_ENDS`, `NOT_EQUALS`, `NOT_EXISTS`, `NOT_STARTS`, `STARTS`",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "Condition value",
			Optional:    true,
		},
	}
}

func (me *Condition) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"env_var":  me.EnvVar,
		"item":     me.Item,
		"operator": me.Operator,
		"value":    me.Value,
	})
}

func (me *Condition) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"env_var":  &me.EnvVar,
		"item":     &me.Item,
		"operator": &me.Operator,
		"value":    &me.Value,
	})
}
