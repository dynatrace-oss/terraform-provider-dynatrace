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

package rest

import (
	"context"
	"io"
)

type Client interface {
	Get(ctx context.Context, url string, expectedStatusCodes ...int) Request
	Post(ctx context.Context, url string, payload any, expectedStatusCodes ...int) Request
	Put(ctx context.Context, url string, payload any, expectedStatusCodes ...int) Request
	Delete(ctx context.Context, url string, expectedStatusCodes ...int) Request
	Upload(ctx context.Context, url string, reader io.ReadCloser, fileName string, expectedStatusCodes ...int) Request
}
