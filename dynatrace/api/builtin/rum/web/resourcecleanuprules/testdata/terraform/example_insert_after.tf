resource "dynatrace_web_app_resource_cleanup" "first-instance" {
  name               = "#name#"
  regular_expression = "([a-z]+)TerraformExample"
  replace_with       = "$1"
}

resource "dynatrace_web_app_resource_cleanup" "second-instance" {
  name               = "#name#-second"
  regular_expression = "([a-z]+)TerraformExample2"
  replace_with       = "$1"
  insert_after = dynatrace_web_app_resource_cleanup.first-instance.id
}