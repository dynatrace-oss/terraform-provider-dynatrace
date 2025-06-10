resource "dynatrace_rpc_based_sampling" "#name#" {
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
