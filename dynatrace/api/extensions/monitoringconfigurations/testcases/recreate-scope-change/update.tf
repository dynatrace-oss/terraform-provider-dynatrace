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
