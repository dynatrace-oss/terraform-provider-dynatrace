package iam

type AccountPolicyServiceClient struct {
	PolicyClient *BasePolicyServiceClient
}

func NewAccountPolicyService(clientID string, accountID string, clientSecret string) *AccountPolicyServiceClient {
	return &AccountPolicyServiceClient{PolicyClient: NewBasePolicyService(clientID, accountID, clientSecret)}
}

func (me *AccountPolicyServiceClient) CREATE(policy *Policy) (string, error) {
	return me.PolicyClient.CREATE(PolicyLevels.Account, me.PolicyClient.accountID, policy)
}

func (me *AccountPolicyServiceClient) GET(levelID string, uuid string) (*Policy, error) {
	return me.PolicyClient.GET(PolicyLevels.Account, me.PolicyClient.accountID, uuid)
}

func (me *AccountPolicyServiceClient) UPDATE(policy *Policy, uuid string) error {
	return me.PolicyClient.UPDATE(PolicyLevels.Account, me.PolicyClient.accountID, policy, uuid)
}

func (me *AccountPolicyServiceClient) List() ([]PolicyStub, error) {
	return me.PolicyClient.List(PolicyLevels.Account, me.PolicyClient.accountID)
}

func (me *AccountPolicyServiceClient) LIST(level PolicyLevel, levelID string) ([]string, error) {
	return me.PolicyClient.LIST(PolicyLevels.Account, me.PolicyClient.accountID)
}

func (me *AccountPolicyServiceClient) DELETE(level PolicyLevel, levelID string, uuid string) error {
	return me.PolicyClient.DELETE(PolicyLevels.Account, me.PolicyClient.accountID, uuid)
}
