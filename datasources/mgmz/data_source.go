package mgmz

import (
	mgmzapi "github.com/dtcookie/dynatrace/api/config/managementzones"
	api20 "github.com/dtcookie/dynatrace/api/config/v2/managementzones"
	"github.com/dtcookie/hcl"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/hcl2sdk"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceRead,
		Schema: hcl2sdk.Convert(map[string]*hcl.Schema{
			"name": {
				Type:     hcl.TypeString,
				Required: true,
			},
			"settings_20_id": {
				Type:     hcl.TypeString,
				Computed: true,
			},
		}),
	}
}

func DataSourceRead(d *schema.ResourceData, m interface{}) error {
	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	conf := m.(*config.ProviderConfiguration)
	apiService := mgmzapi.NewService(conf.DTenvURL, conf.APIToken)
	mgmzList, err := apiService.ListAll()
	if err != nil {
		return err
	}
	found := false
	if len(mgmzList) > 0 {
		for _, mgmz := range mgmzList {
			if name == mgmz.Name {
				d.SetId(mgmz.ID)
				found = true
				break
			}
		}
	}
	if !found {
		id := "not-found-" + uuid.New().String()
		d.SetId(id)
		d.Set("settings_20_id", id)
		return nil
	}

	client20 := api20.NewService(conf.DTApiV2URL, conf.APIToken)
	stubs, err := client20.List()
	if err != nil {
		return err
	}
	if len(stubs) > 0 {
		for _, stub := range stubs {
			if name == stub.Name {
				d.Set("settings_20_id", stub.ID)
			}
		}
	}

	return nil
	// return fmt.Errorf("no management zone with name '%s' found", name)
}
