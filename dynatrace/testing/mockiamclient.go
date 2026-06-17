/*
 * @license
 * Copyright 2026 Dynatrace LLC
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package testing

import (
	"context"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam"
	rest2 "github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
)

type MockIAMClient struct {
	POSTFunc   func(ctx context.Context, url string, payload any, options rest2.RequestOptions, expectedResponseCode int) ([]byte, error)
	PUTFunc    func(ctx context.Context, url string, payload any, options rest2.RequestOptions, expectedResponseCode int) ([]byte, error)
	GETFunc    func(ctx context.Context, url string, options rest2.RequestOptions, expectedResponseCode int) ([]byte, error)
	DELETEFunc func(ctx context.Context, url string, options rest2.RequestOptions, expectedResponseCode int) ([]byte, error)
}

func (me *MockIAMClient) POST(ctx context.Context, url string, payload any, options rest2.RequestOptions, expectedResponseCode int) ([]byte, error) {
	return me.POSTFunc(ctx, url, payload, options, expectedResponseCode)
}

func (me *MockIAMClient) PUT(ctx context.Context, url string, payload any, options rest2.RequestOptions, expectedResponseCode int) ([]byte, error) {
	return me.PUTFunc(ctx, url, payload, options, expectedResponseCode)
}

func (me *MockIAMClient) PUT_MULTI_RESPONSE(ctx context.Context, url string, payload any, options rest2.RequestOptions, expectedResponseCodes []int) ([]byte, error) {
	panic("mock doesnt support PUT_MULTI_RESPONSE")
}

func (me *MockIAMClient) GET(ctx context.Context, url string, options rest2.RequestOptions, expectedResponseCode int) ([]byte, error) {
	return me.GETFunc(ctx, url, options, expectedResponseCode)
}

func (me *MockIAMClient) DELETE(ctx context.Context, url string, options rest2.RequestOptions, expectedResponseCode int) ([]byte, error) {
	return me.DELETEFunc(ctx, url, options, expectedResponseCode)
}

func (me *MockIAMClient) DELETE_MULTI_RESPONSE(ctx context.Context, url string, options rest2.RequestOptions, expectedResponseCodes []int) ([]byte, error) {
	panic("mock doesnt support DELETE_MULTI_RESPONSE")
}

type MockIAMClientGetter struct {
	Client iam.IAMClient
}

func (me *MockIAMClientGetter) New(_ context.Context) iam.IAMClient {
	return me.Client
}
