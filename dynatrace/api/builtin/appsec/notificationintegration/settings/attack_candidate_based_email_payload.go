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

package notificationintegration

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AttackCandidateBasedEmailPayload struct {
	Body    string `json:"body"`    // The template of the email notifications.  \nIn case a value of a security problem is not set, the placeholder will be replaced by an empty string.. **Note:** Security notifications contain sensitive information. Excessive usage of placeholders in the body might leak information to untrusted parties.  \n  \nAvailable placeholders:  \n**{AttackDisplayId}**: The unique identifier assigned by Dynatrace, for example: \"A-1234\".  \n**{Title}**: Location of the attack, for example: \"com.dynatrace.Class.method():120\"  \n**{Type}**: The type of attack, for example: \"SQL Injection\".  \n**{AttackUrl}**: URL of the attack in Dynatrace.  \n**{ProcessGroupId}**: Details about the process group attacked.  \n**{EntryPoint}**: The entry point of the attack into the process, for example: \"/login\". Can be empty.  \n**{Status}**: The status of the attack, for example: \"Exploited\"  \n**{Timestamp}**: When the attack happened.  \n**{VulnerabilityName}**: Name of the associated code-level vulnerability, for example: \"InMemoryDatabaseCaller.getAccountInfo():51\". Can be empty.
	Subject string `json:"subject"` // The subject of the email notifications.. Available placeholders:  \n**{AttackDisplayId}**: The unique identifier assigned by Dynatrace, for example, \"A-1234\".  \n**{Title}**: Location of the attack, for example: \"com.dynatrace.Class.method():120\"  \n**{Type}**: The type of attack, for example: \"SQL Injection\".  \n**{AttackUrl}**: URL of the attack in Dynatrace.  \n**{ProcessGroupId}**: Details about the process group attacked.  \n**{EntryPoint}**: The entry point of the attack into the process, for example: \"/login\". Can be empty.  \n**{Status}**: The status of the attack, for example: \"Exploited\"  \n**{Timestamp}**: When the attack happened.  \n**{VulnerabilityName}**: Name of the associated code-level vulnerability, for example: \"InMemoryDatabaseCaller.getAccountInfo():51\". Can be empty.
}

func (me *AttackCandidateBasedEmailPayload) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"body": {
			Type:        schema.TypeString,
			Description: "The template of the email notifications.  \nIn case a value of a security problem is not set, the placeholder will be replaced by an empty string.. **Note:** Security notifications contain sensitive information. Excessive usage of placeholders in the body might leak information to untrusted parties.  \n  \nAvailable placeholders:  \n**{AttackDisplayId}**: The unique identifier assigned by Dynatrace, for example: \"A-1234\".  \n**{Title}**: Location of the attack, for example: \"com.dynatrace.Class.method():120\"  \n**{Type}**: The type of attack, for example: \"SQL Injection\".  \n**{AttackUrl}**: URL of the attack in Dynatrace.  \n**{ProcessGroupId}**: Details about the process group attacked.  \n**{EntryPoint}**: The entry point of the attack into the process, for example: \"/login\". Can be empty.  \n**{Status}**: The status of the attack, for example: \"Exploited\"  \n**{Timestamp}**: When the attack happened.  \n**{VulnerabilityName}**: Name of the associated code-level vulnerability, for example: \"InMemoryDatabaseCaller.getAccountInfo():51\". Can be empty.",
			Required:    true,
		},
		"subject": {
			Type:        schema.TypeString,
			Description: "The subject of the email notifications.. Available placeholders:  \n**{AttackDisplayId}**: The unique identifier assigned by Dynatrace, for example, \"A-1234\".  \n**{Title}**: Location of the attack, for example: \"com.dynatrace.Class.method():120\"  \n**{Type}**: The type of attack, for example: \"SQL Injection\".  \n**{AttackUrl}**: URL of the attack in Dynatrace.  \n**{ProcessGroupId}**: Details about the process group attacked.  \n**{EntryPoint}**: The entry point of the attack into the process, for example: \"/login\". Can be empty.  \n**{Status}**: The status of the attack, for example: \"Exploited\"  \n**{Timestamp}**: When the attack happened.  \n**{VulnerabilityName}**: Name of the associated code-level vulnerability, for example: \"InMemoryDatabaseCaller.getAccountInfo():51\". Can be empty.",
			Required:    true,
		},
	}
}

func (me *AttackCandidateBasedEmailPayload) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"body":    me.Body,
		"subject": me.Subject,
	})
}

func (me *AttackCandidateBasedEmailPayload) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"body":    &me.Body,
		"subject": &me.Subject,
	})
}
