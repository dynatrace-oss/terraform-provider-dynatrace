resource "dynatrace_custom_units" "#name#" {
  name        = "#name#"
  description = "Created by Terraform"
  plural_name = "TerraformUnits"
  symbol      = "T/u"
}