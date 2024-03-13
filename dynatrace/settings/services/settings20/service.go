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

package settings20

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/httpcache"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/shutdown"

	"net/url"
)

// var NO_REPAIR_INPUT = os.Getenv("DT_NO_REPAIR_INPUT") == "true"
var NO_REPAIR_INPUT = true

func Service[T settings.Settings](credentials *settings.Credentials, schemaID string, schemaVersion string, options ...*ServiceOptions[T]) settings.CRUDService[T] {
	var opts *ServiceOptions[T]
	if len(options) > 0 {
		opts = options[0]
	}
	return &service[T]{
		schemaID: schemaID,
		// schemaVersion: schemaVersion,
		client:  httpcache.DefaultClient(credentials.URL, credentials.Token, schemaID),
		options: opts,
	}
}

type SettingsObjectUpdate struct {
	SchemaVersion string `json:"schemaVersion,omitempty"`
	Value         any    `json:"value"`
	InsertAfter   string `json:"insertAfter,omitempty"`
}

type SettingsObjectCreate struct {
	SchemaVersion string `json:"schemaVersion,omitempty"`
	SchemaID      string `json:"schemaId"`
	Scope         string `json:"scope,omitempty"`
	Value         any    `json:"value"`
	InsertAfter   string `json:"insertAfter,omitempty"`
}

type SettingsObjectCreateResponse struct {
	ObjectID string `json:"objectId"`
}

type service[T settings.Settings] struct {
	schemaID      string
	schemaVersion string
	client        rest.Client
	options       *ServiceOptions[T]
}

func (me *service[T]) LegacyID() func(id string) string {
	if me.options != nil && me.options.LegacyID != nil {
		return me.options.LegacyID
	}
	return nil
}

func (me *service[T]) Get(id string, v T) error {
	var err error
	var settingsObject SettingsObject

	req := me.client.Get(fmt.Sprintf("/api/v2/settings/objects/%s", url.PathEscape(id))).Expect(200)
	if err = req.Finish(&settingsObject); err != nil {
		return err
	}

	if err = json.Unmarshal(settingsObject.Value, v); err != nil {
		return err
	}
	settings.SetScope(v, settingsObject.Scope)
	if me.options != nil && me.options.LegacyID != nil {
		settings.SetLegacyID(id, me.options.LegacyID, v)
	}
	insertBefore, insertAfter, err := me.getInsertIDs(id)
	if err != nil {
		return err
	}
	if insertBefore != nil {
		settings.SetInsertBefore(v, *insertBefore)
	}
	if insertAfter != nil {
		settings.SetInsertAfter(v, *insertAfter)
	}

	return nil
}

func (me *service[T]) getInsertIDs(id string, optIds ...[]string) (*string, *string, error) {
	if len(id) == 0 {
		return nil, nil, nil
	}
	var ids []string
	if len(optIds) > 0 {
		ids = optIds[0]
	} else {
		listedIDs, err := me.listIDs()
		if err != nil {
			return nil, nil, nil
		}
		ids = listedIDs
	}
	insertAfter, err := me.getInsertAfter(ids, id)
	if err != nil {
		return nil, nil, nil
	}
	insertBefore, err := me.getInsertBefore(ids, id)
	if err != nil {
		return nil, nil, nil
	}
	return insertBefore, insertAfter, nil
}

func (me *service[T]) getInsertBefore(ids []string, id string) (*string, error) {
	if len(id) == 0 {
		return nil, nil
	}
	prevIDWasMatching := false
	for _, curID := range ids {
		if prevIDWasMatching {
			insertBeforeID := curID
			if len(insertBeforeID) == 0 {
				return nil, nil
			}
			return &insertBeforeID, nil
		}
		prevIDWasMatching = (curID == id)
	}
	return nil, nil
}

func (me *service[T]) getInsertAfter(ids []string, id string) (*string, error) {
	if len(id) == 0 {
		return nil, nil
	}
	prevID := ""
	for _, curID := range ids {
		if curID == id {
			if len(prevID) == 0 {
				return nil, nil
			}
			return &prevID, nil
		}
		prevID = curID
	}
	return nil, nil
}

func (me *service[T]) listIDs() ([]string, error) {
	var err error

	ids := []string{}
	nextPage := true

	var nextPageKey *string
	for nextPage {
		var sol SettingsObjectList
		var urlStr string
		if nextPageKey != nil {
			urlStr = fmt.Sprintf("/api/v2/settings/objects?nextPageKey=%s", url.QueryEscape(*nextPageKey))
		} else {
			urlStr = fmt.Sprintf("/api/v2/settings/objects?schemaIds=%s&fields=%s&pageSize=100", url.QueryEscape(me.SchemaID()), url.QueryEscape("objectId,scope,schemaVersion"))
		}
		req := me.client.Get(urlStr, 200)
		if err = req.Finish(&sol); err != nil {
			return nil, err
		}
		if shutdown.System.Stopped() {
			return ids, nil
		}

		if len(sol.Items) > 0 {
			for _, item := range sol.Items {
				ids = append(ids, item.ObjectID)
			}
		}
		nextPageKey = sol.NextPageKey
		nextPage = (nextPageKey != nil)
	}

	return ids, nil
}

func (me *service[T]) List() (api.Stubs, error) {
	var err error

	ids, err := me.listIDs()
	if err != nil {
		return api.Stubs{}, err
	}

	stubs := api.Stubs{}
	nextPage := true

	var nextPageKey *string
	for nextPage {
		var sol SettingsObjectList
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
				newItem := settings.NewSettings[T](me)
				if err = json.Unmarshal(item.Value, &newItem); err != nil {
					return nil, err
				}
				if me.options != nil && me.options.LegacyID != nil {
					settings.SetLegacyID(item.ObjectID, me.options.LegacyID, newItem)
				}
				settings.SetScope(newItem, item.Scope)
				insertBefore, insertAfter, err := me.getInsertIDs(item.ObjectID, ids)
				if err != nil {
					return api.Stubs{}, err
				}
				if insertBefore != nil {
					settings.SetInsertBefore(newItem, *insertBefore)
				}
				if insertAfter != nil {
					settings.SetInsertAfter(newItem, *insertAfter)
				}
				var itemName string
				if me.options != nil && me.options.Name != nil {
					if itemName, err = me.options.Name(item.ObjectID, newItem); err != nil {
						itemName = settings.Name(newItem, item.ObjectID)
					}
				} else {
					itemName = settings.Name(newItem, item.ObjectID)
				}
				stub := &api.Stub{ID: item.ObjectID, Name: itemName, Value: newItem, LegacyID: settings.GetLegacyID(newItem)}
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

func (me *service[T]) Validate(v T) error {
	return nil // Settings 2.0 doesn't offer validation
}

func (me *service[T]) Create(v T) (*api.Stub, error) {
	return me.create(v, false, false)
}

type Matcher interface {
	Match(o any) bool
}

func (me *service[T]) create(v T, retry bool, noInsertAfter bool) (*api.Stub, error) {

	if me.options != nil && me.options.Duplicates != nil {
		dupStub, dupErr := me.options.Duplicates(me, v)
		if dupErr != nil {
			return nil, dupErr
		}
		if dupStub != nil {
			return dupStub, nil
		}
	}

	// Special handling for settings with a method named Match(v any) bool
	// It signals that instead of creating a new record the existing ones on the remote
	// should be investigated - and if a match exists the original state should get persisted
	// Upon delete that original state will get reconstructed
	// Note: Such settings also need to contain a field named `RestoreOnDelete` of type `*string`
	if matcher, ok := any(v).(Matcher); ok {
		var stubs api.Stubs
		var err error
		if stubs, err = me.List(); err != nil {
			return nil, err
		}
		for _, stub := range stubs {
			if stub == nil {
				continue
			}
			if stub.Value == nil {
				continue
			}
			if matcher.Match(stub.Value) {
				data, je := json.Marshal(stub.Value)
				if je != nil {
					break
				}
				asjson := string(data)
				settings.SetRestoreOnDelete(asjson, v)
				stub.Value = v
				return stub, me.Update(stub.ID, v)
			}
		}
	}
	soc := SettingsObjectCreate{
		SchemaID:      me.schemaID,
		SchemaVersion: me.schemaVersion,
		Scope:         "environment",
		Value:         v,
	}
	soc.Scope = settings.GetScope(v)
	if !noInsertAfter {
		if insertAfter := settings.GetInsertAfter(v); insertAfter != nil {
			soc.InsertAfter = *insertAfter
		}
	}

	var req rest.Request
	if NO_REPAIR_INPUT {
		req = me.client.Post("/api/v2/settings/objects", []SettingsObjectCreate{soc}).Expect(200)
	} else {
		req = me.client.Post("/api/v2/settings/objects?repairInput=true", []SettingsObjectCreate{soc}).Expect(200)
	}

	objectID := []SettingsObjectCreateResponse{}

	if oerr := req.Finish(&objectID); oerr != nil {
		if isInvalidInsertAfter(oerr) {
			return me.create(v, retry, true)
		}

		if me.options != nil && me.options.CreateRetry != nil && !retry {
			if modifiedPayload := me.options.CreateRetry(v, oerr); !reflect.ValueOf(modifiedPayload).IsNil() {
				return me.create(modifiedPayload, true, noInsertAfter)
			}
		}
		if me.options != nil && me.options.HijackOnCreate != nil {
			var hijackedStub *api.Stub
			var hierr error
			if hijackedStub, hierr = me.options.HijackOnCreate(oerr, me, v); hierr != nil {
				return nil, hierr
			}
			if hijackedStub != nil {
				return hijackedStub, me.Update(hijackedStub.ID, v)
			} else {
				return nil, oerr
			}
		}
		return nil, oerr
	}
	itemName := settings.Name(v, objectID[0].ObjectID)
	stub := &api.Stub{ID: objectID[0].ObjectID, Name: itemName}
	return stub, nil
}

func isInvalidInsertAfter(err error) bool {
	if err == nil {
		return false
	}
	switch resterr := err.(type) {
	case *rest.Error:
		return isInvalidInsertAfterRestErr(resterr)
	case rest.Error:
		return isInvalidInsertAfterRestErr(&resterr)
	default:
		return false
	}
}

func isInvalidInsertAfterRestErr(resterr *rest.Error) bool {
	if resterr == nil {
		return false
	}
	if resterr.Code == 404 && resterr.Message == "Settings not found" {
		return true
	}
	if len(resterr.ConstraintViolations) == 0 {
		return false
	}
	for _, constraingViolation := range resterr.ConstraintViolations {
		if constraingViolation.Message == "Setting value cannot be inserted to the specified position" {
			return true
		}
	}
	return false
}

func (me *service[T]) Update(id string, v T) error {
	return me.update(id, v, false, false)
}

func (me *service[T]) update(id string, v T, retry bool, noInsertAfter bool) error {
	sou := SettingsObjectUpdate{Value: v, SchemaVersion: me.schemaVersion}

	if !noInsertAfter {
		if insertAfter := settings.GetInsertAfter(v); insertAfter != nil {
			sou.InsertAfter = *insertAfter
		}
	}
	var req rest.Request
	if NO_REPAIR_INPUT {
		req = me.client.Put(fmt.Sprintf("/api/v2/settings/objects/%s", url.PathEscape(id)), &sou, 200)
	} else {
		req = me.client.Put(fmt.Sprintf("/api/v2/settings/objects/%s?repairInput=true", url.PathEscape(id)), &sou, 200)
	}

	if err := req.Finish(); err != nil {
		if isInvalidInsertAfter(err) {
			return me.update(id, v, retry, true)
		}
		if me.options != nil && me.options.UpdateRetry != nil && !retry {
			if modifiedPayload := me.options.UpdateRetry(v, err); (any)(modifiedPayload) != (any)(nil) {
				return me.update(id, modifiedPayload, true, noInsertAfter)
			}
		}
		return err
	}
	return nil
}

func (me *service[T]) Delete(id string) error {
	return me.delete(id, 0)
}

func (me *service[T]) delete(id string, numRetries int) error {
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

func (me *service[T]) Name() string {
	return me.SchemaID()
}

func (me *service[T]) SchemaID() string {
	return me.schemaID
}
