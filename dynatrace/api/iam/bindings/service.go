/**
* @license
* Copyright 2025 Dynatrace LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package bindings

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam"
	bindings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/bindings/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/policies"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	coreapi "github.com/dynatrace/dynatrace-configuration-as-code-core/api"
	rest2 "github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
)

type BindingServiceClient struct {
	credentials *rest.Credentials
}

func Service(credentials *rest.Credentials) settings.CRUDService[*bindings.PolicyBinding] {
	return &BindingServiceClient{credentials: credentials}
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
	if err := me.Update(ctx, id, v); err != nil {
		return nil, err
	}
	return &api.Stub{ID: id, Name: "PolicyBindings-" + id}, nil
}

func (me *BindingServiceClient) Get(ctx context.Context, id string, v *bindings.PolicyBinding) error {
	groupID, levelType, levelID, err := splitID(id)
	if err != nil {
		return err
	}
	client := iam.NewIAMClient(ctx, me.credentials)

	response, err := client.GET(ctx, fmt.Sprintf("/iam/v1/repo/%s/%s/bindings/groups/%s", levelType, levelID, groupID), rest2.RequestOptions{})
	if err != nil {
		return err
	}
	if err = json.Unmarshal(response.Data, &v); err != nil {
		return err
	}
	if levelType == "account" {
		v.Account = levelID
	} else if levelType == "environment" {
		v.Environment = levelID
	}
	v.GroupID = groupID

	for idx, policyID := range v.PolicyIDs {
		v.PolicyIDs[idx] = policies.Join(policyID, levelType, levelID)
	}
	return nil
}

func (me *BindingServiceClient) Update(ctx context.Context, id string, bindings *bindings.PolicyBinding) error {
	groupID, levelType, levelID, err := splitID(id)
	if err != nil {
		return err
	}

	client := iam.NewIAMClient(ctx, me.credentials)

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

	if _, err = client.PUT(ctx, fmt.Sprintf("/iam/v1/repo/%s/%s/bindings/groups/%s", levelType, levelID, groupID), bindings, rest2.RequestOptions{}); err != nil {
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

type ListPolicyUUIDsForGroupResponse struct {
	PolicyUUIDs []string `json:"policyUuids"`
}

func (me *BindingServiceClient) GetPolicyUUIDsForGroup(ctx context.Context, groupID string, levelType string, levelID string) ([]string, error) {
	client := iam.NewIAMClient(ctx, me.credentials)
	response, err := client.GET(ctx, fmt.Sprintf("/iam/v1/repo/%s/%s/bindings/groups/%s", levelType, levelID, groupID), rest2.RequestOptions{})
	if err != nil {
		return nil, err
	}

	var listResponse ListPolicyUUIDsForGroupResponse
	if err = json.Unmarshal(response.Data, &listResponse); err != nil {
		return nil, err
	}
	return listResponse.PolicyUUIDs, nil
}

func (me *BindingServiceClient) List(ctx context.Context) (api.Stubs, error) {
	client := iam.NewIAMClient(ctx, me.credentials)

	response, err := client.GET(ctx, fmt.Sprintf("/env/v2/accounts/%s/environments", me.credentials.IAM.AccountID), rest2.RequestOptions{})
	if err != nil {
		return nil, err
	}

	var envResponse ListEnvResponse
	if err = json.Unmarshal(response.Data, &envResponse); err != nil {
		return nil, err
	}

	response, err = client.GET(ctx, fmt.Sprintf("/iam/v1/repo/account/%s/bindings", me.credentials.IAM.AccountID), rest2.RequestOptions{})
	if err != nil {
		return nil, err
	}

	var policyBindingsResponse ListPolicyBindingsResponse
	if err = json.Unmarshal(response.Data, &policyBindingsResponse); err != nil {
		return nil, err
	}

	var stubs api.Stubs
	groupIds := map[string]bool{}
	for _, policy := range policyBindingsResponse.PolicyBindings {
		for _, group := range policy.Groups {
			if _, exists := groupIds[group]; !exists {
				id := join(group, "account", me.credentials.IAM.AccountID)
				stubs = append(stubs, &api.Stub{ID: id, Name: "PolicyBindings-" + id})
				groupIds[group] = true
			}
		}
	}

	for _, environment := range envResponse.Data {
		response, err := client.GET(ctx, fmt.Sprintf("/iam/v1/repo/environment/%s/bindings", environment.ID), rest2.RequestOptions{})
		if err != nil {
			return nil, err
		}

		var policyBindingsResponse ListPolicyBindingsResponse
		if err = json.Unmarshal(response.Data, &policyBindingsResponse); err != nil {
			return nil, err
		}

		groupIds := map[string]bool{}
		for _, policy := range policyBindingsResponse.PolicyBindings {
			for _, group := range policy.Groups {
				if _, exists := groupIds[group]; !exists {
					id := join(group, "environment", environment.ID)
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
		policyUUID, _, _, err := policies.SplitID(policyID, levelType, levelID)
		if err != nil {
			return err
		}
		if _, err = iam.NewIAMClient(ctx, me.credentials).DELETE(ctx, fmt.Sprintf("/iam/v1/repo/%s/%s/bindings/%s/%s", levelType, levelID, policyUUID, groupID), rest2.RequestOptions{}); err != nil {

			// The API currently returns a 400 Bad Request when the binding does not exist
			if apiErr, ok := errors.AsType[coreapi.APIError](err); ok && apiErr.StatusCode == http.StatusBadRequest {
				continue
			}
			return err
		}
	}
	return nil
}

func splitID(id string) (groupID string, levelType string, levelID string, err error) {
	parts := strings.Split(id, "#-#")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("%s is not a valid ID for a policy binding", id)
	}
	return parts[0], parts[1], parts[2], nil
}

func joinID(binding *bindings.PolicyBinding) string {
	levelType, levelID := getLevel(binding)
	return join(binding.GroupID, levelType, levelID)
}

func join(groupID string, levelType string, levelID string) string {
	return fmt.Sprintf("%s#-#%s#-#%s", groupID, levelType, levelID)
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
