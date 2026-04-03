resource "dynatrace_automation_workflow" "empty_values" {
  // description, guide, result not set
  type                   = "STANDARD"
  hourly_execution_limit = 1000
  owner_type             = "USER"
  private                = false
  title                  = "#name#"
  tasks {
    task {
      name = "test_1"
      action = "dynatrace.automations:execute-dql-query"
      input = jsonencode({})
      position {
        x = 0
        y = 1
      }
      // no description
      active = false
    }
  }
  trigger {}
}
