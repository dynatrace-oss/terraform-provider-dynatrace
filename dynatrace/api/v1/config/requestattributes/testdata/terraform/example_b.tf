resource "dynatrace_request_attribute" "attribute" {
  name         = "#name#"
  enabled      = true
  aggregation  = "FIRST"
  confidential = false

  data_type                  = "INTEGER"
  normalization              = "ORIGINAL"
  skip_personal_data_masking = false
  data_sources {
    enabled = true
    source = "SERVER_VARIABLE"
    server_variable_technology = "ASP_NET"
    parameter_name             = "param"
  }
}
