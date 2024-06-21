package users

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	users "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/cluster/v1/users/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

const SchemaID = "accounts:users"

func Service(credentials *settings.Credentials) settings.CRUDService[*users.UserConfig] {
	return &service{
		serviceClient: NewService(fmt.Sprintf("%s%s", credentials.Cluster.URL, "/api/v1.0/onpremise"), credentials.Cluster.Token),
	}
}

// ServiceClient TODO: documentation
type ServiceClient struct {
	client rest.Client
}

func (me *service) Create(ctx context.Context, v *users.UserConfig) (*api.Stub, error) {
	return me.serviceClient.Create(ctx, v)
}

func (me *service) Update(ctx context.Context, id string, v *users.UserConfig) error {
	return me.serviceClient.Update(ctx, v)
}

func (me *service) Delete(ctx context.Context, id string) error {
	return me.serviceClient.Delete(ctx, id)
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return me.serviceClient.List(ctx)
}

func (me *service) Get(ctx context.Context, id string, v *users.UserConfig) error {
	return me.serviceClient.Get(ctx, id, v)
}

func (me *service) SchemaID() string {
	return SchemaID
}

func (cs *ServiceClient) SchemaID() string {
	return SchemaID
}

// NewService creates a new Service Client
// baseURL should look like this: "https://siz65484.live.dynatrace.com/api/config/v1"
// token is an API Token
func NewService(baseURL string, token string) *ServiceClient {
	return &ServiceClient{client: rest.DefaultClient(baseURL, token)}
}

type service struct {
	serviceClient *ServiceClient
}

// Create TODO: documentation
func (cs *ServiceClient) Create(ctx context.Context, userConfig *users.UserConfig) (*api.Stub, error) {
	var err error

	var createdUserConfig users.UserConfig
	if err = cs.client.Post("/users", userConfig, 200).Finish(&createdUserConfig); err != nil {
		return nil, err
	}
	return &api.Stub{ID: createdUserConfig.UserName, Name: createdUserConfig.UserName}, nil
}

// Update TODO: documentation
func (cs *ServiceClient) Update(ctx context.Context, userConfig *users.UserConfig) error {
	return cs.client.Put("/users", userConfig, 200).Finish()
}

// Delete TODO: documentation
func (cs *ServiceClient) Delete(ctx context.Context, id string) error {
	if len(id) == 0 {
		return errors.New("empty ID provided for the Dashboard to delete")
	}
	return cs.client.Delete(fmt.Sprintf("/users/%s", id), 200).Finish()
}

// Get TODO: documentation
func (cs *ServiceClient) Get(ctx context.Context, id string, v *users.UserConfig) error {
	if len(id) == 0 {
		return errors.New("empty ID provided for the Dashboard to fetch")
	}

	var err error
	if err = cs.client.Get(fmt.Sprintf("/users/%s", id), 200).Finish(&v); err != nil {
		if strings.HasPrefix(err.Error(), "Not Found (GET) ") {
			return rest.Error{Code: 404, Message: fmt.Sprintf("user '%s' doesn't exist", id)}
		}
		return err
	}
	return nil
}

// ListAll TODO: documentation
func (cs *ServiceClient) List(ctx context.Context) (api.Stubs, error) {
	var err error
	var stubs api.Stubs
	var users []*users.UserConfig
	if err = cs.client.Get("/users", 200).Finish(&users); err != nil {
		return nil, err
	}
	for _, user := range users {
		stubs = append(stubs, &api.Stub{ID: user.UserName, Name: user.UserName})
	}
	return stubs, nil
}
