resource "dynatrace_rum_advanced_correlation" "#name#" {
  matcher = "STARTS_WITH"
  pattern = "Server-Timing-Terraform-2"
}