//go:build unit

/*
 * @license
 * Copyright 2026 Dynatrace LLC
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package resources_test

import (
	"context"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type service struct {
	Stub *api.Stub
	Err  error
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return nil, nil
}

func (me *service) Get(ctx context.Context, id string, v *mockSchema) error {
	return me.Err
}

func (me *service) SchemaID() string {
	return "mock"
}

func (me *service) Create(ctx context.Context, v *mockSchema) (*api.Stub, error) {
	return me.Stub, me.Err
}

func (me *service) Update(ctx context.Context, id string, v *mockSchema) error {
	return me.Err
}

func (me *service) Delete(ctx context.Context, id string) error {
	return me.Err
}

type mockSchema struct {
	Name string
}

func (me *mockSchema) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

func (me *mockSchema) SchemaId() string {
	return "mock"
}
func (me *mockSchema) MarshalHCL(properties hcl.Properties) error {
	return properties.Encode("name", me.Name)
}

func (me *mockSchema) UnmarshalHCL(properties hcl.Decoder) error {
	return properties.Decode("name", &me.Name)
}

func MockService(stub *api.Stub, err error) func(credentials *rest.Credentials) settings.CRUDService[*mockSchema] {
	return func(credentials *rest.Credentials) settings.CRUDService[*mockSchema] {
		return &service{
			Err:  err,
			Stub: stub,
		}
	}
}
