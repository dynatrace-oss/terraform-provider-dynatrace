/**
* @license
* Copyright 2024 Dynatrace LLC
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

package openpipeline

import (
	"context"
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/automation/httplog"
	openpipeline "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/openpipeline/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/clients"
	caclib "github.com/dynatrace/dynatrace-configuration-as-code-core/clients/openpipeline"
	"golang.org/x/oauth2/clientcredentials"
)

func EventsService(credentials *rest.Credentials) settings.CRUDService[*openpipeline.Configuration] {
	return &service{credentials: credentials, kind: "events", schemaSuffix: "events"}
}

func LogsService(credentials *rest.Credentials) settings.CRUDService[*openpipeline.Configuration] {
	return &service{credentials: credentials, kind: "logs", schemaSuffix: "logs"}
}

func BusinessEventsService(credentials *rest.Credentials) settings.CRUDService[*openpipeline.Configuration] {
	return &service{credentials: credentials, kind: "bizevents", schemaSuffix: "events.business"}
}

func SecurityEventsService(credentials *rest.Credentials) settings.CRUDService[*openpipeline.Configuration] {
	return &service{credentials: credentials, kind: "events.security", schemaSuffix: "events.security"}
}

func SDLCEventsService(credentials *rest.Credentials) settings.CRUDService[*openpipeline.Configuration] {
	return &service{credentials: credentials, kind: "events.sdlc", schemaSuffix: "events.sdlc"}
}

type service struct {
	kind         string
	schemaSuffix string
	credentials  *rest.Credentials
}

func (s *service) createClient(ctx context.Context) (*caclib.Client, error) {
	factory := clients.Factory().
		WithUserAgent("Dynatrace Terraform Provider").
		WithPlatformURL(s.credentials.OAuth.EnvironmentURL).
		WithOAuthCredentials(clientcredentials.Config{
			ClientID:     s.credentials.OAuth.ClientID,
			ClientSecret: s.credentials.OAuth.ClientSecret,
			TokenURL:     s.credentials.OAuth.TokenURL,
		}).
		WithHTTPListener(httplog.HTTPListener)

	return factory.OpenPipelineClient(ctx)
}

func (s *service) List(ctx context.Context) (api.Stubs, error) {
	//create exactly one stub for this ID

	stub := api.Stub{ID: s.kind, Name: s.kind}
	return api.Stubs{&stub}, nil
}

func (s *service) Get(ctx context.Context, id string, v *openpipeline.Configuration) error {
	client, err := s.createClient(ctx)
	if err != nil {
		return err
	}

	response, err := client.Get(ctx, s.kind)
	if err != nil {
		return err
	}

	if !response.IsSuccess() {
		return rest.Envelope(response.Data, s.credentials.OAuth.EnvironmentURL, "GET")
	}

	if err := json.Unmarshal(response.Data, &v); err != nil {
		return err
	}

	v.RemoveFixed()

	return nil
}

func (s *service) SchemaID() string {
	return "platform:openpipeline." + s.schemaSuffix
}

func (s *service) Create(ctx context.Context, v *openpipeline.Configuration) (*api.Stub, error) {
	if err := s.Update(ctx, s.kind, v); err != nil {
		return nil, err
	}

	return &api.Stub{ID: s.kind, Name: s.kind}, nil
}

func (s *service) Update(ctx context.Context, id string, v *openpipeline.Configuration) error {
	client, err := s.createClient(ctx)
	if err != nil {
		return err
	}

	v.Kind = id

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = client.Update(ctx, s.kind, b)
	return err
}

func (s *service) Delete(ctx context.Context, id string) error {
	emptyConfig := openpipeline.Configuration{Routing: &openpipeline.RoutingTable{}}
	return s.Update(ctx, s.kind, &emptyConfig)
}
