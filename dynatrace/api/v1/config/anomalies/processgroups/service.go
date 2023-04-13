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

package processgroups

import (
	"fmt"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/cache"

	processgroups "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/processgroups/settings"
	entities "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/entities"
	entitiesSettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/entities/settings"
)

const SchemaID = "v1:config:anomaly-detection:process-groups"
const BasePath = "/api/config/v1/anomalyDetection/processGroups"

func Service(credentials *settings.Credentials) settings.CRUDService[*processgroups.AnomalyDetection] {
	return &service{client: rest.DefaultClient(credentials.URL, credentials.Token), credentials: credentials}
}

type service struct {
	credentials *settings.Credentials
	client      rest.Client
}

func (me *service) SchemaID() string {
	return SchemaID
}

func (me *service) Create(config *processgroups.AnomalyDetection) (*api.Stub, error) {
	if err := me.Update(config.ProcessGroupId, config); err != nil {
		return nil, err
	}
	return &api.Stub{ID: config.ProcessGroupId + "-anomalydetection", Name: config.ProcessGroupId + "-anomalydetection"}, nil
}

// Update TODO: documentation

func (me *service) Update(id string, config *processgroups.AnomalyDetection) error {
	id = strings.TrimSuffix(id, "-anomalydetection")

	if err := me.client.Put(fmt.Sprintf("/api/config/v1/anomalyDetection/processGroups/%s", id), config, 204).Finish(); err != nil {
		return err
	}

	return nil
}

// Validate TODO: documentation
func (me *service) Validate(config *processgroups.AnomalyDetection) error {
	id := strings.TrimSuffix(config.ProcessGroupId, "-anomalydetection")

	if err := me.client.Put(fmt.Sprintf("/api/config/v1/anomalyDetection/processGroups/%s/validator", id), config, 204).Finish(); err != nil {
		return err
	}

	return nil
}

// Delete TODO: documentation
func (me *service) Delete(id string) error {
	id = strings.TrimSuffix(id, "-anomalydetection")

	if err := me.client.Delete(fmt.Sprintf("/api/config/v1/anomalyDetection/processGroups/%s", id), 204).Finish(); err != nil {
		return err
	}

	return nil
}

// Get TODO: documentation
func (me *service) Get(id string, v *processgroups.AnomalyDetection) error {
	id = strings.TrimSuffix(id, "-anomalydetection")

	if err := me.client.Get(fmt.Sprintf("/api/config/v1/anomalyDetection/processGroups/%s", id), 200).Finish(v); err != nil {
		return err
	}
	v.ProcessGroupId = id

	return nil
}

func (me *service) List() (api.Stubs, error) {
	srv := cache.Read(entities.Service("PROCESS_GROUP", me.credentials), true)
	v := new(entitiesSettings.Settings)
	if err := srv.Get("", v); err != nil {
		return nil, err
	}

	var stubs api.Stubs
	for _, entity := range v.Entities {
		stub := api.Stub{ID: *entity.EntityId, Name: *entity.DisplayName}
		stubs = append(stubs, &stub)
	}

	return stubs, nil
}

func (me *service) New() *processgroups.AnomalyDetection {
	return new(processgroups.AnomalyDetection)
}

func (me *service) Name() string {
	return me.SchemaID()
}
