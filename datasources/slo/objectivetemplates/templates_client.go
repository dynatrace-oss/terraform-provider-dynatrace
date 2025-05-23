/*
 * @license
 * Copyright 2023 Dynatrace LLC
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package objectivetemplates

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
)

const (
	templatesResourcePath = "/platform/slo/v1/objective-templates"
)

type SLOObjectiveTemplatesClient struct {
	client *rest.Client
}

func NewSLOObjectiveTemplatesClient(client *rest.Client) *SLOObjectiveTemplatesClient {
	c := &SLOObjectiveTemplatesClient{
		client: client,
	}
	return c
}

func (c *SLOObjectiveTemplatesClient) List(ctx context.Context, queryParams map[string][]string) (*http.Response, error) {
	resp, err := c.client.GET(ctx, templatesResourcePath, rest.RequestOptions{QueryParams: queryParams, ContentType: "application/json"})
	if err != nil {
		return nil, fmt.Errorf("unable to get templates: %w", err)
	}
	return resp, err
}
