/**
* @license
* Copyright 2020 Dynatrace LLC
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

package ipmappings

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	City        *string  `json:"city,omitempty"`       // The city name of the location.
	CountryCode string   `json:"countryCode"`          // The country code of the location. \n\n Use the alpha-2 code of the [ISO 3166-2 standard](https://dt-url.net/iso3166-2), (for example, `AT` for Austria or `PL` for Poland).
	Ip          string   `json:"ip"`                   // Single IP or IP range start address
	IpTo        *string  `json:"ipTo,omitempty"`       // IP range end
	Latitude    *float64 `json:"latitude,omitempty"`   // Latitude
	Longitude   *float64 `json:"longitude,omitempty"`  // Longitude
	RegionCode  *string  `json:"regionCode,omitempty"` // The region code of the location. \n\n For the [USA](https://dt-url.net/iso3166us) or [Canada](https://dt-url.net/iso3166ca) use ISO 3166-2 state codes without `US-` or `CA-` prefix. \n\n For the rest of the world use [FIPS 10-4 codes](https://dt-url.net/fipscodes) without country prefix.
}

func (me *Settings) Name() string {
	return uuid.NewString()
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"city": {
			Type:        schema.TypeString,
			Description: "The city name of the location.",
			Optional:    true,
		},
		"country_code": {
			Type:        schema.TypeString,
			Description: "The country code of the location. \n\n Use the alpha-2 code of the [ISO 3166-2 standard](https://dt-url.net/iso3166-2), (for example, `AT` for Austria or `PL` for Poland).",
			Required:    true,
		},
		"ip": {
			Type:        schema.TypeString,
			Description: "Single IP or IP range start address",
			Required:    true,
		},
		"ip_to": {
			Type:        schema.TypeString,
			Description: "IP range end",
			Optional:    true,
		},
		"latitude": {
			Type:        schema.TypeFloat,
			Description: "Latitude",
			Optional:    true,
		},
		"longitude": {
			Type:        schema.TypeFloat,
			Description: "Longitude",
			Optional:    true,
		},
		"region_code": {
			Type:        schema.TypeString,
			Description: "The region code of the location. \n\n For the [USA](https://dt-url.net/iso3166us) or [Canada](https://dt-url.net/iso3166ca) use ISO 3166-2 state codes without `US-` or `CA-` prefix. \n\n For the rest of the world use [FIPS 10-4 codes](https://dt-url.net/fipscodes) without country prefix.",
			Optional:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"city":         me.City,
		"country_code": me.CountryCode,
		"ip":           me.Ip,
		"ip_to":        me.IpTo,
		"latitude":     me.Latitude,
		"longitude":    me.Longitude,
		"region_code":  me.RegionCode,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	err := decoder.DecodeAll(map[string]any{
		"city":         &me.City,
		"country_code": &me.CountryCode,
		"ip":           &me.Ip,
		"ip_to":        &me.IpTo,
		"latitude":     &me.Latitude,
		"longitude":    &me.Longitude,
		"region_code":  &me.RegionCode,
	})
	if me.City != nil {
		if me.Latitude == nil {
			me.Latitude = opt.NewFloat64(0)
		}
		if me.Longitude == nil {
			me.Longitude = opt.NewFloat64(0)
		}
	}
	return err
}
