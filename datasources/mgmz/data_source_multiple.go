package mgmz

import (
	"sort"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	managementzonessrv "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/managementzones"
	managementzones "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/managementzones/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/cache"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const staticID = "46465fe6-70cb-4564-864f-c3344caae5c0"

func DataSourceMultiple() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceReadMultiple,
		Schema: map[string]*schema.Schema{
			"values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the Management Zone when referred to as a Settings 2.0 resource (e.g. from within `dynatrace_slack_notification`)",
						},
						"legacy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the Management Zone when referred to as a Configuration API resource (e.g. from within `dynatrace_notification`)",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the Management Zone",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the Management Zone",
						},
					},
				},
			},
		},
	}
}

func DataSourceReadMultiple(d *schema.ResourceData, m any) error {
	service := cache.Read[*managementzones.Settings](managementzonessrv.Service(config.Credentials(m)), true)
	var stubs api.Stubs
	var err error
	if stubs, err = service.List(); err != nil {
		return err
	}
	d.SetId(staticID)
	values := []map[string]any{}
	sort.SliceStable(stubs, func(i, j int) bool {
		return stubs[i].Name < stubs[j].Name
	})
	for _, stub := range stubs {
		stubValue := stub.Value.(*managementzones.Settings)
		description := ""
		if stubValue != nil && stubValue.Description != nil {
			description = *stubValue.Description
		}
		values = append(values, map[string]any{
			"id":          stub.ID,
			"legacy_id":   stubValue.LegacyID,
			"name":        stub.Name,
			"description": description,
		})
	}
	d.Set("values", values)
	return nil
}
