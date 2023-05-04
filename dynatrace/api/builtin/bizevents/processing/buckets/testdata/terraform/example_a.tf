resource "dynatrace_business_events_buckets" "#name#" {
  enabled     = true
  bucket_name = "default_bizevents_medium"
  matcher     = "matchesValue(event.type, \"com.easytrade.buy-assets\")"
  rule_name   = "#name#"
}