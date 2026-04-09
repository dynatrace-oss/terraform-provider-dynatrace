resource "dynatrace_automation_workflow" "empty_values" {
  description            = "my-description"
  guide                  = "my-guide"
  result                 = jsonencode({ "myprop" : "myValue" })
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
      description = "Dummy task"
      active = false
    }
  }
  trigger {}
}
