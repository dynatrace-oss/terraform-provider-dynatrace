terraform {
  required_providers {
    dynatrace = {
      source  = "dynatrace-oss/dynatrace"
      version = ">=1.90.0"
    }
  }
}

provider "dynatrace" {
  dt_env_url       = var.DYNATRACE_ENV_URL
  dt_api_token     = var.DYNATRACE_API_TOKEN
}
