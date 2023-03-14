resource "dynatrace_metric_metadata" "#name#" {
  display_name = "#name#"
  metric_id    = "metric-func:slo.terraform-test"
  unit         = "percent"
}