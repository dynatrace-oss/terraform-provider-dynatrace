data "dynatrace_hub_extension_v2_active_version" "active_version" {
  name = "com.dynatrace.extension.wmi.iis"
}

output "active_version" {
  value = data.dynatrace_hub_extension_v2_active_version.active_version.active_version
}
