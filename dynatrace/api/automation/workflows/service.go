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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"

	cacapi "github.com/dynatrace/dynatrace-configuration-as-code-core/api"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/clients/automation"
	apiClient "github.com/dynatrace/dynatrace-configuration-as-code-core/clients/automation"
)

// AutomationClient defines the subset of the automation client used by this service.
type AutomationClient interface {
	Get(ctx context.Context, resourceType apiClient.ResourceType, id string) (cacapi.Response, error)
	List(ctx context.Context, resourceType apiClient.ResourceType) (cacapi.PagedListResponse, error)
	Create(ctx context.Context, resourceType apiClient.ResourceType, data []byte) (cacapi.Response, error)
	Update(ctx context.Context, resourceType apiClient.ResourceType, id string, data []byte) (cacapi.Response, error)
	Delete(ctx context.Context, resourceType apiClient.ResourceType, id string) (cacapi.Response, error)
}

type service struct {
	credentials  *config.ProviderConfiguration
	clientGetter func(ctx context.Context, credentials *config.ProviderConfiguration) (AutomationClient, error)
}

func Service(credentials *config.ProviderConfiguration) settings.CRUDService[*workflows.Workflow] {
	return &service{credentials: credentials, clientGetter: createCoreClient}
}

func ServiceWithClientGetter(clientGetter func(ctx context.Context, credentials *config.ProviderConfiguration) (AutomationClient, error), credentials *config.ProviderConfiguration) settings.CRUDService[*workflows.Workflow] {
	return &service{credentials: credentials, clientGetter: clientGetter}
}

func createCoreClient(ctx context.Context, credentials *config.ProviderConfiguration) (AutomationClient, error) {
	platformClient, err := tfrest.CreatePlatformClient(ctx, credentials.Platform.EnvironmentURL, credentials.Platform)
	if err != nil {
		return nil, err
	}
	return automation.NewClient(platformClient), nil
}

func (me *service) Get(ctx context.Context, id string, v *workflows.Workflow) error {
	client, err := me.clientGetter(ctx, me.credentials)
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
	// We only remove the position of the API task if the related task in the state and doesn't have positions
	// We keep the permisisons if:
	// - the task is in the state and has a position (we assume that the user set it, so we keep it). Drift detection works.
	// - the task is not in the state
	//   - Either because we removed one or on the API side something was added, and now it needs to be removed (showing with posision). Drift detection doesn't matter much, the task will be removed.
	//   - Or if an import block with generate resource is used (settings.ContextKeyStateConfig is set but empty). Not in state => generate with position. No drift detection, because generate.

	stateConfig, ok := ctx.Value(settings.ContextKeyStateConfig).(*workflows.Workflow)
	if !ok {
		// export doesn't have a state. Keep the position
		return
	}

	positionTaskNames := map[string]bool{}
	for _, st := range stateConfig.Tasks {
		positionTaskNames[st.Name] = st.Position != nil
	}

	for _, task := range v.Tasks {
		if hasPosition, ok := positionTaskNames[task.Name]; ok && !hasPosition {
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
	client, err := me.clientGetter(ctx, me.credentials)
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
	client, err := me.clientGetter(ctx, me.credentials)
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
	client, err := me.clientGetter(ctx, me.credentials)
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
	client, err := me.clientGetter(ctx, me.credentials)
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
