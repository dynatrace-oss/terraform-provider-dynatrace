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

package generic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	generic "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/generic/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/shutdown"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/api/auth"
	crest "github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"net/http"
	"net/url"
)

func Service(credentials *settings.Credentials) settings.CRUDService[*generic.Settings] {
	return &service{credentials: credentials}
}

type service struct {
	credentials *settings.Credentials
}

var httpListener = &crest.HTTPListener{
	Callback: func(response crest.RequestResponse) {
		if response.Request != nil {
			if response.Request.URL != nil {
				if response.Request.Body != nil {
					body, _ := io.ReadAll(response.Request.Body)
					rest.Logger.Println(response.Request.Method, response.Request.URL.String()+"\n    "+string(body))
				} else {
					rest.Logger.Println(response.Request.Method, response.Request.URL)
				}
			}
		}
		if response.Response != nil {
			if response.Response.Body != nil {
				if os.Getenv("DYNATRACE_HTTP_RESPONSE") == "true" {
					body, _ := io.ReadAll(response.Response.Body)
					if body != nil {
						rest.Logger.Println(response.Response.StatusCode, string(body))
					} else {
						rest.Logger.Println(response.Response.StatusCode)
					}
				}
			}
		}
	},
}

type LoggingRoundTripper struct {
	RoundTripper http.RoundTripper
}

func (lrt *LoggingRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	rest.Logger.Println(r.Method, r.URL)
	if r.Body != nil {
		buf := new(bytes.Buffer)
		io.Copy(buf, r.Body)
		rest.Logger.Println("  ", buf.String())
		r.Body = io.NopCloser(bytes.NewBuffer(buf.Bytes()))
	}
	res, err := lrt.RoundTripper.RoundTrip(r)
	if err != nil {
		rest.Logger.Println("  error:", err.Error())
	}
	if res != nil && res.Body != nil {
		buf := new(bytes.Buffer)
		io.Copy(buf, res.Body)
		rest.Logger.Println("  =>", buf.String())
		res.Body = io.NopCloser(bytes.NewBuffer(buf.Bytes()))
	}
	return res, err
}

func (me *service) TokenClient() *crest.Client {
	var parsedURL *url.URL
	parsedURL, _ = url.Parse(me.credentials.URL)

	tokenClient := crest.NewClient(
		parsedURL,
		http.DefaultClient,
		crest.WithHTTPListener(httpListener),
	)

	tokenClient.SetHeader("User-Agent", "Dynatrace Terraform Provider")
	tokenClient.SetHeader("Authorization", "Api-Token "+me.credentials.Token)
	return tokenClient
}

func (me *service) Client(schemaIDs string) *settings20.Client {
	var parsedURL *url.URL
	parsedURL, _ = url.Parse(me.credentials.URL)

	tokenClient := me.TokenClient()

	if os.Getenv("DYNATRACE_DEBUG_GENERIC_SETTINGS") == "true" {
		http.DefaultClient.Transport = &LoggingRoundTripper{http.DefaultTransport}
	}

	oauthClient := crest.NewClient(
		parsedURL,
		auth.NewOAuthBasedClient(
			context.TODO(),
			clientcredentials.Config{
				ClientID:     me.credentials.Automation.ClientID,
				ClientSecret: me.credentials.Automation.ClientSecret,
				TokenURL:     me.credentials.Automation.TokenURL,
				AuthStyle:    oauth2.AuthStyleInParams}),
		crest.WithHTTPListener(httpListener),
	)

	oauthClient.SetHeader("User-Agent", "Dynatrace Terraform Provider")
	oauthClient.SetHeader("Authorization", "Api-Token "+me.credentials.Token)

	return settings20.NewClient(tokenClient, oauthClient, schemaIDs)
}

func (me *service) Get(ctx context.Context, id string, v *generic.Settings) error {
	var err error
	var response settings20.Response
	var settingsObject SettingsObject

	response, err = me.Client("").Get(context.TODO(), id)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		if err := rest.Envelope(response.Data, response.Request.URL, response.Request.Method); err != nil {
			return err
		}
		return fmt.Errorf("status code %d (expected: %d): %s", response.StatusCode, 200, string(response.Data))
	}
	if err := json.Unmarshal(response.Data, &settingsObject); err != nil {
		return err
	}

	v.Value = string(settingsObject.Value)
	v.Scope = settingsObject.Scope
	v.SchemaID = settingsObject.SchemaID

	return nil
}

type schemaStub struct {
	SchemaID string `json:"schemaId"`
}

type schemataResponse struct {
	Items []schemaStub `json:"items"`
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	tokenClient := me.TokenClient()
	response, err := tokenClient.GET(context.TODO(), "api/v2/settings/schemas", crest.RequestOptions{})
	if err != nil {
		return nil, err
	}
	var schemata schemataResponse
	if response.Body != nil {
		json.NewDecoder(response.Body).Decode(&schemata)
	}
	if len(schemata.Items) == 0 {
		return api.Stubs{}, nil
	}
	schemaIDs := []string{}
	for _, schemaStub := range schemata.Items {
		if strings.HasPrefix(schemaStub.SchemaID, "app:") {
			schemaIDs = append(schemaIDs, schemaStub.SchemaID)
		}
	}
	if len(schemaIDs) == 0 {
		return api.Stubs{}, nil
	}
	var stubs api.Stubs
	for _, schemaID := range schemaIDs {
		response, err := me.Client(schemaID).List(context.TODO())
		if response.StatusCode != 200 {
			if err := rest.Envelope(response.Data, response.Request.URL, response.Request.Method); err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("status code %d (expected: %d): %s", response.StatusCode, 200, string(response.Data))
		}
		if err != nil {
			return nil, err
		}
		for _, item := range response.Items {
			newItem := new(generic.Settings)
			newItem.Value = string(item.Value)
			newItem.Scope = item.Scope
			newItem.SchemaID = schemaID
			newItem.Scope = item.Scope
			stubs = append(stubs, &api.Stub{ID: item.ID, Name: item.ID})

		}
	}
	return stubs, nil
}

func (me *service) Validate(v *generic.Settings) error {
	return nil // Settings 2.0 doesn't offer validation
}

func (me *service) Create(ctx context.Context, v *generic.Settings) (*api.Stub, error) {
	return me.create(ctx, v)
}

type Matcher interface {
	Match(o any) bool
}

const errMsgOAuthRequired = "an OAuth Client is required for creating these settings. The configured credentials are currently based on API Tokens only. More information: https://registry.terraform.io/providers/dynatrace-oss/dynatrace/latest/docs/resources/generic_setting"

func (me *service) create(ctx context.Context, v *generic.Settings) (*api.Stub, error) {
	scope := "environment"
	if len(v.Scope) > 0 {
		scope = v.Scope
	}
	response, err := me.Client(v.SchemaID).Create(ctx, scope, []byte(v.Value))
	if response.StatusCode != 200 {
		if err := rest.Envelope(response.Data, response.Request.URL, response.Request.Method); err != nil {
			return nil, err
		}
		if response.StatusCode == 0 {
			return nil, errors.New(errMsgOAuthRequired)
		}
		return nil, fmt.Errorf("status code %d (expected: %d): %s", response.StatusCode, 200, string(response.Data))
	}
	if err != nil {
		return nil, err
	}

	stub := &api.Stub{ID: response.ID, Name: response.ID}
	return stub, nil
}

func (me *service) Update(ctx context.Context, id string, v *generic.Settings) error {
	response, err := me.Client("").Update(ctx, id, []byte(v.Value))
	if response.StatusCode != 200 {
		if err := rest.Envelope(response.Data, response.Request.URL, response.Request.Method); err != nil {
			return err
		}
		if response.StatusCode == 0 {
			return errors.New(errMsgOAuthRequired)
		}
		return fmt.Errorf("status code %d (expected: %d): %s", response.StatusCode, 200, string(response.Data))
	}

	return err
}

func (me *service) Delete(ctx context.Context, id string) error {
	return me.delete(ctx, id, 0)
}

func (me *service) delete(ctx context.Context, id string, numRetries int) error {
	response, err := me.Client("").Delete(ctx, id)
	if response.StatusCode != 204 {
		if err = rest.Envelope(response.Data, response.Request.URL, response.Request.Method); err != nil {
			return err
		}
		err = fmt.Errorf("status code %d (expected: %d): %s", response.StatusCode, 204, string(response.Data))
	}

	if err != nil && strings.Contains(err.Error(), "Deletion of value(s) is not allowed") {
		return nil
	}
	if err != nil && strings.Contains(err.Error(), "Internal Server Error occurred") {
		if numRetries == 10 {
			return err
		}
		time.Sleep(6 * time.Second)
		return me.delete(ctx, id, numRetries+1)
	}
	return err

}

func (me *service) SchemaID() string {
	return "generic"
}

type QueryParams struct {
	Schema string
	Scope  string
	Filter string
}

func (me *service) ListSpecific(query QueryParams) (api.Stubs, error) {
	client := me.TokenClient()

	stubs := api.Stubs{}
	nextPage := true
	var nextPageKey *string
	for nextPage {
		var sol settings20.SettingsObjectList
		var options crest.RequestOptions
		if nextPageKey == nil {
			options.QueryParams = url.Values{
				"fields":    []string{"objectId,scope,value,schemaId"},
				"pageSize":  []string{"100"},
				"schemaIds": []string{query.Schema},
				"scopes":    []string{query.Scope},
				"filter":    []string{query.Filter},
			}
		} else {
			options.QueryParams = url.Values{
				"nextPageKey": []string{*nextPageKey},
			}
		}

		response, err := client.GET(context.TODO(), "api/v2/settings/objects", options)
		if err != nil {
			return nil, err
		}

		if response.Body != nil {
			json.NewDecoder(response.Body).Decode(&sol)
		}
		if len(sol.Items) == 0 {
			return api.Stubs{}, nil
		}
		if shutdown.System.Stopped() {
			return stubs, nil
		}

		if len(sol.Items) > 0 {
			for _, item := range sol.Items {
				newItem := new(generic.Settings)
				newItem.Value = string(item.Value)
				newItem.Scope = item.Scope
				newItem.SchemaID = item.SchemaID
				stubs = append(stubs, &api.Stub{ID: item.ObjectID, Name: item.ObjectID, Value: newItem})
			}
		}
		nextPageKey = sol.NextPageKey
		nextPage = (nextPageKey != nil)
	}

	return stubs, nil
}
