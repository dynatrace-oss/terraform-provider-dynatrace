package hcl

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func equalLineByLine(s1, s2 string) bool {
	parts1 := strings.Split(strings.TrimSpace(s1), "\n")
	parts2 := strings.Split(strings.TrimSpace(s2), "\n")
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

func JSONBytesEqual(b1, b2 []byte) bool {
	var o1 interface{}
	if err := json.Unmarshal(b1, &o1); err != nil {
		return false
	}

	var o2 interface{}
	if err := json.Unmarshal(b2, &o2); err != nil {
		return false
	}

	return reflect.DeepEqual(o1, o2)
}

func SuppressJSONorEOT(k, old, new string, d *schema.ResourceData) bool {
	if JSONStringsEqual(old, new) {
		return true
	}
	return equalLineByLine(old, new)
}
