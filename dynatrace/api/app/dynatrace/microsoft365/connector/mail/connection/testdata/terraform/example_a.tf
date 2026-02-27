resource "dynatrace_ms365_email_connection" "connection" {
  name     = "#name#"
  type     = "client_secret"
  tenant_id     = "00000000-0000-0000-0000-000000000000"
  client_id     = "00000000-0000-0000-0000-000000000000"
  client_secret   = "######"
  from_address    = "random.email@terraform.com"
}
