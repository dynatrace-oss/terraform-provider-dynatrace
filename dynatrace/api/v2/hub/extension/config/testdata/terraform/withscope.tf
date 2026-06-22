resource "dynatrace_hub_extension_config" "com_dynatrace_extension_wmi_iis" {
  name  = "com.dynatrace.extension.wmi.iis"
  scope = "environment"
  value = jsonencode(
    {
      "enabled" : true,
      "description" : "my description",
      "version" : "1.1.1",
      "featureSets" : [
        "IIS Extended Request Metrics"
      ],
      "vars" : {},
      "activationContext" : "LOCAL",
      "activationTags" : []
    }
  )
}
