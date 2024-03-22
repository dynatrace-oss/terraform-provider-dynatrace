resource "dynatrace_automation_workflow_jira" "first-instance"{
  name     = "#name#"
  type     = "basic"
  password = "#######"
  url      = "https://www.example.com"
  user     = "user2"
}

resource "dynatrace_automation_workflow_jira" "second-instance"{
  name     = "#name#-second"
  type     = "basic"
  password = "#######-second"
  url      = "https://www.example-second.com"
  user     = "user2-second"
  insert_after = "${dynatrace_automation_workflow_jira.first-instance.id}"
}

