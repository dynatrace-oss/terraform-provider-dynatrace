resource "dynatrace_openpipeline_primary_grail_tag" "first-instance" {
  enabled                = true
  keep_fields            = true
  matching_condition     = "true"
  primary_tag_field_name = "#name#-first"
  rule_name              = "#name#-first"
  source_fields {
    source_field {
      field_expression = "aws.tags.env"
    }
  }
}

resource "dynatrace_openpipeline_primary_grail_tag" "second-instance" {
  enabled                = true
  keep_fields            = false
  matching_condition     = "isNotNull(k8s.namespace.name)"
  primary_tag_field_name = "#name#-second"
  rule_name              = "#name#-second"
  insert_after           = dynatrace_openpipeline_primary_grail_tag.first-instance.id
  source_fields {
    source_field {
      field_expression = "k8s.namespace.label.env"
    }
  }
}
