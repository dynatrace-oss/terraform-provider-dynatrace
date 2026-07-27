resource "dynatrace_automation_workflow_slack" "default" {
  name  = "#name#"
  token = "#######"
}

resource "dynatrace_automation_workflow_slack" "external_approval" {
  name              = "#name#"
  token             = "#######"
  external_approval = true
  signing_secret    = "#######"
}
