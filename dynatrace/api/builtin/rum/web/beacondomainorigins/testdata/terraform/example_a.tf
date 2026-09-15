resource "dynatrace_web_app_beacon_origins" "origin" {
  matcher = "CONTAINS"
  pattern = "pattern-#name#"
}