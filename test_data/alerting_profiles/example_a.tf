resource "dynatrace_alerting_profile" "#name#" {
	mz_id = ""
	rules {
		tag_filter {
			include_mode = "NONE"
		}
		delay_in_minutes = 0
		severity_level = "AVAILABILITY"
	}
	rules {
		tag_filter {
			include_mode = "NONE"
		}
		delay_in_minutes = 0
		severity_level = "ERROR"
	}
	rules {
		tag_filter {
			include_mode = "NONE"
		}
		delay_in_minutes = 0
		severity_level = "PERFORMANCE"
	}
	rules {
		tag_filter {
			include_mode = "NONE"
		}
		delay_in_minutes = 0
		severity_level = "RESOURCE_CONTENTION"
	}
	rules {
		tag_filter {
			include_mode = "NONE"
		}
		delay_in_minutes = 0
		severity_level = "CUSTOM_ALERT"
	}
	rules {
		tag_filter {
			include_mode = "NONE"
		}
		delay_in_minutes = 0
		severity_level = "MONITORING_UNAVAILABLE"
	}
	display_name = "#name#"
	metadata {
		cluster_version = "1.206.95.20201116-094826"
		current_configuration_versions = [ "0" ]
	}
}
