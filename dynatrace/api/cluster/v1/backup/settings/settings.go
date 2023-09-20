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

package backup

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// The configuration of the user
type Settings struct {
	Enabled                *bool   `json:"enabled,omitempty"`               // Backups are enabled (true) or disabled (false).
	Datacenter             *string `json:"datacenter,omitempty"`            // Datacenter which will create backups
	StoragePath            *string `json:"storagePath,omitempty"`           // A full path to the backup archive
	IncludeRumData         *bool   `json:"includeRumData,omitempty"`        // Include user sessions (true) or GDPR compliance (false)
	IncludeLm20Data        *bool   `json:"includeLm20Data,omitempty"`       // Include (true) or exclude (false) Log Monitoring v2 data
	IncludeTsMetricData    *bool   `json:"includeTsMetricData,omitempty"`   // Include time series metric-data (true) or retain configuration data only (false))
	BandwidthLimitMbits    *int    `json:"bandwidthLimitMbits,omitempty"`   // Cassandra backup bandwidth limit in Mbps
	MaxEsSnapshotsToClean  *int    `json:"maxEsSnapshotsToClean,omitempty"` // Max number of Elasticsearch snapshots to clean. Elasticsearch snapshots won't be created anymore if there will be more backups to clean than this value.
	CassandraScheduledTime int     `json:"cassandraScheduledTime"`          // Hour to start Cassandra backups each day.
	PauseBackups           *bool   `json:"pauseBackups,omitempty"`          // Pauses Elasticsearch and Cassandra backups. In comparison to enable/disable backup, this option does not modify any configuration like Elasticsearch properties.
	CurrentState           string  `json:"-"`
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Backups are enabled (true) or disabled (false).",
			Optional:    true,
		},
		"datacenter": {
			Type:        schema.TypeString,
			Description: "Datacenter which will create backups",
			Optional:    true,
		},
		"storage_path": {
			Type:        schema.TypeString,
			Description: "A full path to the backup archive",
			Optional:    true,
		},
		"include_rum_data": {
			Type:        schema.TypeBool,
			Description: "Include user sessions (true) or GDPR compliance (false)",
			Optional:    true,
		},
		"include_lm20_data": {
			Type:        schema.TypeBool,
			Description: "Include (true) or exclude (false) Log Monitoring v2 data",
			Optional:    true,
		},
		"include_ts_metric_data": {
			Type:        schema.TypeBool,
			Description: "Include time series metric-data (true) or retain configuration data only (false))",
			Optional:    true,
		},
		"bandwidth_limit_mbits": {
			Type:        schema.TypeInt,
			Description: "Cassandra backup bandwidth limit in Mbps",
			Optional:    true,
		},
		"max_es_snapshots_to_clean": {
			Type:        schema.TypeInt,
			Description: "Max number of Elasticsearch snapshots to clean. Elasticsearch snapshots won't be created anymore if there will be more backups to clean than this value.",
			Optional:    true,
		},
		"cassandra_scheduled_time": {
			Type:        schema.TypeInt,
			Description: "Hour to start Cassandra backups each day.",
			Required:    true,
		},
		"pause_backups": {
			Type:        schema.TypeBool,
			Description: "Pauses Elasticsearch and Cassandra backups. In comparison to enable/disable backup, this option does not modify any configuration like Elasticsearch properties.",
			Optional:    true,
		},
		"current_state": {
			Type:        schema.TypeString,
			Description: "For internal use: current state of rules in JSON format",
			Optional:    true,
			Computed:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"include_rum_data":          me.IncludeRumData,
		"include_lm20_data":         me.IncludeLm20Data,
		"include_ts_metric_data":    me.IncludeTsMetricData,
		"bandwidth_limit_mbits":     me.BandwidthLimitMbits,
		"max_es_snapshots_to_clean": me.MaxEsSnapshotsToClean,
		"cassandra_scheduled_time":  me.CassandraScheduledTime,
		"pause_backups":             me.PauseBackups,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"include_rum_data":          &me.IncludeRumData,
		"include_lm20_data":         &me.IncludeLm20Data,
		"include_ts_metric_data":    &me.IncludeTsMetricData,
		"bandwidth_limit_mbits":     &me.BandwidthLimitMbits,
		"max_es_snapshots_to_clean": &me.MaxEsSnapshotsToClean,
		"cassandra_scheduled_time":  &me.CassandraScheduledTime,
		"pause_backups":             &me.PauseBackups,
	})
}
