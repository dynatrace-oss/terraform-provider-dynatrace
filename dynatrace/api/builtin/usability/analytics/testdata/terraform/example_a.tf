resource "dynatrace_usability_analytics" "test" {
    detect_rage_clicks = true
}

resource "dynatrace_usability_analytics" "for_app" {
    application_id = "APPLICATION-EA7C4B59F27D43EB"
    detect_rage_clicks = false
}