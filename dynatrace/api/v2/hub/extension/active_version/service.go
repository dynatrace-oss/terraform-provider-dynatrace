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

package active_environment_config

import (
	"context"
	"fmt"
	"net/url"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	active_version "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/hub/extension/active_version/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/hub/items"
	items_settings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/hub/items/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

func Service(credentials *settings.Credentials) settings.CRUDService[*active_version.Settings] {
	return &service{credentials: credentials, itemsService: items.Service(credentials, items.Options{Type: "EXTENSION2"})}
}

type service struct {
	credentials  *settings.Credentials
	itemsService settings.RService[*items_settings.HubItemList]
}

func (me *service) Get(ctx context.Context, id string, v *active_version.Settings) error {
	var response GetActiveEnvironmentConfigurationResponse
	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	if err := client.Get(fmt.Sprintf("/api/v2/extensions/%s/environmentConfiguration", url.PathEscape(id)), 200).Finish(&response); err != nil {
		return err
	}
	v.Version = response.Version
	v.Name = id
	return nil
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	var stubs api.Stubs
	var hubItemList items_settings.HubItemList
	if err := me.itemsService.Get(ctx, "", &hubItemList); err != nil {
		return nil, err
	}
	for _, item := range hubItemList.Items {
		name := item.ArtifactId
		if len(name) == 0 {
			continue
		}
		sttngs := active_version.Settings{}
		if err := me.Get(ctx, name, &sttngs); err == nil {
			if len(sttngs.Version) > 0 {
				stubs = append(stubs, &api.Stub{ID: name, Name: name})
			}
		}
	}
	return stubs, nil
}

func (me *service) Create(ctx context.Context, v *active_version.Settings) (*api.Stub, error) {
	version := v.Version
	name := v.Name
	if err := me.ensureInstalled(name, version); err != nil {
		return nil, err
	}
	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	createResponse := SetActiveEnvironmentConfigurationResponse{}
	retry := 10
	for retry > 0 {
		payload := SetActiveEnvironmentConfigurationRequest{Version: version}
		if err := client.Put(fmt.Sprintf("/api/v2/extensions/%s/environmentConfiguration", url.PathEscape(name)), &payload, 200).Finish(&createResponse); err != nil {
			return nil, err
		}
		retry = 0
	}
	return &api.Stub{ID: name, Name: name}, nil
}

func (me *service) ensureInstalled(name string, version string) error {
	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	response := struct {
		Name    string `json:"extensionName"`
		Version string `json:"extensionVersion"`
	}{}
	if err := client.Post(fmt.Sprintf("/api/v2/extensions/%s?version=%s", url.PathEscape(name), url.QueryEscape(version)), nil, 200).Finish(&response); err != nil {
		if restErr, ok := err.(rest.Error); ok {
			if (restErr.Code == 400) && restErr.Message == fmt.Sprintf("Extension %s has already been added to environment", name) {
				return nil
			}
		}
		return err
	}
	return nil
}

func (me *service) Update(ctx context.Context, id string, v *active_version.Settings) error {
	_, err := me.Create(ctx, v)
	return err
}

func (me *service) Delete(ctx context.Context, id string) error {
	// we cannot really delete this
	return nil
}

func (me *service) SchemaID() string {
	return "v2:extensions:active-env-config"
}

func (me *service) Validate(v *active_version.Settings) error {
	return nil // no endpoint for that
}

func (me *service) New() *active_version.Settings {
	return new(active_version.Settings)
}

type GetActiveEnvironmentConfigurationResponse struct {
	Version string `json:"version"`
}

type SetActiveEnvironmentConfigurationRequest struct {
	Version string `json:"version"`
}

type SetActiveEnvironmentConfigurationResponse struct {
	Version string `json:"version"`
}
