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

package settings

import (
	"context"
	"reflect"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
)

type RService[T Settings] interface {
	List(ctx context.Context) (api.Stubs, error)
	Get(ctx context.Context, id string, v T) error
	SchemaID() string
}

func FindByName[T Settings](ctx context.Context, service RService[T], name string) (stub *api.Stub, err error) {
	var stubs api.Stubs
	if stubs, err = service.List(ctx); err != nil {
		return nil, err
	}
	for _, stub := range stubs.ToStubs() {
		if stub.Name == name {
			return stub, nil
		}
	}
	return nil, nil
}

type ContextKey string

const ContextKeyStateConfig = ContextKey("state-config")

type CRUDService[T Settings] interface {
	List(ctx context.Context) (api.Stubs, error)
	Get(ctx context.Context, id string, v T) error
	SchemaID() string
	Create(ctx context.Context, v T) (*api.Stub, error)
	Update(ctx context.Context, id string, v T) error
	Delete(ctx context.Context, id string) error
}

type ListIDCRUDService[T Settings] interface {
	CRUDService[T]
	ListIDs(ctx context.Context) (api.Stubs, error)
}

type Validator[T Settings] interface {
	Validate(v T) error
}

func NewSettings[T Settings](service RService[T]) T {
	var proto T
	return reflect.New(reflect.ValueOf(proto).Type().Elem()).Interface().(T)
}
