resource "dynatrace_notification" "#name#" {
  web_hook {
    name                   = "#name#"
    accept_any_certificate = false
    active                 = true
    alerting_profile       = "c21f969b-5f03-333d-83e0-4f8f136e7682"
    payload                = "{\n\"State\":\"{State}\",\n\"ProblemID\":\"{ProblemID}\",\n\"ProblemTitle\":\"{ProblemTitle}\",\n\"ProblemDetails\": {ProblemDetailsJSON}\n}"
    url                    = "http://webhook.site/414b58b0-aafb-4900-a57d-b90d25163cac"
  }
}
