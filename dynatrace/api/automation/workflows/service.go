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
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	tfrest "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	automationerr "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/automation"
	workflows "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/automation/workflows/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/monaco/pkg/client/auth"
	apiClient "github.com/dynatrace/dynatrace-configuration-as-code-core/api/clients/automation"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/clients/automation"
)

func Service(credentials *settings.Credentials) settings.CRUDService[*workflows.Workflow] {
	return &service{credentials}
}

type service struct {
	credentials *settings.Credentials
}

type MyRoundTripper struct {
	RoundTripper http.RoundTripper
}

var lock sync.Mutex

func (rt *MyRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	lock.Lock()
	tfrest.Logger.Println(req.Method, req.URL)
	if req.Body != nil {
		buf := new(bytes.Buffer)
		io.Copy(buf, req.Body)
		data := buf.Bytes()
		tfrest.Logger.Println("  ", string(data))
		req.Body = io.NopCloser(bytes.NewBuffer(data))
	}
	lock.Unlock()
	resp, err := rt.RoundTripper.RoundTrip(req)
	if err != nil {
		tfrest.Logger.Println(err.Error())
	}
	if resp != nil {
		if os.Getenv("DYNATRACE_HTTP_RESPONSE") == "true" {
			if resp.Body != nil {
				buf := new(bytes.Buffer)
				io.Copy(buf, resp.Body)
				data := buf.Bytes()
				resp.Body = io.NopCloser(bytes.NewBuffer(data))
				tfrest.Logger.Println(resp.Status, string(data))
			} else {
				tfrest.Logger.Println(resp.Status)
			}
		}
	}
	return resp, err
}

var R automation.Response

func (me *service) client() *automation.Client {
	if _, ok := http.DefaultClient.Transport.(*MyRoundTripper); !ok {
		if http.DefaultClient.Transport == nil {
			http.DefaultClient.Transport = &MyRoundTripper{http.DefaultTransport}
		} else {
			http.DefaultClient.Transport = &MyRoundTripper{http.DefaultClient.Transport}
		}
	}
	httpClient := auth.NewOAuthClient(context.TODO(), auth.OauthCredentials{
		ClientID:     me.credentials.Automation.ClientID,
		ClientSecret: me.credentials.Automation.ClientSecret,
		TokenURL:     me.credentials.Automation.TokenURL,
	})
	u, _ := url.Parse(me.credentials.Automation.EnvironmentURL)
	return automation.NewClient(rest.NewClient(u, httpClient))
}

func (me *service) Get(id string, v *workflows.Workflow) error {
	var err error
	var response automation.Response
	if response, err = me.client().Get(context.TODO(), apiClient.Workflows, id); err != nil {
		return err
	}
	if response.StatusCode != 200 {
		var e automationerr.ErrorEnvelope
		if e.Unmarshal(response.Data) {
			return e.Err.ToRESTError()
		}
		return tfrest.Error{Code: response.StatusCode, Message: string(response.Data)}
	}

	return json.Unmarshal(response.Data, &v)
}

func (me *service) SchemaID() string {
	return "automation:workflows"
}

type WorkflowStub struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func (me *service) List() (api.Stubs, error) {
	listResponse, err := me.client().List(context.TODO(), apiClient.Workflows)
	if err != nil {
		return nil, err
	}
	if apiErr, ok := listResponse.AsAPIError(); ok {
		return nil, tfrest.Error{Code: apiErr.StatusCode, Message: string(apiErr.Body)}
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

func (me *service) Create(v *workflows.Workflow) (stub *api.Stub, err error) {
	var data []byte
	if data, err = json.Marshal(v); err != nil {
		return nil, err
	}
	var response automation.Response
	if response, err = me.client().Create(context.TODO(), apiClient.Workflows, data); err != nil {
		return nil, err
	}
	if response.StatusCode == 201 {
		var workflowStub WorkflowStub
		if err := json.Unmarshal(response.Data, &workflowStub); err != nil {
			return nil, err
		}
		return &api.Stub{Name: v.Title, ID: workflowStub.ID}, nil
	}
	var e automationerr.ErrorEnvelope
	if e.Unmarshal(response.Data) {
		return nil, e.Err.ToRESTError()
	}
	return nil, tfrest.Error{Code: response.StatusCode, Message: string(response.Data)}
}

func (me *service) Update(id string, v *workflows.Workflow) (err error) {
	var data []byte
	if data, err = json.Marshal(v); err != nil {
		return err
	}
	var response automation.Response
	if response, err = me.client().Update(context.TODO(), apiClient.Workflows, id, data); err != nil {
		return err
	}
	if response.StatusCode == 200 {
		return nil
	}
	var e automationerr.ErrorEnvelope
	if e.Unmarshal(response.Data) {
		return e.Err.ToRESTError()
	}
	return tfrest.Error{Code: response.StatusCode, Message: string(response.Data)}
}

func (me *service) Delete(id string) error {
	response, err := me.client().Delete(context.TODO(), apiClient.Workflows, id)
	if response.StatusCode == 204 || response.StatusCode == 404 {
		return nil
	}
	if response.StatusCode != 0 {
		return tfrest.Error{Code: response.StatusCode, Message: string(response.Data)}
	}
	return err
}

func (me *service) New() *workflows.Workflow {
	return new(workflows.Workflow)
}

func (me *service) Name() string {
	return me.SchemaID()
}
