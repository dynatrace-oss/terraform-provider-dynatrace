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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
)

type ServiceOptions[T Settings] struct {
	Get            func(args ...string) string
	List           func(args ...string) string
	CreateURL      func(v T) string
	ValidateURL    func(v T) string
	UpdateURL      func(id string, v T) string
	DeleteURL      func(id string) string
	Stubs          RecordStubs
	CompleteGet    func(client rest.Client, id string, v T) error
	CreateRetry    func(v T, err error) T
	DeleteRetry    func(id string, err error) (bool, error)
	CreateConfirm  int
	OnChanged      func(rest.Client, string, T) error
	OnBeforeUpdate func(id string, v T) error
	HasNoValidator bool
	Name           func(id string, v T) (string, error)
	HijackOnCreate func(err error, service RService[T], v T) (*Stub, error)
}

func (me *ServiceOptions[T]) Hijack(fn func(err error, service RService[T], v T) (*Stub, error)) *ServiceOptions[T] {
	me.HijackOnCreate = fn
	return me
}

func (me *ServiceOptions[T]) WithCreateRetry(fn func(v T, err error) T) *ServiceOptions[T] {
	me.CreateRetry = fn
	return me
}

func (me *ServiceOptions[T]) WithStubs(stubs RecordStubs) *ServiceOptions[T] {
	me.Stubs = stubs
	return me
}

func (me *ServiceOptions[T]) NoValidator() *ServiceOptions[T] {
	me.HasNoValidator = true
	return me
}

func (me *ServiceOptions[T]) WithOnChanged(onChanged func(rest.Client, string, T) error) *ServiceOptions[T] {
	me.OnChanged = onChanged
	return me
}
func (me *ServiceOptions[T]) WithDeleteRetry(deleteRetry func(id string, err error) (bool, error)) *ServiceOptions[T] {
	me.DeleteRetry = deleteRetry
	return me
}

func DefaultServiceOptions[T Settings](basePath string) *ServiceOptions[T] {
	return &ServiceOptions[T]{
		Get:  Path(basePath + "/%s"),
		List: Path(basePath),
	}
}
