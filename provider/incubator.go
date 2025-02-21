package provider

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/dql"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func incubator(provider *schema.Provider) {
	provider.DataSourcesMap["dynatrace_dql"] = dql.DataSource()
}
