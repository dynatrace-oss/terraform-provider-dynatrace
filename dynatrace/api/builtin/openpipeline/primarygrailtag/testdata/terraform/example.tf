resource "dynatrace_openpipeline_primary_grail_tag" "#name#" {
  enabled                = true
  keep_fields            = true
  matching_condition     = "true"
  primary_tag_field_name = "#name#"
  rule_name              = "#name#"
  source_fields {
    source_field {
      field_expression = "aws.tags.env"
    }
    source_field {
      field_expression = "k8s.namespace.label.env"
    }
  }
}
