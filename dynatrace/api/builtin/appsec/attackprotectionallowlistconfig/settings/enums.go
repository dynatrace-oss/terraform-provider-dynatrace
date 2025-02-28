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

package attackprotectionallowlistconfig

type AgentSideAttributeKey string

var AgentSideAttributeKeies = struct {
	ActorIp                 AgentSideAttributeKey
	DetectionType           AgentSideAttributeKey
	EntryPointPayload       AgentSideAttributeKey
	EntryPointPayloadDomain AgentSideAttributeKey
	EntryPointPayloadPort   AgentSideAttributeKey
	EntryPointUrlPath       AgentSideAttributeKey
}{
	"ACTOR_IP",
	"DETECTION_TYPE",
	"ENTRY_POINT_PAYLOAD",
	"ENTRY_POINT_PAYLOAD_DOMAIN",
	"ENTRY_POINT_PAYLOAD_PORT",
	"ENTRY_POINT_URL_PATH",
}

type AgentSideAttributeMatcher string

var AgentSideAttributeMatchers = struct {
	Contains          AgentSideAttributeMatcher
	DoesNotContain    AgentSideAttributeMatcher
	DoesNotEndWith    AgentSideAttributeMatcher
	DoesNotStartsWith AgentSideAttributeMatcher
	EndsWith          AgentSideAttributeMatcher
	Equals            AgentSideAttributeMatcher
	IpCidr            AgentSideAttributeMatcher
	NotEquals         AgentSideAttributeMatcher
	NotInIpCidr       AgentSideAttributeMatcher
	StartsWith        AgentSideAttributeMatcher
}{
	"CONTAINS",
	"DOES_NOT_CONTAIN",
	"DOES_NOT_END_WITH",
	"DOES_NOT_STARTS_WITH",
	"ENDS_WITH",
	"EQUALS",
	"IP_CIDR",
	"NOT_EQUALS",
	"NOT_IN_IP_CIDR",
	"STARTS_WITH",
}

type BlockingStrategy string

var BlockingStrategies = struct {
	Monitor BlockingStrategy
	Off     BlockingStrategy
}{
	"MONITOR",
	"OFF",
}

type DetectionType string

var DetectionTypes = struct {
	CmdInjection  DetectionType
	JndiInjection DetectionType
	SqlInjection  DetectionType
	Ssrf          DetectionType
}{
	"CMD_INJECTION",
	"JNDI_INJECTION",
	"SQL_INJECTION",
	"SSRF",
}

type ResourceAttributeValueMatcher string

var ResourceAttributeValueMatchers = struct {
	Contains         ResourceAttributeValueMatcher
	DoesNotContain   ResourceAttributeValueMatcher
	DoesNotEndWith   ResourceAttributeValueMatcher
	DoesNotExist     ResourceAttributeValueMatcher
	DoesNotStartWith ResourceAttributeValueMatcher
	EndsWith         ResourceAttributeValueMatcher
	Equals           ResourceAttributeValueMatcher
	Exists           ResourceAttributeValueMatcher
	NotEquals        ResourceAttributeValueMatcher
	StartsWith       ResourceAttributeValueMatcher
}{
	"CONTAINS",
	"DOES_NOT_CONTAIN",
	"DOES_NOT_END_WITH",
	"DOES_NOT_EXIST",
	"DOES_NOT_START_WITH",
	"ENDS_WITH",
	"EQUALS",
	"EXISTS",
	"NOT_EQUALS",
	"STARTS_WITH",
}
