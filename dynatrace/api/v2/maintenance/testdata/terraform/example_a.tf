resource "dynatrace_maintenance" "#name#" {
  enabled = true 
  general_properties {
    name = "#name#" 
    description = "Terraform test execution" 
    type = "PLANNED" 
    disable_synthetic = true 
    suppression = "DETECT_PROBLEMS_AND_ALERT" 
  }
  schedule {
    type = "WEEKLY" 
    weekly_recurrence {
      day_of_week = "MONDAY" 
      recurrence_range {
        end_date = "2022-10-06" 
        start_date = "2022-10-05" 
      }
      time_window {
        end_time = "15:13:00" 
        start_time = "14:13:00" 
        time_zone = "UTC" 
      }
    }
  }
  filters {
    filter {
        entity_type = "HOST" 
        entity_tags = ["KeyTest:ValueTest"] 
    }
  }
}
