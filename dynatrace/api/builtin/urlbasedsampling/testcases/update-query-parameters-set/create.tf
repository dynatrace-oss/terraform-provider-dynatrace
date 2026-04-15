resource "dynatrace_url_based_sampling" "sampling" {
  enabled              = true
  factor               = "ReduceCapturingByFactor16"
  http_method_any      = true
  ignore               = false
  path                 = "/examplepath"
  path_comparison_type = "EQUALS"
  query_parameters {
    parameter {
      name               = "QueryName1"
      value              = "QueryValue1"
      value_is_undefined = false
    }
    # to update
    parameter {
      name               = "QueryName2"
      value              = "QueryValue2"
      value_is_undefined = false
    }
  }
}
