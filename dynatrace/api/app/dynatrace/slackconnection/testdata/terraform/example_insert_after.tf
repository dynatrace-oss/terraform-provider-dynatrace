resource "dynatrace_automation_workflow_slack" "first-instance" {
  name  = "#name#"
  token = "#######"
}

resource "dynatrace_automation_workflow_slack" "second-instance" {
  name  = "#name#-second"
  token = "#######-second"
  insert_after = "${dynatrace_automation_workflow_slack.first-instance.id}"
}
