resource "dynatrace_vulnerability_third_party" "#name#"{
  enabled  = true
  mode     = "MONITORING_OFF"
  operator = "EQUALS"
  property = "PROCESS_TAG"
  value    = "#name#"
}