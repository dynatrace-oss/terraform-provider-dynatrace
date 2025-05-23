package v2bindings

import (
	"context"
	"fmt"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/policies"
	bindings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/v2bindings/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

type BindingServiceClient struct {
	clientID     string
	accountID    string
	clientSecret string
	tokenURL     string
	endpointURL  string
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

func (me *BindingServiceClient) TokenURL() string {
	return me.tokenURL
}

func (me *BindingServiceClient) EndpointURL() string {
	return me.endpointURL
}

func NewPolicyService(clientID string, accountID string, clientSecret string, tokenURL string, endpointURL string) *BindingServiceClient {
	return &BindingServiceClient{clientID: clientID, accountID: accountID, clientSecret: clientSecret, tokenURL: tokenURL, endpointURL: endpointURL}
}

func Service(credentials *rest.Credentials) settings.CRUDService[*bindings.PolicyBinding] {
	return &BindingServiceClient{clientID: credentials.IAM.ClientID, accountID: credentials.IAM.AccountID, clientSecret: credentials.IAM.ClientSecret, tokenURL: credentials.IAM.TokenURL, endpointURL: credentials.IAM.EndpointURL}
}

func (me *BindingServiceClient) SchemaID() string {
	return "accounts:iam:bindings"
}

type PolicyCreateResponse struct {
	UUID string `json:"uuid"`
}

func (me *BindingServiceClient) Create(ctx context.Context, v *bindings.PolicyBinding) (*api.Stub, error) {
	id := joinID(v)
	var err error
	if err = me.Update(ctx, id, v); err != nil {
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
		Boundaries []string          `json:"boundaries"`
	} `json:"policyBindings"`
}

func (me *BindingServiceClient) Get(ctx context.Context, id string, v *bindings.PolicyBinding) error {
	stateConfig := getStateConfig(ctx)
	groupID, levelType, levelID, err := splitID(id)
	if err != nil {
		return err
	}
	client := iam.NewIAMClient(me)

	policyUUIDStruct := struct {
		PolicyUuids []string `json:"policyUuids"`
	}{}
	if err = iam.GET(client, ctx, fmt.Sprintf("%s/iam/v1/repo/%s/%s/bindings/groups/%s", me.endpointURL, levelType, levelID, groupID), 200, false, &policyUUIDStruct); err != nil {
		return err
	}
	if levelType == "account" {
		v.Account = levelID
	} else if levelType == "environment" {
		v.Environment = levelID
	}
	v.GroupID = groupID
	policies := []*bindings.Policy{}

	for _, policyID := range policyUUIDStruct.PolicyUuids {
		var bindingsResponse BindingsResponse
		if err = iam.GET(client, ctx, fmt.Sprintf("%s/iam/v1/repo/%s/%s/bindings/%s/%s", me.endpointURL, levelType, levelID, policyID, groupID), 200, false, &bindingsResponse); err != nil {
			return err
		}
		if len(bindingsResponse.PolicyBindings) == 0 {
			return nil
		}

		resolvedPolicies, err := me.resolvePolicies(ctx, policyID, bindingsResponse, stateConfig)
		if err != nil {
			return err
		}
		policies = append(policies, resolvedPolicies...)
	}
	v.Policies = policies
	return nil
}

func (me *BindingServiceClient) resolvePolicies(ctx context.Context, uuid string, bindingsResponse BindingsResponse, stateConfig *bindings.PolicyBinding) ([]*bindings.Policy, error) {
	results := []*bindings.Policy{}
	for _, policyBinding := range bindingsResponse.PolicyBindings {
		existingPolicies := []*bindings.Policy{}
		if stateConfig != nil {
			existingPolicies = stateConfig.Policies
		}
		levelType, levelID, _, err := policies.ResolvePolicyLevel(ctx, me, uuid)
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
		if len(policyBinding.Boundaries) > 0 {
			policy.Boundaries = append([]string{}, policyBinding.Boundaries...)
		}
		results = append(results, policy)
	}
	return results, nil
}

func (me *BindingServiceClient) Update(ctx context.Context, id string, v *bindings.PolicyBinding) error {
	groupID, levelType, levelID, err := splitID(id)
	if err != nil {
		return err
	}

	client := iam.NewIAMClient(me)

	policiesList := append([]*bindings.Policy{}, v.Policies...)

	for _, policy := range policiesList {
		policyUUID, _, _, _ := policies.SplitID(policy.ID, levelType, levelID)
		if _, err = client.DELETE_MULTI_RESPONSE(ctx, fmt.Sprintf("%s/iam/v1/repo/%s/%s/bindings/%s/%s", me.endpointURL, levelType, levelID, policyUUID, groupID), []int{204, 400, 404}, false); err != nil {
			return err
		}
	}
	for _, policy := range policiesList {
		policyUUID, _, _, _ := policies.SplitID(policy.ID, levelType, levelID)
		payload := struct {
			Parameters map[string]string `json:"parameters"`
			Metadata   map[string]string `json:"metadata"`
			Boundaries []string          `json:"boundaries"`
		}{
			Parameters: policy.Parameters,
			Metadata:   policy.Metadata,
			Boundaries: policy.Boundaries,
		}
		retries := 0
		for retries < 10 {
			if _, err = client.POST(ctx, fmt.Sprintf("%s/iam/v1/repo/%s/%s/bindings/%s/%s", me.endpointURL, levelType, levelID, policyUUID, groupID), payload, 204, false); err != nil {
				return err
			}
			break
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

func (me *BindingServiceClient) FetchAccountBindings(ctx context.Context) chan *api.Stub {
	results := make(chan *api.Stub)
	go func() {
		defer func() {
			close(results)
		}()

		var err error
		var response ListPolicyBindingsResponse
		client := iam.NewIAMClient(me)

		if err = iam.GET(client, ctx, fmt.Sprintf("%s/iam/v1/repo/account/%s/bindings", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:")), 200, false, &response); err != nil {
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

func (me *BindingServiceClient) FetchEnvironmentBindings(ctx context.Context) chan *api.Stub {
	results := make(chan *api.Stub)
	go func() {
		defer func() {
			close(results)
		}()
		var err error
		client := iam.NewIAMClient(me)

		var environmentIDs []string
		if environmentIDs, err = policies.GetEnvironmentIDs(ctx, me); err != nil {
			return
		}

		var stubs api.Stubs
		for _, environmentID := range environmentIDs {
			var response ListPolicyBindingsResponse
			if err = iam.GET(client, ctx, fmt.Sprintf("%s/iam/v1/repo/environment/%s/bindings", me.endpointURL, environmentID), 200, false, &response); err != nil {
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

func (me *BindingServiceClient) List(ctx context.Context) (api.Stubs, error) {
	var stubs api.Stubs
	accountStubs := me.FetchAccountBindings(ctx)
	environmentStubs := me.FetchEnvironmentBindings(ctx)

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

func (me *BindingServiceClient) Delete(ctx context.Context, id string) error {
	groupID, levelType, levelID, err := splitID(id)
	if err != nil {
		return err
	}
	var binding bindings.PolicyBinding
	if err = me.Get(ctx, id, &binding); err != nil {
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
		if _, err = iam.NewIAMClient(me).DELETE(ctx, fmt.Sprintf("%s/iam/v1/repo/%s/%s/bindings/%s/%s", me.endpointURL, levelType, levelID, policyUUID, groupID), 204, false); err != nil {
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
