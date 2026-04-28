data "dynatrace_hub_extension_v2_latest_version" "latest_version" {
  name = "com.dynatrace.extension.jmx-weblogic-cp"
}

output "latest_version" {
  value = data.dynatrace_hub_extension_v2_latest_version.latest_version.latest_version
}
