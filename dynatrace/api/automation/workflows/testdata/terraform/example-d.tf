resource "dynatrace_automation_workflow" "without_tasks" {
  title = "#name#"
  result = "deployment_status"

  tasks {}

  trigger {
    event {
      active = true
      config {
        davis_problem {
          categories {
            error = true
          }
        }
      }
    }
  }
}
