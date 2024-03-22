resource "dynatrace_declarative_grouping" "first-instance" {
  name    = "#name#"
  enabled = true
  scope   = "environment"
  detection {
    process_definition {
      id                 = "PGIdentifierSample"
      process_group_name = "PGDisplayNameSample"
      report             = "always"
      rules {
        rule {
          condition = "$contains(TFExecutableSample)"
          property  = "executable"
        }
        rule {
          condition = "$contains(TFCommandSample)"
          property  = "commandLine"
        }
      }
    }
  }
}

resource "dynatrace_declarative_grouping" "second-instance" {
  name    = "#name#-second"
  enabled = true
  scope   = "environment"
  detection {
    process_definition {
      id                 = "PGIdentifierSample"
      process_group_name = "PGDisplayNameSample"
      report             = "always"
      rules {
        rule {
          condition = "$contains(TFExecutableSample)"
          property  = "executable"
        }
        rule {
          condition = "$contains(TFCommandSample)"
          property  = "commandLine"
        }
      }
    }
  }
  insert_after = dynatrace_declarative_grouping.first-instance.id
}
