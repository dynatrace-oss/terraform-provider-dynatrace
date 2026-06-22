data "dynatrace_entities" "hosts" {
  type = "HOST"
}

resource "dynatrace_hub_extension_v2_config" "config" {
  name  = "com.dynatrace.extension.wmi.iis"
  scope = data.dynatrace_entities.hosts.entities[0].entity_id // update with HOST-... instead of "environment"
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
