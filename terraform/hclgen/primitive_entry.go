/**
* @license
* Copyright 2026 Dynatrace LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package hclgen

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/envutils"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type primitiveEntry struct {
	Indent      string
	Key         string
	Optional    bool
	Computed    bool
	BreadCrumbs string
	Value       any
	DefValTrue  bool
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
	m := map[string]any{}
	json.Unmarshal([]byte(s), &m)
	out, _ := json.MarshalIndent(m, "", indent)
	return " jsonencode(" + escapeHCLTemplate(string(out), indent) + ")"
}

// exportReferencePattern matches the resource and data source references the export
// emits. Any other template sequence is escaped, so exported values can't inject
// arbitrary expressions.
var exportReferencePattern = regexp.MustCompile(`^\$\{(?:data\.)?dynatrace_[a-z0-9_]+\.[\p{L}\p{N}_-]+\.[a-z0-9_]+\}`)

// isEscapedIntroducer reports whether the template introducer at s[i] must be doubled.
func isEscapedIntroducer(s string, i int) bool {
	if i+1 >= len(s) || s[i+1] != '{' {
		return false
	}
	return s[i] == '%' || !exportReferencePattern.MatchString(s[i:])
}

func escapeTemplateSequences(s string) string {
	var sb strings.Builder
	sb.Grow(len(s))

	for i := 0; i < len(s); i++ {
		sb.WriteByte(s[i])
		if (s[i] == '$' || s[i] == '%') && isEscapedIntroducer(s, i) {
			sb.WriteByte(s[i])
		}
	}
	return sb.String()
}

// escapeHCLTemplate prepares a multiline raw string for embedding in HCL (heredoc
// body or jsonencode output): it re-indents every line to match indent and escapes
// template sequences.
func escapeHCLTemplate(s string, indent string) string {
	finalString := strings.ReplaceAll(s, "\r\n", "\n"+indent+"  ")
	finalString = strings.ReplaceAll(finalString, "\n", "\n"+indent+"  ")
	finalString = strings.ReplaceAll(finalString, "\r", "\n"+indent+"  ")
	return escapeTemplateSequences(finalString)
}

// wantHeredoc returns true if EOT should be used (newline or heavily quoted string)
func wantHeredoc(s string) bool {
	return envutils.DynatraceHeredoc.Get() && (strings.Contains(s, "\n") || strings.Count(s, `"`) > 3)
}

func heredocBody(body string, indent string) string {
	escaped := escapeHCLTemplate(body, indent)
	marker := heredocMarker(escaped)
	return "<<-" + marker + "\n" + indent + "  " + escaped + "\n" + indent + marker
}

// heredocMarker returns a marker no body line equals; otherwise that line would
// terminate the heredoc and the rest of the body would be parsed as HCL.
func heredocMarker(body string) string {
	lines := map[string]bool{}
	for line := range strings.SplitSeq(body, "\n") {
		lines[strings.TrimSpace(line)] = true
	}
	marker := "EOT"
	for i := 1; lines[marker]; i++ {
		marker = fmt.Sprintf("EOT%d", i)
	}
	return marker
}

// heredocOrTrimmed renders a multiline s as a heredoc (EOT).
// If the string doesn't end with a newline, EOT is within a "trimsuffix"
func heredocOrTrimmed(s string, indent string) string {
	if before, ok := strings.CutSuffix(s, "\n"); ok {
		return heredocBody(before, indent)
	}
	// The closing arg must sit on its own line; a heredoc terminator can't share
	// its line with following tokens.
	return "trimsuffix(" + heredocBody(s, indent) + "\n" + indent + ", \"\\n\")"
}

/*
 */
func (me *primitiveEntry) Write(w *hclwrite.Body, indent string) error {
	cv := ctyVal(me.Value, indent)
	if cv == cty.NilVal {
		return fmt.Errorf("value of type %T for key %s not supported", me.Value, me.Key)
	}
	if strSliceVal, ok := me.Value.([]string); ok {
		rawVal := "[ "
		sep := ""
		for _, strVal := range strSliceVal {
			rawVal = rawVal + sep + "\"" + string(escapeQuotedStringLit(strVal)) + "\""
			sep = ", "
		}
		rawVal = rawVal + " ]"
		w.SetAttributeRaw(me.Key, hclwrite.Tokens{&hclwrite.Token{Type: hclsyntax.TokenStringLit, Bytes: []byte(" " + rawVal)}})
	} else if strVal, ok := me.Value.(string); ok && isJSON(strVal) {
		w.SetAttributeRaw(me.Key, hclwrite.Tokens{&hclwrite.Token{Type: hclsyntax.TokenStringLit, Bytes: []byte(toJSONencode(strVal, indent))}})
	} else if strValP, ok := me.Value.(*string); ok && strValP != nil && isJSON(*strValP) {
		w.SetAttributeRaw(me.Key, hclwrite.Tokens{&hclwrite.Token{Type: hclsyntax.TokenStringLit, Bytes: []byte(toJSONencode(*strValP, indent))}})
	} else if strVal, ok := me.Value.(string); ok && wantHeredoc(strVal) {
		w.SetAttributeRaw(me.Key, hclwrite.Tokens{&hclwrite.Token{Type: hclsyntax.TokenStringLit, Bytes: []byte(heredocOrTrimmed(strVal, indent))}})
	} else if strValP, ok := me.Value.(*string); ok && strValP != nil && wantHeredoc(*strValP) {
		w.SetAttributeRaw(me.Key, hclwrite.Tokens{&hclwrite.Token{Type: hclsyntax.TokenStringLit, Bytes: []byte(heredocOrTrimmed(*strValP, indent))}})
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
	for range l {
		b = append(b, 0) // make room at the end of our buffer
	}
	ch := b[len(b)-l:]
	utf8.EncodeRune(ch, r)
	return b
}
func escapeQuotedStringLit(s string) []byte {
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
			if isEscapedIntroducer(s, i) {
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

func (me *primitiveEntry) IsComputed() bool {
	return me.Computed
}

func (me *primitiveEntry) IsDefault() bool {
	switch rv := me.Value.(type) {
	case bool:
		if me.DefValTrue {
			return rv
		}
		return !rv
	case string:
		return len(rv) == 0
	case *string:
		return rv == nil
	case *bool:
		return rv == nil
	case *int:
		return rv == nil
	case *int32:
		return rv == nil
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
	case *stringMapEntry:
		return true
	}
	return false
}
