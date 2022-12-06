package hclgen

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

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

func unquote(w *hclwrite.Body, key string, val *string) {
	strVal := strings.TrimPrefix(*val, "HCL-UNQUOTE-")
	w.SetAttributeRaw(key, hclwrite.Tokens{&hclwrite.Token{Type: hclsyntax.TokenStringLit, Bytes: []byte(" " + strVal)}})
}

func isJSON(s string) bool {
	trimmed := strings.TrimSpace(s)
	if !strings.HasPrefix(trimmed, "{") || !strings.HasSuffix(trimmed, "}") {
		return false
	}
	m := map[string]any{}
	return json.Unmarshal([]byte(s), &m) == nil
}

func toJSONencode(s string, indent string) string {
	if strings.Contains(s, "\n") {
		return " jsonencode(\n" + indent + "  " + strings.ReplaceAll(s, "\n", "\n"+indent+"  ") + ")"
	}
	return " jsonencode(" + s + ")"
}

/*
 */
func (me *primitiveEntry) Write(w *hclwrite.Body, indent string) error {
	cv := ctyVal(me.Value, indent)
	if cv == cty.NilVal {
		return fmt.Errorf("value of type %T for key %s not supported", me.Value, me.Key)
	}
	if strVal, ok := me.Value.(string); ok && strings.HasPrefix(strVal, "HCL-UNQUOTE-") {
		unquote(w, me.Key, &strVal)
	} else if strValP, ok := me.Value.(*string); ok && strValP != nil && strings.HasPrefix(*strValP, "HCL-UNQUOTE-") {
		unquote(w, me.Key, strValP)
	} else if strSliceVal, ok := me.Value.([]string); ok {
		rawVal := "[ "
		sep := ""
		for _, strVal := range strSliceVal {
			if strings.HasPrefix(strVal, "HCL-UNQUOTE-") {
				rawVal = rawVal + sep + strings.TrimPrefix(strVal, "HCL-UNQUOTE-")
			} else {
				rawVal = rawVal + sep + "\"" + strings.ReplaceAll(strVal, "\\", "\\\\") + "\""
			}
			sep = ", "
		}
		rawVal = rawVal + " ]"
		w.SetAttributeRaw(me.Key, hclwrite.Tokens{&hclwrite.Token{Type: hclsyntax.TokenStringLit, Bytes: []byte(" " + rawVal)}})
	} else if strVal, ok := me.Value.(string); ok && isJSON(strVal) {
		w.SetAttributeRaw(me.Key, hclwrite.Tokens{&hclwrite.Token{Type: hclsyntax.TokenStringLit, Bytes: []byte(toJSONencode(strVal, indent))}})
	} else if strValP, ok := me.Value.(*string); ok && strValP != nil && isJSON(*strValP) {
		w.SetAttributeRaw(me.Key, hclwrite.Tokens{&hclwrite.Token{Type: hclsyntax.TokenStringLit, Bytes: []byte(toJSONencode(*strValP, indent))}})
	} else if strVal, ok := me.Value.(string); ok && strings.Contains(strVal, "\n") {
		mlstr := "<<-EOT\n" + indent + "  " + strings.ReplaceAll(strVal, "\n", "\n"+indent+"  ") + "\n" + indent + "EOT"
		w.SetAttributeRaw(me.Key, hclwrite.Tokens{&hclwrite.Token{Type: hclsyntax.TokenStringLit, Bytes: []byte(mlstr)}})
	} else if strVal, ok := me.Value.(string); ok && strings.Count(strVal, "\"") > 3 {
		mlstr := "<<-EOT\n" + indent + "  " + strVal + "\n" + indent + "EOT"
		w.SetAttributeRaw(me.Key, hclwrite.Tokens{&hclwrite.Token{Type: hclsyntax.TokenStringLit, Bytes: []byte(mlstr)}})
	} else if strValP, ok := me.Value.(*string); ok && strValP != nil && strings.Count(*strValP, "\"") > 3 {
		mlstr := "<<-EOT\n" + indent + "  " + *strValP + "\n" + indent + "EOT"
		w.SetAttributeRaw(me.Key, hclwrite.Tokens{&hclwrite.Token{Type: hclsyntax.TokenStringLit, Bytes: []byte(mlstr)}})
	} else if strValP, ok := me.Value.(*string); ok && strValP != nil && strings.Contains(*strValP, "\n") {
		mlstr := "<<-EOT\n" + indent + "  " + *strValP + "\n" + indent + "EOT"
		w.SetAttributeRaw(me.Key, hclwrite.Tokens{&hclwrite.Token{Type: hclsyntax.TokenStringLit, Bytes: []byte(mlstr)}})
	} else if strVal, ok := me.Value.(string); ok {
		w.SetAttributeRaw(me.Key, hclwrite.Tokens{
			&hclwrite.Token{Type: hclsyntax.TokenStringLit, Bytes: []byte{' '}},
			&hclwrite.Token{Type: hclsyntax.TokenOQuote, Bytes: []byte{'"'}},
			&hclwrite.Token{Type: hclsyntax.TokenStringLit, Bytes: escapeQuotedStringLit(strVal)},
			&hclwrite.Token{Type: hclsyntax.TokenCQuote, Bytes: []byte{'"'}},
		})
	} else if strValP, ok := me.Value.(*string); ok {
		w.SetAttributeRaw(me.Key, hclwrite.Tokens{
			&hclwrite.Token{Type: hclsyntax.TokenStringLit, Bytes: []byte{' '}},
			&hclwrite.Token{Type: hclsyntax.TokenOQuote, Bytes: []byte{'"'}},
			&hclwrite.Token{Type: hclsyntax.TokenStringLit, Bytes: escapeQuotedStringLit(*strValP)},
			&hclwrite.Token{Type: hclsyntax.TokenCQuote, Bytes: []byte{'"'}},
		})
	} else {
		w.SetAttributeValue(me.Key, cv)
	}
	return nil
}

func appendRune(b []byte, r rune) []byte {
	l := utf8.RuneLen(r)
	for i := 0; i < l; i++ {
		b = append(b, 0) // make room at the end of our buffer
	}
	ch := b[len(b)-l:]
	utf8.EncodeRune(ch, r)
	return b
}
func escapeQuotedStringLit(s string) []byte {
	res := string(escapeQuotedStringLit0(s))
	res = strings.ReplaceAll(res, "$${data.", "${data.")
	res = strings.ReplaceAll(res, "$${dynatrace_.", "${dynatrace_.")
	return []byte(res)
}

func escapeQuotedStringLit0(s string) []byte {
	if len(s) == 0 {
		return nil
	}
	buf := make([]byte, 0, len(s))
	for i, r := range s {
		switch r {
		case '\n':
			buf = append(buf, '\\', 'n')
		case '\r':
			buf = append(buf, '\\', 'r')
		case '\t':
			buf = append(buf, '\\', 't')
		case '"':
			buf = append(buf, '\\', '"')
		case '\\':
			buf = append(buf, '\\', '\\')
		case '$', '%':
			buf = appendRune(buf, r)
			remain := s[i+1:]
			if len(remain) > 0 && remain[0] == '{' {
				// Double up our template introducer symbol to escape it.
				buf = appendRune(buf, r)
			}
		default:
			if !unicode.IsPrint(r) {
				var fmted string
				if r < 65536 {
					fmted = fmt.Sprintf("\\u%04x", r)
				} else {
					fmted = fmt.Sprintf("\\U%08x", r)
				}
				buf = append(buf, fmted...)
			} else {
				buf = appendRune(buf, r)
			}
		}
	}
	return buf
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
