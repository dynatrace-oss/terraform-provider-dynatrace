package service

import (
	"fmt"
	"strings"

	serviceapi "github.com/dtcookie/dynatrace/api/config/topology/service"
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
				Description: "Required tags of the service to find",
				MinItems:    1,
			},
		}),
	}
}

func isSubSetOf(restSource []serviceapi.Tag, input []serviceapi.Tag) bool {
	for _, inputTag := range input {
		found := false
		for _, restTag := range restSource {
			if restTag.Key == inputTag.Key {
				if restTag.Value == nil && inputTag.Value == nil {
					found = true
					break
				} else if restTag.Value != nil && inputTag.Value != nil && *restTag.Value == *inputTag.Value {
					found = true
					break
				}
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

	var tags []serviceapi.Tag
	if v, ok := d.GetOk("tags"); ok {
		sTags := v.(*schema.Set)
		tagList := sTags.List()
		for _, iTag := range tagList {
			var tag serviceapi.Tag
			if strings.Contains(iTag.(string), "=") {
				tagSplit := strings.Split(iTag.(string), "=")
				tag.Key = tagSplit[0]
				tag.Value = &tagSplit[1]
			} else {
				tag.Key = iTag.(string)
			}
			tags = append(tags, tag)
		}
	}

	conf := m.(*config.ProviderConfiguration)
	apiService := serviceapi.NewService(conf.DTNonConfigEnvURL, conf.APIToken)
	serviceList, err := apiService.List()
	if err != nil {
		return err
	}
	if len(serviceList) > 0 {
		for _, service := range serviceList {
			if name == service.DisplayName {
				if isSubSetOf(service.Tags, tags) {
					d.SetId(service.EntityId)
					return nil
				}
			}
		}
	}

	return fmt.Errorf("no service with name '%s' found", name)
}
