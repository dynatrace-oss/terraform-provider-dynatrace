package hcl2json

import (
	"encoding/json"
	"fmt"
	"reflect"

	dtchcl "github.com/dtcookie/hcl"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

func resource2Schema(res *schema.Resource, bc breadCrumbs) *bodySchema {
	var result = bodySchema{Prototypes: map[string]interface{}{}, BlockSchemata: map[string]*bodySchema{}}
	for name, sch := range res.Schema {
		switch sch.Type {
		case schema.TypeBool, schema.TypeInt, schema.TypeFloat, schema.TypeString:
			result.Attributes = append(result.Attributes, hcl.AttributeSchema{Name: name, Required: sch.Required})
			switch sch.Type {
			case schema.TypeBool:
				result.Prototypes[name] = false
			case schema.TypeInt:
				result.Prototypes[name] = int64(0)
			case schema.TypeFloat:
				result.Prototypes[name] = float64(0)
			case schema.TypeString:
				result.Prototypes[name] = ""
			}
			// fmt.Printf("%v = %T\n", bc.dot(name), result.Prototypes[name])
		case schema.TypeList:
			switch elem := sch.Elem.(type) {
			case *schema.Schema:
				switch elem.Type {
				case schema.TypeString:
					result.Prototypes[name] = []string{}
					result.Attributes = append(result.Attributes, hcl.AttributeSchema{Name: name, Required: sch.Required})
					// bc.dot(name).log("[]string")
				default:
					// bc.dot(name).log(fmt.Sprintf(" elem of type *schema.Schema(Type: %v) - not dealing with it yet", elem.Type))
				}
			case *schema.Resource:
				blockHeaderSchema := hcl.BlockHeaderSchema{Type: name}
				result.Blocks = append(result.Blocks, blockHeaderSchema)
				result.BlockSchemata[name] = resource2Schema(elem, bc.dot(name))
			}
		case schema.TypeSet:
			switch elem := sch.Elem.(type) {
			case *schema.Schema:
				switch elem.Type {
				case schema.TypeString:
					result.Prototypes[name] = stringSet{}
					result.Attributes = append(result.Attributes, hcl.AttributeSchema{Name: name, Required: sch.Required})
					// bc.dot(name).log("[]string")
				default:
					// bc.dot(name).log(fmt.Sprintf(" elem of type *schema.Schema(Type: %v) - not dealing with it yet", elem.Type))
				}
			case *schema.Resource:
				blockHeaderSchema := hcl.BlockHeaderSchema{Type: name}
				result.Blocks = append(result.Blocks, blockHeaderSchema)
				result.BlockSchemata[name] = resource2Schema(elem, bc.dot(name))
				result.BlockSchemata[name].TypeSet = true
			}
		}

	}
	return &result
}

type stringSet []string

func (me stringSet) List() []interface{} {
	r := []interface{}{}
	for _, v := range me {
		r = append(r, v)
	}
	return r
}

func (me stringSet) Len() int {
	return len(me)
}

type Set interface {
	Append(interface{}) Set
	List() []interface{}
	Len() int
	Get(int) interface{}
}

type ListSet []interface{}

func (me ListSet) Append(v interface{}) Set {
	return append(me, v)
}

func (me ListSet) Get(idx int) interface{} {
	return me[idx]
}

func (me ListSet) List() []interface{} {
	return me
}

func (me ListSet) Len() int {
	return len(me)
}

type HashSet map[int]interface{}

func (me HashSet) Append(v interface{}) Set {
	me[HashJSON(v)] = v
	return me
}

func (me HashSet) Get(idx int) interface{} {
	return me[idx]
}

func (me HashSet) List() []interface{} {
	res := []interface{}{}
	for _, v := range me {
		res = append(res, v)
	}
	return res
}

func (me HashSet) Len() int {
	return len(me)
}

func HashJSON(v interface{}) int {
	data, _ := json.Marshal(v)
	hash := schema.HashString(string(data))
	// fmt.Println("hash", hash, string(data))
	return hash

}

func translate(body hcl.Body, sch *bodySchema, bc breadCrumbs, target map[string]interface{}) error {
	bodyContent, diag := body.Content(&sch.BodySchema)
	if diag.HasErrors() {
		return &diag
	}

	for _, attribute := range bodyContent.Attributes {

		val, diag := attribute.Expr.Value(&hcl.EvalContext{Variables: map[string]cty.Value{
			"data":                  cty.DynamicVal,
			"dynatrace_dashboard":   cty.DynamicVal,
			"dynatrace_credentials": cty.DynamicVal,
		}})
		if diag.HasErrors() {
			fmt.Println("error in attribute.Expr.Value", diag)
			return &diag
		}
		value := sch.Prototypes[attribute.Name]
		switch typedValue := value.(type) {
		case stringSet:
			if val.Type().IsTupleType() {
				vals := val.AsValueSlice()
				if len(vals) > 0 && !vals[0].IsKnown() {
					val = cty.ListVal([]cty.Value{cty.StringVal("UNKNOWN")})
				} else {
					val = cty.ListVal(vals)
				}
			}
			if err := gocty.FromCtyValue(val, &typedValue); err != nil {
				return fmt.Errorf("%v[%s] - %s: %T, %v", bc, attribute.Name, err.Error(), value, val.GoString())
			}
			target[attribute.Name] = typedValue
		case []string:
			vals := val.AsValueSlice()
			if len(vals) > 0 && !vals[0].IsKnown() {
				val = cty.ListVal([]cty.Value{cty.StringVal("UNKNOWN")})
			} else {
				val = cty.ListVal(vals)
			}
			// val = cty.ListVal(vals)

			// if !val.IsKnown() {
			// 	val = cty.ListVal([]cty.Value{cty.StringVal("UNKNOWN")})
			// }
			// if val.Type().IsTupleType() {
			// 	val = cty.ListVal(val.AsValueSlice())
			// }
			if err := gocty.FromCtyValue(val, &typedValue); err != nil {
				return fmt.Errorf("%v[%s] - %s: %T, %v", bc, attribute.Name, err.Error(), value, val.GoString())
			}
			is := []interface{}{}
			for _, s := range typedValue {
				is = append(is, s)
			}
			target[attribute.Name] = is
		case string:
			if !val.IsKnown() {
				val = cty.StringVal("UNKNOWN")
			}
			if err := gocty.FromCtyValue(val, &typedValue); err != nil {
				return fmt.Errorf("%v[%s] - %s: %T, %v", bc, attribute.Name, err.Error(), value, val.GoString())
			}
			target[attribute.Name] = typedValue
		case bool:
			if err := gocty.FromCtyValue(val, &typedValue); err != nil {
				return fmt.Errorf("%v[%s] - %s: %T, %v", bc, attribute.Name, err.Error(), value, val.GoString())
			}
			target[attribute.Name] = typedValue
		case int64:
			if err := gocty.FromCtyValue(val, &typedValue); err != nil {
				return fmt.Errorf("%v[%s] - %s: %T, %v", bc, attribute.Name, err.Error(), value, val.GoString())
			}
			target[attribute.Name] = int(typedValue)
		case float64:
			if err := gocty.FromCtyValue(val, &typedValue); err != nil {
				return fmt.Errorf("%v[%s] - %s: %T, %v", bc, attribute.Name, err.Error(), value, val.GoString())
			}
			target[attribute.Name] = typedValue
		default:
			return fmt.Errorf("unsupported type %T", typedValue)
		}
	}
	for _, block := range bodyContent.Blocks {
		blockSchema := sch.BlockSchemata[block.Type]
		var blockEntries Set
		untypedBlockEntries, ok := target[block.Type]
		if !ok {
			if blockSchema.TypeSet {
				blockEntries = HashSet{}
			} else {
				blockEntries = ListSet{}
			}
			target[block.Type] = blockEntries
		} else {
			switch typedBlockEntries := untypedBlockEntries.(type) {
			case ListSet:
				blockEntries = typedBlockEntries
			case HashSet:
				blockEntries = typedBlockEntries
			case []interface{}:
				blockEntries = ListSet(typedBlockEntries)
			case Set:
				blockEntries = typedBlockEntries
			}
		}
		blockMap := map[string]interface{}{}
		if err := translate(block.Body, blockSchema, bc.dot(block.Type), blockMap); err != nil {
			return err
		}
		blockEntries = blockEntries.Append(blockMap)
		target[block.Type] = blockEntries
	}
	return nil
}

func buildSchemata() map[string]*bodySchema {
	schemata := map[string]*bodySchema{}
	for k, v := range provider.Provider().ResourcesMap {
		schemata[k] = resource2Schema(v, breadCrumbs(""))
	}
	return schemata
}

func HCL2Config(fileName string) ([]interface{}, error) {
	schemata := buildSchemata()
	parser := hclparse.NewParser()
	hclFile, diag := parser.ParseHCLFile(fileName)
	if diag.HasErrors() {
		return nil, diag
	}
	schema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "resource",
				LabelNames: []string{"resource", "name"},
			},
		},
	}
	bodyContent, diag := hclFile.Body.Content(schema)
	if diag.HasErrors() {
		return nil, diag
	}
	var result []interface{}
	for _, block := range bodyContent.Blocks {
		switch block.Type {
		case "resource":
			m := map[string]interface{}{}
			resource := block.Labels[0]
			bs := schemata[resource]
			if err := translate(block.Body, bs, breadCrumbs(block.Labels[0]), m); err != nil {
				return nil, err
			}
			decoder := NewMapDecoder(m)
			config := reflect.New(reflect.TypeOf(protoTypes[resource]).Elem()).Interface()
			if err := config.(dtchcl.Unmarshaler).UnmarshalHCL(decoder); err != nil {
				return nil, err
			}
			result = append(result, config)
		default:
		}
	}
	return result, nil
}
