package locations

import (
	"fmt"

	"github.com/dtcookie/dynatrace/rest"
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
	var typeLoc *string
	var status *string
	var stage *string
	var cloudPlatform *string
	var ips []string

	if v, ok := d.GetOk("entity_id"); ok {
		d.SetId(v.(string))
		id = opt.NewString(v.(string))
	}
	if v, ok := d.GetOk("name"); ok {
		name = opt.NewString(v.(string))
	}
	if v, ok := d.GetOk("type"); ok {
		typeLoc = opt.NewString(v.(string))
	}
	if v, ok := d.GetOk("status"); ok {
		status = opt.NewString(v.(string))
	}
	if v, ok := d.GetOk("stage"); ok {
		stage = opt.NewString(v.(string))
	}
	if v, ok := d.GetOk("cloud_platform"); ok {
		cloudPlatform = opt.NewString(v.(string))
	}
	if v, ok := d.GetOk("ips"); ok {
		if vt, ok := v.([]string); ok {
			ips = vt
		}
	}

	conf := m.(*config.ProviderConfiguration)
	apiService := NewService(conf.DTNonConfigEnvURL, conf.APIToken)
	rest.Verbose = false
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
		if typeLoc != nil {
			if *typeLoc != string(location.Type) {
				continue
			}
		}
		if status != nil {
			if *status != string(*location.Status) {
				continue
			}
		}
		if stage != nil {
			if *stage != string(*location.Stage) {
				continue
			}
		}
		if cloudPlatform != nil {
			if *cloudPlatform != string(*location.CloudPlatform) {
				continue
			}
		}
		if len(ips) > 0 {
			if !subsetCheck(location.IPs, ips) {
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
			if k == "entity_id" {
				d.SetId(v.(string))
			}
			d.Set(k, v)
			// if k != "id" {
			// 	d.Set(k, v)
			// } else {
			// 	d.SetId(v.(string))
			// 	return nil
			// }
		}

		return nil
	}

	return fmt.Errorf("no matching synthetic location found")
}

// subsetCheck verifies that the input strings are a subset of source strings
// Arguments: source slice of strings, input slice of strings
// Return: true if subset, false if not
func subsetCheck(source []string, input []string) bool {
	for _, inputString := range input {
		found := false
		for _, sourceString := range source {
			if inputString == sourceString {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
