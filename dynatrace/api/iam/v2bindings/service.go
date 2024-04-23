package v2bindings

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/policies"
	bindings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/v2bindings/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

type BindingServiceClient struct {
	clientID     string
	accountID    string
	clientSecret string
}

func (me *BindingServiceClient) ClientID() string {
	return me.clientID
}

func (me *BindingServiceClient) AccountID() string {
	return me.accountID
}

func (me *BindingServiceClient) ClientSecret() string {
	return me.clientSecret
}

func NewPolicyService(clientID string, accountID string, clientSecret string) *BindingServiceClient {
	return &BindingServiceClient{clientID: clientID, accountID: accountID, clientSecret: clientSecret}
}

func Service(credentials *settings.Credentials) settings.CRUDService[*bindings.PolicyBinding] {
	return &BindingServiceClient{clientID: credentials.IAM.ClientID, accountID: credentials.IAM.AccountID, clientSecret: credentials.IAM.ClientSecret}
}

func (me *BindingServiceClient) SchemaID() string {
	return "accounts:iam:bindings"
}

func (me *BindingServiceClient) Name() string {
	return me.SchemaID()
}

type PolicyCreateResponse struct {
	UUID string `json:"uuid"`
}

func (me *BindingServiceClient) Create(v *bindings.PolicyBinding) (*api.Stub, error) {
	id := joinID(v)
	var err error
	if err = me.Update(id, v); err != nil {
		return nil, err
	}
	return &api.Stub{ID: id, Name: "PolicyV2Bindings-" + id}, nil
}

func getStateConfig(ctx context.Context) *bindings.PolicyBinding {
	var stateConfig *bindings.PolicyBinding
	cfg := ctx.Value(settings.ContextKeyStateConfig)
	if typedConfig, ok := cfg.(*bindings.PolicyBinding); ok {
		stateConfig = typedConfig
	}
	return stateConfig
}

func (me *BindingServiceClient) Get(id string, v *bindings.PolicyBinding) error {
	return me.GetWithContext(context.Background(), id, v)
}

func deductPolicyID(policyID string, levelType string, levelID string, existing []*bindings.Policy) string {
	if len(existing) == 0 {
		return fmt.Sprintf("%s#-#%s#-#%s", policyID, levelType, levelID)
	}
	for _, policy := range existing {
		// * uuid#-#leveltype#-#levelID
		// * uuid
		// matches both cases.
		// If the current state contains just the UUID we don't want to concatenate
		if strings.Contains(policy.ID, policyID) {
			return policy.ID
		}
	}
	return fmt.Sprintf("%s#-#%s#-#%s", policyID, levelType, levelID)
}

type BindingsResponse struct {
	LevelType      string `json:"levelType"`
	LevelID        string `json:"levelId"`
	PolicyBindings []struct {
		PolicyUUID string            `json:"policyUuid"`
		GroupUUIDs []string          `json:"groups"`
		Parameters map[string]string `json:"parameters"`
		Metadata   map[string]string `json:"metadata"`
	} `json:"policyBindings"`
}

func (me *BindingServiceClient) GetWithContext(ctx context.Context, id string, v *bindings.PolicyBinding) error {
	stateConfig := getStateConfig(ctx)
	groupID, levelType, levelID, err := splitID(id)
	if err != nil {
		return err
	}
	var responseBytes []byte

	client := iam.NewIAMClient(me)

	if responseBytes, err = client.GET(fmt.Sprintf("https://api.dynatrace.com/iam/v1/repo/%s/%s/bindings/groups/%s", levelType, levelID, groupID), 200, false); err != nil {
		return err
	}
	policyUUIDStruct := struct {
		PolicyUuids []string `json:"policyUuids"`
	}{}
	if err = json.Unmarshal(responseBytes, &policyUUIDStruct); err != nil {
		return err
	}
	if levelType == "account" {
		v.Account = levelID
	} else if levelType == "environment" {
		v.Environment = levelID
	}
	v.GroupID = groupID
	policies := []*bindings.Policy{}

	var wg sync.WaitGroup
	var policyErr error
	var lock sync.Mutex

	for _, policyID := range policyUUIDStruct.PolicyUuids {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var bindingsResponse BindingsResponse
			if err = iam.GET(client, fmt.Sprintf("https://api.dynatrace.com/iam/v1/repo/%s/%s/bindings/%s/%s", levelType, levelID, policyID, groupID), 200, false, &bindingsResponse); err != nil {
				lock.Lock()
				policyErr = err
				lock.Unlock()
				return
			}
			if len(bindingsResponse.PolicyBindings) == 0 {
				return
			}

			policies, err := me.resolvePolicies(policyID, bindingsResponse, stateConfig)
			if err != nil {
				lock.Lock()
				policyErr = err
				lock.Unlock()
				return
			}
			for _, policy := range policies {
				lock.Lock()
				policies = append(policies, policy)
				lock.Unlock()
			}
		}()
	}
	wg.Wait()
	if policyErr != nil {
		return err
	}
	v.Policies = policies
	return nil
}

func (me *BindingServiceClient) resolvePolicies(uuid string, bindingsResponse BindingsResponse, stateConfig *bindings.PolicyBinding) ([]*bindings.Policy, error) {
	results := []*bindings.Policy{}
	for _, policyBinding := range bindingsResponse.PolicyBindings {
		existingPolicies := []*bindings.Policy{}
		if stateConfig != nil {
			existingPolicies = stateConfig.Policies
		}
		levelType, levelID, _, err := policies.ResolvePolicyLevel(me, uuid)
		if err != nil {
			return nil, err
		}
		policy := &bindings.Policy{ID: deductPolicyID(uuid, levelType, levelID, existingPolicies)}
		if len(policyBinding.Parameters) > 0 {
			policy.Parameters = map[string]string{}
			for key, value := range policyBinding.Parameters {
				policy.Parameters[key] = value
			}
		}
		if len(policyBinding.Metadata) > 0 {
			policy.Metadata = map[string]string{}
			for key, value := range policyBinding.Metadata {
				policy.Metadata[key] = value
			}
		}
		results = append(results, policy)
	}
	return results, nil
}

func (me *BindingServiceClient) Update(id string, v *bindings.PolicyBinding) error {
	groupID, levelType, levelID, err := splitID(id)
	if err != nil {
		return err
	}

	client := iam.NewIAMClient(me)

	policiesList := []*bindings.Policy{}

	for _, policy := range policiesList {
		policyUUID, _, _, _ := policies.SplitID(policy.ID, levelType, levelID)
		if _, err = client.DELETE_MULTI_RESPONSE(fmt.Sprintf("https://api.dynatrace.com/iam/v1/repo/%s/%s/bindings/%s/%s", levelType, levelID, policyUUID, groupID), []int{204, 400, 404}, false); err != nil {
			return err
		}
	}
	for _, policy := range policiesList {
		policyUUID, _, _, _ := policies.SplitID(policy.ID, levelType, levelID)
		payload := struct {
			Parameters map[string]string `json:"parameters"`
			Metadata   map[string]string `json:"metadata"`
		}{
			Parameters: policy.Parameters,
			Metadata:   policy.Metadata,
		}
		if _, err = client.POST(fmt.Sprintf("https://api.dynatrace.com/iam/v1/repo/%s/%s/bindings/%s/%s", levelType, levelID, policyUUID, groupID), payload, 204, false); err != nil {
			return err
		}
	}

	return nil
}

type DataStub struct {
	ID string `json:"id"`
}

type ListEnvResponse struct {
	Data []DataStub `json:"data"`
}

type PolicyBindingStub struct {
	Groups []string `json:"groups"`
}

type ListPolicyBindingsResponse struct {
	PolicyBindings []PolicyBindingStub `json:"policyBindings"`
}

func (me *BindingServiceClient) FetchAccountBindings() chan *api.Stub {
	results := make(chan *api.Stub)
	go func() {
		defer func() {
			close(results)
		}()

		var err error
		var responseBytes []byte

		client := iam.NewIAMClient(me)

		if responseBytes, err = client.GET(fmt.Sprintf("https://api.dynatrace.com/iam/v1/repo/account/%s/bindings", strings.TrimPrefix(me.AccountID(), "urn:dtaccount:")), 200, false); err != nil {
			return
		}

		var response ListPolicyBindingsResponse
		if err = json.Unmarshal(responseBytes, &response); err != nil {
			return
		}

		var stubs api.Stubs
		groupIds := map[string]bool{}
		for _, policy := range response.PolicyBindings {
			for _, group := range policy.Groups {
				if _, exists := groupIds[group]; !exists {
					id := fmt.Sprintf("%s#-#%s#-#%s", group, "account", strings.TrimPrefix(me.AccountID(), "urn:dtaccount:"))
					stubs = append(stubs, &api.Stub{ID: id, Name: "PolicyV2Bindings-" + id})
					groupIds[group] = true
				}
			}
		}
		for _, stub := range stubs {
			results <- stub
		}
	}()
	return results
}

func (me *BindingServiceClient) FetchEnvironmentBindings() chan *api.Stub {
	results := make(chan *api.Stub)
	go func() {
		defer func() {
			close(results)
		}()
		var err error
		var responseBytes []byte
		client := iam.NewIAMClient(me)

		var environmentIDs []string
		if environmentIDs, err = policies.GetEnvironmentIDs(me); err != nil {
			return
		}

		var stubs api.Stubs
		for _, environmentID := range environmentIDs {
			if responseBytes, err = client.GET(fmt.Sprintf("https://api.dynatrace.com/iam/v1/repo/environment/%s/bindings", environmentID), 200, false); err != nil {
				return
			}

			var response ListPolicyBindingsResponse
			if err = json.Unmarshal(responseBytes, &response); err != nil {
				return
			}

			groupIds := map[string]bool{}
			for _, policy := range response.PolicyBindings {
				for _, group := range policy.Groups {
					if _, exists := groupIds[group]; !exists {
						id := fmt.Sprintf("%s#-#%s#-#%s", group, "environment", environmentID)
						stubs = append(stubs, &api.Stub{ID: id, Name: "PolicyV2Bindings-" + id})
						groupIds[group] = true
					}
				}
			}
		}
		for _, stub := range stubs {
			results <- stub
		}
	}()
	return results
}

func (me *BindingServiceClient) List() (api.Stubs, error) {
	var stubs api.Stubs
	accountStubs := me.FetchAccountBindings()
	environmentStubs := me.FetchEnvironmentBindings()

	for {
		if accountStubs == nil && environmentStubs == nil {
			break
		}
		select {
		case stub, more := <-accountStubs:
			if more {
				stubs = append(stubs, stub)
			} else {
				accountStubs = nil
				if environmentStubs == nil {
					break
				}
			}
		case stub, more := <-environmentStubs:
			if more {
				stubs = append(stubs, stub)
			} else {
				environmentStubs = nil
				if accountStubs == nil {
					break
				}
			}
		}
		if accountStubs == nil && environmentStubs == nil {
			break
		}
	}
	return stubs, nil
}

func (me *BindingServiceClient) Delete(id string) error {
	groupID, levelType, levelID, err := splitID(id)
	if err != nil {
		return err
	}
	var binding bindings.PolicyBinding
	if err = me.Get(id, &binding); err != nil {
		return err
	}
	policyUUIDs := map[string]string{}
	for _, policy := range binding.Policies {
		policyUUID, _, _, err := policies.SplitID(policy.ID, levelType, levelID)
		if err != nil {
			return err
		}
		policyUUIDs[policyUUID] = policyUUID
	}
	for policyUUID := range policyUUIDs {
		if _, err = iam.NewIAMClient(me).DELETE(fmt.Sprintf("https://api.dynatrace.com/iam/v1/repo/%s/%s/bindings/%s/%s", levelType, levelID, policyUUID, groupID), 204, false); err != nil {
			return err
		}
	}
	return nil
}

func splitID(id string) (groupID string, levelType string, levelID string, err error) {
	parts := strings.Split(id, "#-#")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("%s is not a valid ID for a policy", id)
	}
	return parts[0], parts[1], parts[2], nil
}

func joinID(binding *bindings.PolicyBinding) string {
	levelType, levelID := getLevel(binding)
	return fmt.Sprintf("%s#-#%s#-#%s", binding.GroupID, levelType, levelID)
}

func getLevel(binding *bindings.PolicyBinding) (string, string) {
	if len(binding.Account) > 0 {
		return "account", binding.Account
	}
	if len(binding.Environment) > 0 {
		return "environment", binding.Environment
	}
	return "global", "global"
}
