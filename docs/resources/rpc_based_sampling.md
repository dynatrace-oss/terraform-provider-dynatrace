---
layout: ""
page_title:  Resource - terraform-provider-dynatrace"
subcategory: "Service Monitoring"
description: |-
  The resource `dynatrace_rpc_based_sampling` covers configuration for trace sampling for RPC requests
---

# dynatrace_rpc_based_sampling (Resource)

-> This resource requires the API token scopes **Read settings** (`settings.read`) and **Write settings** (`settings.write`)

## Dynatrace Documentation

- Trace sampling - https://docs.dynatrace.com/docs/shortlink/url-sampling

- Settings API - https://www.dynatrace.com/support/help/dynatrace-api/environment-api/settings (schemaId: `builtin:rpc-based-sampling`)

## Resource Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)
- `endpoint_name_comparison_type` (String) Possible Values: `CONTAINS`, `DOES_NOT_CONTAIN`, `DOES_NOT_END_WITH`, `DOES_NOT_EQUAL`, `DOES_NOT_START_WITH`, `ENDS_WITH`, `EQUALS`, `STARTS_WITH`
- `ignore` (Boolean) No Traces will be captured for matching RPC requests. This applies always, even if Adaptive Traffic Management is inactive.
- `remote_operation_name_comparison_type` (String) Possible Values: `CONTAINS`, `DOES_NOT_CONTAIN`, `DOES_NOT_END_WITH`, `DOES_NOT_EQUAL`, `DOES_NOT_START_WITH`, `ENDS_WITH`, `EQUALS`, `STARTS_WITH`
- `remote_service_name_comparison_type` (String) Possible Values: `CONTAINS`, `DOES_NOT_CONTAIN`, `DOES_NOT_END_WITH`, `DOES_NOT_EQUAL`, `DOES_NOT_START_WITH`, `ENDS_WITH`, `EQUALS`, `STARTS_WITH`
- `wire_protocol_type` (String) Possible Values: `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10` 

  - `1`: ADK

  - `2`: DOTNET_REMOTING

  - `3`: DOTNET_REMOTING_TCP

  - `4`: DOTNET_REMOTING_HTTP

  - `5`: DOTNET_REMOTING_XMLRPC

  - `6`: GRPC

  - `7`: GRPC_BIDI

  - `8`: GRPC_UNARY

  - `9`: GRPC_SERVERSTREAM

  - `10`: GRPC_CLIENTSTREAM

### Optional

- `endpoint_name` (String) Specify the RPC endpoint name. If the endpoint name is empty, either remote operation name or remote service name must be specified that can be used for RPC matching.
- `factor` (String) Possible Values: `0`, `1`, `2`, `3`, `4`, `5`, `6`, `8`, `9`, `10`, `11`, `12`, `13`, `14` 

  - `0`: Increase capturing 128 times

  - `1`: Increase capturing 64 times

  - `2`: Increase capturing 32 times

  - `3`: Increase capturing 16 times

  - `4`: Increase capturing 8 times

  - `5`: Increase capturing 4 times

  - `6`: Increase capturing 2 times

  - `8`: Reduce capturing by factor 2

  - `9`: Reduce capturing by factor 4

  - `10`: Reduce capturing by factor 8

  - `11`: Reduce capturing by factor 16

  - `12`: Reduce capturing by factor 32

  - `13`: Reduce capturing by factor 64

  - `14`: Reduce capturing by factor 128
- `insert_after` (String) Because this resource allows for ordering you may specify the ID of the resource instance that comes before this instance regarding order. If not specified when creating the setting will be added to the end of the list. If not specified during update the order will remain untouched
- `remote_operation_name` (String) Specify the RPC operation name. If the remote operation name is empty, either remote service name or endpoint name must be specified that can be used for RPC matching.
- `remote_service_name` (String) Specify the RPC remote service name. If the remote service name is empty, either remote operation name or endpoint name must be specified that can be used for RPC matching.
- `scope` (String) The scope of this setting (PROCESS_GROUP_INSTANCE, PROCESS_GROUP, CLOUD_APPLICATION, CLOUD_APPLICATION_NAMESPACE, KUBERNETES_CLUSTER, HOST_GROUP). Omit this property if you want to cover the whole environment.

### Read-Only

- `id` (String) The ID of this resource.
