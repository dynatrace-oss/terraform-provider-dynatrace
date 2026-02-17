locals {
  location_name = "Location"
}

data "dynatrace_synthetic_location" "location" {
  name = local.location_name
}

import {
  to = dynatrace_synthetic_location.location
  for_each = try(data.dynatrace_synthetic_location.location.id, null) == null ? [] : [data.dynatrace_synthetic_location.location.id]
  id = each.value
}

resource "dynatrace_synthetic_location" "location" {
  name                                  = local.location_name
  city                                  = "San Francisco de Asis"
  country_code                          = "VE"
  region_code                           = "04"
  deployment_type                       = "STANDARD"
  latitude                              = 10.0756
  location_node_outage_delay_in_minutes = 3
  longitude                             = -67.5442

  lifecycle {
    prevent_destroy = true
  }
}
