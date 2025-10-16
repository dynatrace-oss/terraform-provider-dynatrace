package rest

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/envutil"
	"strings"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/version"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/clients"
	"github.com/google/uuid"
)

type cluster_request request

var clusterClientCache = map[string]*rest.Client{}

var clusterClientCacheMutex sync.Mutex

func clusterClient(baseURL string, apiToken string) (*rest.Client, error) {
	clusterClientCacheMutex.Lock()
	defer clusterClientCacheMutex.Unlock()

	if client, found := clusterClientCache[baseURL]; found {
		return client, nil
	}

	factory := clients.Factory()
	factory = factory.WithUserAgent(version.UserAgent())
	factory = factory.WithClassicURL(baseURL)
	factory = factory.WithAccessToken(apiToken)
	factory = factory.WithHTTPListener(HTTPListener("classic"))

	client, err := factory.CreateClassicClient()

	clusterClientCache[baseURL] = client

	return client, err
}

func (me *cluster_request) Finish(vs ...any) error {
	var v any
	if len(vs) > 0 {
		v = vs[0]
	}

	ctx := context.Background()

	InstallRoundTripper()
	if envutil.GetBoolEnv(envutil.EnvHTTPInsecure, false) {
		if transport, ok := http.DefaultTransport.(*http.Transport); ok {
			transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		}
	}

	classicURL := me.evalClassicURL()

	client, err := clusterClient(classicURL, me.client.Credentials().Cluster.Token)
	if err != nil {
		return err
	}

	fullURL, err := url.Parse(classicURL + me.url)
	if err != nil {
		return err
	}

	var response *http.Response

	ctx = context.WithValue(ctx, "request.id", uuid.NewString())

	switch me.method {
	case http.MethodGet:
		response, err = client.GET(ctx, fullURL.Path, rest.RequestOptions{QueryParams: fullURL.Query()})
	case http.MethodDelete:
		response, err = client.DELETE(ctx, fullURL.Path, rest.RequestOptions{QueryParams: fullURL.Query()})
	case http.MethodPost, http.MethodPatch, http.MethodPut:
		body, err := readerFromPayload(me.payload)
		if err != nil {
			return err
		}
		switch me.method {
		case http.MethodPost:
			response, err = client.POST(ctx, fullURL.Path, body, rest.RequestOptions{QueryParams: fullURL.Query()})
		case http.MethodPut:
			response, err = client.PUT(ctx, fullURL.Path, body, rest.RequestOptions{QueryParams: fullURL.Query()})
		case http.MethodPatch:
			response, err = client.PATCH(ctx, fullURL.Path, body, rest.RequestOptions{QueryParams: fullURL.Query()})
		}
	default:
		return fmt.Errorf("unsupported method %s", me.method)
	}

	if err != nil {
		return err
	}
	defer response.Body.Close()
	if me.onResponse != nil {
		me.onResponse(response)
	}
	var data []byte
	if data, err = io.ReadAll(response.Body); err != nil {
		return err
	}
	_ = data
	if v != nil {
		if err = json.Unmarshal(data, &v); err != nil {
			return fmt.Errorf("%s %s: unmarshal error: %s\n%s", me.method, me.url, err.Error(), string(data))
		}
	}

	return nil
}

func (me *cluster_request) Payload(payload any) Request {
	me.payload = payload
	return me
}

func (me *cluster_request) Expect(codes ...int) Request {
	me.expect = statuscodes(codes)
	return me
}

func (me *cluster_request) OnResponse(onResponse func(resp *http.Response)) Request {
	me.onResponse = onResponse
	return me
}

func (me *cluster_request) evalClassicURL() string {
	envURL := strings.TrimSuffix(strings.TrimSpace(me.client.Credentials().URL), "/")
	if len(envURL) == 0 {
		return envURL
	}
	envURL = strings.ReplaceAll(envURL, ".dev.apps.dynatracelabs.", ".dev.dynatracelabs.")
	envURL = strings.ReplaceAll(envURL, ".sprint.apps.dynatracelabs.", ".sprint.dynatracelabs.")
	envURL = strings.ReplaceAll(envURL, ".apps.dynatrace.", ".live.dynatrace.")
	return envURL
}

func (me *cluster_request) SetHeader(name string, value string) {
	if me.headers == nil {
		me.headers = map[string]string{}
	}
	me.headers[name] = value
}
