package hcl2sdk

import (
	"github.com/dtcookie/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Marshal(config hcl.Marshaler) (map[string]interface{}, error) {
	m, err := config.MarshalHCL()
	if err != nil {
		return nil, err
	}
	after := postProcess(m).(map[string]interface{})
	return after, nil
}

func postProcess(v interface{}) interface{} {
	if v == nil {
		return v
	}
	switch rv := v.(type) {
	case map[string]interface{}:
		for mk, mv := range rv {
			rv[mk] = postProcess(mv)
		}
		return rv
	case []interface{}:
		for idx, av := range rv {
			rv[idx] = postProcess(av)
		}
		return rv
	case hcl.StringSet:
		items := []interface{}{}
		for _, elem := range rv {
			items = append(items, elem)
		}
		return schema.NewSet(schema.HashString, items)
	default:
		return v
	}
}
