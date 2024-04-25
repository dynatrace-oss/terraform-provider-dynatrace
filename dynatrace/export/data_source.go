package export

import (
	"fmt"
	"strings"
)

type DataSourceKind string

var DataSourceKindTenant = DataSourceKind("tenant")
var DataSourceKindEntity = DataSourceKind("entity")
var DataSourceKindPolicy = DataSourceKind("policy")

type DataSource struct {
	ID   string
	Type string
	Name string
	Kind DataSourceKind
}

func AsDataSource(resource *Resource) string {
	if resource == nil {
		return ""
	}
	switch resource.Type {
	case ResourceTypes.ManagementZoneV2:
		return fmt.Sprintf(`data "dynatrace_management_zone_v2" "%s" {
			name = "%s"
		}`, resource.UniqueName, esc(resource.Name))
	case ResourceTypes.Alerting:
		return fmt.Sprintf(`data "dynatrace_alerting_profile" "%s" {
			name = "%s"
		}`, resource.UniqueName, esc(resource.Name))
	case ResourceTypes.RequestAttribute:
		return fmt.Sprintf(`data "dynatrace_request_attribute" "%s" {
			name = "%s"
		}`, resource.UniqueName, esc(resource.Name))
	case ResourceTypes.WebApplication:
		return fmt.Sprintf(`data "dynatrace_application" "%s" {
			name = "%s"
		}`, resource.UniqueName, esc(resource.Name))
	case ResourceTypes.RequestNaming:
		return fmt.Sprintf(`data "dynatrace_request_naming" "%s" {
			name = "%s"
		}`, resource.UniqueName, esc(resource.Name))
	case ResourceTypes.JSONDashboard:
		splitName := strings.Split(resource.Name, " owned by ")
		return fmt.Sprintf(`data "dynatrace_dashboard" "%s" {
			name = "%s"
			owner = "%s"
		}`, resource.UniqueName, esc(splitName[0]), esc(splitName[1]))
	case ResourceTypes.SLO:
		return fmt.Sprintf(`data "dynatrace_slo" "%s" {
			name = "%s"
		}`, resource.UniqueName, esc(resource.Name))
	case ResourceTypes.CalculatedServiceMetric:
		return fmt.Sprintf(`data "dynatrace_calculated_service_metric" "%s" {
			name = "%s"
		}`, resource.UniqueName, esc(resource.Name))
	case ResourceTypes.MobileApplication:
		return fmt.Sprintf(`data "dynatrace_mobile_application" "%s" {
			name = "%s"
		}`, resource.UniqueName, esc(resource.Name))
	case ResourceTypes.Credentials:
		return fmt.Sprintf(`data "dynatrace_credentials" "%s" {
			name = "%s"
		}`, resource.UniqueName, esc(resource.Name))
	case ResourceTypes.HTTPMonitor:
		return fmt.Sprintf(`data "dynatrace_entity" "%s" {
			type = "HTTP_CHECK"
			name = "%s"
		}`, resource.UniqueName, esc(resource.Name))
	case ResourceTypes.BrowserMonitor:
		return fmt.Sprintf(`data "dynatrace_entity" "%s" {
			type = "SYNTHETIC_TEST"
			name = "%s"
		}`, resource.UniqueName, esc(resource.Name))
	case ResourceTypes.SyntheticLocation:
		return fmt.Sprintf(`data "dynatrace_entity" "%s" {
			type = "SYNTHETIC_LOCATION"
			name = "%s"
		}`, resource.UniqueName, esc(resource.Name))
	case ResourceTypes.FailureDetectionParameters:
		return fmt.Sprintf(`data "dynatrace_failure_detection_parameters" "%s" {
			name = "%s"
		}`, resource.UniqueName, esc(resource.Name))
	case ResourceTypes.UpdateWindows:
		return fmt.Sprintf(`data "dynatrace_update_windows" "%s" {
			name = "%s"
		}`, resource.UniqueName, esc(resource.Name))
	case ResourceTypes.AWSCredentials:
		return fmt.Sprintf(`data "dynatrace_aws_credentials" "%s" {
			label = "%s"
		}`, resource.UniqueName, esc(resource.Name))
	case ResourceTypes.AzureCredentials:
		return fmt.Sprintf(`data "dynatrace_azure_credentials" "%s" {
			label = "%s"
		}`, resource.UniqueName, esc(resource.Name))
	case ResourceTypes.IAMGroup:
		return fmt.Sprintf(`data "dynatrace_iam_group" "%s" {
			name = "%s"
		}`, resource.UniqueName, esc(resource.Name))
	case ResourceTypes.IAMPolicy:
		return fmt.Sprintf(`data "dynatrace_iam_policy" "%s" {
			name = "%s"
		}`, resource.UniqueName, esc(resource.Name))
	case ResourceTypes.AppSecVulnerabilityAlerting:
		return fmt.Sprintf(`data "dynatrace_vulnerability_alerting" "%s" {
			name = "%s"
		}`, resource.UniqueName, esc(resource.Name))
	case ResourceTypes.AppSecAttackAlerting:
		return fmt.Sprintf(`data "dynatrace_attack_alerting" "%s" {
			name = "%s"
		}`, resource.UniqueName, esc(resource.Name))
	default:
		return ""
	}
}

func esc(s string) string {
	return strings.ReplaceAll(s, "\"", "\\\"")
}
