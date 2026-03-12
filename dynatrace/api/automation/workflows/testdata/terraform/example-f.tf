resource "dynatrace_automation_workflow" "workflow_with_interval_schedule" {
  description = "#name#"
  private     = true
  title       = "#name#"
  tasks {}
  trigger {
    schedule {
      trigger {
        between_end = "23:59"
        between_start = "23:59"
        interval_minutes = "15"
      }
      active = true
      filter_parameters {
        earliest_start_time = "00:00"
        earliest_start = "2023-07-01"
        count = "1"
        exclude_dates = ["2023-07-02"]
        include_dates = ["2100-01-01"]
        until = "2100-01-02"
      }
    }
  }
}
