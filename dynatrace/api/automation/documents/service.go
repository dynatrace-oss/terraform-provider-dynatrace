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

package documents

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	documents "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/automation/documents/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/monaco/pkg/client/auth"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/monaco/pkg/client/document"
)

func Service(credentials *settings.Credentials) settings.CRUDService[*documents.Document] {
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

func (me *service) client() *document.Client {
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
	return document.NewClient(me.credentials.Automation.EnvironmentURL, httpClient)
}

func (me *service) Get(id string, v *documents.Document) (err error) {
	var result *document.Response
	if result, err = me.client().GET(document.Documents, id); err != nil {
		return err
	}

	v.Actor = result.Actor
	v.Content = result.Content
	v.Name = result.Name
	v.Owner = result.Owner
	v.Type = result.Type
	v.Version = result.Version

	return nil
}

func (me *service) SchemaID() string {
	return "document:documents"
}

func (me *service) List() (api.Stubs, error) {
	result, err := me.client().LIST(document.Documents)
	if err != nil {
		return nil, err
	}
	var stubs api.Stubs
	for _, r := range result {
		var document documents.Document
		document.Actor = r.Actor
		document.Content = r.Content
		document.Name = r.Name
		document.Owner = r.Owner
		document.Type = r.Type
		document.Version = r.Version
		stubs = append(stubs, &api.Stub{ID: r.ID, Name: document.Name})
	}
	return stubs, nil
}

func (me *service) Validate(v *documents.Document) error {
	return nil // no endpoint for that
}

func (me *service) Create(v *documents.Document) (stub *api.Stub, err error) {
	var id string

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("type", v.Type)
	writer.WriteField("name", v.Name)

	part, err := writer.CreatePart(map[string][]string{
		"Content-Type":        {"application/json"},
		"Content-Disposition": {fmt.Sprintf(`form-data; name="content"; filename="%s"`, v.Name)},
	})
	if err != nil {
		panic(err)
	}

	part.Write([]byte(v.Content))
	writer.Close()

	if id, err = me.client().INSERT(document.Documents, body, writer.FormDataContentType()); err != nil {
		return nil, err
	}
	return &api.Stub{ID: id, Name: v.Name}, nil
}

func (me *service) Update(id string, v *documents.Document) (err error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("type", v.Type)
	writer.WriteField("name", v.Name)

	part, err := writer.CreatePart(map[string][]string{
		"Content-Type":        {"application/json"},
		"Content-Disposition": {fmt.Sprintf(`form-data; name="content"; filename="%s"`, v.Name)},
	})
	if err != nil {
		panic(err)
	}

	part.Write([]byte(v.Content))
	writer.Close()

	return me.client().UPDATE(document.Documents, id, body, writer.FormDataContentType())
}

func (me *service) Delete(id string) error {
	return me.client().DELETE(document.Documents, id)
}

func (me *service) New() *documents.Document {
	return new(documents.Document)
}

func (me *service) Name() string {
	return me.SchemaID()
}
