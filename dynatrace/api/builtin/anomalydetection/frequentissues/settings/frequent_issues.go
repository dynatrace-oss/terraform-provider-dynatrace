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

package frequentissues

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type FrequentIssues struct {
	DetectApps  bool  `json:"detectFrequentIssuesInApplications"`            // Detect frequent issues within applications
	DetectEnv   *bool `json:"detectFrequentIssuesInEnvironment,omitempty"`   // Events raised at this level typically occur when no specific topological entity is applicable, often based on data such as logs and metrics. This does not impact the detection of issues within applications, transactions, services, or infrastructure.
	DetectInfra bool  `json:"detectFrequentIssuesInInfrastructure"`          // Detect frequent issues within infrastructure
	DetectTxn   bool  `json:"detectFrequentIssuesInTransactionsAndServices"` // Detect frequent issues within transactions and services
}

func (me *FrequentIssues) Name() string {
	return "frequent_issues"
}

func (me *FrequentIssues) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"detect_apps": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Detect frequent issues within applications, enabled (`true`) or disabled (`false`)",
		},
		"detect_env": {
			Type:        schema.TypeBool,
			Description: "Events raised at this level typically occur when no specific topological entity is applicable, often based on data such as logs and metrics. This does not impact the detection of issues within applications, transactions, services, or infrastructure.",
			Optional:    true, // nullable
		},
		"detect_infra": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Detect frequent issues within infrastructure, enabled (`true`) or disabled (`false`)",
		},
		"detect_txn": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Detect frequent issues within transactions and services, enabled (`true`) or disabled (`false`)",
		},
	}
}

func (me *FrequentIssues) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"detect_apps":  me.DetectApps,
		"detect_env":   me.DetectEnv,
		"detect_infra": me.DetectInfra,
		"detect_txn":   me.DetectTxn,
	})
}

func (me *FrequentIssues) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"detect_apps":  &me.DetectApps,
		"detect_env":   &me.DetectEnv,
		"detect_infra": &me.DetectInfra,
		"detect_txn":   &me.DetectTxn,
	})
}
