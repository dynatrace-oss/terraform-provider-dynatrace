package groups

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
)

// ServiceClient TODO: documentation
type ServiceClient struct {
	client rest.Client
}

// NewService creates a new Service Client
// baseURL should look like this: "https://#######.live.dynatrace.com/api/config/v1"
// token is an API Token
func NewService(baseURL string, token string) *ServiceClient {
	return &ServiceClient{client: rest.DefaultClient(baseURL, token)}
}

// Create TODO: documentation
func (cs *ServiceClient) Create(groupConfig *GroupConfig) (*GroupConfig, error) {
	var err error

	var createdGroupConfig GroupConfig
	if err = cs.client.Post("/groups", groupConfig, 200).Finish(&createdGroupConfig); err != nil {
		return nil, err
	}
	return &createdGroupConfig, nil
}

// Update TODO: documentation
func (cs *ServiceClient) Update(groupConfig *GroupConfig) error {
	return cs.client.Put("/groups", groupConfig, 200).Finish()
}

// Delete TODO: documentation
func (cs *ServiceClient) Delete(id string) error {
	return cs.client.Delete(fmt.Sprintf("/groups/%s", id), 200).Finish()
}

// Get TODO: documentation
func (cs *ServiceClient) Get(id string) (*GroupConfig, error) {
	var err error

	var groupConfig GroupConfig
	if err = cs.client.Get(fmt.Sprintf("/groups/%s", id), 200).Finish(&groupConfig); err != nil {
		return nil, err
	}
	return &groupConfig, nil
}

// ListAll TODO: documentation
func (cs *ServiceClient) ListAll() ([]*GroupConfig, error) {
	var err error
	var groups []*GroupConfig
	if err = cs.client.Get("/groups", 200).Finish(&groups); err != nil {
		return nil, err
	}
	return groups, nil
}
