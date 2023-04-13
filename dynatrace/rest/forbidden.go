package rest

import (
	"fmt"
	"net/http"
)

type Forbidden struct {
	Method string
}

func (me *Forbidden) Raw() ([]byte, error) {
	return nil, fmt.Errorf("%s requests are not supported in migration mode", me.Method)
}

func (me *Forbidden) Finish(vs ...any) error {
	return fmt.Errorf("%s requests are not supported in migration mode", me.Method)
}

func (me *Forbidden) Expect(codes ...int) Request {
	return me
}

func (me *Forbidden) Payload(any) Request {
	return me
}

func (me *Forbidden) OnResponse(func(resp *http.Response)) Request {
	return me
}

type StriclyForbidden struct {
	Method string
	URL    string
}

func (me *StriclyForbidden) Raw() ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func (me *StriclyForbidden) Finish(vs ...any) error {
	fmt.Printf("unexpected request %s %s\n", me.Method, me.URL)
	panic(fmt.Sprintf("unexpected request %s %s", me.Method, me.URL))
	// return fmt.Errorf("unexpected request %s %s", me.method, me.url)
}

func (me *StriclyForbidden) Expect(codes ...int) Request {
	return me
}

func (me *StriclyForbidden) Payload(any) Request {
	return me
}

func (me *StriclyForbidden) OnResponse(func(resp *http.Response)) Request {
	return me
}
