/**
* @license
* Copyright 2025 Dynatrace LLC
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

package ingestsource2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	ingestsource "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/ingestsource2/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
)

const SchemaFormat = "builtin:openpipeline.%s.ingest-sources"

var SchemaRegex = regexp.MustCompile("^builtin:openpipeline.(.*?).ingest-sources$")

func Service(credentials *rest.Credentials) settings.CRUDService[*ingestsource.IngestSource] {
	return &ServiceImpl{credentials: credentials}
}

type ServiceImpl struct {
	credentials            *rest.Credentials
	GenericSettingsService *settings20.Client // settings.CRUDService[*ingestsource.IngestSource]
	SettingsServicePerKind map[string]settings.CRUDService[*ingestsource.IngestSource]
}

func (si *ServiceImpl) SchemaID() string {
	return "openpipelinev2:ingest-source"
}

func (si *ServiceImpl) getClient(ctx context.Context) *settings20.Client { // settings.CRUDService[*ingestsource.IngestSource]
	if si.GenericSettingsService == nil {
		tokenClient, _ := rest.CreateClassicClient(si.credentials.URL, si.credentials.Token)
		oauthClient, _ := rest.CreateClassicOAuthBasedClient(ctx, si.credentials)
		si.GenericSettingsService = settings20.NewClient(tokenClient, oauthClient, "")
	}

	return si.GenericSettingsService
}

func (si *ServiceImpl) getClientForKind(kind string) settings.CRUDService[*ingestsource.IngestSource] {
	if si.SettingsServicePerKind == nil {
		si.SettingsServicePerKind = make(map[string]settings.CRUDService[*ingestsource.IngestSource])
	}

	client, ok := si.SettingsServicePerKind[kind]
	if ok && client != nil {
		return client
	}

	client = settings20.Service[*ingestsource.IngestSource](si.credentials, fmt.Sprintf(SchemaFormat, kind), "")
	si.SettingsServicePerKind[kind] = client
	return client
}

func (si *ServiceImpl) Get(ctx context.Context, objectID string, v *ingestsource.IngestSource) error {
	return si.get(ctx, objectID, v)
}

func (si *ServiceImpl) get(ctx context.Context, objectID string, v *ingestsource.IngestSource) error {
	type SettingsObject struct {
		SchemaVersion string          `json:"schemaVersion"`
		SchemaID      string          `json:"schemaId"`
		Scope         string          `json:"scope"`
		Value         json.RawMessage `json:"value"`
	}

	var response settings20.Response
	var settingsObject SettingsObject

	response, err := si.getClient(ctx).Get(ctx, objectID)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		if err = rest.Envelope(response.Data, response.Request.URL, response.Request.Method); err != nil {
			return err
		}
		return fmt.Errorf("status code %d (expected: %d): %s", response.StatusCode, 200, string(response.Data))
	}
	if err = json.Unmarshal(response.Data, &settingsObject); err != nil {
		return err
	}

	if err = json.Unmarshal(settingsObject.Value, v); err != nil {
		return err
	}

	res := SchemaRegex.FindStringSubmatch(settingsObject.SchemaID)
	// TODO: error handling, validation, tests.
	v.Kind = res[1]
	return nil
}

func (si *ServiceImpl) List(ctx context.Context) (api.Stubs, error) {
	var stubs api.Stubs
	stubs = make(api.Stubs, 0)
	var errs []error
	for _, kind := range ingestsource.AllowedKinds {
		client := si.getClientForKind(kind)
		res, err := client.List(ctx)

		errs = append(errs, err)
		stubs = append(stubs, res...)
	}

	if len(errs) > 0 {
		return stubs, errors.Join(errs...)
	}
	return stubs, nil
}

func (si *ServiceImpl) Create(ctx context.Context, v *ingestsource.IngestSource) (*api.Stub, error) {
	client := si.getClientForKind(v.Kind)
	return client.Create(ctx, v)
}

// Update will not be called if "kind" changes, because of ForceNew: true
// Therefore we do not need to care about changes in "kind" here, which would otherwise
// entail a settings deletion and re-creation
func (si *ServiceImpl) Update(ctx context.Context, id string, v *ingestsource.IngestSource) error {
	return si.getClientForKind(v.Kind).Update(ctx, id, v)
}

func (si *ServiceImpl) Delete(ctx context.Context, id string) error {
	client := si.getClient(ctx)

	_, err := client.Delete(ctx, id)
	return err
}
