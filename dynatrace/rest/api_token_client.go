package rest

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest/logging"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/clients"
	"github.com/google/uuid"
)

func APITokenClient(credentials *Credentials) Client {
	return &api_token_client{credentials: credentials}
}

type api_token_client struct {
	credentials *Credentials
}

func (me *api_token_client) Credentials() *Credentials {
	return me.credentials
}

func (me *api_token_client) Get(ctx context.Context, url string, expectedStatusCodes ...int) Request {
	req := &api_token_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: http.MethodGet}
	if len(expectedStatusCodes) > 0 {
		req.expect = statuscodes(expectedStatusCodes)
	}
	return req
}

func (me *api_token_client) Post(ctx context.Context, url string, payload any, expectedStatusCodes ...int) Request {
	req := &api_token_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: http.MethodPost, payload: payload, headers: map[string]string{"Content-Type": "application/json"}}
	if len(expectedStatusCodes) > 0 {
		req.expect = statuscodes(expectedStatusCodes)
	}
	return req
}

func (me *api_token_client) Put(ctx context.Context, url string, payload any, expectedStatusCodes ...int) Request {
	req := &api_token_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: http.MethodPut, payload: payload, headers: map[string]string{"Content-Type": "application/json"}}
	if len(expectedStatusCodes) > 0 {
		req.expect = statuscodes(expectedStatusCodes)
	}
	return req
}

func (me *api_token_client) Delete(ctx context.Context, url string, expectedStatusCodes ...int) Request {
	req := &api_token_request{id: uuid.NewString(), ctx: ctx, client: me, url: url, method: http.MethodDelete}
	if len(expectedStatusCodes) > 0 {
		req.expect = statuscodes(expectedStatusCodes)
	}
	return req
}

func (me *api_token_request) Finish(vs ...any) error {
	credentials := me.client.Credentials()
	if !credentials.ContainsAPIToken() {
		return NoAPITokenError
	}
	if DYNATRACE_HTTP_LEGACY {
		legacyRequest := legacy_request(*me)
		if credentials.URL == TestCaseEnvURL {
			return errors.New("legacy")
		}
		return legacyRequest.Finish(vs...)
	}
	classicRequest := classic_request(*me)
	if credentials.URL == TestCaseEnvURL {
		return errors.New("classic")
	}
	return classicRequest.Finish(vs...)
}

type api_token_request request

func (me *api_token_request) Expect(codes ...int) Request {
	me.expect = statuscodes(codes)
	return me
}

func (me *api_token_request) OnResponse(onResponse func(resp *http.Response)) Request {
	me.onResponse = onResponse
	return me
}

type classic_request request

var classicClientCache = map[string]*rest.Client{}

var classicClientCacheMutex sync.Mutex

func classicClient(classicURL string, apiToken string) (*rest.Client, error) {
	classicClientCacheMutex.Lock()
	defer classicClientCacheMutex.Unlock()

	if client, found := classicClientCache[classicURL]; found {
		return client, nil
	}

	factory := clients.Factory()
	factory = factory.WithUserAgent("Dynatrace Terraform Provider")
	factory = factory.WithClassicURL(classicURL)
	factory = factory.WithAccessToken(apiToken)
	factory = factory.WithHTTPListener(logging.HTTPListener("classic "))

	client, err := factory.CreateClassicClient()

	classicClientCache[classicURL] = client

	return client, err
}

func (me *classic_request) Finish(optionalTarget ...any) error {
	var target any
	if len(optionalTarget) > 0 {
		target = optionalTarget[0]
	}

	preRequest()

	classicURL := request(*me).evalClassicURL()

	client, err := classicClient(classicURL, me.client.Credentials().Token)
	if err != nil {
		return err
	}

	fullURL, err := url.Parse(classicURL + me.url)
	if err != nil {
		return err
	}
	return request(*me).HandleResponse(client, fullURL, target)
}
