data "dynatrace_hub_extension_v2_configs" "config1" {
  name               = "com.dynatrace.extension.wmi.iis"
  description        = "#name#"
  version            = "1.1.1"
  activation_context = "LOCAL"

  depends_on = [dynatrace_hub_extension_v2_config.config1, dynatrace_hub_extension_v2_config.config2]
}

data "dynatrace_hub_extension_v2_configs" "config2" {
  name               = "com.dynatrace.extension.wmi.iis"
  description        = "#name#-2"
  version            = "2.0.1"
  activation_context = "LOCAL"

  depends_on = [dynatrace_hub_extension_v2_config.config1, dynatrace_hub_extension_v2_config.config2]
}

output "description1" {
  value = jsondecode(data.dynatrace_hub_extension_v2_configs.config1.items[0].value).description
}

output "description2" {
  value = jsondecode(data.dynatrace_hub_extension_v2_configs.config2.items[0].value).description
}

resource "dynatrace_hub_extension_v2_config" "config1" {
  name  = "com.dynatrace.extension.wmi.iis"
  scope = "environment"
  value = jsonencode(
    {
      "enabled" : true,
      "description" : "#name#",
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

resource "dynatrace_hub_extension_v2_config" "config2" {
  name  = "com.dynatrace.extension.wmi.iis"
  scope = "environment"
  value = jsonencode(
    {
      "enabled" : false,
      "description" : "#name#-2",
      "version" : "2.0.1"
      "featureSets" : [
        "IIS Extended Request Metrics"
      ],
      "vars" : {
        "iis_app_pool" : "Name != '_Total'",
        "iis_site" : "Name != '_Total'"
      },
      "activationContext" : "LOCAL",
      "activationTags" : []
    }
  )
}
