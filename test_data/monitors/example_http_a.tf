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
