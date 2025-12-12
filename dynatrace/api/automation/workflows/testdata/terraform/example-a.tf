resource "dynatrace_automation_workflow" "#name#" {
  description = "#name#"
  actor       = "703d65c0-4aff-45d9-8b34-2c6f5f17bb8e"
  owner       = "703d65c0-4aff-45d9-8b34-2c6f5f17bb8e"
  private     = true
  title       = "#name#"
  tasks {
    task {
      name        = "http_request_3"
      description = "Issue an HTTP request to any API"
      action      = "dynatrace.automations:http-function"
      active      = false
      input       = jsonencode({
              "method": "GET",
              "url": "https://www.third-task.com"
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
    task {
      name        = "http_request_1"
      description = "Issue an HTTP request to any API"
      action      = "dynatrace.automations:http-function"
      active      = true
      input       = jsonencode({
              "method": "GET",
              "url": "https://www.google.at/"
        })
      position {
        x = 0
        y = 1
      }
    }
    task {
      name        = "http_request_2"
      description = "Issue an HTTP request to any API"
      action      = "dynatrace.automations:http-function"
      active      = false
      input       = jsonencode({
              "method": "GET",
              "url": "https://www.second-task.com/"
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
  }
  trigger {
    event {
      # active = false
      config {
        davis_event {
          entity_tags = {
            asdf = ""
          }
          entity_tags_match  = "all"
          # on_problem_close = false
          types              = [ "CUSTOM_ANNOTATION" ]
          custom_filter = "matchesPhrase(custom.event.type, \"DEPLOY\")"
        }
      }
    }
  }
}
