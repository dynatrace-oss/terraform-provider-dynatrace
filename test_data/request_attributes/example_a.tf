resource "dynatrace_request_attribute" "#name#" {
  skip_personal_data_masking = true
  confidential               = true
  data_sources {
    scope {
      service_technology   = "JAVA"
      tag_of_process_group = "Mail"
    }
    methods {
      capture            = "THIS"
      deep_object_access = ".getClass().getName()"
      method {
        return_type    = "void"
        visibility     = "PUBLIC"
        argument_types = ["java.lang.String[]"]
        class_name     = "idler.Idler"
        method_name    = "main"
        modifiers      = ["STATIC"]
      }
    }
    technology = "JAVA"
    value_processing {
      split_at = ""
      trim     = false
    }
    enabled = true
    source  = "METHOD_PARAM"
  }
  data_type = "STRING"
  metadata {
    configuration_versions = [3]
    cluster_version        = "1.206.95.20201116-094826"
  }
  normalization = "ORIGINAL"
  enabled       = true
  name          = "#name#"
  aggregation   = "FIRST"
}

