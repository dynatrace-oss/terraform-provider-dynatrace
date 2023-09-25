package list

import (
	"fmt"
	"net/url"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
)

type GetCustomTagsResponse struct {
	TotalCount int   `json:"totalCount"`
	Tags       []Tag `json:"tags"`
}

type Tag struct {
	StringRepresentation string `json:"stringRepresentation"`
	Value                string `json:"value"`
	Key                  string `json:"key"`
	Context              string `json:"context"`
}

func GETCustomTags(entityID string, client rest.Client) ([]Tag, error) {
	u := fmt.Sprintf("/api/v2/tags?entitySelector=%s&from=%s", url.QueryEscape(fmt.Sprintf("entityId(%s)", entityID)), url.QueryEscape("now-6M"))
	var response GetCustomTagsResponse
	err := client.Get(u, 200).Finish(&response)
	if err != nil {
		return nil, err
	}
	return response.Tags, nil
}
