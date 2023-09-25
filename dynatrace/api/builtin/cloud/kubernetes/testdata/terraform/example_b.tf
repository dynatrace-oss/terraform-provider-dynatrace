resource "dynatrace_kubernetes" "#name#" {
  enabled                            = true
  cluster_id                         = "#name#"
  cluster_id_enabled                 = true
  label                              = "#name#"
}