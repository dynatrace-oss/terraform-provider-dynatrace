resource "dynatrace_iam_service_user" "wf_user" {
  name = "#name#"
}

resource "dynatrace_automation_workflow" "workflow_with_davis_event_trigger" {
  description = "Desc"
  actor       = dynatrace_iam_service_user.wf_user.id
  # The owner can only be changed if the workflow is public (private = false).
  owner       = dynatrace_iam_service_user.wf_user.id
  private     = false
  title       = "#name#"
  tasks {
    task {
      name        = "http_request_1"
      description = "Issue an HTTP request to any API"
      action      = "dynatrace.automations:http-function"
      active      = true
      input       = jsonencode({
        "method": "GET",
        "url": "https://www.example.com/"
      })
      position {
        x = 0
        y = 1
      }
      retry {
        count = "3"
        delay = "1000"
        failed_loop_iterations_only = false
      }
    }
    task {
      name        = "http_request_2"
      description = "Issue an HTTP request to any API"
      action      = "dynatrace.automations:http-function"
      active      = false
      timeout     = 50000
      input       = jsonencode({
        "method": "GET",
        "url": "https://www.example.com/"
      })
      conditions {
        custom = ""
        states = {
          http_request_1   = "SUCCESS"
          run_javascript_1 = "OK"
        }
      }
      position {
        x = -1
        y = 2
      }
    }
    task {
      name        = "http_request_3"
      description = "Issue an HTTP request to any API"
      action      = "dynatrace.automations:http-function"
      active      = false
      input       = jsonencode({
        "method": "GET",
        "url": "https://www.example.com"
      })
      conditions {
        custom = "{{http_request_1}}"
        states = {
          http_request_2 = "OK"
        }
      }
      position {
        x = 0
        y = 3
      }
    }
    task {
      name        = "run_javascript_1"
      description = "Build a custom task running js Code"
      action      = "dynatrace.automations:run-javascript"
      active      = false
      input       = jsonencode({
        "script": "// optional import of sdk modules\nimport { execution } from '@dynatrace-sdk/automation-utils';\n\nexport default async function ({ execution_id }) {\n  // your code goes here\n  // e.g. get the current execution\n  const ex = await execution(execution_id);\n  console.log('Automated script execution on behalf of', ex.trigger);\n  \n  return { triggeredBy: ex.trigger };\n}"
      })
      position {
        x = -2
        y = 1
      }
    }
  }
  trigger {
    event {
      active = false
      config {
        davis_event {
          entity_tags = {
            asdf = ""
          }
          entity_tags_match  = "all"
          on_problem_close = false
          custom_filter = "matchesPhrase(custom.event.type, \"DEPLOY\")"
        }
      }
    }
  }
}
