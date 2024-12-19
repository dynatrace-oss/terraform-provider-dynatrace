resource "dynatrace_ms365_email_connection" "#name#"{
  name     = "#name#"
  type     = "client_secret"
  tenant_id     = "123e4567-e89b-12d3-a456-426614174000"
  client_id     = "123e4567-e89b-12d3-a456-426614174000"
  client_secret   = "######"
  from_address    = "random.email@terraform.com"
}