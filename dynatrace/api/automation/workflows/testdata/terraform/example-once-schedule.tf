# Exercises the "once" schedule trigger type via the `at` field

resource "dynatrace_automation_workflow" "workflow_with_once_schedule" {
  description = "#name#"
  private     = true
  title       = "#name#"
  tasks {}
  trigger {
    schedule {
      active = true
      trigger {
        at = "2100-01-01T14:30:00"
      }
    }
  }
}
