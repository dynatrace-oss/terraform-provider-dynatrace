package iam

import (
	"encoding/json"
	"fmt"
)

type BasePolicyServiceClient struct {
	clientID     string
	accountID    string
	clientSecret string
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

func NewBasePolicyService(clientID string, accountID string, clientSecret string) *BasePolicyServiceClient {
	return &BasePolicyServiceClient{clientID: clientID, accountID: accountID, clientSecret: clientSecret}
}

func (me *BasePolicyServiceClient) CREATE(level PolicyLevel, levelID string, policy *Policy) (string, error) {
	var err error
	var responseBytes []byte

	client := NewIAMClient(me)
	if responseBytes, err = client.POST(fmt.Sprintf("https://api.dynatrace.com/iam/v1/repo/%s/%s/policies", level, levelID), policy, 201, false); err != nil {
		return "", err
	}

	stub := PolicyStub{}
	if err = json.Unmarshal(responseBytes, &stub); err != nil {
		return "", err
	}
	return stub.UUID, nil
}

func (me *BasePolicyServiceClient) GET(level PolicyLevel, levelID string, uuid string) (*Policy, error) {
	var err error
	var responseBytes []byte

	client := NewIAMClient(me)

	if responseBytes, err = client.GET(fmt.Sprintf("https://api.dynatrace.com/iam/v1/repo/%s/%s/policies/%s", level, levelID, uuid), 200, false); err != nil {
		return nil, err
	}

	var response Policy
	if err = json.Unmarshal(responseBytes, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (me *BasePolicyServiceClient) UPDATE(level PolicyLevel, levelID string, policy *Policy, uuid string) error {
	var err error

	client := NewIAMClient(me)

	if _, err = client.PUT(fmt.Sprintf("https://api.dynatrace.com/iam/v1/repo/%s/%s/policies/%s", level, levelID, uuid), policy, 200, false); err != nil {
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

func (me *BasePolicyServiceClient) List(level PolicyLevel, levelID string) ([]PolicyStub, error) {
	var err error
	var responseBytes []byte

	if responseBytes, err = NewIAMClient(me).GET(fmt.Sprintf("https://api.dynatrace.com/iam/v1/repo/%s/%s/policies", level, levelID), 200, false); err != nil {
		return nil, err
	}

	var response ListPoliciesResponse
	if err = json.Unmarshal(responseBytes, &response); err != nil {
		return nil, err
	}
	return response.Items, nil
}

func (me *BasePolicyServiceClient) LIST(level PolicyLevel, levelID string) ([]string, error) {
	var err error

	var userStubs []PolicyStub
	if userStubs, err = me.List(level, levelID); err != nil {
		return nil, err
	}
	ids := []string{}
	for _, stub := range userStubs {
		ids = append(ids, stub.UUID)
	}
	return ids, nil
}

func (me *BasePolicyServiceClient) DELETE(level PolicyLevel, levelID string, uuid string) error {
	_, err := NewIAMClient(me).DELETE(fmt.Sprintf("https://api.dynatrace.com/iam/v1/repo/%s/%s/policies/%s", level, levelID, uuid), 204, false)
	return err
}
