package download

import (
	"github.com/dtcookie/hcl"
)

type StandardClient interface {
	GET(id string) (interface{}, error)
	LIST() ([]string, error)
}

type NoListClient interface {
	GET() (interface{}, error)
}

type Resource struct {
	ID         string
	Name       string
	RESTObject hcl.Marshaler
	ReqInter   ReqInter
	UniqueName string
	Variables  map[string]string
}

type ReqInter struct {
	Type    Type
	Message []string
}

func (resource *Resource) Dedup() {
	m := map[string]string{}
	if len(resource.ReqInter.Message) > 0 {
		for _, msg := range resource.ReqInter.Message {
			m[msg] = ""
		}
	}
	res := []string{}
	for k := range m {
		res = append(res, k)
	}
	resource.ReqInter.Message = res
}

type Type string

var InterventionTypes = struct {
	Flawed  Type
	ReqAttn Type
}{
	".flawed",
	".requires_attention",
}

type Resources map[string]*Resource

type ResourceData map[string]Resources
