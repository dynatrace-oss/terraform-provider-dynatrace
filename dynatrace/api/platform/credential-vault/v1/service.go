/**
* @license
* Copyright 2026 Dynatrace LLC
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

package vault

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	vault "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/platform/credential-vault/v1/settings"
	smservice "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/http"
	http "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/http/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	"github.com/dynatrace/dynatrace-configuration-as-code-core/clients/credentialvault"
)

func Service(clientSet rest.ClientSet) (settings.CRUDService[*vault.Credentials], error) {
	platformClient, err := clientSet.PlatformClient()
	if err != nil {
		return nil, err
	}

	smService, err := smservice.ServiceSynchronized(clientSet)
	if err != nil {
		return nil, err
	}

	return &service{
		client:    credentialvault.NewClient(platformClient),
		smService: smService,
		envURL:    clientSet.Credentials().Platform.EnvironmentURL,
	}, nil
}

type service struct {
	client    *credentialvault.Client
	smService settings.CRUDService[*http.SyntheticMonitor]
	envURL    string
}

type entryReference struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (me *service) SchemaID() string {
	return "platform:credential-vault"
}

func (me *service) Create(ctx context.Context, v *vault.Credentials) (*api.Stub, error) {
	data, err := json.Marshal(v)

	if err != nil {
		return nil, err
	}

	raw, err := me.client.Create(ctx, data)
	if err != nil {
		return nil, err
	}

	var response entryReference
	if err = json.Unmarshal(raw.Data, &response); err != nil {
		return nil, err
	}

	return &api.Stub{Name: v.Name, ID: response.ID}, nil
}

func (me *service) Update(ctx context.Context, id string, v *vault.Credentials) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	if _, err = me.client.Update(ctx, id, data); err != nil {
		return err
	}

	return nil
}

func (me *service) Get(ctx context.Context, id string, v *vault.Credentials) error {
	response, err := me.client.Get(ctx, id)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(response.Data, v); err != nil {
		return err
	}

	return nil
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	result, err := me.client.List(ctx, credentialvault.ListOptions{})
	if err != nil {
		return nil, err
	}
	var stubs api.Stubs
	for _, raw := range result.All() {
		var entry entryReference
		if err := json.Unmarshal(raw, &entry); err != nil {
			return nil, err
		}
		stubs = append(stubs, &api.Stub{ID: entry.ID, Name: entry.Name})
	}
	return stubs, nil
}

func (me *service) Delete(ctx context.Context, id string) error {
	// Corner case (delete right after create): The monitor may not be available right after creation
	syntheticErr := me.deleteSyntheticMonitor(ctx, id)
	if syntheticErr != nil {
		syntheticErr = fmt.Errorf("associated HTTP monitor couldn't be deleted: %w", syntheticErr)
	}

	if _, err := me.client.Delete(ctx, id); err != nil {
		return errors.Join(err, syntheticErr)
	}
	if syntheticErr != nil {
		return rest.Warning{Message: syntheticErr.Error()}
	}
	return nil
}

// deleteSyntheticMonitor deletes the Synthetic HTTP monitor that external vault synchronization creates.
// Its name ends with " (<credential-id>)"
func (me *service) deleteSyntheticMonitor(ctx context.Context, id string) error {
	nameSuffixes := []string{fmt.Sprintf(" (%s)", id), fmt.Sprintf(" (%s)", id)}
	monitors, err := me.smService.List(ctx)
	if err != nil {
		return err
	}
	for _, monitor := range monitors {
		if slices.ContainsFunc(nameSuffixes, func(suffix string) bool { return strings.HasSuffix(monitor.Name, suffix) }) {
			if err = me.smService.Delete(ctx, monitor.ID); err != nil && !rest.IsNotFoundError(err) {
				return err
			}
			return nil
		}
	}
	return nil
}
