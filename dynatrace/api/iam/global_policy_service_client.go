package iam

type GlobalPolicyServiceClient struct {
	PolicyClient *BasePolicyServiceClient
}

func NewGlobalPolicyService(clientID string, accountID string, clientSecret string) *GlobalPolicyServiceClient {
	return &GlobalPolicyServiceClient{PolicyClient: NewBasePolicyService(clientID, accountID, clientSecret)}
}

func (me *GlobalPolicyServiceClient) GET(levelID string, uuid string) (*Policy, error) {
	return me.PolicyClient.GET(PolicyLevels.Environment, me.PolicyClient.accountID, uuid)
}

func (me *GlobalPolicyServiceClient) List() ([]PolicyStub, error) {
	return me.PolicyClient.List(PolicyLevels.Environment, me.PolicyClient.accountID)
}

func (me *GlobalPolicyServiceClient) LIST(level PolicyLevel, levelID string) ([]string, error) {
	return me.PolicyClient.LIST(PolicyLevels.Environment, me.PolicyClient.accountID)
}

func (me *GlobalPolicyServiceClient) DELETE(level PolicyLevel, levelID string, uuid string) error {
	return me.PolicyClient.DELETE(PolicyLevels.Environment, me.PolicyClient.accountID, uuid)
}
