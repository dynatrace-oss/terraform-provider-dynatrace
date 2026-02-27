resource "dynatrace_iam_service_user" "wf_user" {
  name = "#name#"
}

resource "dynatrace_automation_workflow" "workflow_with_davis_problem_trigger" {
  description = "#name#"
  actor       = dynatrace_iam_service_user.wf_user.id
  private     = true
  title       = "#name#"
  tasks {
    task {
      name = "test_1"
      action = "dynatrace.automations:execute-dql-query"
      position {
        x = 0
        y = 1
      }
      description = "Dummy task"
      active = false
    }
  }
  trigger {
    event {
      active = true
      config {
        davis_problem {
          categories {
            error = true
          }
          entity_tags = {
            unknowntag = "this-tag-does-not-exist"
          }
          entity_tags_match  = "all"
          custom_filter = "matchesPhrase(custom.event.type, \"DEPLOY\")"
          analysis_ready = true
        }
      }
    }
  }
}
