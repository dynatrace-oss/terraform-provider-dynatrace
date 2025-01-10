package appdetection

import (
	"context"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	appdetectionsrv "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/appdetection"
	appdetection "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/appdetection/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/cache"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: logging.EnableDSCtx(DataSourceRead),
		Schema: map[string]*schema.Schema{
			"values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application Detection Rule ID",
						},
						"application_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application ID",
						},
						"matcher": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Matcher",
						},
						"pattern": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Pattern",
						},
					},
				},
			},
		},
	}
}

func DataSourceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	d.SetId("dynatrace_application_detection_rules")
	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}

	service := cache.Read[*appdetection.Settings](appdetectionsrv.Service(creds), true)
	var stubs api.Stubs
	if stubs, err = service.List(ctx); err != nil {
		return diag.FromErr(err)
	}
	values := []map[string]any{}
	for _, stub := range stubs {
		stubValue := stub.Value.(*appdetection.Settings)
		values = append(values, map[string]any{
			"id":             stub.ID,
			"application_id": stubValue.ApplicationID,
			"matcher":        stubValue.Matcher,
			"pattern":        stubValue.Pattern,
		})
	}
	d.Set("values", values)
	return diag.Diagnostics{}
}
