package hclgen

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type primitiveEntry struct {
	Indent      string
	Key         string
	Optional    bool
	BreadCrumbs string
	Value       interface{}
}

func (me *primitiveEntry) Write(w *hclwrite.Body, indent string) error {
	cv := ctyVal(me.Value, indent)
	if cv == cty.NilVal {
		return fmt.Errorf("value of type %T for key %s not supported", me.Value, me.Key)
	}
	if strVal, ok := me.Value.(string); ok && strings.HasPrefix(strVal, "HCL-UNQUOTE-") {
		strVal = strings.TrimPrefix(strVal, "HCL-UNQUOTE-")
		w.SetAttributeRaw(me.Key, hclwrite.Tokens{&hclwrite.Token{Type: hclsyntax.TokenStringLit, Bytes: []byte(" " + strVal)}})
	} else if strValP, ok := me.Value.(*string); ok && strValP != nil && strings.HasPrefix(*strValP, "HCL-UNQUOTE-") {
		strVal = strings.TrimPrefix(*strValP, "HCL-UNQUOTE-")
		w.SetAttributeRaw(me.Key, hclwrite.Tokens{&hclwrite.Token{Type: hclsyntax.TokenStringLit, Bytes: []byte(" " + strVal)}})
	} else {
		w.SetAttributeValue(me.Key, cv)
	}
	return nil
}

func (me *primitiveEntry) IsOptional() bool {
	return me.Optional
}

func (me *primitiveEntry) IsDefault() bool {
	switch rv := me.Value.(type) {
	case bool:
		return !rv
	case string:
		return len(rv) == 0
	default:
		return false
	}
}

func (me *primitiveEntry) IsLessThan(other exportEntry) bool {
	switch ro := other.(type) {
	case *primitiveEntry:
		return strings.Compare(sk(me.Key), sk(ro.Key)) < 0
	case *resourceEntry:
		return true
	}
	return false
}
