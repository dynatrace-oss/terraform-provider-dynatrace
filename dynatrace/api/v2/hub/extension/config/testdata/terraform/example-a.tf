resource "dynatrace_hub_extension_config" "com_dynatrace_extension_jmx-weblogic-cp" {
  name = "com.dynatrace.extension.jmx-weblogic-cp"
    value = jsonencode(
    {
      "activationContext": "LOCAL",
      "activationTags": [],
      "enabled" : true,
      "description" : "jj",
      "version" : "2.0.4",
      "featureSets" : [
        "cache",
        "connections",
        "capacity"
      ]
    }
  )
}
