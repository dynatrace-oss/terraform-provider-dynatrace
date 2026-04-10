resource "dynatrace_automation_business_calendar" "calendar" {
  description = "#name#"
  title         = "#name#"
  valid_from    = "2023-07-31"
  valid_to      = "2033-07-31"
  week_days     = [1, 2, 3, 4, 5]
  week_start    = 1
  holidays {
    holiday {
      date  = "2023-08-15"
      title = "Holiday 1"
    }
    // edit => recreate due to set-hash change
    holiday {
      date  = "2023-10-26"
      title = "Holiday 2 edit"
    }
    // removed #3
    // add
    holiday {
      date  = "2023-10-27"
      title = "Holiday 4"
    }
  }
}
