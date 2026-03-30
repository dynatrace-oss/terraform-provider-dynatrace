/**
* @license
* Copyright 2020 Dynatrace LLC
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

package workflows

import (
	"context"
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	workflows "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/automation/workflows/settings"
	tfrest "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	cacapi "github.com/dynatrace/dynatrace-configuration-as-code-core/api"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/clients/automation"
	apiClient "github.com/dynatrace/dynatrace-configuration-as-code-core/clients/automation"
)

func Service(credentials *tfrest.Credentials) settings.CRUDService[*workflows.Workflow] {
	return &service{credentials}
}

type service struct {
	credentials *tfrest.Credentials
}

func (me *service) client(ctx context.Context) (*automation.Client, error) {
	platformClient, err := tfrest.CreatePlatformClient(ctx, me.credentials.OAuth.EnvironmentURL, me.credentials)
	if err != nil {
		return nil, err
	}
	return automation.NewClient(platformClient), nil
}

func (me *service) Get(ctx context.Context, id string, v *workflows.Workflow) error {
	client, err := me.client(ctx)
	if err != nil {
		return err
	}
	var response cacapi.Response
	if response, err = client.Get(ctx, apiClient.Workflows, id); err != nil {
		return err
	}

	if err := json.Unmarshal(response.Data, &v); err != nil {
		return err
	}
	deComputeTaskPosition(ctx, v)

	return nil
}

// deComputeTaskPosition removes API set defaults of task positions.
// The reason here is simply because the TypeSet hash changes every time if a field of an item is set to computed even if nothing changed.
// If the position is set by the user in the Terraform file, then it's correctly compared to the API set one.
// DiffSuppress can't help here because even if the task doesn't change, the bug of TypeSet happens, which adds a new empty item which leads to a non-empty plan
func deComputeTaskPosition(ctx context.Context, v *workflows.Workflow) {
	positionTaskNames := map[string]bool{}
	if stateConfig, ok := ctx.Value(settings.ContextKeyStateConfig).(*workflows.Workflow); ok {
		for _, st := range stateConfig.Tasks {
			if st.Position != nil {
				positionTaskNames[st.Name] = true
			}
		}
	}
	for _, task := range v.Tasks {
		if !positionTaskNames[task.Name] {
			// => if not in state: set position (returned from API) to nil
			task.Position = nil
		}
	}
}

func (me *service) SchemaID() string {
	return "automation:workflows"
}

type WorkflowStub struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	client, err := me.client(ctx)
	if err != nil {
		return nil, err
	}
	listResponse, err := client.List(ctx, apiClient.Workflows)
	if err != nil {
		return nil, err
	}

	var stubs api.Stubs
	for _, r := range listResponse.All() {
		var workflowStub WorkflowStub
		if err := json.Unmarshal(r, &workflowStub); err != nil {
			return nil, err
		}
		stubs = append(stubs, &api.Stub{ID: workflowStub.ID, Name: workflowStub.Title})
	}
	return stubs, nil
}

func (me *service) Validate(v *workflows.Workflow) error {
	return nil // no endpoint for that
}

func (me *service) Create(ctx context.Context, v *workflows.Workflow) (stub *api.Stub, err error) {
	client, err := me.client(ctx)
	if err != nil {
		return nil, err
	}
	var data []byte
	if data, err = json.Marshal(v); err != nil {
		return nil, err
	}
	var response cacapi.Response
	if response, err = client.Create(ctx, apiClient.Workflows, data); err != nil {
		return nil, err
	}

	var workflowStub WorkflowStub
	if err := json.Unmarshal(response.Data, &workflowStub); err != nil {
		return nil, err
	}
	return &api.Stub{Name: v.Title, ID: workflowStub.ID}, nil
}

func (me *service) Update(ctx context.Context, id string, v *workflows.Workflow) (err error) {
	client, err := me.client(ctx)
	if err != nil {
		return err
	}
	var data []byte
	if data, err = json.Marshal(v); err != nil {
		return err
	}
	_, err = client.Update(ctx, apiClient.Workflows, id, data)
	return err
}

func (me *service) Delete(ctx context.Context, id string) error {
	client, err := me.client(ctx)
	if err != nil {
		return err
	}
	if _, err = client.Delete(ctx, apiClient.Workflows, id); err != nil && !cacapi.IsNotFoundError(err) {
		return err
	}

	return nil
}

func (me *service) New() *workflows.Workflow {
	return new(workflows.Workflow)
}
