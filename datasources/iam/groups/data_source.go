package groups

import (
	"fmt"

	"github.com/dtcookie/dynatrace/api/accounts/iam"
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
		}),
	}
}

func DataSourceRead(d *schema.ResourceData, m interface{}) error {
	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	conf := m.(*config.ProviderConfiguration)
	apiService := iam.NewGroupService(conf.IAM.ClientID, conf.IAM.AccountID, conf.IAM.ClientSecret)
	stubs, err := apiService.List()
	if err != nil {
		return err
	}
	if len(stubs) > 0 {
		for _, stub := range stubs {
			if name == stub.Name {
				d.SetId(stub.UUID)
				return nil
			}
		}
	}

	return fmt.Errorf("no iam user group with name '%s' found", name)
}
