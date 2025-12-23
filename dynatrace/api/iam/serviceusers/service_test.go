//go:build unit

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

package serviceusers

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam"
	serviceusers "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/serviceusers/settings"

	"github.com/stretchr/testify/assert"
)

const testAccountID = "test-account-id"
const testEndpointURL = "https://api-test.dynatrace.com"

type mockIAMClient struct {
	POSTFunc   func(ctx context.Context, url string, payload any, expectedResponseCode int, forceNewBearer bool) ([]byte, error)
	PUTFunc    func(ctx context.Context, url string, payload any, expectedResponseCode int, forceNewBearer bool) ([]byte, error)
	GETFunc    func(ctx context.Context, url string, expectedResponseCode int, forceNewBearer bool) ([]byte, error)
	DELETEFunc func(ctx context.Context, url string, expectedResponseCode int, forceNewBearer bool) ([]byte, error)
}

func (me *mockIAMClient) POST(ctx context.Context, url string, payload any, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
	return me.POSTFunc(ctx, url, payload, expectedResponseCode, forceNewBearer)
}

func (me *mockIAMClient) PUT(ctx context.Context, url string, payload any, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
	return me.PUTFunc(ctx, url, payload, expectedResponseCode, forceNewBearer)
}

func (me *mockIAMClient) PUT_MULTI_RESPONSE(ctx context.Context, url string, payload any, expectedResponseCodes []int, forceNewBearer bool) ([]byte, error) {
	panic("mock doesnt support PUT_MULTI_RESPONSE")
}

func (me *mockIAMClient) GET(ctx context.Context, url string, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
	return me.GETFunc(ctx, url, expectedResponseCode, forceNewBearer)
}

func (me *mockIAMClient) DELETE(ctx context.Context, url string, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
	return me.DELETEFunc(ctx, url, expectedResponseCode, forceNewBearer)
}

func (me *mockIAMClient) DELETE_MULTI_RESPONSE(ctx context.Context, url string, expectedResponseCodes []int, forceNewBearer bool) ([]byte, error) {
	panic("mock doesnt support DELETE_MULTI_RESPONSE")
}

func createTestServiceUserServiceClient(client *mockIAMClient) *ServiceUserServiceClient {
	return &ServiceUserServiceClient{
		iamClient: &mockIAMClientGetter{
			client: client,
		},
		accountID:   testAccountID,
		endpointURL: testEndpointURL,
	}
}

type mockIAMClientGetter struct {
	client *mockIAMClient
}

func (me *mockIAMClientGetter) NewIAMClient() iam.IAMClient {
	return me.client
}

func TestService_Create(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		mockClient := &mockIAMClient{
			POSTFunc: func(ctx context.Context, url string, payload any, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
				assert.Equal(t, fmt.Sprintf("%s/iam/v1/accounts/%s/service-users", testEndpointURL, testAccountID), url)
				return []byte(`{"uid":"test-uid","email":"test@example.com","name":"Test User"}`), nil
			},
			PUTFunc: func(ctx context.Context, url string, payload any, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
				assert.Equal(t, fmt.Sprintf("%s/iam/v1/accounts/%s/users/%s/groups", testEndpointURL, testAccountID, "test@example.com"), url)
				return nil, nil
			},
		}

		client := createTestServiceUserServiceClient(mockClient)
		serviceUser := &serviceusers.ServiceUser{
			Name:        "Test User",
			Email:       "test@example.com",
			Description: "Test description",
			Groups:      []string{"group-1", "group-2"},
		}

		stub, err := client.Create(t.Context(), serviceUser)
		assert.NoError(t, err)
		assert.Equal(t, "test-uid", stub.ID)
		assert.Equal(t, "Test User", stub.Name)
	})

	t.Run("creation fails on POST", func(t *testing.T) {
		mockClient := &mockIAMClient{
			POSTFunc: func(ctx context.Context, url string, payload any, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
				return nil, errors.New("POST failed")
			},
		}

		client := createTestServiceUserServiceClient(mockClient)
		serviceUser := &serviceusers.ServiceUser{
			Name:  "Test User",
			Email: "test@example.com",
		}

		_, err := client.Create(t.Context(), serviceUser)
		assert.EqualError(t, err, "POST failed")
	})

	t.Run("creation fails on group assignment and cleanup succeeds", func(t *testing.T) {
		mockClient := &mockIAMClient{
			POSTFunc: func(ctx context.Context, url string, payload any, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
				return []byte(`{"uid":"test-uid","email":"test@example.com","name":"Test User"}`), nil
			},
			PUTFunc: func(ctx context.Context, url string, payload any, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
				return nil, errors.New("group assignment failed")
			},
			DELETEFunc: func(ctx context.Context, url string, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
				return nil, nil
			},
		}

		client := createTestServiceUserServiceClient(mockClient)
		serviceUser := &serviceusers.ServiceUser{
			Name:   "Test User",
			Email:  "test@example.com",
			Groups: []string{"group-1"},
		}

		_, err := client.Create(t.Context(), serviceUser)
		assert.EqualError(t, err, "group assignment failed")
	})

	t.Run("creation fails on group assignment and cleanup fails", func(t *testing.T) {
		mockClient := &mockIAMClient{
			POSTFunc: func(ctx context.Context, url string, payload any, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
				return []byte(`{"uid":"test-uid","email":"test@example.com","name":"Test User"}`), nil
			},
			PUTFunc: func(ctx context.Context, url string, payload any, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
				return nil, errors.New("group assignment failed")
			},
			DELETEFunc: func(ctx context.Context, url string, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
				return nil, errors.New("delete failed")
			},
		}

		client := createTestServiceUserServiceClient(mockClient)
		serviceUser := &serviceusers.ServiceUser{
			Name:   "Test User",
			Email:  "test@example.com",
			Groups: []string{"group-1"},
		}

		_, err := client.Create(t.Context(), serviceUser)
		assert.EqualError(t, err, "failed to create service user: group assignment failed; additionally failed to clean up service user: delete failed")
	})
}

func TestService_Get(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		mockClient := &mockIAMClient{
			GETFunc: func(ctx context.Context, url string, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
				expectedURL := fmt.Sprintf("%s/iam/v1/accounts/%s/service-users/%s", testEndpointURL, testAccountID, "test-uid")
				assert.Equal(t, expectedURL, url)
				return []byte(`{
					"uid": "test-uid",
					"email": "test@example.com",
					"name": "Test User",
					"description": "Test description",
					"groups": [
						{"groupName": "Group 1", "uuid": "group-uuid-1"},
						{"groupName": "Group 2", "uuid": "group-uuid-2"}
					]
				}`), nil
			},
		}

		client := createTestServiceUserServiceClient(mockClient)
		serviceUser := serviceusers.ServiceUser{}

		err := client.Get(t.Context(), "test-uid", &serviceUser)
		assert.NoError(t, err)

		assert.Equal(t, serviceusers.ServiceUser{
			UID:         "test-uid",
			Email:       "test@example.com",
			Name:        "Test User",
			Description: "Test description",
			Groups:      []string{"group-uuid-1", "group-uuid-2"},
		}, serviceUser)
	})

	t.Run("get fails", func(t *testing.T) {
		mockClient := &mockIAMClient{
			GETFunc: func(ctx context.Context, url string, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
				return nil, errors.New("GET failed")
			},
		}

		client := createTestServiceUserServiceClient(mockClient)
		serviceUser := &serviceusers.ServiceUser{}

		err := client.Get(t.Context(), "test-uid", serviceUser)
		assert.EqualError(t, err, "GET failed")
	})

	t.Run("get with empty groups", func(t *testing.T) {
		mockClient := &mockIAMClient{
			GETFunc: func(ctx context.Context, url string, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
				return []byte(`{
					"uid": "test-uid",
					"email": "test@example.com",
					"name": "Test User",
					"description": "",
					"groups": []
				}`), nil
			},
		}

		client := createTestServiceUserServiceClient(mockClient)
		serviceUser := &serviceusers.ServiceUser{}

		err := client.Get(t.Context(), "test-uid", serviceUser)
		assert.NoError(t, err)
		assert.Empty(t, serviceUser.Groups)
	})
}

func TestService_List(t *testing.T) {
	t.Run("successful list without pagination", func(t *testing.T) {
		mockClient := &mockIAMClient{
			GETFunc: func(ctx context.Context, url string, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
				return []byte(`{
					"count": 2,
					"results": [
						{"uid": "uid-1", "email": "user1@example.com", "name": "User 1"},
						{"uid": "uid-2", "email": "user2@example.com", "name": "User 2"}
					]
				}`), nil
			},
		}

		client := createTestServiceUserServiceClient(mockClient)
		stubs, err := client.List(t.Context())
		assert.NoError(t, err)
		assert.Len(t, stubs, 2)
		assert.Equal(t, "uid-1", stubs[0].ID)
		assert.Equal(t, "User 1", stubs[0].Name)
		assert.Equal(t, "uid-2", stubs[1].ID)
		assert.Equal(t, "User 2", stubs[1].Name)
	})

	t.Run("successful list with pagination", func(t *testing.T) {
		callCount := 0
		mockClient := &mockIAMClient{
			GETFunc: func(ctx context.Context, url string, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
				callCount++
				if callCount == 1 {
					assert.Equal(t, fmt.Sprintf("%s/iam/v1/accounts/%s/service-users", testEndpointURL, testAccountID), url)
					return []byte(`{
						"count": 2,
						"results": [
							{"uid": "uid-1", "email": "user1@example.com", "name": "User 1"}
						],
						"nextPageKey": "page2"
					}`), nil
				}
				if callCount == 2 {
					assert.Equal(t, fmt.Sprintf("%s/iam/v1/accounts/%s/service-users?nextPageKey=page2", testEndpointURL, testAccountID), url)
					return []byte(`{
					"count": 2,
					"results": [
						{"uid": "uid-2", "email": "user2@example.com", "name": "User 2"}
					]
				}`), nil
				}
				assert.FailNow(t, "unexpected call to GET")
				return nil, errors.New("unexpected call to GET")
			},
		}

		client := createTestServiceUserServiceClient(mockClient)
		stubs, err := client.List(t.Context())
		assert.NoError(t, err)
		assert.Equal(t, 2, callCount)
		assert.Len(t, stubs, 2)
	})

	t.Run("list fails", func(t *testing.T) {
		mockClient := &mockIAMClient{
			GETFunc: func(ctx context.Context, url string, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
				return nil, errors.New("GET failed")
			},
		}

		client := createTestServiceUserServiceClient(mockClient)
		_, err := client.List(t.Context())
		assert.EqualError(t, err, "GET failed")
	})

	t.Run("list empty", func(t *testing.T) {
		mockClient := &mockIAMClient{
			GETFunc: func(ctx context.Context, url string, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
				return []byte(`{
					"count": 0,
					"results": []
				}`), nil
			},
		}

		client := createTestServiceUserServiceClient(mockClient)
		stubs, err := client.List(t.Context())
		assert.NoError(t, err)
		assert.Empty(t, stubs)
	})
}

func TestService_Update(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		putCallCount := 0
		mockClient := &mockIAMClient{
			PUTFunc: func(ctx context.Context, url string, payload any, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
				putCallCount++
				if putCallCount == 1 {
					assert.Equal(t, fmt.Sprintf("%s/iam/v1/accounts/%s/service-users/%s", testEndpointURL, testAccountID, "test-uid"), url)
					return nil, nil
				}

				if putCallCount == 2 {
					assert.Equal(t, fmt.Sprintf("%s/iam/v1/accounts/%s/users/%s/groups", testEndpointURL, testAccountID, "test@example.com"), url)
					return nil, nil
				}

				assert.FailNow(t, "unexpected call to PUT")
				return nil, errors.New("unexpected call to PUT")
			},
		}

		client := createTestServiceUserServiceClient(mockClient)
		serviceUser := &serviceusers.ServiceUser{
			Name:        "Updated User",
			Email:       "test@example.com",
			Description: "Updated description",
			Groups:      []string{"group-1", "group-2"},
		}

		err := client.Update(t.Context(), "test-uid", serviceUser)
		assert.NoError(t, err)
		assert.Equal(t, 2, putCallCount)
	})

	t.Run("update fails on user details", func(t *testing.T) {
		mockClient := &mockIAMClient{
			PUTFunc: func(ctx context.Context, url string, payload any, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
				return nil, errors.New("PUT failed")
			},
		}

		client := createTestServiceUserServiceClient(mockClient)
		serviceUser := &serviceusers.ServiceUser{
			Name:  "Updated User",
			Email: "test@example.com",
		}

		err := client.Update(t.Context(), "test-uid", serviceUser)
		assert.EqualError(t, err, "PUT failed")
	})

	t.Run("update fails on group assignment", func(t *testing.T) {
		putCallCount := 0
		mockClient := &mockIAMClient{
			PUTFunc: func(ctx context.Context, url string, payload any, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
				putCallCount++
				if putCallCount == 1 {
					assert.Equal(t, fmt.Sprintf("%s/iam/v1/accounts/%s/service-users/%s", testEndpointURL, testAccountID, "test-uid"), url)
					return nil, nil
				}
				if putCallCount == 2 {
					assert.Equal(t, fmt.Sprintf("%s/iam/v1/accounts/%s/users/%s/groups", testEndpointURL, testAccountID, "test@example.com"), url)
					return nil, errors.New("group assignment failed")
				}
				assert.FailNow(t, "unexpected call to PUT")
				return nil, errors.New("unexpected call to PUT")
			},
		}

		client := createTestServiceUserServiceClient(mockClient)
		serviceUser := &serviceusers.ServiceUser{
			Name:   "Updated User",
			Email:  "test@example.com",
			Groups: []string{"group-1"},
		}

		err := client.Update(t.Context(), "test-uid", serviceUser)
		assert.EqualError(t, err, "group assignment failed")
	})
}

func TestService_Delete(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		mockClient := &mockIAMClient{
			DELETEFunc: func(ctx context.Context, url string, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
				assert.Equal(t, fmt.Sprintf("%s/iam/v1/accounts/%s/service-users/%s", testEndpointURL, testAccountID, "test-uid"), url)
				return nil, nil
			},
		}

		client := createTestServiceUserServiceClient(mockClient)
		err := client.Delete(t.Context(), "test-uid")
		assert.NoError(t, err)
	})

	t.Run("delete fails", func(t *testing.T) {
		mockClient := &mockIAMClient{
			DELETEFunc: func(ctx context.Context, url string, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
				return nil, errors.New("DELETE failed")
			},
		}

		client := createTestServiceUserServiceClient(mockClient)
		err := client.Delete(t.Context(), "test-uid")
		assert.EqualError(t, err, "DELETE failed")
	})

	t.Run("delete ignores user does not exist error", func(t *testing.T) {
		mockClient := &mockIAMClient{
			DELETEFunc: func(ctx context.Context, url string, expectedResponseCode int, forceNewBearer bool) ([]byte, error) {
				return nil, errors.New("User test-uid does not exist")
			},
		}

		client := createTestServiceUserServiceClient(mockClient)
		err := client.Delete(t.Context(), "test-uid")
		assert.NoError(t, err)
	})
}

func TestService_SchemaID(t *testing.T) {
	client := createTestServiceUserServiceClient(&mockIAMClient{})
	schemaID := client.SchemaID()
	assert.Equal(t, "accounts:iam:serviceusers", schemaID)
}
