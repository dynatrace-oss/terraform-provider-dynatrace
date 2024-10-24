package policies

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
)

var cachedEnvironmentIDs []string
var envIDFetcMutex sync.Mutex

// GetEnvironmentIDs retrieves all environmentIDs reachable via the given accountID
// The operation is guarded by a mutex
func GetEnvironmentIDs(ctx context.Context, auth iam.Authenticator) ([]string, error) {
	envIDFetcMutex.Lock()
	defer envIDFetcMutex.Unlock()
	if cachedEnvironmentIDs != nil {
		return cachedEnvironmentIDs, nil
	}
	return getEnvironmentIDs(ctx, auth)
}

// CheckPolicyExists attempts to fetch the details of a policy, identified by the given `policyUUID`
// and the presumed `levelType` and `levelID`. If the policy is not defined at the given `levelType`
// and `levelID` it returns false, without returning an error
// An error is returned ONLY if querying for the existence of the policy failed for another reason than `404 Not Found`
func CheckPolicyExists(ctx context.Context, auth iam.Authenticator, levelType string, levelID string, policyUUID string) (bool, string, error) {
	var err error

	response := struct {
		UUID string `json:"uuid"`
		Name string `json:"name"`
	}{}
	client := iam.NewIAMClient(auth)
	if err = iam.GET(client, ctx, fmt.Sprintf("%s/iam/v1/repo/%s/%s/policies/%s", auth.EndpointURL(), levelType, levelID, policyUUID), 200, false, &response); err != nil {
		// TODO: this is dirty. The IAM client unfortunately doesn't produce special kinds errors. string compare is the only option atm
		if strings.HasPrefix(err.Error(), "response code 404") {
			return false, "", nil
		}
		return false, "", err
	}
	return true, response.Name, nil
}

// FetchPolicyLevel determines the `levelType` and `levelID` of a policy identified by its UUID
// by trial and error, i.e. requests the policy from the REST API using all known combinations
//
// Option 1: The policy is a global policy (levelType = global, levelID = global)
// Option 2: The policy is on the account level, identified by the argument `accountID` (levelType = account, levelID = `accountID` argument)
// Option 3: The policy is on the environment level. ALL environments reachable via the account are being taken into consideration
//
// # If all attempts fail the returned error contains the UUID in its message
//
// This operation is guarded by a mutex
func FetchPolicyLevel(ctx context.Context, auth iam.Authenticator, uuid string) (levelType string, levelID string, name string, err error) {
	allPoliciesMutex.Lock()
	defer allPoliciesMutex.Unlock()
	return fetchPolicyLevel(ctx, auth, uuid)
}

// fetchPolicyLevel determines the `levelType` and `levelID` of a policy identified by its UUID
// by trial and error, i.e. requests the policy from the REST API using all known combinations
//
// Option 1: The policy is a global policy (levelType = global, levelID = global)
// Option 2: The policy is on the account level, identified by the argument `accountID` (levelType = account, levelID = `accountID` argument)
// Option 3: The policy is on the environment level. ALL environments reachable via the account are being taken into consideration
//
// # If all attempts fail the returned error contains the UUID in its message
//
// This operation is NOT guarded by a mutex. See `FetchPolicyLevel` for a guarded version
func fetchPolicyLevel(ctx context.Context, auth iam.Authenticator, uuid string) (levelType string, levelID string, name string, err error) {
	var exists bool
	if exists, name, err = CheckPolicyExists(ctx, auth, "global", "global", uuid); err != nil {
		return "", "", "", err
	}
	if exists {
		return "global", "global", name, nil
	}
	accountID := strings.TrimPrefix(auth.AccountID(), "urn:dtaccount:")
	if exists, name, err = CheckPolicyExists(ctx, auth, "account", accountID, uuid); err != nil {
		return "", "", name, err
	}
	if exists {
		return "account", accountID, name, nil
	}

	var environmentIDs []string
	if environmentIDs, err = GetEnvironmentIDs(ctx, auth); err != nil {
		return "", "", name, err
	}
	for _, environmentID := range environmentIDs {
		if exists, name, err = CheckPolicyExists(ctx, auth, "environment", environmentID, uuid); err != nil {
			return "", "", name, err
		}
		if exists {
			return "environment", environmentID, name, nil
		}
	}

	return "", "", name, rest.Error{Code: 404, Message: fmt.Sprintf("unable to resolve levelType and levelID of policy `%s`", uuid)}
}

// ResolvePolicyLevel determines the `levelType` and `levelID` of a policy using different strategies
//   - If it hasn't happened yet, all known policies are getting retrieved from the REST API
//   - In case none of the results contains the given `policyUUID` the `levelType` and `levelID` are getting
//     resolved using trial and error (see `fetchPolicyLevel`)
//
// This operation is guarded by a mutex
func ResolvePolicyLevel(ctx context.Context, auth iam.Authenticator, uuid string) (levelType string, levelID string, name string, err error) {
	allPoliciesMutex.Lock()
	defer allPoliciesMutex.Unlock()
	return resolvePolicyLevel(ctx, auth, uuid)
}

// ResolvePolicyLevel determines the `levelType` and `levelID` of a policy using different strategies
//   - If it hasn't happened yet, all known policies are getting retrieved from the REST API
//   - In case none of the results contains the given `uuid` the `levelType` and `levelID` are getting
//     resolved using trial and error (see `fetchPolicyLevel`)
//
// This operation is NOT guarded by a mutex. See `ResolvePolicyLevel` for a guarded version
func resolvePolicyLevel(ctx context.Context, auth iam.Authenticator, uuid string) (levelType string, levelID string, name string, err error) {
	allPolicyLevels, err := fetchAllPolicyLevels(ctx, auth)
	if err != nil {
		return "", "", "", err
	}
	pl, found := allPolicyLevels[uuid]
	if found {
		return pl.LevelType, pl.LevelID, pl.Name, nil
	}
	if levelType, levelID, name, err = fetchPolicyLevel(ctx, auth, uuid); err == nil {
		if err2 := registerPolicyLevel(ctx, auth, PolicyLevel{UUID: uuid, LevelType: levelType, LevelID: levelID, Name: name}); err2 != nil {
			return levelType, levelID, name, err2
		}
	}
	return levelType, levelID, name, err
}

type PolicyLevel struct {
	UUID      string
	LevelType string
	LevelID   string
	Name      string
}

var globalAllPolicyLevels map[string]PolicyLevel
var allPoliciesMutex sync.Mutex

// RegisterPolicyLevel notes down the `levelType` and `levelID` of the policy identified by the given `uuid`.
// Prior to that, if it hasn't happened yet, all known polices are getting pulled from the REST API
// In other words: Registering the `levelType` and `levelID` avoids that just partial information about
// the policys is stored locally. It's all nor nothing.
//
// # An error will be returned in case loading all known polices from the REST API fails for some reason
//
// This operation is guarded by a mutex.
func RegisterPolicyLevel(ctx context.Context, auth iam.Authenticator, level PolicyLevel) error {
	allPoliciesMutex.Lock()
	defer allPoliciesMutex.Unlock()
	return registerPolicyLevel(ctx, auth, level)
}

// registerPolicyLevel notes down the `levelType` and `levelID` of the policy identified by the given `uuid`.
// Prior to that, if it hasn't happened yet, all known polices are getting pulled from the REST API
// In other words: Registering the `levelType` and `levelID` avoids that just partial information about
// the policys is stored locally. It's all nor nothing.
//
// # An error will be returned in case loading all known polices from the REST API fails for some reason
//
// This operation is NOT guarded by a mutex. See `RegisterPolicyLevel` for a guarded version
func registerPolicyLevel(ctx context.Context, auth iam.Authenticator, level PolicyLevel) error {
	// fmt.Println("[POLICY-LEVEL]", "[REGISTER]", "["+level.UUID+"]", "BEGIN")
	// start := time.Now()
	// defer func() {
	// 	fmt.Println("[POLICY-LEVEL]", "[REGISTER]", "["+level.UUID+"]", fmt.Sprintf("... LASTED %v seconds", int64(time.Since(start).Seconds())))
	// }()

	if globalAllPolicyLevels == nil {
		_, err := fetchAllPolicyLevels(ctx, auth)
		if err != nil {
			return err
		}
	}
	globalAllPolicyLevels[level.UUID] = PolicyLevel{UUID: level.UUID, LevelType: level.LevelType, LevelID: level.LevelID, Name: level.Name}
	return nil
}

// FetchAllPolicyLevels pulls all known polices reachable via the given IAM Client from the REST API
// and notes down the `levelType` and `levelID` for these polices (identified via a UUID only).
//
// # You should use `ResolvePolicyLevel` to look up the `levelType` and `levelID` of a policy identifed by a UUID only
//
// This operation is guarded by a mutext
func FetchAllPolicyLevels(ctx context.Context, auth iam.Authenticator) (map[string]PolicyLevel, error) {
	allPoliciesMutex.Lock()
	defer allPoliciesMutex.Unlock()
	return fetchAllPolicyLevels(ctx, auth)
}

func fetchGlobalPolicies(ctx context.Context, auth iam.Authenticator) (results chan *api.Stub) {
	client := iam.NewIAMClient(auth)
	results = make(chan *api.Stub)
	go func() {
		defer func() {
			defer close(results)
		}()

		var response ListPoliciesResponse
		if err := iam.GET(client, ctx, fmt.Sprintf("%s/iam/v1/repo/global/global/policies", auth.EndpointURL()), 200, false, &response); err != nil {
			return
		}

		for _, policy := range response.Policies {
			results <- &api.Stub{ID: fmt.Sprintf("%s#-#%s#-#%s", policy.UUID, "global", "global"), Name: policy.Name}
		}
	}()

	return results
}

// fetchAllPolicyLevels pulls all known polices reachable via the given IAM Client from the REST API
// and notes down the `levelType` and `levelID` for these polices (identified via a UUID only).
//
// # You should use `ResolvePolicyLevel` to look up the `levelType` and `levelID` of a policy identifed by a UUID only
//
// This operation is NOT guarded by a mutext. See `FetchAllPolicyLevels` for a guarded version
func fetchAllPolicyLevels(ctx context.Context, auth iam.Authenticator) (m map[string]PolicyLevel, err error) {
	if globalAllPolicyLevels != nil {
		return globalAllPolicyLevels, nil
	}
	// start := time.Now()
	// defer func() {
	// 	fmt.Println("[POLICY-LEVEL]", "[FETCH-ALL]", fmt.Sprintf("... LASTED %v seconds", int64(time.Since(start).Seconds())))
	// }()
	nonGlobalStubs, err := list(ctx, auth)
	if err != nil {
		return nil, err
	}
	globalStubs := fetchGlobalPolicies(ctx, auth)

	m = map[string]PolicyLevel{}

	handleStub := func(stub *api.Stub) {
		if stub == nil {
			return
		}
		uuid, levelType, levelID, _ := SplitIDNoDefaults(stub.ID)
		m[uuid] = PolicyLevel{LevelType: levelType, LevelID: levelID, Name: stub.Name}
	}

	for {
		if nonGlobalStubs == nil && globalStubs == nil {
			break
		}
		select {
		case stub, more := <-nonGlobalStubs:
			handleStub(stub)
			if !more {
				nonGlobalStubs = nil
				if globalStubs == nil {
					break
				}
			}
		case stub, more := <-globalStubs:
			handleStub(stub)
			if !more {
				globalStubs = nil
				if nonGlobalStubs == nil {
					break
				}
			}
		}
		if nonGlobalStubs == nil && globalStubs == nil {
			break
		}
	}

	globalAllPolicyLevels = m

	return globalAllPolicyLevels, nil

}

// GetEnvironmentIDs retrieves all environmentIDs reachable via the given IAM Client
// The operation is NOT guarded by a mutex. See `GetEnvironmentIDs` for a guarded version
func getEnvironmentIDs(ctx context.Context, auth iam.Authenticator) ([]string, error) {
	client := iam.NewIAMClient(auth)
	accountID := strings.TrimPrefix(auth.AccountID(), "urn:dtaccount:")
	var err error

	var envResponse ListEnvResponse
	if err = iam.GET(client, ctx, fmt.Sprintf("%s/env/v2/accounts/%s/environments", auth.EndpointURL(), accountID), 200, false, &envResponse); err != nil {
		return nil, err
	}

	environmentIDs := []string{}

	for _, dataStub := range envResponse.Data {
		if len(dataStub.ID) > 0 {
			environmentIDs = append(environmentIDs, dataStub.ID)
		}
	}
	cachedEnvironmentIDs = environmentIDs
	return environmentIDs, nil
}
