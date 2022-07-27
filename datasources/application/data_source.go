package application

import (
	"fmt"

	applicationapi "github.com/dtcookie/dynatrace/api/config/topology/application"
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
				Description: "Required tags of the application to find",
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
	apiService := applicationapi.NewService(conf.DTNonConfigEnvURL, conf.APIToken)
	applicationList, err := apiService.List()
	if err != nil {
		return err
	}
	if len(applicationList) > 0 {
		for _, application := range applicationList {
			if name == application.DisplayName {
				keys := []string{}
				for _, tag := range application.Tags {
					keys = append(keys, tag.Key)
				}
				if isSubSetOf(keys, tags) {
					d.SetId(application.EntityId)
					return nil
				}
			}
		}
	}

	return fmt.Errorf("no application with name '%s' and tag(s) %v found", name, tags)
}
