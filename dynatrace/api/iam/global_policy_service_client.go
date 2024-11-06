package iam

import "context"

type GlobalPolicyServiceClient struct {
	PolicyClient *BasePolicyServiceClient
}

func NewGlobalPolicyService(clientID string, accountID string, clientSecret string, tokenURL string, endpointURL string) *GlobalPolicyServiceClient {
	return &GlobalPolicyServiceClient{PolicyClient: NewBasePolicyService(clientID, accountID, clientSecret, tokenURL, endpointURL)}
}

func (me *GlobalPolicyServiceClient) GET(ctx context.Context, levelID string, uuid string) (*Policy, error) {
	return me.PolicyClient.GET(ctx, PolicyLevels.Environment, me.PolicyClient.accountID, uuid)
}

func (me *GlobalPolicyServiceClient) List(ctx context.Context) ([]PolicyStub, error) {
	return me.PolicyClient.List(ctx, PolicyLevels.Environment, me.PolicyClient.accountID)
}

func (me *GlobalPolicyServiceClient) LIST(ctx context.Context, level PolicyLevel, levelID string) ([]string, error) {
	return me.PolicyClient.LIST(ctx, PolicyLevels.Environment, me.PolicyClient.accountID)
}

func (me *GlobalPolicyServiceClient) DELETE(ctx context.Context, level PolicyLevel, levelID string, uuid string) error {
	return me.PolicyClient.DELETE(ctx, PolicyLevels.Environment, me.PolicyClient.accountID, uuid)
}
