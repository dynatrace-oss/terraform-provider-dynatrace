resource "dynatrace_automation_workflow" "workflow_with_empty_json_encode" {
  description = "#name#"
  private     = true
  title       = "#name#"
  input = jsonencode({})
  tasks {
    task {
      name = "test_1"
      action = "dynatrace.automations:execute-dql-query"
      input = jsonencode({})
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
