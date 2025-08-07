package permissions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dynatrace/dynatrace-configuration-as-code-core/api"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
)

const (
	settingsEndpoint = "platform/classic/environment-api/v2/settings"
	schemasEndpoint  = settingsEndpoint + "/schemas"
	objectsEndpoint  = settingsEndpoint + "/objects"
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

type SettingsClient interface {
	GetSchemaIDsWithOwnerBasedAccessControl(ctx context.Context) ([]string, error)
	ListObjectsIDsOfSchema(ctx context.Context, schemaID string) ([]string, error)
}

type settingsClient struct {
	client *rest.Client
}

// NewSettingsClient creates a new settings client for fetching schema IDs that support owner-based access control as well as listing object IDs of a given schema.
func NewSettingsClient(client *rest.Client) *settingsClient {
	return &settingsClient{client: client}
}

// GetSchemaIDsWithOwnerBasedAccessControl retrieves schema IDs that support owner-based access control.
func (c *settingsClient) GetSchemaIDsWithOwnerBasedAccessControl(ctx context.Context) ([]string, error) {
	response, err := c.getSchemas(ctx, "schemaId,ownerBasedAccessControl")
	if err != nil {
		return nil, fmt.Errorf("failed to get schemas: %w", err)
	}

	var schemata schemaList
	if json.Unmarshal(response.Data, &schemata) != nil {
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

func (c *settingsClient) getSchemas(ctx context.Context, fields string) (api.Response, error) {
	response, err := c.client.GET(ctx, schemasEndpoint, rest.RequestOptions{QueryParams: url.Values{"fields": []string{fields}}})
	if err != nil {
		return api.Response{}, err
	}

	return api.NewResponseFromHTTPResponse(response)
}

// ListObjectsIDsOfSchema retrieves the IDs of objects for a given schema ID.
// It tries first with admin access and then without, if the first request fails with a forbidden error.
func (c *settingsClient) ListObjectsIDsOfSchema(ctx context.Context, schemaID string) ([]string, error) {
	response, err := c.getObjectIDsFirstPageWithAdminRetry(ctx, schemaID)
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

	for {
		if ol.NextPageKey == "" {
			break
		}

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

const fields = "objectId"
const pageSize = 100

func (c *settingsClient) getObjectIDsFirstPageWithAdminRetry(ctx context.Context, schemaID string) (api.Response, error) {
	response, err := c.getObjectsFirstPage(ctx, schemaID, fields, pageSize, true)
	if err == nil {
		return response, nil
	}

	var apiErr api.APIError
	if !errors.As(err, &apiErr) || (apiErr.StatusCode != http.StatusForbidden) {
		return api.Response{}, err
	}
	return c.getObjectsFirstPage(ctx, schemaID, fields, pageSize, false)
}

func (c *settingsClient) getObjectsFirstPage(ctx context.Context, schemaID string, fields string, pageSize int, adminAccess bool) (api.Response, error) {
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

func (c *settingsClient) getObjectsNextPage(ctx context.Context, nextPageKey string) (api.Response, error) {
	queryParams := url.Values{
		"nextPageKey": []string{nextPageKey},
	}
	response, err := c.client.GET(ctx, objectsEndpoint, rest.RequestOptions{QueryParams: queryParams})
	if err != nil {
		return api.Response{}, err
	}

	return api.NewResponseFromHTTPResponse(response)
}
