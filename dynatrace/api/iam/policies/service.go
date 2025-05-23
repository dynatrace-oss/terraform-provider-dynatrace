package policies

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam"
	policies "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/policies/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

type PolicyServiceClient struct {
	clientID     string
	accountID    string
	clientSecret string
	tokenURL     string
	endpointURL  string
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

func (me *PolicyServiceClient) TokenURL() string {
	return me.tokenURL
}

func (me *PolicyServiceClient) EndpointURL() string {
	return me.endpointURL
}

func NewPolicyService(clientID string, accountID string, clientSecret string, tokenURL string, endpointURL string) *PolicyServiceClient {
	return &PolicyServiceClient{clientID: clientID, accountID: accountID, clientSecret: clientSecret, tokenURL: tokenURL, endpointURL: endpointURL}
}

func Service(credentials *rest.Credentials) settings.CRUDService[*policies.Policy] {
	return &PolicyServiceClient{clientID: credentials.IAM.ClientID, accountID: credentials.IAM.AccountID, clientSecret: credentials.IAM.ClientSecret, tokenURL: credentials.IAM.TokenURL, endpointURL: credentials.IAM.EndpointURL}
}

func ServiceWithGloabals(credentials *rest.Credentials) *PolicyServiceClient {
	return &PolicyServiceClient{clientID: credentials.IAM.ClientID, accountID: credentials.IAM.AccountID, clientSecret: credentials.IAM.ClientSecret, tokenURL: credentials.IAM.TokenURL, endpointURL: credentials.IAM.EndpointURL}
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

func (me *PolicyServiceClient) Create(ctx context.Context, v *policies.Policy) (*api.Stub, error) {
	var err error
	var responseBytes []byte

	levelType, levelID := getLevel(v)

	client := iam.NewIAMClient(me)
	if responseBytes, err = client.POST(ctx, fmt.Sprintf("%s/iam/v1/repo/%s/%s/policies", me.endpointURL, levelType, levelID), v, 201, false); err != nil {
		return nil, err
	}
	var pcr PolicyCreateResponse
	if err = json.Unmarshal(responseBytes, &pcr); err != nil {
		return nil, err
	}
	v.UUID = pcr.UUID
	RegisterPolicyLevel(ctx, me, PolicyLevel{UUID: v.UUID, LevelType: levelType, LevelID: levelID})
	return &api.Stub{ID: joinID(pcr.UUID, v), Name: v.Name}, nil
}

func (me *PolicyServiceClient) Get(ctx context.Context, id string, v *policies.Policy) error {
	err := me.get(ctx, id, v)
	if err != nil {
		return err
	}
	if len(v.Account) == 0 && len(v.Environment) == 0 && !strings.HasSuffix(id, "#-#global#-#global") {
		return nil // TODO: investigate whether this can ever happen
	}
	return err
}

func UNUSED_PARAMETER(v any) {}

func (me *PolicyServiceClient) get(ctx context.Context, id string, v *policies.Policy) error {
	UNUSED_PARAMETER(ctx)
	var levelType string
	var levelID string

	uuid, _, _, err := SplitIDNoDefaults(id)
	if err != nil {
		return err
	}
	client := iam.NewIAMClient(me)
	var policyName string

	levelType, levelID, policyName, err = ResolvePolicyLevel(ctx, me, uuid)
	if err != nil {
		return err
	}
	if levelType == "global" && levelID == "global" {
		v.UUID = uuid
		v.Name = policyName
		return nil
	}

	if err = iam.GET(client, ctx, fmt.Sprintf("%s/iam/v1/repo/%s/%s/policies/%s", me.endpointURL, levelType, levelID, uuid), 200, false, &v); err != nil {
		return err
	}
	if levelType == "account" {
		v.Account = levelID
		if len(levelID) == 0 {
			return errors.New(fmt.Sprintf("Policy `%s` has level type `%s`, but level id is empty", id, levelType))
		}
	} else if levelType == "environment" {
		v.Environment = levelID
		if len(levelID) == 0 {
			return errors.New(fmt.Sprintf("Policy `%s` has level type `%s`, but level id is empty", id, levelType))
		}
	}
	v.UUID = uuid
	RegisterPolicyLevel(ctx, me, PolicyLevel{UUID: v.UUID, LevelType: levelType, LevelID: levelID})
	return nil
}

func (me *PolicyServiceClient) Update(ctx context.Context, id string, user *policies.Policy) error {
	var levelType string
	var levelID string
	uuid, _, _, err := SplitIDNoDefaults(id)
	if err != nil {
		return err
	}
	levelType, levelID, _, err = ResolvePolicyLevel(ctx, me, uuid)
	if err != nil {
		return err
	}
	client := iam.NewIAMClient(me)

	if _, err = client.PUT(ctx, fmt.Sprintf("%s/iam/v1/repo/%s/%s/policies/%s", me.endpointURL, levelType, levelID, uuid), user, 204, false); err != nil {
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

func listForEnvironment(ctx context.Context, auth iam.Authenticator, environmentID string) (results chan *api.Stub, err error) {
	results = make(chan *api.Stub)
	go func() {
		defer close(results)
		client := iam.NewIAMClient(auth)

		var response ListPoliciesResponse
		if err = iam.GET(client, ctx, fmt.Sprintf("%s/iam/v1/repo/environment/%s/policies", auth.EndpointURL(), environmentID), 200, false, &response); err != nil {
			return
		}

		for _, policy := range response.Policies {
			results <- &api.Stub{ID: Join(policy.UUID, "environment", environmentID), Name: policy.Name}
		}
	}()
	return results, nil
}

func listForEnvironments(ctx context.Context, auth iam.Authenticator) (results chan *api.Stub, err error) {
	results = make(chan *api.Stub)
	var environmentIDs []string
	if environmentIDs, err = GetEnvironmentIDs(ctx, auth); err != nil {
		return nil, err
	}
	var wg sync.WaitGroup

	for _, environmentID := range environmentIDs {
		wg.Add(1)
		var envIdxChan chan *api.Stub
		if envIdxChan, err = listForEnvironment(ctx, auth, environmentID); err != nil {
			wg.Done()
			return nil, err
		}
		go func(source, target chan *api.Stub) {
			defer wg.Done()
			for elem := range source {
				target <- elem
			}
		}(envIdxChan, results)
	}
	// once all goroutines are done close results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	return results, nil
}

func listForAccount(ctx context.Context, auth iam.Authenticator) (results chan *api.Stub, err error) {
	client := iam.NewIAMClient(auth)
	accountID := strings.TrimPrefix(auth.AccountID(), "urn:dtaccount:")
	results = make(chan *api.Stub)
	go func() {
		defer close(results)
		var response ListPoliciesResponse
		if err = iam.GET(client, ctx, fmt.Sprintf("%s/iam/v1/repo/account/%s/policies", auth.EndpointURL(), accountID), 200, false, &response); err != nil {
			return
		}

		for _, policy := range response.Policies {
			results <- &api.Stub{ID: Join(policy.UUID, "account", accountID), Name: policy.Name}
		}
	}()
	return results, nil
}

func list(ctx context.Context, auth iam.Authenticator) (results chan *api.Stub, err error) {
	results = make(chan *api.Stub)

	var envChan chan *api.Stub
	var accChan chan *api.Stub

	if envChan, err = listForEnvironments(ctx, auth); err != nil {
		return nil, err
	}
	if accChan, err = listForAccount(ctx, auth); err != nil {
		return nil, err
	}

	go func() {
		defer close(results)
		for {
			if envChan == nil && accChan == nil {
				break
			}
			select {
			case stub, more := <-envChan:
				if stub != nil {
					results <- stub
				}
				if !more {
					envChan = nil
					if accChan == nil {
						break
					}
				}
			case stub, more := <-accChan:
				if stub != nil {
					results <- stub
				}
				if !more {
					accChan = nil
					if envChan == nil {
						break
					}
				}
			}
			if envChan == nil && accChan == nil {
				break
			}
		}
	}()

	return results, nil
}

func (me *PolicyServiceClient) List(ctx context.Context) (api.Stubs, error) {
	stubs := api.Stubs{}
	policyLevels, err := FetchAllPolicyLevels(ctx, me)
	if err != nil {
		return stubs, err
	}
	for uuid, level := range policyLevels {
		if level.LevelType == "global" && level.LevelID == "global" {
			continue
		}
		stubs = append(stubs, &api.Stub{ID: Join(uuid, level.LevelType, level.LevelID), Name: level.Name})
	}
	return stubs, nil
}

func (me *PolicyServiceClient) ListWithGlobals(ctx context.Context) (api.Stubs, error) {
	stubs := api.Stubs{}
	policyLevels, err := FetchAllPolicyLevels(ctx, me)
	if err != nil {
		return stubs, err
	}
	for uuid, level := range policyLevels {
		stubs = append(stubs, &api.Stub{ID: Join(uuid, level.LevelType, level.LevelID), Name: level.Name})
	}
	return stubs, nil
}

func (me *PolicyServiceClient) Delete(ctx context.Context, id string) error {
	var levelType string
	var levelID string
	var err error
	var uuid string

	uuid, _, _, err = SplitIDNoDefaults(id)
	if err != nil {
		return err
	}
	levelType, levelID, _, err = ResolvePolicyLevel(ctx, me, uuid)
	if err != nil {
		return err
	}

	_, err = iam.NewIAMClient(me).DELETE(ctx, fmt.Sprintf("%s/iam/v1/repo/%s/%s/policies/%s", me.endpointURL, levelType, levelID, uuid), 204, false)
	return err
}

var uuidRegexp = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

func IsValidUUID(uuid string) bool {
	return uuidRegexp.MatchString(uuid)
}

// defLevelType and devLevelID are getting used in case the passed policyID is just its UUID
//
// In such a case the caller needs to have access to other configuration with these two strings
// e.g. the config object the policyIDs are embedded in
func SplitID(id string, defLevelType string, defLevelID string) (uuid string, levelType string, levelID string, err error) {
	if IsValidUUID(id) {
		return id, defLevelType, defLevelID, nil
	}
	return SplitIDNoDefaults(id)
}

func SplitIDNoDefaults(id string) (uuid string, levelType string, levelID string, err error) {
	if IsValidUUID(id) {
		return id, "", "", nil
	}
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

func Join(uuid string, levelType string, levelID string) string {
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
