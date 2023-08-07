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

package generic

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	generic "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/generic/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/shutdown"

	"net/url"
)

// var NO_REPAIR_INPUT = os.Getenv("DT_NO_REPAIR_INPUT") == "true"
var NO_REPAIR_INPUT = true

func Service(credentials *settings.Credentials) settings.CRUDService[*generic.Settings] {
	return &service{client: rest.DefaultClient(credentials.URL, credentials.Token)}
}

type SettingsObjectUpdate struct {
	SchemaVersion string          `json:"schemaVersion,omitempty"`
	Value         json.RawMessage `json:"value"`
}

type SettingsObjectCreate struct {
	SchemaVersion string          `json:"schemaVersion,omitempty"`
	SchemaID      string          `json:"schemaId"`
	Scope         string          `json:"scope"`
	Value         json.RawMessage `json:"value"`
}

type SettingsObjectCreateResponse struct {
	ObjectID string `json:"objectId"`
}

type service struct {
	client rest.Client
}

func (me *service) Get(id string, v *generic.Settings) error {
	var err error
	var settingsObject SettingsObject

	req := me.client.Get(fmt.Sprintf("/api/v2/settings/objects/%s", url.PathEscape(id))).Expect(200)
	if err = req.Finish(&settingsObject); err != nil {
		return err
	}
	v.Value = string(settingsObject.Value)
	v.Scope = settingsObject.Scope
	v.SchemaID = settingsObject.SchemaID

	return nil
}

func (me *service) List() (api.Stubs, error) {
	schemaIDs := strings.TrimSpace(os.Getenv("DYNATRACE_SCHEMA_IDS"))
	if len(schemaIDs) == 0 {
		return api.Stubs{}, nil
	}
	var stubs api.Stubs
	for _, schemaID := range strings.Split(schemaIDs, ",") {
		var err error
		nextPage := true

		var nextPageKey *string
		for nextPage {
			var sol SettingsObjectList
			var urlStr string
			if nextPageKey != nil {
				urlStr = fmt.Sprintf("/api/v2/settings/objects?nextPageKey=%s", url.QueryEscape(*nextPageKey))
			} else {
				urlStr = fmt.Sprintf("/api/v2/settings/objects?schemaIds=%s&fields=%s&pageSize=100", url.QueryEscape(strings.TrimSpace(schemaID)), url.QueryEscape("objectId,value,scope,schemaVersion"))
			}
			req := me.client.Get(urlStr, 200)
			if err = req.Finish(&sol); err != nil {
				return nil, err
			}
			if shutdown.System.Stopped() {
				return stubs, nil
			}

			if len(sol.Items) > 0 {
				for _, item := range sol.Items {
					newItem := new(generic.Settings)
					newItem.Value = string(item.Value)
					newItem.Scope = item.Scope
					newItem.SchemaID = schemaID
					newItem.Scope = item.Scope
					stubs = append(stubs, &api.Stub{ID: item.ObjectID, Name: item.ObjectID, Value: newItem})
				}
			}
			nextPageKey = sol.NextPageKey
			nextPage = (nextPageKey != nil)
		}
	}

	return stubs, nil
}

func (me *service) Validate(v *generic.Settings) error {
	return nil // Settings 2.0 doesn't offer validation
}

func (me *service) Create(v *generic.Settings) (*api.Stub, error) {
	return me.create(v, false)
}

type Matcher interface {
	Match(o any) bool
}

func (me *service) create(v *generic.Settings, retry bool) (*api.Stub, error) {
	vdata, verr := json.Marshal(v.Value)
	if verr != nil {
		return nil, verr
	}
	soc := SettingsObjectCreate{
		SchemaID: v.SchemaID,
		Scope:    "environment",
		Value:    vdata,
	}
	soc.Scope = v.Scope

	var req rest.Request
	if NO_REPAIR_INPUT {
		req = me.client.Post("/api/v2/settings/objects", []SettingsObjectCreate{soc}).Expect(200)
	} else {
		req = me.client.Post("/api/v2/settings/objects?repairInput=true", []SettingsObjectCreate{soc}).Expect(200)
	}

	objectID := []SettingsObjectCreateResponse{}

	if oerr := req.Finish(&objectID); oerr != nil {
		return nil, oerr
	}
	itemName := objectID[0].ObjectID
	stub := &api.Stub{ID: objectID[0].ObjectID, Name: itemName}
	return stub, nil
}

func (me *service) Update(id string, v *generic.Settings) error {
	vdata, verr := json.Marshal(v)
	if verr != nil {
		return verr
	}

	sou := SettingsObjectUpdate{Value: vdata}
	var req rest.Request
	if NO_REPAIR_INPUT {
		req = me.client.Put(fmt.Sprintf("/api/v2/settings/objects/%s", url.PathEscape(id)), &sou, 200)
	} else {
		req = me.client.Put(fmt.Sprintf("/api/v2/settings/objects/%s?repairInput=true", url.PathEscape(id)), &sou, 200)
	}

	if err := req.Finish(); err != nil {
		return err
	}
	return nil
}

func (me *service) Delete(id string) error {
	return me.delete(id, 0)
}

func (me *service) delete(id string, numRetries int) error {
	err := me.client.Delete(fmt.Sprintf("/api/v2/settings/objects/%s", url.PathEscape(id)), 204).Finish()
	if err != nil && strings.Contains(err.Error(), "Deletion of value(s) is not allowed") {
		return nil
	}
	if err != nil && strings.Contains(err.Error(), "Internal Server Error occurred") {
		if numRetries == 10 {
			return err
		}
		time.Sleep(6 * time.Second)
		return me.delete(id, numRetries+1)
	}
	return err

}

func (me *service) Name() string {
	return me.SchemaID()
}

func (me *service) SchemaID() string {
	return "generic"
}
