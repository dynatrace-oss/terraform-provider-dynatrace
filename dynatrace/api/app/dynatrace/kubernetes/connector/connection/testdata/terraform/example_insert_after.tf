resource "dynatrace_automation_workflow_k8s_connections" "first-instance" {
  name      = "terraform1"
  uid       = "c6a3798e-634e-49fd-a3ca-00e5c16a1bf3"
  namespace = "terraform1"
  token     = "dt0e01.000000000000000000000000.0000000000000000000000000000000000000000000000000000000000000000"
}

resource "dynatrace_automation_workflow_k8s_connections" "second-instance" {
  name      = "terraform2"
  uid       = "c6a3798e-634e-49fd-a3ca-00e5c16a1bf4"
  namespace = "terraform2"
  token     = "dt0e01.000000000000000000000000.0000000000000000000000000000000000000000000000000000000000000000"
  insert_after = "${dynatrace_automation_workflow_k8s_connections.first-instance.id}"
}