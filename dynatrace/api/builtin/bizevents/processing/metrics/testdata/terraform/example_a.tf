resource "dynatrace_business_events_metrics" "#name#" {
  enabled           = true
  key               = "bizevents.easyTrade.#name#"
  matcher           = "matchesValue(event.type, \"com.easytrade.buy-assets\")"
  measure           = "ATTRIBUTE"
  measure_attribute = "trading_volume"
}
