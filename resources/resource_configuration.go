package resources

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

type ResourceConfiguration struct {
	Resource func() *schema.Resource
}
