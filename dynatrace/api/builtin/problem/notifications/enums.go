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

package notifications

type Type string

var Types = struct {
	Email        Type
	Slack        Type
	Jira         Type
	AnsibleTower Type
	OpsGenie     Type
	PagerDuty    Type
	VictorOps    Type
	WebHook      Type
	XMatters     Type
	Trello       Type
	ServiceNow   Type
}{
	Type("EMAIL"),
	Type("SLACK"),
	Type("JIRA"),
	Type("ANSIBLETOWER"),
	Type("OPS_GENIE"),
	Type("PAGER_DUTY"),
	Type("VICTOROPS"),
	Type("WEBHOOK"),
	Type("XMATTERS"),
	Type("TRELLO"),
	Type("SERVICE_NOW"),
}
