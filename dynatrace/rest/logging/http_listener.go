/**
* @license
* Copyright 2026 Dynatrace LLC
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

package logging

import (
	"context"
	"io"
	"net/http"
	"os"

	crest "github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
)

func logRequest(ctx context.Context, id string, request *http.Request, prefix string) {
	if request == nil || request.URL == nil {
		return
	}

	Logger.Printf(ctx, "[%s] [%s] [REQUEST] %s %s", prefix, id, request.Method, request.URL.String())
	if request.Body == nil {
		return
	}

	body, _ := io.ReadAll(request.Body)
	Logger.Printf(ctx, "[%s] [%s] [PAYLOAD] %s", prefix, id, string(body))
}

func logResponse(ctx context.Context, id string, response *http.Response, prefix string) {
	if response == nil {
		return
	}

	var bodyStr string
	if response.Body != nil && os.Getenv("DYNATRACE_HTTP_RESPONSE") == "true" {
		body, _ := io.ReadAll(response.Body)
		if body != nil {
			bodyStr = " " + string(body)
		}
	}
	Logger.Printf(ctx, "[%s] [%s] [RESPONSE] %d%s", prefix, id, response.StatusCode, bodyStr)
}

func requestContext(response crest.RequestResponse) context.Context {
	if response.Request != nil {
		return response.Request.Context()
	}
	if response.Response != nil && response.Response.Request != nil {
		return response.Response.Request.Context()
	}
	return context.Background()
}

func HTTPListener(prefix string) *crest.HTTPListener {
	return &crest.HTTPListener{
		Callback: func(response crest.RequestResponse) {
			ctx := requestContext(response)
			logRequest(ctx, response.ID, response.Request, prefix)
			logResponse(ctx, response.ID, response.Response, prefix)
		},
	}
}
