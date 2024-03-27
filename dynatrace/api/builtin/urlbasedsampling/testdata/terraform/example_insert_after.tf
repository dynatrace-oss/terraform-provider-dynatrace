resource "dynatrace_url_based_sampling" "first-instance" {
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

resource "dynatrace_url_based_sampling" "second-instance" {
  enabled              = true
  factor               = "ReduceCapturingByFactor16"
  http_method_any      = true
  ignore               = false
  path                 = "/examplepath-2"
  path_comparison_type = "EQUALS"
  query_parameters {
    parameter {
      name               = "QueryName"
      value              = "QueryValue2"
      value_is_undefined = false
    }
  }
  insert_after = dynatrace_url_based_sampling.first-instance.id
}
