package policies

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam"
	policies "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/policies/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

type PolicyServiceClient struct {
	clientID     string
	accountID    string
	clientSecret string
}

func (me *PolicyServiceClient) ClientID() string {
	return me.clientID
}

func (me *PolicyServiceClient) AccountID() string {
	return me.accountID
}

func (me *PolicyServiceClient) ClientSecret() string {
	return me.clientSecret
}

func NewPolicyService(clientID string, accountID string, clientSecret string) *PolicyServiceClient {
	return &PolicyServiceClient{clientID: clientID, accountID: accountID, clientSecret: clientSecret}
}

func Service(credentials *settings.Credentials) settings.CRUDService[*policies.Policy] {
	return &PolicyServiceClient{clientID: credentials.IAM.ClientID, accountID: credentials.IAM.AccountID, clientSecret: credentials.IAM.ClientSecret}
}

func (me *PolicyServiceClient) SchemaID() string {
	return "accounts:iam:policies"
}

type PolicyCreateResponse struct {
	UUID string `json:"uuid"`
}

func (me *PolicyServiceClient) Create(v *policies.Policy) (*settings.Stub, error) {
	var err error
	var responseBytes []byte

	levelType, levelID := getLevel(v)

	client := iam.NewIAMClient(me)
	if responseBytes, err = client.POST(fmt.Sprintf("https://api.dynatrace.com/iam/v1/repo/%s/%s/policies", levelType, levelID), v, 201, false); err != nil {
		return nil, err
	}
	var pcr PolicyCreateResponse
	if err = json.Unmarshal(responseBytes, &pcr); err != nil {
		return nil, err
	}
	return &settings.Stub{ID: joinID(pcr.UUID, v), Name: v.Name}, nil
}

func (me *PolicyServiceClient) Get(id string, v *policies.Policy) error {
	uuid, levelType, levelID, err := SplitID(id)
	if err != nil {
		return err
	}
	var responseBytes []byte

	client := iam.NewIAMClient(me)

	if responseBytes, err = client.GET(fmt.Sprintf("https://api.dynatrace.com/iam/v1/repo/%s/%s/policies/%s", levelType, levelID, uuid), 200, false); err != nil {
		return err
	}
	if err = json.Unmarshal(responseBytes, &v); err != nil {
		return err
	}
	if levelType == "acount" {
		v.Account = levelID
	} else if levelType == "environment" {
		v.Environment = levelID
	}
	return nil
}

func (me *PolicyServiceClient) Update(id string, user *policies.Policy) error {
	uuid, levelType, levelID, err := SplitID(id)
	if err != nil {
		return err
	}

	client := iam.NewIAMClient(me)

	if _, err = client.PUT(fmt.Sprintf("https://api.dynatrace.com/iam/v1/repo/%s/%s/policies/%s", levelType, levelID, uuid), user, 204, false); err != nil {
		return err
	}
	return nil
}

func (me *PolicyServiceClient) List() (settings.Stubs, error) {
	return settings.Stubs{}, nil
	// var err error
	// var responseBytes []byte

	// if responseBytes, err = iam.NewIAMClient(me).GET(fmt.Sprintf("https://api.dynatrace.com/iam/v1/accounts/%s/users", me.AccountID()), 200, false); err != nil {
	// 	return nil, err
	// }

	// var response ListUsersResponse
	// if err = json.Unmarshal(responseBytes, &response); err != nil {
	// 	return nil, err
	// }
	// var stubs settings.Stubs
	// for _, item := range response.Items {
	// 	stubs = append(stubs, &settings.Stub{ID: item.UID, Name: item.Email})
	// }
	// return stubs, nil
}

func (me *PolicyServiceClient) Delete(id string) error {
	uuid, levelType, levelID, err := SplitID(id)
	if err != nil {
		return err
	}
	_, err = iam.NewIAMClient(me).DELETE(fmt.Sprintf("https://api.dynatrace.com/iam/v1/repo/%s/%s/policies/%s", levelType, levelID, uuid), 204, false)
	return err
}

func SplitID(id string) (uuid string, levelType string, levelID string, err error) {
	parts := strings.Split(id, "#-#")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("%s is not a valid ID for a policy", id)
	}
	return parts[0], parts[1], parts[2], nil
}

func joinID(uuid string, policy *policies.Policy) string {
	levelType, levelID := getLevel(policy)
	return fmt.Sprintf("%s#-#%s#-#%s", uuid, levelType, levelID)
}

func getLevel(policy *policies.Policy) (string, string) {
	if len(policy.Account) > 0 {
		return "account", policy.Account
	}
	if len(policy.Environment) > 0 {
		return "environment", policy.Environment
	}
	return "global", "global"
}
