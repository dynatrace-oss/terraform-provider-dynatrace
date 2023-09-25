resource "dynatrace_ag_token" "#name#" {
  type = "ENVIRONMENT"
  expiration_date = "now+3d"
  name = "#name#"  
}