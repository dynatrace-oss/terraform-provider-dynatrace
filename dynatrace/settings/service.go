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

type Service[T Settings] interface {
	List(ctx context.Context) (api.Stubs, error)
	Get(ctx context.Context, id string, v T) error
	Create(ctx context.Context, v T) (*api.Stub, error)
	Update(ctx context.Context, id string, v T) error
	Delete(ctx context.Context, id string) error
	Validate(ctx context.Context, v T) error
	SchemaID() string
}
