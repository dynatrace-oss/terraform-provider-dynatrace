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

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AttributeConditions []*AttributeCondition

func (me *AttributeConditions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"condition": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "Attribute conditions",
			Elem:        &schema.Resource{Schema: new(AttributeCondition).Schema()},
		},
	}
}

func (me AttributeConditions) MarshalHCL(properties hcl.Properties) error {
	if len(me) > 0 {
		if err := properties.EncodeSlice("condition", me); err != nil {
			return err
		}
	}
	return nil
}

func (me *AttributeConditions) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("condition"); ok {

		entrySet := value.(*schema.Set)

		for _, entryMap := range entrySet.List() {
			hash := entrySet.F(entryMap)
			entry := new(AttributeCondition)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "condition", hash)); err != nil {
				return err
			}
			*me = append(*me, entry)
		}
	}
	return nil
}

// No documentation available
type AttributeCondition struct {
	EnumValue        *string   `json:"enumValue,omitempty"`        // Value
	Tag              *string   `json:"tag,omitempty"`              // Tag. Format: `[CONTEXT]tagKey:tagValue`
	Key              Attribute `json:"key"`                        // Property
	DynamicKey       *string   `json:"dynamicKey,omitempty"`       // Dynamic key
	CaseSensitive    *bool     `json:"caseSensitive,omitempty"`    // Case sensitive
	IntegerValue     *int      `json:"integerValue,omitempty"`     // Value
	EntityID         *string   `json:"entityId,omitempty"`         // Value
	Operator         Operator  `json:"operator"`                   // Operator
	StringValue      *string   `json:"stringValue,omitempty"`      // Value
	DynamicKeySource *string   `json:"dynamicKeySource,omitempty"` // Key source
}

func (me *AttributeCondition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enum_value": {
			Type:        schema.TypeString,
			Description: "Value",
			Optional:    true,
		},
		"tag": {
			Type:        schema.TypeString,
			Description: "Tag. Format: `[CONTEXT]tagKey:tagValue`",
			Optional:    true,
		},
		"key": {
			Type:        schema.TypeString,
			Description: "Property",
			Required:    true,
		},
		"dynamic_key": {
			Type:        schema.TypeString,
			Description: "Dynamic key",
			Optional:    true,
		},
		"case_sensitive": {
			Type:        schema.TypeBool,
			Description: "Case sensitive",
			Optional:    true,
		},
		"integer_value": {
			Type:        schema.TypeInt,
			Description: "Value",
			Optional:    true,
		},
		"entity_id": {
			Type:        schema.TypeString,
			Description: "Value",
			Optional:    true,
		},
		"operator": {
			Type:        schema.TypeString,
			Description: "Operator",
			Required:    true,
		},
		"string_value": {
			Type:        schema.TypeString,
			Description: "Value",
			Optional:    true,
		},
		"dynamic_key_source": {
			Type:        schema.TypeString,
			Description: "Key source",
			Optional:    true,
		},
	}
}

func (me *AttributeCondition) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enum_value":         me.EnumValue,
		"tag":                me.Tag,
		"key":                me.Key,
		"dynamic_key":        me.DynamicKey,
		"case_sensitive":     me.CaseSensitive,
		"integer_value":      me.IntegerValue,
		"entity_id":          me.EntityID,
		"operator":           me.Operator,
		"string_value":       me.StringValue,
		"dynamic_key_source": me.DynamicKeySource,
	})
}

func (me *AttributeCondition) UnmarshalHCL(decoder hcl.Decoder) error {
	caseSensitiveIn := []string{"CLOUD_APPLICATION_LABELS", "CLOUD_APPLICATION_NAMESPACE_LABELS", "HOST_KUBERNETES_LABELS", "PROCESS_GROUP_PREDEFINED_METADATA", "CUSTOM_DEVICE_METADATA", "ENTERPRISE_APPLICATION_METADATA", "DATA_CENTER_SERVICE_METADATA", "HOST_CUSTOM_METADATA", "PROCESS_GROUP_CUSTOM_METADATA", "KUBERNETES_SERVICE_NAME", "HOST_AZURE_WEB_APPLICATION_HOST_NAMES", "HOST_AZURE_WEB_APPLICATION_SITE_NAMES", "HOST_DETECTED_NAME", "HOST_NAME", "HOST_OS_VERSION", "HOST_BOSH_NAME", "HOST_BOSH_INSTANCE_ID", "HOST_BOSH_INSTANCE_NAME", "HOST_BOSH_AVAILABILITY_ZONE", "HOST_BOSH_DEPLOYMENT_ID", "HOST_BOSH_STEMCELL_VERSION", "HOST_AWS_NAME_TAG", "HOST_ONEAGENT_CUSTOM_HOST_NAME", "KUBERNETES_CLUSTER_NAME", "KUBERNETES_NODE_NAME", "CLOUD_APPLICATION_NAMESPACE_NAME", "CLOUD_APPLICATION_NAME", "PROCESS_GROUP_AZURE_HOST_NAME", "PROCESS_GROUP_AZURE_SITE_NAME", "PROCESS_GROUP_DETECTED_NAME", "PROCESS_GROUP_NAME", "PROCESS_GROUP_TECHNOLOGY_EDITION", "PROCESS_GROUP_TECHNOLOGY_VERSION", "SERVICE_AKKA_ACTOR_SYSTEM", "SERVICE_DATABASE_NAME", "SERVICE_DATABASE_VENDOR", "SERVICE_DATABASE_HOST_NAME", "SERVICE_DETECTED_NAME", "SERVICE_WEB_SERVER_ENDPOINT", "SERVICE_IBM_CTG_GATEWAY_URL", "SERVICE_MESSAGING_LISTENER_CLASS_NAME", "SERVICE_NAME", "SERVICE_PUBLIC_DOMAIN_NAME", "SERVICE_REMOTE_ENDPOINT", "SERVICE_REMOTE_SERVICE_NAME", "SERVICE_TECHNOLOGY_EDITION", "SERVICE_TECHNOLOGY_VERSION", "SERVICE_WEB_APPLICATION_ID", "SERVICE_WEB_CONTEXT_ROOT", "SERVICE_WEB_SERVER_NAME", "SERVICE_WEB_SERVICE_NAME", "SERVICE_WEB_SERVICE_NAMESPACE", "SERVICE_CTG_SERVICE_NAME", "SERVICE_ESB_APPLICATION_NAME", "CUSTOM_DEVICE_DNS_ADDRESS", "CUSTOM_DEVICE_NAME", "CUSTOM_DEVICE_GROUP_NAME", "WEB_APPLICATION_NAME", "WEB_APPLICATION_NAME_PATTERN", "MOBILE_APPLICATION_NAME", "CUSTOM_APPLICATION_NAME", "ENTERPRISE_APPLICATION_NAME", "DATA_CENTER_SERVICE_NAME", "BROWSER_MONITOR_NAME", "EXTERNAL_MONITOR_NAME", "EXTERNAL_MONITOR_ENGINE_NAME", "EXTERNAL_MONITOR_ENGINE_DESCRIPTION", "HTTP_MONITOR_NAME", "DOCKER_CONTAINER_NAME", "DOCKER_FULL_IMAGE_NAME", "DOCKER_IMAGE_VERSION", "DOCKER_STRIPPED_IMAGE_NAME", "ESXI_HOST_HARDWARE_MODEL", "ESXI_HOST_HARDWARE_VENDOR", "ESXI_HOST_NAME", "ESXI_HOST_CLUSTER_NAME", "ESXI_HOST_PRODUCT_NAME", "ESXI_HOST_PRODUCT_VERSION", "NAME_OF_COMPUTE_NODE", "EC2_INSTANCE_NAME", "EC2_INSTANCE_AMI_ID", "EC2_INSTANCE_BEANSTALK_ENV_NAME", "EC2_INSTANCE_AWS_INSTANCE_TYPE", "EC2_INSTANCE_ID", "EC2_INSTANCE_PRIVATE_HOST_NAME", "EC2_INSTANCE_PUBLIC_HOST_NAME", "EC2_INSTANCE_AWS_SECURITY_GROUP", "OPENSTACK_VM_INSTANCE_TYPE", "OPENSTACK_VM_NAME", "OPENSTACK_VM_SECURITY_GROUP", "VMWARE_VM_NAME", "GOOGLE_COMPUTE_INSTANCE_ID", "GOOGLE_COMPUTE_INSTANCE_NAME", "GOOGLE_COMPUTE_INSTANCE_MACHINE_TYPE", "GOOGLE_COMPUTE_INSTANCE_PROJECT", "GOOGLE_COMPUTE_INSTANCE_PROJECT_ID", "AWS_AVAILABILITY_ZONE_NAME", "AZURE_REGION_NAME", "OPENSTACK_REGION_NAME", "OPENSTACK_AVAILABILITY_ZONE_NAME", "GEOLOCATION_SITE_NAME", "VMWARE_DATACENTER_NAME", "GOOGLE_CLOUD_PLATFORM_ZONE_NAME", "AWS_AUTO_SCALING_GROUP_NAME", "AWS_CLASSIC_LOAD_BALANCER_NAME", "AWS_APPLICATION_LOAD_BALANCER_NAME", "AWS_NETWORK_LOAD_BALANCER_NAME", "AWS_RELATIONAL_DATABASE_SERVICE_NAME", "AWS_RELATIONAL_DATABASE_SERVICE_INSTANCE_CLASS", "AWS_RELATIONAL_DATABASE_SERVICE_ENDPOINT", "AWS_RELATIONAL_DATABASE_SERVICE_ENGINE", "AWS_RELATIONAL_DATABASE_SERVICE_DB_NAME", "AZURE_SCALE_SET_NAME", "AZURE_VM_NAME", "OPENSTACK_PROJECT_NAME", "HOST_GROUP_NAME", "AWS_ACCOUNT_ID", "AWS_ACCOUNT_NAME", "OPENSTACK_ACCOUNT_NAME", "OPENSTACK_ACCOUNT_PROJECT_NAME", "CLOUD_FOUNDRY_ORG_NAME", "CLOUD_FOUNDRY_FOUNDATION_NAME", "APPMON_SERVER_NAME", "APPMON_SYSTEM_PROFILE_NAME", "QUEUE_NAME", "QUEUE_VENDOR"}

	err := decoder.DecodeAll(map[string]any{
		"enum_value":         &me.EnumValue,
		"tag":                &me.Tag,
		"key":                &me.Key,
		"dynamic_key":        &me.DynamicKey,
		"case_sensitive":     &me.CaseSensitive,
		"integer_value":      &me.IntegerValue,
		"entity_id":          &me.EntityID,
		"operator":           &me.Operator,
		"string_value":       &me.StringValue,
		"dynamic_key_source": &me.DynamicKeySource,
	})

	if me.CaseSensitive == nil && contains(string(me.Key), caseSensitiveIn) && me.Operator != Operators.Exists && me.Operator != Operators.NotExists {
		me.CaseSensitive = opt.NewBool(false)
	}
	return err
}

func contains(s string, array []string) bool {
	for _, str := range array {
		if s == str {
			return true
		}
	}
	return false
}
