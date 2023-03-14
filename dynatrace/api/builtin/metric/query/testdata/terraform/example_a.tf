resource "dynatrace_metric_query" "#name#" {
  metric_id       = "metric-func:slo.terraform-test"
  metric_selector =<<-EOT
    ((100*(builtin:service.requestCount.server:filter(in("dt.entity.service",entitySelector("type(SERVICE),mzId(0000000000000000000),serviceType(WEB_SERVICE,WEB_REQUEST_SERVICE)"))):splitBy())/(builtin:service.requestCount.server:filter(in("dt.entity.service",entitySelector("type(SERVICE),mzId(0000000000000000000),serviceType(WEB_SERVICE,WEB_REQUEST_SERVICE)"))):splitBy())) - (95.0))
  EOT
}
