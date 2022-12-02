package hcl2json

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	dtchcl "github.com/dtcookie/hcl"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

type Module struct {
	Variables map[string]string `json:"variables,omitempty"`
	Records   []*Record         `json:"records"`
}

type Record struct {
	ID         string      `json:"id"`
	Schema     string      `json:"schema,omitempty"`
	Resource   string      `json:"resource,omitempty"`
	Scope      string      `json:"scope"`
	ScopeClass string      `json:"-"`
	IDV1       string      `json:"idv1,omitempty"`
	Value      interface{} `json:"value"`
}

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
	return hash

}

func cty2str(m map[string]cty.Value) map[string]any {
	res := map[string]any{}
	for k, v := range m {
		if v.Type() == cty.String {
			res[k] = v.AsString()
		} else {
			res[k] = cty2str(v.AsValueMap())
		}
	}
	return res
}

func translate(body hcl.Body, sch *bodySchema, bc breadCrumbs, target map[string]interface{}, variables map[string]cty.Value) error {
	bodyContent, diag := body.Content(&sch.BodySchema)
	if diag.HasErrors() {
		return &diag
	}

	for _, attribute := range bodyContent.Attributes {

		val, diag := attribute.Expr.Value(&hcl.EvalContext{Variables: variables})
		if diag.HasErrors() {
			bo, _ := json.MarshalIndent(cty2str(variables), "", "  ")
			fmt.Println(string(bo))
			return &diag
		}
		value := sch.Prototypes[attribute.Name]
		switch typedValue := value.(type) {
		case stringSet:
			if val.Type().IsTupleType() {
				val = cty.ListVal(val.AsValueSlice())
			}
			if err := gocty.FromCtyValue(val, &typedValue); err != nil {
				return fmt.Errorf("%v[%s] - %s: %T, %v", bc, attribute.Name, err.Error(), value, val.GoString())
			}
			target[attribute.Name] = typedValue
		case []string:
			if val.Type().IsTupleType() {
				val = cty.ListVal(val.AsValueSlice())
			}
			if err := gocty.FromCtyValue(val, &typedValue); err != nil {
				return fmt.Errorf("%v[%s] - %s: %T, %v", bc, attribute.Name, err.Error(), value, val.GoString())
			}
			is := []interface{}{}
			for _, s := range typedValue {
				is = append(is, s)
			}
			target[attribute.Name] = is
		case string:
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
		if err := translate(block.Body, blockSchema, bc.dot(block.Type), blockMap, variables); err != nil {
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

func readHCLVariables(fileName string) (map[string]cty.Value, map[string]string, error) {
	properties := map[string]string{}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, nil, err
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "# DEFINE ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "# DEFINE "))
			keyValuePair := strings.Split(line, " = ")
			if len(keyValuePair) == 2 {
				key := strings.TrimSpace(keyValuePair[0])
				value := strings.TrimSpace(keyValuePair[1])
				properties[key] = value
			}
		}
	}

	file.Close()

	uvars := map[string]any{}

	for k, v := range properties {
		store(k, v, uvars)
	}

	ctyMapVal := any2cty(uvars)
	if ctyMapVal == cty.NilVal {
		return nil, nil, nil
	}
	return ctyMapVal.AsValueMap(), properties, nil
}

func store(key string, value string, m map[string]any) {
	parts := strings.Split(key, ".")
	if len(parts) == 1 {
		m[key] = value
		return
	}
	var valueMap map[string]any
	valueMap, exists := m[parts[0]].(map[string]any)
	if !exists {
		valueMap = map[string]any{}
	}

	store(strings.TrimPrefix(key, parts[0]+"."), value, valueMap)
	m[parts[0]] = valueMap

}

func any2cty(v any) cty.Value {
	switch tv := v.(type) {
	case string:
		return cty.StringVal(tv)
	case map[string]any:
		m := map[string]cty.Value{}
		for mk, mv := range tv {
			m[mk] = any2cty(mv)
		}
		if len(m) == 0 {
			m["pseudo"] = cty.StringVal("pseudo")
		}
		return cty.ObjectVal(m)
	default:
		panic(fmt.Sprintf("unexpected type %T", tv))
	}
}

func HCL2Config(fileName string) (*Module, error) {
	variables, properties, err := readHCLVariables(fileName)
	if err != nil {
		return nil, err
	}
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
	var result []*Record
	for _, block := range bodyContent.Blocks {
		switch block.Type {
		case "resource":
			m := map[string]interface{}{}
			resource := block.Labels[0]
			resourceName := block.Labels[1]

			bs := schemata[resource]
			if err := translate(block.Body, bs, breadCrumbs(block.Labels[0]), m, variables); err != nil {
				return nil, err
			}
			decoder := NewMapDecoder(m)
			config := reflect.New(reflect.TypeOf(protoTypes[resource]).Elem()).Interface()
			if err := config.(dtchcl.Unmarshaler).UnmarshalHCL(decoder); err != nil {
				return nil, err
			}
			record := &Record{
				ID:       variables[resource].AsValueMap()[resourceName].AsValueMap()["id"].AsString(),
				Value:    config,
				Schema:   resource,
				Resource: resource + "." + resourceName,
				Scope:    "environment",
			}
			objID := &ObjectID{ID: record.ID}
			if err := objID.Decode(); err == nil {
				if len(objID.SchemaID) > 0 {
					record.Schema = objID.SchemaID
				}
				record.ScopeClass = objID.Scope.Class
			} else {
				record.Schema = err.Error()
			}
			result = append(result, record)
		default:
		}
	}
	return &Module{
		Records:   result,
		Variables: properties,
	}, nil
}
