package locations

import (
	"github.com/dtcookie/hcl"
	"github.com/dtcookie/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/hcl2sdk"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceRead,
		Schema: hcl2sdk.Convert(map[string]*hcl.Schema{
			"id": {
				Type:     hcl.TypeString,
				Optional: true,
			},
			"name": {
				Type:     hcl.TypeString,
				Optional: true,
			},
			"locations": {
				Type:     hcl.TypeList,
				MaxItems: 1,
				Elem:     &hcl.Resource{Schema: new(Locations).Schema()},
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func DataSourceRead(d *schema.ResourceData, m interface{}) error {
	var id *string
	var name *string

	if v, ok := d.GetOk("id"); ok {
		d.SetId(v.(string))
		id = opt.NewString(v.(string))
	} else {
		d.SetId(uuid.New().String())
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

	locs := Locations{}
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
		locs = append(locs, &Location{
			ID:   location.ID,
			Name: location.Name,
			Type: LocationType(string(location.Type)),
		})
	}
	marshalled, err := locs.MarshalHCL()
	if err != nil {
		return err
	}
	them := map[string]interface{}{}
	for k, v := range marshalled {
		them[k] = v
	}
	d.Set("locations", []interface{}{them})

	return nil
}
