resource "dynatrace_servicenow_connection" "client_credentials" {
  name          = "#name#"
  url           = "https://www.example.com"
  type          = "client-credentials"
  client_id     = "#name#"
  client_secret = "#######"
}
