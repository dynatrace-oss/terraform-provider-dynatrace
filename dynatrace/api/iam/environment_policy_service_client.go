package iam

type EnvironmentPolicyServiceClient struct {
	PolicyClient *BasePolicyServiceClient
}

func NewEnvironmentPolicyService(clientID string, accountID string, clientSecret string) *EnvironmentPolicyServiceClient {
	return &EnvironmentPolicyServiceClient{PolicyClient: NewBasePolicyService(clientID, accountID, clientSecret)}
}

func (me *EnvironmentPolicyServiceClient) CREATE(policy *Policy) (string, error) {
	return me.PolicyClient.CREATE(PolicyLevels.Environment, me.PolicyClient.accountID, policy)
}

func (me *EnvironmentPolicyServiceClient) GET(levelID string, uuid string) (*Policy, error) {
	return me.PolicyClient.GET(PolicyLevels.Environment, me.PolicyClient.accountID, uuid)
}

func (me *EnvironmentPolicyServiceClient) UPDATE(policy *Policy, uuid string) error {
	return me.PolicyClient.UPDATE(PolicyLevels.Environment, me.PolicyClient.accountID, policy, uuid)
}

func (me *EnvironmentPolicyServiceClient) List() ([]PolicyStub, error) {
	return me.PolicyClient.List(PolicyLevels.Environment, me.PolicyClient.accountID)
}

func (me *EnvironmentPolicyServiceClient) LIST(level PolicyLevel, levelID string) ([]string, error) {
	return me.PolicyClient.LIST(PolicyLevels.Environment, me.PolicyClient.accountID)
}

func (me *EnvironmentPolicyServiceClient) DELETE(level PolicyLevel, levelID string, uuid string) error {
	return me.PolicyClient.DELETE(PolicyLevels.Environment, me.PolicyClient.accountID, uuid)
}
