resource "dynatrace_notification" "#name#" {
  pager_duty {
    name             = "#name#"
    account          = ";l;ll;"
    active           = true
    alerting_profile = dynatrace_alerting_profile.Default.id
    service_api_key  = "#######"
    service_name     = "lklklk"
  }
}

resource "dynatrace_alerting_profile" "Default" {
  display_name = "#name#"
}
