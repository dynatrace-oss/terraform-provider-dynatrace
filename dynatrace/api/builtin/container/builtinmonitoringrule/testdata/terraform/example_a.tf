resource "dynatrace_container_builtin_rule" "#name#" {
  ignore_docker_pause_container     = false
  ignore_kubernetes_pause_container = true
  ignore_open_shift_build_pod_name  = false
  ignore_open_shift_sdn_namespace   = true
}
