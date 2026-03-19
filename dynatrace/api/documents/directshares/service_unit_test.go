//go:build unit

/**
* @license
* Copyright 2026 Dynatrace LLC
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

package directshares

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	coreapi "github.com/dynatrace/dynatrace-configuration-as-code-core/api"

	serviceSettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/documents/directshares/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockDirectSharesClient struct {
	listFn             func(ctx context.Context) (coreapi.PagedListResponse, error)
	getFn              func(ctx context.Context, id string) (coreapi.Response, error)
	createFn           func(ctx context.Context, data []byte) (coreapi.Response, error)
	deleteFn           func(ctx context.Context, id string) error
	getRecipientsFn    func(ctx context.Context, id string) (coreapi.PagedListResponse, error)
	addRecipientsFn    func(ctx context.Context, id string, data []byte) error
	removeRecipientsFn func(ctx context.Context, id string, data []byte) error
}

func (m *mockDirectSharesClient) List(ctx context.Context) (coreapi.PagedListResponse, error) {
	return m.listFn(ctx)
}
func (m *mockDirectSharesClient) Get(ctx context.Context, id string) (coreapi.Response, error) {
	return m.getFn(ctx, id)
}
func (m *mockDirectSharesClient) Create(ctx context.Context, data []byte) (coreapi.Response, error) {
	return m.createFn(ctx, data)
}
func (m *mockDirectSharesClient) Delete(ctx context.Context, id string) error {
	return m.deleteFn(ctx, id)
}
func (m *mockDirectSharesClient) GetRecipients(ctx context.Context, id string) (coreapi.PagedListResponse, error) {
	return m.getRecipientsFn(ctx, id)
}
func (m *mockDirectSharesClient) AddRecipients(ctx context.Context, id string, data []byte) error {
	return m.addRecipientsFn(ctx, id, data)
}
func (m *mockDirectSharesClient) RemoveRecipients(ctx context.Context, id string, data []byte) error {
	return m.removeRecipientsFn(ctx, id, data)
}

func mockClientGetter(client *mockDirectSharesClient) func(ctx context.Context, credentials *rest.Credentials) (directSharesClient, error) {
	return func(ctx context.Context, credentials *rest.Credentials) (directSharesClient, error) {
		return client, nil
	}
}

func failingClientGetter(err error) func(ctx context.Context, credentials *rest.Credentials) (directSharesClient, error) {
	return func(ctx context.Context, credentials *rest.Credentials) (directSharesClient, error) {
		return nil, err
	}
}

func pagedResponse(objects ...[]byte) coreapi.PagedListResponse {
	return coreapi.PagedListResponse{
		{Objects: objects},
	}
}

func TestService_Get(t *testing.T) {
	t.Run("Returns error when client creation fails", func(t *testing.T) {
		clientErr := errors.New("client creation failed")
		svc := ServiceWithClientGetter(failingClientGetter(clientErr), &rest.Credentials{})
		v := &serviceSettings.DirectShare{}
		err := svc.Get(t.Context(), "id-1", v)
		assert.ErrorIs(t, err, clientErr)
	})

	t.Run("Returns error when Get call fails", func(t *testing.T) {
		getErr := errors.New("get failed")
		mock := &mockDirectSharesClient{
			getFn: func(ctx context.Context, id string) (coreapi.Response, error) {
				return coreapi.Response{}, getErr
			},
		}
		svc := ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		v := &serviceSettings.DirectShare{}
		err := svc.Get(t.Context(), "id-1", v)
		assert.ErrorIs(t, err, getErr)
	})

	t.Run("Returns error on invalid Get response JSON", func(t *testing.T) {
		mock := &mockDirectSharesClient{
			getFn: func(ctx context.Context, id string) (coreapi.Response, error) {
				return coreapi.Response{Data: []byte("not-json")}, nil
			},
		}
		svc := ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		v := &serviceSettings.DirectShare{}
		err := svc.Get(t.Context(), "id-1", v)
		assert.Error(t, err)
	})

	t.Run("Returns error when GetRecipients fails", func(t *testing.T) {
		recipientsErr := errors.New("get recipients failed")
		dsJSON, _ := json.Marshal(directShareDTO{ID: "id-1", DocumentId: "doc-1", Access: []string{"read"}})
		mock := &mockDirectSharesClient{
			getFn: func(ctx context.Context, id string) (coreapi.Response, error) {
				return coreapi.Response{Data: dsJSON}, nil
			},
			getRecipientsFn: func(ctx context.Context, id string) (coreapi.PagedListResponse, error) {
				return nil, recipientsErr
			},
		}
		svc := ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		v := &serviceSettings.DirectShare{}
		err := svc.Get(t.Context(), "id-1", v)
		assert.ErrorIs(t, err, recipientsErr)
	})

	t.Run("Returns error on invalid recipient JSON", func(t *testing.T) {
		dsJSON, _ := json.Marshal(directShareDTO{ID: "id-1", DocumentId: "doc-1", Access: []string{"read"}})
		mock := &mockDirectSharesClient{
			getFn: func(ctx context.Context, id string) (coreapi.Response, error) {
				return coreapi.Response{Data: dsJSON}, nil
			},
			getRecipientsFn: func(ctx context.Context, id string) (coreapi.PagedListResponse, error) {
				return pagedResponse([]byte("bad-json")), nil
			},
		}
		svc := ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		v := &serviceSettings.DirectShare{}
		err := svc.Get(t.Context(), "id-1", v)
		assert.Error(t, err)
	})

	t.Run("Populates DirectShare on success", func(t *testing.T) {
		dsJSON, _ := json.Marshal(directShareDTO{ID: "id-1", DocumentId: "doc-1", Access: []string{"read", "write"}})
		recipientJSON, _ := json.Marshal(recipientDTO{ID: "user-1", Type: "user"})
		mock := &mockDirectSharesClient{
			getFn: func(ctx context.Context, id string) (coreapi.Response, error) {
				assert.Equal(t, "id-1", id)
				return coreapi.Response{Data: dsJSON}, nil
			},
			getRecipientsFn: func(ctx context.Context, id string) (coreapi.PagedListResponse, error) {
				assert.Equal(t, "id-1", id)
				return pagedResponse(recipientJSON), nil
			},
		}
		svc := ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		v := &serviceSettings.DirectShare{}
		err := svc.Get(t.Context(), "id-1", v)
		require.NoError(t, err)
		assert.Equal(t, "id-1", v.ID)
		assert.Equal(t, "doc-1", v.DocumentId)
		assert.Equal(t, "read-write", v.Access)
		require.Len(t, v.Recipients, 1)
		assert.Equal(t, "user-1", v.Recipients[0].ID)
		assert.Equal(t, "user", v.Recipients[0].Type)
	})
}

func TestService_List(t *testing.T) {
	t.Run("Returns error when client creation fails", func(t *testing.T) {
		clientErr := errors.New("client creation failed")
		svc := ServiceWithClientGetter(failingClientGetter(clientErr), &rest.Credentials{})
		_, err := svc.List(t.Context())
		assert.ErrorIs(t, err, clientErr)
	})

	t.Run("Returns error when List call fails", func(t *testing.T) {
		listErr := errors.New("list failed")
		mock := &mockDirectSharesClient{
			listFn: func(ctx context.Context) (coreapi.PagedListResponse, error) {
				return nil, listErr
			},
		}
		svc := ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		_, err := svc.List(t.Context())
		assert.ErrorIs(t, err, listErr)
	})

	t.Run("Returns error on invalid JSON", func(t *testing.T) {
		mock := &mockDirectSharesClient{
			listFn: func(ctx context.Context) (coreapi.PagedListResponse, error) {
				return pagedResponse([]byte("bad-json")), nil
			},
		}
		svc := ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		_, err := svc.List(t.Context())
		assert.Error(t, err)
	})

	t.Run("Returns stubs on success", func(t *testing.T) {
		ds1, _ := json.Marshal(directShareDTO{ID: "id-1", DocumentId: "doc-1"})
		ds2, _ := json.Marshal(directShareDTO{ID: "id-2", DocumentId: "doc-2"})
		mock := &mockDirectSharesClient{
			listFn: func(ctx context.Context) (coreapi.PagedListResponse, error) {
				return pagedResponse(ds1, ds2), nil
			},
		}
		svc := ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		stubs, err := svc.List(t.Context())
		require.NoError(t, err)
		require.Len(t, stubs, 2)
		assert.Equal(t, "id-1", stubs[0].ID)
		assert.Equal(t, "doc-1", stubs[0].Name)
		assert.Equal(t, "id-2", stubs[1].ID)
		assert.Equal(t, "doc-2", stubs[1].Name)
	})

	t.Run("Returns nil stubs for empty list", func(t *testing.T) {
		mock := &mockDirectSharesClient{
			listFn: func(ctx context.Context) (coreapi.PagedListResponse, error) {
				return pagedResponse(), nil
			},
		}
		svc := ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		stubs, err := svc.List(t.Context())
		require.NoError(t, err)
		assert.Nil(t, stubs)
	})
}

func TestService_Create(t *testing.T) {
	t.Run("Returns error when client creation fails", func(t *testing.T) {
		clientErr := errors.New("client creation failed")
		svc := ServiceWithClientGetter(failingClientGetter(clientErr), &rest.Credentials{})
		_, err := svc.Create(t.Context(), &serviceSettings.DirectShare{})
		assert.ErrorIs(t, err, clientErr)
	})

	t.Run("Returns error when Create call fails", func(t *testing.T) {
		createErr := errors.New("create failed")
		mock := &mockDirectSharesClient{
			createFn: func(ctx context.Context, data []byte) (coreapi.Response, error) {
				return coreapi.Response{}, createErr
			},
		}
		svc := ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		_, err := svc.Create(t.Context(), &serviceSettings.DirectShare{
			DocumentId: "doc-1",
			Access:     "read",
		})
		assert.ErrorIs(t, err, createErr)
	})

	t.Run("Returns error on invalid response JSON", func(t *testing.T) {
		mock := &mockDirectSharesClient{
			createFn: func(ctx context.Context, data []byte) (coreapi.Response, error) {
				return coreapi.Response{Data: []byte("bad-json")}, nil
			},
		}
		svc := ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		_, err := svc.Create(t.Context(), &serviceSettings.DirectShare{
			DocumentId: "doc-1",
			Access:     "read",
		})
		assert.Error(t, err)
	})

	t.Run("Sends correct payload and returns stub on success", func(t *testing.T) {
		var capturedPayload createDirectShareDTO
		mock := &mockDirectSharesClient{
			createFn: func(ctx context.Context, data []byte) (coreapi.Response, error) {
				err := json.Unmarshal(data, &capturedPayload)
				require.NoError(t, err)
				respJSON, _ := json.Marshal(map[string]string{"id": "new-id"})
				return coreapi.Response{Data: respJSON}, nil
			},
		}
		svc := ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		stub, err := svc.Create(t.Context(), &serviceSettings.DirectShare{
			DocumentId: "doc-1",
			Access:     "read",
			Recipients: serviceSettings.Recipients{
				{ID: "user-1", Type: "user"},
				{ID: "group-1", Type: "group"},
			},
		})
		require.NoError(t, err)
		assert.Equal(t, "new-id", stub.ID)

		assert.Equal(t, "doc-1", capturedPayload.DocumentId)
		assert.Equal(t, "read", capturedPayload.Access)
		require.Len(t, capturedPayload.Recipients, 2)
		assert.Equal(t, "user-1", capturedPayload.Recipients[0].ID)
		assert.Equal(t, "user", capturedPayload.Recipients[0].Type)
		assert.Equal(t, "group-1", capturedPayload.Recipients[1].ID)
		assert.Equal(t, "group", capturedPayload.Recipients[1].Type)
	})
}

func TestService_Update(t *testing.T) {
	t.Run("Returns error when client creation fails", func(t *testing.T) {
		clientErr := errors.New("client creation failed")
		svc := ServiceWithClientGetter(failingClientGetter(clientErr), &rest.Credentials{})
		err := svc.Update(t.Context(), "id-1", &serviceSettings.DirectShare{})
		assert.ErrorIs(t, err, clientErr)
	})

	t.Run("Returns error when GetRecipients fails", func(t *testing.T) {
		recipientsErr := errors.New("get recipients failed")
		mock := &mockDirectSharesClient{
			getRecipientsFn: func(ctx context.Context, id string) (coreapi.PagedListResponse, error) {
				return nil, recipientsErr
			},
		}
		svc := ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		err := svc.Update(t.Context(), "id-1", &serviceSettings.DirectShare{})
		assert.ErrorIs(t, err, recipientsErr)
	})

	t.Run("Adds new recipients and removes stale ones", func(t *testing.T) {
		remoteRecipient1, _ := json.Marshal(recipientDTO{ID: "existing-user", Type: "user"})
		remoteRecipient2, _ := json.Marshal(recipientDTO{ID: "stale-user", Type: "user"})

		var addedPayload addDirectShareRecipientsDTO
		var removedPayload removeDirectShareRecipientsDTO
		mock := &mockDirectSharesClient{
			getRecipientsFn: func(ctx context.Context, id string) (coreapi.PagedListResponse, error) {
				return pagedResponse(remoteRecipient1, remoteRecipient2), nil
			},
			addRecipientsFn: func(ctx context.Context, id string, data []byte) error {
				return json.Unmarshal(data, &addedPayload)
			},
			removeRecipientsFn: func(ctx context.Context, id string, data []byte) error {
				return json.Unmarshal(data, &removedPayload)
			},
		}
		svc := ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		err := svc.Update(t.Context(), "id-1", &serviceSettings.DirectShare{
			Recipients: serviceSettings.Recipients{
				{ID: "existing-user", Type: "user"},
				{ID: "new-user", Type: "user"},
			},
		})
		require.NoError(t, err)

		require.Len(t, addedPayload.Recipients, 1)
		assert.Equal(t, "new-user", addedPayload.Recipients[0].ID)

		require.Len(t, removedPayload.Ids, 1)
		assert.Equal(t, "stale-user", removedPayload.Ids[0])
	})

	t.Run("Skips add when no new recipients", func(t *testing.T) {
		remoteRecipient, _ := json.Marshal(recipientDTO{ID: "user-1", Type: "user"})
		addCalled := false
		mock := &mockDirectSharesClient{
			getRecipientsFn: func(ctx context.Context, id string) (coreapi.PagedListResponse, error) {
				return pagedResponse(remoteRecipient), nil
			},
			addRecipientsFn: func(ctx context.Context, id string, data []byte) error {
				addCalled = true
				return nil
			},
			removeRecipientsFn: func(ctx context.Context, id string, data []byte) error {
				return nil
			},
		}
		svc := ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		err := svc.Update(t.Context(), "id-1", &serviceSettings.DirectShare{
			Recipients: serviceSettings.Recipients{
				{ID: "user-1", Type: "user"},
			},
		})
		require.NoError(t, err)
		assert.False(t, addCalled, "AddRecipients should not be called when no new recipients")
	})

	t.Run("Skips remove when no stale recipients", func(t *testing.T) {
		removeCalled := false
		mock := &mockDirectSharesClient{
			getRecipientsFn: func(ctx context.Context, id string) (coreapi.PagedListResponse, error) {
				return pagedResponse(), nil
			},
			addRecipientsFn: func(ctx context.Context, id string, data []byte) error {
				return nil
			},
			removeRecipientsFn: func(ctx context.Context, id string, data []byte) error {
				removeCalled = true
				return nil
			},
		}
		svc := ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		err := svc.Update(t.Context(), "id-1", &serviceSettings.DirectShare{
			Recipients: serviceSettings.Recipients{
				{ID: "new-user", Type: "user"},
			},
		})
		require.NoError(t, err)
		assert.False(t, removeCalled, "RemoveRecipients should not be called when no stale recipients")
	})

	t.Run("Returns error when AddRecipients fails", func(t *testing.T) {
		addErr := errors.New("add recipients failed")
		mock := &mockDirectSharesClient{
			getRecipientsFn: func(ctx context.Context, id string) (coreapi.PagedListResponse, error) {
				return pagedResponse(), nil
			},
			addRecipientsFn: func(ctx context.Context, id string, data []byte) error {
				return addErr
			},
		}
		svc := ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		err := svc.Update(t.Context(), "id-1", &serviceSettings.DirectShare{
			Recipients: serviceSettings.Recipients{
				{ID: "user-1", Type: "user"},
			},
		})
		assert.ErrorIs(t, err, addErr)
	})

	t.Run("Returns error when RemoveRecipients fails", func(t *testing.T) {
		remoteRecipient, _ := json.Marshal(recipientDTO{ID: "stale-user", Type: "user"})
		removeErr := errors.New("remove recipients failed")
		mock := &mockDirectSharesClient{
			getRecipientsFn: func(ctx context.Context, id string) (coreapi.PagedListResponse, error) {
				return pagedResponse(remoteRecipient), nil
			},
			addRecipientsFn: func(ctx context.Context, id string, data []byte) error {
				return nil
			},
			removeRecipientsFn: func(ctx context.Context, id string, data []byte) error {
				return removeErr
			},
		}
		svc := ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		err := svc.Update(t.Context(), "id-1", &serviceSettings.DirectShare{})
		assert.ErrorIs(t, err, removeErr)
	})
}

func TestService_Delete(t *testing.T) {
	t.Run("Returns error when client creation fails", func(t *testing.T) {
		clientErr := errors.New("client creation failed")
		svc := ServiceWithClientGetter(failingClientGetter(clientErr), &rest.Credentials{})
		err := svc.Delete(t.Context(), "id-1")
		assert.ErrorIs(t, err, clientErr)
	})

	t.Run("Returns error when Delete call fails", func(t *testing.T) {
		deleteErr := errors.New("delete failed")
		mock := &mockDirectSharesClient{
			deleteFn: func(ctx context.Context, id string) error {
				return deleteErr
			},
		}
		svc := ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		err := svc.Delete(t.Context(), "id-1")
		assert.ErrorIs(t, err, deleteErr)
	})

	t.Run("Returns nil on success", func(t *testing.T) {
		var capturedID string
		mock := &mockDirectSharesClient{
			deleteFn: func(ctx context.Context, id string) error {
				capturedID = id
				return nil
			},
		}
		svc := ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		err := svc.Delete(t.Context(), "id-1")
		require.NoError(t, err)
		assert.Equal(t, "id-1", capturedID)
	})
}
