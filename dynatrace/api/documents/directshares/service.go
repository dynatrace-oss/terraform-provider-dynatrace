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

package directshares

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	directshares "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/documents/directshares/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/monaco/pkg/client/auth"
	directshare "github.com/dynatrace-oss/terraform-provider-dynatrace/monaco/pkg/client/document/directshares"
)

func Service(credentials *settings.Credentials) settings.CRUDService[*directshares.DirectShare] {
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
	rest.Logger.Println(req.Method, req.URL)
	if req.Body != nil {
		buf := new(bytes.Buffer)
		io.Copy(buf, req.Body)
		data := buf.Bytes()
		rest.Logger.Println("  ", string(data))
		req.Body = io.NopCloser(bytes.NewBuffer(data))
	}
	lock.Unlock()
	resp, err := rt.RoundTripper.RoundTrip(req)
	if err != nil {
		rest.Logger.Println(err.Error())
	}
	if resp != nil {
		if os.Getenv("DYNATRACE_HTTP_RESPONSE") == "true" {
			if resp.Body != nil {
				buf := new(bytes.Buffer)
				io.Copy(buf, resp.Body)
				data := buf.Bytes()
				resp.Body = io.NopCloser(bytes.NewBuffer(data))
				rest.Logger.Println(resp.Status, string(data))
			} else {
				rest.Logger.Println(resp.Status)
			}
		}
	}
	return resp, err
}

func (me *service) client() *directshare.Client {
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
	return directshare.NewClient(me.credentials.Automation.EnvironmentURL, httpClient)
}

func (me *service) Get(id string, v *directshares.DirectShare) (err error) {
	var result *directshare.Response
	if result, err = me.client().GET(directshare.DirectShares, id); err != nil {
		return err
	}

	v.DocumentId = result.DocumentId
	v.ID = result.ID
	v.Access = strings.Join(result.Access, "-")
	v.Recipients = result.Recipients

	return nil
}

func (me *service) SchemaID() string {
	return "document:direct-shares"
}

func (me *service) List() (api.Stubs, error) {
	var stubs api.Stubs
	return stubs, nil // not implemented

}

func (me *service) Validate(v *directshares.DirectShare) error {
	return nil // no endpoint for that
}

func (me *service) Create(v *directshares.DirectShare) (stub *api.Stub, err error) {
	var id string

	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	if id, err = me.client().INSERT(directshare.DirectShares, data); err != nil {
		return nil, err
	}
	return &api.Stub{ID: id}, nil
}

func (me *service) Update(id string, v *directshares.DirectShare) (err error) {

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return me.client().UPDATE(directshare.DirectShares, id, data)
}

func (me *service) Delete(id string) error {
	return me.client().DELETE(directshare.DirectShares, id)
}

func (me *service) New() *directshares.DirectShare {
	return new(directshares.DirectShare)
}
