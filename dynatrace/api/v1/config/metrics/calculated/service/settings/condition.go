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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/metrics/calculated/service/settings/comparisoninfo"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Condition A condition of a rule usage.
type Condition struct {
	Attribute      Attribute                     `json:"attribute"`      // The attribute to be matched.  Note that for a service property attribute you must use the comparison of the `FAST_STRING` type.
	ComparisonInfo comparisoninfo.ComparisonInfo `json:"comparisonInfo"` // Type-specific comparison for attributes. The actual set of fields depends on the `type` of the comparison.  See the [Service metrics API - JSON models](https://dt-url.net/9803svb) help topic for example models of every notification type.
	Unknowns       map[string]json.RawMessage    `json:"-"`
}

func (me *Condition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"attribute": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The attribute to be matched.  Note that for a service property attribute you must use the comparison of the `FAST_STRING` type. Possible values are `ACTOR_SYSTEM`, `AKKA_ACTOR_CLASS_NAME`, `AKKA_ACTOR_MESSAGE_TYPE`, `AKKA_ACTOR_PATH`, `APPLICATION_BUILD_VERSION`, `APPLICATION_RELEASE_VERSION`, `AZURE_FUNCTIONS_FUNCTION_NAME`, `AZURE_FUNCTIONS_SITE_NAME`, `CICS_PROGRAM_NAME`, `CICS_SYSTEM_ID`, `CICS_TASK_ID`, `CICS_TRANSACTION_ID`, `CICS_USER_ID`, `CPU_TIME`, `CTG_GATEWAY_URL`, `CTG_PROGRAM`, `CTG_SERVER_NAME`, `CTG_TRANSACTION_ID`, `CUSTOMSERVICE_CLASS`, `CUSTOMSERVICE_METHOD`, `DATABASE_CHILD_CALL_COUNT`, `DATABASE_CHILD_CALL_TIME`, `DATABASE_HOST`, `DATABASE_NAME`, `DATABASE_TYPE`, `DATABASE_URL`, `DISK_IO_TIME`, `ERROR_COUNT`, `ESB_APPLICATION_NAME`, `ESB_INPUT_TYPE`, `ESB_LIBRARY_NAME`, `ESB_MESSAGE_FLOW_NAME`, `EXCEPTION_CLASS`, `EXCEPTION_MESSAGE`, `FAILED_STATE`, `FAILURE_REASON`, `FLAW_STATE`, `HTTP_REQUEST_METHOD`, `HTTP_STATUS`, `HTTP_STATUS_CLASS`, `IMS_PROGRAM_NAME`, `IMS_TRANSACTION_ID`, `IMS_USER_ID`, `IO_TIME`, `IS_KEY_REQUEST`, `LAMBDA_COLDSTART`, `LOCK_TIME`, `MESSAGING_DESTINATION_TYPE`, `MESSAGING_IS_TEMPORARY_QUEUE`, `MESSAGING_QUEUE_NAME`, `MESSAGING_QUEUE_VENDOR`, `NETWORK_IO_TIME`, `NON_DATABASE_CHILD_CALL_COUNT`, `NON_DATABASE_CHILD_CALL_TIME`, `PROCESS_GROUP_NAME`, `PROCESS_GROUP_TAG`, `REMOTE_ENDPOINT`, `REMOTE_METHOD`, `REMOTE_SERVICE_NAME`, `REQUEST_NAME`, `REQUEST_TYPE`, `RESPONSE_TIME`, `RESPONSE_TIME_CLIENT`, `RMI_CLASS`, `RMI_METHOD`, `SERVICE_DISPLAY_NAME`, `SERVICE_NAME`, `SERVICE_PORT`, `SERVICE_PUBLIC_DOMAIN_NAME`, `SERVICE_REQUEST_ATTRIBUTE`, `SERVICE_TAG`, `SERVICE_TYPE`, `SERVICE_WEB_APPLICATION_ID`, `SERVICE_WEB_CONTEXT_ROOT`, `SERVICE_WEB_SERVER_NAME`, `SERVICE_WEB_SERVICE_NAME`, `SERVICE_WEB_SERVICE_NAMESPACE`, `SUSPENSION_TIME`, `TOTAL_PROCESSING_TIME`, `WAIT_TIME`, `WEBREQUEST_QUERY`, `WEBREQUEST_RELATIVE_URL`, `WEBREQUEST_URL`, `WEBREQUEST_URL_HOST`, `WEBREQUEST_URL_PATH`, `WEBREQUEST_URL_PORT`, `WEBSERVICE_ENDPOINT`, `WEBSERVICE_METHOD` and `ZOS_CALL_TYPE`",
		},
		"comparison": {
			Type:        schema.TypeList,
			Required:    true,
			MaxItems:    1,
			MinItems:    1,
			Description: "Type-specific comparison for attributes",
			Elem:        &schema.Resource{Schema: new(comparisoninfo.Wrapper).Schema()},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		}}
}

func (me *Condition) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	return properties.EncodeAll(map[string]any{
		"attribute":  me.Attribute,
		"comparison": &comparisoninfo.Wrapper{Comparison: me.ComparisonInfo},
		"unknowns":   me.Unknowns,
	})
}

func (me *Condition) UnmarshalHCL(decoder hcl.Decoder) error {
	compWrap := comparisoninfo.Wrapper{Comparison: me.ComparisonInfo}
	if err := decoder.DecodeAll(map[string]any{
		"attribute":  &me.Attribute,
		"comparison": &compWrap,
		"unknowns":   &me.Unknowns,
	}); err != nil {
		return err
	}
	me.ComparisonInfo = compWrap.Comparison
	return nil
}

func (me *Condition) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"attribute":      me.Attribute,
		"comparisonInfo": me.ComparisonInfo,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *Condition) UnmarshalJSON(data []byte) error {
	properties := xjson.NewProperties(me.Unknowns)
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	compWrap := comparisoninfo.Wrapper{Comparison: me.ComparisonInfo}
	if err := properties.UnmarshalAll(map[string]any{
		"attribute":      &me.Attribute,
		"comparisonInfo": &compWrap,
	}); err != nil {
		return err
	}
	me.ComparisonInfo = compWrap.Comparison
	return nil
}

// Attribute The attribute to extract from. You can only use attributes of the **string** type.
type Attribute string

// Attributes offers the known enum values
var Attributes = struct {
	ActorSystem                Attribute
	AkkaActorClassName         Attribute
	AkkaActorMessageType       Attribute
	AkkaActorPath              Attribute
	ApplicationBuildVersion    Attribute
	ApplicationReleaseVersion  Attribute
	AzureFunctionsFunctionName Attribute
	AzureFunctionsSiteName     Attribute
	CICSProgramName            Attribute
	CICSSystemID               Attribute
	CICSTaskID                 Attribute
	CICSTransactionID          Attribute
	CICSUserID                 Attribute
	CPUTime                    Attribute
	CTGGatewayURL              Attribute
	CTGProgram                 Attribute
	CTGServerName              Attribute
	CTGTransactionID           Attribute
	CustomserviceClass         Attribute
	CustomserviceMethod        Attribute
	DatabaseChildCallCount     Attribute
	DatabaseChildCallTime      Attribute
	DatabaseHost               Attribute
	DatabaseName               Attribute
	DatabaseType               Attribute
	DatabaseURL                Attribute
	DiskIoTime                 Attribute
	ErrorCount                 Attribute
	EsbApplicationName         Attribute
	EsbInputType               Attribute
	EsbLibraryName             Attribute
	EsbMessageFlowName         Attribute
	ExceptionClass             Attribute
	ExceptionMessage           Attribute
	FailedState                Attribute
	FailureReason              Attribute
	FlawState                  Attribute
	HTTPRequestMethod          Attribute
	HTTPStatus                 Attribute
	HTTPStatusClass            Attribute
	IMSProgramName             Attribute
	IMSTransactionID           Attribute
	IMSUserID                  Attribute
	IoTime                     Attribute
	IsKeyRequest               Attribute
	LambdaColdstart            Attribute
	LockTime                   Attribute
	MessagingDestinationType   Attribute
	MessagingIsTemporaryQueue  Attribute
	MessagingQueueName         Attribute
	MessagingQueueVendor       Attribute
	NetworkIoTime              Attribute
	NonDatabaseChildCallCount  Attribute
	NonDatabaseChildCallTime   Attribute
	ProcessGroupName           Attribute
	ProcessGroupTag            Attribute
	RemoteEndpoint             Attribute
	RemoteMethod               Attribute
	RemoteServiceName          Attribute
	RequestName                Attribute
	RequestType                Attribute
	ResponseTime               Attribute
	ResponseTimeClient         Attribute
	RMIClass                   Attribute
	RMIMethod                  Attribute
	ServiceDisplayName         Attribute
	ServiceName                Attribute
	ServicePort                Attribute
	ServicePublicDomainName    Attribute
	ServiceRequestAttribute    Attribute
	ServiceTag                 Attribute
	ServiceType                Attribute
	ServiceWebApplicationID    Attribute
	ServiceWebContextRoot      Attribute
	ServiceWebServerName       Attribute
	ServiceWebServiceName      Attribute
	ServiceWebServiceNamespace Attribute
	SuspensionTime             Attribute
	TotalProcessingTime        Attribute
	WaitTime                   Attribute
	WebrequestQuery            Attribute
	WebrequestRelativeURL      Attribute
	WebrequestURL              Attribute
	WebrequestURLHost          Attribute
	WebrequestURLPath          Attribute
	WebrequestURLPort          Attribute
	WebserviceEndpoint         Attribute
	WebserviceMethod           Attribute
	ZosCallType                Attribute
}{
	"ACTOR_SYSTEM",
	"AKKA_ACTOR_CLASS_NAME",
	"AKKA_ACTOR_MESSAGE_TYPE",
	"AKKA_ACTOR_PATH",
	"APPLICATION_BUILD_VERSION",
	"APPLICATION_RELEASE_VERSION",
	"AZURE_FUNCTIONS_FUNCTION_NAME",
	"AZURE_FUNCTIONS_SITE_NAME",
	"CICS_PROGRAM_NAME",
	"CICS_SYSTEM_ID",
	"CICS_TASK_ID",
	"CICS_TRANSACTION_ID",
	"CICS_USER_ID",
	"CPU_TIME",
	"CTG_GATEWAY_URL",
	"CTG_PROGRAM",
	"CTG_SERVER_NAME",
	"CTG_TRANSACTION_ID",
	"CUSTOMSERVICE_CLASS",
	"CUSTOMSERVICE_METHOD",
	"DATABASE_CHILD_CALL_COUNT",
	"DATABASE_CHILD_CALL_TIME",
	"DATABASE_HOST",
	"DATABASE_NAME",
	"DATABASE_TYPE",
	"DATABASE_URL",
	"DISK_IO_TIME",
	"ERROR_COUNT",
	"ESB_APPLICATION_NAME",
	"ESB_INPUT_TYPE",
	"ESB_LIBRARY_NAME",
	"ESB_MESSAGE_FLOW_NAME",
	"EXCEPTION_CLASS",
	"EXCEPTION_MESSAGE",
	"FAILED_STATE",
	"FAILURE_REASON",
	"FLAW_STATE",
	"HTTP_REQUEST_METHOD",
	"HTTP_STATUS",
	"HTTP_STATUS_CLASS",
	"IMS_PROGRAM_NAME",
	"IMS_TRANSACTION_ID",
	"IMS_USER_ID",
	"IO_TIME",
	"IS_KEY_REQUEST",
	"LAMBDA_COLDSTART",
	"LOCK_TIME",
	"MESSAGING_DESTINATION_TYPE",
	"MESSAGING_IS_TEMPORARY_QUEUE",
	"MESSAGING_QUEUE_NAME",
	"MESSAGING_QUEUE_VENDOR",
	"NETWORK_IO_TIME",
	"NON_DATABASE_CHILD_CALL_COUNT",
	"NON_DATABASE_CHILD_CALL_TIME",
	"PROCESS_GROUP_NAME",
	"PROCESS_GROUP_TAG",
	"REMOTE_ENDPOINT",
	"REMOTE_METHOD",
	"REMOTE_SERVICE_NAME",
	"REQUEST_NAME",
	"REQUEST_TYPE",
	"RESPONSE_TIME",
	"RESPONSE_TIME_CLIENT",
	"RMI_CLASS",
	"RMI_METHOD",
	"SERVICE_DISPLAY_NAME",
	"SERVICE_NAME",
	"SERVICE_PORT",
	"SERVICE_PUBLIC_DOMAIN_NAME",
	"SERVICE_REQUEST_ATTRIBUTE",
	"SERVICE_TAG",
	"SERVICE_TYPE",
	"SERVICE_WEB_APPLICATION_ID",
	"SERVICE_WEB_CONTEXT_ROOT",
	"SERVICE_WEB_SERVER_NAME",
	"SERVICE_WEB_SERVICE_NAME",
	"SERVICE_WEB_SERVICE_NAMESPACE",
	"SUSPENSION_TIME",
	"TOTAL_PROCESSING_TIME",
	"WAIT_TIME",
	"WEBREQUEST_QUERY",
	"WEBREQUEST_RELATIVE_URL",
	"WEBREQUEST_URL",
	"WEBREQUEST_URL_HOST",
	"WEBREQUEST_URL_PATH",
	"WEBREQUEST_URL_PORT",
	"WEBSERVICE_ENDPOINT",
	"WEBSERVICE_METHOD",
	"ZOS_CALL_TYPE",
}
