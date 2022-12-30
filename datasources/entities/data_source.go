package entities

import (
	"github.com/dtcookie/dynatrace/api/config/v2/entities"
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
			"type": {
				Type:     hcl.TypeString,
				Required: true,
			},
			"entities": {
				Type:     hcl.TypeList,
				MaxItems: 1,
				Elem:     &hcl.Resource{Schema: new(entities.Entities).Schema()},
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func DataSourceRead(d *schema.ResourceData, m interface{}) error {
	var entityType string
	if v, ok := d.GetOk("type"); ok {
		entityType = v.(string)
	}

	conf := m.(*config.ProviderConfiguration)
	apiService := entities.NewService(conf.DTApiV2URL, conf.APIToken)
	settings, err := apiService.List(entityType)
	if err != nil {
		return err
	}
	d.SetId(uuid.New().String())
	if settings.Entities != nil {
		marshalled, err := settings.Entities.MarshalHCL()
		if err != nil {
			return err
		}
		them := map[string]interface{}{}
		for k, v := range marshalled {
			them[k] = v
		}
		d.Set("entities", []interface{}{them})
	}
	return nil
}
