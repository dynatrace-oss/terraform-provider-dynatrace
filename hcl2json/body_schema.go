package hcl2json

import (
	"encoding/json"

	"github.com/hashicorp/hcl/v2"
)

type bodySchema struct {
	hcl.BodySchema
	Prototypes    map[string]interface{}
	BlockSchemata map[string]*bodySchema
	TypeSet       bool
}

func (me *bodySchema) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{}
	if len(me.Prototypes) > 0 {
		for k, prototype := range me.Prototypes {
			m[k] = prototype
		}
	}
	if len(me.BlockSchemata) > 0 {
		for k, bs := range me.BlockSchemata {
			m[k] = bs
		}
	}
	return json.Marshal(m)
}
