resource "dynatrace_container_builtin_rule" "rules" {
  ignore_docker_pause_container                       = false
  ignore_kubernetes_pause_container                   = true
  ignore_open_shift_build_pod_name                    = false
  ignore_open_shift_sdn_namespace                     = true
  ignore_open_shift_etcd_namespace                    = false
  ignore_open_shift_ingress_canary_namespace          = false
  ignore_open_shift_kube_apiserver_namespace          = false
  ignore_open_shift_machine_config_operator_namespace = false
  ignore_open_shift_monitoring_namespace              = false
  ignore_open_shift_ovn_kubernetes_namespace          = false
}
