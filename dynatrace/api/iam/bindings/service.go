package bindings

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam"
	bindings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/bindings/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/policies"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

type BindingServiceClient struct {
	clientID     string
	accountID    string
	clientSecret string
	tokenURL     string
	endpointURL  string
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

func (me *BindingServiceClient) TokenURL() string {
	return me.tokenURL
}

func (me *BindingServiceClient) EndpointURL() string {
	return me.endpointURL
}

func NewPolicyService(clientID string, accountID string, clientSecret string, tokenURL string, endpointURL string) *BindingServiceClient {
	return &BindingServiceClient{clientID: clientID, accountID: accountID, clientSecret: clientSecret, tokenURL: tokenURL, endpointURL: endpointURL}
}

func Service(credentials *settings.Credentials) settings.CRUDService[*bindings.PolicyBinding] {
	return &BindingServiceClient{clientID: credentials.IAM.ClientID, accountID: credentials.IAM.AccountID, clientSecret: credentials.IAM.ClientSecret, tokenURL: credentials.IAM.TokenURL, endpointURL: credentials.IAM.EndpointURL}
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

func (me *BindingServiceClient) Create(ctx context.Context, v *bindings.PolicyBinding) (*api.Stub, error) {
	id := joinID(v)
	var err error
	if err = me.Update(ctx, id, v); err != nil {
		return nil, err
	}
	return &api.Stub{ID: id, Name: "PolicyBindings-" + id}, nil
}

func (me *BindingServiceClient) Get(ctx context.Context, id string, v *bindings.PolicyBinding) error {
	groupID, levelType, levelID, err := splitID(id)
	if err != nil {
		return err
	}
	var responseBytes []byte

	client := iam.NewIAMClient(me)

	if responseBytes, err = client.GET(fmt.Sprintf("%s/iam/v1/repo/%s/%s/bindings/groups/%s", me.endpointURL, levelType, levelID, groupID), 200, false); err != nil {
		return err
	}
	if err = json.Unmarshal(responseBytes, &v); err != nil {
		return err
	}
	if levelType == "account" {
		v.Account = levelID
	} else if levelType == "environment" {
		v.Environment = levelID
	}
	v.GroupID = groupID
	for idx, policyID := range v.PolicyIDs {
		v.PolicyIDs[idx] = fmt.Sprintf("%s#-#%s#-#%s", policyID, levelType, levelID)
	}
	return nil
}

func (me *BindingServiceClient) Update(ctx context.Context, id string, bindings *bindings.PolicyBinding) error {
	groupID, levelType, levelID, err := splitID(id)
	if err != nil {
		return err
	}

	client := iam.NewIAMClient(me)

	policyIDs := []string{}
	for _, policyID := range bindings.PolicyIDs {
		uuid, policyLevelType, policyLevelID, err := policies.SplitID(policyID, levelType, levelID)
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

	if _, err = client.PUT(fmt.Sprintf("%s/iam/v1/repo/%s/%s/bindings/groups/%s", me.endpointURL, levelType, levelID, groupID), bindings, 204, false); err != nil {
		return err
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

func (me *BindingServiceClient) List(ctx context.Context) (api.Stubs, error) {
	var err error
	var responseBytes []byte
	client := iam.NewIAMClient(me)

	if responseBytes, err = client.GET(fmt.Sprintf("%s/env/v2/accounts/%s/environments", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:")), 200, false); err != nil {
		return nil, err
	}

	var envResponse ListEnvResponse
	if err = json.Unmarshal(responseBytes, &envResponse); err != nil {
		return nil, err
	}

	if responseBytes, err = client.GET(fmt.Sprintf("%s/iam/v1/repo/account/%s/bindings", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:")), 200, false); err != nil {
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
				id := fmt.Sprintf("%s#-#%s#-#%s", group, "account", me.AccountID())
				stubs = append(stubs, &api.Stub{ID: id, Name: "PolicyBindings-" + id})
				groupIds[group] = true
			}
		}
	}

	for _, environment := range envResponse.Data {
		if responseBytes, err = client.GET(fmt.Sprintf("%s/iam/v1/repo/environment/%s/bindings", me.endpointURL, environment.ID), 200, false); err != nil {
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
					stubs = append(stubs, &api.Stub{ID: id, Name: "PolicyBindings-" + id})
					groupIds[group] = true
				}
			}
		}
	}
	return stubs, nil
}

func (me *BindingServiceClient) Delete(ctx context.Context, id string) error {
	groupID, levelType, levelID, err := splitID(id)
	if err != nil {
		return err
	}
	var binding bindings.PolicyBinding
	if err = me.Get(ctx, id, &binding); err != nil {
		return err
	}
	for _, policyID := range binding.PolicyIDs {
		if _, err = iam.NewIAMClient(me).DELETE_MULTI_RESPONSE(fmt.Sprintf("%s/iam/v1/repo/%s/%s/bindings/%s/%s", me.endpointURL, levelType, levelID, policyID, groupID), []int{204, 400}, false); err != nil {
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
