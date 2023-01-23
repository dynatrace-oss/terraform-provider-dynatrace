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

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// String The key for dynamic attributes of the `STRING` type.
type String struct {
	BaseConditionKey
	DynamicKey string `json:"dynamicKey"` // The key of the attribute, which need dynamic keys. Not applicable otherwise, as the attibute itself acts as a key.
}

func (sck *String) GetType() *ConditionKeyType {
	return &ConditionKeyTypes.String
}

func (sck *String) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"attribute": {
			Type:        schema.TypeString,
			Description: "The attribute to be used for comparision",
			Required:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "if specified, needs to be `STRING`",
			Optional:    true,
			Deprecated:  "The value of the attribute type is implicit, therefore shouldn't get specified",
		},
		"dynamic_key": {
			Type:        schema.TypeString,
			Description: "The key of the attribute, which need dynamic keys. Not applicable otherwise, as the attibute itself acts as a key. Possible values are\n   - `AMAZON_ECR_IMAGE_ACCOUNT_ID`\n   - `AMAZON_ECR_IMAGE_REGION`\n   - `AMAZON_LAMBDA_FUNCTION_NAME`\n   - `AMAZON_REGION`\n   - `APACHE_CONFIG_PATH`\n   - `APACHE_SPARK_MASTER_IP_ADDRESS`\n   - `ASP_DOT_NET_CORE_APPLICATION_PATH`\n   - `AWS_ECS_CLUSTER`\n   - `AWS_ECS_CONTAINERNAME`\n   - `AWS_ECS_FAMILY`\n   - `AWS_ECS_REVISION`\n   - `CASSANDRA_CLUSTER_NAME`\n   - `CATALINA_BASE`\n   - `CATALINA_HOME`\n   - `CLOUD_FOUNDRY_APP_ID`\n   - `CLOUD_FOUNDRY_APP_NAME`\n   - `CLOUD_FOUNDRY_INSTANCE_INDEX`\n   - `CLOUD_FOUNDRY_SPACE_ID`\n   - `CLOUD_FOUNDRY_SPACE_NAME`\n   - `COLDFUSION_JVM_CONFIG_FILE`\n   - `COLDFUSION_SERVICE_NAME`\n   - `COMMAND_LINE_ARGS`\n   - `DOTNET_COMMAND`\n   - `DOTNET_COMMAND_PATH`\n   - `DYNATRACE_CLUSTER_ID`\n   - `DYNATRACE_NODE_ID`\n   - `ELASTICSEARCH_CLUSTER_NAME`\n   - `ELASTICSEARCH_NODE_NAME`\n   - `EQUINOX_CONFIG_PATH`\n   - `EXE_NAME`\n   - `EXE_PATH`\n   - `GLASS_FISH_DOMAIN_NAME`\n   - `GLASS_FISH_INSTANCE_NAME`\n   - `GOOGLE_APP_ENGINE_INSTANCE`\n   - `GOOGLE_APP_ENGINE_SERVICE`\n   - `GOOGLE_CLOUD_PROJECT`\n   - `HYBRIS_BIN_DIRECTORY`\n   - `HYBRIS_CONFIG_DIRECTORY`\n   - `HYBRIS_DATA_DIRECTORY`\n   - `IBM_CICS_REGION`\n   - `IBM_CTG_NAME`\n   - `IBM_IMS_CONNECT_REGION`\n   - `IBM_IMS_CONTROL_REGION`\n   - `IBM_IMS_MESSAGE_PROCESSING_REGION`\n   - `IBM_IMS_SOAP_GW_NAME`\n   - `IBM_INTEGRATION_NODE_NAME`\n   - `IBM_INTEGRATION_SERVER_NAME`\n   - `IIS_APP_POOL`\n   - `IIS_ROLE_NAME`\n   - `JAVA_JAR_FILE`\n   - `JAVA_JAR_PATH`\n   - `JAVA_MAIN_CLASS`\n   - `JAVA_MAIN_MODULE`\n   - `JBOSS_HOME`\n   - `JBOSS_MODE`\n   - `JBOSS_SERVER_NAME`\n   - `KUBERNETES_BASE_POD_NAME`\n   - `KUBERNETES_CONTAINER_NAME`\n   - `KUBERNETES_FULL_POD_NAME`\n   - `KUBERNETES_NAMESPACE`\n   - `KUBERNETES_POD_UID`\n   - `MSSQL_INSTANCE_NAME`\n   - `NODE_JS_APP_BASE_DIRECTORY`\n   - `NODE_JS_APP_NAME`\n   - `NODE_JS_SCRIPT_NAME`\n   - `ORACLE_SID`\n   - `PG_ID_CALC_INPUT_KEY_LINKAGE`\n   - `PHP_SCRIPT_PATH`\n   - `PHP_WORKING_DIRECTORY`\n   - `RUBY_APP_ROOT_PATH`\n   - `RUBY_SCRIPT_PATH`\n   - `RULE_RESULT`\n   - `SOFTWAREAG_INSTALL_ROOT`\n   - `SOFTWAREAG_PRODUCTPROPNAME`\n   - `SPRINGBOOT_APP_NAME`\n   - `SPRINGBOOT_PROFILE_NAME`\n   - `SPRINGBOOT_STARTUP_CLASS`\n   - `TIBCO_BUSINESSWORKS_CE_APP_NAME`\n   - `TIBCO_BUSINESSWORKS_CE_VERSION`\n   - `TIBCO_BUSINESS_WORKS_APP_NODE_NAME`\n   - `TIBCO_BUSINESS_WORKS_APP_SPACE_NAME`\n   - `TIBCO_BUSINESS_WORKS_DOMAIN_NAME`\n   - `TIBCO_BUSINESS_WORKS_ENGINE_PROPERTY_FILE`\n   - `TIBCO_BUSINESS_WORKS_ENGINE_PROPERTY_FILE_PATH`\n   - `TIBCO_BUSINESS_WORKS_HOME`\n   - `VARNISH_INSTANCE_NAME`\n   - `WEB_LOGIC_CLUSTER_NAME`\n   - `WEB_LOGIC_DOMAIN_NAME`\n   - `WEB_LOGIC_HOME`\n   - `WEB_LOGIC_NAME`\n   - `WEB_SPHERE_CELL_NAME`\n   - `WEB_SPHERE_CLUSTER_NAME`\n   - `WEB_SPHERE_NODE_NAME and WEB_SPHERE_SERVER_NAME`",
			Required:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider but have meanwhile gotten introduced by a newer version of the Dynatrace REST API",
			Optional:    true,
		},
	}
}

func (sck *String) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(sck.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("attribute", string(sck.Attribute)); err != nil {
		return err
	}
	if err := properties.Encode("dynamic_key", sck.DynamicKey); err != nil {
		return err
	}
	return nil
}

func (sck *String) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), sck); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &sck.Unknowns); err != nil {
			return err
		}
		delete(sck.Unknowns, "attribute")
		delete(sck.Unknowns, "dynamic_key")
		delete(sck.Unknowns, "type")
		if len(sck.Unknowns) == 0 {
			sck.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("attribute"); ok {
		sck.Attribute = Attribute(value.(string))
	}
	if value, ok := decoder.GetOk("dynamic_key"); ok {
		sck.DynamicKey = value.(string)
	}
	return nil
}

func (sck *String) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(sck.Unknowns) > 0 {
		for k, v := range sck.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(sck.Attribute)
		if err != nil {
			return nil, err
		}
		m["attribute"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(ConditionKeyTypes.String)
		if err != nil {
			return nil, err
		}
		m["type"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(sck.DynamicKey)
		if err != nil {
			return nil, err
		}
		m["dynamicKey"] = rawMessage
	}
	return json.Marshal(m)
}

func (sck *String) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	sck.Type = sck.GetType()
	if v, found := m["attribute"]; found {
		if err := json.Unmarshal(v, &sck.Attribute); err != nil {
			return err
		}
	}
	if v, found := m["dynamicKey"]; found {
		if err := json.Unmarshal(v, &sck.DynamicKey); err != nil {
			return err
		}
	}
	delete(m, "attribute")
	delete(m, "dynamicKey")
	delete(m, "type")
	if len(m) > 0 {
		sck.Unknowns = m
	}
	return nil
}
