resource "dynatrace_hub_extension_v2_config" "config" {
  name  = "com.dynatrace.extension.jmx-weblogic-cp"
  scope = "environment"
  value = jsonencode(
    {
      "activationContext" : "LOCAL",
      "activationTags" : [],
      "enabled" : true,
      "description" : "update",
      "version" : "2.1.1", // version and description update
      "featureSets" : [
        "cache",
        "connections",
        "capacity"
      ]
    }
  )
}
