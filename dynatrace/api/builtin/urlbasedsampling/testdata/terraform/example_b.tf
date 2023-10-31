resource "dynatrace_url_based_sampling" "#name#" {
  enabled              = true
  http_method          = [ "PATCH", "TRACE", "OPTIONS", "CONNECT", "HEAD" ]
  http_method_any      = false
  ignore               = true
  path                 = "/ignore"
  path_comparison_type = "STARTS_WITH"
  scope                = "environment"
  query_parameters {
    parameter {
      name               = "QueryName"
      value_is_undefined = true
    }
  }
}