package processgroupalerting

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Enabled                  bool          `json:"enabled"`                            // Enable process group availability monitoring
	AlertingMode             *AlertingMode `json:"alertingMode,omitempty"`             // **if any process becomes unavailable:**\nDynatrace will open a new problem if a single process in this group shuts down or crashes. \n\n**if minimum threshold is not met:**\nDynatrace tracks the number of process instances that comprise this process group and treats the group as a cluster. This setting enables you to define a minimum number of process instances that must be available. A problem will be opened if this process group has fewer than the minimum number of required process instances. \n\n Details of the related impact on service requests will be included in the problem summary.\n\n**Note:** If a process is intentionally shutdown or retired while this setting is active, you'll need to manually close the problem.
	MinimumInstanceThreshold *int          `json:"minimumInstanceThreshold,omitempty"` // Open a new problem if the number of active process instances in the group is fewer than:
	ProcessGroupId           string        `json:"-"`
}

func (me *Settings) Name() string {
	return me.ProcessGroupId
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"process_group": {
			Type:        schema.TypeString,
			Description: "The process group ID for availability monitoring",
			Required:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Enable process group availability monitoring",
			Required:    true,
		},
		"alerting_mode": {
			Type:        schema.TypeString,
			Description: "**if any process becomes unavailable:**\nDynatrace will open a new problem if a single process in this group shuts down or crashes. \n\n**if minimum threshold is not met:**\nDynatrace tracks the number of process instances that comprise this process group and treats the group as a cluster. This setting enables you to define a minimum number of process instances that must be available. A problem will be opened if this process group has fewer than the minimum number of required process instances. \n\n Details of the related impact on service requests will be included in the problem summary.\n\n**Note:** If a process is intentionally shutdown or retired while this setting is active, you'll need to manually close the problem.",
			Optional:    true,
		},
		"minimum_instance_threshold": {
			Type:        schema.TypeInt,
			Description: "Open a new problem if the number of active process instances in the group is fewer than:",
			Optional:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	if err := properties.EncodeAll(map[string]any{
		"process_group":              me.ProcessGroupId,
		"enabled":                    me.Enabled,
		"alerting_mode":              me.AlertingMode,
		"minimum_instance_threshold": me.MinimumInstanceThreshold,
	}); err != nil {
		return err
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"process_group":              &me.ProcessGroupId,
		"enabled":                    &me.Enabled,
		"alerting_mode":              &me.AlertingMode,
		"minimum_instance_threshold": &me.MinimumInstanceThreshold,
	})
}

func (me *Settings) SetScope(scope string) {
	me.ProcessGroupId = scope
}

func (me *Settings) GetScope() string {
	return me.ProcessGroupId
}

func (me *Settings) Store() ([]byte, error) {
	var data []byte
	var err error
	if data, err = json.Marshal(me); err != nil {
		return nil, err
	}
	m := map[string]json.RawMessage{}
	if err = json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	if data, err = json.Marshal(me.ProcessGroupId); err != nil {
		return nil, err
	}
	m["process_group"] = data
	return json.MarshalIndent(m, "", "  ")
}

func (me *Settings) Load(data []byte) error {
	if err := json.Unmarshal(data, &me); err != nil {
		return err
	}

	c := struct {
		ProcessGroupId string `json:"process_group"`
	}{}
	if err := json.Unmarshal(data, &c); err != nil {
		return err
	}
	me.ProcessGroupId = c.ProcessGroupId

	return nil
}
