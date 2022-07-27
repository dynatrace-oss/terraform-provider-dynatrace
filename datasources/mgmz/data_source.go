package mgmz

import (
	"fmt"

	mgmzapi "github.com/dtcookie/dynatrace/api/config/managementzones"
	"github.com/dtcookie/hcl"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/hcl2sdk"
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
	if len(mgmzList) > 0 {
		for _, mgmz := range mgmzList {
			if name == mgmz.Name {
				d.SetId(mgmz.ID)
				return nil
			}
		}
	}

	return fmt.Errorf("no management zone with name '%s' found", name)
}
