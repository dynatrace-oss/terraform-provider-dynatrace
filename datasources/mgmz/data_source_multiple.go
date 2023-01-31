package mgmz

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
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
					},
				},
			},
		},
	}
}

func DataSourceReadMultiple(d *schema.ResourceData, m any) error {
	service := export.Service(config.Credentials(m), export.ResourceTypes.ManagementZoneV2)
	var stubs settings.Stubs
	var err error
	if stubs, err = service.List(); err != nil {
		return err
	}
	d.SetId(staticID)
	values := []map[string]any{}
	for _, stub := range stubs {
		values = append(values, map[string]any{
			"id":        stub.ID,
			"legacy_id": settings.LegacyID(stub.ID),
			"name":      stub.Name,
		})
	}
	d.Set("values", values)
	return nil
}
