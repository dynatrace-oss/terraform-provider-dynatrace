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

package ingestsource

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	ingestsource "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/ingestsource/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
)

const schemaFormat = "builtin:openpipeline.%s.ingest-sources"

var schemaRegex = regexp.MustCompile("^builtin:openpipeline\\.(.*?)\\.ingest-sources$")

func Service(credentials *rest.Credentials) settings.CRUDService[*ingestsource.IngestSource] {
	return &ServiceImpl{Credentials: credentials}
}

type SettingsObject struct {
	SchemaID string          `json:"schemaId"`
	Value    json.RawMessage `json:"value"`
}

type GenericSettingsClient interface {
	Get(ctx context.Context, id string) (settings20.Response, error)
	Delete(ctx context.Context, id string) (settings20.Response, error)
}

type SettingsClientForKind interface {
	List(ctx context.Context) (api.Stubs, error)
	Create(ctx context.Context, v *ingestsource.IngestSource) (*api.Stub, error)
	Update(ctx context.Context, id string, v *ingestsource.IngestSource) error
}

type ServiceImpl struct {
	Credentials            *rest.Credentials
	GenericSettingsClient  GenericSettingsClient
	SettingsClientsPerKind map[string]SettingsClientForKind
}

func (si *ServiceImpl) getClient(ctx context.Context) GenericSettingsClient {
	if si.GenericSettingsClient == nil {
		tokenClient, _ := rest.CreateClassicClient(si.Credentials.URL, si.Credentials.Token)
		oauthClient, _ := rest.CreateClassicOAuthBasedClient(ctx, si.Credentials)
		si.GenericSettingsClient = settings20.NewClient(tokenClient, oauthClient, "")
	}

	return si.GenericSettingsClient
}

func (si *ServiceImpl) getClientForKind(kind string) SettingsClientForKind {
	if si.SettingsClientsPerKind == nil {
		si.SettingsClientsPerKind = make(map[string]SettingsClientForKind)
	}

	client, ok := si.SettingsClientsPerKind[kind]
	if ok && client != nil {
		return client
	}

	client = settings20.Service[*ingestsource.IngestSource](si.Credentials, fmt.Sprintf(schemaFormat, kind), "")
	si.SettingsClientsPerKind[kind] = client
	return client
}

func parseKindFromSchemaID(schemaID string) (kind string, found bool) {
	res := schemaRegex.FindStringSubmatch(schemaID)
	if res == nil {
		return "", false
	}

	return res[1], true
}

func (si *ServiceImpl) Get(ctx context.Context, objectID string, v *ingestsource.IngestSource) error {
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

	kind, ok := parseKindFromSchemaID(settingsObject.SchemaID)
	if !ok {
		return fmt.Errorf("could not parse kind from schema id '%s'", settingsObject.SchemaID)
	}

	v.Kind = kind
	return nil
}

func (si *ServiceImpl) List(ctx context.Context) (api.Stubs, error) {
	stubs := make(api.Stubs, 0)
	var errs []error
	for _, kind := range ingestsource.AllowedKinds {
		client := si.getClientForKind(kind)
		res, err := client.List(ctx)

		if err == nil {
			stubs = append(stubs, res...)
		} else {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return stubs, nil
}

func (si *ServiceImpl) Create(ctx context.Context, v *ingestsource.IngestSource) (*api.Stub, error) {
	return si.getClientForKind(v.Kind).Create(ctx, v)
}

// Update will not be called if "kind" changes, because of ForceNew: true
// Therefore we do not need to care about changes in "kind" here, which would otherwise
// entail a settings deletion and re-creation
func (si *ServiceImpl) Update(ctx context.Context, id string, v *ingestsource.IngestSource) error {
	return si.getClientForKind(v.Kind).Update(ctx, id, v)
}

func (si *ServiceImpl) Delete(ctx context.Context, id string) error {
	_, err := si.getClient(ctx).Delete(ctx, id)
	return err
}

func (si *ServiceImpl) SchemaID() string {
	return "openpipelinev2:ingest-source"
}
