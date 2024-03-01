---
layout: ""
page_title: dynatrace_http_monitor Resource - terraform-provider-dynatrace"
subcategory: "HTTP Monitors"
description: |-
  The resource `dynatrace_http_monitor` covers configuration for HTTP monitors
---

# dynatrace_http_monitor (Resource)

-> This resource requires the API token scope **Create and read synthetic monitors, locations, and nodes** (`ExternalSyntheticIntegration`)

## Dynatrace Documentation

- Synthetic Monitoring - HTTP monitors - https://www.dynatrace.com/support/help/shortlink/synthetic-hub#http-monitors

- Synthetic Monitors API - https://www.dynatrace.com/support/help/dynatrace-api/environment-api/synthetic/synthetic-monitors

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_http_monitor` downloads all existing HTTP monitor configuration

The full documentation of the export feature is available [here](https://registry.terraform.io/providers/dynatrace-oss/dynatrace/latest/docs/guides/export-v2).

## Resource Example Usage

```terraform
resource "dynatrace_http_monitor" "#name#" {
  name      = "#name#"
  frequency = 1
  locations = ["GEOLOCATION-F3E06A526BE3B4C4"]
  anomaly_detection {
    loading_time_thresholds {
    }
    outage_handling {
      global_outage = true
      global_outage_policy {
        consecutive_runs = 1
      }
    }
  }
  script {
    request {
      description     = "getOffice365ActiveUserCounts"
      method          = "GET"
      post_processing = <<-EOT
        if (response.getStatusCode() != 200) {
            api.fail("HTTP error: " + response.getStatusCode());
        }
        var register = function(obj, key, value) {
          if (key in obj) {
              return;
          }
          value = value.trim();
          if (value.length === 0) {
              return;
          }
          var iValue = parseInt(value);
          if (isNaN(iValue)) {
              return;
          }
          obj[key] = iValue;
        };
        var lines = response.getResponseBody().trim().split("\n");
        var idx = 0;
        
        var counts = {};
        for (idx = lines.length - 1; idx >= 0; idx--) {
            var line = lines[idx].trim();
            if (line.length === 0) {
                continue;
            }
            if (line.startsWith("Report Refresh Date")) {
                continue;
            }
            var values = line.split(",");
            register(counts, "office365", values[1]);
            register(counts, "exchange", values[2]);
            register(counts, "onedrive", values[3]);
            register(counts, "sharepoint", values[4]);
            register(counts, "skype", values[5]);
            register(counts, "yammer", values[6]);
            register(counts, "teams", values[7]);
        }
        
        if ("office365" in counts) api.setValue("office365", counts.office365);
        if ("exchange" in counts) api.setValue("exchange", counts.exchange);
        if ("onedrive" in counts) api.setValue("onedrive", counts.onedrive);
        if ("sharepoint" in counts) api.setValue("sharepoint", counts.sharepoint);
        if ("skype" in counts) api.setValue("skype", counts.skype);
        if ("yammer" in counts) api.setValue("yammer", counts.yammer);
        if ("teams" in counts) api.setValue("teams", counts.teams);
      EOT
      url             = "https://graph.microsoft.com/v1.0/reports/getOffice365ActiveUserCounts(period='D7')"
      configuration {
        accept_any_certificate = true
        follow_redirects       = true
        headers {
          header {
            name  = "Authorization"
            value = "Bearer {CREDENTIALS_VAULT-93C49382ACCA047B|token}"
          }
          header {
            name  = "Accept"
            value = "application/json"
          }
        }
      }
    }
    request {
      description     = "getMailboxUsageQuotaStatusMailboxCounts"
      method          = "GET"
      post_processing = <<-EOT
        if (response.getStatusCode() != 200) {
            api.fail("HTTP error: " + response.getStatusCode());
        }
        var register = function(obj, key, value) {
          if (key in obj) {
              return;
          }
          value = value.trim();
          if (value.length === 0) {
              return;
          }
          var iValue = parseInt(value);
          if (isNaN(iValue)) {
              return;
          }
          obj[key] = iValue;
        };
        var lines = response.getResponseBody().trim().split("\n");
        var idx = 0;
        
        var counts = {};
        for (idx = lines.length - 1; idx >= 0; idx--) {
            var line = lines[idx].trim();
            if (line.length === 0) {
                continue;
            }
            if (line.startsWith("Report Refresh Date")) {
                continue;
            }
            var values = line.split(",");
            register(counts, "under_limit", values[1]);
            register(counts, "warning_issued", values[2]);
            register(counts, "send_prohibited", values[3]);
            register(counts, "send_receive_prohibited", values[4]);
            register(counts, "indeterminate", values[5]);
        }
        
        api.setValue("under_limit", ("under_limit" in counts) ? counts.under_limit : 0);
        api.setValue("warning_issued", ("warning_issued" in counts) ? counts.warning_issued : 0);
        api.setValue("send_prohibited", ("send_prohibited" in counts) ? counts.send_prohibited : 0);
        api.setValue("send_receive_prohibited", ("send_receive_prohibited" in counts) ? counts.send_receive_prohibited : 0);
        api.setValue("indeterminate", ("indeterminate" in counts) ? counts.send_receive_prohibited : 0);
      EOT
      url             = "https://graph.microsoft.com/v1.0/reports/getMailboxUsageQuotaStatusMailboxCounts(period='D7')"
      configuration {
        accept_any_certificate = true
        follow_redirects       = true
        headers {
          header {
            name  = "Authorization"
            value = "Bearer {CREDENTIALS_VAULT-93C49382ACCA047B|token}"
          }
        }
      }
      validation {
        rule {
          type  = "httpStatusesList"
          value = "\u003e=400"
        }
      }
    }
    request {
      description     = "getMailboxUsageStorage"
      method          = "GET"
      post_processing = <<-EOT
        if (response.getStatusCode() != 200) {
            api.fail("HTTP error: " + response.getStatusCode());
        }
        var register = function(obj, key, value) {
          if (key in obj) {
              return;
          }
          value = value.trim();
          if (value.length === 0) {
              return;
          }
          var iValue = parseInt(value);
          if (isNaN(iValue)) {
              return;
          }
          obj[key] = iValue;
        };
        var lines = response.getResponseBody().trim().split("\n");
        var idx = 0;
        
        var counts = {};
        for (idx = lines.length - 1; idx >= 0; idx--) {
            var line = lines[idx].trim();
            if (line.length === 0) {
                continue;
            }
            if (line.startsWith("Report Refresh Date")) {
                continue;
            }
            var values = line.split(",");
            register(counts, "storage_used", values[1]);
        }
        
        api.setValue("storage_used", ("storage_used" in counts) ? counts.storage_used / 1024 / 1024 / 1024 : 0);
        api.setValue("storage_used_mailbox", ("storage_used" in counts) ? counts.storage_used / 1024 / 1024 / 1024 : 0);
      EOT
      url             = "https://graph.microsoft.com/v1.0/reports/getMailboxUsageStorage(period='D7')"
      configuration {
        accept_any_certificate = true
        follow_redirects       = true
        headers {
          header {
            name  = "Authorization"
            value = "Bearer {CREDENTIALS_VAULT-93C49382ACCA047B|token}"
          }
        }
      }
      validation {
        rule {
          type  = "httpStatusesList"
          value = "\u003e=400"
        }
      }
    }
    request {
      description     = "getSharePointSiteUsageStorage"
      method          = "GET"
      post_processing = <<-EOT
        if (response.getStatusCode() != 200) {
            api.fail("HTTP error: " + response.getStatusCode());
        }
        var register = function(obj, key, value) {
          if (key in obj) {
              return;
          }
          value = value.trim();
          if (value.length === 0) {
              return;
          }
          var iValue = parseInt(value);
          if (isNaN(iValue)) {
              return;
          }
          obj[key] = iValue;
        };
        var lines = response.getResponseBody().trim().split("\n");
        var idx = 0;
        
        var counts = {};
        for (idx = lines.length - 1; idx >= 0; idx--) {
            var line = lines[idx].trim();
            if (line.length === 0) {
                continue;
            }
            if (line.startsWith("Report Refresh Date")) {
                continue;
            }
            if (!line.includes(",All,")) {
                continue;
            }
            var values = line.split(",");
            register(counts, "storage_used", values[2]);
        }
        
        api.setValue("storage_used_sharepoint", ("storage_used" in counts) ? counts.storage_used / 1024 / 1024 / 1024 : 0);
      EOT
      url             = "https://graph.microsoft.com/v1.0/reports/getSharePointSiteUsageStorage(period='D7')"
      configuration {
        accept_any_certificate = true
        follow_redirects       = true
        headers {
          header {
            name  = "Authorization"
            value = "Bearer {CREDENTIALS_VAULT-93C49382ACCA047B|token}"
          }
        }
      }
      validation {
        rule {
          type  = "httpStatusesList"
          value = "\u003e=400"
        }
      }
    }
    request {
      description     = "getOneDriveUsageStorage"
      method          = "GET"
      post_processing = <<-EOT
        if (response.getStatusCode() != 200) {
            api.fail("HTTP error: " + response.getStatusCode());
        }
        var register = function(obj, key, value) {
          if (key in obj) {
              return;
          }
          value = value.trim();
          if (value.length === 0) {
              return;
          }
          var iValue = parseInt(value);
          if (isNaN(iValue)) {
              return;
          }
          obj[key] = iValue;
        };
        var lines = response.getResponseBody().trim().split("\n");
        var idx = 0;
        
        var counts = {};
        for (idx = lines.length - 1; idx >= 0; idx--) {
            var line = lines[idx].trim();
            if (line.length === 0) {
                continue;
            }
            if (line.startsWith("Report Refresh Date")) {
                continue;
            }
            if (!line.includes(",All,")) {
                continue;
            }
            var values = line.split(",");
            register(counts, "storage_used", values[2]);
        }
        
        api.setValue("storage_used_onedrive", ("storage_used" in counts) ? counts.storage_used / 1024 / 1024 / 1024 : 0);
      EOT
      url             = "https://graph.microsoft.com/v1.0/reports/getOneDriveUsageStorage(period='D7')"
      configuration {
        accept_any_certificate = true
        follow_redirects       = true
        headers {
          header {
            name  = "Authorization"
            value = "Bearer {CREDENTIALS_VAULT-93C49382ACCA047B|token}"
          }
        }
      }
      validation {
        rule {
          type  = "httpStatusesList"
          value = "\u003e=400"
        }
      }
    }
    request {
      description     = "ServiceComms/CurrentStatus"
      method          = "GET"
      post_processing = <<-EOT
        var healthyStates = [
            "PostIncidentReviewPublished",
            "ServiceRestored",
            "ServiceOperational",
            "FalsePositive"
        ];
        /* Work load status per https://docs.microsoft.com/en-us/office/office-365-management-api/office-365-service-communications-api-reference#status-definitions
        Investigating
        ServiceDegradation
        ServiceInterruption
        RestoringService
        ExtendedRecovery
        InvestigationSuspended
        ServiceRestored
        FalsePositive
        PostIncidentReportPublished
        ServiceOperational
        */
        
        json = JSON.parse(response.getResponseBody());
        
        var payload = "office365.service.status.queried 1";
        json.value.forEach(element => {
            payload = payload + "\noffice365.service.status,workload=" + element.Workload + ",status=" + element.Status + ",healthy=" + (healthyStates.indexOf(element.Status) >= 0) + " 1";
        });
        api.setValue("service_status", payload);
      EOT
      url             = "https://manage.office.com/api/v1.0/{CREDENTIALS_VAULT-1A8E917381883F54|token}/ServiceComms/CurrentStatus"
      configuration {
        accept_any_certificate = true
        follow_redirects       = true
        headers {
          header {
            name  = "Authorization"
            value = "Bearer {CREDENTIALS_VAULT-CE4EA27BA94C9061|token}"
          }
          header {
            name  = "Accept"
            value = "application/json"
          }
        }
      }
      validation {
        rule {
          type  = "httpStatusesList"
          value = "\u003e=400"
        }
      }
    }
    request {
      description = "api/v2/metrics/ingest"
      body        = <<-EOT
        office365.user.count,product=sharepoint {sharepoint}
        office365.user.count,product=onedrive {onedrive}
        office365.user.count,product=yammer {yammer}
        office365.user.count,product=office365 {office365}
        office365.user.count,product=skype {skype}
        office365.user.count,product=exchange {exchange}
        office365.user.count,product=teams {teams}
        office365.mailbox.quota.count,category=under_limit {under_limit}
        office365.mailbox.quota.count,category=warning_issued {warning_issued}
        office365.mailbox.quota.count,category=send_prohibited {send_prohibited}
        office365.mailbox.quota.count,category=send_receive_prohibited {send_receive_prohibited}
        office365.mailbox.quota.count,category=indeterminate {indeterminate}
        office365.storage.used.bytes,site=outlook {storage_used_mailbox}
        office365.storage.used.bytes,site=sharepoint {storage_used_sharepoint}
        office365.storage.used.bytes,site=onedrive {storage_used_onedrive}
        {service_status}
      EOT
      method      = "POST"
      url         = "https://siz65484.live.dynatrace.com/api/v2/metrics/ingest"
      configuration {
        accept_any_certificate = true
        follow_redirects       = true
        headers {
          header {
            name  = "Content-Type"
            value = "text/plain"
          }
          header {
            name  = "Authorization"
            value = "Api-Token {CREDENTIALS_VAULT-55F1E51535993619|token}"
          }
        }
      }
      validation {
        rule {
          type  = "httpStatusesList"
          value = "\u003e=400"
        }
      }
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `frequency` (Number) The frequency of the monitor, in minutes.

You can use one of the following values: `5`, `10`, `15`, `30`, and `60`.
- `name` (String) The name of the monitor.

### Optional

- `anomaly_detection` (Block List) The anomaly detection configuration. (see [below for nested schema](#nestedblock--anomaly_detection))
- `enabled` (Boolean) The monitor is enabled (`true`) or disabled (`false`).
- `locations` (Set of String) A list of locations from which the monitor is executed.

To specify a location, use its entity ID.
- `manually_assigned_apps` (Set of String) A set of manually assigned applications.
- `no_script` (Boolean) No script block - handle requests via `dynatrace_http_monitor_script` resource
- `script` (Block List, Max: 1) The HTTP Script (see [below for nested schema](#nestedblock--script))
- `tags` (Block List) A set of tags assigned to the monitor.

You can specify only the value of the tag here and the `CONTEXTLESS` context and source 'USER' will be added automatically. (see [below for nested schema](#nestedblock--tags))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--anomaly_detection"></a>
### Nested Schema for `anomaly_detection`

Optional:

- `loading_time_thresholds` (Block List) Thresholds for loading times (see [below for nested schema](#nestedblock--anomaly_detection--loading_time_thresholds))
- `outage_handling` (Block List) Outage handling configuration (see [below for nested schema](#nestedblock--anomaly_detection--outage_handling))

<a id="nestedblock--anomaly_detection--loading_time_thresholds"></a>
### Nested Schema for `anomaly_detection.loading_time_thresholds`

Optional:

- `enabled` (Boolean) Performance threshold is enabled (`true`) or disabled (`false`)
- `thresholds` (Block List) The list of performance threshold rules (see [below for nested schema](#nestedblock--anomaly_detection--loading_time_thresholds--thresholds))

<a id="nestedblock--anomaly_detection--loading_time_thresholds--thresholds"></a>
### Nested Schema for `anomaly_detection.loading_time_thresholds.thresholds`

Required:

- `threshold` (Block List, Min: 1) The list of performance threshold rules (see [below for nested schema](#nestedblock--anomaly_detection--loading_time_thresholds--thresholds--threshold))

<a id="nestedblock--anomaly_detection--loading_time_thresholds--thresholds--threshold"></a>
### Nested Schema for `anomaly_detection.loading_time_thresholds.thresholds.threshold`

Required:

- `value_ms` (Number) Notify if monitor takes longer than *X* milliseconds to load

Optional:

- `event_index` (Number) Specify the event to which an ACTION threshold applies
- `request_index` (Number) Specify the request to which an ACTION threshold applies
- `type` (String) The type of the threshold: `TOTAL` (total loading time) or `ACTION` (action loading time)




<a id="nestedblock--anomaly_detection--outage_handling"></a>
### Nested Schema for `anomaly_detection.outage_handling`

Optional:

- `global_outage` (Boolean) (Field has overlap with `dynatrace_browser_monitor_outage` and `dynatrace_http_monitor_outage`) When enabled (`true`), generate a problem and send an alert when the monitor is unavailable at all configured locations
- `global_outage_policy` (Block List) (Field has overlap with `dynatrace_browser_monitor_outage` and `dynatrace_http_monitor_outage`) Global outage handling configuration. 

 Alert if **consecutiveRuns** times consecutively (see [below for nested schema](#nestedblock--anomaly_detection--outage_handling--global_outage_policy))
- `local_outage` (Boolean) (Field has overlap with `dynatrace_browser_monitor_outage` and `dynatrace_http_monitor_outage`) When enabled (`true`), generate a problem and send an alert when the monitor is unavailable for one or more consecutive runs at any location
- `local_outage_policy` (Block List) (Field has overlap with `dynatrace_browser_monitor_outage` and `dynatrace_http_monitor_outage`) Local outage handling configuration. 

 Alert if **affectedLocations** of locations are unable to access the web application **consecutiveRuns** times consecutively (see [below for nested schema](#nestedblock--anomaly_detection--outage_handling--local_outage_policy))
- `retry_on_error` (Boolean) (Field has overlap with `dynatrace_browser_monitor_outage` and `dynatrace_http_monitor_outage`) Schedule retry if browser monitor execution results in a fail. For HTTP monitors this property is ignored

<a id="nestedblock--anomaly_detection--outage_handling--global_outage_policy"></a>
### Nested Schema for `anomaly_detection.outage_handling.global_outage_policy`

Required:

- `consecutive_runs` (Number) The number of consecutive fails to trigger an alert


<a id="nestedblock--anomaly_detection--outage_handling--local_outage_policy"></a>
### Nested Schema for `anomaly_detection.outage_handling.local_outage_policy`

Required:

- `affected_locations` (Number) The number of affected locations to trigger an alert
- `consecutive_runs` (Number) The number of consecutive fails to trigger an alert




<a id="nestedblock--script"></a>
### Nested Schema for `script`

Required:

- `request` (Block List, Min: 1) A HTTP request to be performed by the monitor. (see [below for nested schema](#nestedblock--script--request))

<a id="nestedblock--script--request"></a>
### Nested Schema for `script.request`

Required:

- `method` (String) The HTTP method of the request.
- `url` (String) The URL to check.

Optional:

- `authentication` (Block List, Max: 1) Authentication options for this request (see [below for nested schema](#nestedblock--script--request--authentication))
- `body` (String) The body of the HTTP request.
- `configuration` (Block List, Max: 1) The setup of the monitor (see [below for nested schema](#nestedblock--script--request--configuration))
- `description` (String) A short description of the event to appear in the web UI.
- `post_processing` (String) Javascript code to execute after sending the request.
- `pre_processing` (String) Javascript code to execute before sending the request.
- `request_timeout` (Number) Adapt request timeout option - the maximum time this request is allowed to consume. Keep in mind the maximum timeout of the complete monitor is 60 seconds
- `validation` (Block List, Max: 1) Validation helps you verify that your HTTP monitor loads the expected content (see [below for nested schema](#nestedblock--script--request--validation))

<a id="nestedblock--script--request--authentication"></a>
### Nested Schema for `script.request.authentication`

Required:

- `credentials` (String) The ID of the credentials within the Dynatrace Credentials Vault.
- `type` (String) The type of authentication. Possible values are `BASIC_AUTHENTICATION`, `NTLM` and `KERBEROS`.

Optional:

- `kdc_ip` (String) The KDC IP. Valid and required only if the type of authentication is `KERBEROS`.
- `realm_name` (String) The Realm Name. Valid and required only if the type of authentication is `KERBEROS`.


<a id="nestedblock--script--request--configuration"></a>
### Nested Schema for `script.request.configuration`

Optional:

- `accept_any_certificate` (Boolean) If set to `false`, then the monitor fails with invalid SSL certificates.

If not set, the `false` option is used
- `client_certificate` (String, Sensitive) The client certificate, if applicable - eg. CREDENTIALS_VAULT-XXXXXXXXXXXXXXXX
- `follow_redirects` (Boolean) If set to `false`, redirects are reported as successful requests with response code 3xx.

If not set, the `false` option is used.
- `headers` (Block List, Max: 1) The setup of the monitor (see [below for nested schema](#nestedblock--script--request--configuration--headers))
- `sensitive_data` (Boolean) Option not to store and display request and response bodies and header values in execution details, `true` or `false`. If not set, `false`.
- `user_agent` (String) The User agent of the request

<a id="nestedblock--script--request--configuration--headers"></a>
### Nested Schema for `script.request.configuration.headers`

Required:

- `header` (Block Set, Min: 1) contains an HTTP header of the request (see [below for nested schema](#nestedblock--script--request--configuration--headers--header))

<a id="nestedblock--script--request--configuration--headers--header"></a>
### Nested Schema for `script.request.configuration.headers.header`

Required:

- `name` (String) The key of the header
- `value` (String) The value of the header




<a id="nestedblock--script--request--validation"></a>
### Nested Schema for `script.request.validation`

Required:

- `rule` (Block List, Min: 1) A list of validation rules (see [below for nested schema](#nestedblock--script--request--validation--rule))

<a id="nestedblock--script--request--validation--rule"></a>
### Nested Schema for `script.request.validation.rule`

Required:

- `type` (String) The type of the rule. Possible values are `patternConstraint`, `regexConstraint`, `httpStatusesList` and `certificateExpiryDateConstraint`
- `value` (String) The content to look for

Optional:

- `pass_if_found` (Boolean) The validation condition. `true` means validation succeeds if the specified content/element is found. `false` means validation fails if the specified content/element is found. Always specify `false` for `certificateExpiryDateConstraint` to fail the monitor if SSL certificate expiry is within the specified number of days





<a id="nestedblock--tags"></a>
### Nested Schema for `tags`

Optional:

- `tag` (Block Set) Tag with source of a Dynatrace entity. (see [below for nested schema](#nestedblock--tags--tag))

<a id="nestedblock--tags--tag"></a>
### Nested Schema for `tags.tag`

Required:

- `context` (String) The origin of the tag. Supported values are `AWS`, `AWS_GENERIC`, `AZURE`, `CLOUD_FOUNDRY`, `CONTEXTLESS`, `ENVIRONMENT`, `GOOGLE_CLOUD` and `KUBERNETES`.

Custom tags use the `CONTEXTLESS` value.
- `key` (String) The key of the tag.

Custom tags have the tag value here.

Optional:

- `source` (String) The source of the tag. Supported values are `USER`, `RULE_BASED` and `AUTO`.
- `value` (String) The value of the tag.

Not applicable to custom tags.
 