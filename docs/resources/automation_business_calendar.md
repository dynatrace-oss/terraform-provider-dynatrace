---
layout: ""
page_title: "dynatrace_automation_business_calendar Resource - terraform-provider-dynatrace"
subcategory: "Automation"
description: |-
  The resource `dynatrace_automation_business_calendar` covers configuration of Business Calendars for Workflows
---

# dynatrace_automation_business_calendar (Resource)

-> **Dynatrace SaaS only**

-> To utilize this resource, please define the environment variables `DT_CLIENT_ID`, `DT_CLIENT_SECRET`, `DT_ACCOUNT_ID` with an OAuth client including the following permissions: **View calendars** (`automation:calendars:read`) and **Create and edit calendars** (`automation:calendars:write`).

-> This resource is excluded by default in the export utility, please explicitly specify the resource to retrieve existing configuration.

## Dynatrace Documentation

- Dynatrace Workflows - https://www.dynatrace.com/support/help/platform-modules/cloud-automation/workflows

## Resource Example Usage

```terraform
resource "dynatrace_automation_business_calendar" "#name#" {
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
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `title` (String) The title / name of the Business Calendar

### Optional

- `description` (String) An optional description for the Business Calendar
- `holidays` (Block List, Max: 1) A list of holidays valid in this calendar (see [below for nested schema](#nestedblock--holidays))
- `valid_from` (String) The date from when on this calendar is valid from. Example: `2023-07-04` for July 4th 2023
- `valid_to` (String) The date until when on this calendar is valid to. Example: `2023-07-04` for July 4th 2023
- `week_days` (Set of Number) The days to be considered week days in this calendar. `1' = `Monday`, `2` = `Tuesday`, `3` = `Wednesday`, `4` = `Thursday`, `5` = `Friday`, `6` = `Saturday`, `7` = `Sunday`
- `week_start` (Number) Specifies the day of the week that's considered to be the first day in the week. `1` for Monday, `7` for Sunday

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--holidays"></a>
### Nested Schema for `holidays`

Required:

- `holiday` (Block Set, Min: 1) A (unordered) list of holidays valid in this calendar (see [below for nested schema](#nestedblock--holidays--holiday))

<a id="nestedblock--holidays--holiday"></a>
### Nested Schema for `holidays.holiday`

Required:

- `date` (String) The date of this holiday: Example `2017-07-04` for July 4th 2017
- `title` (String) An official name for this holiday
