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

package teams

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Responsibilities struct {
	Development    bool `json:"development"`    // Responsible for developing and maintaining high quality software. Development teams are responsible for making code changes to address performance regressions, errors, or security vulnerabilities.
	Infrastructure bool `json:"infrastructure"` // Responsible for the administration, management, and support of the IT infrastructure including physical servers, virtualization, and cloud. Teams with infrastructure responsibility are responsible for addressing hardware issues, resource limits, and operating system vulnerabilities.
	LineOfBusiness bool `json:"lineOfBusiness"` // Responsible for ensuring that applications in development align with business needs and meet the usability requirements of users, stakeholders, customers, and external partners. Teams with line of business responsibility are responsible for understanding the customer experience and how it affects business goals.
	Operations     bool `json:"operations"`     // Responsible for deploying and managing software, with a focus on high availability and performance. Teams with operations responsibilities needs to understand the impact, priority, and team responsible for addressing problems detected by Dynatrace.
	Security       bool `json:"security"`       // Responsible for the security posture of the organization. Teams with security responsibility must understand the impact, priority, and team responsible for addressing security vulnerabilities.
}

func (me *Responsibilities) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"development": {
			Type:        schema.TypeBool,
			Description: "Responsible for developing and maintaining high quality software. Development teams are responsible for making code changes to address performance regressions, errors, or security vulnerabilities.",
			Required:    true,
		},
		"infrastructure": {
			Type:        schema.TypeBool,
			Description: "Responsible for the administration, management, and support of the IT infrastructure including physical servers, virtualization, and cloud. Teams with infrastructure responsibility are responsible for addressing hardware issues, resource limits, and operating system vulnerabilities.",
			Required:    true,
		},
		"line_of_business": {
			Type:        schema.TypeBool,
			Description: "Responsible for ensuring that applications in development align with business needs and meet the usability requirements of users, stakeholders, customers, and external partners. Teams with line of business responsibility are responsible for understanding the customer experience and how it affects business goals.",
			Required:    true,
		},
		"operations": {
			Type:        schema.TypeBool,
			Description: "Responsible for deploying and managing software, with a focus on high availability and performance. Teams with operations responsibilities needs to understand the impact, priority, and team responsible for addressing problems detected by Dynatrace.",
			Required:    true,
		},
		"security": {
			Type:        schema.TypeBool,
			Description: "Responsible for the security posture of the organization. Teams with security responsibility must understand the impact, priority, and team responsible for addressing security vulnerabilities.",
			Required:    true,
		},
	}
}

func (me *Responsibilities) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"development":      me.Development,
		"infrastructure":   me.Infrastructure,
		"line_of_business": me.LineOfBusiness,
		"operations":       me.Operations,
		"security":         me.Security,
	})
}

func (me *Responsibilities) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"development":      &me.Development,
		"infrastructure":   &me.Infrastructure,
		"line_of_business": &me.LineOfBusiness,
		"operations":       &me.Operations,
		"security":         &me.Security,
	})
}
