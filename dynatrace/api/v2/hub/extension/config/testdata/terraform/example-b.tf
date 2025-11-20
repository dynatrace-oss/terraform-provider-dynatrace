resource "dynatrace_hub_extension_config" "com_dynatrace_extension_jmx-weblogic-cp2" {
  name = "com.dynatrace.extension.jmx-weblogic-cp"
  scope = "environment"
    value = jsonencode(
    {
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