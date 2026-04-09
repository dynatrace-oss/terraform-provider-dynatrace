resource "dynatrace_azure_credentials" "Example" {
  active                                     = true
  app_id                                     = "123456789"
  auto_tagging                               = true
  directory_id                               = "123456789"
  key                                        = "123456789"
  label                                      = "#name#"
  monitor_only_tagged_entities               = false
}

resource "dynatrace_azure_service" "ContainerService" {
  name           = "cloud:azure:containerservice:managedcluster"
  credentials_id = dynatrace_azure_credentials.Example.id
  // update => recreate due to set-hash change
  metric {
    name = "kube_pod_status_ready"
    dimensions = ["condition"]
  }
  metric {
    name       = "kube_node_status_condition"
    dimensions = [ "condition", "status", "node" ]
  }
  // one metric removed
}
