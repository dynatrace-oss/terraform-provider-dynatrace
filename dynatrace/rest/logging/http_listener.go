package logging

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/envutil"
	"strings"
	"sync"

	crest "github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
	"github.com/google/uuid"
)

var DYNATRACE_HTTP_OAUTH = envutil.GetBoolEnv(envutil.EnvHTTPOAuth, false)

func logResponse(ctx context.Context, id string, response *http.Response) {
	if response == nil {
		return
	}

	if response.Body == nil {
		Logger.Printf(ctx, "[%s] [RESPONSE] %d", id, response.StatusCode)
		return
	}
	if !envutil.GetBoolEnv(envutil.EnvHTTPResponse, false) {
		Logger.Printf(ctx, "[%s] [RESPONSE] %d", id, response.StatusCode)
		return
	}
	body, _ := io.ReadAll(response.Body)
	if body != nil {
		Logger.Printf(ctx, "           [%s] [RESPONSE] %d %s", id, response.StatusCode, string(body))
		return
	}
	Logger.Printf(ctx, "           [%s] [RESPONSE] %d", id, response.StatusCode)
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
	logRequest := func(ctx context.Context, id string, request *http.Request) {
		if request == nil {
			return
		}
		if request.URL == nil {
			return
		}
		if request.Body == nil {
			Logger.Printf(ctx, "[%s] [%s] [REQUEST ] %s %s", prefix, id, request.Method, request.URL.String())
			return
		}
		// if len(request.Header) > 0 {
		// 	for headerName, headerValue := range request.Header {
		// 		if len(headerValue) > 0 {
		// 			Logger.Printf(ctx, "[%s] [%s] [HEADER  ] %s => %s", prefix, id, headerName, headerValue[0])
		// 		}
		// 	}
		// }

		body, _ := io.ReadAll(request.Body)
		Logger.Printf(ctx, "[%s] [%s] [REQUEST ] %s %s", prefix, id, request.Method, request.URL.String())
		Logger.Printf(ctx, "           [%s] [PAYLOAD ] %s", id, string(body))
	}

	return &crest.HTTPListener{
		Callback: func(response crest.RequestResponse) {
			ctx := requestContext(response)
			logRequest(ctx, response.ID, response.Request)
			logResponse(ctx, response.ID, response.Response)
		},
	}
}

var lock sync.Mutex

func InstallRoundTripper() {
	lock.Lock()
	defer lock.Unlock()
	if _, ok := http.DefaultClient.Transport.(*RoundTripper); !ok {
		if http.DefaultClient.Transport == nil {
			http.DefaultClient.Transport = &RoundTripper{RoundTripper: http.DefaultTransport}
		} else {
			http.DefaultClient.Transport = &RoundTripper{RoundTripper: http.DefaultClient.Transport}
		}
	}
}

type RoundTripper struct {
	RoundTripper http.RoundTripper
	lock         sync.Mutex
}

func (rt *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if !DYNATRACE_HTTP_OAUTH || !strings.Contains(req.URL.String(), "oauth2") {
		return rt.RoundTripper.RoundTrip(req)
	}
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
		Logger.Printf(ctx, "[%s]%s [PAYLOAD ] %s", id, category, string(data))
		req.Body = io.NopCloser(bytes.NewBuffer(data))
	}
	rt.lock.Unlock()
	resp, err := rt.RoundTripper.RoundTrip(req)
	if err != nil {
		Logger.Printf(ctx, "[%s]%s [ERROR   ] %s", id, category, err.Error())
	}
	if resp != nil {
		if os.Getenv("DYNATRACE_HTTP_RESPONSE") == "true" {
			if resp.Body != nil {
				buf := new(bytes.Buffer)
				io.Copy(buf, resp.Body)
				data := buf.Bytes()
				resp.Body = io.NopCloser(bytes.NewBuffer(data))
				if os.Getenv("DT_DEBUG_IAM_BEARER") == "true" || category != " [OAUTH]" {
					Logger.Printf(ctx, "[%s]%s [RESPONSE] %v %v", id, category, resp.Status, string(data))
				}
			} else {
				Logger.Printf(ctx, "[%s]%s [RESPONSE] %s", id, category, resp.Status)
			}
		}
	}
	return resp, err
}
