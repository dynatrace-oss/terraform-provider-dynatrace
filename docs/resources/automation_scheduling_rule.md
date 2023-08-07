---
layout: ""
page_title: "dynatrace_automation_scheduling_rule Resource - terraform-provider-dynatrace"
subcategory: "Automation"
description: |-
  The resource `dynatrace_automation_scheduling_rule` covers configuration of Scheduling Rules for Workflows
---

# dynatrace_automation_scheduling_rule (Resource)

-> This resource is excluded by default in the export utility. You can, of course, specify that resource explicitly in order to export it. In that case, don't forget to specify the environment variables `DYNATRACE_AUTOMATION_CLIENT_ID` and `DYNATRACE_AUTOMATION_CLIENT_SECRET` for authentication.

## Dynatrace Documentation

- Dynatrace Workflows - https://www.dynatrace.com/support/help/platform-modules/cloud-automation/workflows

## Prerequisites

Using this resource requires an OAuth client to be configured within your account settings.
The scopes of the OAuth Client need to include `View rules (automation:rules:read)` and `Create and edit rules (automation:rules:write)`.

Finally the provider configuration requires the credentials for that OAuth Client.
The configuration section of your provider needs to look like this.
```terraform
provider "dynatrace" {
  dt_env_url   = "https://########.live.dynatrace.com/"  
  dt_api_token = "######.########################.################################################################"  

  # Usually not required. Terraform will deduct it if `dt_env_url` has been specified
  # automation_env_url = "https://########.apps.dynatrace.com/" 
  automation_client_id = "######.########"
  automation_client_secret = "######.########.################################################################"  
}
```
-> In order to handle credentials in a secure manner we recommend to use the environment variables `DYNATRACE_AUTOMATION_CLIENT_ID` and `DYNATRACE_AUTOMATION_CLIENT_SECRET` as an alternative.

## Resource Examples

### Recurrence Rule

```terraform
resource "dynatrace_automation_business_calendar" "calendar" {
  description = "#name#"
  title         = "#name#"
  valid_from    = "2023-07-31"
  valid_to      = "2033-07-31"
  week_days     = [1, 2, 3, 4, 5]
  week_start    = 1
  holidays {
    holiday {
      date  = "2023-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2023-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2023-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2023-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2023-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2023-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2024-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2024-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2024-04-01"
      title = "Ostermontag"
    }
    holiday {
      date  = "2024-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2024-05-09"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2024-05-20"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2024-05-30"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2024-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2024-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2024-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2024-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2024-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2024-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2025-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2025-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2025-04-21"
      title = "Ostermontag"
    }
    holiday {
      date  = "2025-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2025-05-29"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2025-06-09"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2025-06-19"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2025-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2025-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2025-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2025-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2025-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2025-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2026-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2026-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2026-04-06"
      title = "Ostermontag"
    }
    holiday {
      date  = "2026-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2026-05-14"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2026-05-25"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2026-06-04"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2026-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2026-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2026-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2026-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2026-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2026-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2027-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2027-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2027-03-29"
      title = "Ostermontag"
    }
    holiday {
      date  = "2027-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2027-05-06"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2027-05-17"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2027-05-27"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2027-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2027-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2027-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2027-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2027-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2027-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2028-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2028-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2028-04-17"
      title = "Ostermontag"
    }
    holiday {
      date  = "2028-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2028-05-25"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2028-06-05"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2028-06-15"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2028-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2028-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2028-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2028-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2028-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2028-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2029-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2029-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2029-04-02"
      title = "Ostermontag"
    }
    holiday {
      date  = "2029-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2029-05-10"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2029-05-21"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2029-05-31"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2029-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2029-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2029-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2029-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2029-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2029-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2030-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2030-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2030-04-22"
      title = "Ostermontag"
    }
    holiday {
      date  = "2030-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2030-05-30"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2030-06-10"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2030-06-20"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2030-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2030-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2030-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2030-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2030-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2030-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2031-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2031-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2031-04-14"
      title = "Ostermontag"
    }
    holiday {
      date  = "2031-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2031-05-22"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2031-06-02"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2031-06-12"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2031-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2031-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2031-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2031-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2031-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2031-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2032-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2032-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2032-03-29"
      title = "Ostermontag"
    }
    holiday {
      date  = "2032-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2032-05-06"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2032-05-17"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2032-05-27"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2032-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2032-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2032-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2032-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2032-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2032-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2033-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2033-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2033-04-18"
      title = "Ostermontag"
    }
    holiday {
      date  = "2033-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2033-05-26"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2033-06-06"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2033-06-16"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2023-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2024-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2025-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2026-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2027-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2028-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2029-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2030-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2031-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2032-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2033-07-31"
      title = "Reini Day"
    }
  }
}

resource "dynatrace_automation_scheduling_rule" "#name#" {
  business_calendar = dynatrace_automation_business_calendar.calendar.id
  title             = "#name#"
  recurrence {
    datestart     = "2023-07-31"
    days_in_month = [-1]
    days_in_year  = [-2, -1, 1, 2, 3]
    frequency     = "WEEKLY"
    interval      = 33
    months        = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12]
    weekdays      = [ "MO", "TU", "WE" ]
    weeks         = [-2, -1, 1, 2, 3]
    workdays      = "WORKING"
  }
}
```

### Fixed Offset Rule

```terraform
resource "dynatrace_automation_business_calendar" "calendar" {
  description = "#name#"
  title         = "#name#"
  valid_from    = "2023-07-31"
  valid_to      = "2033-07-31"
  week_days     = [1, 2, 3, 4, 5]
  week_start    = 1
  holidays {
    holiday {
      date  = "2023-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2023-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2023-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2023-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2023-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2023-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2024-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2024-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2024-04-01"
      title = "Ostermontag"
    }
    holiday {
      date  = "2024-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2024-05-09"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2024-05-20"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2024-05-30"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2024-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2024-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2024-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2024-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2024-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2024-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2025-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2025-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2025-04-21"
      title = "Ostermontag"
    }
    holiday {
      date  = "2025-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2025-05-29"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2025-06-09"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2025-06-19"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2025-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2025-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2025-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2025-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2025-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2025-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2026-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2026-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2026-04-06"
      title = "Ostermontag"
    }
    holiday {
      date  = "2026-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2026-05-14"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2026-05-25"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2026-06-04"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2026-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2026-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2026-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2026-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2026-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2026-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2027-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2027-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2027-03-29"
      title = "Ostermontag"
    }
    holiday {
      date  = "2027-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2027-05-06"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2027-05-17"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2027-05-27"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2027-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2027-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2027-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2027-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2027-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2027-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2028-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2028-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2028-04-17"
      title = "Ostermontag"
    }
    holiday {
      date  = "2028-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2028-05-25"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2028-06-05"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2028-06-15"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2028-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2028-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2028-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2028-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2028-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2028-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2029-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2029-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2029-04-02"
      title = "Ostermontag"
    }
    holiday {
      date  = "2029-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2029-05-10"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2029-05-21"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2029-05-31"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2029-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2029-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2029-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2029-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2029-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2029-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2030-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2030-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2030-04-22"
      title = "Ostermontag"
    }
    holiday {
      date  = "2030-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2030-05-30"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2030-06-10"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2030-06-20"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2030-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2030-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2030-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2030-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2030-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2030-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2031-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2031-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2031-04-14"
      title = "Ostermontag"
    }
    holiday {
      date  = "2031-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2031-05-22"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2031-06-02"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2031-06-12"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2031-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2031-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2031-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2031-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2031-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2031-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2032-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2032-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2032-03-29"
      title = "Ostermontag"
    }
    holiday {
      date  = "2032-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2032-05-06"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2032-05-17"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2032-05-27"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2032-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2032-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2032-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2032-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2032-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2032-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2033-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2033-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2033-04-18"
      title = "Ostermontag"
    }
    holiday {
      date  = "2033-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2033-05-26"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2033-06-06"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2033-06-16"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2023-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2024-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2025-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2026-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2027-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2028-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2029-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2030-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2031-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2032-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2033-07-31"
      title = "Reini Day"
    }
  }
}

resource "dynatrace_automation_scheduling_rule" "base" {
  business_calendar = dynatrace_automation_business_calendar.calendar.id
  title             = "#name#"
  recurrence {
    datestart     = "2023-07-31"
    days_in_month = [-1]
    days_in_year  = [-2, -1, 1, 2, 3]
    frequency     = "WEEKLY"
    interval      = 33
    months        = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12]
    weekdays      = [ "MO", "TU", "WE" ]
    weeks         = [-2, -1, 1, 2, 3]
    workdays      = "WORKING"
  }
}

resource "dynatrace_automation_scheduling_rule" "#name#" {
  title         = "#name#"
  fixed_offset {
    offset = 50
    rule   = dynatrace_automation_scheduling_rule.base.id
  }
}
```

### Relative Offset rule

```terraform
resource "dynatrace_automation_business_calendar" "calendar" {
  description = "#name#"
  title         = "#name#"
  valid_from    = "2023-07-31"
  valid_to      = "2033-07-31"
  week_days     = [1, 2, 3, 4, 5]
  week_start    = 1
  holidays {
    holiday {
      date  = "2023-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2023-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2023-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2023-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2023-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2023-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2024-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2024-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2024-04-01"
      title = "Ostermontag"
    }
    holiday {
      date  = "2024-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2024-05-09"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2024-05-20"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2024-05-30"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2024-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2024-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2024-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2024-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2024-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2024-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2025-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2025-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2025-04-21"
      title = "Ostermontag"
    }
    holiday {
      date  = "2025-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2025-05-29"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2025-06-09"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2025-06-19"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2025-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2025-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2025-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2025-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2025-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2025-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2026-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2026-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2026-04-06"
      title = "Ostermontag"
    }
    holiday {
      date  = "2026-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2026-05-14"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2026-05-25"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2026-06-04"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2026-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2026-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2026-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2026-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2026-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2026-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2027-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2027-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2027-03-29"
      title = "Ostermontag"
    }
    holiday {
      date  = "2027-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2027-05-06"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2027-05-17"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2027-05-27"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2027-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2027-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2027-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2027-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2027-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2027-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2028-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2028-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2028-04-17"
      title = "Ostermontag"
    }
    holiday {
      date  = "2028-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2028-05-25"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2028-06-05"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2028-06-15"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2028-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2028-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2028-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2028-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2028-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2028-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2029-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2029-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2029-04-02"
      title = "Ostermontag"
    }
    holiday {
      date  = "2029-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2029-05-10"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2029-05-21"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2029-05-31"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2029-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2029-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2029-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2029-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2029-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2029-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2030-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2030-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2030-04-22"
      title = "Ostermontag"
    }
    holiday {
      date  = "2030-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2030-05-30"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2030-06-10"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2030-06-20"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2030-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2030-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2030-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2030-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2030-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2030-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2031-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2031-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2031-04-14"
      title = "Ostermontag"
    }
    holiday {
      date  = "2031-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2031-05-22"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2031-06-02"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2031-06-12"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2031-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2031-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2031-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2031-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2031-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2031-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2032-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2032-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2032-03-29"
      title = "Ostermontag"
    }
    holiday {
      date  = "2032-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2032-05-06"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2032-05-17"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2032-05-27"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2032-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2032-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2032-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2032-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2032-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2032-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2033-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2033-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2033-04-18"
      title = "Ostermontag"
    }
    holiday {
      date  = "2033-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2033-05-26"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2033-06-06"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2033-06-16"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2023-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2024-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2025-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2026-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2027-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2028-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2029-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2030-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2031-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2032-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2033-07-31"
      title = "Reini Day"
    }
  }
}

resource "dynatrace_automation_scheduling_rule" "base" {
  business_calendar = dynatrace_automation_business_calendar.calendar.id
  title             = "#name#"
  recurrence {
    datestart     = "2023-07-31"
    days_in_month = [-1]
    days_in_year  = [-2, -1, 1, 2, 3]
    frequency     = "WEEKLY"
    interval      = 33
    months        = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12]
    weekdays      = [ "MO", "TU", "WE" ]
    weeks         = [-2, -1, 1, 2, 3]
    workdays      = "WORKING"
  }
}

resource "dynatrace_automation_scheduling_rule" "source" {
  business_calendar = dynatrace_automation_business_calendar.calendar.id
  title             = "#name#"
  recurrence {
    datestart     = "2023-07-31"
    days_in_month = [-1]
    days_in_year  = [-2, -1, 1, 2, 3]
    frequency     = "WEEKLY"
    interval      = 33
    months        = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12]
    weekdays      = [ "MO", "TU", "WE" ]
    weeks         = [-2, -1, 1, 2, 3]
    workdays      = "WORKING"
  }
}

resource "dynatrace_automation_scheduling_rule" "target" {
  title         = "#name#"
  fixed_offset {
    offset = 50
    rule   = dynatrace_automation_scheduling_rule.base.id
  }
}

resource "dynatrace_automation_scheduling_rule" "#name#" {
  title         = "#name#"
  relative_offset {
    direction   = "previous"
    source_rule = dynatrace_automation_scheduling_rule.source.id
    target_rule = dynatrace_automation_scheduling_rule.target.id
  }
}
```


### Grouping Rule

```terraform
resource "dynatrace_automation_business_calendar" "calendar" {
  description = "#name#"
  title         = "#name#"
  valid_from    = "2023-07-31"
  valid_to      = "2033-07-31"
  week_days     = [1, 2, 3, 4, 5]
  week_start    = 1
  holidays {
    holiday {
      date  = "2023-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2023-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2023-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2023-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2023-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2023-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2024-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2024-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2024-04-01"
      title = "Ostermontag"
    }
    holiday {
      date  = "2024-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2024-05-09"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2024-05-20"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2024-05-30"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2024-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2024-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2024-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2024-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2024-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2024-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2025-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2025-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2025-04-21"
      title = "Ostermontag"
    }
    holiday {
      date  = "2025-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2025-05-29"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2025-06-09"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2025-06-19"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2025-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2025-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2025-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2025-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2025-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2025-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2026-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2026-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2026-04-06"
      title = "Ostermontag"
    }
    holiday {
      date  = "2026-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2026-05-14"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2026-05-25"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2026-06-04"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2026-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2026-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2026-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2026-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2026-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2026-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2027-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2027-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2027-03-29"
      title = "Ostermontag"
    }
    holiday {
      date  = "2027-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2027-05-06"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2027-05-17"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2027-05-27"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2027-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2027-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2027-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2027-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2027-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2027-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2028-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2028-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2028-04-17"
      title = "Ostermontag"
    }
    holiday {
      date  = "2028-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2028-05-25"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2028-06-05"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2028-06-15"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2028-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2028-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2028-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2028-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2028-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2028-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2029-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2029-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2029-04-02"
      title = "Ostermontag"
    }
    holiday {
      date  = "2029-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2029-05-10"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2029-05-21"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2029-05-31"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2029-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2029-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2029-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2029-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2029-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2029-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2030-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2030-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2030-04-22"
      title = "Ostermontag"
    }
    holiday {
      date  = "2030-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2030-05-30"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2030-06-10"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2030-06-20"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2030-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2030-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2030-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2030-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2030-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2030-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2031-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2031-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2031-04-14"
      title = "Ostermontag"
    }
    holiday {
      date  = "2031-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2031-05-22"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2031-06-02"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2031-06-12"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2031-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2031-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2031-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2031-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2031-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2031-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2032-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2032-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2032-03-29"
      title = "Ostermontag"
    }
    holiday {
      date  = "2032-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2032-05-06"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2032-05-17"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2032-05-27"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2032-08-15"
      title = "Mariä Himmelfahrt"
    }
    holiday {
      date  = "2032-10-26"
      title = "Nationalfeiertag"
    }
    holiday {
      date  = "2032-11-01"
      title = "Allerheiligen"
    }
    holiday {
      date  = "2032-12-08"
      title = "Mariä Empfängnis"
    }
    holiday {
      date  = "2032-12-25"
      title = "Christtag"
    }
    holiday {
      date  = "2032-12-26"
      title = "Stefanitag"
    }
    holiday {
      date  = "2033-01-01"
      title = "Neujahr"
    }
    holiday {
      date  = "2033-01-06"
      title = "Heilige Drei Könige"
    }
    holiday {
      date  = "2033-04-18"
      title = "Ostermontag"
    }
    holiday {
      date  = "2033-05-01"
      title = "Staatsfeiertag"
    }
    holiday {
      date  = "2033-05-26"
      title = "Christi Himmelfahrt"
    }
    holiday {
      date  = "2033-06-06"
      title = "Pfingstmontag"
    }
    holiday {
      date  = "2033-06-16"
      title = "Fronleichnam"
    }
    holiday {
      date  = "2023-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2024-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2025-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2026-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2027-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2028-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2029-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2030-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2031-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2032-07-31"
      title = "Reini Day"
    }
    holiday {
      date  = "2033-07-31"
      title = "Reini Day"
    }
  }
}


resource "dynatrace_automation_scheduling_rule" "subtract" {
  business_calendar = dynatrace_automation_business_calendar.calendar.id
  title             = "#name#"
  recurrence {
    datestart     = "2029-07-31"
    frequency     = "DAILY"
    interval      = 1
    workdays      = "WORKING"
  }
}

resource "dynatrace_automation_scheduling_rule" "comba" {
  business_calendar = dynatrace_automation_business_calendar.calendar.id
  title             = "#name#"
  recurrence {
    datestart     = "2023-07-31"
    frequency     = "DAILY"
    interval      = 1
    workdays      = "WORKING"
  }
}

resource "dynatrace_automation_scheduling_rule" "combb" {
  business_calendar = dynatrace_automation_business_calendar.calendar.id
  title             = "#name#"
  recurrence {
    datestart     = "2023-07-31"
    frequency     = "DAILY"
    interval      = 1
    workdays      = "WORKING"
  }
}

resource "dynatrace_automation_scheduling_rule" "intersect" {
  business_calendar = dynatrace_automation_business_calendar.calendar.id
  title             = "#name#"
  recurrence {
    datestart     = "2023-07-31"
    frequency     = "DAILY"
    interval      = 1
    workdays      = "WORKING"
  }
}

resource "dynatrace_automation_scheduling_rule" "#name#" {
  title         = "#name#"
  grouping {
    combine   = [ dynatrace_automation_scheduling_rule.comba.id, dynatrace_automation_scheduling_rule.combb.id ]
    intersect = [ dynatrace_automation_scheduling_rule.intersect.id ]
    subtract  = [ dynatrace_automation_scheduling_rule.subtract.id ]
  }
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `title` (String) The title / name of the scheduling rule

### Optional

- `business_calendar` (String)
- `description` (String) An optional description for the scheduling rule
- `fixed_offset` (Block List, Max: 1) (see [below for nested schema](#nestedblock--fixed_offset))
- `grouping` (Block List, Max: 1) (see [below for nested schema](#nestedblock--grouping))
- `recurrence` (Block List, Max: 1) (see [below for nested schema](#nestedblock--recurrence))
- `relative_offset` (Block List, Max: 1) (see [below for nested schema](#nestedblock--relative_offset))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--fixed_offset"></a>
### Nested Schema for `fixed_offset`

Required:

- `offset` (Number) Every day of the scheduling rule referred to with `rule` will be offset by this amount of days
- `rule` (String) Refers to a scheduling rule for which to produce valid days with an offset


<a id="nestedblock--grouping"></a>
### Nested Schema for `grouping`

Required:

- `combine` (Set of String) The IDs of scheduling rules determining the days the schedule should apply to

Optional:

- `intersect` (Set of String) The IDs of scheduling rules determining the days the schedule is allowed apply to. If specified, only days that are covered by `combine` and `intersect` are valid days for the schedule
- `subtract` (Set of String) The IDs of scheduling rules determing the days the schedule must not apply. If specified it reduces down the set of days covered by `combine` and `intersect`


<a id="nestedblock--recurrence"></a>
### Nested Schema for `recurrence`

Required:

- `datestart` (String) The recurrence start. Example: `2017-07-04` represents July 4th 2017
- `frequency` (String) Possible values are `YEARLY`, `MONTHLY`, `WEEKLY`, `DAILY`, `HOURLY`, `MINUTELY` and `SECONDLY`. Example: `frequency` = `DAILY` and `interval` = `2` schedules for every other day
- `workdays` (String) Possible values are `WORKING` (Work days), `HOLIDAYS` (Holidays) and `OFF` (Weekends + Holidays)

Optional:

- `days_in_month` (Set of Number) Restricts the recurrence to specific days within a month. `1`, `2`, `3`, ... refers to the first, second, third day in the month. You can also specify negative values to refer to values relative to the last day. `-1` refers to the last day, `-2` refers to the second to the last day, ...
- `days_in_year` (Set of Number) Restricts the recurrence to specific days within a year. `1`, `2`, `3`, ... refers to the first, second, third day of the year. You can also specify negative values to refer to values relative to the last day. `-1` refers to the last day, `-2` refers to the second to the last day, ...
- `easter` (Set of Number) Restricts the recurrence to specific days relative to Easter Sunday. `0` will yield the Easter Sunday itself
- `interval` (Number) The interval between each iteration. Default: 1. Example: `frequency` = `DAILY` and `interval` = `2` schedules for every other day
- `months` (Set of Number) Restricts the recurrence to specific months. `1` for `January`, `2` for `February`, ..., `12` for `December`
- `weekdays` (Set of String) Restricts the recurrence to specific week days. Possible values are `MO`, `TU`, `WE`, `TH`, `FR`, `SA` and `SU`
- `weeks` (Set of Number) Restricts the recurrence to specific weeks within a year. `1`, `2`, `3`, ... refers to the first, second, third week of the year. You can also specify negative values to refer to values relative to the last week. `-1` refers to the last week, `-2` refers to the second to the last week, ...


<a id="nestedblock--relative_offset"></a>
### Nested Schema for `relative_offset`

Required:

- `direction` (String)
- `source_rule` (String)
- `target_rule` (String)
