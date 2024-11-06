package httplog

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/google/uuid"
)

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
	id := uuid.NewString()
	rt.lock.Lock()
	ctx := req.Context()
	category := ""
	if strings.Contains(req.URL.String(), "oauth2") {
		category = " [OAUTH]"
	}

	rest.Logger.Println(ctx, fmt.Sprintf("[%s]%s %s %s", id, category, req.Method, req.URL.String()))
	if req.Body != nil {
		buf := new(bytes.Buffer)
		io.Copy(buf, req.Body)
		data := buf.Bytes()
		rest.Logger.Printf(ctx, "[%s]%s [PAYLOAD] %s", id, category, string(data))
		req.Body = io.NopCloser(bytes.NewBuffer(data))
	}
	rt.lock.Unlock()
	resp, err := rt.RoundTripper.RoundTrip(req)
	if err != nil {
		rest.Logger.Printf(ctx, "[%s]%s [ERROR] %s", id, category, err.Error())
	}
	if resp != nil {
		if os.Getenv("DYNATRACE_HTTP_RESPONSE") == "true" {
			if resp.Body != nil {
				buf := new(bytes.Buffer)
				io.Copy(buf, resp.Body)
				data := buf.Bytes()
				resp.Body = io.NopCloser(bytes.NewBuffer(data))
				if os.Getenv("DT_DEBUG_IAM_BEARER") == "true" || category != " [OAUTH]" {
					rest.Logger.Printf(ctx, "[%s]%s [RESPONSE] %v %v", id, category, resp.Status, string(data))
				}
			} else {
				rest.Logger.Printf(ctx, "[%s]%s [RESPONSE] %s", id, category, resp.Status)
			}
		}
	}
	return resp, err
}
