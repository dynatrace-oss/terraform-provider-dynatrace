resource "dynatrace_maintenance_window" "#name#" {
  suppression = "DONT_DETECT_PROBLEMS"
  type        = "PLANNED"
  name        = "#name#"
  schedule {
    end             = "2021-05-11 14:41"
    zone_id         = "Europe/Vienna"
    recurrence_type = "ONCE"
    start           = "2021-05-11 13:41"
  }
  suppress_synth_mon_exec = true
  scope {
    matches {
      tag_combination = "AND"
      tags {
        key     = "bggtedgxen"
        context = "CONTEXTLESS"
      }
      tags {
        context = "CONTEXTLESS"
        key     = "deldel1"
      }
    }
  }
}
