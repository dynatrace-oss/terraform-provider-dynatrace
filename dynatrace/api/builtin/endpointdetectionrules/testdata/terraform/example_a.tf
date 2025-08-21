resource "dynatrace_endpoint_detection_rules" "first-instance" {
  enabled      = true
  scope        = "environment"
  rule {
    description            = "This is a sample description"
    condition              = "matchesValue(k8s.cluster.name,\"#name#-1\")"
    endpoint_name_template = "#name#"
    if_condition_matches   = "DETECT_REQUEST_ON_ENDPOINT"
    rule_name              = "#name#-1"
  }
}

resource "dynatrace_endpoint_detection_rules" "second-instance" {
  insert_after = dynatrace_endpoint_detection_rules.first-instance.id
  enabled      = false
  scope        = "environment"
  rule {
    condition            = "matchesValue(k8s.cluster.name,\"#name#-2\")"
    if_condition_matches = "SUPPRESS_REQUEST"
    rule_name            = "#name#-2"
  }
}
