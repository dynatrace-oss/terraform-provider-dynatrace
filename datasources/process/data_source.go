package process

import (
	"fmt"

	processapi "github.com/dtcookie/dynatrace/api/config/topology/process"
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
			"tags": {
				Type:        hcl.TypeSet,
				Elem:        &hcl.Schema{Type: hcl.TypeString},
				Optional:    true,
				Description: "Required tags of the process to find",
				MinItems:    1,
			},
		}),
	}
}

func isSubSetOf(source []string, input []string) bool {
	for _, i := range input {
		found := false
		for _, s := range source {
			if s == i {
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

func DataSourceRead(d *schema.ResourceData, m interface{}) error {
	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	var tags []string
	if v, ok := d.GetOk("tags"); ok {
		sTags := v.(*schema.Set)
		tagList := sTags.List()
		for _, iTag := range tagList {
			tags = append(tags, iTag.(string))
		}
	}

	conf := m.(*config.ProviderConfiguration)
	apiService := processapi.NewService(conf.DTNonConfigEnvURL, conf.APIToken)
	processList, err := apiService.List()
	if err != nil {
		return err
	}
	if len(processList) > 0 {
		for _, process := range processList {
			if name == process.DisplayName {
				keys := []string{}
				for _, tag := range process.Tags {
					keys = append(keys, tag.Key)
				}
				if isSubSetOf(keys, tags) {
					d.SetId(process.EntityId)
					return nil
				}
			}
		}
	}

	return fmt.Errorf("no process with name '%s' and tag(s) %v found", name, tags)
}
