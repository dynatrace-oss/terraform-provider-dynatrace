terraform {
  required_providers {
    dynatrace = {
      versions = ["0.1.0"]
      source = "dynatrace.com/com/dynatrace"
    }
  }
}

provider "dynatrace" {
    dt_env_url    = "https://########.live.dynatrace.com"
    dt_api_token  = "################"
}

resource "dynatrace_custom_service" "BODYHandler" {
    name = "BODYHandler"
    technology = "java"
    enabled = true
    queue_entry_point = false
    process_groups = []
    rule {
      enabled = true
      class {
        name = "dtcookie.vertx.MainVerticle$BodyHandler"
        match = "EQUALS"
      }
      // annotations = [] 
      method {
        name = "handle"
        returns = "void"
        arguments = ["io.vertx.core.buffer.Buffer"]
      }
    }
}

/*
resource "dynatrace_sample" "sample-a" {
    legend_shown = true
    result_metadata {
      key = "this-idioticly-long-string-our-dashboards-contain-as-map-keys-1"
      last_modified = 1234
      custom_color = "#000000"
    }
    result_metadata {
      key = "this-idioticly-long-string-our-dashboards-contain-as-map-keys-2"
      last_modified = 4321
      custom_color = "#FFFFFF"
    }
}
*/