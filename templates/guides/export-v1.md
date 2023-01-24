---
layout: ""
page_title: "Export Utility (Legacy method)"
description: |-
  The export utility queries the Dynatrace Environment specified and fetches all currently supported configuration
---

## Export Utility (Legacy method)

### Command Line Syntax
Invoking the export functionality requires
* The environment variable `DYNATRACE_ENV_URL` as the URL of your Dynatrace environment
* The environment variable `DYNATRACE_API_TOKEN` as the API Token of your Dynatrace environment
* Optionally the environment variable `DYNATRACE_TARGET_FOLDER`. If it's not set, the output folder `./configuration` is assumed

Windows: `terraform-provider-dynatrace.exe exportv1 *[<resourcename>[=<id>]]`

Linux: `./terraform-provider-dynatrace exportv1 *[<resourcename>[=<id>]]`

### Usage Examples
* `./terraform-provider-dynatrace exportv1` downloads all available configuration settings
* `./terraform-provider-dynatrace exportv1 dynatrace_dashboard` downloads all available dashboards
* `./terraform-provider-dynatrace exportv1 dynatrace_dashboard dynatrace_slo` downloads all available dashboards and all available SLOs
* `./terraform-provider-dynatrace exportv1 dynatrace_dashboard=4f5942d4-3450-40a8-818f-c5faeb3563d0` downloads only the dashboard with the id `4f5942d4-3450-40a8-818f-c5faeb3563d0`
* `./terraform-provider-dynatrace exportv1 dynatrace_dashboard=4f5942d4-3450-40a8-818f-c5faeb3563d0 dynatrace_dashboard=9c4b75f1-9a64-4b44-a8e4-149154fd5325` downloads only the dashboards with the ids `4f5942d4-3450-40a8-818f-c5faeb3563d0` and `9c4b75f1-9a64-4b44-a8e4-149154fd5325`
* `./terraform-provider-dynatrace exportv1 dynatrace_slo dynatrace_dashboard=4f5942d4-3450-40a8-818f-c5faeb3563d0 dynatrace_dashboard=9c4b75f1-9a64-4b44-a8e4-149154fd5325` downloads all available SLOs and only the dashboards with the ids `4f5942d4-3450-40a8-818f-c5faeb3563d0` and `9c4b75f1-9a64-4b44-a8e4-149154fd5
