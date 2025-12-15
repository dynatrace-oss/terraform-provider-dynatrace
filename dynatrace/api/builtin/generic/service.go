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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	generic "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/generic/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/shutdown"
)

const settingsObjectEndpoint = "/api/v2/settings/objects"
const settingsSchemaEndpoint = "/api/v2/settings/schemas"

func Service(credentials *rest.Credentials) settings.CRUDService[*generic.Settings] {
	return &service{credentials: credentials}
}

type service struct {
	credentials *rest.Credentials
}

func (me *service) Client() rest.Client {
	return rest.HybridClient(me.credentials)
}

func (me *service) Get(ctx context.Context, id string, v *generic.Settings) error {
	var settingsObject settings20.SettingsObject
	client := me.Client()

	u, err := url.JoinPath(settingsObjectEndpoint, id)
	if err != nil {
		return newSettingsURLJoinError(err)
	}
	err = client.Get(ctx, u, 200).Finish(&settingsObject)
	if err != nil {
		return err
	}

	v.Value = string(settingsObject.Value)
	v.Scope = settingsObject.Scope
	v.SchemaID = settingsObject.SchemaID

	return nil
}

type schemaStub struct {
	SchemaID string `json:"schemaId"`
}

type schemataResponse struct {
	Items []schemaStub `json:"items"`
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	client := me.Client()
	var schemata schemataResponse
	err := client.Get(ctx, settingsSchemaEndpoint, 200).Finish(&schemata)
	if err != nil {
		return nil, err
	}
	if len(schemata.Items) == 0 {
		return api.Stubs{}, nil
	}
	schemaIDs := []string{}
	for _, schemaStub := range schemata.Items {
		if strings.HasPrefix(schemaStub.SchemaID, "app:") {
			schemaIDs = append(schemaIDs, schemaStub.SchemaID)
		}
	}
	if len(schemaIDs) == 0 {
		return api.Stubs{}, nil
	}
	var stubs api.Stubs
	for _, schemaID := range schemaIDs {
		newStubs, err := me.list(ctx, client, schemaID, "")
		if err != nil {
			return nil, err
		}

		stubs = append(stubs, newStubs...)
	}
	return stubs, nil
}

func (me *service) list(ctx context.Context, client rest.Client, schemaID string, nextPageKey string) (api.Stubs, error) {
	var queryParams url.Values
	if nextPageKey == "" {
		queryParams = url.Values{
			"schemaIds": []string{schemaID},
			"fields":    []string{"objectId"},
			"pageSize":  []string{"500"},
		}
	} else {
		queryParams = url.Values{
			"nextPageKey": []string{nextPageKey},
		}
	}

	u := url.URL{
		Path:     settingsObjectEndpoint,
		RawQuery: queryParams.Encode(),
	}

	var objectsResponse settings20.SettingsObjectList
	err := client.Get(ctx, u.String(), 200).Finish(&objectsResponse)

	if err != nil {
		return nil, err
	}

	var stubs api.Stubs
	for _, item := range objectsResponse.Items {
		stubs = append(stubs, &api.Stub{ID: item.ObjectID, Name: item.ObjectID})
	}

	if objectsResponse.NextPageKey != nil && *objectsResponse.NextPageKey != "" {
		newStubs, err := me.list(ctx, client, schemaID, *objectsResponse.NextPageKey)
		if err != nil {
			return nil, err
		}
		stubs = append(stubs, newStubs...)
	}

	return stubs, nil
}

func (me *service) Validate(v *generic.Settings) error {
	return nil // Settings 2.0 doesn't offer validation
}

func (me *service) Create(ctx context.Context, v *generic.Settings) (*api.Stub, error) {
	stubs, err := me.create(ctx, v)
	if err == nil {
		return stubs, nil
	}
	if rest.IsRequiresOAuthError(err) && me.credentials.ContainsOAuthOrPlatformToken() {
		ctx := rest.NewPreferOAuthContext(ctx)
		return me.create(ctx, v)
	}
	return nil, err
}

func (me *service) create(ctx context.Context, v *generic.Settings) (*api.Stub, error) {
	scope := "environment"
	if len(v.Scope) > 0 {
		scope = v.Scope
	}
	obj := []settings20.SettingsObjectCreate{{
		SchemaID: v.SchemaID,
		Scope:    scope,
		Value:    json.RawMessage(v.Value),
	}}

	var response []settings20.SettingsObjectCreateResponse
	err := me.Client().Post(ctx, settingsObjectEndpoint, obj).Finish(&response)
	if err != nil {
		return nil, err
	}

	if len(response) == 0 {
		return nil, errors.New("failed to create settings object: no ID returned")
	}

	id := response[0].ObjectID
	stub := &api.Stub{ID: id, Name: id}
	return stub, nil
}

func (me *service) Update(ctx context.Context, id string, v *generic.Settings) error {
	err := me.update(ctx, id, v)
	if err == nil {
		return nil
	}
	if rest.IsRequiresOAuthError(err) && me.credentials.ContainsOAuthOrPlatformToken() {
		ctx := rest.NewPreferOAuthContext(ctx)
		return me.update(ctx, id, v)
	}
	return err
}

func (me *service) update(ctx context.Context, id string, v *generic.Settings) error {
	obj := settings20.SettingsObjectUpdate{
		Value: json.RawMessage(v.Value),
	}

	u, err := url.JoinPath(settingsObjectEndpoint, id)
	if err != nil {
		return newSettingsURLJoinError(err)
	}
	return me.Client().Put(ctx, u, obj).Finish(nil)
}

func (me *service) Delete(ctx context.Context, id string) error {
	client := me.Client()

	u, err := url.JoinPath(settingsObjectEndpoint, id)
	if err != nil {
		return newSettingsURLJoinError(err)
	}

	err = client.Delete(ctx, u).Finish(nil)
	if err != nil {
		return err
	}

	return nil
}

func (me *service) SchemaID() string {
	return "generic"
}

type QueryParams struct {
	Schema string
	Scope  string
	Filter string
}

func (me *service) ListSpecific(ctx context.Context, query QueryParams) (api.Stubs, error) {
	client := me.Client()

	stubs := api.Stubs{}
	nextPage := true
	var nextPageKey *string
	for nextPage {
		var sol settings20.SettingsObjectList
		var queryValues url.Values
		if nextPageKey == nil {
			queryValues = url.Values{
				"fields":    []string{"objectId,scope,value,schemaId"},
				"pageSize":  []string{"100"},
				"schemaIds": []string{query.Schema},
				"scopes":    []string{query.Scope},
				"filter":    []string{query.Filter},
			}
		} else {
			queryValues = url.Values{
				"nextPageKey": []string{*nextPageKey},
			}
		}
		u := url.URL{
			Path:     settingsObjectEndpoint,
			RawQuery: queryValues.Encode(),
		}

		err := client.Get(ctx, u.String()).Finish(&sol)
		if err != nil {
			return nil, err
		}

		if len(sol.Items) == 0 {
			return api.Stubs{}, nil
		}
		if shutdown.System.Stopped() {
			return stubs, nil
		}

		if len(sol.Items) > 0 {
			for _, item := range sol.Items {
				newItem := new(generic.Settings)
				newItem.Value = string(item.Value)
				newItem.Scope = item.Scope
				newItem.SchemaID = item.SchemaID
				stubs = append(stubs, &api.Stub{ID: item.ObjectID, Name: item.ObjectID, Value: newItem})
			}
		}
		nextPageKey = sol.NextPageKey
		nextPage = (nextPageKey != nil)
	}

	return stubs, nil
}

func newSettingsURLJoinError(err error) error {
	return fmt.Errorf("failed to create settings object URL: %w", err)
}
