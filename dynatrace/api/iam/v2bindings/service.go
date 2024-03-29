package v2bindings

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/policies"
	bindings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/v2bindings/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

type BindingServiceClient struct {
	clientID     string
	accountID    string
	clientSecret string
}

func (me *BindingServiceClient) ClientID() string {
	return me.clientID
}

func (me *BindingServiceClient) AccountID() string {
	return me.accountID
}

func (me *BindingServiceClient) ClientSecret() string {
	return me.clientSecret
}

func NewPolicyService(clientID string, accountID string, clientSecret string) *BindingServiceClient {
	return &BindingServiceClient{clientID: clientID, accountID: accountID, clientSecret: clientSecret}
}

func Service(credentials *settings.Credentials) settings.CRUDService[*bindings.PolicyBinding] {
	return &BindingServiceClient{clientID: credentials.IAM.ClientID, accountID: credentials.IAM.AccountID, clientSecret: credentials.IAM.ClientSecret}
}

func (me *BindingServiceClient) SchemaID() string {
	return "accounts:iam:bindings"
}

func (me *BindingServiceClient) Name() string {
	return me.SchemaID()
}

type PolicyCreateResponse struct {
	UUID string `json:"uuid"`
}

func (me *BindingServiceClient) Create(v *bindings.PolicyBinding) (*api.Stub, error) {
	id := joinID(v)
	var err error
	if err = me.Update(id, v); err != nil {
		return nil, err
	}
	return &api.Stub{ID: id, Name: "PolicyV2Bindings-" + id}, nil
}

func (me *BindingServiceClient) Get(id string, v *bindings.PolicyBinding) error {
	groupID, levelType, levelID, err := splitID(id)
	if err != nil {
		return err
	}
	var responseBytes []byte

	client := iam.NewIAMClient(me)

	if responseBytes, err = client.GET(fmt.Sprintf("https://api.dynatrace.com/iam/v1/repo/%s/%s/bindings/groups/%s", levelType, levelID, groupID), 200, false); err != nil {
		return err
	}
	policyUUIDStruct := struct {
		PolicyUuids []string `json:"policyUuids"`
	}{}
	if err = json.Unmarshal(responseBytes, &policyUUIDStruct); err != nil {
		return err
	}
	if levelType == "account" {
		v.Account = levelID
	} else if levelType == "environment" {
		v.Environment = levelID
	}
	v.GroupID = groupID
	policies := []*bindings.Policy{}
	for _, policyID := range policyUUIDStruct.PolicyUuids {
		if responseBytes, err = client.GET(fmt.Sprintf("https://api.dynatrace.com/iam/v1/repo/%s/%s/bindings/%s/%s", levelType, levelID, policyID, groupID), 200, false); err != nil {
			return err
		}
		bindingsResponse := struct {
			LevelType      string `json:"levelType"`
			LevelID        string `json:"levelId"`
			PolicyBindings []struct {
				PolicyUUID string            `json:"policyUuid"`
				GroupUUIDs []string          `json:"groups"`
				Parameters map[string]string `json:"parameters"`
				Metadata   map[string]string `json:"metadata"`
			} `json:"policyBindings"`
		}{}
		if err = json.Unmarshal(responseBytes, &bindingsResponse); err != nil {
			return err
		}
		if len(bindingsResponse.PolicyBindings) == 0 {
			continue
		}
		for _, policyBinding := range bindingsResponse.PolicyBindings {
			policy := &bindings.Policy{ID: fmt.Sprintf("%s#-#%s#-#%s", policyID, levelType, levelID)}
			policies = append(policies, policy)
			if len(policyBinding.Parameters) > 0 {
				policy.Parameters = map[string]string{}
				for key, value := range policyBinding.Parameters {
					policy.Parameters[key] = value
				}
			}
			if len(policyBinding.Metadata) > 0 {
				policy.Metadata = map[string]string{}
				for key, value := range policyBinding.Metadata {
					policy.Metadata[key] = value
				}
			}
		}
	}
	v.Policies = policies
	return nil
}

func (me *BindingServiceClient) Update(id string, v *bindings.PolicyBinding) error {
	groupID, levelType, levelID, err := splitID(id)
	if err != nil {
		return err
	}

	client := iam.NewIAMClient(me)

	policiesList := []*bindings.Policy{}
	for _, policy := range v.Policies {
		policyID := policy.ID
		policyUUID, policyLevelType, policyLevelID, err := policies.SplitID(policyID)
		if err != nil {
			return err
		}
		if policyLevelID != levelID || policyLevelType != levelType {
			return fmt.Errorf("The policy %s is defined for %s = %s. It cannot be used within the scope %s = %s", policyUUID, policyLevelType, policyLevelID, levelType, levelID)
		}
		policiesList = append(policiesList, policy)
	}

	for _, policy := range policiesList {
		policyUUID, _, _, _ := policies.SplitID(policy.ID)
		if _, err = client.DELETE_MULTI_RESPONSE(fmt.Sprintf("https://api.dynatrace.com/iam/v1/repo/%s/%s/bindings/%s/%s", levelType, levelID, policyUUID, groupID), []int{204, 400, 404}, false); err != nil {
			return err
		}
	}
	for _, policy := range policiesList {
		policyUUID, _, _, _ := policies.SplitID(policy.ID)
		payload := struct {
			Parameters map[string]string `json:"parameters"`
			Metadata   map[string]string `json:"metadata"`
		}{
			Parameters: policy.Parameters,
			Metadata:   policy.Metadata,
		}
		if _, err = client.POST(fmt.Sprintf("https://api.dynatrace.com/iam/v1/repo/%s/%s/bindings/%s/%s", levelType, levelID, policyUUID, groupID), payload, 204, false); err != nil {
			return err
		}
	}

	return nil
}

type DataStub struct {
	ID string `json:"id"`
}

type ListEnvResponse struct {
	Data []DataStub `json:"data"`
}

type PolicyBindingStub struct {
	Groups []string `json:"groups"`
}

type ListPolicyBindingsResponse struct {
	PolicyBindings []PolicyBindingStub `json:"policyBindings"`
}

func (me *BindingServiceClient) List() (api.Stubs, error) {
	var err error
	var responseBytes []byte
	client := iam.NewIAMClient(me)

	if responseBytes, err = client.GET(fmt.Sprintf("https://api.dynatrace.com/env/v2/accounts/%s/environments", strings.TrimPrefix(me.AccountID(), "urn:dtaccount:")), 200, false); err != nil {
		return nil, err
	}

	var envResponse ListEnvResponse
	if err = json.Unmarshal(responseBytes, &envResponse); err != nil {
		return nil, err
	}

	if responseBytes, err = client.GET(fmt.Sprintf("https://api.dynatrace.com/iam/v1/repo/account/%s/bindings", strings.TrimPrefix(me.AccountID(), "urn:dtaccount:")), 200, false); err != nil {
		return nil, err
	}

	var response ListPolicyBindingsResponse
	if err = json.Unmarshal(responseBytes, &response); err != nil {
		return nil, err
	}

	var stubs api.Stubs
	groupIds := map[string]bool{}
	for _, policy := range response.PolicyBindings {
		for _, group := range policy.Groups {
			if _, exists := groupIds[group]; !exists {
				id := fmt.Sprintf("%s#-#%s#-#%s", group, "account", strings.TrimPrefix(me.AccountID(), "urn:dtaccount:"))
				stubs = append(stubs, &api.Stub{ID: id, Name: "PolicyV2Bindings-" + id})
				groupIds[group] = true
			}
		}
	}

	for _, environment := range envResponse.Data {
		if responseBytes, err = client.GET(fmt.Sprintf("https://api.dynatrace.com/iam/v1/repo/environment/%s/bindings", environment.ID), 200, false); err != nil {
			return nil, err
		}

		var response ListPolicyBindingsResponse
		if err = json.Unmarshal(responseBytes, &response); err != nil {
			return nil, err
		}

		groupIds := map[string]bool{}
		for _, policy := range response.PolicyBindings {
			for _, group := range policy.Groups {
				if _, exists := groupIds[group]; !exists {
					id := fmt.Sprintf("%s#-#%s#-#%s", group, "environment", environment.ID)
					stubs = append(stubs, &api.Stub{ID: id, Name: "PolicyV2Bindings-" + id})
					groupIds[group] = true
				}
			}
		}
	}
	return stubs, nil
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
	policyUUIDs := map[string]string{}
	for _, policy := range binding.Policies {
		policyUUID, _, _, err := policies.SplitID(policy.ID)
		if err != nil {
			return err
		}
		policyUUIDs[policyUUID] = policyUUID
	}
	for policyUUID, _ := range policyUUIDs {
		if _, err = iam.NewIAMClient(me).DELETE(fmt.Sprintf("https://api.dynatrace.com/iam/v1/repo/%s/%s/bindings/%s/%s", levelType, levelID, policyUUID, groupID), 204, false); err != nil {
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
	if len(binding.Account) > 0 {
		return "account", binding.Account
	}
	if len(binding.Environment) > 0 {
		return "environment", binding.Environment
	}
	return "global", "global"
}
