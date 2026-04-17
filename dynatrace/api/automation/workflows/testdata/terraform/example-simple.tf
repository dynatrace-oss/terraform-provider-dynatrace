resource "dynatrace_automation_workflow" "simple_workflow" {
  type                   = "SIMPLE"
  hourly_execution_limit = 1000
  owner_type             = "USER"
  private                = false
  title                  = "#name#"
  tasks {
    task {
      name        = "http_request_1"
      description = "Issue an HTTP request to any API change"
      action      = "dynatrace.automations:http-function"
      active      = true
    }
  }
  trigger {
  }
}
