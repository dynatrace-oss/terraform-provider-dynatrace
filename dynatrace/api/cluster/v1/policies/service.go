package policies

import (
	"fmt"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	policies "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/cluster/v1/policies/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

type PolicyServiceClient struct {
	client rest.Client
}

func NewPolicyService(baseURL string, apiToken string) *PolicyServiceClient {
	return &PolicyServiceClient{client: rest.DefaultClient(baseURL, apiToken)}
}

func Service(credentials *settings.Credentials) settings.CRUDService[*policies.Policy] {
	return &service{
		serviceClient: NewPolicyService(fmt.Sprintf("%s%s", credentials.Cluster.URL, "/api/cluster/v2"), credentials.Cluster.Token),
	}
}

type service struct {
	serviceClient *PolicyServiceClient
}

func (me *service) List() (api.Stubs, error) {
	return me.serviceClient.List()
}

func (me *service) Get(id string, v *policies.Policy) error {
	return me.serviceClient.Get(id, v)
}

func (me *service) SchemaID() string {
	return "accounts:policies"
}

func (me *service) Create(v *policies.Policy) (*api.Stub, error) {
	return me.serviceClient.Create(v)
}

func (me *service) Update(id string, v *policies.Policy) error {
	return me.serviceClient.Update(id, v)
}

func (me *service) Delete(id string) error {
	return me.serviceClient.Delete(id)
}

func (me *PolicyServiceClient) SchemaID() string {
	return "accounts:policies"
}

type PolicyCreateResponse struct {
	UUID string `json:"uuid"`
}

func (me *PolicyServiceClient) Create(v *policies.Policy) (*api.Stub, error) {
	var err error
	levelType, levelID := getLevel(v)

	var pcr PolicyCreateResponse
	if err = me.client.Post(fmt.Sprintf("/iam/repo/%s/%s/policies", levelType, levelID), v, 201).Finish(&pcr); err != nil {
		return nil, err
	}
	return &api.Stub{ID: joinID(pcr.UUID, v), Name: v.Name}, nil
}

func (me *PolicyServiceClient) Get(id string, v *policies.Policy) error {
	uuid, levelType, levelID, err := SplitID(id)
	if err != nil {
		return err
	}

	if err = me.client.Get(fmt.Sprintf("/iam/repo/%s/%s/policies/%s", levelType, levelID, uuid), 200).Finish(&v); err != nil {
		return err
	}
	if levelType == "cluster" {
		v.Cluster = levelID
	} else if levelType == "environment" {
		v.Environment = levelID
	}
	return nil
}

func (me *PolicyServiceClient) Update(id string, user *policies.Policy) error {
	uuid, levelType, levelID, err := SplitID(id)
	if err != nil {
		return err
	}

	if err = me.client.Put(fmt.Sprintf("/iam/repo/%s/%s/policies/%s", levelType, levelID, uuid), user, 204).Finish(); err != nil {
		return err
	}
	return nil
}

func (me *PolicyServiceClient) List() (api.Stubs, error) {
	var err error
	var stubs api.Stubs

	clusterInfoResponse := struct {
		ClusterUUID string `json:"clusterUuid"`
	}{}

	if err = me.client.Get("/license/consumption/hour", 200).Finish(&clusterInfoResponse); err != nil {
		return stubs, err
	}
	policiesResponse := struct {
		Policies []struct {
			UUID string `json:"uuid"`
			Name string `json:"name"`
		} `json:"policies"`
	}{}
	if err = me.client.Get(fmt.Sprintf("/iam/repo/cluster/%s/policies", clusterInfoResponse.ClusterUUID), 200).Finish(&policiesResponse); err != nil {
		return stubs, err
	}
	for _, policy := range policiesResponse.Policies {
		stubs = append(stubs, &api.Stub{ID: fmt.Sprintf("%s#-#%s#-#%s", policy.UUID, "cluster", clusterInfoResponse.ClusterUUID), Name: policy.Name})
	}

	environmentsResponse := struct {
		Environments []struct {
			ID string `json:"id"`
		} `json:"environments"`
	}{}
	if err = me.client.Get("/environments?pageSize=1000", 200).Finish(&environmentsResponse); err != nil {
		return stubs, err
	}
	for _, environment := range environmentsResponse.Environments {
		policiesResponse := struct {
			Policies []struct {
				UUID string `json:"uuid"`
				Name string `json:"name"`
			} `json:"policies"`
		}{}
		if err = me.client.Get(fmt.Sprintf("/iam/repo/environment/%s/policies", environment.ID), 200).Finish(&policiesResponse); err != nil {
			return stubs, err
		}
		for _, policy := range policiesResponse.Policies {
			stubs = append(stubs, &api.Stub{ID: fmt.Sprintf("%s#-#%s#-#%s", policy.UUID, "environment", environment.ID), Name: policy.Name})
		}
	}

	return stubs, nil
}

func (me *PolicyServiceClient) Delete(id string) error {
	uuid, levelType, levelID, err := SplitID(id)
	if err != nil {
		return err
	}
	return me.client.Delete(fmt.Sprintf("/iam/repo/%s/%s/policies/%s", levelType, levelID, uuid), 204).Finish()
}

func SplitID(id string) (uuid string, levelType string, levelID string, err error) {
	parts := strings.Split(id, "#-#")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("%s is not a valid ID for a policy", id)
	}
	return parts[0], parts[1], parts[2], nil
}

func joinID(uuid string, policy *policies.Policy) string {
	levelType, levelID := getLevel(policy)
	return fmt.Sprintf("%s#-#%s#-#%s", uuid, levelType, levelID)
}

func getLevel(policy *policies.Policy) (string, string) {
	if len(policy.Cluster) > 0 {
		return "cluster", policy.Cluster
	}
	if len(policy.Environment) > 0 {
		return "environment", policy.Environment
	}
	return "global", "global"
}
