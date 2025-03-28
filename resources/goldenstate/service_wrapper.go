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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
)

type ServiceFunc func(*rest.Credentials) BasicService

type BasicService interface {
	List(ctx context.Context) (api.Stubs, error)
	Delete(ctx context.Context, id string) error
}

type ListIDsService interface {
	ListIDs(ctx context.Context) (api.Stubs, error)
}

func Wrap[T settings.Settings](fn func(credentials *rest.Credentials) settings.CRUDService[T]) func(credentials *rest.Credentials) BasicService {
	return func(credentials *rest.Credentials) BasicService {
		return &GenericService[T]{Service: fn(credentials)}
	}
}

type GenericService[T settings.Settings] struct {
	Service settings.CRUDService[T]
}

func (me *GenericService[T]) List(ctx context.Context) (api.Stubs, error) {
	if lister, ok := me.Service.(ListIDsService); ok {
		return lister.ListIDs(context.WithValue(ctx, settings20.ContextKeyDeleteableOnly, true))
	}
	return me.Service.List(ctx)
}

func (me *GenericService[T]) Delete(ctx context.Context, id string) error {
	return me.Service.Delete(ctx, id)
}
