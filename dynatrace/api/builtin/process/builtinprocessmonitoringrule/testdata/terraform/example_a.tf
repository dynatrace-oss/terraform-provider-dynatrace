resource "dynatrace_builtin_process_monitoring" "#name#" {
  host_group_id                         = "environment"
  aspnetcore                            = false
  cf_appsmanagerjs                      = false
  container                             = false
  docker_pauseamd64                     = false
  exe_bbs                               = false
  exe_caddy                             = false
  exe_schedular                         = false
  exe_silkdaemon                        = false
  go_static                             = false
  node_nodegyp                          = false
  cmd_foreverbinmonitor                 = false
}
