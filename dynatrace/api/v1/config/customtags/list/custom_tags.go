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

func GETCustomTags(ctx context.Context, entityID string, client rest.Client) ([]Tag, error) {
	u := fmt.Sprintf("/api/v2/tags?entitySelector=%s&from=%s", url.QueryEscape(fmt.Sprintf("entityId(%s)", entityID)), url.QueryEscape("now-6M"))
	var response GetCustomTagsResponse
	err := client.Get(ctx, u, 200).Finish(&response)
	if err != nil {
		return nil, err
	}
	return response.Tags, nil
}
