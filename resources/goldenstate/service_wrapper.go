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

package goldenstate

import (
	"context"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

type ServiceFunc func(*settings.Credentials) BasicService

type BasicService interface {
	List(ctx context.Context) (api.Stubs, error)
	Delete(ctx context.Context, id string) error
}

func Wrap[T settings.Settings](fn func(credentials *settings.Credentials) settings.CRUDService[T]) func(credentials *settings.Credentials) BasicService {
	return func(credentials *settings.Credentials) BasicService {
		return &GenericService[T]{Service: fn(credentials)}
	}
}

type GenericService[T settings.Settings] struct {
	Service settings.CRUDService[T]
}

func (me *GenericService[T]) List(ctx context.Context) (api.Stubs, error) {
	return me.Service.List(ctx)
}

func (me *GenericService[T]) Delete(ctx context.Context, id string) error {
	return me.Service.Delete(ctx, id)
}
