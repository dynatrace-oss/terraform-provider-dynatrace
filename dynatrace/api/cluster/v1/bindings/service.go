package bindings

import (
	"fmt"
	"strings"

	bindings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/cluster/v1/bindings/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/policies"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

type BindingServiceClient struct {
	client rest.Client
}

func NewPolicyService(baseURL string, apiToken string) *BindingServiceClient {
	return &BindingServiceClient{client: rest.DefaultClient(baseURL, apiToken)}
}

func (me *BindingServiceClient) SchemaID() string {
	return "accounts:policy-bindings"
}

type PolicyCreateResponse struct {
	UUID string `json:"uuid"`
}

func (me *BindingServiceClient) Create(v *bindings.PolicyBinding) (*settings.Stub, error) {
	id := joinID(v)
	var err error
	if err = me.Update(id, v); err != nil {
		return nil, err
	}
	return &settings.Stub{ID: id, Name: "PolicyBindings-" + id}, nil
}

func (me *BindingServiceClient) Get(id string, v *bindings.PolicyBinding) error {
	groupID, levelType, levelID, err := splitID(id)
	if err != nil {
		return err
	}

	if err = me.client.Get(fmt.Sprintf("/iam/repo/%s/%s/bindings/groups/%s", levelType, levelID, groupID), 200).Finish(&v); err != nil {
		return err
	}
	if levelType == "cluster" {
		v.Cluster = levelID
	} else if levelType == "environment" {
		v.Environment = levelID
	}
	v.GroupID = groupID
	for idx, policyID := range v.PolicyIDs {
		v.PolicyIDs[idx] = fmt.Sprintf("%s#-#%s#-#%s", policyID, levelType, levelID)
	}
	return nil
}

func (me *BindingServiceClient) Update(id string, bindings *bindings.PolicyBinding) error {
	groupID, levelType, levelID, err := splitID(id)
	if err != nil {
		return err
	}

	policyIDs := []string{}
	for _, policyID := range bindings.PolicyIDs {
		uuid, policyLevelType, policyLevelID, err := policies.SplitID(policyID)
		if policyLevelID != levelID || policyLevelType != levelType {
			return fmt.Errorf("The policy %s is defined for %s = %s. It cannot be used within the scope %s = %s", uuid, policyLevelType, policyLevelID, levelType, levelID)
		}
		if err == nil {
			policyIDs = append(policyIDs, uuid)
		} else {
			policyIDs = append(policyIDs, policyID)
		}
	}
	bindings.PolicyIDs = policyIDs

	if err = me.client.Put(fmt.Sprintf("/iam/repo/%s/%s/bindings/groups/%s", levelType, levelID, groupID), bindings, 204).Finish(); err != nil {
		return err
	}
	return nil
}

func (me *BindingServiceClient) List() (settings.Stubs, error) {
	return settings.Stubs{}, nil
}

func (me *BindingServiceClient) Delete(id string) error {
	groupID, levelType, levelID, err := splitID(id)
	if err != nil {
		return err
	}
	var binding bindings.PolicyBinding
	if err = me.Get(id, &binding); err != nil {
		return err
	}
	for _, policyID := range binding.PolicyIDs {
		if err = me.client.Delete(fmt.Sprintf("/iam/repo/%s/%s/bindings/%s/%s", levelType, levelID, policyID, groupID), 204).Finish(); err != nil {
			return err
		}
	}
	return nil
}

func splitID(id string) (groupID string, levelType string, levelID string, err error) {
	parts := strings.Split(id, "#-#")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("%s is not a valid ID for a policy", id)
	}
	return parts[0], parts[1], parts[2], nil
}

func joinID(binding *bindings.PolicyBinding) string {
	levelType, levelID := getLevel(binding)
	return fmt.Sprintf("%s#-#%s#-#%s", binding.GroupID, levelType, levelID)
}

func getLevel(binding *bindings.PolicyBinding) (string, string) {
	if len(binding.Cluster) > 0 {
		return "cluster", binding.Cluster
	}
	if len(binding.Environment) > 0 {
		return "environment", binding.Environment
	}
	return "global", "global"
}
