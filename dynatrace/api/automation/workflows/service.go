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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	workflows "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/automation/workflows/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/monaco/pkg/client/auth"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/monaco/pkg/client/automation"
)

func Service(credentials *settings.Credentials) settings.CRUDService[*workflows.Workflow] {
	return &service{credentials}
}

type service struct {
	credentials *settings.Credentials
}

func (me *service) client() *automation.Client {
	httpClient := auth.NewOAuthClient(context.TODO(), auth.OauthCredentials{
		ClientID:     me.credentials.Automation.ClientID,
		ClientSecret: me.credentials.Automation.ClientSecret,
		TokenURL:     me.credentials.Automation.TokenURL,
	})
	return automation.NewClient(me.credentials.Automation.EnvironmentURL, httpClient)
}

func (me *service) Get(id string, v *workflows.Workflow) (err error) {
	var result *automation.Response
	if result, err = me.client().GET(automation.Workflows, id); err != nil {
		return err
	}
	return json.Unmarshal(result.Data, &v)
}

func (me *service) SchemaID() string {
	return "automation:workflows"
}

func (me *service) List() (api.Stubs, error) {
	result, err := me.client().LIST(automation.Workflows)
	if err != nil {
		return nil, err
	}
	var stubs api.Stubs
	for _, r := range result {
		var workflow workflows.Workflow
		if err := json.Unmarshal(r.Data, &workflow); err != nil {
			return nil, err
		}
		stubs = append(stubs, &api.Stub{ID: r.ID, Name: workflow.Title})
	}
	return stubs, nil
}

func (me *service) Validate(v *workflows.Workflow) error {
	return nil // no endpoint for that
}

func (me *service) Create(v *workflows.Workflow) (stub *api.Stub, err error) {
	var data []byte
	var id string
	if data, err = json.Marshal(v); err != nil {
		return nil, err
	}
	if id, err = me.client().INSERT(automation.Workflows, data); err != nil {
		return nil, err
	}
	return &api.Stub{Name: v.Title, ID: id}, nil
}

func (me *service) Update(id string, v *workflows.Workflow) (err error) {
	var data []byte
	if data, err = json.Marshal(v); err != nil {
		return err
	}
	return me.client().UPDATE(automation.Workflows, id, data)
}

func (me *service) Delete(id string) error {
	return me.client().DELETE(automation.Workflows, id)
}

func (me *service) New() *workflows.Workflow {
	return new(workflows.Workflow)
}

func (me *service) Name() string {
	return me.SchemaID()
}
