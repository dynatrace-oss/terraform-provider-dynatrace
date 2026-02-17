variable "DYNATRACE_API_TOKEN" {
  type        = string
  description = "Dynatrace API Token with appropriate permissions."
  sensitive   = true
}

variable "DYNATRACE_ENV_URL" {
  type        = string
  description = "Dynatrace Environment URL."
}
