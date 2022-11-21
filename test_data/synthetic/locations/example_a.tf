resource "dynatrace_synthetic_location" "#name#" {
  name = "#name#" 
  auto_update_chromium = true 
  availability_location_outage = true 
  availability_node_outage = true 
  availability_notifications_enabled = true 
  city = "San Francisco de Asis" 
  country_code = "VE" 
  deployment_type = "STANDARD" 
  latitude = 10.0758 
  location_node_outage_delay_in_minutes = 3 
  longitude = -67.5442 
  region_code = "04" 
}
