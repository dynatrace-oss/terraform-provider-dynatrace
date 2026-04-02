resource "dynatrace_hub_extension_v2_config" "config" {
  name  = "com.dynatrace.extension.jmx-weblogic-cp"
  scope = "environment"
  value = jsonencode(
    {
      "activationContext" : "LOCAL",
      "activationTags" : [],
      "enabled" : true,
      "description" : "jj",
      "version" : "2.1.1",
      "featureSets" : [
        "cache",
        "connections",
        "capacity"
      ]
    }
  )
}
