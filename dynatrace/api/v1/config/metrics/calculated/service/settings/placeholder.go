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

package service

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/metrics/calculated/service/settings/propagation"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Placeholder The custom placeholder to be used as a naming value pattern.
//
//	It enables you to extract a request attribute value or other request attribute and use it in the naming pattern.
type Placeholder struct {
	Name              string                     `json:"name"`                        // The name of the placeholder. Use it in the naming pattern as `{name}`.
	Aggregation       *Aggregation               `json:"aggregation,omitempty"`       // Which value of the request attribute must be used when it occurs across multiple child requests.  Only applicable for the `SERVICE_REQUEST_ATTRIBUTE` attribute, when **useFromChildCalls** is `true`.  For the `COUNT` aggregation, the **kind** field is not applicable.
	Attribute         Attribute                  `json:"attribute"`                   // The attribute to extract from. You can only use attributes of the **string** type.
	Source            *propagation.Source        `json:"source,omitempty"`            // Defines valid sources of request attributes for conditions or placeholders.
	Normalization     *Normalization             `json:"normalization,omitempty"`     // The format of the extracted string.
	DelimiterOrRegex  *string                    `json:"delimiterOrRegex,omitempty"`  // Depending on the **type** value:   * `REGEX_EXTRACTION`: The regular expression.   * `BETWEEN_DELIMITER`: The opening delimiter string to look for.   * All other values: The delimiter string to look for.
	UseFromChildCalls *bool                      `json:"useFromChildCalls,omitempty"` // If `true` request attribute will be taken from a child service call.   Only applicable for the `SERVICE_REQUEST_ATTRIBUTE` attribute. Defaults to `false`.
	RequestAttribute  *string                    `json:"requestAttribute,omitempty"`  // The request attribute to extract from.   Required if the **kind** value is `SERVICE_REQUEST_ATTRIBUTE`. Not applicable otherwise.
	EndDelimiter      *string                    `json:"endDelimiter,omitempty"`      // The closing delimiter string to look for.   Required if the **kind** value is `BETWEEN_DELIMITER`. Not applicable otherwise.
	Kind              Kind                       `json:"kind"`                        // The type of extraction.   Defines either usage of regular expression (`regex`) or the position of request attribute value to be extracted.  When the **attribute** is `SERVICE_REQUEST_ATTRIBUTE` attribute and **aggregation** is `COUNT`, needs to be set to `ORIGINAL_TEXT`
	Unknowns          map[string]json.RawMessage `json:"-"`
}

func (me *Placeholder) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the placeholder. Use it in the naming pattern as `{name}`",
		},
		"aggregation": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Which value of the request attribute must be used when it occurs across multiple child requests. Only applicable for the `SERVICE_REQUEST_ATTRIBUTE` attribute, when **useFromChildCalls** is `true`. For the `COUNT` aggregation, the **kind** field is not applicable. Possible values are `COUNT`, `FIRST` and `LAST`.",
		},
		"attribute": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The attribute to extract from. You can only use attributes of the **string** type. Possible values are `ACTOR_SYSTEM`, `AKKA_ACTOR_CLASS_NAME`, `AKKA_ACTOR_MESSAGE_TYPE`, `AKKA_ACTOR_PATH`, `APPLICATION_BUILD_VERSION`, `APPLICATION_RELEASE_VERSION`, `AZURE_FUNCTIONS_FUNCTION_NAME`, `AZURE_FUNCTIONS_SITE_NAME`, `CICS_PROGRAM_NAME`, `CICS_SYSTEM_ID`, `CICS_TASK_ID`, `CICS_TRANSACTION_ID`, `CICS_USER_ID`, `CPU_TIME`, `CTG_GATEWAY_URL`, `CTG_PROGRAM`, `CTG_SERVER_NAME`, `CTG_TRANSACTION_ID`, `CUSTOMSERVICE_CLASS`, `CUSTOMSERVICE_METHOD`, `DATABASE_CHILD_CALL_COUNT`, `DATABASE_CHILD_CALL_TIME`, `DATABASE_HOST`, `DATABASE_NAME`, `DATABASE_TYPE`, `DATABASE_URL`, `DISK_IO_TIME`, `ERROR_COUNT`, `ESB_APPLICATION_NAME`, `ESB_INPUT_TYPE`, `ESB_LIBRARY_NAME`, `ESB_MESSAGE_FLOW_NAME`, `EXCEPTION_CLASS`, `EXCEPTION_MESSAGE`, `FAILED_STATE`, `FAILURE_REASON`, `FLAW_STATE`, `HTTP_REQUEST_METHOD`, `HTTP_STATUS`, `HTTP_STATUS_CLASS`, `IMS_PROGRAM_NAME`, `IMS_TRANSACTION_ID`, `IMS_USER_ID`, `IO_TIME`, `IS_KEY_REQUEST`, `LAMBDA_COLDSTART`, `LOCK_TIME`, `MESSAGING_DESTINATION_TYPE`, `MESSAGING_IS_TEMPORARY_QUEUE`, `MESSAGING_QUEUE_NAME`, `MESSAGING_QUEUE_VENDOR`, `NETWORK_IO_TIME`, `NON_DATABASE_CHILD_CALL_COUNT`, `NON_DATABASE_CHILD_CALL_TIME`, `PROCESS_GROUP_NAME`, `PROCESS_GROUP_TAG`, `REMOTE_ENDPOINT`, `REMOTE_METHOD`, `REMOTE_SERVICE_NAME`, `REQUEST_NAME`, `REQUEST_TYPE`, `RESPONSE_TIME`, `RESPONSE_TIME_CLIENT`, `RMI_CLASS`, `RMI_METHOD`, `SERVICE_DISPLAY_NAME`, `SERVICE_NAME`, `SERVICE_PORT`, `SERVICE_PUBLIC_DOMAIN_NAME`, `SERVICE_REQUEST_ATTRIBUTE`, `SERVICE_TAG`, `SERVICE_TYPE`, `SERVICE_WEB_APPLICATION_ID`, `SERVICE_WEB_CONTEXT_ROOT`, `SERVICE_WEB_SERVER_NAME`, `SERVICE_WEB_SERVICE_NAME`, `SERVICE_WEB_SERVICE_NAMESPACE`, `SUSPENSION_TIME`, `TOTAL_PROCESSING_TIME`, `WAIT_TIME`, `WEBREQUEST_QUERY`, `WEBREQUEST_RELATIVE_URL`, `WEBREQUEST_URL`, `WEBREQUEST_URL_HOST`, `WEBREQUEST_URL_PATH`, `WEBREQUEST_URL_PORT`, `WEBSERVICE_ENDPOINT`, `WEBSERVICE_METHOD` and `ZOS_CALL_TYPE`",
		},
		"normalization": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The format of the extracted string. Possible values are `ORIGINAL`, `TO_LOWER_CASE` and `TO_UPPER_CASE`",
		},
		"delimiter_or_regex": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Depending on the `kind` value:\n\n\n* `REGEX_EXTRACTION`: The regular expression.\n\n\n* `BETWEEN_DELIMITER`: The opening delimiter string to look for.\n\n\n* All other values: The delimiter string to look for",
		},
		"use_from_child_calls": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "If `true` request attribute will be taken from a child service call. Only applicable for the `SERVICE_REQUEST_ATTRIBUTE` attribute. Defaults to `false`",
		},
		"request_attribute": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The request attribute to extract from. Required if the `kind` value is `SERVICE_REQUEST_ATTRIBUTE`. Not applicable otherwise",
		},
		"end_delimiter": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The closing delimiter string to look for. Required if the `kind` value is `BETWEEN_DELIMITER`. Not applicable otherwise",
		},
		"kind": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The type of extraction. Defines either usage of regular expression (`regex`) or the position of request attribute value to be extracted. When the `attribute` is `SERVICE_REQUEST_ATTRIBUTE` attribute and `aggregation` is `COUNT`, needs to be set to `ORIGINAL_TEXT`. Possible values are 	`AFTER_DELIMITER`, `BEFORE_DELIMITER`, `BETWEEN_DELIMITER`, `ORIGINAL_TEXT` and `REGEX_EXTRACTION`",
		},
		"source": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			MinItems:    1,
			Description: "Defines valid sources of request attributes for conditions or placeholders",
			Elem:        &schema.Resource{Schema: new(propagation.Source).Schema()},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *Placeholder) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	return properties.EncodeAll(map[string]any{
		"name":                 me.Name,
		"aggregation":          me.Aggregation,
		"attribute":            me.Attribute,
		"normalization":        me.Normalization,
		"delimiter_or_regex":   me.DelimiterOrRegex,
		"use_from_child_calls": me.UseFromChildCalls,
		"request_attribute":    me.RequestAttribute,
		"end_delimiter":        me.EndDelimiter,
		"kind":                 me.Kind,
		"source":               me.Source,
		"unknowns":             me.Unknowns,
	})
}

func (me *Placeholder) UnmarshalHCL(decoder hcl.Decoder) error {
	err := decoder.DecodeAll(map[string]any{
		"name":                 &me.Name,
		"aggregation":          &me.Aggregation,
		"attribute":            &me.Attribute,
		"normalization":        &me.Normalization,
		"delimiter_or_regex":   &me.DelimiterOrRegex,
		"use_from_child_calls": &me.UseFromChildCalls,
		"request_attribute":    &me.RequestAttribute,
		"end_delimiter":        &me.EndDelimiter,
		"kind":                 &me.Kind,
		"source":               &me.Source,
		"unknowns":             &me.Unknowns,
	})
	return err
}

func (me *Placeholder) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"name":              me.Name,
		"aggregation":       me.Aggregation,
		"attribute":         me.Attribute,
		"normalization":     me.Normalization,
		"delimiterOrRegex":  me.DelimiterOrRegex,
		"useFromChildCalls": me.UseFromChildCalls,
		"requestAttribute":  me.RequestAttribute,
		"endDelimiter":      me.EndDelimiter,
		"kind":              me.Kind,
		"source":            me.Source,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *Placeholder) UnmarshalJSON(data []byte) error {
	properties := xjson.NewProperties(me.Unknowns)
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	return properties.UnmarshalAll(map[string]any{
		"name":              &me.Name,
		"aggregation":       &me.Aggregation,
		"attribute":         &me.Attribute,
		"normalization":     &me.Normalization,
		"delimiterOrRegex":  &me.DelimiterOrRegex,
		"useFromChildCalls": &me.UseFromChildCalls,
		"requestAttribute":  &me.RequestAttribute,
		"endDelimiter":      &me.EndDelimiter,
		"kind":              &me.Kind,
		"source":            &me.Source,
	})
}

// Aggregation Which value of the request attribute must be used when it occurs across multiple child requests.
// Only applicable for the `SERVICE_REQUEST_ATTRIBUTE` attribute, when **useFromChildCalls** is `true`.
// For the `COUNT` aggregation, the **kind** field is not applicable.
type Aggregation string

// Aggregations offers the known enum values
var Aggregations = struct {
	Count Aggregation
	First Aggregation
	Last  Aggregation
}{
	"COUNT",
	"FIRST",
	"LAST",
}

// Normalization The format of the extracted string.
type Normalization string

// Normalizations offers the known enum values
var Normalizations = struct {
	Original    Normalization
	ToLowerCase Normalization
	ToUpperCase Normalization
}{
	"ORIGINAL",
	"TO_LOWER_CASE",
	"TO_UPPER_CASE",
}

// Kind The type of extraction.
//
//	Defines either usage of regular expression (`regex`) or the position of request attribute value to be extracted.
//
// When the **attribute** is `SERVICE_REQUEST_ATTRIBUTE` attribute and **aggregation** is `COUNT`, needs to be set to `ORIGINAL_TEXT`
type Kind string

// Kinds offers the known enum values
var Kinds = struct {
	AfterDelimiter   Kind
	BeforeDelimiter  Kind
	BetweenDelimiter Kind
	OriginalText     Kind
	RegexExtraction  Kind
}{
	"AFTER_DELIMITER",
	"BEFORE_DELIMITER",
	"BETWEEN_DELIMITER",
	"ORIGINAL_TEXT",
	"REGEX_EXTRACTION",
}
