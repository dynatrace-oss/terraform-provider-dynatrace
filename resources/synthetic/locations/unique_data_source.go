package locations

import (
	"github.com/dtcookie/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/hcl2sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func UniqueDataSource() *schema.Resource {
	return &schema.Resource{
		Read:   UniqueDataSourceRead,
		Schema: hcl2sdk.Convert(new(Location).Schema()),
	}
}

func UniqueDataSourceRead(d *schema.ResourceData, m interface{}) error {
	var id *string
	var name *string

	if v, ok := d.GetOk("id"); ok {
		d.SetId(v.(string))
		id = opt.NewString(v.(string))
		// } else {
		// 	d.SetId(uuid.New().String())
	}

	if v, ok := d.GetOk("name"); ok {
		name = opt.NewString(v.(string))
	}

	conf := m.(*config.ProviderConfiguration)
	apiService := NewService(conf.DTNonConfigEnvURL, conf.APIToken)
	locationList, err := apiService.List()
	if err != nil {
		return err
	}

	for _, location := range locationList.Locations {
		if id != nil {
			if *id != location.ID {
				continue
			}
		}
		if name != nil {
			if *name != location.Name {
				continue
			}
		}
		loc := &Location{
			ID:            location.ID,
			Name:          location.Name,
			Type:          LocationType(string(location.Type)),
			Status:        location.Status,
			CloudPlatform: location.CloudPlatform,
			IPs:           location.IPs,
		}
		marshalled, err := loc.MarshalHCL()
		if err != nil {
			return err
		}
		for k, v := range marshalled {
			if k != "id" {
				d.Set(k, v)
			} else {
				d.SetId(v.(string))
			}
		}

		break
	}
	return nil
}
