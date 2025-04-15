data "dynatrace_platform_slo_template" "HostCPU" {
  name = "Host CPU usage utilization"
}

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
    template_id = data.dynatrace_platform_slo_template.HostCPU.id
    variables {
      sli_reference_variable {
        name  = "hosts"
        value = "\"HOST-1234567890000000\""
      }
    }
  }
}
