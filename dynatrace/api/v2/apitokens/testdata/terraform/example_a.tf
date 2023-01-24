resource "dynatrace_api_token" "#name#" {
  name                    = "#name#"
  enabled                 = false
  # personal_access_token = false
  scopes                  = [ "geographicRegions.read" ]
}
