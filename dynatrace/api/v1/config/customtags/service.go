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

package customtags

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	customtags "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/customtags/settings"
)

func Service(credentials *settings.Credentials) settings.CRUDService[*customtags.Settings] {
	return &service{credentials: credentials}
}

type service struct {
	credentials *settings.Credentials
}

func (me *service) Get(selector string, v *customtags.Settings) (err error) {
	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	if err = client.Get(fmt.Sprintf("/api/v2/tags?entitySelector=%s", url.QueryEscape(selector)), 200).Finish(v); err != nil {
		return err
	}
	v.EntitySelector = selector

	return nil
}

func (me *service) SchemaID() string {
	return "v2:environment:customtags"
}

func (me *service) List() (settings.Stubs, error) {
	return settings.Stubs{}, nil
}

func (me *service) Validate(v *customtags.Settings) error {
	return nil // no endpoint for that
}

func (me *service) Create(v *customtags.Settings) (*settings.Stub, error) {
	if err := me.Update(v.EntitySelector, v); err != nil {
		return nil, err
	}
	return &settings.Stub{ID: v.EntitySelector, Name: v.EntitySelector}, nil
}

func (me *service) Update(id string, v *customtags.Settings) error {
	var err error

	var settingsObj customtags.Settings
	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	if err = client.Post(fmt.Sprintf("/api/v2/tags?entitySelector=%s", url.QueryEscape(v.EntitySelector)), v, 200).Finish(&settingsObj); err != nil {
		return err
	}
	v.MatchedEntities = settingsObj.MatchedEntities

	return nil
}

func (me *service) Delete(id string) error {
	return errors.New("not implemented")
}

func (me *service) DeleteValue(v *customtags.Settings) error {
	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	for _, tag := range v.Tags {
		if tag.Value == nil || len(*tag.Value) == 0 {
			if err := client.Delete(fmt.Sprintf("/api/v2/tags?key=%s&entitySelector=%s", url.QueryEscape(tag.Key), url.QueryEscape(v.EntitySelector)), 200).Finish(); err != nil {
				return err
			}
		} else {
			if err := client.Delete(fmt.Sprintf("/api/v2/tags?key=%s&value=%s&entitySelector=%s", url.QueryEscape(tag.Key), url.QueryEscape(*tag.Value), url.QueryEscape(v.EntitySelector)), 200).Finish(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (me *service) New() *customtags.Settings {
	return new(customtags.Settings)
}

func (me *service) Name() string {
	return me.SchemaID()
}
