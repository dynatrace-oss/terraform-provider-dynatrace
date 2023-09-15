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

package preferences

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// The configuration of the user
type Settings struct {
	CertificateManagementEnabled   bool   `json:"certificateManagementEnabled,omitempty"`
	CertificateManagementPossible  bool   `json:"certificateManagementPossible,omitempty"`
	SupportSendBilling             bool   `json:"supportSendBilling"`                 // If true, usage and billing information will be reported.
	SuppressNonBillingRelevantData bool   `json:"suppressNonBillingRelevantData"`     // If true, usage and billing information will NOT be reported.
	SupportSendClusterHealth       bool   `json:"supportSendClusterHealth"`           // If true, Dynatrace cluster health will be reported.
	SupportSendEvents              bool   `json:"supportSendEvents,omitempty"`        // If true, Dynatrace cluster health and OneAgent events will be reported.
	SupportAllowRemoteAccess       bool   `json:"supportAllowRemoteAccess,omitempty"` // If true, audited remote-access to your Dynatrace configuration is allowed settings.
	RemoteAccessOnDemandOnly       bool   `json:"remoteAccessOnDemandOnly,omitempty"` // If true, audited access to your Dynatrace cluster is allowed by approved Dynatrace employees otherwise by privileged Dynatrace employees.
	CommunityCreateUser            bool   `json:"communityCreateUser,omitempty"`      // If true, each new user will get an invitation to set up a Dynatrace user account to access Dynatrace support resources user upon first login.
	CommunityExternalSearch        bool   `json:"communityExternalSearch,omitempty"`  // If true, results from Documentation are included in platform search.
	RuxitMonitorsRuxit             bool   `json:"ruxitMonitorsRuxit"`                 // If true, Dynatrace OneAgent monitors Dynatrace.
	TelemetrySharing               bool   `json:"telemetrySharing,omitempty"`
	HelpChatEnabled                bool   `json:"helpChatEnabled,omitempty"`             // If true, live, in-product assistance with our Product Experts is enabled.
	ReadOnlyRemoteAccessAllowed    bool   `json:"readOnlyRemoteAccessAllowed,omitempty"` // If true, audited, read-only remote access to your Dynatrace configuration settings is allowed.
	OriginalConfig                 string `json:"-"`
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"certificate_management_enabled": {
			Type:        schema.TypeBool,
			Description: "Default: `false`",
			Optional:    true,
		},
		"certificate_management_possible": {
			Type:        schema.TypeBool,
			Description: "Default: `true`",
			Optional:    true,
			Default:     true,
		},
		"support_send_billing": {
			Type:        schema.TypeBool,
			Description: "If true, usage and billing information will be reported.",
			Required:    true,
		},
		"suppress_non_billing_relevant_data": {
			Type:        schema.TypeBool,
			Description: "If true, usage and billing information will NOT be reported.",
			Required:    true,
		},
		"support_send_cluster_health": {
			Type:        schema.TypeBool,
			Description: "If true, Dynatrace cluster health will be reported.",
			Required:    true,
		},
		"support_send_events": {
			Type:        schema.TypeBool,
			Description: "If true, Dynatrace cluster health and OneAgent events will be reported. Default: `true`",
			Optional:    true,
			Default:     true,
		},
		"support_allow_remote_access": {
			Type:        schema.TypeBool,
			Description: "If true, audited remote-access to your Dynatrace configuration is allowed settings. Default: `true`",
			Optional:    true,
			Default:     true,
		},
		"remote_access_on_demand_only": {
			Type:        schema.TypeBool,
			Description: "If true, audited access to your Dynatrace cluster is allowed by approved Dynatrace employees otherwise by privileged Dynatrace employees. Default: `false`",
			Optional:    true,
		},
		"community_create_user": {
			Type:        schema.TypeBool,
			Description: "If true, each new user will get an invitation to set up a Dynatrace user account to access Dynatrace support resources user upon first login. Default: `false`",
			Optional:    true,
		},
		"community_external_search": {
			Type:        schema.TypeBool,
			Description: "If true, results from Documentation are included in platform search. Default: `false`",
			Optional:    true,
		},
		"ruxit_monitors_ruxit": {
			Type:        schema.TypeBool,
			Description: "If true, Dynatrace OneAgent monitors Dynatrace. Default: `true`",
			Required:    true,
		},
		"telemetry_sharing": {
			Type:        schema.TypeBool,
			Description: "Default: `false`",
			Optional:    true,
		},
		"help_chat_enabled": {
			Type:        schema.TypeBool,
			Description: "If true, live, in-product assistance with our Product Experts is enabled. Default: `false`",
			Optional:    true,
		},
		"read_only_remote_access_allowed": {
			Type:        schema.TypeBool,
			Description: "If true, audited, read-only remote access to your Dynatrace configuration settings is allowed. Default: `false`",
			Optional:    true,
		},
		"original_config": {
			Type:        schema.TypeString,
			Description: "For internal use: original config in JSON format to be used on destroy",
			Optional:    true,
			Computed:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"certificate_management_enabled":     me.CertificateManagementEnabled,
		"certificate_management_possible":    me.CertificateManagementPossible,
		"support_send_billing":               me.SupportSendBilling,
		"suppress_non_billing_relevant_data": me.SuppressNonBillingRelevantData,
		"support_send_cluster_health":        me.SupportSendClusterHealth,
		"support_send_events":                me.SupportSendEvents,
		"support_allow_remote_access":        me.SupportAllowRemoteAccess,
		"remote_access_on_demand_only":       me.RemoteAccessOnDemandOnly,
		"community_create_user":              me.CommunityCreateUser,
		"community_external_search":          me.CommunityExternalSearch,
		"ruxit_monitors_ruxit":               me.RuxitMonitorsRuxit,
		"telemetry_sharing":                  me.TelemetrySharing,
		"help_chat_enabled":                  me.HelpChatEnabled,
		"read_only_remote_access_allowed":    me.ReadOnlyRemoteAccessAllowed,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"certificate_management_enabled":     &me.CertificateManagementEnabled,
		"certificate_management_possible":    &me.CertificateManagementPossible,
		"support_send_billing":               &me.SupportSendBilling,
		"suppress_non_billing_relevant_data": &me.SuppressNonBillingRelevantData,
		"support_send_cluster_health":        &me.SupportSendClusterHealth,
		"support_send_events":                &me.SupportSendEvents,
		"support_allow_remote_access":        &me.SupportAllowRemoteAccess,
		"remote_access_on_demand_only":       &me.RemoteAccessOnDemandOnly,
		"community_create_user":              &me.CommunityCreateUser,
		"community_external_search":          &me.CommunityExternalSearch,
		"ruxit_monitors_ruxit":               &me.RuxitMonitorsRuxit,
		"telemetry_sharing":                  &me.TelemetrySharing,
		"help_chat_enabled":                  &me.HelpChatEnabled,
		"read_only_remote_access_allowed":    &me.ReadOnlyRemoteAccessAllowed,
	})
}
