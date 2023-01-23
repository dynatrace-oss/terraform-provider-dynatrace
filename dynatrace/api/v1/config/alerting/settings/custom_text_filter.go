package alerting

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CustomTextFilter Configuration of a matching filter.
type CustomTextFilter struct {
	Enabled         bool                       `json:"enabled"`         // The filter is enabled (`true`) or disabled (`false`).
	Negate          bool                       `json:"negate"`          // Reverses the comparison **operator**. For example it turns the **begins with** into **does not begin with**.
	Operator        Operator                   `json:"operator"`        // Operator of the comparison.   You can reverse it by setting **negate** to `true`.
	Value           string                     `json:"value"`           // The value to compare to.
	CaseInsensitive bool                       `json:"caseInsensitive"` // The condition is case sensitive (`false`) or case insensitive (`true`).   If not set, then `false` is used, making the condition case sensitive.
	Unknowns        map[string]json.RawMessage `json:"-"`
}

func (me *CustomTextFilter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "The filter is enabled (`true`) or disabled (`false`)",
			Optional:    true,
		},
		"negate": {
			Type:        schema.TypeBool,
			Description: "Reverses the comparison **operator**. For example it turns the **begins with** into **does not begin with**",
			Optional:    true,
		},
		"operator": {
			Type:        schema.TypeString,
			Description: "Operator of the comparison.   You can reverse it by setting **negate** to `true`. Possible values are `BEGINS_WITH`, `CONTAINS`, `CONTAINS_REGEX`, `ENDS_WITH` and `EQUALS`",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "The value to compare to",
			Required:    true,
		},
		"case_insensitive": {
			Type:        schema.TypeBool,
			Description: "The condition is case sensitive (`false`) or case insensitive (`true`).   If not set, then `false` is used, making the condition case sensitive",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *CustomTextFilter) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("enabled", me.Enabled); err != nil {
		return err
	}
	if err := properties.Encode("negate", me.Negate); err != nil {
		return err
	}
	if err := properties.Encode("operator", string(me.Operator)); err != nil {
		return err
	}
	if err := properties.Encode("value", me.Value); err != nil {
		return err
	}
	if err := properties.Encode("case_insensitive", me.CaseInsensitive); err != nil {
		return err
	}

	return nil
}

func (me *CustomTextFilter) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "enabled")
		delete(me.Unknowns, "negate")
		delete(me.Unknowns, "operator")
		delete(me.Unknowns, "value")
		delete(me.Unknowns, "case_insensitive")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("enabled"); ok {
		me.Enabled = value.(bool)
	}
	if value, ok := decoder.GetOk("negate"); ok {
		me.Negate = value.(bool)
	}
	if value, ok := decoder.GetOk("case_insensitive"); ok {
		me.CaseInsensitive = value.(bool)
	}
	if value, ok := decoder.GetOk("operator"); ok {
		me.Operator = Operator(value.(string))
	}
	if value, ok := decoder.GetOk("value"); ok {
		me.Value = value.(string)
	}
	return nil
}

func (me *CustomTextFilter) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(me.Unknowns) > 0 {
		for k, v := range me.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(me.Enabled)
		if err != nil {
			return nil, err
		}
		m["enabled"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(me.Negate)
		if err != nil {
			return nil, err
		}
		m["negate"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(me.Operator)
		if err != nil {
			return nil, err
		}
		m["operator"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(me.Value)
		if err != nil {
			return nil, err
		}
		m["value"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(me.CaseInsensitive)
		if err != nil {
			return nil, err
		}
		m["caseInsensitive"] = rawMessage
	}
	return json.Marshal(m)
}

func (me *CustomTextFilter) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["enabled"]; found {
		if err := json.Unmarshal(v, &me.Enabled); err != nil {
			return err
		}
	}
	if v, found := m["negate"]; found {
		if err := json.Unmarshal(v, &me.Negate); err != nil {
			return err
		}
	}
	if v, found := m["operator"]; found {
		if err := json.Unmarshal(v, &me.Operator); err != nil {
			return err
		}
	}
	if v, found := m["value"]; found {
		if err := json.Unmarshal(v, &me.Value); err != nil {
			return err
		}
	}
	if v, found := m["caseInsensitive"]; found {
		if err := json.Unmarshal(v, &me.CaseInsensitive); err != nil {
			return err
		}
	}

	delete(m, "enabled")
	delete(m, "negate")
	delete(m, "operator")
	delete(m, "value")
	delete(m, "caseInsensitive")

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}

// Operator Operator of the comparison.
//
//	You can reverse it by setting **negate** to `true`.
type Operator string

// Operators offers the known enum values
var Operators = struct {
	BeginsWith    Operator
	Contains      Operator
	ContainsRegex Operator
	EndsWith      Operator
	Equals        Operator
}{
	"BEGINS_WITH",
	"CONTAINS",
	"CONTAINS_REGEX",
	"ENDS_WITH",
	"EQUALS",
}
