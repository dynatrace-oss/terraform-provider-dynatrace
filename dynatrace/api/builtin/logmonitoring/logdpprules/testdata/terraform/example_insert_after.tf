resource "dynatrace_log_processing" "first-instance" {
  enabled   = false
  query     = "content=\"terraform\""
  rule_name = "#name#"
  processor_definition {
    rule = <<-EOT
      PARSE(content,"LD 'total: ' INT:total '; failed: ' INT:failed")
      | FIELDS_ADD(failed.percentage: 100.0 * failed / total + '%')
      | FIELDS_REMOVE(total, failed)
    EOT
  }
  rule_testing {
    sample_log = jsonencode({
      "content" : ""
    })
  }
}

resource "dynatrace_log_processing" "second-instance" {
  enabled   = false
  query     = "content=\"terraform\""
  rule_name = "#name#-second"
  processor_definition {
    rule = <<-EOT
      PARSE(content,"LD 'total: ' INT:total '; failed: ' INT:failed")
      | FIELDS_ADD(failed.percentage: 100.0 * failed / total + '%')
      | FIELDS_REMOVE(total, failed)
    EOT
  }
  rule_testing {
    sample_log = jsonencode({
      "content" : ""
    })
  }
  insert_after = dynatrace_log_processing.first-instance.id
}
