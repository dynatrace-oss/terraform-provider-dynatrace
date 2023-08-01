resource "dynatrace_appsec_notification" "Terraform_Security_Problem_Webhook_Test" {
  type                                    = "WEBHOOK"
  enabled                                 = true
  display_name                            = "Terraform Security Problem Webhook Test"
  security_problem_based_alerting_profile = "vu9U3hXa3q0AAAABACxidWlsdGluOmFwcHNlYy5ub3RpZmljYXRpb24tYWxlcnRpbmctcHJvZmlsZQAGdGVuYW50AAZ0ZW5hbnQAJDMyMDhkNWMyLTFlZmYtMzk5My1iNjMwLWI0MjQ5N2U4MDQ2Nr7vVN4V2t6t"
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