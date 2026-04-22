resource "dynatrace_appsec_notification" "Terraform_Security_Problem_Webhook_Test" {
  type                                    = "WEBHOOK"
  enabled                                 = true
  display_name                            = "Terraform Security Problem Webhook Test"
  security_problem_based_alerting_profile = dynatrace_vulnerability_alerting.alert.id
  trigger                                 = "SECURITY_PROBLEM"
  security_problem_based_webhook_payload {
    payload = jsonencode({
          "DavisSecurityScore": "{DavisSecurityScore}",
          "SecurityProblemId": "{SecurityProblemId}",
          "SecurityProblemUrl": "{SecurityProblemUrl}",
          "Severity": "{Severity}"
      })
  }
  webhook_configuration {
    accept_any_certificate = false
    url                    = "https://www.dynatrace.com"
  }
}

resource "dynatrace_vulnerability_alerting" "alert" {
  name                   = "#name#"
  enabled                = true
  enabled_risk_levels    = ["LOW", "MEDIUM", "HIGH", "CRITICAL"]
  enabled_trigger_events = ["SECURITY_PROBLEM_OPENED"]
}
