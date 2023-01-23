resource "dynatrace_notification" "#name#" {
  pager_duty {
    name             = "#name#"
    account          = ";l;ll;"
    active           = true
    alerting_profile = "f75e68ef-aca7-3a07-9c21-94eb00ecfc56"
    service_api_key  = "#######"
    service_name     = "lklklk"
  }
}