resource "dynatrace_servicenow_connection" "#name#"{
  name    = "#name#"
  url     = "https://www.#name#.com"
  type    = "client-credentials"
  client_id    = "#name#"
  client_secret   = "#######"
  }