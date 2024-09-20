resource "dynatrace_devobs_data_masking" "#name#" {
  enabled             = false
  replacement_pattern = "#name#"
  replacement_type    = "STRING"
  rule_name           = "#name#"
  rule_type           = "VAR_NAME"
  rule_var_name       = "#name#"
}
