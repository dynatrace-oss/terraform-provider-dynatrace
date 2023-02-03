resource "dynatrace_advanced_process_group_detection_rule" "_47d495a8-5577-436d-a3b3-777924c2d103" {
  enabled = true
  group_extraction {
    property        = "COMMAND_LINE_ARGS"
    standalone_rule = false
    delimiter {
      remove_ids = true
    }
  }
  instance_extraction {
    property = "AWS_ECS_FAMILY"
    delimiter {
      remove_ids = true
    }
  }
  process_detection {
    contained_string = "-config"
    property         = "COMMAND_LINE_ARGS"
    restrict_to_process_type = "PROCESS_TYPE_APACHE_HTTPD"
  }
}
