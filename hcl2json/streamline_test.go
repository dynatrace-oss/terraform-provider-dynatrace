package hcl2json_test

import (
	"encoding/json"
	"sort"
)

func streamline(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	switch tv := v.(type) {
	case bool:
		return nil
	case string:
		if len(tv) == 0 {
			return nil
		}
	case map[string]interface{}:
		for k, v := range tv {
			av := streamline(v)
			if av == nil {
				delete(tv, k)
			} else {
				tv[k] = av
			}
		}
		if len(tv) == 0 {
			return nil
		}
	case []string:
		if len(tv) == 0 {
			return nil
		}
		sort.Strings(tv)
	case []int:
		if len(tv) == 0 {
			return nil
		}
	case []float64:
		if len(tv) == 0 {
			return nil
		}
	case []interface{}:
		if len(tv) == 0 {
			return nil
		}
		jsons := []string{}
		for _, v := range tv {
			byt, _ := json.Marshal(streamline(v))
			jsons = append(jsons, string(byt))
		}
		sort.Strings(jsons)
		for idx, je := range jsons {
			json.Unmarshal([]byte(je), &tv[idx])
		}
	}
	return v
}
