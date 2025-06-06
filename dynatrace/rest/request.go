package rest

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest/logging"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
	"github.com/google/uuid"
)

var DYNATRACE_HTTP_LEGACY = (os.Getenv("DYNATRACE_HTTP_LEGACY") == "true")
var DYNATRACE_HTTP_OAUTH_PREFERENCE = (os.Getenv("DYNATRACE_HTTP_OAUTH_PREFERENCE") == "true")

type Request interface {
	Finish(v ...any) error
	Expect(codes ...int) Request
	OnResponse(func(resp *http.Response)) Request
}

type request struct {
	id         string
	ctx        context.Context
	client     Client
	url        string
	expect     statuscodes
	method     string
	payload    any
	fileName   string
	headers    map[string]string
	onResponse func(resp *http.Response)
}

func (me request) evalClassicURL() string {
	envURL := strings.TrimSuffix(strings.TrimSpace(me.client.Credentials().URL), "/")
	if len(envURL) == 0 {
		return envURL
	}
	envURL = strings.ReplaceAll(envURL, ".dev.apps.dynatracelabs.", ".dev.dynatracelabs.")
	envURL = strings.ReplaceAll(envURL, ".sprint.apps.dynatracelabs.", ".sprint.dynatracelabs.")
	envURL = strings.ReplaceAll(envURL, ".apps.dynatrace.", ".live.dynatrace.")
	return envURL
}

func PreRequest() {
	logging.InstallRoundTripper()
	if strings.TrimSpace(os.Getenv("DYNATRACE_HTTP_INSECURE")) == "true" {
		if transport, ok := http.DefaultTransport.(*http.Transport); ok {
			transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		}
	}
}

func readerFromPayload(payload any) (io.Reader, error) {
	if payload == nil {
		return nil, nil
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(data), nil
}

func UnmarshalError(method string, url string, data []byte, response *http.Response) (err error) {
	if response == nil {
		return nil
	}
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return nil
	}
	var env errorEnvelope
	if err = json.Unmarshal(data, &env); err == nil && env.Error != nil {
		return Error{Code: env.Error.Code, Method: method, URL: url, Message: env.Error.Message, ConstraintViolations: env.Error.ConstraintViolations}
	} else {
		var envs []errorEnvelope
		if err = json.Unmarshal(data, &envs); err == nil && len(envs) > 0 {
			env = envs[0]
			return Error{Code: env.Error.Code, Method: method, URL: url, Message: env.Error.Message, ConstraintViolations: env.Error.ConstraintViolations}
		}
	}
	if len(data) > 0 {
		return Error{Code: response.StatusCode, Method: method, URL: url, Message: string(data)}
	}
	return Error{Code: response.StatusCode, Method: method, URL: url, Message: response.Status}
}

var reManaged = regexp.MustCompile(`^/e/[0-9a-fA-F\-]{36}/`)

func (me request) HandleResponse(client *rest.Client, u *url.URL, target any) (err error) {
	var response *http.Response

	ctx := context.WithValue(me.ctx, "request.id", uuid.NewString())

	uPath := reManaged.ReplaceAllString(u.Path, "/")

	switch me.method {
	case http.MethodGet:
		response, err = client.GET(ctx, uPath, rest.RequestOptions{QueryParams: u.Query()})
	case http.MethodDelete:
		response, err = client.DELETE(ctx, uPath, rest.RequestOptions{QueryParams: u.Query()})
	case http.MethodPost, http.MethodPatch, http.MethodPut:
		body, err := readerFromPayload(me.payload)
		if err != nil {
			return err
		}
		switch me.method {
		case http.MethodPost:
			response, err = client.POST(ctx, uPath, body, rest.RequestOptions{QueryParams: u.Query()})
		case http.MethodPut:
			response, err = client.PUT(ctx, uPath, body, rest.RequestOptions{QueryParams: u.Query()})
		case http.MethodPatch:
			response, err = client.PATCH(ctx, uPath, body, rest.RequestOptions{QueryParams: u.Query()})
		}
	default:
		return fmt.Errorf("unsupported method %s", me.method)
	}

	if err != nil {
		return err
	}
	response = nil
	defer func() {
		if response != nil && response.Body != nil {
			response.Body.Close()
		}
	}()
	if me.onResponse != nil && response != nil {
		me.onResponse(response)
	}

	var data []byte
	if response != nil && response.Body != nil {
		if data, err = io.ReadAll(response.Body); err != nil {
			return err
		}
	}

	if err = UnmarshalError(me.method, u.String(), data, response); err != nil {
		return err
	}

	if target != nil {
		if err = json.Unmarshal(data, &target); err != nil {
			return fmt.Errorf("%s %s: unmarshal error: %s\n%s", me.method, u.String(), err.Error(), string(data))
		}
	}

	return nil
}

var Headers = struct {
	ContentType struct{ ApplicationJSON map[string]string }
}{
	ContentType: struct{ ApplicationJSON map[string]string }{ApplicationJSON: map[string]string{"Content-Type": "application/json"}},
}
