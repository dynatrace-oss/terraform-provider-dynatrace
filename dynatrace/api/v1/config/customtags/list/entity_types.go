package list

import (
	"fmt"
	"net/url"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
)

type GetEntityTypesResponse struct {
	TotalCount  int          `json:"totalCount"`
	NextPageKey string       `json:"nextPageKey"`
	Types       []EntityType `json:"types"`
}

type EntityType struct {
	DimensionKey        string               `json:"dimensionKey,omitempty"`
	Type                string               `json:"type,omitempty"`
	FromRelationships   any                  `json:"fromRelationships,omitempty"`
	ToRelationships     any                  `json:"toRelationships,omitempty"`
	EntityLimitExceeded bool                 `json:"entityLimitExceeded,omitempty"`
	ManagementZones     string               `json:"managementZones,omitempty"`
	Tags                string               `json:"tags,omitempty"`
	DisplayName         string               `json:"displayName,omitempty"`
	Properties          []EntityTypeProperty `json:"properties,omitempty"`
}

type EntityTypeProperty struct {
	ID          string `json:"id,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Type        string `json:"type,omitempty"`
}

func GETEntityTypes(client rest.Client) ([]EntityType, error) {
	entityTypes := []EntityType{}
	var nextPageKey string
	for {
		u := fmt.Sprintf("/api/v2/entityTypes?pageSize=500")
		if len(nextPageKey) > 0 {
			u = fmt.Sprintf("/api/v2/entityTypes?nextPageKey=%s", url.QueryEscape(nextPageKey))
		}
		var response GetEntityTypesResponse
		err := client.Get(u, 200).Finish(&response)
		if err != nil {
			return nil, err
		}
		entityTypes = append(entityTypes, response.Types...)
		nextPageKey = response.NextPageKey
		if len(nextPageKey) == 0 {
			break
		}
	}

	return entityTypes, nil
}

func (et EntityType) GetCustomTags(client rest.Client) ([]Tag, error) {
	u := fmt.Sprintf("/api/v2/tags?entitySelector=%s&from=%s", url.QueryEscape(fmt.Sprintf("type(%s)", et.Type)), url.QueryEscape(TIME_FRAME))
	var response GetCustomTagsResponse
	err := client.Get(u, 200).Finish(&response)
	if err != nil {
		return nil, err
	}
	return response.Tags, nil
}
