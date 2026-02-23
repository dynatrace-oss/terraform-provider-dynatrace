resource "dynatrace_automation_workflow" "#name#" {
  description = "#name#"
  actor       = "703d65c0-4aff-45d9-8b34-2c6f5f17bb8e"
  owner       = "703d65c0-4aff-45d9-8b34-2c6f5f17bb8e"
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
      # active = false
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
