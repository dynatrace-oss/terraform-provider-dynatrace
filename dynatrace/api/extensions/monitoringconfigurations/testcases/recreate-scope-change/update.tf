data "dynatrace_entities" "hosts" {
  type = "HOST"
}

resource "dynatrace_hub_extension_v2_config" "config" {
  name  = "com.dynatrace.extension.jmx-weblogic-cp"
  scope = data.dynatrace_entities.hosts.entities[0].entity_id // update with HOST-... instead of "environment"
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
