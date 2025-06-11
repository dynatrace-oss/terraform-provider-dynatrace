resource "dynatrace_rpc_based_sampling" "first-instance" {
  enabled                               = false
  endpoint_name                         = "#name#-endpoint"
  endpoint_name_comparison_type         = "DOES_NOT_END_WITH"
  ignore                                = true
  remote_operation_name                 = "#name#-operation"
  remote_operation_name_comparison_type = "CONTAINS"
  remote_service_name                   = "#name#-service"
  remote_service_name_comparison_type   = "STARTS_WITH"
  scope                                 = "environment"
  wire_protocol_type                    = "8"
}

resource "dynatrace_rpc_based_sampling" "second-instance" {
  enabled                               = true
  endpoint_name_comparison_type         = "EQUALS"
  factor                                = "13"
  ignore                                = false
  remote_operation_name                 = "TerraformTest2Operation"
  remote_operation_name_comparison_type = "EQUALS"
  remote_service_name_comparison_type   = "EQUALS"
  scope                                 = "environment"
  wire_protocol_type                    = "1"
  insert_after = dynatrace_rpc_based_sampling.first-instance.id
}
