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

type SecurityProblemBasedEmailPayload struct {
	Body    string `json:"body"`    // The template of the email notifications.  \nIn case a value of a security problem is not set, the placeholder will be replaced by an empty string.. **Note:** Security notifications contain sensitive information. Excessive usage of placeholders in the description might leak information to untrusted parties.  \n  \nAvailable placeholders:  \n**{SecurityProblemId}**: The unique identifier assigned by Dynatrace, for example, \"S-1234\".  \n**{Title}**: A short summary of the type of vulnerability that was found, for example, \"Remote Code Execution\".  \n**{Description}**: A more detailed description of the vulnerability.  \n**{CvssScore}**: CVSS score of the identified vulnerability, for example, \"10.0\". Can be empty. \n**{DavisSecurityScore}**: [Davis Security Score](https://www.dynatrace.com/support/help/how-to-use-dynatrace/application-security/davis-security-score/) is an enhanced risk-calculation score based on the CVSS, for example, \"10.0\".  \n**{Severity}**: The security problem severity, for example, \"Critical\" or \"Medium\".  \n**{SecurityProblemUrl}**: URL of the security problem in Dynatrace.  \n**{AffectedEntities}**: Details about the entities affected by the security problem in a json array.  \n**{ManagementZones}**: Comma-separated list of all management zones affected by the vulnerability at the time of detection.  \n**{Tags}**: Comma-separated list of tags that are defined for a vulnerability's affected entities. For example: \"PROD, owner:John\". Assign the tag's key in square brackets: **{Tags[key]}** to get all associated values. For example: \"{Tags[owner]}\" will be resolved as \"John\". Tags without an assigned value will be resolved as empty string.  \n**{Exposed}**: Describes whether one or more affected process is exposed to the public Internet. Can be \"true\" or \"false\".  \n**{DataAssetsReachable}**: Describes whether one or more affected process can reach data assets. Can be \"true\" or \"false\".  \n**{ExploitAvailable}**: Describes whether there's an exploit available for the vulnerability. Can be \"true\" or \"false\".
	Subject string `json:"subject"` // The subject of the email notifications.. Available placeholders:  \n**{SecurityProblemId}**: The unique identifier assigned by Dynatrace, for example, \"S-1234\".  \n**{Title}**: A short summary of the type of vulnerability that was found, for example, \"Remote Code Execution\".  \n**{CvssScore}**: CVSS score of the identified vulnerability, for example, \"10.0\". Can be empty. \n**{DavisSecurityScore}**: [Davis Security Score](https://www.dynatrace.com/support/help/how-to-use-dynatrace/application-security/davis-security-score/) is an enhanced risk-calculation score based on the CVSS, for example, \"10.0\".  \n**{Severity}**: The security problem severity, for example, \"Critical\" or \"Medium\".  \n**{SecurityProblemUrl}**: URL of the security problem in Dynatrace.  \n**{Exposed}**: Describes whether one or more affected process is exposed to the public Internet. Can be \"true\" or \"false\".  \n**{DataAssetsReachable}**: Describes whether one or more affected process can reach data assets. Can be \"true\" or \"false\".  \n**{ExploitAvailable}**: Describes whether there's an exploit available for the vulnerability. Can be \"true\" or \"false\".
}

func (me *SecurityProblemBasedEmailPayload) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"body": {
			Type:        schema.TypeString,
			Description: "The template of the email notifications.  \nIn case a value of a security problem is not set, the placeholder will be replaced by an empty string.. **Note:** Security notifications contain sensitive information. Excessive usage of placeholders in the description might leak information to untrusted parties.  \n  \nAvailable placeholders:  \n**{SecurityProblemId}**: The unique identifier assigned by Dynatrace, for example, \"S-1234\".  \n**{Title}**: A short summary of the type of vulnerability that was found, for example, \"Remote Code Execution\".  \n**{Description}**: A more detailed description of the vulnerability.  \n**{CvssScore}**: CVSS score of the identified vulnerability, for example, \"10.0\". Can be empty. \n**{DavisSecurityScore}**: [Davis Security Score](https://www.dynatrace.com/support/help/how-to-use-dynatrace/application-security/davis-security-score/) is an enhanced risk-calculation score based on the CVSS, for example, \"10.0\".  \n**{Severity}**: The security problem severity, for example, \"Critical\" or \"Medium\".  \n**{SecurityProblemUrl}**: URL of the security problem in Dynatrace.  \n**{AffectedEntities}**: Details about the entities affected by the security problem in a json array.  \n**{ManagementZones}**: Comma-separated list of all management zones affected by the vulnerability at the time of detection.  \n**{Tags}**: Comma-separated list of tags that are defined for a vulnerability's affected entities. For example: \"PROD, owner:John\". Assign the tag's key in square brackets: **{Tags[key]}** to get all associated values. For example: \"{Tags[owner]}\" will be resolved as \"John\". Tags without an assigned value will be resolved as empty string.  \n**{Exposed}**: Describes whether one or more affected process is exposed to the public Internet. Can be \"true\" or \"false\".  \n**{DataAssetsReachable}**: Describes whether one or more affected process can reach data assets. Can be \"true\" or \"false\".  \n**{ExploitAvailable}**: Describes whether there's an exploit available for the vulnerability. Can be \"true\" or \"false\".",
			Required:    true,
		},
		"subject": {
			Type:        schema.TypeString,
			Description: "The subject of the email notifications.. Available placeholders:  \n**{SecurityProblemId}**: The unique identifier assigned by Dynatrace, for example, \"S-1234\".  \n**{Title}**: A short summary of the type of vulnerability that was found, for example, \"Remote Code Execution\".  \n**{CvssScore}**: CVSS score of the identified vulnerability, for example, \"10.0\". Can be empty. \n**{DavisSecurityScore}**: [Davis Security Score](https://www.dynatrace.com/support/help/how-to-use-dynatrace/application-security/davis-security-score/) is an enhanced risk-calculation score based on the CVSS, for example, \"10.0\".  \n**{Severity}**: The security problem severity, for example, \"Critical\" or \"Medium\".  \n**{SecurityProblemUrl}**: URL of the security problem in Dynatrace.  \n**{Exposed}**: Describes whether one or more affected process is exposed to the public Internet. Can be \"true\" or \"false\".  \n**{DataAssetsReachable}**: Describes whether one or more affected process can reach data assets. Can be \"true\" or \"false\".  \n**{ExploitAvailable}**: Describes whether there's an exploit available for the vulnerability. Can be \"true\" or \"false\".",
			Required:    true,
		},
	}
}

func (me *SecurityProblemBasedEmailPayload) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"body":    me.Body,
		"subject": me.Subject,
	})
}

func (me *SecurityProblemBasedEmailPayload) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"body":    &me.Body,
		"subject": &me.Subject,
	})
}
