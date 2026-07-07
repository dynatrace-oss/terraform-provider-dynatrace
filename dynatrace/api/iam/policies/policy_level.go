/**
* @license
* Copyright 2025 Dynatrace LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package policies

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	coreapi "github.com/dynatrace/dynatrace-configuration-as-code-core/api"
	rest2 "github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
)

var cachedEnvironmentIDs []string
var envIDFetcMutex sync.Mutex

// GetEnvironmentIDs retrieves all environmentIDs reachable via the given accountID
// The operation is guarded by a mutex
func GetEnvironmentIDs(ctx context.Context, credentials *rest.Credentials) ([]string, error) {
	envIDFetcMutex.Lock()
	defer envIDFetcMutex.Unlock()
	if cachedEnvironmentIDs != nil {
		return cachedEnvironmentIDs, nil
	}
	return getEnvironmentIDs(ctx, credentials)
}

// CheckPolicyExists attempts to fetch the details of a policy, identified by the given `policyUUID`
// and the presumed `levelType` and `levelID`. If the policy is not defined at the given `levelType`
// and `levelID` it returns false, without returning an error
// An error is returned ONLY if querying for the existence of the policy failed for another reason than `404 Not Found`
func CheckPolicyExists(ctx context.Context, credentials *rest.Credentials, levelType string, levelID string, policyUUID string) (bool, string, error) {
	levelPolicy := struct {
		UUID string `json:"uuid"`
		Name string `json:"name"`
	}{}
	client := iam.NewIAMClient(ctx, credentials)
	response, err := client.GET(ctx, fmt.Sprintf("/iam/v1/repo/%s/%s/policies/%s", levelType, levelID, policyUUID), rest2.RequestOptions{})
	if err != nil {
		// A 404 Not Found means the policy is not defined at the given level - that is not an error here
		if coreapi.IsNotFoundError(err) {
			return false, "", nil
		}
		return false, "", err
	}
	if err = json.Unmarshal(response.Data, &levelPolicy); err != nil {
		return false, "", err
	}
	return true, levelPolicy.Name, nil
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
func FetchPolicyLevel(ctx context.Context, credentials *rest.Credentials, uuid string) (levelType string, levelID string, name string, err error) {
	allPoliciesMutex.Lock()
	defer allPoliciesMutex.Unlock()
	return fetchPolicyLevel(ctx, credentials, uuid)
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
func fetchPolicyLevel(ctx context.Context, credentials *rest.Credentials, uuid string) (levelType string, levelID string, name string, err error) {
	var exists bool
	if exists, name, err = CheckPolicyExists(ctx, credentials, "global", "global", uuid); err != nil {
		return "", "", "", err
	}
	if exists {
		return "global", "global", name, nil
	}

	if exists, name, err = CheckPolicyExists(ctx, credentials, "account", credentials.IAM.AccountID, uuid); err != nil {
		return "", "", name, err
	}
	if exists {
		return "account", credentials.IAM.AccountID, name, nil
	}

	var environmentIDs []string
	if environmentIDs, err = GetEnvironmentIDs(ctx, credentials); err != nil {
		return "", "", name, err
	}
	for _, environmentID := range environmentIDs {
		if exists, name, err = CheckPolicyExists(ctx, credentials, "environment", environmentID, uuid); err != nil {
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
func ResolvePolicyLevel(ctx context.Context, credentials *rest.Credentials, uuid string) (levelType string, levelID string, name string, err error) {
	allPoliciesMutex.Lock()
	defer allPoliciesMutex.Unlock()
	return resolvePolicyLevel(ctx, credentials, uuid)
}

// ResolvePolicyLevel determines the `levelType` and `levelID` of a policy using different strategies
//   - If it hasn't happened yet, all known policies are getting retrieved from the REST API
//   - In case none of the results contains the given `uuid` the `levelType` and `levelID` are getting
//     resolved using trial and error (see `fetchPolicyLevel`)
//
// This operation is NOT guarded by a mutex. See `ResolvePolicyLevel` for a guarded version
func resolvePolicyLevel(ctx context.Context, credentials *rest.Credentials, uuid string) (levelType string, levelID string, name string, err error) {
	allPolicyLevels, err := fetchAllPolicyLevels(ctx, credentials)
	if err != nil {
		return "", "", "", err
	}
	pl, found := allPolicyLevels[uuid]
	if found {
		return pl.LevelType, pl.LevelID, pl.Name, nil
	}
	if levelType, levelID, name, err = fetchPolicyLevel(ctx, credentials, uuid); err == nil {
		if err2 := registerPolicyLevel(ctx, credentials, PolicyLevel{UUID: uuid, LevelType: levelType, LevelID: levelID, Name: name}); err2 != nil {
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
func RegisterPolicyLevel(ctx context.Context, credentials *rest.Credentials, level PolicyLevel) error {
	allPoliciesMutex.Lock()
	defer allPoliciesMutex.Unlock()
	return registerPolicyLevel(ctx, credentials, level)
}

// registerPolicyLevel notes down the `levelType` and `levelID` of the policy identified by the given `uuid`.
// Prior to that, if it hasn't happened yet, all known polices are getting pulled from the REST API
// In other words: Registering the `levelType` and `levelID` avoids that just partial information about
// the policys is stored locally. It's all nor nothing.
//
// # An error will be returned in case loading all known polices from the REST API fails for some reason
//
// This operation is NOT guarded by a mutex. See `RegisterPolicyLevel` for a guarded version
func registerPolicyLevel(ctx context.Context, credentials *rest.Credentials, level PolicyLevel) error {
	// fmt.Println("[POLICY-LEVEL]", "[REGISTER]", "["+level.UUID+"]", "BEGIN")
	// start := time.Now()
	// defer func() {
	// 	fmt.Println("[POLICY-LEVEL]", "[REGISTER]", "["+level.UUID+"]", fmt.Sprintf("... LASTED %v seconds", int64(time.Since(start).Seconds())))
	// }()

	if globalAllPolicyLevels == nil {
		_, err := fetchAllPolicyLevels(ctx, credentials)
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
func FetchAllPolicyLevels(ctx context.Context, credentials *rest.Credentials) (map[string]PolicyLevel, error) {
	allPoliciesMutex.Lock()
	defer allPoliciesMutex.Unlock()
	return fetchAllPolicyLevels(ctx, credentials)
}

func fetchGlobalPolicies(ctx context.Context, credentials *rest.Credentials) (results chan *api.Stub) {
	client := iam.NewIAMClient(ctx, credentials)
	results = make(chan *api.Stub)
	go func() {
		defer func() {
			defer close(results)
		}()

		var policyList ListPoliciesResponse
		response, err := client.GET(ctx, "/iam/v1/repo/global/global/policies", rest2.RequestOptions{})
		if err != nil {
			return
		}
		if err := json.Unmarshal(response.Data, &policyList); err != nil {
			return
		}

		for _, policy := range policyList.Policies {
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
func fetchAllPolicyLevels(ctx context.Context, credentials *rest.Credentials) (m map[string]PolicyLevel, err error) {
	if globalAllPolicyLevels != nil {
		return globalAllPolicyLevels, nil
	}
	// start := time.Now()
	// defer func() {
	// 	fmt.Println("[POLICY-LEVEL]", "[FETCH-ALL]", fmt.Sprintf("... LASTED %v seconds", int64(time.Since(start).Seconds())))
	// }()
	nonGlobalStubs, err := list(ctx, credentials)
	if err != nil {
		return nil, err
	}
	globalStubs := fetchGlobalPolicies(ctx, credentials)

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
func getEnvironmentIDs(ctx context.Context, credentials *rest.Credentials) ([]string, error) {
	client := iam.NewIAMClient(ctx, credentials)

	var envResponse ListEnvResponse
	response, err := client.GET(ctx, fmt.Sprintf("/env/v2/accounts/%s/environments", credentials.IAM.AccountID), rest2.RequestOptions{})
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(response.Data, &envResponse); err != nil {
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
