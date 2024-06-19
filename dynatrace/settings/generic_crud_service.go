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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
)

type GenericCRUDService[T Settings] struct {
	Service CRUDService[T]
}

func (me *GenericCRUDService[T]) NoCache() bool {
	if ncs, ok := me.Service.(NoCacheService); ok {
		return ncs.NoCache()
	}
	return false
}

func (me *GenericCRUDService[T]) List() (api.Stubs, error) {
	return me.Service.List()
}

func (me *GenericCRUDService[T]) GetWithContext(ctx context.Context, id string, v Settings) error {
	if contextGetter, ok := me.Service.(ContextGetter[T]); ok {
		return contextGetter.GetWithContext(ctx, id, v.(T))
	}
	return me.Service.Get(id, v.(T))
}

func (me *GenericCRUDService[T]) Get(id string, v Settings) error {
	return me.Service.Get(id, v.(T))
}

func (me *GenericCRUDService[T]) Create(v Settings) (*api.Stub, error) {
	return me.Service.Create(v.(T))
}

func (me *GenericCRUDService[T]) CreateWithContext(ctx context.Context, v Settings) (*api.Stub, error) {
	if contextCreator, ok := me.Service.(ContextCreator[T]); ok {
		return contextCreator.CreateWithContext(ctx, v.(T))
	}
	return me.Service.Create(v.(T))
}

func (me *GenericCRUDService[T]) Validate(v Settings) error {
	if validator, ok := me.Service.(Validator[T]); ok {
		return validator.Validate(v.(T))
	}
	return nil
}

func (me *GenericCRUDService[T]) Update(id string, v Settings) error {
	return me.Service.Update(id, v.(T))
}

func (me *GenericCRUDService[T]) UpdateWithContext(ctx context.Context, id string, v Settings) error {
	if contextUpdater, ok := me.Service.(ContextUpdater[T]); ok {
		return contextUpdater.UpdateWithContext(ctx, id, v.(T))
	}
	return me.Service.Update(id, v.(T))
}

func (me *GenericCRUDService[T]) Delete(id string) error {
	return me.Service.Delete(id)
}

func (me *GenericCRUDService[T]) DeleteWithContext(ctx context.Context, id string) error {
	if contextDeleter, ok := me.Service.(ContextDeleter[T]); ok {
		return contextDeleter.DeleteWithContext(ctx, id)
	}
	return me.Service.Delete(id)
}

func (me *GenericCRUDService[T]) SchemaID() string {
	return me.Service.SchemaID()
}
