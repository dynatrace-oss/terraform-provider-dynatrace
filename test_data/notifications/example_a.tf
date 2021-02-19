resource "dynatrace_notification" "#name#" {
  service_now_notification_config {
    name             = "#name#"
    type             = "SERVICE_NOW"
    active           = true
    alerting_profile = "c21f969b-5f03-333d-83e0-4f8f136e7682"
    send_events      = true
    send_incidents   = true
    url              = ""
    username         = "admin"
    instance_name    = "dev87541"
    message          = "{State} {ProblemImpact} Problem {ProblemID}: {ProblemTitle}"
    password         = "pw"
  }
}
