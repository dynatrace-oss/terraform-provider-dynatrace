package iam

import "context"

type EnvironmentPolicyServiceClient struct {
	PolicyClient *BasePolicyServiceClient
}

func NewEnvironmentPolicyService(clientID string, accountID string, clientSecret string, tokenURL string, endpointURL string) *EnvironmentPolicyServiceClient {
	return &EnvironmentPolicyServiceClient{PolicyClient: NewBasePolicyService(clientID, accountID, clientSecret, tokenURL, endpointURL)}
}

func (me *EnvironmentPolicyServiceClient) CREATE(ctx context.Context, policy *Policy) (string, error) {
	return me.PolicyClient.CREATE(ctx, PolicyLevels.Environment, me.PolicyClient.accountID, policy)
}

func (me *EnvironmentPolicyServiceClient) GET(ctx context.Context, levelID string, uuid string) (*Policy, error) {
	return me.PolicyClient.GET(ctx, PolicyLevels.Environment, me.PolicyClient.accountID, uuid)
}

func (me *EnvironmentPolicyServiceClient) UPDATE(ctx context.Context, policy *Policy, uuid string) error {
	return me.PolicyClient.UPDATE(ctx, PolicyLevels.Environment, me.PolicyClient.accountID, policy, uuid)
}

func (me *EnvironmentPolicyServiceClient) List(ctx context.Context) ([]PolicyStub, error) {
	return me.PolicyClient.List(ctx, PolicyLevels.Environment, me.PolicyClient.accountID)
}

func (me *EnvironmentPolicyServiceClient) LIST(ctx context.Context, level PolicyLevel, levelID string) ([]string, error) {
	return me.PolicyClient.LIST(ctx, PolicyLevels.Environment, me.PolicyClient.accountID)
}

func (me *EnvironmentPolicyServiceClient) DELETE(ctx context.Context, level PolicyLevel, levelID string, uuid string) error {
	return me.PolicyClient.DELETE(ctx, PolicyLevels.Environment, me.PolicyClient.accountID, uuid)
}
