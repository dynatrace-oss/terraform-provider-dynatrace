package quota

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Settings represents environment level consumption and quotas information. Only returned if includeConsumptionInfo or includeUncachedConsumptionInfo param is true. If skipped when editing via PUT method then already set quotas will remain
type Settings struct {
	HostUnits     *HostUnits      `json:"hostUnits"`         // Host units consumption and quota information on environment level. If skipped when editing via PUT method then already set quota will remain
	DEMUnits      *DEMUnits       `json:"demUnits"`          // DEM units consumption and quota information on environment level. Not set (and not editable) if DEM units is not enabled. If skipped when editing via PUT method then already set quotas will remain
	UserSessions  *UserSessions   `json:"userSessions"`      // User sessions consumption and quota information on environment level. If skipped when editing via PUT method then already set quotas will remain
	Synthetic     *Synthetic      `json:"syntheticMonitors"` // Synthetic monitors consumption and quota information on environment level. Not set (and not editable) if neither Synthetic nor DEM units is enabled. If skipped when editing via PUT method then already set quotas will remain
	DDUs          *DavisDataUnits `json:"davisDataUnits"`    // Davis Data Units consumption and quota information on environment level. Not set (and not editable) if Davis data units is not enabled. If skipped when editing via PUT method then already set quotas will remain
	LogMonitoring *LogMonitoring  `json:"logMonitoring"`     // Log Monitoring consumption and quota information on environment level. Not set (and not editable) if Log monitoring is not enabled. Not set (and not editable) if Log monitoring is migrated to Davis data on license level. If skipped when editing via PUT method then already set quotas will remain
}

func (me *Settings) IsEmpty() bool {
	if me == nil {
		return true
	}
	return me.HostUnits.IsEmpty() &&
		me.DEMUnits.IsEmpty() &&
		me.UserSessions.IsEmpty() &&
		me.Synthetic.IsEmpty() &&
		me.DDUs.IsEmpty() &&
		me.LogMonitoring.IsEmpty()
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"host_units": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Host units consumption and quota information on environment level",
		},
		"dem_units": {
			Type:        schema.TypeList,
			MinItems:    1,
			MaxItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(DEMUnits).Schema()},
			Description: "DEM units consumption and quota information on environment level",
		},
		"user_sessions": {
			Type:        schema.TypeList,
			MinItems:    1,
			MaxItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(UserSessions).Schema()},
			Description: "User sessions consumption and quota information on environment level",
		},
		"synthetic": {
			Type:        schema.TypeList,
			MinItems:    1,
			MaxItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(Synthetic).Schema()},
			Description: "Synthetic monitors consumption and quota information on environment level. Not set (and not editable) if neither Synthetic nor DEM units is enabled",
		},
		"ddus": {
			Type:        schema.TypeList,
			MinItems:    1,
			MaxItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(DavisDataUnits).Schema()},
			Description: "Davis Data Units consumption and quota information on environment level. Not set (and not editable) if Davis data units is not enabled",
		},
		"logs": {
			Type:        schema.TypeList,
			MinItems:    1,
			MaxItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(LogMonitoring).Schema()},
			Description: "Log Monitoring consumption and quota information on environment level. Not set (and not editable) if Log monitoring is not enabled. Not set (and not editable) if Log monitoring is migrated to Davis data on license level",
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	if !me.HostUnits.IsEmpty() {
		if err := properties.Encode("host_units", me.HostUnits.MaxLimit); err != nil {
			return err
		}
	}
	if !me.DEMUnits.IsEmpty() {
		if err := properties.Encode("dem_units", me.DEMUnits); err != nil {
			return err
		}
	}
	if !me.UserSessions.IsEmpty() {
		if err := properties.Encode("user_sessions", me.UserSessions); err != nil {
			return err
		}
	}
	if !me.Synthetic.IsEmpty() {
		if err := properties.Encode("synthetic", me.Synthetic); err != nil {
			return err
		}
	}
	if me.DDUs != nil && !me.DDUs.IsEmpty() {
		if err := properties.Encode("ddus", me.DDUs); err != nil {
			return err
		}
	}
	if me.LogMonitoring != nil && !me.LogMonitoring.IsEmpty() {
		if err := properties.Encode("logs", me.LogMonitoring); err != nil {
			return err
		}
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	me.HostUnits = new(HostUnits)
	if err := decoder.Decode("host_units", &me.HostUnits.MaxLimit); err != nil {
		return err
	}
	if err := decoder.Decode("dem_units", &me.DEMUnits); err != nil {
		return err
	}
	if me.DEMUnits == nil {
		me.DEMUnits = &DEMUnits{MonthlyLimit: nil, AnnualLimit: nil}
	}
	if err := decoder.Decode("user_sessions", &me.UserSessions); err != nil {
		return err
	}
	if me.UserSessions == nil {
		me.UserSessions = &UserSessions{TotalAnnualLimit: nil, TotalMonthlyLimit: nil}
	}
	if err := decoder.Decode("synthetic", &me.Synthetic); err != nil {
		return err
	}
	if me.Synthetic == nil {
		me.Synthetic = &Synthetic{MonthlyLimit: nil, AnnualLimit: nil}
	}
	if err := decoder.Decode("ddus", &me.DDUs); err != nil {
		return err
	}
	if me.DDUs == nil {
		me.DDUs = &DavisDataUnits{MonthlyLimit: nil, AnnualLimit: nil}
	}
	if err := decoder.Decode("logs", &me.LogMonitoring); err != nil {
		return err
	}
	if me.LogMonitoring == nil {
		me.LogMonitoring = &LogMonitoring{MonthlyLimit: nil, AnnualLimit: nil}
	}
	return nil
}
