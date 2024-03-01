package groups

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	groups "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/cluster/v1/groups/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

const SchemaID = "accounts:groups"

func Service(credentials *settings.Credentials) settings.CRUDService[*groups.GroupConfig] {
	return &service{
		serviceClient: NewService(fmt.Sprintf("%s%s", credentials.Cluster.URL, "/api/v1.0/onpremise"), credentials.Cluster.Token),
	}
}

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

type service struct {
	serviceClient *ServiceClient
}

func (me *service) Create(v *groups.GroupConfig) (*api.Stub, error) {
	return me.serviceClient.Create(v)
}

func (me *service) Update(id string, v *groups.GroupConfig) error {
	return me.serviceClient.Update(v)
}

func (me *service) Delete(id string) error {
	return me.serviceClient.Delete(id)
}

func (me *service) List() (api.Stubs, error) {
	return me.serviceClient.List()
}

func (me *service) Get(id string, v *groups.GroupConfig) error {
	return me.serviceClient.Get(id, v)
}

func (me *service) Name() string {
	return me.SchemaID()
}

func (me *service) SchemaID() string {
	return SchemaID
}

func (cs *ServiceClient) SchemaID() string {
	return SchemaID
}

// Create TODO: documentation
func (cs *ServiceClient) Create(groupConfig *groups.GroupConfig) (*api.Stub, error) {
	var err error

	var createdGroupConfig groups.GroupConfig
	if err = cs.client.Post("/groups", groupConfig, 200).Finish(&createdGroupConfig); err != nil {
		return nil, err
	}
	return &api.Stub{ID: *createdGroupConfig.ID, Name: createdGroupConfig.Name}, nil
}

// Update TODO: documentation
func (cs *ServiceClient) Update(groupConfig *groups.GroupConfig) error {
	return cs.client.Put("/groups", groupConfig, 200).Finish()
}

// Delete TODO: documentation
func (cs *ServiceClient) Delete(id string) error {
	return cs.client.Delete(fmt.Sprintf("/groups/%s", id), 200).Finish()
}

// Get TODO: documentation
func (cs *ServiceClient) Get(id string, v *groups.GroupConfig) error {
	var err error

	if err = cs.client.Get(fmt.Sprintf("/groups/%s", id), 200).Finish(&v); err != nil {
		return err
	}
	return nil
}

// ListAll TODO: documentation
func (cs *ServiceClient) List() (api.Stubs, error) {
	var err error
	var stubs api.Stubs
	var groups []*groups.GroupConfig
	if err = cs.client.Get("/groups", 200).Finish(&groups); err != nil {
		return nil, err
	}
	for _, group := range groups {
		stubs = append(stubs, &api.Stub{ID: *group.ID, Name: group.Name})
	}
	return stubs, nil
}
