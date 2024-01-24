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

package attribute

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/httpcache"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/shutdown"

	attribute "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/span/attribute/settings"

	allowlist "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/attribute/allowlist"
	allowlistsettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/attribute/allowlist/settings"
	masking "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/attribute/masking"
	maskingsettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/attribute/masking/settings"
)

const SchemaID = "builtin:span-attribute"
const SchemaVersion = "0.0.42"

func Service(credentials *settings.Credentials) settings.CRUDService[*attribute.Settings] {
	return &service{
		credentials: credentials,
		client:      httpcache.DefaultClient(credentials.URL, credentials.Token, SchemaID),
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

func (me *service) Create(v *attribute.Settings) (*api.Stub, error) {
	soc := settings20.SettingsObjectCreate{
		SchemaID:      SchemaID,
		SchemaVersion: SchemaVersion,
		Scope:         "environment",
		Value:         v,
	}

	req := me.client.Post("/api/v2/settings/objects", []settings20.SettingsObjectCreate{soc}).Expect(200)
	resp := []settings20.SettingsObjectCreateResponse{}

	var stub *api.Stub
	if err := req.Finish(&resp); err != nil {
		if strings.Contains(err.Error(), "builtin:span-attribute settings have been replaced by builtin:attribute-allow-list and builtin:attribute-masking") {
			allowlistSettings := allowlistsettings.Settings{
				Enabled: true,
				Key:     v.Key,
			}
			stub, err = allowlist.Service(me.credentials).Create(&allowlistSettings)
			if err != nil {
				return nil, err
			}

			if v.Masking != attribute.MaskingTypes.NotMasked {
				maskingSettings := maskingsettings.Settings{
					Enabled: true,
					Key:     v.Key,
					Masking: maskingsettings.MaskingType(v.Masking),
				}
				if _, err = masking.Service(me.credentials).Create(&maskingSettings); err != nil {
					return nil, err
				}
			}
		} else {
			return nil, err
		}
	} else {
		stub = &api.Stub{ID: resp[0].ObjectID, Name: v.Key}
	}
	return stub, nil
}

func (me *service) UpdateWithContext(ctx context.Context, id string, v *attribute.Settings) error {
	if err := me.client.Get("/api/v2/settings/schemas/builtin%3Aspan-attribute", 200).Finish(); err == nil {
		sou := settings20.SettingsObjectUpdate{Value: v, SchemaVersion: SchemaVersion}
		if err := me.client.Put(fmt.Sprintf("/api/v2/settings/objects/%s", url.PathEscape(id)), &sou, 200).Finish(); err != nil {
			return err
		}
	} else {
		if strings.Contains(err.Error(), "builtin:span-attribute settings have been replaced by builtin:attribute-allow-list and builtin:attribute-masking") {
			var stateKey string
			cfg := ctx.Value(settings.ContextKeyStateConfig)
			if attributeConfig, ok := cfg.(*attribute.Settings); ok {
				stateKey = attributeConfig.Key
			}

			allowlistStubs, err := allowlist.Service(me.credentials).List()
			if err != nil {
				return err
			}
			var allowlistId *string
			for _, stub := range allowlistStubs {
				allowlistStub := stub.Value.(*allowlistsettings.Settings)
				if allowlistStub.Key == stateKey {
					allowlistId = &stub.ID
					break
				}
			}

			if allowlistId != nil {
				allowlistSettings := allowlistsettings.Settings{
					Enabled: true,
					Key:     v.Key,
				}
				if err = allowlist.Service(me.credentials).Update(*allowlistId, &allowlistSettings); err != nil {
					return err
				}
			}

			maskingStubs, err := masking.Service(me.credentials).List()
			if err != nil {
				return err
			}
			var maskingId *string
			for _, stub := range maskingStubs {
				maskingStub := stub.Value.(*maskingsettings.Settings)
				if maskingStub.Key == stateKey {
					maskingId = &stub.ID
					break
				}
			}

			if maskingId != nil {
				if v.Masking != attribute.MaskingTypes.NotMasked {
					maskingSettings := maskingsettings.Settings{
						Enabled: true,
						Key:     v.Key,
						Masking: maskingsettings.MaskingType(v.Masking),
					}
					if err = masking.Service(me.credentials).Update(*maskingId, &maskingSettings); err != nil {
						return err
					}
				} else {
					if err = masking.Service(me.credentials).Delete(*maskingId); err != nil {
						return err
					}
				}
			} else {
				if v.Masking != attribute.MaskingTypes.NotMasked {
					maskingSettings := maskingsettings.Settings{
						Enabled: true,
						Key:     v.Key,
						Masking: maskingsettings.MaskingType(v.Masking),
					}
					if _, err = masking.Service(me.credentials).Create(&maskingSettings); err != nil {
						return err
					}
				}
			}
		} else {
			return err
		}
	}
	return nil
}

func (me *service) Update(id string, v *attribute.Settings) error {
	return errors.New("`builtin:span-attribute` Update function should not be called, please create a GitHub Issue")
}

func (me *service) Validate(v *attribute.Settings) error {
	return nil // Settings 2.0 doesn't offer validation
}

func (me *service) DeleteWithContext(ctx context.Context, id string) error {
	var err error
	if err = me.client.Get("/api/v2/settings/schemas/builtin%3Aspan-attribute", 200).Finish(); err == nil {
		if err = me.client.Delete(fmt.Sprintf("/api/v2/settings/objects/%s", url.PathEscape(id)), 204).Finish(); err != nil {
			return err
		}
	} else {
		if strings.Contains(err.Error(), "builtin:span-attribute settings have been replaced by builtin:attribute-allow-list and builtin:attribute-masking") {
			var stateKey string
			cfg := ctx.Value(settings.ContextKeyStateConfig)
			if attributeConfig, ok := cfg.(*attribute.Settings); ok {
				stateKey = attributeConfig.Key
			}

			allowlistStubs, err := allowlist.Service(me.credentials).List()
			if err != nil {
				return err
			}
			var allowlistId *string
			for _, stub := range allowlistStubs {
				allowlistStub := stub.Value.(*allowlistsettings.Settings)
				if allowlistStub.Key == stateKey {
					allowlistId = &stub.ID
					break
				}
			}
			if allowlistId != nil {
				if err = me.client.Delete(fmt.Sprintf("/api/v2/settings/objects/%s", url.PathEscape(*allowlistId)), 204).Finish(); err != nil {
					return err
				}
			}

			maskingStubs, err := masking.Service(me.credentials).List()
			if err != nil {
				return err
			}
			var maskingId *string
			for _, stub := range maskingStubs {
				maskingStub := stub.Value.(*maskingsettings.Settings)
				if maskingStub.Key == stateKey {
					maskingId = &stub.ID
					break
				}
			}
			if maskingId != nil {
				if err = me.client.Delete(fmt.Sprintf("/api/v2/settings/objects/%s", url.PathEscape(*maskingId)), 204).Finish(); err != nil {
					return err
				}
			}
		} else {
			return err
		}
	}

	return nil
}

func (me *service) Delete(id string) error {
	return errors.New("`builtin:span-attribute` Delete function should not be called, please create a GitHub Issue")
}

func (me *service) GetWithContext(ctx context.Context, id string, v *attribute.Settings) error {
	var err error
	var settingsObject settings20.SettingsObject
	if err = me.client.Get("/api/v2/settings/schemas/builtin%3Aspan-attribute", 200).Finish(); err == nil {
		req := me.client.Get(fmt.Sprintf("/api/v2/settings/objects/%s", url.PathEscape(id))).Expect(200)
		if err = req.Finish(&settingsObject); err != nil {
			return err
		}
		if err = json.Unmarshal(settingsObject.Value, v); err != nil {
			return err
		}
	} else {
		if strings.Contains(err.Error(), "builtin:span-attribute settings have been replaced by builtin:attribute-allow-list and builtin:attribute-masking") {
			var stateKey string
			cfg := ctx.Value(settings.ContextKeyStateConfig)
			if attributeConfig, ok := cfg.(*attribute.Settings); ok {
				stateKey = attributeConfig.Key
			}

			allowlistStubs, err := allowlist.Service(me.credentials).List()
			if err != nil {
				return err
			}
			for _, stub := range allowlistStubs {
				allowlistStub := stub.Value.(*allowlistsettings.Settings)
				if allowlistStub.Key == stateKey {
					v.Key = stateKey
					break
				}
			}
			if v.Key == "" {
				return errors.New("re-run with confighcl")
			}

			v.Masking = attribute.MaskingTypes.NotMasked
			maskingStubs, err := masking.Service(me.credentials).List()
			if err != nil {
				return err
			}
			for _, stub := range maskingStubs {
				maskingStub := stub.Value.(*maskingsettings.Settings)
				if maskingStub.Key == stateKey {
					v.Masking = attribute.MaskingType(maskingStub.Masking)
					break
				}
			}
		} else {
			return err
		}
	}

	settings.SetScope(v, settingsObject.Scope)

	return nil
}

func (me *service) Get(id string, v *attribute.Settings) error {
	return errors.New("`builtin:span-attribute` Get function should not be called, please create a GitHub Issue")
}

func (me *service) List() (api.Stubs, error) {
	var err error

	stubs := api.Stubs{}
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
				newItem := attribute.Settings{}
				if err = json.Unmarshal(item.Value, &newItem); err != nil {
					return nil, err
				}
				settings.SetScope(&newItem, item.Scope)
				itemName := newItem.Name()
				stub := &api.Stub{ID: item.ObjectID, Name: itemName, Value: newItem}
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
