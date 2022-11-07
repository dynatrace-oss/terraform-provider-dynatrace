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
	RESTObject hcl.Marshaler
	Name       string
}

type Resources []Resource

type ResourceData map[string]Resources
