resource "dynatrace_process_group_detection_flags" "#name#" {
  scope                                       = "HOST_GROUP-1234567890000000"
  add_node_js_script_name                     = false
  auto_detect_cassandra_clusters              = true
  auto_detect_spring_boot                     = true
  auto_detect_tibco_container_edition_engines = false
  auto_detect_tibco_engines                   = false
  auto_detect_web_methods_integration_server  = false
  auto_detect_web_sphere_liberty_application  = false
  group_ibmmqby_instance_name                 = false
  identify_jboss_server_by_system_property    = true
  ignore_unique_identifiers                   = true
  short_lived_processes_monitoring            = true
  split_oracle_database_pg                    = false
  split_oracle_listener_pg                    = false
  use_catalina_base                           = false
  use_docker_container_name                   = false
}
