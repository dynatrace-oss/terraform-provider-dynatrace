resource "dynatrace_hub_extension_v2_config" "config" {
  name  = "com.dynatrace.extension.wmi.iis"
  scope = "environment"
  value = jsonencode(
    {
      "enabled" : true,
      "description" : "update",
      "version" : "2.0.1", // version and description update
      "featureSets" : [
        "IIS Extended Request Metrics"
      ],
      "vars" : {
        iis_app_pool = "Name != '_Total'"
        iis_site     = "Name != '_Total'"
      },
      "activationContext" : "LOCAL",
      "activationTags" : []
    }
  )
}
