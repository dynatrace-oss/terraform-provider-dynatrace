package challenge

import (
	"math/rand"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SupportedService struct {
	Name    string    `json:"name"`
	Metrics []*Metric `json:"metrics"`
}

func (me *SupportedService) Shuffle() {
	if len(me.Metrics) > 0 {
		rand.Shuffle(len(me.Metrics), func(i, j int) {
			ei := me.Metrics[i]
			ej := me.Metrics[j]
			me.Metrics[i] = ej
			me.Metrics[j] = ei
		})
		for _, elem := range me.Metrics {
			elem.Shuffle()
		}
	}
}

func (me *SupportedService) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "No documentation available",
			Required:    true,
		},
		"metric": {
			Type:        schema.TypeSet,
			Description: "No documentation available",
			Optional:    true,
			Set: func(i any) int {
				m := i.(map[string]any)
				if h, ok := m["__hash__"]; ok {
					return h.(int)
				}
				h := schema.HashString(m["name"])
				m["__hash__"] = h
				return h
			},
			Elem: &schema.Resource{Schema: new(Metric).Schema()},
		},
	}
}

func (me *SupportedService) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("name", me.Name); err != nil {
		return err
	}
	if err := properties.EncodeSlice("metric", me.Metrics); err != nil {
		return err
	}
	return nil
}

func (me *SupportedService) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("name", &me.Name); err != nil {
		return err
	}
	if err := decoder.DecodeSlice("metric", &me.Metrics); err != nil {
		return err
	}
	return nil
}
