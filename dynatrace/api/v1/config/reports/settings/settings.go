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

package reports

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ReportStubList struct {
	Values []ReportStub `json:"values"`
}

type ReportStub struct {
	DashboardID string `json:"dashboardId"`
	ID          string `json:"id"`
}

type ReportStubs []*ReportStub

func (me *ReportStubList) ToStubs() api.Stubs {
	res := []*api.Stub{}
	for _, setting := range me.Values {
		logging.File.Println("DashboardId: " + setting.DashboardID)
		res = append(res, &api.Stub{ID: setting.ID, Name: setting.DashboardID})
	}
	return res
}

type Settings struct {
	DashboardId        string         `json:"dashboardId"`       // The ID of the associated dashboard
	EmailNotifications *bool          `json:"enabled,omitempty"` // The email notifications for the dashboard report are enabled (true) or disabled (false).
	Subscriptions      *Subscriptions `json:"subscriptions"`     // A list of the report subscribers
	Type               Type           `json:"type"`              // The type of report
}

func (me *Settings) Name() string {
	return me.DashboardId
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"dashboard_id": {
			Type:        schema.TypeString,
			Description: "The ID of the associated dashboard",
			Required:    true,
		},
		"email_notifications": {
			Type:        schema.TypeBool,
			Description: "The email notifications for the dashboard report are enabled (true) or disabled (false).",
			Optional:    true,
		},
		"subscriptions": {
			Type:        schema.TypeList,
			Description: "A list of the report subscribers",
			Optional:    true,
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Subscriptions).Schema()},
		},
		"type": {
			Type:        schema.TypeString,
			Description: "The type of report",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"dashboard_id":        me.DashboardId,
		"email_notifications": me.EmailNotifications,
		"subscriptions":       me.Subscriptions,
		"type":                me.Type,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"dashboard_id":        &me.DashboardId,
		"email_notifications": &me.EmailNotifications,
		"subscriptions":       &me.Subscriptions,
		"type":                &me.Type,
	})
}
