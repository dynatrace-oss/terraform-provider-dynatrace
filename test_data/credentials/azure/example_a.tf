resource "dynatrace_azure_credentials" "#name#" {
  active                       = false
  app_id                       = "ABCDE"
  auto_tagging                 = true
  directory_id                 = "ABCDE"
  label                        = "#name#"
  key                          = "aaaa"
  monitor_only_tagged_entities = true

  monitor_only_tag_pairs {
    name  = "string"
    value = "string"
  }
}
