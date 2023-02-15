resource "dynatrace_process_group_simple_detection" "#name#" {
  enabled             = false
  group_identifier    = "GroupIdentifierExample"
  instance_identifier = "InstanceIdentifierExample"
  process_type        = "PROCESS_TYPE_GO"
  rule_type           = "prop"
}