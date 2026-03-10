# Exercises the new fields: hourly_execution_limit, analysis_ready,
# is_deployed, owner_type, input, task_defaults, guide, and result

resource "dynatrace_iam_group" "group" {
  name = "#name#"
}

resource "dynatrace_automation_workflow" "with_group_owner" {
  # Required
  title = "#name#"

  # Optional — new fields under test
  hourly_execution_limit = 500
  is_deployed            = false
  owner = dynatrace_iam_group.group.id
  owner_type             = "GROUP"
  input = jsonencode({
    "environment" : "production",
    "threshold" : 42
  })
  task_defaults = jsonencode({
    "timeout" : 3600
  })
  guide = "##My Guide"
  result = "deployment_status"

  # Optional
  description = "Validates all newly added workflow fields"
  private     = false

  tasks {
    task {
      # Required
      action = "dynatrace.automations:http-function"
      name   = "http_request_1"

      # Optional
      active      = false
      description = "Simple HTTP task for testing"
      input = jsonencode({
        "method" : "GET",
        "url" : "https://www.example.com/"
      })
      position {
        x = 0
        y = 1
      }
    }
  }

  trigger {
    event {
      active = true
      config {
        davis_problem {
          # Required
          categories {
            error = true
          }

          # Optional — new field under test
          analysis_ready = true

          # Optional
          on_problem_close = false
        }
      }
    }
  }
}
