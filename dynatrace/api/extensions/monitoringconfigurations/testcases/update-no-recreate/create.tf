resource "dynatrace_hub_extension_v2_config" "config" {
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
