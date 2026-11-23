resource "dynatrace_openpipeline_v2_primary_grail_tag" "example" {
  rule_name              = "#name#"
  enabled                = true
  primary_tag_field_name = "#name#"
  matching_condition     = "isNotNull(k8s.namespace.name)"
  keep_fields            = true
  source_fields {
    source_field {
      field_expression = "k8s.namespace.labels.team"
    }
    source_field {
      field_expression = "team"
    }
  }
}
