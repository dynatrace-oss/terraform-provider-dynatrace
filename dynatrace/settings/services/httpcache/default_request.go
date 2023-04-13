package httpcache

import (
	"errors"
	"net/http"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
)

func Request(f Finisher) rest.Request {
	return &request{Finisher: f}
}

func List(schemaID string) rest.Request {
	return Request(&ListV1{schemaID})
}

func EmptyList(schemaID string) rest.Request {
	return Request(&NonSupportedV1{schemaID})
}

func Get(schemaID string, id string) rest.Request {
	return Request(&GetV1{SchemaID: schemaID, ID: id})
}

type request struct {
	Finisher Finisher
}

func (me *request) Raw() ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (me *request) Finish(vs ...any) error {
	var v any
	if len(vs) > 0 {
		v = vs[0]
	}
	return me.Finisher.Finish(v)
}

func (me *request) Expect(codes ...int) rest.Request {
	return me
}

func (me *request) Payload(any) rest.Request {
	return me
}

func (me *request) OnResponse(func(resp *http.Response)) rest.Request {
	return me
}
