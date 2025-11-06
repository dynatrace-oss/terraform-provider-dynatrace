---
layout: ""
page_title: "dynatrace_golden_state Resource - terraform-provider-dynatrace"
subcategory: "Incubator"
description: |-
  The resource `dynatrace_golden_state` allows for identifying configuration within a Dynatrace environment that isn't managed by Terraform
---

# dynatrace_golden_state (Resource)

!> This resource is currently in an experimental phase. It is disabled by default. If you would like to get early access please reach out to us via GitHub ticket. Dynatrace Support will not yet be able to assist you here.

The resource `dynatrace_golden_state` doesn't represent an actual setting that can get maintained within a Dynatrace Environment or on a Dynatrace Cluster.
Purpose of this resource is to easily identify whether there exist setttings on a Dynatrace environment that are not maintained by Terraform.

The resource allows you to specify for every supported resource a set of IDs. Any settings that are not among that set of IDs are then considered to be maintained outside of Terraform.
`dynatrace_golden_state` in that case either prints out warnings for these "non-terraform" settings or optionally automatically deletes them.

## Usage Examples

### Example A 

```terraform
resource "dynatrace_golden_state" "golden_state" {
  # mode = "WARN"
  dynatrace_management_zone_v2 = [
    dynatrace_management_zone_v2.team-mainframe.id,
    dynatrace_management_zone_v2.frontend.id,
    dynatrace_management_zone_v2.team-hawaiian-pizza.id
  ]
}
```

All Management Zones except the ones referred to with the IDs `dynatrace_management_zone_v2.team-mainframe.id`, `dynatrace_management_zone_v2.frontend.id` and `dynatrace_management_zone_v2.team-hawaiian-pizza.id` will be treated as "not maintained by Terraform'.
`terraform apply` will result with output similar to this:

```
dynatrace_golden_state.golden_state: Modifications complete after 2s [id=60109201-c276-4f31-8272-8662fc15818b]
╷
│ Warning: There exist resources of type `dynatrace_management_zone_v2` not managed by terraform
│
│   with dynatrace_golden_state.golden_state,
│   on golden_state.tf line 1, in resource "dynatrace_golden_state" "golden_state":
│    1: resource "dynatrace_golden_state" "golden_state" {
│
│ [ deleteme                 ] vu9U3hXa3q0AAAABABhidWlsdGluOm1hbmFnZW1lbnQtem9uZXMABnRlbmFudAAGdGVuYW50ACQ4ZDY0YzA3OC03YjdiLTMzZDMtYmQxNy1lNWI5MTU1OTgyMzO-71TeFdrerQ
│ [ uncoveredmgmz            ] vu9U3hXa3q0AAAABABhidWlsdGluOm1hbmFnZW1lbnQtem9uZXMABnRlbmFudAAGdGVuYW50ACQ5ZjkzNTNiZi1hNGQ3LTMxN2UtYjM5NS1jYzg5ZDUxNWExMjm-71TeFdrerQ
```

Resources other than `dynatrace_management_zone_v2` will not be taken into consideration.

-> Specifying `mode = "WARN"` is optional because `WARN` is the default value for the attribute `MODE`.

### Example B 

```terraform
resource "dynatrace_golden_state" "golden_state" {
  mode = "DELETE"
  dynatrace_management_zone_v2 = [
    dynatrace_management_zone_v2.team-mainframe.id,
    dynatrace_management_zone_v2.frontend.id,
    dynatrace_management_zone_v2.team-hawaiian-pizza.id
  ]
}
```
All Management Zones except the ones referred to with the IDs `dynatrace_management_zone_v2.team-mainframe.id`, `dynatrace_management_zone_v2.frontend.id` and `dynatrace_management_zone_v2.team-hawaiian-pizza.id` will be treated as "not maintained by Terraform'.

`terraform apply` will result in all Management Zones which's IDs are not within that given set of IDs to get automatically deleted.

Resources other than `dynatrace_management_zone_v2` will not be affected.

### Example C

```terraform
resource "dynatrace_golden_state" "golden_state" {
  # mode = "WARN"
  dynatrace_management_zone_v2 = [ ]
}
```
You're signaling to `dynatrace_golden_state` that you expect no Management Zones at all are configured on the Dynatrace Environment.

`terraform apply` will automatically delete existing Management Zones. Resources other than `dynatrace_management_zone_v2` will not be affected.

Resources other than `dynatrace_management_zone_v2` will not be taken into consideration.

### Example D

```terraform
resource "dynatrace_golden_state" "golden_state" {
  mode = "DELETE"
  dynatrace_management_zone_v2 = [ ]
}
```
You're signaling to `dynatrace_golden_state` that you expect no Management Zones at all are configured on the Dynatrace Environment.

`terraform apply` will print out a warning message like in Example A for all existing Management Zones. Resources other than `dynatrace_management_zone_v2` will not be affected.

### Example E

```terraform
resource "dynatrace_golden_state" "golden_state" {
  # mode = "WARN"
  dynatrace_management_zone_v2 = [
    dynatrace_management_zone_v2.team-mainframe.id,
    dynatrace_management_zone_v2.frontend.id,
    dynatrace_management_zone_v2.team-hawaiian-pizza.id
  ]
  dynatrace_alerting = [
    dynatrace_alerting.quick.id,
    dynatrace_alerting.slow.id
  ]
  dynatrace_autotag_v2 = [ ]  
}
```
You're signaling to `dynatrace_golden_state` that
* You expect only the Management Zones referred to with `dynatrace_management_zone_v2.team-mainframe.id`, `dynatrace_management_zone_v2.frontend.id`, `dynatrace_management_zone_v2.team-hawaiian-pizza.id` to exist the Dynatrace Environment
* You expect only the Alerting Profiles referred to with `dynatrace_alerting.quick.id` and `dynatrace_alerting.slow.id` to exist on the Dynatrace Environment
* You expect to see no Auto Tags to exist on the Dynatrace Environment

`terraform apply` will print out a warning message like in Example A for all Management Zones, Alerting Profiles and Auto Tags that don't match the IDs provided.

Resources other than `dynatrace_management_zone_v2`, `dynatrace_alerting` and `dynatrace_autotag_v2` will not be taken into consideration.

### Example F

```terraform
resource "dynatrace_golden_state" "golden_state" {
  mode = "DELETE"
  dynatrace_management_zone_v2 = [
    dynatrace_management_zone_v2.team-mainframe.id,
    dynatrace_management_zone_v2.frontend.id,
    dynatrace_management_zone_v2.team-hawaiian-pizza.id
  ]
  dynatrace_alerting = [
    dynatrace_alerting.quick.id,
    dynatrace_alerting.slow.id
  ]
  dynatrace_autotag_v2 = [ ]  
}
```
You're signaling to `dynatrace_golden_state` that
* You wish to delete all Management Zones except the ones referred to with `dynatrace_management_zone_v2.team-mainframe.id`, `dynatrace_management_zone_v2.frontend.id`, `dynatrace_management_zone_v2.team-hawaiian-pizza.id`
* You with to delete all Alerting Profiles expect the ones referred to with `dynatrace_alerting.quick.id` and `dynatrace_alerting.slow.id`
* You wish to delete all Auto Tags from the Dynatrace Environment

`terraform apply` will automatically delete all Management Zones, Alerting Profiles and Auto Tags that don't match the IDs provided.

Resources other than `dynatrace_management_zone_v2`, `dynatrace_alerting` and `dynatrace_autotag_v2` will not be taken into consideration.

### Example G

```terraform
resource "dynatrace_golden_state" "golden_state_auto_delete" {
  mode = "DELETE"
  dynatrace_management_zone_v2 = [
    dynatrace_management_zone_v2.team-mainframe.id,
    dynatrace_management_zone_v2.frontend.id,
    dynatrace_management_zone_v2.team-hawaiian-pizza.id
  ]
  dynatrace_alerting = [
    dynatrace_alerting.quick.id,
    dynatrace_alerting.slow.id
  ]
}

resource "dynatrace_golden_state" "golden_state_warn" {
  # mode = "WARN"
  dynatrace_autotag_v2 = [ ]  
}
```

Utilizing two separate `dynatrace_golden_state` resource blocks allows you to specify different behavior for different resources.
For some resources you may want Terraform to automatically delete them.
For some resources you may want to just get notified via warning messages.

You're signaling that
* You wish to delete all Management Zones except the ones referred to with `dynatrace_management_zone_v2.team-mainframe.id`, `dynatrace_management_zone_v2.frontend.id`, `dynatrace_management_zone_v2.team-hawaiian-pizza.id`
* You with to delete all Alerting Profiles expect the ones referred to with `dynatrace_alerting.quick.id` and `dynatrace_alerting.slow.id`
* You expect to see no Auto Tags to exist on the Dynatrace Environment

`terraform apply` will automatically delete all Management Zones, Alerting Profiles that don't match the IDs provided. For any Auto Tags that are configured it will print out a warning like in Example A.

!> Specifying the same resource in BOTH `dynatrace_golden_state` resource blocks may lead to unpredictable results. Make sure that e.g. `dynatrace_autotag_v2` is just specified in ONE of these two blocks.

Resources other than `dynatrace_management_zone_v2`, `dynatrace_alerting` and `dynatrace_autotag_v2` will not be taken into consideration.

## Frequently Asked Questions

### Wouldn't `terraform plan` reveal the same as the warning messages?
Running `terraform plan` in a lot of cases matches up with the warning messages produced when specifying `mode = "WARN"`. Unfortunately the plan doesn't exactly reflect the differences detected by `dynatrace_golden_state`.
In some cases the plan for resource `dynatrace_golden_state` shows difference that eventually don't result in any warnings or automatic deletions.

-> Be aware of the fact that `dynatrace_golden_state` shouldn't be considered a `classic` Terraform Resource. There doesn't exist a setting on the Dynatrace Environment that reflects exactly what's configured for that resource. The "remote state` needs to get queried for and deducted internally by the Provider.

### Why do I have to specify the "known" IDs explicitly?
We've initially aimed for functionality where you don't have to specify the "known" IDs explicitly - and that essentially automatically covers all resources within the same Terraform Module.
That's unfortunately easier said than done. Terraform doesn't offer any mechanisms that allows a Provider to query for what other resources exist within the same module.
We've also discussed evaluating the State file for that purpose. But of course that would immediately exclude environments where the state is getting managed remote - which is the case for most production environments.

### Warnings for settings not maintained by Terraform don't consistently contain the name of the setting in `[]`
Where applicable (e.g. Management Zones, Alerting Profiles, ...) a user readable name for a settings will show up within the error message - allowing you to identify it as easily as possible.
A lot of Settings offered by a Dynatrace Environment unfortunately cannot get identified by a name. Here you will have to work with the ID.

### `dynatrace_golden_state` doesn't recognize resource `dynatrace_xyz`, but that resource is supported by the Provider
It depends on the specific resource why that's the case
* Some Settings (e.g. `dynatrace_ddu_pool`) exist just as a singleton on a Dynatrace Environment. It makes much more sense to configure the settings you would like to enforce using a normal resource block. `dynatrace_golden_state` wouldn't bring in any benefits as it just focuses on deleting.
* Some Settings (e.g. `dynatrace_span_entry_point`) come prepopulated on every freshly provisioned Dynatrace Environment. Moreover thes "defaults" may even change over time (during cluster updates). We've decided to not mess around with these resources in order to avoid inconsistencies when `dynatrace_golden_state` is configured to automatically delete settings.
* The resource `dynatrace_golden_state` is not yet feature complete. Every resource it currently supports needs to get manually tested against a brand new environment. Essentially in order to ensure that neither of the two previous reasons are applicable. About 60% of all resources have been tested through already. We're adding support for additional resources gradually. 

### What resources are currently supported by `dynatrace_golden_state`?
  * dynatrace_management_zone_v2
  * dynatrace_alerting
  * dynatrace_autotag_v2
  * dynatrace_request_attribute
  * dynatrace_queue_manager
  * dynatrace_ims_bridges
  * dynatrace_custom_service
  * dynatrace_aws_credentials
  * dynatrace_azure_credentials
  * dynatrace_span_capture_rule
  * dynatrace_span_context_propagation
  * dynatrace_slo_v2
  * dynatrace_web_application
  * dynatrace_mobile_application
  * dynatrace_jira_notification
  * dynatrace_webhook_notification
  * dynatrace_ansible_tower_notification
  * dynatrace_email_notification
  * dynatrace_ops_genie_notification
  * dynatrace_pager_duty_notification
  * dynatrace_service_now_notification
  * dynatrace_slack_notification
  * dynatrace_trello_notification
  * dynatrace_victor_ops_notification
  * dynatrace_xmatters_notification
  * dynatrace_maintenance
  * dynatrace_metric_events
  * dynatrace_key_requests
  * dynatrace_credentials
  * dynatrace_calculated_service_metric
  * dynatrace_calculated_web_metric
  * dynatrace_calculated_mobile_metric
  * dynatrace_http_monitor
  * dynatrace_browser_monitor
  * dynatrace_calculated_synthetic_metric
  * dynatrace_host_naming
  * dynatrace_processgroup_naming
  * dynatrace_service_naming
  * dynatrace_request_naming
  * dynatrace_application_detection_rule
  * dynatrace_application_error_rules
  * dynatrace_synthetic_location
  * dynatrace_queue_sharing_groups
  * dynatrace_pg_alerting
  * dynatrace_database_anomalies_v2
  * dynatrace_process_monitoring_rule
  * dynatrace_disk_anomalies_v2
  * dynatrace_disk_specific_anomalies_v2
  * dynatrace_host_anomalies_v2
  * dynatrace_custom_app_anomalies
  * dynatrace_custom_app_crash_rate
  * dynatrace_process_monitoring
  * dynatrace_process_availability
  * dynatrace_process_group_detection
  * dynatrace_mobile_app_anomalies
  * dynatrace_mobile_app_crash_rate
  * dynatrace_web_app_anomalies
  * dynatrace_muted_requests
  * dynatrace_declarative_grouping
  * dynatrace_host_process_group_monitoring
  * dynatrace_rum_ip_locations
  * dynatrace_custom_app_enablement
  * dynatrace_mobile_app_enablement
  * dynatrace_web_app_enablement
  * dynatrace_process_group_rum
  * dynatrace_rum_provider_breakdown
  * dynatrace_web_app_resource_cleanup
  * dynatrace_update_windows
  * dynatrace_process_group_detection_flags
  * dynatrace_process_group_monitoring
  * dynatrace_process_group_simple_detection
  * dynatrace_log_metrics
  * dynatrace_browser_monitor_performance
  * dynatrace_session_replay_web_privacy
  * dynatrace_monitored_technologies_apache
  * dynatrace_monitored_technologies_dotnet
  * dynatrace_monitored_technologies_go
  * dynatrace_monitored_technologies_iis
  * dynatrace_monitored_technologies_java
  * dynatrace_monitored_technologies_nginx
  * dynatrace_monitored_technologies_nodejs
  * dynatrace_monitored_technologies_opentracing
  * dynatrace_monitored_technologies_php
  * dynatrace_monitored_technologies_varnish
  * dynatrace_monitored_technologies_wsmb
  * dynatrace_process_visibility
  * dynatrace_oneagent_features
  * dynatrace_rum_advanced_correlation
  * dynatrace_web_app_beacon_origins
  * dynatrace_web_app_resource_types
  * dynatrace_generic_types
  * dynatrace_data_privacy
  * dynatrace_service_failure
  * dynatrace_service_http_failure
  * dynatrace_disk_options
  * dynatrace_extension_execution_controller
  * dynatrace_nettracer
  * dynatrace_aix_extension
  * dynatrace_k8s_namespace_anomalies
