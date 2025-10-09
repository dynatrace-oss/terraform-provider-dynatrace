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
	"errors"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	openpipeline "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/openpipeline/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	cacapi "github.com/dynatrace/dynatrace-configuration-as-code-core/api"
	caclib "github.com/dynatrace/dynatrace-configuration-as-code-core/clients/openpipeline"
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

func MetricsService(credentials *rest.Credentials) settings.CRUDService[*openpipeline.Configuration] {
	return &service{credentials: credentials, kind: "metrics", schemaSuffix: "metrics"}
}

func UserSessionsService(credentials *rest.Credentials) settings.CRUDService[*openpipeline.Configuration] {
	return &service{credentials: credentials, kind: "usersessions", schemaSuffix: "user.sessions"}
}

func DavisProblemsService(credentials *rest.Credentials) settings.CRUDService[*openpipeline.Configuration] {
	return &service{credentials: credentials, kind: "davis.problems", schemaSuffix: "davis.problems"}
}

func DavisEventsService(credentials *rest.Credentials) settings.CRUDService[*openpipeline.Configuration] {
	return &service{credentials: credentials, kind: "davis.events", schemaSuffix: "davis.events"}
}

func SystemEventsService(credentials *rest.Credentials) settings.CRUDService[*openpipeline.Configuration] {
	return &service{credentials: credentials, kind: "system.events", schemaSuffix: "system.events"}
}

func UserEventsService(credentials *rest.Credentials) settings.CRUDService[*openpipeline.Configuration] {
	return &service{credentials: credentials, kind: "user.events", schemaSuffix: "user.events"}
}

func SpansService(credentials *rest.Credentials) settings.CRUDService[*openpipeline.Configuration] {
	return &service{credentials: credentials, kind: "spans", schemaSuffix: "spans"}
}

type service struct {
	kind         string
	schemaSuffix string
	credentials  *rest.Credentials
}

func (s *service) createClient(ctx context.Context) (*caclib.Client, error) {
	platformClient, err := rest.CreatePlatformClient(ctx, s.credentials.OAuth.EnvironmentURL, s.credentials)
	if err != nil {
		return nil, err
	}
	return caclib.NewClient(platformClient), nil
}

func (s *service) List(ctx context.Context) (api.Stubs, error) {
	client, err := s.createClient(ctx)
	if err != nil {
		return nil, err
	}
	// after migration the resource isn't available anymore (404)
	if _, err = client.Get(ctx, s.kind); err != nil {
		return nil, err
	}

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
		apiErr := cacapi.APIError{}
		if errors.As(err, &apiErr) {
			return rest.Envelope(apiErr.Body, s.credentials.OAuth.EnvironmentURL, "GET")
		}
		return err
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
