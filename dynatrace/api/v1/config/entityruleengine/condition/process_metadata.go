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

// ProcessMetadata The key for dynamic attributes of the `PROCESS_PREDEFINED_METADATA_KEY` type.
type ProcessMetadata struct {
	BaseConditionKey
	DynamicKey *DynamicKey `json:"dynamicKey"` // The key of the attribute, which need dynamic keys. Not applicable otherwise, as the attibute itself acts as a key.
}

func (pmck *ProcessMetadata) GetType() *ConditionKeyType {
	return &ConditionKeyTypes.ProcessPredefinedMetadataKey
}

func (pmck *ProcessMetadata) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"attribute": {
			Type:        schema.TypeString,
			Description: "The attribute to be used for comparision",
			Required:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "if specified, needs to be PROCESS_PREDEFINED_METADATA_KEY",
			Optional:    true,
			Deprecated:  "The value of the attribute type is implicit, therefore shouldn't get specified",
		},
		"dynamic_key": {
			Type:        schema.TypeString,
			Description: "The key of the attribute, which need dynamic keys. Not applicable otherwise, as the attibute itself acts as a key. Possible values are AMAZON_ECR_IMAGE_ACCOUNT_ID,AMAZON_ECR_IMAGE_REGION, AMAZON_LAMBDA_FUNCTION_NAME, AMAZON_REGION, APACHE_CONFIG_PATH, APACHE_SPARK_MASTER_IP_ADDRESS, ASP_DOT_NET_CORE_APPLICATION_PATH, AWS_ECS_CLUSTER, AWS_ECS_CONTAINERNAME, AWS_ECS_FAMILY, AWS_ECS_REVISION, CASSANDRA_CLUSTER_NAME, CATALINA_BASE, CATALINA_HOME, CLOUD_FOUNDRY_APP_ID, CLOUD_FOUNDRY_APP_NAME, CLOUD_FOUNDRY_INSTANCE_INDEX, CLOUD_FOUNDRY_SPACE_ID, CLOUD_FOUNDRY_SPACE_NAME, COLDFUSION_JVM_CONFIG_FILE, COLDFUSION_SERVICE_NAME, COMMAND_LINE_ARGS, DOTNET_COMMAND, DOTNET_COMMAND_PATH, DYNATRACE_CLUSTER_ID, DYNATRACE_NODE_ID, ELASTICSEARCH_CLUSTER_NAME, ELASTICSEARCH_NODE_NAME, EQUINOX_CONFIG_PATH, EXE_NAME, EXE_PATH, GLASS_FISH_DOMAIN_NAME, GLASS_FISH_INSTANCE_NAME, GOOGLE_APP_ENGINE_INSTANCE, GOOGLE_APP_ENGINE_SERVICE, GOOGLE_CLOUD_PROJECT, HYBRIS_BIN_DIRECTORY, HYBRIS_CONFIG_DIRECTORY, HYBRIS_DATA_DIRECTORY, IBM_CICS_REGION, IBM_CTG_NAME, IBM_IMS_CONNECT_REGION, IBM_IMS_CONTROL_REGION, IBM_IMS_MESSAGE_PROCESSING_REGION, IBM_IMS_SOAP_GW_NAME, IBM_INTEGRATION_NODE_NAME, IBM_INTEGRATION_SERVER_NAME, IIS_APP_POOL, IIS_ROLE_NAME, JAVA_JAR_FILE, JAVA_JAR_PATH, JAVA_MAIN_CLASS, JAVA_MAIN_MODULE, JBOSS_HOME, JBOSS_MODE, JBOSS_SERVER_NAME, KUBERNETES_BASE_POD_NAME, KUBERNETES_CONTAINER_NAME, KUBERNETES_FULL_POD_NAME, KUBERNETES_NAMESPACE, KUBERNETES_POD_UID, MSSQL_INSTANCE_NAME, NODE_JS_APP_BASE_DIRECTORY, NODE_JS_APP_NAME, NODE_JS_SCRIPT_NAME, ORACLE_SID, PG_ID_CALC_INPUT_KEY_LINKAGE, PHP_SCRIPT_PATH, PHP_WORKING_DIRECTORY, RUBY_APP_ROOT_PATH, RUBY_SCRIPT_PATH, RULE_RESULT, SOFTWAREAG_INSTALL_ROOT, SOFTWAREAG_PRODUCTPROPNAME, SPRINGBOOT_APP_NAME, SPRINGBOOT_PROFILE_NAME, SPRINGBOOT_STARTUP_CLASS, TIBCO_BUSINESSWORKS_CE_APP_NAME, TIBCO_BUSINESSWORKS_CE_VERSION, TIBCO_BUSINESS_WORKS_APP_NODE_NAME, TIBCO_BUSINESS_WORKS_APP_SPACE_NAME, TIBCO_BUSINESS_WORKS_DOMAIN_NAME, TIBCO_BUSINESS_WORKS_ENGINE_PROPERTY_FILE, TIBCO_BUSINESS_WORKS_ENGINE_PROPERTY_FILE_PATH, TIBCO_BUSINESS_WORKS_HOME, VARNISH_INSTANCE_NAME, WEB_LOGIC_CLUSTER_NAME, WEB_LOGIC_DOMAIN_NAME, WEB_LOGIC_HOME, WEB_LOGIC_NAME, WEB_SPHERE_CELL_NAME, WEB_SPHERE_CLUSTER_NAME, WEB_SPHERE_NODE_NAME and WEB_SPHERE_SERVER_NAME",
			Required:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider but have meanwhile gotten introduced by a newer version of the Dynatrace REST API",
			Optional:    true,
		},
	}
}

func (pmck *ProcessMetadata) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(pmck.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("attribute", string(pmck.Attribute)); err != nil {
		return err
	}
	if err := properties.Encode("dynamic_key", pmck.DynamicKey); err != nil {
		return err
	}
	return nil
}

func (pmck *ProcessMetadata) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), pmck); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &pmck.Unknowns); err != nil {
			return err
		}
		delete(pmck.Unknowns, "attribute")
		delete(pmck.Unknowns, "dynamic_key")
		delete(pmck.Unknowns, "type")
		if len(pmck.Unknowns) == 0 {
			pmck.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("attribute"); ok {
		pmck.Attribute = Attribute(value.(string))
	}
	if value, ok := decoder.GetOk("dynamic_key"); ok {
		pmck.DynamicKey = DynamicKey(value.(string)).Ref()
	}
	return nil
}

func (pmck *ProcessMetadata) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(pmck.Unknowns) > 0 {
		for k, v := range pmck.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(pmck.Attribute)
		if err != nil {
			return nil, err
		}
		m["attribute"] = rawMessage
	}
	if pmck.GetType() != nil {
		rawMessage, err := json.Marshal(ConditionKeyTypes.ProcessPredefinedMetadataKey)
		if err != nil {
			return nil, err
		}
		m["type"] = rawMessage
	}
	if pmck.DynamicKey != nil {
		rawMessage, err := json.Marshal(pmck.DynamicKey)
		if err != nil {
			return nil, err
		}
		m["dynamicKey"] = rawMessage
	}
	return json.Marshal(m)
}

func (pmck *ProcessMetadata) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	pmck.Type = pmck.GetType()
	if v, found := m["attribute"]; found {
		if err := json.Unmarshal(v, &pmck.Attribute); err != nil {
			return err
		}
	}
	if v, found := m["dynamicKey"]; found {
		if err := json.Unmarshal(v, &pmck.DynamicKey); err != nil {
			return err
		}
	}
	delete(m, "attribute")
	delete(m, "dynamicKey")
	delete(m, "type")
	if len(m) > 0 {
		pmck.Unknowns = m
	}
	return nil
}
