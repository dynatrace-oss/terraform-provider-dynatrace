package alerting

import (
	alertingapi "github.com/dtcookie/dynatrace/api/config/alerting"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/config"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceRead,
		Schema: map[string]*schema.Schema{
			"profiles": {
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

func DataSourceRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(uuid.New().String())
	conf := m.(*config.ProviderConfiguration)
	apiService := alertingapi.NewService(conf.DTenvURL, conf.APIToken)
	stubList, err := apiService.List()
	if err != nil {
		return err
	}
	if len(stubList.Values) > 0 {
		profiles := map[string]interface{}{}
		for _, stub := range stubList.Values {
			profiles[stub.Name] = stub.ID
		}
		d.Set("profiles", profiles)
	}
	return nil
}
