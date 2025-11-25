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

package list

import (
	"context"
	"fmt"
	"net/url"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
)

type GETEntitiesResponse struct {
	TotalCount  int      `json:"totalCount"`
	NextPageKey string   `json:"nextPageKey"`
	Entities    []Entity `json:"entities"`
}

type Entity struct {
	ID   string `json:"entityId"`
	Name string `json:"displayName"`
	Tags []Tag  `json:"tags"`
}

type EntityTag struct {
	StringRepresentation string `json:"stringRepresentation"`
	Value                string `json:"value"`
	Key                  string `json:"key"`
	Context              string `json:"context"`
}

func GETEntities(ctx context.Context, entityType string, client rest.Client, c chan Entity) error {
	var nextPageKey string
	for {
		u := fmt.Sprintf("/api/v2/entities?entitySelector=%s&pageSize=500&from=%s&fields=tags", url.QueryEscape(fmt.Sprintf("type(%s)", entityType)), url.QueryEscape("now-6M"))
		if len(nextPageKey) > 0 {
			u = fmt.Sprintf("/api/v2/entities?nextPageKey=%s", url.QueryEscape(nextPageKey))
		}
		var response GETEntitiesResponse
		err := client.Get(ctx, u, 200).Finish(&response)
		if err != nil {
			close(c)
			return err
		}
		for _, entity := range response.Entities {
			c <- entity
		}
		nextPageKey = response.NextPageKey
		if len(nextPageKey) == 0 {
			break
		}
	}
	close(c)
	return nil
}

func GETEntitiesWithTags(ctx context.Context, entityType string, tags []Tag, client rest.Client, c chan Entity) error {
	entityIds := map[string]string{}
	for _, tag := range tags {
		nextPageKey := ""
		for {
			u := fmt.Sprintf("/api/v2/entities?entitySelector=%s&pageSize=500&from=%s&fields=tags", url.QueryEscape(fmt.Sprintf(`type(%s),tag("%s")`, entityType, tag.StringRepresentation)), url.QueryEscape("now-6M"))
			if len(nextPageKey) > 0 {
				u = fmt.Sprintf("/api/v2/entities?nextPageKey=%s", url.QueryEscape(nextPageKey))
			}
			var response GETEntitiesResponse
			err := client.Get(ctx, u, 200).Finish(&response)
			if err != nil {
				close(c)
				return err
			}
			for _, entity := range response.Entities {
				if _, found := entityIds[entity.ID]; !found {
					entityIds[entity.ID] = entity.ID
					c <- entity
				}
			}
			nextPageKey = response.NextPageKey
			if len(nextPageKey) == 0 {
				break
			}
		}
	}
	close(c)
	return nil
}
