data "dynatrace_dql" "records" {
  default_sampling_ratio = 1.0
  default_scan_limit_gbytes = -1
  timezone = "UTC"
  fetch_timeout_seconds = 60
  locale = "en"
  max_result_records = 10
  max_result_bytes = 10485760
  query = <<-EOT
    data record(id = 1) | fields id
  EOT
}

output "records" {
  value = data.dynatrace_dql.records.records
}
