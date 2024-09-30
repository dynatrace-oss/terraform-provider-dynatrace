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

package aws

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export/sensitive"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AWSCredentialsConfig Configuration of an AWS credentials.
type AWSCredentialsConfig struct {
	Label                            string                     `json:"label"`                            // The name of the credentials.
	TaggedOnly                       *bool                      `json:"taggedOnly"`                       // Monitor only resources which have specified AWS tags (`true`) or all resources (`false`).
	AuthenticationData               *AWSAuthenticationData     `json:"authenticationData"`               // A credentials for the AWS authentication.
	PartitionType                    PartitionType              `json:"partitionType"`                    // The type of the AWS partition.
	TagsToMonitor                    []*AWSConfigTag            `json:"tagsToMonitor"`                    // A list of AWS tags to be monitored.  You can specify up to 10 tags.  Only applicable when the **taggedOnly** parameter is set to `true`.
	ConnectionStatus                 *ConnectionStatus          `json:"connectionStatus,omitempty"`       // The status of the connection to the AWS environment.   * `CONNECTED`: There was a connection within last 10 minutes.  * `DISCONNECTED`: A problem occurred with establishing connection using these credentials. Check whether the data is correct.  * `UNINITIALIZED`: The successful connection has never been established for these credentials.
	CredentialsEnabled               bool                       `json:"credentialsEnabled"`               // Enable monitoring of specified AWS credentials
	RunningOnDynatraceInfrastructure bool                       `json:"runningOnDynatraceInfrastructure"` // Run credentials on Dynatrace infrastructure
	RemoveDefaults                   bool                       `json:"-"`
	Unknowns                         map[string]json.RawMessage `json:"-"`
}

func (awscc *AWSCredentialsConfig) Name() string {
	return awscc.Label
}

func (awscc *AWSCredentialsConfig) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"label": {
			Type:        schema.TypeString,
			Description: "The name of the credentials",
			Optional:    true,
		},
		"tagged_only": {
			Type:        schema.TypeBool,
			Description: "Monitor only resources which have specified AWS tags (`true`) or all resources (`false`)",
			Required:    true,
		},
		"partition_type": {
			Type:        schema.TypeString,
			Description: "The type of the AWS partition",
			Required:    true,
		},
		"authentication_data": {
			Type:        schema.TypeList,
			Description: "credentials for the AWS authentication",
			Required:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: new(AWSAuthenticationData).Schema(),
			},
		},
		"tags_to_monitor": {
			Type:        schema.TypeSet,
			Description: "AWS tags to be monitored. You can specify up to 10 tags. Only applicable when the **tagged_only** parameter is set to `true`",
			Optional:    true,
			MaxItems:    10,
			Elem:        &schema.Resource{Schema: new(AWSConfigTag).Schema()},
		},
		"supporting_services_to_monitor": {
			Type:             schema.TypeSet,
			Description:      "supporting services to be monitored",
			Optional:         true,
			Elem:             &schema.Resource{Schema: new(AWSSupportingServiceConfig).Schema()},
			DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool { return true },
			Deprecated:       "Managing supporting services directly within AWS Credentials has been deprecated within the REST API. This attribute just exists for backwards compatibility. It no longer has an effect. For managing services use the resource `dynatrace_aws_service`",
		},
		"supporting_services_managed_in_dynatrace": {
			Type:             schema.TypeBool,
			Description:      "If enabled (`true`) the attribute `supporting_services` will not get synchronized with Dynatrace. You will be able to manage them via WebUI without interference by Terraform.",
			Optional:         true,
			Default:          false,
			DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool { return true },
			Deprecated:       "Supporting Services are no longer getting managed via this resource. Regardless of the value set here, this resource won't affect the supporting services during updates",
		},
		"remove_defaults": {
			Type:        schema.TypeBool,
			Description: "Instructs the provider to remove the supporting services Dynatrace applies by default to newly created AWS Credentials. Supporting Services applied by via `dynatrace_aws_service` subsequently won't get touched.\nNote: This attribute is only getting considered during creation of the resource. Changing it afterwards won't have an effect",
			Optional:    true,
			Default:     false,
		},
		"credentials_enabled": {
			Type:        schema.TypeBool,
			Description: "Enable monitoring of specified AWS credentials",
			Optional:    true,
			Default:     true,
		},
		"running_on_dynatrace_infrastructure": {
			Type:        schema.TypeBool,
			Description: "Run credentials on Dynatrace infrastructure",
			Optional:    true,
			Default:     false,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider",
			Optional:    true,
		},
	}
}

func (awscc *AWSCredentialsConfig) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(awscc.Unknowns) > 0 {
		for k, v := range awscc.Unknowns {
			m[k] = v
		}
	}
	delete(m, "metadata")
	delete(m, "id")
	delete(m, "connectionStatus")
	{
		rawMessage, err := json.Marshal(awscc.Label)
		if err != nil {
			return nil, err
		}
		m["label"] = rawMessage
	}
	{
		rawMessage, _ := json.Marshal([]string{})
		m["supportingServicesToMonitor"] = rawMessage
	}
	if rawMessage, err := json.Marshal(opt.Bool(awscc.TaggedOnly)); err == nil {
		m["taggedOnly"] = rawMessage
	} else {
		return nil, err
	}

	if awscc.AuthenticationData != nil {
		rawMessage, err := json.Marshal(awscc.AuthenticationData)
		if err != nil {
			return nil, err
		}
		m["authenticationData"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(awscc.PartitionType)
		if err != nil {
			return nil, err
		}
		m["partitionType"] = rawMessage
	}
	if awscc.TagsToMonitor != nil {
		rawMessage, err := json.Marshal(awscc.TagsToMonitor)
		if err != nil {
			return nil, err
		}
		m["tagsToMonitor"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(awscc.CredentialsEnabled)
		if err != nil {
			return nil, err
		}
		m["credentialsEnabled"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(awscc.RunningOnDynatraceInfrastructure)
		if err != nil {
			return nil, err
		}
		m["runningOnDynatraceInfrastructure"] = rawMessage
	}
	return json.Marshal(m)
}

// UnmarshalJSON provides custom JSON deserialization
func (awscc *AWSCredentialsConfig) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["label"]; found {
		if err := json.Unmarshal(v, &awscc.Label); err != nil {
			return err
		}
	}
	if v, found := m["taggedOnly"]; found {
		if err := json.Unmarshal(v, &awscc.TaggedOnly); err != nil {
			return err
		}
	} else {
		awscc.TaggedOnly = opt.NewBool(false)
	}
	if v, found := m["authenticationData"]; found {
		if err := json.Unmarshal(v, &awscc.AuthenticationData); err != nil {
			return err
		}
	}
	if v, found := m["partitionType"]; found {
		if err := json.Unmarshal(v, &awscc.PartitionType); err != nil {
			return err
		}
	}
	if v, found := m["tagsToMonitor"]; found {
		if err := json.Unmarshal(v, &awscc.TagsToMonitor); err != nil {
			return err
		}
	}
	if v, found := m["credentialsEnabled"]; found {
		if err := json.Unmarshal(v, &awscc.CredentialsEnabled); err != nil {
			return err
		}
	} else {
		awscc.CredentialsEnabled = true
	}
	if v, found := m["runningOnDynatraceInfrastructure"]; found {
		if err := json.Unmarshal(v, &awscc.RunningOnDynatraceInfrastructure); err != nil {
			return err
		}
	} else {
		awscc.RunningOnDynatraceInfrastructure = false
	}

	delete(m, "id")
	delete(m, "label")
	delete(m, "metadata")
	delete(m, "supportingServicesToMonitor")
	delete(m, "taggedOnly")
	delete(m, "connectionStatus")
	delete(m, "authenticationData")
	delete(m, "partitionType")
	delete(m, "tagsToMonitor")
	delete(m, "connectionStatus")
	delete(m, "credentialsEnabled")
	delete(m, "runningOnDynatraceInfrastructure")
	if len(m) > 0 {
		awscc.Unknowns = m
	}
	return nil
}

func (awscc *AWSCredentialsConfig) PrepareMarshalHCL(decoder hcl.Decoder) error {
	if awscc.AuthenticationData != nil {
		if awscc.AuthenticationData.KeyBasedAuthentication != nil {
			if awscc.AuthenticationData.KeyBasedAuthentication.SecretKey == nil {
				if secretKey, ok := decoder.GetOk("authentication_data.0.secret_key"); ok && len(secretKey.(string)) > 0 {
					awscc.AuthenticationData.KeyBasedAuthentication.SecretKey = opt.NewString(secretKey.(string))
				}
			}
		}
	}
	return nil
}

func (awscc *AWSCredentialsConfig) MarshalHCL(properties hcl.Properties) error {
	if awscc.Unknowns != nil {
		data, err := json.Marshal(awscc.Unknowns)
		if err != nil {
			return err
		}
		if err := properties.Encode("unknowns", string(data)); err != nil {
			return err
		}
	}
	if err := properties.Encode("label", awscc.Label); err != nil {
		return err
	}
	if err := properties.Encode("tagged_only", opt.Bool(awscc.TaggedOnly)); err != nil {
		return err
	}

	if err := properties.Encode("authentication_data", awscc.AuthenticationData); err != nil {
		return err
	}
	if err := sensitive.ConditionalIgnoreChangesSinglePlus(
		awscc.Schema(),
		&properties,
		[]string{"authentication_data"},
	); err != nil {
		return err
	}
	if err := properties.Encode("partition_type", awscc.PartitionType); err != nil {
		return err
	}
	if err := properties.Encode("tags_to_monitor", awscc.TagsToMonitor); err != nil {
		return err
	}
	if err := properties.Encode("credentials_enabled", awscc.CredentialsEnabled); err != nil {
		return err
	}
	if err := properties.Encode("running_on_dynatrace_infrastructure", awscc.RunningOnDynatraceInfrastructure); err != nil {
		return err
	}

	return nil
}

func (awscc *AWSCredentialsConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), awscc); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &awscc.Unknowns); err != nil {
			return err
		}
		delete(awscc.Unknowns, "supporting_services_to_monitor")
		delete(awscc.Unknowns, "label")
		delete(awscc.Unknowns, "tagged_only")
		delete(awscc.Unknowns, "authentication_data")
		delete(awscc.Unknowns, "partition_type")
		delete(awscc.Unknowns, "tags_to_monitor")
		if len(awscc.Unknowns) == 0 {
			awscc.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("label"); ok {
		awscc.Label = value.(string)
	}
	if value, ok := decoder.GetOk("tagged_only"); ok {
		awscc.TaggedOnly = opt.NewBool(value.(bool))
	}
	if value, ok := decoder.GetOk("credentials_enabled"); ok {
		awscc.CredentialsEnabled = value.(bool)
	}
	if value, ok := decoder.GetOk("running_on_dynatrace_infrastructure"); ok {
		awscc.RunningOnDynatraceInfrastructure = value.(bool)
	}
	if _, ok := decoder.GetOk("authentication_data.#"); ok {
		awscc.AuthenticationData = new(AWSAuthenticationData)
		if err := awscc.AuthenticationData.UnmarshalHCL(hcl.NewDecoder(decoder, "authentication_data", 0)); err != nil {
			return err
		}
	}
	if value, ok := decoder.GetOk("partition_type"); ok {
		awscc.PartitionType = PartitionType(value.(string))
	}
	if err := decoder.DecodeSlice("tags_to_monitor", &awscc.TagsToMonitor); err != nil {
		return err
	}
	if value, ok := decoder.GetOk("remove_defaults"); ok {
		awscc.RemoveDefaults = value.(bool)
	}
	if awscc.TaggedOnly != nil && *awscc.TaggedOnly && awscc.TagsToMonitor == nil {
		awscc.TagsToMonitor = []*AWSConfigTag{}
	}
	if awscc.TaggedOnly == nil || !(*awscc.TaggedOnly) {
		awscc.TagsToMonitor = nil
	}
	return nil
}

const credsNotProvided = "REST API didn't provide credential data"

func (me *AWSCredentialsConfig) FillDemoValues() []string {
	if me.AuthenticationData != nil && me.AuthenticationData.KeyBasedAuthentication != nil {
		me.AuthenticationData.KeyBasedAuthentication.SecretKey = opt.NewString("################")
		return []string{credsNotProvided}
	}
	return nil
}

func (me *AWSCredentialsConfig) Anonymize() {
	if me.AuthenticationData == nil {
		return
	}
	if me.AuthenticationData.RoleBasedAuthentication == nil {
		return
	}
	me.AuthenticationData.RoleBasedAuthentication.ExternalID = nil
}
