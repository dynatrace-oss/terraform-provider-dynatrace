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

package hosts

import (
	"encoding/json"
	"fmt"
	"strings"

	hosts "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/anomalydetection/infrastructure/hosts/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/shutdown"

	"net/url"
)

const SchemaVersion = "1.2.8"
const SchemaID = "builtin:anomaly-detection.infrastructure-hosts"

func Service(credentials *settings.Credentials) settings.CRUDService[*hosts.Settings] {
	return &service{
		credentials: credentials,
		client:      rest.DefaultClient(credentials.URL, credentials.Token),
	}
}

type service struct {
	credentials *settings.Credentials
	client      rest.Client
}

func (me *service) Name() string {
	return me.SchemaID()
}

func (me *service) SchemaID() string {
	return SchemaID
}

func (me *service) Create(v *hosts.Settings) (*settings.Stub, error) {
	soc := settings20.SettingsObjectCreate{
		SchemaID:      SchemaID,
		SchemaVersion: SchemaVersion,
		Scope:         v.Scope,
		Value:         v,
	}

	req := me.client.Post("/api/v2/settings/objects", []settings20.SettingsObjectCreate{soc}).Expect(200)
	objectID := []settings20.SettingsObjectCreateResponse{}

	if err := req.Finish(&objectID); err != nil {
		if strings.Contains(err.Error(), "Given property 'highSystemLoadDetection' with value: 'null' does not comply with required NonNull of schema") {
			v.Host.HighSystemLoadDetection = &hosts.HighSystemLoadDetection{
				Enabled:       true,
				DetectionMode: &hosts.DetectionModes.Auto,
			}
			soc.Value = v

			req = me.client.Post("/api/v2/settings/objects", []settings20.SettingsObjectCreate{soc}).Expect(200)
			objectID = []settings20.SettingsObjectCreateResponse{}
			if err := req.Finish(&objectID); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	itemName := settings.Name(v, objectID[0].ObjectID)
	stub := &settings.Stub{ID: objectID[0].ObjectID, Name: itemName}
	return stub, nil
}

func (me *service) Update(id string, v *hosts.Settings) error {
	sou := settings20.SettingsObjectUpdate{Value: v, SchemaVersion: SchemaVersion}
	if err := me.client.Put(fmt.Sprintf("/api/v2/settings/objects/%s", url.PathEscape(id)), &sou, 200).Finish(); err != nil {
		if strings.Contains(err.Error(), "Given property 'highSystemLoadDetection' with value: 'null' does not comply with required NonNull of schema") {
			v.Host.HighSystemLoadDetection = &hosts.HighSystemLoadDetection{
				Enabled:       true,
				DetectionMode: &hosts.DetectionModes.Auto,
			}
			sou.Value = v

			if err = me.client.Put(fmt.Sprintf("/api/v2/settings/objects/%s", url.PathEscape(id)), &sou, 200).Finish(); err == nil {
				return nil
			}
		}
		return err
	}
	return nil
}

func (me *service) Validate(v *hosts.Settings) error {
	return nil // Settings 2.0 doesn't offer validation
}

func (me *service) Delete(id string) error {
	err := me.client.Delete(fmt.Sprintf("/api/v2/settings/objects/%s", url.PathEscape(id)), 204).Finish()
	return err
}

func (me *service) Get(id string, v *hosts.Settings) error {
	var err error
	var settingsObject settings20.SettingsObject

	req := me.client.Get(fmt.Sprintf("/api/v2/settings/objects/%s", url.PathEscape(id))).Expect(200)
	if err = req.Finish(&settingsObject); err != nil {
		return err
	}

	if err = json.Unmarshal(settingsObject.Value, v); err != nil {
		return err
	}
	settings.SetScope(v, settingsObject.Scope)

	return nil
}

func (me *service) List() (settings.Stubs, error) {
	var err error

	stubs := settings.Stubs{}
	nextPage := true

	var nextPageKey *string
	for nextPage {
		var sol settings20.SettingsObjectList
		var urlStr string
		if nextPageKey != nil {
			urlStr = fmt.Sprintf("/api/v2/settings/objects?nextPageKey=%s", url.QueryEscape(*nextPageKey))
		} else {
			urlStr = fmt.Sprintf("/api/v2/settings/objects?schemaIds=%s&fields=%s&pageSize=100", url.QueryEscape(me.SchemaID()), url.QueryEscape("objectId,value,scope,schemaVersion"))
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
				newItem := hosts.Settings{}
				if err = json.Unmarshal(item.Value, &newItem); err != nil {
					return nil, err
				}
				settings.SetScope(&newItem, item.Scope)
				itemName := newItem.Name()
				stub := &settings.Stub{ID: item.ObjectID, Name: itemName, Value: newItem}
				if len(itemName) > 0 {
					stubs = append(stubs, stub)
				}
			}
		}
		nextPageKey = sol.NextPageKey
		nextPage = (nextPageKey != nil)
	}

	return stubs, nil
}
