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

package smtp

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// The configuration of the user
type Settings struct {
	HostName                       string             `json:"hostName"`                                 // Host name
	Port                           int                `json:"port"`                                     // Integer value of port. Default: `25`
	UserName                       string             `json:"userName"`                                 // User name
	Password                       string             `json:"password"`                                 // Password
	IsPasswordConfigured           bool               `json:"isPasswordConfigured,omitempty"`           // If true, a password has been configured. Default: `false`.
	ConnectionSecurity             ConnectionSecurity `json:"connectionSecurity"`                       // Connection security, possible values: `NO_ENCRYPTION`, `OPTIONAL_STARTTLS`, `REQUIRE_STARTTLS`, `TLS`. Default: `NO_ENCRYPTION`
	SenderEmailAddress             string             `json:"senderEmailAddress"`                       // Sender email address
	AllowFallbackViaMissionControl bool               `json:"allowFallbackViaMissionControl,omitempty"` // If true, we will send e-mails via Mission Control in case of problems with SMTP server.
	UseSmtpServer                  bool               `json:"useSmtpServer,omitempty"`                  // If true, we will send e-mails via SMTP server.
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"host_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Host Name",
		},
		"port": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     25,
			Description: "Integer value of port. Default: `25`",
		},
		"user_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "User Name",
		},
		"password": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Password",
			Sensitive:   true,
		},
		"is_password_configured": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "If true, a password has been configured. Default: `false`.",
		},
		"connection_security": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     ConnectionSecurities.NoEncryption,
			Description: "Connection security, possible values: `NO_ENCRYPTION`, `OPTIONAL_STARTTLS`, `REQUIRE_STARTTLS`, `TLS`. Default: `NO_ENCRYPTION`",
		},
		"sender_email_address": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Sender email address",
		},
		"allow_fallback_via_mission_control": {
			Type:         schema.TypeBool,
			Optional:     true,
			AtLeastOneOf: []string{"allow_fallback_via_mission_control", "use_smtp_server"},
			Description:  "If true, we will send e-mails via Mission Control in case of problems with SMTP server.",
		},
		"use_smtp_server": {
			Type:         schema.TypeBool,
			Optional:     true,
			AtLeastOneOf: []string{"allow_fallback_via_mission_control", "use_smtp_server"},
			Description:  "If true, we will send e-mails via SMTP server.",
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"host_name":                          me.HostName,
		"port":                               me.Port,
		"user_name":                          me.UserName,
		"is_password_configured":             me.IsPasswordConfigured,
		"connection_security":                me.ConnectionSecurity,
		"sender_email_address":               me.SenderEmailAddress,
		"allow_fallback_via_mission_control": me.AllowFallbackViaMissionControl,
		"use_smtp_server":                    me.UseSmtpServer,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"host_name":                          &me.HostName,
		"port":                               &me.Port,
		"user_name":                          &me.UserName,
		"password":                           &me.Password,
		"is_password_configured":             &me.IsPasswordConfigured,
		"connection_security":                &me.ConnectionSecurity,
		"sender_email_address":               &me.SenderEmailAddress,
		"allow_fallback_via_mission_control": &me.AllowFallbackViaMissionControl,
		"use_smtp_server":                    &me.UseSmtpServer,
	})
}
