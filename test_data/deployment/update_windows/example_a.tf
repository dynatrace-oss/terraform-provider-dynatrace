resource "dynatrace_update_windows" "#name#" {
  name       = "#name#"
  enabled    = true
  recurrence = "ONCE"
  once_recurrence {
    recurrence_range {
      end   = "2023-02-15T04:00:00Z"
      start = "2023-02-15T02:00:00Z"
    }
  }
}