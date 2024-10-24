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

package httplog

import (
	"context"
	"io"
	"os"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	crest "github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
	"github.com/google/uuid"
)

var HTTPListener = &crest.HTTPListener{
	Callback: func(response crest.RequestResponse) {
		id := uuid.NewString()
		ctx := context.TODO()
		if response.Request != nil && response.Request.Context() != nil {
			ctx = response.Request.Context()
		}
		if response.Request != nil {
			if response.Request.URL != nil {
				if response.Request.Body != nil {
					body, _ := io.ReadAll(response.Request.Body)
					rest.Logger.Printf(ctx, "[%s] %s %s", id, response.Request.Method, response.Request.URL.String())
					rest.Logger.Printf(ctx, "[%s] [PAYLOAD] %s", id, string(body))
				} else {
					rest.Logger.Printf(ctx, "[%s] %s %s", id, response.Request.Method, response.Request.URL.String())
				}
			}
		}
		if response.Response != nil {
			if response.Response.Body != nil {
				if os.Getenv("DYNATRACE_HTTP_RESPONSE") == "true" {
					body, _ := io.ReadAll(response.Response.Body)
					if body != nil {
						rest.Logger.Printf(ctx, "[%s] [RESPONSE] %d %s", id, response.Response.StatusCode, string(body))
					} else {
						rest.Logger.Printf(ctx, "[%s] [RESPONSE] %d", id, response.Response.StatusCode)
					}
				}
			}
		}
	},
}
