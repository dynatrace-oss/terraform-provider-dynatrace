package httplog

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	rt.lock.Lock()
	rest.Logger.Println(req.Method, req.URL)
	if rest.StdoutLog {
		tflog.Debug(req.Context(), fmt.Sprintf("%s %s", req.Method, req.URL.String()))
	}
	if req.Body != nil {
		buf := new(bytes.Buffer)
		io.Copy(buf, req.Body)
		data := buf.Bytes()
		rest.Logger.Println("  ", string(data))
		if rest.StdoutLog {
			tflog.Debug(req.Context(), fmt.Sprintf("  %s", string(data)))
		}
		req.Body = io.NopCloser(bytes.NewBuffer(data))
	}
	rt.lock.Unlock()
	resp, err := rt.RoundTripper.RoundTrip(req)
	if err != nil {
		rest.Logger.Println(err.Error())
	}
	if resp != nil {
		if os.Getenv("DYNATRACE_HTTP_RESPONSE") == "true" {
			if resp.Body != nil {
				buf := new(bytes.Buffer)
				io.Copy(buf, resp.Body)
				data := buf.Bytes()
				resp.Body = io.NopCloser(bytes.NewBuffer(data))
				rest.Logger.Println(resp.Status, string(data))
				if rest.StdoutLog {
					tflog.Debug(req.Context(), fmt.Sprintf("%v %v", resp.Status, string(data)))
				}
			} else {
				rest.Logger.Println(resp.Status)
				if rest.StdoutLog {
					tflog.Debug(req.Context(), resp.Status)
				}
			}
		}
	}
	return resp, err
}
