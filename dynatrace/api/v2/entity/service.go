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

package entity

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	entity "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/entity/settings"
)

const SchemaID = "v2:environment:entity"

func Service(credentials *settings.Credentials) settings.RService[*entity.Entity] {
	return &service{client: rest.DefaultClient(credentials.URL, credentials.Token)}
}

type service struct {
	client rest.Client
}

func (me *service) Get(ctx context.Context, id string, v *entity.Entity) error {
	if len(os.Getenv("DYNATRACE_MIGRATION_CACHE_FOLDER")) > 0 {
		return errors.New("entity calls disabled when using Migration Cache")
	}
	return me.client.Get(fmt.Sprintf(`/api/v2/entities/%s?from=%s`, url.PathEscape(id), url.QueryEscape("now-3y")), 200).Finish(v)
}

func (me *service) SchemaID() string {
	return SchemaID
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return api.Stubs{&api.Stub{ID: me.SchemaID(), Name: me.SchemaID()}}, nil
}

func DataSourceService(credentials *settings.Credentials) settings.RService[*entity.Entity] {
	return &dataSourceService{client: rest.DefaultClient(credentials.URL, credentials.Token)}
}

type dataSourceService struct {
	client rest.Client
}

func (me *dataSourceService) Get(ctx context.Context, id string, v *entity.Entity) error {
	if len(os.Getenv("DYNATRACE_MIGRATION_CACHE_FOLDER")) > 0 {
		return errors.New("entity calls disabled when using Migration Cache")
	}

	entityType := evalEntityType(id)
	if len(entityType) == 0 {
		return me.client.Get(fmt.Sprintf(`/api/v2/entities/%s?from=%s&fields=tags`, url.PathEscape(id), url.QueryEscape("now-3y")), 200).Finish(v)
	}

	result := getEntity(id, me.client, getEntitiesRecord(entityType))
	if result == nil {
		return rest.Error{Code: 404, Message: fmt.Sprintf("Unable to find entity with id %s", id)}
	}
	v.DisplayName = result.DisplayName
	v.EntityId = result.EntityId
	v.Type = result.Type
	return nil

}

func (me *dataSourceService) SchemaID() string {
	return SchemaID
}

func (me *dataSourceService) List(ctx context.Context) (api.Stubs, error) {
	return api.Stubs{&api.Stub{ID: me.SchemaID(), Name: me.SchemaID()}}, nil
}

var entityCache = EntityCache{
	Lock:     sync.RWMutex{},
	Entities: map[string]*EntitiesRecord{},
}

type EntityCache struct {
	Lock     sync.RWMutex
	Entities map[string]*EntitiesRecord
}

type EntitiesRecord struct {
	Lock       sync.Mutex
	EntityType string
	Entities   map[string]*entity.Entity
}

func getEntitiesRecord(entityType string) *EntitiesRecord {
	entityCache.Lock.RLock()
	defer entityCache.Lock.RUnlock()
	record, found := entityCache.Entities[entityType]
	if !found {
		record := &EntitiesRecord{
			Lock:       sync.Mutex{},
			EntityType: entityType,
			Entities:   nil,
		}
		entityCache.Entities[entityType] = record
		return record
	}
	return record
}

type EntitiesListResponse struct {
	TotalCount int `json:"totalCount"`
	PageSize   int `json:"pageSize"`
	Entities   []struct {
		ID          string `json:"entityId"`
		Type        string `json:"type"`
		DisplayName string `json:"displayName"`
		Tags        []struct {
			Context              string  `json:"context"`
			Key                  string  `json:"key"`
			Value                *string `json:"value,omitempty"`
			StringRepresentation *string `json:"stringRepresentation,omitempty"`
		} `json:"tags,omitempty"`
	} `json:"entities"`
	NextPageKey string `json:"nextPageKey"`
}

func getEntity(id string, client rest.Client, record *EntitiesRecord) *entity.Entity {
	record.Lock.Lock()
	defer record.Lock.Unlock()
	if record.Entities != nil {
		entity, found := record.Entities[id]
		if !found {
			return nil
		}
		return entity
	}
	entities := fetchEntities(client, record)
	if entities == nil {
		return nil
	}
	record.Entities = entities
	entity, found := record.Entities[id]
	if !found {
		return nil
	}
	return entity
}

func fetchEntities(client rest.Client, record *EntitiesRecord) map[string]*entity.Entity {
	results := map[string]*entity.Entity{}
	nextPageKey := "-"
	for len(nextPageKey) > 0 {
		entitySelector := fmt.Sprintf("type(%s)", record.EntityType)
		var u string
		if nextPageKey != "-" {
			u = fmt.Sprintf("/api/v2/entities?nextPageKey=%s", url.QueryEscape(nextPageKey))
		} else {
			u = fmt.Sprintf("/api/v2/entities?pageSize=4000&entitySelector=%s&from=%s&fields=tags", url.QueryEscape(entitySelector), url.QueryEscape("now-3y"))
		}
		var response EntitiesListResponse
		if err := client.Get(u, 200).Finish(&response); err != nil {
			return nil
		}
		for _, elem := range response.Entities {
			entity := &entity.Entity{
				EntityId:    &elem.ID,
				DisplayName: &elem.DisplayName,
				Type:        &elem.Type,
			}
			results[elem.ID] = entity
		}
		nextPageKey = response.NextPageKey
	}
	return results
}

func evalEntityType(id string) string {
	if strings.HasPrefix(id, "HOST-") {
		return "HOST"
	}
	if strings.HasPrefix(id, "HOST_GROUP-") {
		return "HOST_GROUP"
	}
	if strings.HasPrefix(id, "PROCESS_GROUP-") {
		return "PROCESS_GROUP"
	}
	if strings.HasPrefix(id, "PROCESS_GROUP_INSTANCE-") {
		return "PROCESS_GROUP_INSTANCE"
	}
	if strings.HasPrefix(id, "SERVICE-") {
		return "SERVICE"
	}
	if strings.HasPrefix(id, "KUBERNETES_CLUSTER-") {
		return "KUBERNETES_CLUSTER"
	}
	if strings.HasPrefix(id, "APPLICATION_METHOD-") {
		return "APPLICATION_METHOD"
	}
	if strings.HasPrefix(id, "DISK-") {
		return "DISK"
	}

	return ""
}
