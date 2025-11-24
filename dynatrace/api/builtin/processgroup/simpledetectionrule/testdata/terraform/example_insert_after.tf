resource "dynatrace_process_group_simple_detection" "first-instance" {
  enabled             = false
  group_identifier    = "GroupIdentifierExample"
  instance_identifier = "InstanceIdentifierExample_#name#"
  process_type        = "PROCESS_TYPE_GO"
  rule_type           = "prop"
}

resource "dynatrace_process_group_simple_detection" "second-instance" {
  enabled             = false
  group_identifier    = "GroupIdentifierExample2"
  instance_identifier = "InstanceIdentifierExample_#name#"
  process_type        = "PROCESS_TYPE_GO"
  rule_type           = "prop"
  insert_after        = dynatrace_process_group_simple_detection.first-instance.id
}
