package hcl2json

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/dtcookie/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type mapdecoder struct {
	m map[string]interface{}
}

func NewMapDecoder(m map[string]interface{}) hcl.Decoder {
	return &mapdecoder{m}
}

type resolver struct {
	source interface{}
}

func (me *resolver) GetOk(key string) (interface{}, bool) {
	source := me.source
	found := false
	for _, part := range strings.Split(key, ".") {
		switch s := source.(type) {
		case map[string]interface{}:
			if source, found = s[part]; !found {
				return nil, false
			}
		case HashSet:
			if part == "#" {
				return s.Len(), true
			} else if idx, err := strconv.Atoi(part); err == nil {
				source = s.Get(idx)
			}
		case ListSet:
			if part == "#" {
				return s.Len(), true
			} else if idx, err := strconv.Atoi(part); err == nil {
				source = s.Get(idx)
			}
		case []interface{}:
			if part == "#" {
				return len(s), true
			} else if idx, err := strconv.Atoi(part); err == nil {
				source = s[idx]
			}
		case []string:
			if part == "#" {
				return len(s), true
			} else if idx, err := strconv.Atoi(part); err == nil {
				source = s[idx]
			}
		case []int:
			if part == "#" {
				return len(s), true
			} else if idx, err := strconv.Atoi(part); err == nil {
				source = s[idx]
			}
		case []float64:
			if part == "#" {
				return len(s), true
			} else if idx, err := strconv.Atoi(part); err == nil {
				source = s[idx]
			}
		default:
			return nil, false
		}

	}
	return source, true
}

func (me *mapdecoder) GetOk(key string) (interface{}, bool) {
	res := &resolver{me.m}
	r, e := res.GetOk(key)
	switch tr := r.(type) {
	case ListSet:
		r = tr.List()
	case HashSet:
		r = schema.NewSet(HashJSON, tr.List())
	}
	// if !strings.Contains(key, "tile") {
	// fmt.Println("GetOk", key, "=>", r, e)
	// }
	return r, e
}

func (me *mapdecoder) Get(key string) interface{} {
	panic("mapdecoder.Get not implemented")
}

func (me *mapdecoder) GetChange(key string) (interface{}, interface{}) {
	if r, ok := me.GetOk(key); ok {
		return r, r
	}
	return nil, nil
}

func (me *mapdecoder) GetStringSet(key string) []string {
	result := []string{}
	if value, ok := me.GetOk(key); ok {
		if arr, ok := value.([]interface{}); ok {
			for _, elem := range arr {
				result = append(result, elem.(string))
			}
		} else if set, ok := value.(hcl.Set); ok {
			if set.Len() > 0 {
				for _, elem := range set.List() {
					result = append(result, elem.(string))
				}
			}
		}
	}
	return result
}

func (me *mapdecoder) GetOkExists(key string) (interface{}, bool) {
	return me.GetOk(key)
}

func (me *mapdecoder) Reader(unkowns ...map[string]json.RawMessage) hcl.Reader {
	panic("mapdecoder.Reader not implemented")
}

func (me *mapdecoder) HasChange(key string) bool {
	panic("mapdecoder.HasChange not implemented")
}

func (me *mapdecoder) MarshalAll(items map[string]interface{}) (hcl.Properties, error) {
	panic("mapdecoder.MarshalAll not implemented")
}

func (me *mapdecoder) Decode(key string, v interface{}) error {
	return hcl.NewDecoder(me).Decode(key, v)
}

func (me *mapdecoder) DecodeAll(m map[string]interface{}) error {
	return hcl.NewDecoder(me).DecodeAll(m)
}

func (me *mapdecoder) DecodeAny(map[string]interface{}) (interface{}, error) {
	panic("mapdecoder.DecodeAny not implemented")
}

func (me *mapdecoder) DecodeSlice(key string, v interface{}) error {
	panic("mapdecoder.DecodeSlice not implemented")
}
