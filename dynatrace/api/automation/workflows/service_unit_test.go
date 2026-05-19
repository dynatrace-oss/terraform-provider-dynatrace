//go:build unit

/*
 * @license
 * Copyright 2026 Dynatrace LLC
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package workflows_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/automation/workflows"
	workflowSettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/automation/workflows/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	coreapi "github.com/dynatrace/dynatrace-configuration-as-code-core/api"
	apiClient "github.com/dynatrace/dynatrace-configuration-as-code-core/clients/automation"
)

// mockAutomationClient implements AutomationClient for testing.
type mockAutomationClient struct {
	getFn    func(ctx context.Context, resourceType apiClient.ResourceType, id string) (coreapi.Response, error)
	listFn   func(ctx context.Context, resourceType apiClient.ResourceType) (coreapi.PagedListResponse, error)
	createFn func(ctx context.Context, resourceType apiClient.ResourceType, data []byte) (coreapi.Response, error)
	updateFn func(ctx context.Context, resourceType apiClient.ResourceType, id string, data []byte) (coreapi.Response, error)
	deleteFn func(ctx context.Context, resourceType apiClient.ResourceType, id string) (coreapi.Response, error)
}

func (m *mockAutomationClient) Get(ctx context.Context, resourceType apiClient.ResourceType, id string) (coreapi.Response, error) {
	return m.getFn(ctx, resourceType, id)
}
func (m *mockAutomationClient) List(ctx context.Context, resourceType apiClient.ResourceType) (coreapi.PagedListResponse, error) {
	return m.listFn(ctx, resourceType)
}
func (m *mockAutomationClient) Create(ctx context.Context, resourceType apiClient.ResourceType, data []byte) (coreapi.Response, error) {
	return m.createFn(ctx, resourceType, data)
}
func (m *mockAutomationClient) Update(ctx context.Context, resourceType apiClient.ResourceType, id string, data []byte) (coreapi.Response, error) {
	return m.updateFn(ctx, resourceType, id, data)
}
func (m *mockAutomationClient) Delete(ctx context.Context, resourceType apiClient.ResourceType, id string) (coreapi.Response, error) {
	return m.deleteFn(ctx, resourceType, id)
}

func mockClientGetter(client *mockAutomationClient) func(ctx context.Context, credentials *rest.Credentials) (workflows.AutomationClient, error) {
	return func(ctx context.Context, credentials *rest.Credentials) (workflows.AutomationClient, error) {
		return client, nil
	}
}

func failingClientGetter(err error) func(ctx context.Context, credentials *rest.Credentials) (workflows.AutomationClient, error) {
	return func(ctx context.Context, credentials *rest.Credentials) (workflows.AutomationClient, error) {
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
		svc := workflows.ServiceWithClientGetter(failingClientGetter(assert.AnError), &rest.Credentials{})
		err := svc.Get(t.Context(), "wf-1", &workflowSettings.Workflow{})
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("Returns error when client Get fails", func(t *testing.T) {
		mock := &mockAutomationClient{
			getFn: func(ctx context.Context, resourceType apiClient.ResourceType, id string) (coreapi.Response, error) {
				return coreapi.Response{}, assert.AnError
			},
		}
		svc := workflows.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		err := svc.Get(t.Context(), "wf-1", &workflowSettings.Workflow{})
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("Returns error on invalid response JSON", func(t *testing.T) {
		mock := &mockAutomationClient{
			getFn: func(ctx context.Context, resourceType apiClient.ResourceType, id string) (coreapi.Response, error) {
				return coreapi.Response{Data: []byte("not-json")}, nil
			},
		}
		svc := workflows.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		err := svc.Get(t.Context(), "wf-1", &workflowSettings.Workflow{})
		var jsonErr *json.SyntaxError
		assert.ErrorAs(t, err, &jsonErr)
	})

	t.Run("Populates Workflow on success and passes Workflows resource type", func(t *testing.T) {
		responseData := []byte(`{"title":"My Workflow","description":"desc","tasks":{"task1":{"action":"namespace:type","input":{},"position":{"x":10,"y":20}}}}`)
		var capturedResourceType apiClient.ResourceType
		var capturedID string
		mock := &mockAutomationClient{
			getFn: func(ctx context.Context, resourceType apiClient.ResourceType, id string) (coreapi.Response, error) {
				capturedResourceType = resourceType
				capturedID = id
				return coreapi.Response{Data: responseData}, nil
			},
		}
		svc := workflows.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		v := &workflowSettings.Workflow{}
		err := svc.Get(t.Context(), "wf-1", v)
		require.NoError(t, err)
		assert.Equal(t, apiClient.Workflows, capturedResourceType)
		assert.Equal(t, "wf-1", capturedID)
		assert.Equal(t, &workflowSettings.Workflow{
			Title:       "My Workflow",
			Description: "desc",
			Tasks: workflowSettings.Tasks{
				&workflowSettings.Task{
					Name:     "task1",
					Position: &workflowSettings.TaskPosition{X: 10, Y: 20},
					Action:   "namespace:type",
					Input:    map[string]any{},
				},
			},
		}, v)
	})

	t.Run("Keeps task positions when no state config is present (export)", func(t *testing.T) {
		responseData := []byte(`{"title":"wf","tasks":{"task1":{"action":"namespace:type","input":{},"position":{"x":1,"y":2}}}}`)
		mock := &mockAutomationClient{
			getFn: func(ctx context.Context, resourceType apiClient.ResourceType, id string) (coreapi.Response, error) {
				return coreapi.Response{Data: responseData}, nil
			},
		}
		svc := workflows.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		v := &workflowSettings.Workflow{}
		// context without a state config (export scenario)
		err := svc.Get(t.Context(), "wf-1", v)
		require.NoError(t, err)
		require.Len(t, v.Tasks, 1)
		assert.Equal(t, &workflowSettings.TaskPosition{X: 1, Y: 2}, v.Tasks[0].Position)
	})

	t.Run("Keeps task positions when matching state task has a position", func(t *testing.T) {
		responseData := []byte(`{"title":"wf","tasks":{"task1":{"action":"namespace:type","input":{},"position":{"x":5,"y":6}}}}`)
		mock := &mockAutomationClient{
			getFn: func(ctx context.Context, resourceType apiClient.ResourceType, id string) (coreapi.Response, error) {
				return coreapi.Response{Data: responseData}, nil
			},
		}
		svc := workflows.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		stateConfig := &workflowSettings.Workflow{
			Tasks: workflowSettings.Tasks{
				{Name: "task1", Position: &workflowSettings.TaskPosition{X: 5, Y: 6}},
			},
		}
		ctx := context.WithValue(t.Context(), settings.ContextKeyStateConfig, stateConfig)
		v := &workflowSettings.Workflow{}
		err := svc.Get(ctx, "wf-1", v)
		require.NoError(t, err)
		require.Len(t, v.Tasks, 1)
		assert.Equal(t, &workflowSettings.TaskPosition{X: 5, Y: 6}, v.Tasks[0].Position)
	})

	t.Run("Keeps task positions if the API response task does not match a local one", func(t *testing.T) {
		responseData := []byte(`{"title":"wf","tasks":{"task1":{"action":"namespace:type","input":{},"position":{"x":1,"y":2}}, "task2":{"action":"namespace:type","input":{},"position":{"x":5,"y":6}}}}`)
		mock := &mockAutomationClient{
			getFn: func(ctx context.Context, resourceType apiClient.ResourceType, id string) (coreapi.Response, error) {
				return coreapi.Response{Data: responseData}, nil
			},
		}
		svc := workflows.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		stateConfig := &workflowSettings.Workflow{
			Tasks: workflowSettings.Tasks{
				{Name: "task1", Position: &workflowSettings.TaskPosition{X: 5, Y: 6}},
			},
		}
		ctx := context.WithValue(t.Context(), settings.ContextKeyStateConfig, stateConfig)
		v := &workflowSettings.Workflow{}
		err := svc.Get(ctx, "wf-1", v)
		require.NoError(t, err)
		require.Len(t, v.Tasks, 2)

		var addedTask *workflowSettings.Task
		// maps Response from map to list, therefore we should look for it as the array order is not guaranteed
		for _, task := range v.Tasks {
			if task.Name == "task2" {
				addedTask = task
				break
			}
		}
		require.NotNil(t, addedTask)
		// not in state: keep position (e.g., import block with generate)
		assert.Equal(t, addedTask.Name, "task2")
		assert.Equal(t, addedTask.Position, &workflowSettings.TaskPosition{X: 5, Y: 6})
	})

	t.Run("Removes task position when matching state task does not have a position", func(t *testing.T) {
		responseData := []byte(`{"title":"wf","tasks":{"task1":{"action":"namespace:type","input":{},"position":{"x":5,"y":6}}}}`)
		mock := &mockAutomationClient{
			getFn: func(ctx context.Context, resourceType apiClient.ResourceType, id string) (coreapi.Response, error) {
				return coreapi.Response{Data: responseData}, nil
			},
		}
		svc := workflows.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		stateConfig := &workflowSettings.Workflow{
			Tasks: workflowSettings.Tasks{
				{Name: "task1"},
			},
		}
		ctx := context.WithValue(t.Context(), settings.ContextKeyStateConfig, stateConfig)
		v := &workflowSettings.Workflow{}
		err := svc.Get(ctx, "wf-1", v)
		require.NoError(t, err)
		require.Len(t, v.Tasks, 1)
		assert.Nil(t, v.Tasks[0].Position)
	})
}

func TestService_List(t *testing.T) {
	t.Run("Returns error when client creation fails", func(t *testing.T) {
		svc := workflows.ServiceWithClientGetter(failingClientGetter(assert.AnError), &rest.Credentials{})
		_, err := svc.List(t.Context())
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("Returns error when client List fails", func(t *testing.T) {
		mock := &mockAutomationClient{
			listFn: func(ctx context.Context, resourceType apiClient.ResourceType) (coreapi.PagedListResponse, error) {
				return nil, assert.AnError
			},
		}
		svc := workflows.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		_, err := svc.List(t.Context())
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("Returns error on invalid workflow stub JSON", func(t *testing.T) {
		mock := &mockAutomationClient{
			listFn: func(ctx context.Context, resourceType apiClient.ResourceType) (coreapi.PagedListResponse, error) {
				return pagedResponse([]byte("bad-json")), nil
			},
		}
		svc := workflows.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		_, err := svc.List(t.Context())
		var jsonErr *json.SyntaxError
		assert.ErrorAs(t, err, &jsonErr)
	})

	t.Run("Returns stubs on success and passes Workflows resource type", func(t *testing.T) {
		wf1JSON, _ := json.Marshal(workflows.WorkflowStub{ID: "wf-1", Title: "Workflow 1"})
		wf2JSON, _ := json.Marshal(workflows.WorkflowStub{ID: "wf-2", Title: "Workflow 2"})
		var capturedResourceType apiClient.ResourceType
		mock := &mockAutomationClient{
			listFn: func(ctx context.Context, resourceType apiClient.ResourceType) (coreapi.PagedListResponse, error) {
				capturedResourceType = resourceType
				return pagedResponse(wf1JSON, wf2JSON), nil
			},
		}
		svc := workflows.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		stubs, err := svc.List(t.Context())
		require.NoError(t, err)
		assert.Equal(t, apiClient.Workflows, capturedResourceType)
		assert.ElementsMatch(t, []*api.Stub{{ID: "wf-1", Name: "Workflow 1"}, {ID: "wf-2", Name: "Workflow 2"}}, stubs)
	})

	t.Run("Returns empty stubs when no workflows exist", func(t *testing.T) {
		mock := &mockAutomationClient{
			listFn: func(ctx context.Context, resourceType apiClient.ResourceType) (coreapi.PagedListResponse, error) {
				return pagedResponse(), nil
			},
		}
		svc := workflows.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		stubs, err := svc.List(t.Context())
		require.NoError(t, err)
		assert.Empty(t, stubs)
	})
}

func TestService_Create(t *testing.T) {
	t.Run("Returns error when client creation fails", func(t *testing.T) {
		svc := workflows.ServiceWithClientGetter(failingClientGetter(assert.AnError), &rest.Credentials{})
		_, err := svc.Create(t.Context(), &workflowSettings.Workflow{Title: "wf"})
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("Returns error when client Create fails", func(t *testing.T) {
		mock := &mockAutomationClient{
			createFn: func(ctx context.Context, resourceType apiClient.ResourceType, data []byte) (coreapi.Response, error) {
				return coreapi.Response{}, assert.AnError
			},
		}
		svc := workflows.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		_, err := svc.Create(t.Context(), &workflowSettings.Workflow{Title: "wf"})
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("Returns error on invalid response JSON", func(t *testing.T) {
		mock := &mockAutomationClient{
			createFn: func(ctx context.Context, resourceType apiClient.ResourceType, data []byte) (coreapi.Response, error) {
				return coreapi.Response{Data: []byte("bad-json")}, nil
			},
		}
		svc := workflows.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		_, err := svc.Create(t.Context(), &workflowSettings.Workflow{Title: "wf"})
		var jsonErr *json.SyntaxError
		assert.ErrorAs(t, err, &jsonErr)
	})

	t.Run("Sends correct payload and returns stub on success", func(t *testing.T) {
		var capturedResourceType apiClient.ResourceType
		var capturedPayload map[string]any
		mock := &mockAutomationClient{
			createFn: func(ctx context.Context, resourceType apiClient.ResourceType, data []byte) (coreapi.Response, error) {
				capturedResourceType = resourceType
				err := json.Unmarshal(data, &capturedPayload)
				require.NoError(t, err)
				respJSON, _ := json.Marshal(workflows.WorkflowStub{ID: "new-wf-id"})
				return coreapi.Response{Data: respJSON}, nil
			},
		}
		svc := workflows.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		stub, err := svc.Create(t.Context(), &workflowSettings.Workflow{
			Title:       "My Workflow",
			Description: "desc",
		})
		require.NoError(t, err)
		assert.Equal(t, apiClient.Workflows, capturedResourceType)
		assert.Equal(t, "My Workflow", capturedPayload["title"])
		assert.Equal(t, "desc", capturedPayload["description"])
		assert.Equal(t, &api.Stub{ID: "new-wf-id", Name: "My Workflow"}, stub)
	})
}

func TestService_Update(t *testing.T) {
	t.Run("Returns error when client creation fails", func(t *testing.T) {
		svc := workflows.ServiceWithClientGetter(failingClientGetter(assert.AnError), &rest.Credentials{})
		err := svc.Update(t.Context(), "wf-1", &workflowSettings.Workflow{})
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("Returns error when client Update fails", func(t *testing.T) {
		mock := &mockAutomationClient{
			updateFn: func(ctx context.Context, resourceType apiClient.ResourceType, id string, data []byte) (coreapi.Response, error) {
				return coreapi.Response{}, assert.AnError
			},
		}
		svc := workflows.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		err := svc.Update(t.Context(), "wf-1", &workflowSettings.Workflow{})
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("Sends correct id, resource type and payload on success", func(t *testing.T) {
		var capturedResourceType apiClient.ResourceType
		var capturedID string
		var capturedPayload map[string]any
		mock := &mockAutomationClient{
			updateFn: func(ctx context.Context, resourceType apiClient.ResourceType, id string, data []byte) (coreapi.Response, error) {
				capturedResourceType = resourceType
				capturedID = id
				err := json.Unmarshal(data, &capturedPayload)
				require.NoError(t, err)
				return coreapi.Response{}, nil
			},
		}
		svc := workflows.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		err := svc.Update(t.Context(), "wf-1", &workflowSettings.Workflow{
			Title:       "Updated",
			Description: "updated desc",
		})
		require.NoError(t, err)
		assert.Equal(t, apiClient.Workflows, capturedResourceType)
		assert.Equal(t, "wf-1", capturedID)
		assert.Equal(t, "Updated", capturedPayload["title"])
		assert.Equal(t, "updated desc", capturedPayload["description"])
	})
}

func TestService_Delete(t *testing.T) {
	t.Run("Returns error when client creation fails", func(t *testing.T) {
		svc := workflows.ServiceWithClientGetter(failingClientGetter(assert.AnError), &rest.Credentials{})
		err := svc.Delete(t.Context(), "wf-1")
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("Returns error when client Delete fails with non-NotFound error", func(t *testing.T) {
		mock := &mockAutomationClient{
			deleteFn: func(ctx context.Context, resourceType apiClient.ResourceType, id string) (coreapi.Response, error) {
				return coreapi.Response{}, assert.AnError
			},
		}
		svc := workflows.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		err := svc.Delete(t.Context(), "wf-1")
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("Returns nil when client Delete fails with NotFound error", func(t *testing.T) {
		mock := &mockAutomationClient{
			deleteFn: func(ctx context.Context, resourceType apiClient.ResourceType, id string) (coreapi.Response, error) {
				return coreapi.Response{}, coreapi.APIError{StatusCode: http.StatusNotFound}
			},
		}
		svc := workflows.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		err := svc.Delete(t.Context(), "wf-1")
		assert.NoError(t, err)
	})

	t.Run("Returns nil and passes correct id and resource type on success", func(t *testing.T) {
		var capturedResourceType apiClient.ResourceType
		var capturedID string
		mock := &mockAutomationClient{
			deleteFn: func(ctx context.Context, resourceType apiClient.ResourceType, id string) (coreapi.Response, error) {
				capturedResourceType = resourceType
				capturedID = id
				return coreapi.Response{}, nil
			},
		}
		svc := workflows.ServiceWithClientGetter(mockClientGetter(mock), &rest.Credentials{})
		err := svc.Delete(t.Context(), "wf-1")
		require.NoError(t, err)
		assert.Equal(t, apiClient.Workflows, capturedResourceType)
		assert.Equal(t, "wf-1", capturedID)
	})
}
