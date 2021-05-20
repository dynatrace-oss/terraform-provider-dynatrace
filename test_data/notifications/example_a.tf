resource "dynatrace_notification" "#name#" {
  service_now {
    name             = "#name#"
    active           = true
    alerting_profile = "c21f969b-5f03-333d-83e0-4f8f136e7682"
    send_events      = false
    send_incidents   = false
    username         = "admin"
    instance_name    = "dev87541"
    message          = "{State} {ProblemImpact} Problem {ProblemID}: {ProblemTitle}"
    password         = "pw2"
  }
}
