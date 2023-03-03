resource "dynatrace_web_app_beacon_origins" "#name#" {
  matcher = "CONTAINS"
  pattern = "Terraform-Origin"
}