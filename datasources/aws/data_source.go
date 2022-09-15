package aws

import (
	awsapi "github.com/dtcookie/dynatrace/api/config/credentials/aws"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read:   DataSourceRead,
		Schema: map[string]*schema.Schema{},
	}
}

func DataSourceRead(d *schema.ResourceData, m interface{}) error {
	conf := m.(*config.ProviderConfiguration)
	apiService := awsapi.NewService(conf.DTenvURL, conf.APIToken)
	awsIAMExternalID, err := apiService.GetIAMExternalID()
	if err != nil {
		return err
	}
	d.SetId(awsIAMExternalID)
	return nil
}
