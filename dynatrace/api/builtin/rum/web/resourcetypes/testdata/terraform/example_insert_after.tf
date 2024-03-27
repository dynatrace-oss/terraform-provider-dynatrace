resource "dynatrace_web_app_resource_types" "first-instance" {
  primary_resource_type = "OTHER"
  regular_expression    = ".Terraform"
}

resource "dynatrace_web_app_resource_types" "second-instance" {
  primary_resource_type = "OTHER"
  regular_expression    = ".Terraform-2"
  insert_after = dynatrace_web_app_resource_types.first-instance.id
}