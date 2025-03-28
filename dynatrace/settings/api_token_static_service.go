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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
)

type staticService[T Settings] struct {
	schemaID string
	client   rest.Client
	url      string
	stub     api.Stub
}

func APITokenStaticService[T Settings](credentials *rest.Credentials, schemaID string, url string, stub api.Stub) Service[T] {
	return &staticService[T]{
		schemaID: schemaID,
		client:   rest.APITokenClient(credentials),
		url:      url,
		stub:     stub,
	}
}

func (me *staticService[T]) Get(ctx context.Context, id string, v T) error {
	return me.client.Get(ctx, me.url, 200).Finish(v)
}

func (me *staticService[T]) List(ctx context.Context) (api.Stubs, error) {
	return api.Stubs{&me.stub}, nil
}

func (me *staticService[T]) Create(ctx context.Context, v T) (*api.Stub, error) {
	return &me.stub, me.Update(ctx, me.stub.ID, v)
}

func (me *staticService[T]) Validate(ctx context.Context, v T) error {
	return me.client.Post(ctx, me.url+"/validator", v, 204).Finish()
}

func (me *staticService[T]) Delete(ctx context.Context, id string) error {
	return nil
}

func (me *staticService[T]) Update(ctx context.Context, id string, v T) error {
	return me.client.Put(ctx, me.url, v, 204).Finish()
}

func (me *staticService[T]) SchemaID() string {
	return me.schemaID
}

func (me *staticService[T]) Name() string {
	return me.schemaID
}
