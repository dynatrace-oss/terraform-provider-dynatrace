resource "dynatrace_openpipeline_v2_primary_grail_tag" "first" {
  rule_name              = "#name#-1"
  enabled                = true
  primary_tag_field_name = "#name#-1"
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


resource "dynatrace_openpipeline_v2_primary_grail_tag" "second" {
  rule_name              = "#name#-2"
  enabled                = true
  primary_tag_field_name = "#name#-2"
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
  insert_after = dynatrace_openpipeline_v2_primary_grail_tag.first.id
}
