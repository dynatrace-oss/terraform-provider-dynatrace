package host

import (
	"fmt"

	hostapi "github.com/dtcookie/dynatrace/api/config/topology/host"
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
				Description: "Required tags of the host to find",
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
	apiService := hostapi.NewService(conf.DTNonConfigEnvURL, conf.APIToken)
	hostList, err := apiService.List()
	if err != nil {
		return err
	}
	if len(hostList) > 0 {
		for _, host := range hostList {
			if name == host.DisplayName {
				if dscommon.TagSubsetCheck(host.Tags, tags) {
					d.SetId(host.EntityId)
					return nil
				}
			}
		}
	}

	return fmt.Errorf("no host with name '%s' and tag(s) %v found", name, tagList)
}
