package users

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
)

// ServiceClient TODO: documentation
type ServiceClient struct {
	client rest.Client
}

// NewService creates a new Service Client
// baseURL should look like this: "https://siz65484.live.dynatrace.com/api/config/v1"
// token is an API Token
func NewService(baseURL string, token string) *ServiceClient {
	return &ServiceClient{client: rest.DefaultClient(baseURL, token)}
}

// Create TODO: documentation
func (cs *ServiceClient) Create(userConfig *UserConfig) (*UserConfig, error) {
	var err error

	var createdUserConfig UserConfig
	if err = cs.client.Post("/users", userConfig, 200).Finish(&createdUserConfig); err != nil {
		return nil, err
	}
	return &createdUserConfig, nil
}

// Update TODO: documentation
func (cs *ServiceClient) Update(userConfig *UserConfig) error {
	return cs.client.Put("/users", userConfig, 200).Finish()
}

// Delete TODO: documentation
func (cs *ServiceClient) Delete(id string) error {
	if len(id) == 0 {
		return errors.New("empty ID provided for the Dashboard to delete")
	}
	return cs.client.Delete(fmt.Sprintf("/users/%s", id), 200).Finish()
}

// Get TODO: documentation
func (cs *ServiceClient) Get(id string) (*UserConfig, error) {
	if len(id) == 0 {
		return nil, errors.New("empty ID provided for the Dashboard to fetch")
	}

	var err error
	var userConfig UserConfig
	if err = cs.client.Get(fmt.Sprintf("/users/%s", id), 200).Finish(&userConfig); err != nil {
		if strings.HasPrefix(err.Error(), "Not Found (GET) ") {
			return nil, rest.Error{Code: 404, Message: fmt.Sprintf("user '%s' doesn't exist", id)}
		}
		return nil, err
	}
	return &userConfig, nil
}

// ListAll TODO: documentation
func (cs *ServiceClient) ListAll() ([]*UserConfig, error) {
	var err error

	var users []*UserConfig
	if err = cs.client.Get("/users", 200).Finish(&users); err != nil {
		return nil, err
	}
	return users, nil
}
