resource "dynatrace_log_processing" "#name#" {
  enabled   = false
  query     = "matchesPhrase(content, \"terratest\")"
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