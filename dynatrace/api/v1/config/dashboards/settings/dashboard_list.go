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

package dashboards

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

// DashboardList is a list of short representations of dashboards
type DashboardList struct {
	Dashboards []*DashboardStub `json:"dashboards,omitempty"` // the short representations of the dashboards
}

func (me *DashboardList) ToStubs() settings.Stubs {
	stubs := settings.Stubs{}
	for _, dbstub := range me.Dashboards {
		if dbstub.Name == nil {
			panic(dbstub.ID)
		}
		if dbstub.Owner != nil {
			if *dbstub.Owner != "Dynatrace" {
				stubs = append(stubs, &settings.Stub{ID: dbstub.ID, Name: fmt.Sprintf("%s owned by %s", *dbstub.Name, *dbstub.Owner)})
			}
		} else {
			stubs = append(stubs, &settings.Stub{ID: dbstub.ID, Name: *dbstub.Name})
		}
	}
	return stubs
}

// DashboardStub is a short representation of a dashboard
type DashboardStub struct {
	ID    string  `json:"id"`              // the ID of the dashboard
	Name  *string `json:"name,omitempty"`  // the name of the dashboard
	Owner *string `json:"owner,omitempty"` // the owner of the dashboard
}
