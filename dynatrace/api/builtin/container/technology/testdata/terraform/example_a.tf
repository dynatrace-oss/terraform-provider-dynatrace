resource "dynatrace_container_technology" "environment" {
  bosh_process_manager = true
  containerd           = true
  crio                 = true
  docker               = true
  docker_windows       = true
  garden               = true
  scope                = "environment"
  winc                 = true
}