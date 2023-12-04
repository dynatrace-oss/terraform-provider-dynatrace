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

package logagentconfiguration

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	ContainerTimezoneHeuristicEnabled bool    `json:"LAConfigContainerTimezoneHeuristicEnabled"` // Enables automatic detection of timezone in container's logs if it is not explicitly defined in content or configured.
	ContainersLogsDetectionEnabled    bool    `json:"LAConfigContainersLogsDetectionEnabled"`    // Allows detection of log messages written to the containerized application's stdout/stderr streams.
	DateSearchLimit_Bytes             int     `json:"LAConfigDateSearchLimit_Bytes"`             // Defines the number of characters in every log line (starting from the first character in the line) where the timestamp is searched.
	DefaultTimezone                   string  `json:"LAConfigDefaultTimezone"`                   // Default timezone for agent if more specific configurations is not defined.
	EventLogQueryTimeout_Sec          int     `json:"LAConfigEventLogQueryTimeout_Sec"`          // Defines the maximum timeout value, in seconds, for the query extracting Windows Event Logs
	IISDetectionEnabled               bool    `json:"LAConfigIISDetectionEnabled"`               // Allows detection of logs and event logs written by IIS server.
	LogScannerLinuxNfsEnabled         bool    `json:"LAConfigLogScannerLinuxNfsEnabled"`         // Allows detection of logs written to mounted network storage drives.
	MaxLgisPerEntityCount             int     `json:"LAConfigMaxLgisPerEntityCount"`             // Defines the maximum number of log group instances per entity after which, the new automatic ones wouldn't be added.
	MinBinaryDetectionLimit_Bytes     int     `json:"LAConfigMinBinaryDetectionLimit_Bytes"`     // Defines the minimum number of bytes in log file required for binary detection.
	MonitorOwnLogsEnabled             bool    `json:"LAConfigMonitorOwnLogsEnabled"`             // Enabling this option may affect your licensing costs. For more details, see [documentation](https://dt-url.net/4l02yi8).
	OpenLogFilesDetectionEnabled      bool    `json:"LAConfigOpenLogFilesDetectionEnabled"`      // Automatically detect logs written by important processes. For more details, check our [documentation](https://dt-url.net/7v02z76)
	SeverityDetectionLimit_Bytes      int     `json:"LAConfigSeverityDetectionLimit_Bytes"`      // Defines the number of characters in every log line (starting from the first character in the line) where severity is searched.
	SeverityDetectionLinesLimit       int     `json:"LAConfigSeverityDetectionLinesLimit"`       // Defines the number of the first lines of every log entry where severity is searched.
	SystemLogsDetectionEnabled        bool    `json:"LAConfigSystemLogsDetectionEnabled"`        // Linux: syslog, message log Windows: system, application, security event logs
	Scope                             *string `json:"-" scope:"scope"`                           // The scope of this setting (HOST, HOST_GROUP). Omit this property if you want to cover the whole environment.
}

func (me *Settings) Name() string {
	return *me.Scope
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"container_timezone_heuristic_enabled": {
			Type:        schema.TypeBool,
			Description: "Enables automatic detection of timezone in container's logs if it is not explicitly defined in content or configured.",
			Required:    true,
		},
		"containers_logs_detection_enabled": {
			Type:        schema.TypeBool,
			Description: "Allows detection of log messages written to the containerized application's stdout/stderr streams.",
			Required:    true,
		},
		"date_search_limit_bytes": {
			Type:        schema.TypeInt,
			Description: "Defines the number of characters in every log line (starting from the first character in the line) where the timestamp is searched.",
			Required:    true,
		},
		"default_timezone": {
			Type:        schema.TypeString,
			Description: "Default timezone for agent if more specific configurations is not defined.",
			Required:    true,
		},
		"event_log_query_timeout_sec": {
			Type:        schema.TypeInt,
			Description: "Defines the maximum timeout value, in seconds, for the query extracting Windows Event Logs",
			Required:    true,
		},
		"iisdetection_enabled": {
			Type:        schema.TypeBool,
			Description: "Allows detection of logs and event logs written by IIS server.",
			Required:    true,
		},
		"log_scanner_linux_nfs_enabled": {
			Type:        schema.TypeBool,
			Description: "Allows detection of logs written to mounted network storage drives.",
			Required:    true,
		},
		"max_lgis_per_entity_count": {
			Type:        schema.TypeInt,
			Description: "Defines the maximum number of log group instances per entity after which, the new automatic ones wouldn't be added.",
			Required:    true,
		},
		"min_binary_detection_limit_bytes": {
			Type:        schema.TypeInt,
			Description: "Defines the minimum number of bytes in log file required for binary detection.",
			Required:    true,
		},
		"monitor_own_logs_enabled": {
			Type:        schema.TypeBool,
			Description: "Enabling this option may affect your licensing costs. For more details, see [documentation](https://dt-url.net/4l02yi8).",
			Required:    true,
		},
		"open_log_files_detection_enabled": {
			Type:        schema.TypeBool,
			Description: "Automatically detect logs written by important processes. For more details, check our [documentation](https://dt-url.net/7v02z76)",
			Required:    true,
		},
		"severity_detection_limit_bytes": {
			Type:        schema.TypeInt,
			Description: "Defines the number of characters in every log line (starting from the first character in the line) where severity is searched.",
			Required:    true,
		},
		"severity_detection_lines_limit": {
			Type:        schema.TypeInt,
			Description: "Defines the number of the first lines of every log entry where severity is searched.",
			Required:    true,
		},
		"system_logs_detection_enabled": {
			Type:        schema.TypeBool,
			Description: "Linux: syslog, message log Windows: system, application, security event logs",
			Required:    true,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (HOST, HOST_GROUP). Omit this property if you want to cover the whole environment.",
			Optional:    true,
			Default:     "environment",
			ForceNew:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"container_timezone_heuristic_enabled": me.ContainerTimezoneHeuristicEnabled,
		"containers_logs_detection_enabled":    me.ContainersLogsDetectionEnabled,
		"date_search_limit_bytes":              me.DateSearchLimit_Bytes,
		"default_timezone":                     me.DefaultTimezone,
		"event_log_query_timeout_sec":          me.EventLogQueryTimeout_Sec,
		"iisdetection_enabled":                 me.IISDetectionEnabled,
		"log_scanner_linux_nfs_enabled":        me.LogScannerLinuxNfsEnabled,
		"max_lgis_per_entity_count":            me.MaxLgisPerEntityCount,
		"min_binary_detection_limit_bytes":     me.MinBinaryDetectionLimit_Bytes,
		"monitor_own_logs_enabled":             me.MonitorOwnLogsEnabled,
		"open_log_files_detection_enabled":     me.OpenLogFilesDetectionEnabled,
		"severity_detection_limit_bytes":       me.SeverityDetectionLimit_Bytes,
		"severity_detection_lines_limit":       me.SeverityDetectionLinesLimit,
		"system_logs_detection_enabled":        me.SystemLogsDetectionEnabled,
		"scope":                                me.Scope,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"container_timezone_heuristic_enabled": &me.ContainerTimezoneHeuristicEnabled,
		"containers_logs_detection_enabled":    &me.ContainersLogsDetectionEnabled,
		"date_search_limit_bytes":              &me.DateSearchLimit_Bytes,
		"default_timezone":                     &me.DefaultTimezone,
		"event_log_query_timeout_sec":          &me.EventLogQueryTimeout_Sec,
		"iisdetection_enabled":                 &me.IISDetectionEnabled,
		"log_scanner_linux_nfs_enabled":        &me.LogScannerLinuxNfsEnabled,
		"max_lgis_per_entity_count":            &me.MaxLgisPerEntityCount,
		"min_binary_detection_limit_bytes":     &me.MinBinaryDetectionLimit_Bytes,
		"monitor_own_logs_enabled":             &me.MonitorOwnLogsEnabled,
		"open_log_files_detection_enabled":     &me.OpenLogFilesDetectionEnabled,
		"severity_detection_limit_bytes":       &me.SeverityDetectionLimit_Bytes,
		"severity_detection_lines_limit":       &me.SeverityDetectionLinesLimit,
		"system_logs_detection_enabled":        &me.SystemLogsDetectionEnabled,
		"scope":                                &me.Scope,
	})
}
