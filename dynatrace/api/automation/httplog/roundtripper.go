package httplog

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
)

var lock sync.Mutex

type RoundTripper struct {
	RoundTripper http.RoundTripper
}

func (rt *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	lock.Lock()
	rest.Logger.Println(req.Method, req.URL)
	if req.Body != nil {
		buf := new(bytes.Buffer)
		io.Copy(buf, req.Body)
		data := buf.Bytes()
		rest.Logger.Println("  ", string(data))
		req.Body = io.NopCloser(bytes.NewBuffer(data))
	}
	lock.Unlock()
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
			} else {
				rest.Logger.Println(resp.Status)
			}
		}
	}
	return resp, err
}
