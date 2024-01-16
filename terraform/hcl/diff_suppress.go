package hcl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func equalLineByLine(s1, s2 string) bool {
	if len(s1) == 0 && len(s2) == 0 {
		return false
	}
	s1 = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(s1, "\r\n", "\n"), "\r", "\n"))
	s1 = strings.TrimSuffix(s1, "\n")
	s2 = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(s2, "\r\n", "\n"), "\r", "\n"))
	s2 = strings.TrimSuffix(s2, "\n")

	parts1 := strings.Split(s1, "\n")
	parts2 := strings.Split(s2, "\n")
	if len(parts1) != len(parts2) {
		return false
	}
	for idx := range parts1 {
		part1 := strings.TrimSpace(parts1[idx])
		part2 := strings.TrimSpace(parts2[idx])
		if part1 != part2 {
			return false
		}
	}
	return true
}

func SuppressEOT(k, old, new string, d *schema.ResourceData) bool {
	return equalLineByLine(old, new)
}

func JSONStringsEqual(s1, s2 string) bool {
	if len(s1) == 0 && len(s2) == 0 {
		return true
	}
	b1 := bytes.NewBufferString("")
	if err := json.Compact(b1, []byte(s1)); err != nil {
		return false
	}

	b2 := bytes.NewBufferString("")
	if err := json.Compact(b2, []byte(s2)); err != nil {
		return false
	}

	return JSONBytesEqual(b1.Bytes(), b2.Bytes())
}

func compact(v any) any {
	if v == nil {
		return nil
	}
	switch typedV := v.(type) {
	case map[string]any:
		m := map[string]any{}
		for mk, mv := range typedV {
			if mv != nil {
				if cv := compact(mv); cv != nil {
					m[mk] = cv
				}
			}
		}
		if len(m) == 0 {
			return nil
		}
		return m
	case []any:
		s := []any{}
		for _, sv := range typedV {
			if cv := compact(sv); cv != nil {
				s = append(s, cv)
			}
		}
		if len(s) == 0 {
			return nil
		}
		return s
	default:
		return typedV
	}
}

func JSONBytesEqual(b1, b2 []byte) bool {
	var o1 any
	if err := json.Unmarshal(b1, &o1); err != nil {
		return false
	}
	o1 = compact(o1)

	var o2 any
	if err := json.Unmarshal(b2, &o2); err != nil {
		return false
	}
	o2 = compact(o2)
	return DeepEqual(o1, o2, "")
}

func SuppressJSONorEOT(k, old, new string, d *schema.ResourceData) bool {
	if JSONStringsEqual(old, new) {
		return true
	}
	return equalLineByLine(old, new)
}

var DYNATRACE_DASHBOARD_TESTS = len(os.Getenv("DYNATRACE_DASHBOARD_TESTS")) > 0

func Println(path string, message string, b ...bool) {
	if !DYNATRACE_DASHBOARD_TESTS {
		return
	}
	if len(b) > 0 {
		if b[0] {
			return
		}
	}
	logging.File.Println(path, ":", message)
}

func DeepEqual(a any, b any, path string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil && b != nil {
		Println(path, "a == nil && b != nil")
		return false
	}
	if a != nil && b == nil {
		Println(path, "a != nil && b == nil")
		return false
	}
	ta := reflect.TypeOf(a)
	tb := reflect.TypeOf(b)

	if ta != tb {
		Println(path, "reflect.TypeOf(a) != reflect.TypeOf(b)")
		return false
	}

	switch aTyped := a.(type) {
	case bool:
		Println(path, "a != b", aTyped == b.(bool))
		return aTyped == b.(bool)
	case string:
		Println(path, "a != b", aTyped == b.(string))
		return aTyped == b.(string)
	case int:
		Println(path, "a != b", aTyped == b.(int))
		return aTyped == b.(int)
	case int8:
		Println(path, "a != b", aTyped == b.(int8))
		return aTyped == b.(int8)
	case int16:
		Println(path, "a != b", aTyped == b.(int16))
		return aTyped == b.(int16)
	case int32:
		Println(path, "a != b", aTyped == b.(int32))
		return aTyped == b.(int32)
	case int64:
		Println(path, "a != b", aTyped == b.(int64))
		return aTyped == b.(int64)
	case uint:
		Println(path, "a != b", aTyped == b.(uint))
		return aTyped == b.(uint)
	case uint8:
		Println(path, "a != b", aTyped == b.(uint8))
		return aTyped == b.(uint8)
	case uint16:
		Println(path, "a != b", aTyped == b.(uint16))
		return aTyped == b.(uint16)
	case uint32:
		Println(path, "a != b", aTyped == b.(uint32))
		return aTyped == b.(uint32)
	case uint64:
		Println(path, "a != b", aTyped == b.(uint64))
		return aTyped == b.(uint64)
	case float32:
		Println(path, "a != b", aTyped == b.(float32))
		return aTyped == b.(float32)
	case float64:
		Println(path, "a != b", aTyped == b.(float64))
		return aTyped == b.(float64)
	case map[string]any:
		bTyped := b.(map[string]any)
		delete(aTyped, "metricExpressions")
		delete(bTyped, "metricExpressions")
		if len(aTyped) != len(bTyped) {
			Println(path, "MAP len(a) != len(b)")
			if DYNATRACE_DASHBOARD_TESTS {
				aData, _ := json.Marshal(a)
				Println(path, "--- THIS IS A ---")
				Println(path, string(aData))
				bData, _ := json.Marshal(b)
				Println(path, "--- THIS IS B ---")
				Println(path, string(bData))
			}
			return false
		}
		for ka, va := range aTyped {
			// TODO: this is cheating. Essentially we're suppressing the diff even if metricExpressions are different
			// Reason: The REST API ALWAYS delivers that string in a different format
			// Solution: We would need to implement a metric expression parser
			if ka == "metricExpressions" {
				continue
			}
			if _, ok := bTyped[ka]; !ok {
				Println(fmt.Sprintf("%s[%s]", path, ka), fmt.Sprintf("%s not contained in b", ka))
				return false
			}

			if vb, ok := bTyped[ka]; !ok || !DeepEqual(va, vb, fmt.Sprintf("%s[%s]", path, ka)) {
				Println(fmt.Sprintf("%s[%s]", path, ka), "DeepEqual failed")
				return false
			}
		}
		return true
	case []any:
		bTyped := b.([]any)
		if len(aTyped) != len(bTyped) {
			Println(path, "SLICE len(a) != len(b)")
			return false
		}
		m := map[int]bool{}
		for i := 0; i < len(aTyped); i++ {
			m[i] = false
		}
		for ai, aElem := range aTyped {
			found := false
			for i, bElem := range bTyped {
				if m[i] {
					continue
				}
				if DeepEqual(aElem, bElem, fmt.Sprintf("%s[%d]", path, i)) {
					found = true
					m[i] = true
					break
				}
			}
			if !found {
				Println(fmt.Sprintf("%s[%d]", path, ai), "not found in b")
				if DYNATRACE_DASHBOARD_TESTS {
					aData, _ := json.Marshal(a)
					Println(path, "--- THIS IS A ---")
					Println(path, string(aData))
					bData, _ := json.Marshal(b)
					Println(path, "--- THIS IS B ---")
					Println(path, string(bData))
				}
				return false
			}
		}
		return true
	default:
		panic(fmt.Sprintf("Type `%T` not supported yet", aTyped))
	}
}
