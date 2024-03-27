resource "dynatrace_process_group_detection" "first-instance" {
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

resource "dynatrace_process_group_detection" "second-instance" {
  enabled = true
  group_extraction {
    property        = "COMMAND_LINE_ARGS"
    standalone_rule = false
    delimiter {
      remove_ids = false
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
  insert_after = dynatrace_process_group_detection.first-instance.id
}
