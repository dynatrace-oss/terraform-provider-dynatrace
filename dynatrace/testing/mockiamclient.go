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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/api"
	rest2 "github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
)

type MockIAMClient struct {
	POSTFunc       func(ctx context.Context, url string, payload any, options rest2.RequestOptions) (api.Response, error)
	PUTFunc        func(ctx context.Context, url string, payload any, options rest2.RequestOptions) (api.Response, error)
	GETFunc        func(ctx context.Context, url string, options rest2.RequestOptions) (api.Response, error)
	DELETEFunc     func(ctx context.Context, url string, options rest2.RequestOptions) (api.Response, error)
	AccountIDValue string
}

func (me *MockIAMClient) AccountID() string {
	return me.AccountIDValue
}

func (me *MockIAMClient) POST(ctx context.Context, url string, payload any, options rest2.RequestOptions) (api.Response, error) {
	return me.POSTFunc(ctx, url, payload, options)
}

func (me *MockIAMClient) PUT(ctx context.Context, url string, payload any, options rest2.RequestOptions) (api.Response, error) {
	return me.PUTFunc(ctx, url, payload, options)
}

func (me *MockIAMClient) GET(ctx context.Context, url string, options rest2.RequestOptions) (api.Response, error) {
	return me.GETFunc(ctx, url, options)
}

func (me *MockIAMClient) DELETE(ctx context.Context, url string, options rest2.RequestOptions) (api.Response, error) {
	return me.DELETEFunc(ctx, url, options)
}

// MockClientSet is a rest.ClientSet whose IAM client and credentials are supplied directly, for
// injecting a MockIAMClient into IAM services under test.
type MockClientSet struct {
	IAMClientValue   rest.IAMClient
	CredentialsValue *rest.Credentials
}

func (me *MockClientSet) IAMClient() rest.IAMClient { return me.IAMClientValue }

func (me *MockClientSet) Credentials() *rest.Credentials { return me.CredentialsValue }
