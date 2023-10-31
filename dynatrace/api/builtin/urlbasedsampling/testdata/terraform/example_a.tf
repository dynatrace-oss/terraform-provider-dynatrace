resource "dynatrace_url_based_sampling" "#name#" {
  enabled              = true
  factor               = "ReduceCapturingByFactor16"
  http_method_any      = true
  ignore               = false
  path                 = "/examplepath"
  path_comparison_type = "EQUALS"
  query_parameters {
    parameter {
      name               = "QueryName"
      value              = "QueryValue"
      value_is_undefined = false
    }
  }
}
