resource "dynatrace_platform_slo" "#name#" {
  name        = "#name#"
  description = "Measures the CPU usage of selected hosts over time."
  criteria {
    criteria_detail {
      target         = 95
      timeframe_from = "now-7d"
      timeframe_to   = "now"
    }
  }
  sli_reference {
    template_id = "e2J1aWx0aW46aW50ZXJuYWwuc2VydmljZS5sZXZlbC5vYmplY3RpdmUudGVtcGxhdGVzLCBpZDogNTJjNzRmMGItNzY1ZS00NWRiLWFmMGQtN2E1MDdjNGY0YjRlfQ=="
    variables {
      sli_reference_variable {
        name  = "hosts"
        value = "\"HOST-1234567890000000\""
      }
    }
  }
}
