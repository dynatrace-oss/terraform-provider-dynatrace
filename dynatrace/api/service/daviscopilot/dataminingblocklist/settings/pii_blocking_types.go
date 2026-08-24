/*
 * @license
 * Copyright 2026 Dynatrace LLC
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dataminingblocklist

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PiiBlockingTypes struct {
	CanadianSocialInsuranceNumber bool `json:"canadianSocialInsuranceNumber"` // Canadian social insurance number
	CreditCardNumber              bool `json:"creditCardNumber"`              // Credit card number
	EmailAddress                  bool `json:"emailAddress"`                  // Email address
	IbanBankAccount               bool `json:"ibanBankAccount"`               // IBAN bank account
	IpAddress                     bool `json:"ipAddress"`                     // IP address
	PhoneNumber                   bool `json:"phoneNumber"`                   // Phone number
	UrlQueryParameters            bool `json:"urlQueryParameters"`            // URL query parameters
	UsBankNumber                  bool `json:"usBankNumber"`                  // US bank number
	UsSocialSecurityNumber        bool `json:"usSocialSecurityNumber"`        // US social security number
}

func (me *PiiBlockingTypes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"canadian_social_insurance_number": {
			Type:        schema.TypeBool,
			Description: "Canadian social insurance number",
			Required:    true,
		},
		"credit_card_number": {
			Type:        schema.TypeBool,
			Description: "Credit card number",
			Required:    true,
		},
		"email_address": {
			Type:        schema.TypeBool,
			Description: "Email address",
			Required:    true,
		},
		"iban_bank_account": {
			Type:        schema.TypeBool,
			Description: "IBAN bank account",
			Required:    true,
		},
		"ip_address": {
			Type:        schema.TypeBool,
			Description: "IP address",
			Required:    true,
		},
		"phone_number": {
			Type:        schema.TypeBool,
			Description: "Phone number",
			Required:    true,
		},
		"url_query_parameters": {
			Type:        schema.TypeBool,
			Description: "URL query parameters",
			Required:    true,
		},
		"us_bank_number": {
			Type:        schema.TypeBool,
			Description: "US bank number",
			Required:    true,
		},
		"us_social_security_number": {
			Type:        schema.TypeBool,
			Description: "US social security number",
			Required:    true,
		},
	}
}

func (me *PiiBlockingTypes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"canadian_social_insurance_number": me.CanadianSocialInsuranceNumber,
		"credit_card_number":               me.CreditCardNumber,
		"email_address":                    me.EmailAddress,
		"iban_bank_account":                me.IbanBankAccount,
		"ip_address":                       me.IpAddress,
		"phone_number":                     me.PhoneNumber,
		"url_query_parameters":             me.UrlQueryParameters,
		"us_bank_number":                   me.UsBankNumber,
		"us_social_security_number":        me.UsSocialSecurityNumber,
	})
}

func (me *PiiBlockingTypes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"canadian_social_insurance_number": &me.CanadianSocialInsuranceNumber,
		"credit_card_number":               &me.CreditCardNumber,
		"email_address":                    &me.EmailAddress,
		"iban_bank_account":                &me.IbanBankAccount,
		"ip_address":                       &me.IpAddress,
		"phone_number":                     &me.PhoneNumber,
		"url_query_parameters":             &me.UrlQueryParameters,
		"us_bank_number":                   &me.UsBankNumber,
		"us_social_security_number":        &me.UsSocialSecurityNumber,
	})
}
