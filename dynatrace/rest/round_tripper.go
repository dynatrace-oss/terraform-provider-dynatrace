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

package rest

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/envutils"
	"github.com/google/uuid"
)

var lock sync.Mutex

func InstallRoundTripper() {
	lock.Lock()
	defer lock.Unlock()

	// Avoid installing this round tripper multiple times
	if _, ok := http.DefaultClient.Transport.(*RoundTripper); !ok {
		return
	}

	if http.DefaultClient.Transport == nil {
		http.DefaultClient.Transport = &RoundTripper{RoundTripper: http.DefaultTransport}
	} else {
		http.DefaultClient.Transport = &RoundTripper{RoundTripper: http.DefaultClient.Transport}
	}
}

type RoundTripper struct {
	RoundTripper http.RoundTripper
	lock         sync.Mutex
}

func (rt *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	id := uuid.NewString()
	if v := req.Context().Value("request.id"); v != nil {
		if sv, ok := v.(string); ok {
			id = sv
		}
	}
	rt.lock.Lock()
	ctx := req.Context()
	category := ""
	if strings.Contains(req.URL.String(), "oauth2") {
		category = " [OAUTH]"
	}

	Logger.Println(ctx, fmt.Sprintf("[%s]%s %s %s", id, category, req.Method, req.URL.String()))
	if req.Body != nil {
		buf := new(bytes.Buffer)
		io.Copy(buf, req.Body)
		data := buf.Bytes()
		Logger.Printf(ctx, "[%s]%s [PAYLOAD] %s", id, category, string(data))
		req.Body = io.NopCloser(bytes.NewBuffer(data))
	}
	rt.lock.Unlock()
	resp, err := rt.RoundTripper.RoundTrip(req)
	if err != nil {
		Logger.Printf(ctx, "[%s]%s [ERROR] %s", id, category, err.Error())
	}
	if resp != nil {
		if envutils.DynatraceHTTPResponse.Get() {
			if resp.Body != nil {
				buf := new(bytes.Buffer)
				io.Copy(buf, resp.Body)
				data := buf.Bytes()
				resp.Body = io.NopCloser(bytes.NewBuffer(data))
				if envutils.DTDebugIAMBearer.Get() || category != " [OAUTH]" {
					Logger.Printf(ctx, "[%s]%s [RESPONSE] %v %v", id, category, resp.Status, string(data))
				}
			} else {
				Logger.Printf(ctx, "[%s]%s [RESPONSE] %s", id, category, resp.Status)
			}
		}
	}
	return resp, err
}
