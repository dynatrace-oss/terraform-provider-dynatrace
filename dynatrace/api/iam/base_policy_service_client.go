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

package iam

import (
	"context"
	"encoding/json"
	"fmt"

	rest2 "github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
)

type BasePolicyServiceClient struct {
	clientID     string
	accountID    string
	clientSecret string
	tokenURL     string
	endpointURL  string
}

func (me *BasePolicyServiceClient) ClientID() string {
	return me.clientID
}

func (me *BasePolicyServiceClient) AccountID() string {
	return me.accountID
}

func (me *BasePolicyServiceClient) ClientSecret() string {
	return me.clientSecret
}

func (me *BasePolicyServiceClient) TokenURL() string {
	return me.tokenURL
}

func (me *BasePolicyServiceClient) EndpointURL() string {
	return me.endpointURL
}

func NewBasePolicyService(clientID string, accountID string, clientSecret string, tokenURL string, endpointURL string) *BasePolicyServiceClient {
	return &BasePolicyServiceClient{clientID: clientID, accountID: accountID, clientSecret: clientSecret, tokenURL: tokenURL, endpointURL: endpointURL}
}

func (me *BasePolicyServiceClient) CREATE(ctx context.Context, level PolicyLevel, levelID string, policy *Policy) (string, error) {
	client := NewIAMClient(ctx, me)
	response, err := client.POST(ctx, fmt.Sprintf("/iam/v1/repo/%s/%s/policies", level, levelID), policy, rest2.RequestOptions{})
	if err != nil {
		return "", err
	}

	stub := PolicyStub{}
	if err = json.Unmarshal(response.Data, &stub); err != nil {
		return "", err
	}
	return stub.UUID, nil
}

func (me *BasePolicyServiceClient) GET(ctx context.Context, level PolicyLevel, levelID string, uuid string) (*Policy, error) {
	client := NewIAMClient(ctx, me)

	response, err := client.GET(ctx, fmt.Sprintf("/iam/v1/repo/%s/%s/policies/%s", level, levelID, uuid), rest2.RequestOptions{})
	if err != nil {
		return nil, err
	}

	var policyResponse Policy
	if err = json.Unmarshal(response.Data, &policyResponse); err != nil {
		return nil, err
	}
	return &policyResponse, nil
}

func (me *BasePolicyServiceClient) UPDATE(ctx context.Context, level PolicyLevel, levelID string, policy *Policy, uuid string) error {
	client := NewIAMClient(ctx, me)

	if _, err := client.PUT(ctx, fmt.Sprintf("/iam/v1/repo/%s/%s/policies/%s", level, levelID, uuid), policy, rest2.RequestOptions{}); err != nil {
		return err
	}
	return nil
}

type PolicyStub struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ListPoliciesResponse struct {
	Items []PolicyStub `json:"policies"`
}

func (me *BasePolicyServiceClient) List(ctx context.Context, level PolicyLevel, levelID string) ([]PolicyStub, error) {
	response, err := NewIAMClient(ctx, me).GET(ctx, fmt.Sprintf("/iam/v1/repo/%s/%s/policies", level, levelID), rest2.RequestOptions{})
	if err != nil {
		return nil, err
	}

	var policiesResponse ListPoliciesResponse
	if err = json.Unmarshal(response.Data, &policiesResponse); err != nil {
		return nil, err
	}
	return policiesResponse.Items, nil
}

func (me *BasePolicyServiceClient) LIST(ctx context.Context, level PolicyLevel, levelID string) ([]string, error) {
	userStubs, err := me.List(ctx, level, levelID)
	if err != nil {
		return nil, err
	}
	ids := []string{}
	for _, stub := range userStubs {
		ids = append(ids, stub.UUID)
	}
	return ids, nil
}

func (me *BasePolicyServiceClient) DELETE(ctx context.Context, level PolicyLevel, levelID string, uuid string) error {
	_, err := NewIAMClient(ctx, me).DELETE(ctx, fmt.Sprintf("/iam/v1/repo/%s/%s/policies/%s", level, levelID, uuid), rest2.RequestOptions{})
	return err
}
