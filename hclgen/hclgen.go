package hclgen

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"

	"github.com/dtcookie/hcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type HCLGen struct {
	file *hclwrite.File
}

func New() *HCLGen {
	return &HCLGen{file: hclwrite.NewEmptyFile()}
}

func sk(s string) string {
	if s == "name" {
		return "00" + s
	}
	if s == "description" {
		return "10" + s
	}
	if s == "type" {
		return "20" + s
	}
	if s == "enabled" {
		return "30" + s
	}
	return "90" + s
}

func resOpt(bc string, sch map[string]*hcl.Schema) bool {
	bc = strings.TrimPrefix(bc, ".")
	if strings.Contains(bc, ".") {
		idx := strings.Index(bc, ".")
		return resOpt0(bc[:idx], bc[idx+1:], sch[bc[:idx]])
	}
	return resOpt0(bc, "", sch[bc])
}

func resOpt0(key string, bc string, sch *hcl.Schema) bool {
	if sch == nil {
		return false
	}
	switch sch.Type {
	case hcl.TypeBool:
		return sch.Optional
	case hcl.TypeInt:
		return sch.Optional
	case hcl.TypeFloat:
		return sch.Optional
	case hcl.TypeString:
		return sch.Optional
	case hcl.TypeList:
		switch v := sch.Elem.(type) {
		case *hcl.Resource:
			return resOpt(bc, v.Schema)
		default:
			return sch.Optional
		}

	case hcl.TypeMap:
		return false
	case hcl.TypeSet:
		return false
	default:
		return false

	}
}

type exportEntry interface {
	Write(w *hclwrite.Body, indent string) error
	IsOptional() bool
	IsDefault() bool
	IsLessThan(other exportEntry) bool
}

type exportEntries []exportEntry

func (e exportEntries) Less(i, j int) bool {
	return e[i].IsLessThan(e[j])
}

func (e *exportEntries) eval(key string, value interface{}, breadCrumbs string, schema map[string]*hcl.Schema) {
	if value == nil {
		return
	}
	switch v := value.(type) {
	case string, bool, int, int32, int64, int8, int16, uint, uint32, uint64, uint8, uint16, float32, float64:
		entry := &primitiveEntry{Key: key, Value: value, BreadCrumbs: breadCrumbs, Optional: resOpt(breadCrumbs, schema)}
		*e = append(*e, entry)
	case *string, *bool, *int, *int32, *int64, *int8, *int16, *uint, *uint32, *uint64, *uint8, *uint16, *float32, *float64:
		if v == nil {
			return
		}
		entry := &primitiveEntry{Key: key, Value: v, BreadCrumbs: breadCrumbs, Optional: resOpt(breadCrumbs, schema)}
		*e = append(*e, entry)
	case []interface{}:
		if len(v) == 0 {
			return
		}
		switch typedElem := v[0].(type) {
		case map[string]interface{}:
			for _, elem := range v {
				entry := &resourceEntry{Key: key, Entries: exportEntries{}}
				entry.Entries.handle(elem.(map[string]interface{}), breadCrumbs, schema)
				*e = append(*e, entry)
			}
		case string:
			vs := []string{}
			for _, elem := range v {
				vs = append(vs, elem.(string))
			}
			entry := &primitiveEntry{Key: key, Value: vs}
			*e = append(*e, entry)
		case bool:
			vs := []bool{}
			for _, elem := range v {
				vs = append(vs, elem.(bool))
			}
			entry := &primitiveEntry{Key: key, Value: vs}
			*e = append(*e, entry)
		case int:
			vs := []int{}
			for _, elem := range v {
				vs = append(vs, elem.(int))
			}
			entry := &primitiveEntry{Key: key, Value: vs}
			*e = append(*e, entry)
		case int32:
			vs := []int32{}
			for _, elem := range v {
				vs = append(vs, elem.(int32))
			}
			entry := &primitiveEntry{Key: key, Value: vs}
			*e = append(*e, entry)
		case int64:
			vs := []int64{}
			for _, elem := range v {
				vs = append(vs, elem.(int64))
			}
			entry := &primitiveEntry{Key: key, Value: vs}
			*e = append(*e, entry)
		case int8:
			vs := []int8{}
			for _, elem := range v {
				vs = append(vs, elem.(int8))
			}
			entry := &primitiveEntry{Key: key, Value: vs}
			*e = append(*e, entry)
		case int16:
			vs := []int16{}
			for _, elem := range v {
				vs = append(vs, elem.(int16))
			}
			entry := &primitiveEntry{Key: key, Value: vs}
			*e = append(*e, entry)
		case uint:
			vs := []uint{}
			for _, elem := range v {
				vs = append(vs, elem.(uint))
			}
			entry := &primitiveEntry{Key: key, Value: vs}
			*e = append(*e, entry)
		case uint32:
			vs := []uint32{}
			for _, elem := range v {
				vs = append(vs, elem.(uint32))
			}
			entry := &primitiveEntry{Key: key, Value: vs}
			*e = append(*e, entry)
		case uint64:
			vs := []uint64{}
			for _, elem := range v {
				vs = append(vs, elem.(uint64))
			}
			entry := &primitiveEntry{Key: key, Value: vs}
			*e = append(*e, entry)
		case uint8:
			vs := []uint8{}
			for _, elem := range v {
				vs = append(vs, elem.(uint8))
			}
			entry := &primitiveEntry{Key: key, Value: vs}
			*e = append(*e, entry)
		case uint16:
			vs := []uint16{}
			for _, elem := range v {
				vs = append(vs, elem.(uint16))
			}
			entry := &primitiveEntry{Key: key, Value: vs}
			*e = append(*e, entry)
		case float32:
			vs := []float32{}
			for _, elem := range v {
				vs = append(vs, elem.(float32))
			}
			entry := &primitiveEntry{Key: key, Value: vs}
			*e = append(*e, entry)
		case float64:
			vs := []float64{}
			for _, elem := range v {
				vs = append(vs, elem.(float64))
			}
			entry := &primitiveEntry{Key: key, Value: vs}
			*e = append(*e, entry)
		default:
			panic(fmt.Sprintf("unsupported elem type %T", typedElem))
		}
	case []string:
		if len(v) == 0 {
			return
		}
		entry := &primitiveEntry{Key: key, Value: value, Optional: resOpt(breadCrumbs, schema)}
		*e = append(*e, entry)
	case hcl.StringSet:
		if len(v) == 0 {
			return
		}
		entry := &primitiveEntry{Key: key, Value: value, Optional: resOpt(breadCrumbs, schema)}
		*e = append(*e, entry)
	case []float64:
		if len(v) == 0 {
			return
		}
		entry := &primitiveEntry{Key: key, Value: value, Optional: resOpt(breadCrumbs, schema)}
		*e = append(*e, entry)
	case map[string]interface{}:
		if len(v) == 0 {
			return
		}
		entry := &resourceEntry{Key: key, Entries: exportEntries{}}
		for xk, xv := range v {
			entry.Entries = append(entry.Entries, &primitiveEntry{Key: xk, Value: xv, Optional: resOpt(breadCrumbs, schema)})
		}
		*e = append(*e, entry)
	default:
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.String:
			e.eval(key, fmt.Sprintf("%v", v), breadCrumbs, schema)
		default:
			panic(fmt.Sprintf(">>>>> [%v] type %T not supported yet\n", key, v))
		}

	}
}

func (e *exportEntries) handle(m map[string]interface{}, breadCrumbs string, schema map[string]*hcl.Schema) {
	for k, v := range m {
		e.eval(k, v, breadCrumbs+"."+k, schema)
	}
}

func Export(marshaler hcl.Marshaler, w io.Writer, resourceType string, resourceName string, comments ...string) error {
	return New().Export(marshaler, w, resourceType, resourceName, comments...)
}

func (me *HCLGen) Export(marshaler hcl.Marshaler, w io.Writer, resourceType string, resourceName string, comments ...string) error {
	var m map[string]interface{}
	var err error
	if m, err = marshaler.MarshalHCL(); err != nil {
		return err
	}
	var schema map[string]*hcl.Schema
	if schemer, ok := marshaler.(hcl.Schemer); ok {
		schema = schemer.Schema()
	}
	return me.export(m, schema, w, resourceType, resourceName, comments...)
}

func (me *HCLGen) export(m map[string]interface{}, schema map[string]*hcl.Schema, w io.Writer, resourceType string, resourceName string, comments ...string) error {
	var err error
	ents := exportEntries{}
	ents.handle(m, "", schema)
	sort.SliceStable(ents, ents.Less)
	rootBody := me.file.Body()
	tokens := hclwrite.Tokens{}
	if len(comments) > 0 {
		for _, comment := range comments {
			comment = strings.TrimSpace(comment)
			if len(comment) == 0 {
				tokens = append(tokens, &hclwrite.Token{Type: hclsyntax.TokenComment, Bytes: []byte("# Exported by HCLGEN\n")})
			}
		}
	}
	if len(tokens) > 0 {
		rootBody.AppendUnstructuredTokens(tokens)
	}
	bs := rootBody.AppendNewBlock(
		"resource",
		[]string{
			resourceType,
			resourceName,
		},
	)
	body := bs.Body()
	for _, entry := range ents {
		if !(entry.IsOptional() && entry.IsDefault()) {
			if err := entry.Write(body, "  "); err != nil {
				return err
			}
		} else {
			body.AppendUnstructuredTokens(hclwrite.Tokens{
				&hclwrite.Token{Type: hclsyntax.TokenComment, Bytes: []byte("#")},
			})
			if err := entry.Write(body, "  "); err != nil {
				return err
			}
		}
	}
	w.Write(me.file.Bytes())
	return err
}
