resource "dynatrace_process_grouping_rules" "rulesAfter" {
  enabled = true
  custom_technology_name = "custom_technology_name1"
  scope = "environment"
  insert_after = dynatrace_process_grouping_rules.rules.id
  pg_extraction {
    process_group_extraction {
      report = "auto"
      name = "My process group 1"
      process_type = "PROCESS_TYPE_GO"
      detection {
        detection_condition {
          condition = "$contains(TFExecutableSample)"
          case_sensitive = true
          property = "EXE_NAME"
        }
      }
      pg_id_source {
        property        = "COMMAND_LINE_ARGS"
        standalone_rule = false
        type            = "EXISTING"
        advanced_settings {
          ignore_numbers = false
        }
      }
      pgi_id_source {
        property = "AWS_ECS_FAMILY"
        advanced_settings {
          ignore_numbers = false
        }
      }
    }
  }
}


resource "dynatrace_process_grouping_rules" "rules" {
  enabled = true
  custom_technology_name = "custom_technology_name2"
  scope = "environment"
  pg_extraction {
    process_group_extraction {
      report = "auto"
      name = "My process group 1"
      process_type = "PROCESS_TYPE_GO"
      detection {
        detection_condition {
          condition = "$contains(TFExecutableSample)"
          case_sensitive = true
          property = "EXE_NAME"
        }
      }
      pg_id_source {
        property        = "COMMAND_LINE_ARGS"
        standalone_rule = false
        type            = "EXISTING"
        advanced_settings {
          ignore_numbers = false
        }
      }
      pgi_id_source {
        property = "AWS_ECS_FAMILY"
        advanced_settings {
          ignore_numbers = false
        }
      }
    }
  }
}
