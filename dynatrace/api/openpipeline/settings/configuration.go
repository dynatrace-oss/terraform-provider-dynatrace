package openpipeline

import (
	"encoding/json"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Configuration struct {
	Kind           string        `json:"id"`
	Editable       *bool         `json:"editable,omitempty"`
	Version        string        `json:"version"`
	CustomBasePath string        `json:"customBasePath"`
	Endpoints      *Endpoints    `json:"-"`
	Pipelines      *Pipelines    `json:"-"`
	Routing        *RoutingTable `json:"routing"`
}

func (d *Configuration) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"endpoints": {
			Type:        schema.TypeList,
			Description: "List of all ingest sources of the configuration",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Endpoints).Schema()},
			Optional:    true,
		},
		"routing": {
			Type:        schema.TypeList,
			Description: "Dynamic routing definition",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(RoutingTable).Schema()},
			Optional:    true,
		},
		"pipelines": {
			Type:        schema.TypeList,
			Description: "List of all pipelines of the configuration",
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Pipelines).Schema()},
			Optional:    true,
		},
	}
}

func (d *Configuration) MarshalHCL(properties hcl.Properties) error {
	if d.Endpoints != nil && len(d.Endpoints.Endpoints) > 0 {
		if err := properties.Encode("endpoints", d.Endpoints); err != nil {
			return err
		}
	}

	if d.Pipelines != nil && len(d.Pipelines.Pipelines) > 0 {
		if err := properties.Encode("pipelines", d.Pipelines); err != nil {
			return err
		}
	}

	if d.Routing != nil && len(d.Routing.Entries) > 0 {
		if err := properties.Encode("routing", d.Routing); err != nil {
			return err
		}
	}

	return nil
}

func (d *Configuration) UnmarshalHCL(decoder hcl.Decoder) error {
	err := decoder.DecodeAll(map[string]any{
		"endpoints": &d.Endpoints,
		"pipelines": &d.Pipelines,
		"routing":   &d.Routing,
	})

	if d.Endpoints == nil {
		d.Endpoints = &Endpoints{}
	}

	if d.Pipelines == nil {
		d.Pipelines = &Pipelines{}
	}

	if d.Routing == nil {
		d.Routing = &RoutingTable{}
	}

	return err
}

func (d Configuration) MarshalJSON() ([]byte, error) {
	rawEndpoints, err := json.Marshal(d.Endpoints)
	if err != nil {
		return nil, err
	}

	rawPipelines, err := json.Marshal(d.Pipelines)
	if err != nil {
		return nil, err
	}

	type configuration Configuration
	conf := struct {
		configuration
		RawEndpoints json.RawMessage `json:"endpoints"`
		RawPipelines json.RawMessage `json:"pipelines"`
	}{
		configuration: (configuration)(d),
		RawEndpoints:  rawEndpoints,
		RawPipelines:  rawPipelines,
	}

	return json.Marshal(conf)
}

func (d *Configuration) UnmarshalJSON(b []byte) error {
	type configuration Configuration

	conf := struct {
		configuration
		RawEndpoints json.RawMessage `json:"endpoints"`
		RawPipelines json.RawMessage `json:"pipelines"`
	}{}

	if err := json.Unmarshal(b, &conf); err != nil {
		return err
	}

	*d = Configuration(conf.configuration)

	if err := json.Unmarshal(conf.RawEndpoints, &d.Endpoints); err != nil {
		return fmt.Errorf("error while reading endpoints field: %w", err)
	}

	if err := json.Unmarshal(conf.RawPipelines, &d.Pipelines); err != nil {
		return fmt.Errorf("error while reading pipelines field: %w", err)
	}
	return nil
}

func (d *Configuration) Name() string {
	return d.Kind
}

func (d *Configuration) RemoveFixed() {
	if d.Endpoints != nil {
		d.Endpoints.RemoveFixed()
	}

	if d.Routing != nil {
		d.Routing.RemoveFixed()
	}

	if d.Pipelines != nil {
		d.Pipelines.RemoveFixed()
	}
}
