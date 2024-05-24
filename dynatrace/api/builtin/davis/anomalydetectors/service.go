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

package anomalydetectors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	anomalydetectors "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/davis/anomalydetectors/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/api/auth"
	crest "github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const SchemaVersion = "1.0.1"
const SchemaID = "builtin:davis.anomaly-detectors"

func Service(credentials *settings.Credentials) settings.CRUDService[*anomalydetectors.Settings] {
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
	if _, ok := http.DefaultClient.Transport.(*LoggingRoundTripper); !ok {
		if http.DefaultClient.Transport == nil {
			http.DefaultClient.Transport = &LoggingRoundTripper{http.DefaultTransport}
		} else {
			http.DefaultClient.Transport = &LoggingRoundTripper{http.DefaultClient.Transport}
		}
	}

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
	logging.File.Println(me.credentials.Automation.TokenURL)

	oauthClient.SetHeader("User-Agent", "Dynatrace Terraform Provider")
	oauthClient.SetHeader("Authorization", "Api-Token "+me.credentials.Token)

	return settings20.NewClient(tokenClient, oauthClient, schemaIDs)
}

func (me *service) Create(v *anomalydetectors.Settings) (*api.Stub, error) {
	scope := "environment"
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	response, err := me.Client(SchemaID).Create(context.TODO(), scope, data)
	if response.StatusCode != 200 {
		if err := rest.Envelope(response.Data, response.Request.URL, response.Request.Method); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("status code %d (expected: %d): %s", response.StatusCode, 200, string(response.Data))
	}
	if err != nil {
		return nil, err
	}

	stub := &api.Stub{ID: response.ID, Name: response.ID}
	return stub, nil
}

func (me *service) Update(id string, v *anomalydetectors.Settings) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	response, err := me.Client("").Update(context.TODO(), id, data)
	if response.StatusCode != 200 {
		if err := rest.Envelope(response.Data, response.Request.URL, response.Request.Method); err != nil {
			return err
		}
		return fmt.Errorf("status code %d (expected: %d): %s", response.StatusCode, 200, string(response.Data))
	}

	return err
}

func (me *service) Validate(v *anomalydetectors.Settings) error {
	return nil // Settings 2.0 doesn't offer validation
}

func (me *service) Delete(id string) error {
	response, err := me.Client("").Delete(context.TODO(), id)
	if response.StatusCode != 204 {
		if err = rest.Envelope(response.Data, response.Request.URL, response.Request.Method); err != nil {
			return err
		}
		err = fmt.Errorf("status code %d (expected: %d): %s", response.StatusCode, 204, string(response.Data))
	}

	if err != nil && strings.Contains(err.Error(), "Deletion of value(s) is not allowed") {
		return nil
	}

	return err
}

type SettingsObject struct {
	SchemaVersion string          `json:"schemaVersion"`
	SchemaID      string          `json:"schemaId"`
	Scope         string          `json:"scope"`
	Value         json.RawMessage `json:"value"`
}

func (me *service) Get(id string, v *anomalydetectors.Settings) error {
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
	if err := json.Unmarshal(settingsObject.Value, &v); err != nil {
		return err
	}

	return nil
}

func (me *service) List() (api.Stubs, error) {
	var stubs api.Stubs
	response, err := me.Client(SchemaID).List(context.TODO())
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
		stubs = append(stubs, &api.Stub{ID: item.ID, Name: item.ID})

	}
	return stubs, nil
}

func (me *service) Name() string {
	return me.SchemaID()
}

func (me *service) SchemaID() string {
	return SchemaID
}
