# ID vu9U3hXa3q0AAAABABlidWlsdGluOm9uZWFnZW50LmZlYXR1cmVzAAZ0ZW5hbnQABnRlbmFudAAkMWQzYjY4ODMtOWViZi0zMDljLTg1YjktNjg4OTcxYzE3NDM1vu9U3hXa3q0
resource "dynatrace_oneagent_features" "SENSOR_DOTNET_ASPNET" {
  enabled         = true
  instrumentation = true
  key             = "SENSOR_DOTNET_ASPNET"
  scope           = "environment"
}
