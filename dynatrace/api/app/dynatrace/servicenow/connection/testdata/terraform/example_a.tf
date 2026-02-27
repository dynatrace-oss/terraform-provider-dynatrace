resource "dynatrace_servicenow_connection" "user_password" {
  name     = "#name#"
  url      = "https://www.example.com"
  type     = "basic"
  user     = "#name#"
  password = "#######"
}
