package processgroup

import (
	"fmt"

	processgroupapi "github.com/dtcookie/dynatrace/api/config/topology/processgroup"
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
				Description: "Required tags of the process group to find",
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
	apiService := processgroupapi.NewService(conf.DTNonConfigEnvURL, conf.APIToken)
	processGroupList, err := apiService.List()
	if err != nil {
		return err
	}
	if len(processGroupList) > 0 {
		for _, processGroup := range processGroupList {
			if name == processGroup.DisplayName {
				keys := []string{}
				for _, tag := range processGroup.Tags {
					keys = append(keys, tag.Key)
				}
				if isSubSetOf(keys, tags) {
					d.SetId(processGroup.EntityId)
					return nil
				}
			}
		}
	}

	return fmt.Errorf("no process group with name '%s' and tag(s) %v found", name, tags)
}
