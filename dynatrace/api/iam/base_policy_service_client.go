package iam

import (
	"context"
	"encoding/json"
	"fmt"
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
	var err error
	var responseBytes []byte

	client := NewIAMClient(me)
	if responseBytes, err = client.POST(ctx, fmt.Sprintf("%s/iam/v1/repo/%s/%s/policies", me.endpointURL, level, levelID), policy, 201, false); err != nil {
		return "", err
	}

	stub := PolicyStub{}
	if err = json.Unmarshal(responseBytes, &stub); err != nil {
		return "", err
	}
	return stub.UUID, nil
}

func (me *BasePolicyServiceClient) GET(ctx context.Context, level PolicyLevel, levelID string, uuid string) (*Policy, error) {
	var err error
	var responseBytes []byte

	client := NewIAMClient(me)

	if responseBytes, err = client.GET(ctx, fmt.Sprintf("%s/iam/v1/repo/%s/%s/policies/%s", me.endpointURL, level, levelID, uuid), 200, false); err != nil {
		return nil, err
	}

	var response Policy
	if err = json.Unmarshal(responseBytes, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (me *BasePolicyServiceClient) UPDATE(ctx context.Context, level PolicyLevel, levelID string, policy *Policy, uuid string) error {
	var err error

	client := NewIAMClient(me)

	if _, err = client.PUT(ctx, fmt.Sprintf("%s/iam/v1/repo/%s/%s/policies/%s", me.endpointURL, level, levelID, uuid), policy, 200, false); err != nil {
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
	var err error
	var responseBytes []byte

	if responseBytes, err = NewIAMClient(me).GET(ctx, fmt.Sprintf("%s/iam/v1/repo/%s/%s/policies", me.endpointURL, level, levelID), 200, false); err != nil {
		return nil, err
	}

	var response ListPoliciesResponse
	if err = json.Unmarshal(responseBytes, &response); err != nil {
		return nil, err
	}
	return response.Items, nil
}

func (me *BasePolicyServiceClient) LIST(ctx context.Context, level PolicyLevel, levelID string) ([]string, error) {
	var err error

	var userStubs []PolicyStub
	if userStubs, err = me.List(ctx, level, levelID); err != nil {
		return nil, err
	}
	ids := []string{}
	for _, stub := range userStubs {
		ids = append(ids, stub.UUID)
	}
	return ids, nil
}

func (me *BasePolicyServiceClient) DELETE(ctx context.Context, level PolicyLevel, levelID string, uuid string) error {
	_, err := NewIAMClient(me).DELETE(ctx, fmt.Sprintf("%s/iam/v1/repo/%s/%s/policies/%s", me.endpointURL, level, levelID, uuid), 204, false)
	return err
}
