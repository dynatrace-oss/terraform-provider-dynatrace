package entity

import (
	"fmt"

	"github.com/dtcookie/dynatrace/api/config/v2/entities"
	"github.com/dtcookie/hcl"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/hcl2sdk"
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
	if len(*settings.Entities) > 0 {
		for _, entity := range *settings.Entities {
			if name == *entity.DisplayName {
				d.SetId(*entity.EntityId)
				return nil
			}
		}
	}

	return fmt.Errorf("no entity with name '%s' found", name)
}
