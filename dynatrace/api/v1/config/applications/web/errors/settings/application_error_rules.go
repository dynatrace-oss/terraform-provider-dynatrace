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

package errors

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Rules represents configuration of error rules in the web application
type Rules struct {
	WebApplicationID                         string           `json:"-"` // The EntityID of the the WebApplication
	Name                                     string           `json:"-"`
	IgnoreJavaScriptErrorsInApdexCalculation bool             `json:"ignoreJavaScriptErrorsInApdexCalculation"` // Exclude (`true`) or include (`false`) JavaScript errors in Apdex calculation
	IgnoreHttpErrorsInApdexCalculation       bool             `json:"ignoreHttpErrorsInApdexCalculation"`       // Exclude (`true`) or include (`false`) HTTP errors listed in **httpErrorRules** in Apdex calculation
	IgnoreCustomErrorsInApdexCalculation     bool             `json:"ignoreCustomErrorsInApdexCalculation"`     // (Field has overlap with `dynatrace_web_app_custom_errors`) Exclude (`true`) or include (`false`) custom errors listed in **customErrorRules** in Apdex calculation
	HTTPErrors                               HTTPErrorRules   `json:"httpErrorRules"`                           // An ordered list of HTTP errors.\n\n Rules are evaluated from top to bottom; the first matching rule applies
	CustomErrors                             CustomErrorRules `json:"customErrorRules"`                         // (Field has overlap with `dynatrace_web_app_custom_errors`) An ordered list of custom errors.\n\n Rules are evaluated from top to bottom; the first matching rule applies
}

func (me *Rules) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"web_application_id": {
			Type:        schema.TypeString,
			Description: "The EntityID of the the WebApplication",
			Optional:    true,
		},
		"ignore_js_errors_apdex": {
			Type:        schema.TypeBool,
			Description: "Exclude (`true`) or include (`false`) JavaScript errors in Apdex calculation",
			Optional:    true,
		},
		"ignore_http_errors_apdex": {
			Type:        schema.TypeBool,
			Description: "Exclude (`true`) or include (`false`) HTTP errors listed in **httpErrorRules** in Apdex calculation",
			Optional:    true,
		},
		"ignore_custom_errors_apdex": {
			Type:        schema.TypeBool,
			Description: "Exclude (`true`) or include (`false`) custom errors listed in **customErrorRules** in Apdex calculation",
			Optional:    true,
		},
		"http_errors": {
			Type:        schema.TypeList,
			Description: "An ordered list of HTTP errors.\n\n Rules are evaluated from top to bottom; the first matching rule applies",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(HTTPErrorRules).Schema()},
		},
		"custom_errors": {
			Type:        schema.TypeList,
			Description: "An ordered list of HTTP errors.\n\n Rules are evaluated from top to bottom; the first matching rule applies",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(CustomErrorRules).Schema()},
		},
	}
}

func (me *Rules) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"web_application_id":         me.WebApplicationID,
		"ignore_js_errors_apdex":     me.IgnoreJavaScriptErrorsInApdexCalculation,
		"ignore_http_errors_apdex":   me.IgnoreHttpErrorsInApdexCalculation,
		"ignore_custom_errors_apdex": me.IgnoreCustomErrorsInApdexCalculation,
		"http_errors":                me.HTTPErrors,
		"custom_errors":              me.CustomErrors,
	})
}

func (me *Rules) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"web_application_id":         &me.WebApplicationID,
		"ignore_js_errors_apdex":     &me.IgnoreJavaScriptErrorsInApdexCalculation,
		"ignore_http_errors_apdex":   &me.IgnoreHttpErrorsInApdexCalculation,
		"ignore_custom_errors_apdex": &me.IgnoreCustomErrorsInApdexCalculation,
		"http_errors":                &me.HTTPErrors,
		"custom_errors":              &me.CustomErrors,
	})
}

func (me *Rules) Store() ([]byte, error) {
	var data []byte
	var err error
	if data, err = json.Marshal(me); err != nil {
		return nil, err
	}
	m := map[string]json.RawMessage{}
	if err = json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	if data, err = json.Marshal(me.Name); err != nil {
		return nil, err
	}
	m["name"] = data
	if data, err = json.Marshal(me.WebApplicationID); err != nil {
		return nil, err
	}
	m["webApplicationID"] = data
	return json.MarshalIndent(m, "", "  ")
}

func (me *Rules) Load(data []byte) error {
	if err := json.Unmarshal(data, &me); err != nil {
		return err
	}

	c := struct {
		Name             string `json:"name"`
		WebApplicationID string `json:"webApplicationID"`
	}{}
	if err := json.Unmarshal(data, &c); err != nil {
		return err
	}
	me.Name = c.Name
	me.WebApplicationID = c.WebApplicationID

	return nil
}
