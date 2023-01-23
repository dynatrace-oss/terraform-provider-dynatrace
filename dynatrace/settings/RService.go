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
	"reflect"
)

type RService[T Settings] interface {
	List() (Stubs, error)
	Get(id string, v T) error
	SchemaID() string
}

func FindByName[T Settings](service RService[T], name string) (stub *Stub, err error) {
	var stubs Stubs
	if stubs, err = service.List(); err != nil {
		return nil, err
	}
	for _, stub := range stubs.ToStubs() {
		if stub.Name == name {
			return stub, nil
		}
	}
	return nil, nil
}

type CRUDService[T Settings] interface {
	List() (Stubs, error)
	Get(id string, v T) error
	SchemaID() string
	Create(v T) (*Stub, error)
	Update(id string, v T) error
	Delete(id string) error
}

type Validator[T Settings] interface {
	Validate(v T) error
}

func NewSettings[T Settings](service RService[T]) T {
	var proto T
	return reflect.New(reflect.ValueOf(proto).Type().Elem()).Interface().(T)
}
