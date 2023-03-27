package policies

import (
	"fmt"
	"strings"

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

func (me *PolicyServiceClient) SchemaID() string {
	return "accounts:policies"
}

type PolicyCreateResponse struct {
	UUID string `json:"uuid"`
}

func (me *PolicyServiceClient) Create(v *policies.Policy) (*settings.Stub, error) {
	var err error
	levelType, levelID := getLevel(v)

	var pcr PolicyCreateResponse
	if err = me.client.Post(fmt.Sprintf("/iam/repo/%s/%s/policies", levelType, levelID), v, 201).Finish(&pcr); err != nil {
		return nil, err
	}
	return &settings.Stub{ID: joinID(pcr.UUID, v), Name: v.Name}, nil
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

func (me *PolicyServiceClient) List() (settings.Stubs, error) {
	return settings.Stubs{}, nil
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
