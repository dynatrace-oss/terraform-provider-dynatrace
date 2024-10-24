package iam

import "context"

type AccountPolicyServiceClient struct {
	PolicyClient *BasePolicyServiceClient
}

func NewAccountPolicyService(clientID string, accountID string, clientSecret string, tokenURL string, endpointURL string) *AccountPolicyServiceClient {
	return &AccountPolicyServiceClient{PolicyClient: NewBasePolicyService(clientID, accountID, clientSecret, tokenURL, endpointURL)}
}

func (me *AccountPolicyServiceClient) CREATE(ctx context.Context, policy *Policy) (string, error) {
	return me.PolicyClient.CREATE(ctx, PolicyLevels.Account, me.PolicyClient.accountID, policy)
}

func (me *AccountPolicyServiceClient) GET(ctx context.Context, levelID string, uuid string) (*Policy, error) {
	return me.PolicyClient.GET(ctx, PolicyLevels.Account, me.PolicyClient.accountID, uuid)
}

func (me *AccountPolicyServiceClient) UPDATE(ctx context.Context, policy *Policy, uuid string) error {
	return me.PolicyClient.UPDATE(ctx, PolicyLevels.Account, me.PolicyClient.accountID, policy, uuid)
}

func (me *AccountPolicyServiceClient) List(ctx context.Context) ([]PolicyStub, error) {
	return me.PolicyClient.List(ctx, PolicyLevels.Account, me.PolicyClient.accountID)
}

func (me *AccountPolicyServiceClient) LIST(ctx context.Context, level PolicyLevel, levelID string) ([]string, error) {
	return me.PolicyClient.LIST(ctx, PolicyLevels.Account, me.PolicyClient.accountID)
}

func (me *AccountPolicyServiceClient) DELETE(ctx context.Context, level PolicyLevel, levelID string, uuid string) error {
	return me.PolicyClient.DELETE(ctx, PolicyLevels.Account, me.PolicyClient.accountID, uuid)
}
