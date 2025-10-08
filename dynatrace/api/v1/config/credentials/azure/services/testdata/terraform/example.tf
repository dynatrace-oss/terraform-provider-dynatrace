resource "dynatrace_azure_credentials" "Example" {
  active                                     = true
  app_id                                     = "123456789"
  auto_tagging                               = true
  directory_id                               = "123456789"
  key                                        = "123456789"
  label                                      = "Example"
  monitor_only_tagged_entities               = false
}


resource "dynatrace_azure_service" "ContainerService" {
  name           = "cloud:azure:containerservice:managedcluster"
  credentials_id = dynatrace_azure_credentials.Example.id
  metric {
    name = "kube_pod_status_ready"
    dimensions = []
  }
  metric {
    name       = "kube_node_status_condition"
    dimensions = [ "condition", "status", "node" ]
  }
  metric {
    name       = "kube_pod_status_phase"
    dimensions = [ "phase", "namespace" ]
  }
}