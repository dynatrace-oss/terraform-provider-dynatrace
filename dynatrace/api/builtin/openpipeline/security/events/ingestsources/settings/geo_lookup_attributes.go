/**
* @license
* Copyright 2025 Dynatrace LLC
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

package ingestsources

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type GeoLookupAttributes struct {
	GeoFieldPrefix *string          `json:"geoFieldPrefix,omitempty"` // Optional prefix for all output geo fields. If specified, output fields will be prefixed as <prefix>.geo.<field>. If omitted, output fields will be geo.<field>.
	IpFieldKey     string           `json:"ipFieldKey"`               // The field key that contains the IP address to be resolved to a geo location.
	OutputFields   []GeoOutputField `json:"outputFields,omitempty"`   // The geo fields to enrich the record with. If empty or not specified, the default fields (city name, country ISO code, country name, location) are used. Possible values: `cityName`, `continentIsoCode`, `continentName`, `countryIsoCode`, `countryName`, `location`, `postalCode`, `regionIsoCode`, `regionName`, `subdivisionIsoCodes`
}

func (me *GeoLookupAttributes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"geo_field_prefix": {
			Type:        schema.TypeString,
			Description: "Optional prefix for all output geo fields. If specified, output fields will be prefixed as <prefix>.geo.<field>. If omitted, output fields will be geo.<field>.",
			Optional:    true, // nullable
		},
		"ip_field_key": {
			Type:        schema.TypeString,
			Description: "The field key that contains the IP address to be resolved to a geo location.",
			Required:    true,
		},
		"output_fields": {
			Type:        schema.TypeSet,
			Description: "The geo fields to enrich the record with. If empty or not specified, the default fields (city name, country ISO code, country name, location) are used. Possible values: `cityName`, `continentIsoCode`, `continentName`, `countryIsoCode`, `countryName`, `location`, `postalCode`, `regionIsoCode`, `regionName`, `subdivisionIsoCodes`",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *GeoLookupAttributes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"geo_field_prefix": me.GeoFieldPrefix,
		"ip_field_key":     me.IpFieldKey,
		"output_fields":    me.OutputFields,
	})
}

func (me *GeoLookupAttributes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"geo_field_prefix": &me.GeoFieldPrefix,
		"ip_field_key":     &me.IpFieldKey,
		"output_fields":    &me.OutputFields,
	})
}
