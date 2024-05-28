resource "dynatrace_report" "#name#" {
  type                = "DASHBOARD"
  dashboard_id        = "41eae96d-4930-4f44-bbd8-3699f21a8bbf"
  email_notifications = true
  subscriptions {
    month = ["terraform1@dynatrace.com", "terraform2@dynatrace.com"]
    week = ["terraform3@dynatrace.com", "terraform4@dynatrace.com", "terraform5@dynatrace.com"]
  }
}