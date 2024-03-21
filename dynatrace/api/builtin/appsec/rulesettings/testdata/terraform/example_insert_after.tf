resource "dynatrace_vulnerability_third_party" "first-instance" {
  enabled  = true
  mode     = "MONITORING_OFF"
  operator = "EQUALS"
  property = "PROCESS_TAG"
  value    = "#name#"
}

resource "dynatrace_vulnerability_third_party" "second-instance" {
  enabled      = true
  mode         = "MONITORING_OFF"
  operator     = "EQUALS"
  property     = "PROCESS_TAG"
  value        = "#name#-second"
  insert_after = dynatrace_vulnerability_third_party.first-instance.id
}
