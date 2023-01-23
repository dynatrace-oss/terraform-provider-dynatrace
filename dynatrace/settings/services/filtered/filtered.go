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

package filtered

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/cache"
)

type Filter[T settings.Settings] interface {
	Filter(v T) (bool, error)
	Suffix() string
}

type FilterService[T settings.Settings] struct {
	Service settings.CRUDService[T]
	Filter  Filter[T]
}

func (me *FilterService[T]) List() (settings.Stubs, error) {
	var err error
	var stubs settings.Stubs
	var filteredStubs settings.Stubs
	if stubs, err = me.Service.List(); err != nil {
		return nil, err
	}
	for _, stub := range stubs {
		if stub.Value == nil {
			return filteredStubs, nil
		}
		var allowed bool
		if allowed, err = me.Filter.Filter(stub.Value.(T)); err != nil {
			return nil, err
		} else if allowed {
			filteredStubs = append(filteredStubs, stub)
		}
	}

	return filteredStubs, nil
}

func (me *FilterService[T]) Validate(v T) error {
	if validator, ok := me.Service.(settings.Validator[T]); ok {
		return validator.Validate(v)
	}
	return nil
}

func (me *FilterService[T]) Get(id string, v T) error {
	return me.Service.Get(id, v)
}

func (me *FilterService[T]) Create(v T) (*settings.Stub, error) {
	return me.Service.Create(v)
}

func (me *FilterService[T]) Update(id string, v T) error {
	return me.Service.Update(id, v)
}

func (me *FilterService[T]) Delete(id string) error {
	return me.Service.Delete(id)
}

func (me *FilterService[T]) SchemaID() string {
	return me.Service.SchemaID() + ":" + me.Filter.Suffix()
}

func (me *FilterService[T]) NoCache() bool {
	return true
}

func Service[T settings.Settings](service settings.CRUDService[T], f Filter[T]) settings.CRUDService[T] {
	return &FilterService[T]{Service: cache.CRUD(service), Filter: f}
}
