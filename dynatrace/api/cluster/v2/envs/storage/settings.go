package storage

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/cluster/v2/envs/storage/retention"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Settings represents environment level storage usage and limit information. Not returned if includeStorageInfo param is not true. If skipped when editing via PUT method then already set limits will remain
type Settings struct {
	TransactionTrafficQuota *TransactionTrafficQuota `json:"transactionTrafficQuota"` // Maximum number of newly monitored entry point PurePaths captured per process/minute on environment level. Can be set to any value from 100 to 100000. If skipped when editing via PUT method then already set limit will remain
	UserActionsPerMinute    *UserActionsPerMinute    `json:"userActionsPerMinute"`    // Maximum number of user actions generated per minute on environment level. Can be set to any value from 1 to 2147483646 or left unlimited. If skipped when editing via PUT method then already set limit will remain

	Transactions              *Transactions              `json:"transactionStorage"`        // Transaction storage usage and limit information on environment level. If skipped when editing via PUT method then already set limit will remain
	SessionReplayStorage      *SessionReplayStorage      `json:"sessionReplayStorage"`      // Session replay storage usage and limit information on environment level. If skipped when editing via PUT method then already set limit will remain
	SymbolFilesFromMobileApps *SymbolFilesFromMobileApps `json:"symbolFilesFromMobileApps"` // Symbol files from mobile apps storage usage and limit information on environment level. If skipped when editing via PUT method then already set limit will remain
	LogMonitoringStorage      *LogMonitoringStorage      `json:"logMonitoringStorage"`      // Log monitoring storage usage and limit information on environment level. Not editable when Log monitoring is not allowed by license or not configured on cluster level. If skipped when editing via PUT method then already set limit will remain

	ServiceRequestLevelRetention *retention.ServiceRequestLevel `json:"serviceRequestLevelRetention"` // Service request level retention settings on environment level. Service code level retention time can't be greater than service request level retention time and both can't exceed one year.If skipped when editing via PUT method then already set limit will remain
	ServiceCodeLevelRetention    *retention.ServiceCodeLevel    `json:"serviceCodeLevelRetention"`    // Service code level retention settings on environment level. Service code level retention time can't be greater than service request level retention time and both can't exceed one year.If skipped when editing via PUT method then already set limit will remain
	RealUserMonitoringRetention  *retention.RealUserMonitoring  `json:"realUserMonitoringRetention"`  // Real user monitoring retention settings on environment level. Can be set to any value from 1 to 35 days. If skipped when editing via PUT method then already set limit will remain
	SyntheticMonitoringRetention *retention.SyntheticMonitoring `json:"syntheticMonitoringRetention"` // Synthetic monitoring retention settings on environment level. Can be set to any value from 1 to 35 days. If skipped when editing via PUT method then already set limit will remain
	SessionReplayRetention       *retention.SessionReplay       `json:"sessionReplayRetention"`       // Session replay retention settings on environment level. Can be set to any value from 1 to 35 days. If skipped when editing via PUT method then already set limit will remain
	LogMonitoringRetention       *retention.LogMonitoring       `json:"logMonitoringRetention"`       // Log monitoring retention settings on environment level. Not editable when Log monitoring is not allowed by license or not configured on cluster level. Can be set to any value from 5 to 90 days. If skipped when editing via PUT method then already set limit will remain
}

type limits struct {
	Transactions  *int64
	SessionReplay *int64
	SymbolFiles   *int64
	Logs          *int64
}

func (me *limits) IsEmpty() bool {
	return me == nil || (me.Transactions == nil && me.SessionReplay == nil && me.SymbolFiles == nil && me.Logs == nil)
}

func (me *limits) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"transactions": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Transaction storage usage and limit information on environment level in bytes. 0 for unlimited.",
		},
		"session_replay": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Session replay storage usage and limit information on environment level in bytes. 0 for unlimited.",
		},
		"symbol_files": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Session replay storage usage and limit information on environment level in bytes. 0 for unlimited.",
		},
		"logs": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Log monitoring storage usage and limit information on environment level in bytes. Not editable when Log monitoring is not allowed by license or not configured on cluster level. 0 for unlimited.",
		}}
}

func (me *limits) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("transactions", me.Transactions); err != nil {
		return err
	}
	if err := properties.Encode("session_replay", me.SessionReplay); err != nil {
		return err
	}
	if err := properties.Encode("symbol_files", me.SymbolFiles); err != nil {
		return err
	}
	if err := properties.Encode("logs", me.Logs); err != nil {
		return err
	}
	return nil
}

func (me *limits) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"transactions":   &me.Transactions,
		"session_replay": &me.SessionReplay,
		"symbol_files":   &me.SymbolFiles,
		"logs":           &me.Logs,
	})
}

type retent struct {
	ServiceRequestLevel int64
	ServiceCodeLevel    int64
	RUM                 int64
	Synthetic           int64
	SessionReplay       int64
	Logs                int64
}

func (me *retent) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"service_request_level": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Service request level retention settings on environment level in days. Service code level retention time can't be greater than service request level retention time and both can't exceed one year",
		},
		"service_code_level": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Service code level retention settings on environment level in days. Service code level retention time can't be greater than service request level retention time and both can't exceed one year",
		},
		"rum": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Real user monitoring retention settings on environment level in days. Can be set to any value from 1 to 35 days",
		},
		"synthetic": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Synthetic monitoring retention settings on environment level in days. Can be set to any value from 1 to 35 days",
		},
		"session_replay": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Session replay retention settings on environment level in days. Can be set to any value from 1 to 35 days",
		},
		"logs": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Log monitoring retention settings on environment level in days. Not editable when Log monitoring is not allowed by license or not configured on cluster level. Can be set to any value from 5 to 90 days",
		}}
}

func (me *retent) MarshalHCL(properties hcl.Properties) error {
	if me.ServiceCodeLevel != 0 {
		if err := properties.Encode("service_code_level", me.ServiceCodeLevel); err != nil {
			return err
		}
	}
	if me.ServiceCodeLevel != 0 {
		if err := properties.Encode("service_request_level", me.ServiceRequestLevel); err != nil {
			return err
		}
	}
	if me.RUM != 0 {
		if err := properties.Encode("rum", me.RUM); err != nil {
			return err
		}
	}
	if me.Synthetic != 0 {
		if err := properties.Encode("synthetic", me.Synthetic); err != nil {
			return err
		}
	}
	if me.SessionReplay != 0 {
		if err := properties.Encode("session_replay", me.SessionReplay); err != nil {
			return err
		}
	}
	if me.Logs != 0 {
		if err := properties.Encode("logs", me.Logs); err != nil {
			return err
		}
	}
	return nil
}

func (me *retent) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"service_request_level": &me.ServiceRequestLevel,
		"service_code_level":    &me.ServiceCodeLevel,
		"rum":                   &me.RUM,
		"synthetic":             &me.Synthetic,
		"session_replay":        &me.SessionReplay,
		"logs":                  &me.Logs,
	})
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"user_actions": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Maximum number of user actions generated per minute on environment level. Can be set to any value from 1 to 2147483646 or left unlimited by omitting this property",
		},
		"transactions": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Maximum number of newly monitored entry point PurePaths captured per process/minute on environment level. Can be set to any value from 100 to 100000",
		},
		"limits": {
			Type:     schema.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 1,
			Elem:     &schema.Resource{Schema: new(limits).Schema()},
		},
		"retention": {
			Type:     schema.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 1,
			Elem:     &schema.Resource{Schema: new(retent).Schema()},
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	if !me.UserActionsPerMinute.IsEmpty() {
		if err := properties.Encode("user_actions", me.UserActionsPerMinute.MaxLimit); err != nil {
			return err
		}
	}
	if !me.TransactionTrafficQuota.IsEmpty() {
		if err := properties.Encode("transactions", me.TransactionTrafficQuota.MaxLimit); err != nil {
			return err
		}
	}
	vLimits := new(limits)

	if !me.Transactions.IsEmpty() {
		vLimits.Transactions = me.Transactions.MaxLimit
	}
	if !me.SessionReplayStorage.IsEmpty() {
		vLimits.SessionReplay = me.SessionReplayStorage.MaxLimit
	}
	if !me.SymbolFilesFromMobileApps.IsEmpty() {
		vLimits.SymbolFiles = me.SymbolFilesFromMobileApps.MaxLimit
	}
	if !me.LogMonitoringStorage.IsEmpty() {
		vLimits.Logs = me.LogMonitoringStorage.MaxLimit
	}
	if !vLimits.IsEmpty() {
		if err := properties.Encode("limits", vLimits); err != nil {
			return err
		}
	}

	if err := properties.Encode("retention", &retent{
		Logs:                me.LogMonitoringRetention.MaxLimitInDays,
		SessionReplay:       me.SessionReplayRetention.MaxLimitInDays,
		ServiceCodeLevel:    me.ServiceCodeLevelRetention.MaxLimitInDays,
		ServiceRequestLevel: me.ServiceRequestLevelRetention.MaxLimitInDays,
		RUM:                 me.RealUserMonitoringRetention.MaxLimitInDays,
		Synthetic:           me.SyntheticMonitoringRetention.MaxLimitInDays,
	}); err != nil {
		return err
	}

	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	me.UserActionsPerMinute = new(UserActionsPerMinute)
	if err := decoder.Decode("user_actions", &me.UserActionsPerMinute.MaxLimit); err != nil {
		return err
	}
	me.TransactionTrafficQuota = new(TransactionTrafficQuota)
	if err := decoder.Decode("transactions", &me.TransactionTrafficQuota.MaxLimit); err != nil {
		return err
	}

	vLimits := new(limits)
	if err := decoder.Decode("limits", &vLimits); err != nil {
		return err
	}

	me.Transactions = &Transactions{MaxLimit: vLimits.Transactions}
	me.SessionReplayStorage = &SessionReplayStorage{MaxLimit: vLimits.SessionReplay}
	me.SymbolFilesFromMobileApps = &SymbolFilesFromMobileApps{MaxLimit: vLimits.SymbolFiles}
	me.LogMonitoringStorage = &LogMonitoringStorage{MaxLimit: vLimits.Logs}

	ret := new(retent)
	if err := decoder.Decode("retention", &ret); err != nil {
		return err
	}
	me.LogMonitoringRetention = &retention.LogMonitoring{MaxLimitInDays: ret.Logs}
	me.SessionReplayRetention = &retention.SessionReplay{MaxLimitInDays: ret.SessionReplay}
	me.SyntheticMonitoringRetention = &retention.SyntheticMonitoring{MaxLimitInDays: ret.Synthetic}
	me.RealUserMonitoringRetention = &retention.RealUserMonitoring{MaxLimitInDays: ret.RUM}
	me.ServiceCodeLevelRetention = &retention.ServiceCodeLevel{MaxLimitInDays: ret.ServiceCodeLevel}
	me.ServiceRequestLevelRetention = &retention.ServiceRequestLevel{MaxLimitInDays: ret.ServiceRequestLevel}
	return nil
}
