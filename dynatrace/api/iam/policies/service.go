package policies

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
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

func (me *PolicyServiceClient) Name() string {
	return me.SchemaID()
}

type PolicyCreateResponse struct {
	UUID string `json:"uuid"`
}

func (me *PolicyServiceClient) Create(v *policies.Policy) (*api.Stub, error) {
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
	v.UUID = pcr.UUID
	return &api.Stub{ID: joinID(pcr.UUID, v), Name: v.Name}, nil
}

func (me *PolicyServiceClient) Get(id string, v *policies.Policy) error {
	uuid, levelType, levelID, err := SplitIDNoDefaults(id)
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
	if levelType == "account" {
		v.Account = levelID
	} else if levelType == "environment" {
		v.Environment = levelID
	}
	v.UUID = uuid
	return nil
}

func (me *PolicyServiceClient) Update(id string, user *policies.Policy) error {
	uuid, levelType, levelID, err := SplitIDNoDefaults(id)
	if err != nil {
		return err
	}

	client := iam.NewIAMClient(me)

	if _, err = client.PUT(fmt.Sprintf("https://api.dynatrace.com/iam/v1/repo/%s/%s/policies/%s", levelType, levelID, uuid), user, 204, false); err != nil {
		return err
	}
	return nil
}

type DataStub struct {
	ID string `json:"id"`
}

type ListEnvResponse struct {
	Data []DataStub `json:"data"`
}

type PolicyStub struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ListPoliciesResponse struct {
	Policies []PolicyStub `json:"policies"`
}

func (me *PolicyServiceClient) List() (api.Stubs, error) {
	var err error
	var responseBytes []byte
	client := iam.NewIAMClient(me)

	if responseBytes, err = client.GET(fmt.Sprintf("https://api.dynatrace.com/env/v2/accounts/%s/environments", strings.TrimPrefix(me.AccountID(), "urn:dtaccount:")), 200, false); err != nil {
		return nil, err
	}

	var envResponse ListEnvResponse
	if err = json.Unmarshal(responseBytes, &envResponse); err != nil {
		return nil, err
	}

	if responseBytes, err = client.GET(fmt.Sprintf("https://api.dynatrace.com/iam/v1/repo/account/%s/policies", strings.TrimPrefix(me.AccountID(), "urn:dtaccount:")), 200, false); err != nil {
		return nil, err
	}

	var response ListPoliciesResponse
	if err = json.Unmarshal(responseBytes, &response); err != nil {
		return nil, err
	}

	var stubs api.Stubs
	for _, policy := range response.Policies {
		stubs = append(stubs, &api.Stub{ID: fmt.Sprintf("%s#-#%s#-#%s", policy.UUID, "account", strings.TrimPrefix(me.AccountID(), "urn:dtaccount:")), Name: policy.Name})
	}

	for _, environment := range envResponse.Data {
		if responseBytes, err = client.GET(fmt.Sprintf("https://api.dynatrace.com/iam/v1/repo/environment/%s/policies", environment.ID), 200, false); err != nil {
			return nil, err
		}

		var response ListPoliciesResponse
		if err = json.Unmarshal(responseBytes, &response); err != nil {
			return nil, err
		}

		for _, policy := range response.Policies {
			stubs = append(stubs, &api.Stub{ID: fmt.Sprintf("%s#-#%s#-#%s", policy.UUID, "environment", environment.ID), Name: policy.Name})
		}
	}

	if responseBytes, err = client.GET("https://api.dynatrace.com/iam/v1/repo/global/global/policies", 200, false); err != nil {
		return nil, err
	}

	// ------ global policies ------
	var globalResponse ListPoliciesResponse
	if err = json.Unmarshal(responseBytes, &globalResponse); err != nil {
		return nil, err
	}

	for _, policy := range globalResponse.Policies {
		stubs = append(stubs, &api.Stub{ID: fmt.Sprintf("%s#-#%s#-#%s", policy.UUID, "global", "global"), Name: policy.Name})
	}
	return stubs, nil
}

func (me *PolicyServiceClient) Delete(id string) error {
	uuid, levelType, levelID, err := SplitIDNoDefaults(id)
	if err != nil {
		return err
	}
	_, err = iam.NewIAMClient(me).DELETE(fmt.Sprintf("https://api.dynatrace.com/iam/v1/repo/%s/%s/policies/%s", levelType, levelID, uuid), 204, false)
	return err
}

var uuidRegexp = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

func IsValidUUID(uuid string) bool {
	return uuidRegexp.MatchString(uuid)
}

// defLevelType and devLevelID are getting used
// in case the passed policyID is just its UUID
//
// In such a case the caller needs to have access
// to other configuration with these two strings
// e.g. the config object the policyIDs are embedded in
func SplitID(id string, defLevelType string, defLevelID string) (uuid string, levelType string, levelID string, err error) {
	if IsValidUUID(id) {
		return id, defLevelType, defLevelID, nil
	}
	return SplitIDNoDefaults(id)
}

func SplitIDNoDefaults(id string) (uuid string, levelType string, levelID string, err error) {
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
