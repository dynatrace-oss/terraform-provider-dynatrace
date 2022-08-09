package processgroup

import (
	"fmt"

	processgroupapi "github.com/dtcookie/dynatrace/api/config/topology/processgroup"
	tagapi "github.com/dtcookie/dynatrace/api/config/topology/tag"
	"github.com/dtcookie/hcl"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/config"
	dscommon "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources"
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

func DataSourceRead(d *schema.ResourceData, m interface{}) error {
	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	var tagList []interface{}
	var tags []tagapi.Tag
	if v, ok := d.GetOk("tags"); ok {
		sTags := v.(*schema.Set)
		tagList = sTags.List()
		dscommon.StringsToTags(tagList, &tags)
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
				if dscommon.TagSubsetCheck(processGroup.Tags, tags) {
					d.SetId(processGroup.EntityId)
					return nil
				}
			}
		}
	}

	return fmt.Errorf("no process group with name '%s' and tag(s) %v found", name, tagList)
}
