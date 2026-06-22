data "dynatrace_hub_extension_v2_latest_version" "latest_version" {
  name = "com.dynatrace.extension.wmi.iis"
}

output "latest_version" {
  value = data.dynatrace_hub_extension_v2_latest_version.latest_version.latest_version
}
