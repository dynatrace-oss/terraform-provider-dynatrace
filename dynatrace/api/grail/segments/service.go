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

package segments

import (
	"context"
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	segmentsapi "github.com/dynatrace/dynatrace-configuration-as-code-core/api"
	segmentsclient "github.com/dynatrace/dynatrace-configuration-as-code-core/clients/segments"

	segments "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/grail/segments/settings"
)

func Service(credentials *rest.Credentials) settings.CRUDService[*segments.Segment] {
	return &service{credentials}
}

type service struct {
	credentials *rest.Credentials
}

func (me *service) client(ctx context.Context) (*segmentsclient.Client, error) {
	platformClient, err := rest.CreatePlatformClient(ctx, me.credentials.OAuth.EnvironmentURL, me.credentials)
	if err != nil {
		return nil, err
	}
	return segmentsclient.NewClient(platformClient), nil
}

func (me *service) Get(ctx context.Context, id string, v *segments.Segment) (err error) {
	client, err := me.client(ctx)
	if err != nil {
		return err
	}
	response, err := client.Get(ctx, id)
	if err != nil {
		if apiError, ok := err.(segmentsapi.APIError); ok {
			return rest.Error{Code: apiError.StatusCode, Message: apiError.Error()}
		}
		return err
	}

	return json.Unmarshal(response.Data, &v)
}

func (me *service) SchemaID() string {
	return "storage:filter-segments"
}

type SegmentStub struct {
	UID         string `json:"uid"`
	Name        string `json:"name"`
	IsReadyMade bool   `json:"isReadyMade"`
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	client, err := me.client(ctx)
	if err != nil {
		return nil, err
	}
	listResponse, err := client.List(ctx)
	if err != nil {
		if apiError, ok := err.(segmentsapi.APIError); ok {
			return nil, rest.Error{Code: apiError.StatusCode, Message: apiError.Error()}
		}
		return nil, err
	}
	var segments []SegmentStub
	if err := json.Unmarshal(listResponse.Data, &segments); err != nil {
		return nil, err
	}
	var stubs api.Stubs
	for _, segment := range segments {
		// Ready-made segments are unmodifiable, so we skip them in the listing
		if segment.IsReadyMade {
			continue
		}
		stubs = append(stubs, &api.Stub{ID: segment.UID, Name: segment.Name})
	}

	return stubs, nil
}

func (me *service) Validate(_ *segments.Segment) error {
	return nil // no endpoint for that
}

func (me *service) Create(ctx context.Context, v *segments.Segment) (stub *api.Stub, err error) {
	client, err := me.client(ctx)
	if err != nil {
		return nil, err
	}
	var data []byte
	if data, err = json.Marshal(v); err != nil {
		return nil, err
	}
	response, err := client.Create(ctx, data)
	if err != nil {
		if apiError, ok := err.(segmentsapi.APIError); ok {
			return nil, rest.Error{Code: apiError.StatusCode, Message: apiError.Error()}
		}
		return nil, err
	}

	var segmentStub SegmentStub
	if err := json.Unmarshal(response.Data, &segmentStub); err != nil {
		return nil, err
	}

	return &api.Stub{Name: v.Name, ID: segmentStub.UID}, nil
}

func (me *service) Update(ctx context.Context, id string, v *segments.Segment) (err error) {
	client, err := me.client(ctx)
	if err != nil {
		return err
	}
	var data []byte
	if data, err = json.Marshal(v); err != nil {
		return err
	}

	_, err = client.Update(ctx, id, data)
	if err != nil {
		if apiError, ok := err.(segmentsapi.APIError); ok {
			return rest.Error{Code: apiError.StatusCode, Message: apiError.Error()}
		}
		return err
	}

	return nil
}

func (me *service) Delete(ctx context.Context, id string) error {
	client, err := me.client(ctx)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, id)
	if err != nil {
		if apiError, ok := err.(segmentsapi.APIError); ok {
			return rest.Error{Code: apiError.StatusCode, Message: apiError.Error()}
		}
		return err
	}

	return nil
}

func (me *service) New() *segments.Segment {
	return new(segments.Segment)
}
