resource "dynatrace_web_app_resource_cleanup" "#name#" {
  name               = "#name#"
  regular_expression = "([a-z]+)TerraformExample"
  replace_with       = "$1"
}