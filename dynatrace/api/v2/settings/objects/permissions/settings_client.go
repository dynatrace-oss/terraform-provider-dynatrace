/*
 * @license
 * Copyright 2025 Dynatrace LLC
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

package permissions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	trest "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/api"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
)

const (
	settingsEndpoint = "platform/classic/environment-api/v2/settings"
	schemasEndpoint  = settingsEndpoint + "/schemas"
	objectsEndpoint  = settingsEndpoint + "/objects"

	fieldsObjectID = "objectId"
	pageSize       = 100
)

type schemaList struct {
	Items []schemaStub `json:"items"`
}

type schemaStub struct {
	SchemaID                string `json:"schemaId"`
	OwnerBasedAccessControl bool   `json:"ownerBasedAccessControl"`
}

type objectsList struct {
	Items       []objectStub `json:"items"`
	NextPageKey string       `json:"nextPageKey"`
}

type objectStub struct {
	ObjectID string `json:"objectId"`
}

type PlatformSettingsClient interface {
	GetSchemaIDsWithOwnerBasedAccessControl(ctx context.Context) ([]string, error)
	ListObjectsIDsOfSchema(ctx context.Context, schemaID string) ([]string, error)
}

type platformSettingsClient struct {
	client *rest.Client
}

// NewPlatformSettingsClient creates a new settings client for fetching schema IDs that support owner-based access control as well as listing object IDs of a given schema.
func NewPlatformSettingsClient(client *rest.Client) *platformSettingsClient {
	return &platformSettingsClient{client: client}
}

// GetSchemaIDsWithOwnerBasedAccessControl retrieves schema IDs that support owner-based access control.
func (c *platformSettingsClient) GetSchemaIDsWithOwnerBasedAccessControl(ctx context.Context) ([]string, error) {
	response, err := c.getSchemas(ctx, "schemaId,ownerBasedAccessControl")
	if err != nil {
		return nil, fmt.Errorf("failed to get schemas: %w", err)
	}

	var schemata schemaList
	if err := json.Unmarshal(response.Data, &schemata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schema response: %w", err)
	}

	schemaIDs := []string{}
	for _, schemaStub := range schemata.Items {
		if schemaStub.OwnerBasedAccessControl {
			schemaIDs = append(schemaIDs, schemaStub.SchemaID)
		}
	}
	return schemaIDs, nil
}

func (c *platformSettingsClient) getSchemas(ctx context.Context, fields string) (api.Response, error) {
	response, err := c.client.GET(ctx, schemasEndpoint, rest.RequestOptions{QueryParams: url.Values{"fields": []string{fields}}})
	if err != nil {
		return api.Response{}, err
	}

	return api.NewResponseFromHTTPResponse(response)
}

// ListObjectsIDsOfSchema retrieves the IDs of objects for a given schema ID.
// It tries first with admin access and then without, if the first request fails with a forbidden error.
func (c *platformSettingsClient) ListObjectsIDsOfSchema(ctx context.Context, schemaID string) ([]string, error) {
	response, err, _ := trest.DoWithAdminAccessRetry(func(adminAccess bool) (api.Response, error) {
		return c.getObjectsFirstPage(ctx, schemaID, fieldsObjectID, pageSize, adminAccess)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get first page of objects: %w", err)
	}

	var ol objectsList
	if err := json.Unmarshal(response.Data, &ol); err != nil {
		return nil, fmt.Errorf("failed to unmarshal objects response: %w", err)
	}

	objectIDs := []string{}
	for _, objectStub := range ol.Items {
		objectIDs = append(objectIDs, objectStub.ObjectID)
	}

	for ol.NextPageKey != "" {
		response, err := c.getObjectsNextPage(ctx, ol.NextPageKey)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page of objects: %w", err)
		}

		ol = objectsList{}
		if err := json.Unmarshal(response.Data, &ol); err != nil {
			return nil, fmt.Errorf("failed to unmarshal objects response: %w", err)
		}

		for _, objectStub := range ol.Items {
			objectIDs = append(objectIDs, objectStub.ObjectID)
		}
	}

	return objectIDs, nil
}

func (c *platformSettingsClient) getObjectsFirstPage(ctx context.Context, schemaID string, fields string, pageSize int, adminAccess bool) (api.Response, error) {
	queryParams := url.Values{
		"schemaIds":   []string{schemaID},
		"fields":      []string{fields},
		"pageSize":    []string{strconv.Itoa(pageSize)},
		"adminAccess": []string{strconv.FormatBool(adminAccess)},
	}
	response, err := c.client.GET(ctx, objectsEndpoint, rest.RequestOptions{QueryParams: queryParams})
	if err != nil {
		return api.Response{}, err
	}

	return api.NewResponseFromHTTPResponse(response)
}

func (c *platformSettingsClient) getObjectsNextPage(ctx context.Context, nextPageKey string) (api.Response, error) {
	queryParams := url.Values{
		"nextPageKey": []string{nextPageKey},
	}
	response, err := c.client.GET(ctx, objectsEndpoint, rest.RequestOptions{QueryParams: queryParams})
	if err != nil {
		return api.Response{}, err
	}

	return api.NewResponseFromHTTPResponse(response)
}
