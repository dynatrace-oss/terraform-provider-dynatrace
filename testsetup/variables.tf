variable "DYNATRACE_API_TOKEN" {
  type        = string
  description = "Dynatrace API Token with appropriate permissions."
  sensitive   = true
  default     = null
}

variable "DYNATRACE_ENV_URL" {
  type        = string
  description = "Dynatrace Environment URL."
}
