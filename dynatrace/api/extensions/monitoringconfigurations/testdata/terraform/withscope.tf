data "dynatrace_entities" "hosts" {
  type = "HOST"
}

resource "dynatrace_hub_extension_v2_config" "com_dynatrace_extension_jmx-weblogic-cp" {
  name  = "com.dynatrace.extension.jmx-weblogic-cp"
  scope = data.dynatrace_entities.hosts.entities[0].entity_id // or "environment"
  value = jsonencode(
    {
      "activationContext" : "LOCAL",
      "activationTags" : [],
      "enabled" : true,
      "description" : "my description",
      "version" : "2.1.1",
      "featureSets" : [
        "cache",
        "connections",
        "capacity"
      ]
    }
  )
}
