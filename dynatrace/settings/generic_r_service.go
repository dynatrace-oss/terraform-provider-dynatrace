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

type GenericRService[T Settings] struct {
	Service RService[T]
}

func (me *GenericRService[T]) NoCache() bool {
	if ncs, ok := me.Service.(NoCacheService); ok {
		return ncs.NoCache()
	}
	return false
}

func (me *GenericRService[T]) List() (Stubs, error) {
	return me.Service.List()
}

func (me *GenericRService[T]) Get(id string, v Settings) error {
	return me.Service.Get(id, v.(T))
}

func (me *GenericRService[T]) SchemaID() string {
	return me.Service.SchemaID()
}
