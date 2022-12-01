package download

import "github.com/dtcookie/hcl"

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
