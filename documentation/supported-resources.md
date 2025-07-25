# Supported Resources

Dynatrace Configuration as Code via Terraform supports the following resources, please reference the resource specific pages for additional information. 

## API Token Resources
| Resource name | API endpoint | API token permissions |
| ------------- | ----------- | --------------------- |
| dynatrace_activegate_token | /api/v2/settings/objects (schema: builtin:activegate-token) | settings.read, settings.write |
| dynatrace_activegate_updates | /api/v2/settings/objects (schema: builtin:deployment.activegate.updates) | settings.read, settings.write |
| dynatrace_ag_token | /api/v2/activeGateTokens | activeGateTokenManagement.create, activeGateTokenManagement.read, activeGateTokenManagement.write |
| dynatrace_aix_extension | /api/v2/settings/objects (schema: builtin:host.monitoring.aix-kernel-extension) | settings.read, settings.write |
| dynatrace_alerting | /api/v2/settings/objects (schema: builtin:alerting.profile) | settings.read, settings.write |
| dynatrace_ansible_tower_notification | /api/v2/settings/objects (schema: builtin:problem.notifications) | settings.read, settings.write |
| dynatrace_api_detection | /api/v2/settings/objects (schema: builtin:apis.detection-rules) | settings.read, settings.write |
| dynatrace_api_token | /api/v2/slo | apiTokens.read, apiTokens.write |
| dynatrace_app_monitoring | /api/v2/settings/objects (schema: builtin:dt-javascript-runtime.app-monitoring) | settings.read, settings.write |
| dynatrace_application_detection_rule | /api/config/v1/applicationDetectionRules | ReadConfig, WriteConfig |
| dynatrace_application_error_rules | /api/config/v1/applications/web/{id}/errorRules | ReadConfig, WriteConfig |
| dynatrace_appsec_notification | /api/v2/settings/objects (schema: builtin:appsec.notification-integration) | securityProblems.read, securityProblems.write |
| dynatrace_attack_alerting | /api/v2/settings/objects (schema: builtin:appsec.notification-attack-alerting-profile) | securityProblems.read, securityProblems.write |
| dynatrace_attack_allowlist | /api/v2/settings/objects (schema: builtin:appsec.attack-protection-allowlist-config) | attacks.read, attacks.write |
| dynatrace_attack_rules | /api/v2/settings/objects (schema: builtin:appsec.attack-protection-advanced-config) | attacks.read, attacks.write |
| dynatrace_attack_settings | /api/v2/settings/objects (schema: builtin:appsec.attack-protection-settings) | attacks.read, attacks.write |
| dynatrace_attribute_allow_list | /api/v2/settings/objects (schema: builtin:attribute-allow-list) | settings.read, settings.write |
| dynatrace_attribute_block_list | /api/v2/settings/objects (schema: builtin:attribute-block-list) | settings.read, settings.write |
| dynatrace_attribute_masking | /api/v2/settings/objects (schema: builtin:attribute-masking) | settings.read, settings.write |
| dynatrace_attributes_preferences | /api/v2/settings/objects (schema: builtin:attributes-preferences) | settings.read, settings.write |
| dynatrace_audit_log | /api/v2/settings/objects (schema: builtin:audit-log) | settings.read, settings.write |
| dynatrace_automation_controller_connections | /api/v2/settings/objects (schema: app:dynatrace.redhat.ansible:automation-controller.connection) | settings.read, settings.write |
| dynatrace_automation_workflow_aws_connections | /api/v2/settings/objects (schema: builtin:hyperscaler-authentication.aws.connection) | settings.read, settings.write |
| dynatrace_automation_workflow_jira | /api/v2/settings/objects (schema: app:dynatrace.jira:connection) | settings.read, settings.write |
| dynatrace_automation_workflow_k8s_connections | /api/v2/settings/objects (schema: app:dynatrace.kubernetes.connector:connection) | settings.read, settings.write |
| dynatrace_automation_workflow_slack | /api/v2/settings/objects (schema: app:dynatrace.slack:connection) | settings.read, settings.write |
| dynatrace_autotag_rules | /api/v2/settings/objects (schema: builtin:tags.auto-tagging-rules) | settings.read, settings.write |
| dynatrace_autotag_v2 | /api/v2/settings/objects (schema: builtin:tags.auto-tagging) | settings.read, settings.write |
| dynatrace_aws_anomalies | /api/v2/settings/objects (schema: builtin:anomaly-detection.infrastructure-aws) | settings.read, settings.write |
| dynatrace_aws_credentials | /api/config/v1/aws/credentials | ReadConfig, WriteConfig |
| dynatrace_aws_service | /api/config/v1/aws/credentials/{id}/services | ReadConfig, WriteConfig |
| dynatrace_azure_credentials | /api/config/v1/azure/credentials | ReadConfig, WriteConfig |
| dynatrace_azure_service | /api/config/v1/azure/credentials/{id}/services | ReadConfig, WriteConfig |
| dynatrace_browser_monitor | /api/v1/synthetic/monitors | ExternalSyntheticIntegration |
| dynatrace_browser_monitor_outage | /api/v2/settings/objects (schema: builtin:synthetic.browser.outage-handling) | settings.read, settings.write |
| dynatrace_browser_monitor_performance | /api/v2/settings/objects (schema: builtin:synthetic.browser.performance-thresholds) | settings.read, settings.write |
| dynatrace_builtin_process_monitoring | /api/v2/settings/objects (schema: builtin:process.built-in-process-monitoring-rule) | settings.read, settings.write |
| dynatrace_business_events_buckets | /api/v2/settings/objects (schema: builtin:bizevents-processing-buckets.rule) | settings.read, settings.write |
| dynatrace_business_events_capturing_variants | /api/v2/settings/objects (schema: builtin:bizevents.http.capturing-variants) | settings.read, settings.write |
| dynatrace_business_events_metrics | /api/v2/settings/objects (schema: builtin:bizevents-processing-metrics.rule) | settings.read, settings.write |
| dynatrace_business_events_oneagent | /api/v2/settings/objects (schema: builtin:bizevents.http.incoming) | settings.read, settings.write |
| dynatrace_business_events_oneagent_outgoing | /api/v2/settings/objects (schema: builtin:bizevents.http.outgoing) | settings.read, settings.write |
| dynatrace_business_events_processing | /api/v2/settings/objects (schema: builtin:bizevents-processing-pipelines.rule) | settings.read, settings.write |
| dynatrace_business_events_security_context | /api/v2/settings/objects (schema: builtin:bizevents-security-context-rules) | settings.read, settings.write |
| dynatrace_calculated_mobile_metric | /api/config/v1/calculatedMetrics/mobile | ReadConfig, WriteConfig |
| dynatrace_calculated_service_metric | /api/config/v1/calculatedMetrics/service | ReadConfig, WriteConfig |
| dynatrace_calculated_synthetic_metric | /api/config/v1/calculatedMetrics/synthetic | ReadConfig, WriteConfig |
| dynatrace_calculated_web_metric | /api/config/v1/calculatedMetrics/rum | ReadConfig, WriteConfig |
| dynatrace_cloud_development_environments | /api/v2/settings/objects (schema: builtin:app-engine-registry.cloud-development-environments) | settings.read, settings.write |
| dynatrace_cloud_foundry | /api/v2/settings/objects (schema: builtin:cloud.cloudfoundry) | settings.read, settings.write |
| dynatrace_cloudapp_workloaddetection | /api/v2/settings/objects (schema: builtin:process-group.cloud-application-workload-detection) | settings.read, settings.write |
| dynatrace_connectivity_alerts | /api/v2/settings/objects (schema: builtin:alerting.connectivity-alerts) | settings.read, settings.write |
| dynatrace_container_builtin_rule | /api/v2/settings/objects (schema: builtin:container.built-in-monitoring-rule) | settings.read, settings.write |
| dynatrace_container_registry | /api/v2/settings/objects (schema: builtin:container-registry) | settings.read, settings.write |
| dynatrace_container_rule | /api/v2/settings/objects (schema: builtin:container.monitoring-rule) | settings.read, settings.write |
| dynatrace_container_technology | /api/v2/settings/objects (schema: builtin:container.technology) | settings.read, settings.write |
| dynatrace_crashdump_analytics | /api/v2/settings/objects (schema: builtin:crashdump.analytics) | settings.read, settings.write |
| dynatrace_credentials | /api/config/v1/credentials | credentialVault.read, credentialVault.write |
| dynatrace_custom_app_anomalies | /api/v2/settings/objects (schema: builtin:anomaly-detection.rum-custom) | settings.read, settings.write |
| dynatrace_custom_app_crash_rate | /api/v2/settings/objects (schema: builtin:anomaly-detection.rum-custom-crash-rate-increase) | settings.read, settings.write |
| dynatrace_custom_app_enablement | /api/v2/settings/objects (schema: builtin:rum.custom.enablement) | settings.read, settings.write |
| dynatrace_custom_device | /api/v2/entities/custom | entities.read, entities.write |
| dynatrace_custom_service | /api/config/v1/service/customServices/{technology} | ReadConfig, WriteConfig |
| dynatrace_custom_tags | /api/v2/tags | entities.read, entities.write |
| dynatrace_custom_units | /api/v2/settings/objects (schema: builtin:custom-unit) | settings.read, settings.write |
| dynatrace_dashboard_sharing | /api/config/v1/dashboards/{id}/shareSettings | ReadConfig, WriteConfig |
| dynatrace_dashboards_allowlist | /api/v2/settings/objects (schema: builtin:dashboards.image.allowlist) | settings.read, settings.write |
| dynatrace_dashboards_general | /api/v2/settings/objects (schema: builtin:dashboards.general) | settings.read, settings.write |
| dynatrace_dashboards_presets | /api/v2/settings/objects (schema: builtin:dashboards.presets) | settings.read, settings.write |
| dynatrace_data_privacy | /api/v2/settings/objects (schema: builtin:preferences.privacy) | settings.read, settings.write |
| dynatrace_database_anomalies_v2 | /api/v2/settings/objects (schema: builtin:anomaly-detection.databases) | settings.read, settings.write |
| dynatrace_davis_anomaly_detectors | /api/v2/settings/objects (schema: builtin:davis.anomaly-detectors) | settings:objects:read, settings:objects:write, storage:bizevents:read |
| dynatrace_davis_copilot | /api/v2/settings/objects (schema: service:davis.copilot.datamining-blocklist) | settings.read, settings.write |
| dynatrace_db_app_feature_flags | /api/v2/settings/objects (schema: app:dynatrace.database.overview:feature-flags) | settings.read, settings.write |
| dynatrace_declarative_grouping | /api/v2/settings/objects (schema: builtin:declarativegrouping) | settings.read, settings.write |
| dynatrace_default_launchpad | /api/v2/settings/objects (schema: app:dynatrace.launcher:default.launchpad) | settings.read, settings.write |
| dynatrace_devobs_agent_optin | /api/v2/settings/objects (schema: builtin:devobs.agent.optin) | settings.read, settings.write |
| dynatrace_devobs_data_masking | /api/v2/settings/objects (schema: builtin:devobs.sensitive.data.masking) | settings.read, settings.write |
| dynatrace_devobs_git_onprem | /api/v2/settings/objects (schema: app:dynatrace.devobs.debugger:git.on.prem) | settings.read, settings.write |
| dynatrace_discovery_default_rules | /api/v2/settings/objects (schema: app:dynatrace.discovery.coverage:discovery.findings.default.rules.schema) | settings.read, settings.write |
| dynatrace_discovery_feature_flags | /api/v2/settings/objects (schema: app:dynatrace.discovery.coverage:feature-flags) | settings.read, settings.write |
| dynatrace_disk_analytics | /api/v2/settings/objects (schema: builtin:disk.analytics.extension) | settings.read, settings.write |
| dynatrace_disk_anomalies_v2 | /api/v2/settings/objects (schema: builtin:anomaly-detection.infrastructure-disks) | settings.read, settings.write |
| dynatrace_disk_anomaly_rules | /api/v2/settings/objects (schema: builtin:anomaly-detection.disk-rules) | settings.read, settings.write |
| dynatrace_disk_edge_anomaly_detectors | /api/v2/settings/objects (schema: builtin:infrastructure.disk.edge.anomaly-detectors) | settings.read, settings.write |
| dynatrace_disk_options | /api/v2/settings/objects (schema: builtin:disk.options) | settings.read, settings.write |
| dynatrace_disk_specific_anomalies_v2 | /api/v2/settings/objects (schema: builtin:anomaly-detection.infrastructure-disks.per-disk-override) | settings.read, settings.write |
| dynatrace_ebpf_service_discovery | /api/v2/settings/objects (schema: builtin:ebpf.service.discovery) | settings.read, settings.write |
| dynatrace_email_notification | /api/v2/settings/objects (schema: builtin:problem.notifications) | settings.read, settings.write |
| dynatrace_environment | /api/cluster/v2/environments | ServiceProviderAPI |
| dynatrace_eula_settings | /api/v2/settings/objects (schema: builtin:eula-settings) | settings.read, settings.write |
| dynatrace_event_driven_ansible_connections | /api/v2/settings/objects (schema: app:dynatrace.redhat.ansible:eda-webhook.connection) | settings.read, settings.write |
| dynatrace_extension_execution_controller | /api/v2/settings/objects (schema: builtin:eec.local) | settings.read, settings.write |
| dynatrace_extension_execution_remote | /api/v2/settings/objects (schema: builtin:eec.remote) | settings.read, settings.write |
| dynatrace_failure_detection_parameters | /api/v2/settings/objects (schema: builtin:failure-detection.environment.parameters) | settings.read, settings.write |
| dynatrace_failure_detection_rules | /api/v2/settings/objects (schema: builtin:failure-detection.environment.rules) | settings.read, settings.write |
| dynatrace_frequent_issues | /api/v2/settings/objects (schema: builtin:anomaly-detection.frequent-issues) | settings.read, settings.write |
| dynatrace_generic_relationships | /api/v2/settings/objects (schema: builtin:monitoredentities.generic.relation) | settings.read, settings.write |
| dynatrace_generic_setting | /api/v2/settings/objects (schema: generic) | settings.read, settings.write |
| dynatrace_generic_types | /api/v2/settings/objects (schema: builtin:monitoredentities.generic.type) | settings.read, settings.write |
| dynatrace_geolocation | /api/v2/settings/objects (schema: builtin:geo-settings) | settings.read, settings.write |
| dynatrace_github_connection | /api/v2/settings/objects (schema: app:dynatrace.github.connector:connection) | settings.read, settings.write |
| dynatrace_grail_metrics_allowall | /api/v2/settings/objects (schema: builtin:grail.metrics.allow-all) | settings.read, settings.write |
| dynatrace_grail_metrics_allowlist | /api/v2/settings/objects (schema: builtin:grail.metrics.allow-list) | settings.read, settings.write |
| dynatrace_grail_security_context | /api/v2/settings/objects (schema: builtin:monitoredentities.grail.security.context) | settings.read, settings.write |
| dynatrace_histogram_metrics | /api/v2/settings/objects (schema: builtin:histogram-metrics) | settings.read, settings.write |
| dynatrace_host_anomalies_v2 | /api/v2/settings/objects (schema: builtin:anomaly-detection.infrastructure-hosts) | settings.read, settings.write |
| dynatrace_host_monitoring | /api/v2/settings/objects (schema: builtin:host.monitoring) | settings.read, settings.write |
| dynatrace_host_monitoring_advanced | /api/v2/settings/objects (schema: builtin:host.monitoring.advanced) | settings.read, settings.write |
| dynatrace_host_monitoring_mode | /api/v2/settings/objects (schema: builtin:host.monitoring.mode) | settings.read, settings.write |
| dynatrace_host_naming | /api/config/v1/conditionalNaming/host | ReadConfig, WriteConfig |
| dynatrace_host_process_group_monitoring | /api/v2/settings/objects (schema: builtin:host.process-groups.monitoring-state) | settings.read, settings.write |
| dynatrace_http_monitor | /api/v1/synthetic/monitors | ExternalSyntheticIntegration |
| dynatrace_http_monitor_cookies | /api/v2/settings/objects (schema: builtin:synthetic.http.cookies) | settings.read, settings.write |
| dynatrace_http_monitor_outage | /api/v2/settings/objects (schema: builtin:synthetic.http.outage-handling) | settings.read, settings.write |
| dynatrace_http_monitor_performance | /api/v2/settings/objects (schema: builtin:synthetic.http.performance-thresholds) | settings.read, settings.write |
| dynatrace_http_monitor_script | /api/v1/synthetic/monitors | ExternalSyntheticIntegration |
| dynatrace_hub_extension_active_version | /api/v2/extensions/{extensionName}/environmentConfiguration | extensions.write, extensionEnvironment.read, extensionEnvironment.write |
| dynatrace_hub_extension_config | /api/v2/extensions/{extensionName}/monitoringConfigurations | extensions.read, extensions.write, hub.read |
| dynatrace_hub_permissions | /api/v2/settings/objects (schema: app:dynatrace.hub:manage.permissions) | settings.read, settings.write |
| dynatrace_hub_subscriptions | /api/v2/settings/objects (schema: builtin:hub-channel.subscriptions) | settings.read, settings.write |
| dynatrace_ibm_mq_filters | /api/v2/settings/objects (schema: builtin:mainframe.mqfilters) | settings.read, settings.write |
| dynatrace_ims_bridges | /api/v2/settings/objects (schema: builtin:ibmmq.ims-bridges) | settings.read, settings.write |
| dynatrace_infraops_app_feature_flags | /api/v2/settings/objects (schema: app:dynatrace.infraops:feature-flags) | settings.read, settings.write |
| dynatrace_infraops_app_settings | /api/v2/settings/objects (schema: app:dynatrace.infraops:settings) | settings.read, settings.write |
| dynatrace_ip_address_masking | /api/v2/settings/objects (schema: builtin:preferences.ipaddressmasking) | settings.read, settings.write |
| dynatrace_issue_tracking | /api/v2/settings/objects (schema: builtin:issue-tracking.integration) | settings.read, settings.write |
| dynatrace_jenkins_connection | /api/v2/settings/objects (schema: app:dynatrace.jenkins.connector:connection) | settings.read, settings.write |
| dynatrace_jira_notification | /api/v2/settings/objects (schema: builtin:problem.notifications) | settings.read, settings.write |
| dynatrace_json_dashboard | /api/config/v1/dashboards | ReadConfig, WriteConfig |
| dynatrace_json_dashboard_base | /api/config/v1/dashboards | ReadConfig, WriteConfig |
| dynatrace_k8s_cluster_anomalies | /api/v2/settings/objects (schema: builtin:anomaly-detection.kubernetes.cluster) | settings.read, settings.write |
| dynatrace_k8s_monitoring | /api/v2/settings/objects (schema: builtin:cloud.kubernetes.monitoring) | settings.read, settings.write |
| dynatrace_k8s_namespace_anomalies | /api/v2/settings/objects (schema: builtin:anomaly-detection.kubernetes.namespace) | settings.read, settings.write |
| dynatrace_k8s_node_anomalies | /api/v2/settings/objects (schema: builtin:anomaly-detection.kubernetes.node) | settings.read, settings.write |
| dynatrace_k8s_pvc_anomalies | /api/v2/settings/objects (schema: builtin:anomaly-detection.kubernetes.pvc) | settings.read, settings.write |
| dynatrace_k8s_workload_anomalies | /api/v2/settings/objects (schema: builtin:anomaly-detection.kubernetes.workload) | settings.read, settings.write |
| dynatrace_key_requests | /api/v2/settings/objects (schema: builtin:settings.subscriptions.service) | settings.read, settings.write |
| dynatrace_key_user_action | /api/config/v1/applications/web/{id}/keyUserActions | ReadConfig, WriteConfig |
| dynatrace_kubernetes | /api/v2/settings/objects (schema: builtin:cloud.kubernetes) | settings.read, settings.write |
| dynatrace_kubernetes_app | /api/v2/settings/objects (schema: builtin:app-transition.kubernetes) | settings.read, settings.write |
| dynatrace_kubernetes_enrichment | /api/v2/settings/objects (schema: builtin:kubernetes.generic.metadata.enrichment) | settings.read, settings.write |
| dynatrace_kubernetes_spm | /api/v2/settings/objects (schema: builtin:kubernetes.security-posture-management) | settings.read, settings.write |
| dynatrace_limit_outbound_connections | /api/v2/settings/objects (schema: builtin:dt-javascript-runtime.allowed-outbound-connections) | settings.read, settings.write |
| dynatrace_log_agent_feature_flags | /api/v2/settings/objects (schema: builtin:logmonitoring.log-agent-feature-flags) | settings.read, settings.write |
| dynatrace_log_buckets | /api/v2/settings/objects (schema: builtin:logmonitoring.log-buckets-rules) | settings.read, settings.write |
| dynatrace_log_custom_attribute | /api/v2/settings/objects (schema: builtin:logmonitoring.log-custom-attributes) | settings.read, settings.write |
| dynatrace_log_custom_source | /api/v2/settings/objects (schema: builtin:logmonitoring.custom-log-source-settings) | settings.read, settings.write |
| dynatrace_log_debug_settings | /api/v2/settings/objects (schema: builtin:logmonitoring.log-debug-settings) | settings.read, settings.write |
| dynatrace_log_events | /api/v2/settings/objects (schema: builtin:logmonitoring.log-events) | settings.read, settings.write |
| dynatrace_log_metrics | /api/v2/settings/objects (schema: builtin:logmonitoring.schemaless-log-metric) | settings.read, settings.write |
| dynatrace_log_oneagent | /api/v2/settings/objects (schema: builtin:logmonitoring.log-agent-configuration) | settings.read, settings.write |
| dynatrace_log_processing | /api/v2/settings/objects (schema: builtin:logmonitoring.log-dpp-rules) | settings.read, settings.write |
| dynatrace_log_security_context | /api/v2/settings/objects (schema: builtin:logmonitoring.log-security-context-rules) | settings.read, settings.write |
| dynatrace_log_sensitive_data_masking | /api/v2/settings/objects (schema: builtin:logmonitoring.sensitive-data-masking-settings) | settings.read, settings.write |
| dynatrace_log_storage | /api/v2/settings/objects (schema: builtin:logmonitoring.log-storage-settings) | settings.read, settings.write |
| dynatrace_log_timestamp | /api/v2/settings/objects (schema: builtin:logmonitoring.timestamp-configuration) | settings.read, settings.write |
| dynatrace_mainframe_transaction_monitoring | /api/v2/settings/objects (schema: builtin:mainframe.txmonitoring) | settings.read, settings.write |
| dynatrace_maintenance | /api/v2/settings/objects (schema: builtin:alerting.maintenance-window) | settings.read, settings.write |
| dynatrace_managed_backup | /api/v1.0/onpremise/backup/config | ServiceProviderAPI |
| dynatrace_managed_internet_proxy | /api/v1.0/onpremise/proxy/configuration | ControlManagement, ServiceProviderAPI, or UnattendedInstall |
| dynatrace_managed_network_zones | /api/cluster/v2/networkZones/{id} | ServiceProviderAPI |
| dynatrace_managed_preferences | /api/v1.0/onpremise/preferences | ServiceProviderAPI |
| dynatrace_managed_public_endpoints | /api/v1.0/onpremise/endpoint/additionalWebUiAddresses, /api/v1.0/onpremise/endpoint/beaconForwarderAddress, /api/v1.0/onpremise/endpoint/cdnAddress, /api/v1.0/onpremise/endpoint/webUiAddress | ServiceProviderAPI |
| dynatrace_managed_remote_access | /api/cluster/v2/remoteaccess/requests | ServiceProviderAPI |
| dynatrace_managed_smtp | /api/v1.0/onpremise/smtp | ServiceProviderAPI |
| dynatrace_management_zone_v2 | /api/v2/settings/objects (schema: builtin:management-zones) | settings.read, settings.write |
| dynatrace_metric_events | /api/v2/settings/objects (schema: builtin:anomaly-detection.metric-events) | settings.read, settings.write |
| dynatrace_metric_metadata | /api/v2/settings/objects (schema: builtin:metric.metadata) | settings.read, settings.write |
| dynatrace_metric_query | /api/v2/settings/objects (schema: builtin:metric.query) | settings.read, settings.write |
| dynatrace_mgmz_permission | /api/v1.0/onpremise/groups/managementZones | ServiceProviderAPI |
| dynatrace_mobile_app_anomalies | /api/v2/settings/objects (schema: builtin:anomaly-detection.rum-mobile) | settings.read, settings.write |
| dynatrace_mobile_app_crash_rate | /api/v2/settings/objects (schema: builtin:anomaly-detection.rum-mobile-crash-rate-increase) | settings.read, settings.write |
| dynatrace_mobile_app_enablement | /api/v2/settings/objects (schema: builtin:rum.mobile.enablement) | settings.read, settings.write |
| dynatrace_mobile_app_key_performance | /api/v2/settings/objects (schema: builtin:rum.mobile.key-performance-metrics) | settings.read, settings.write |
| dynatrace_mobile_application | /api/config/v1/applications/mobile | ReadConfig, WriteConfig |
| dynatrace_mobile_notifications | /api/v2/settings/objects (schema: builtin:mobile.notifications) | settings.read, settings.write |
| dynatrace_monitored_technologies_apache | /api/v2/settings/objects (schema: builtin:monitored-technologies.apache) | settings.read, settings.write |
| dynatrace_monitored_technologies_dotnet | /api/v2/settings/objects (schema: builtin:monitored-technologies.dotnet) | settings.read, settings.write |
| dynatrace_monitored_technologies_go | /api/v2/settings/objects (schema: builtin:monitored-technologies.go) | settings.read, settings.write |
| dynatrace_monitored_technologies_iis | /api/v2/settings/objects (schema: builtin:monitored-technologies.iis) | settings.read, settings.write |
| dynatrace_monitored_technologies_java | /api/v2/settings/objects (schema: builtin:monitored-technologies.java) | settings.read, settings.write |
| dynatrace_monitored_technologies_nginx | /api/v2/settings/objects (schema: builtin:monitored-technologies.nginx) | settings.read, settings.write |
| dynatrace_monitored_technologies_nodejs | /api/v2/settings/objects (schema: builtin:monitored-technologies.nodejs) | settings.read, settings.write |
| dynatrace_monitored_technologies_opentracing | /api/v2/settings/objects (schema: builtin:monitored-technologies.open-tracing-native) | settings.read, settings.write |
| dynatrace_monitored_technologies_php | /api/v2/settings/objects (schema: builtin:monitored-technologies.php) | settings.read, settings.write |
| dynatrace_monitored_technologies_python | /api/v2/settings/objects (schema: builtin:monitored-technologies.python) | settings.read, settings.write |
| dynatrace_monitored_technologies_varnish | /api/v2/settings/objects (schema: builtin:monitored-technologies.varnish) | settings.read, settings.write |
| dynatrace_monitored_technologies_wsmb | /api/v2/settings/objects (schema: builtin:monitored-technologies.wsmb) | settings.read, settings.write |
| dynatrace_ms365_email_connection | /api/v2/settings/objects (schema: app:dynatrace.microsoft365.connector:mail.connection) | settings.read, settings.write |
| dynatrace_msteams_connection | /api/v2/settings/objects (schema: app:dynatrace.msteams:connection) | settings.read, settings.write |
| dynatrace_muted_requests | /api/v2/settings/objects (schema: builtin:settings.mutedrequests) | settings.read, settings.write |
| dynatrace_nettracer | /api/v2/settings/objects (schema: builtin:nettracer.traffic) | settings.read, settings.write |
| dynatrace_network_monitor | /api/v2/settings/objects (schema: v2:synthetic:monitors:network) | settings.read, settings.write |
| dynatrace_network_monitor_outage | /api/v2/settings/objects (schema: builtin:synthetic.multiprotocol.outage-handling) | settings.read, settings.write |
| dynatrace_network_traffic | /api/v2/settings/objects (schema: builtin:exclude.network.traffic) | settings.read, settings.write |
| dynatrace_network_zone | /api/v2/networkZones/{id} | networkZones.read, networkZones.write |
| dynatrace_network_zones | /api/v2/settings/objects (schema: builtin:networkzones) | settings.read, settings.write |
| dynatrace_oneagent_default_mode | /api/v2/settings/objects (schema: builtin:deployment.oneagent.default-mode) | settings.read, settings.write |
| dynatrace_oneagent_default_version | /api/v2/settings/objects (schema: builtin:deployment.oneagent.default-version) | settings.read, settings.write |
| dynatrace_oneagent_features | /api/v2/settings/objects (schema: builtin:oneagent.features) | settings.read, settings.write |
| dynatrace_oneagent_side_masking | /api/v2/settings/objects (schema: builtin:oneagent.side.masking.settings) | settings.read, settings.write |
| dynatrace_oneagent_updates | /api/v2/settings/objects (schema: builtin:deployment.oneagent.updates) | settings.read, settings.write |
| dynatrace_opentelemetry_metrics | /api/v2/settings/objects (schema: builtin:opentelemetry-metrics) | settings.read, settings.write |
| dynatrace_ops_genie_notification | /api/v2/settings/objects (schema: builtin:problem.notifications) | settings.read, settings.write |
| dynatrace_os_services | /api/v2/settings/objects (schema: builtin:os-services-monitoring) | settings.read, settings.write |
| dynatrace_ownership_config | /api/v2/settings/objects (schema: builtin:ownership.config) | settings.read, settings.write |
| dynatrace_ownership_teams | /api/v2/settings/objects (schema: builtin:ownership.teams) | settings.read, settings.write |
| dynatrace_pager_duty_notification | /api/v2/settings/objects (schema: builtin:problem.notifications) | settings.read, settings.write |
| dynatrace_pagerduty_connection | /api/v2/settings/objects (schema: app:dynatrace.pagerduty:connection) | settings.read, settings.write |
| dynatrace_pg_alerting | /api/v2/settings/objects (schema: builtin:availability.process-group-alerting) | settings.read, settings.write |
| dynatrace_policy | /api/cluster/v2/iam/repo/{levelType}/{levelId}/policies | ServiceProviderAPI |
| dynatrace_policy_bindings | /api/cluster/v2/iam/repo/{levelType}/{levelId}/bindings/{policyUuid} | ServiceProviderAPI |
| dynatrace_problem_fields | /api/v2/settings/objects (schema: builtin:problem.fields) | settings.read, settings.write |
| dynatrace_problem_record_propagation_rules | /api/v2/settings/objects (schema: builtin:problem.record.propagation.rules) | settings.read, settings.write |
| dynatrace_process_availability | /api/v2/settings/objects (schema: builtin:processavailability) | settings.read, settings.write |
| dynatrace_process_group_detection | /api/v2/settings/objects (schema: builtin:process-group.advanced-detection-rule) | settings.read, settings.write |
| dynatrace_process_group_detection_flags | /api/v2/settings/objects (schema: builtin:process-group.detection-flags) | settings.read, settings.write |
| dynatrace_process_group_monitoring | /api/v2/settings/objects (schema: builtin:process-group.monitoring.state) | settings.read, settings.write |
| dynatrace_process_group_rum | /api/v2/settings/objects (schema: builtin:rum.processgroup) | settings.read, settings.write |
| dynatrace_process_group_simple_detection | /api/v2/settings/objects (schema: builtin:process-group.simple-detection-rule) | settings.read, settings.write |
| dynatrace_process_monitoring | /api/v2/settings/objects (schema: builtin:process.process-monitoring) | settings.read, settings.write |
| dynatrace_process_monitoring_rule | /api/v2/settings/objects (schema: builtin:process.custom-process-monitoring-rule) | settings.read, settings.write |
| dynatrace_process_visibility | /api/v2/settings/objects (schema: builtin:process-visibility) | settings.read, settings.write |
| dynatrace_processgroup_naming | /api/config/v1/conditionalNaming/processGroup | ReadConfig, WriteConfig |
| dynatrace_queue_manager | /api/v2/settings/objects (schema: builtin:ibmmq.queue-managers) | settings.read, settings.write |
| dynatrace_queue_sharing_groups | /api/v2/settings/objects (schema: builtin:ibmmq.queue-sharing-group) | settings.read, settings.write |
| dynatrace_remote_environments | /api/v2/settings/objects (schema: builtin:remote.environment) | settings.read, settings.write |
| dynatrace_report | /api/config/v1/reports | ReadConfig, WriteConfig |
| dynatrace_request_attribute | /api/config/v1/service/requestAttributes | ReadConfig, CaptureRequestData |
| dynatrace_request_naming | /api/config/v1/service/requestNaming | ReadConfig, WriteConfig |
| dynatrace_request_namings | /api/config/v1/service/requestNaming/order | ReadConfig, WriteConfig |
| dynatrace_resource_attributes | /api/v2/settings/objects (schema: builtin:resource-attribute) | settings.read, settings.write |
| dynatrace_rpc_based_sampling | /api/v2/settings/objects (schema: builtin:rpc-based-sampling) | settings.read, settings.write |
| dynatrace_rum_advanced_correlation | /api/v2/settings/objects (schema: builtin:rum.resource-timing-origins) | settings.read, settings.write |
| dynatrace_rum_host_headers | /api/v2/settings/objects (schema: builtin:rum.host-headers) | settings.read, settings.write |
| dynatrace_rum_ip_determination | /api/v2/settings/objects (schema: builtin:rum.ip-determination) | settings.read, settings.write |
| dynatrace_rum_ip_locations | /api/v2/settings/objects (schema: builtin:rum.ip-mappings) | settings.read, settings.write |
| dynatrace_rum_overload_prevention | /api/v2/settings/objects (schema: builtin:rum.overload-prevention) | settings.read, settings.write |
| dynatrace_rum_provider_breakdown | /api/v2/settings/objects (schema: builtin:rum.provider-breakdown) | settings.read, settings.write |
| dynatrace_security_context | /api/v2/settings/objects (schema: builtin:security-context) | settings.read, settings.write |
| dynatrace_service_anomalies_v2 | /api/v2/settings/objects (schema: builtin:anomaly-detection.services) | settings.read, settings.write |
| dynatrace_service_detection_rules | /api/v2/settings/objects (schema: builtin:service-detection-rules) | settings.read, settings.write |
| dynatrace_service_external_web_request | /api/v2/settings/objects (schema: builtin:service-detection.external-web-request) | settings.read, settings.write |
| dynatrace_service_external_web_service | /api/v2/settings/objects (schema: builtin:service-detection.external-web-service) | settings.read, settings.write |
| dynatrace_service_failure | /api/v2/settings/objects (schema: builtin:failure-detection.service.general-parameters) | settings.read, settings.write |
| dynatrace_service_full_web_request | /api/v2/settings/objects (schema: builtin:service-detection.full-web-request) | settings.read, settings.write |
| dynatrace_service_full_web_service | /api/v2/settings/objects (schema: builtin:service-detection.full-web-service) | settings.read, settings.write |
| dynatrace_service_http_failure | /api/v2/settings/objects (schema: builtin:failure-detection.service.http-parameters) | settings.read, settings.write |
| dynatrace_service_naming | /api/config/v1/conditionalNaming/service | ReadConfig, WriteConfig |
| dynatrace_service_now_notification | /api/v2/settings/objects (schema: builtin:problem.notifications) | settings.read, settings.write |
| dynatrace_service_splitting | /api/v2/settings/objects (schema: builtin:service-splitting-rules) | settings.read, settings.write |
| dynatrace_servicenow_connection | /api/v2/settings/objects (schema: app:dynatrace.servicenow:connection) | settings.read, settings.write |
| dynatrace_session_replay_resource_capture | /api/v2/settings/objects (schema: builtin:sessionreplay.web.resource-capturing) | settings.read, settings.write |
| dynatrace_session_replay_web_privacy | /api/v2/settings/objects (schema: builtin:sessionreplay.web.privacy-preferences) | settings.read, settings.write |
| dynatrace_site_reliability_guardian | /api/v2/settings/objects (schema: app:dynatrace.site.reliability.guardian:guardians) | settings.read, settings.write |
| dynatrace_slack_notification | /api/v2/settings/objects (schema: builtin:problem.notifications) | settings.read, settings.write |
| dynatrace_slo | /api/v2/apiTokens | slo.read, slo.write |
| dynatrace_slo_normalization | /api/v2/settings/objects (schema: builtin:monitoring.slo.normalization) | settings.read, settings.write |
| dynatrace_slo_v2 | /api/v2/settings/objects (schema: builtin:monitoring.slo) | settings.read, settings.write |
| dynatrace_span_attribute | /api/v2/settings/objects (schema: builtin:span-attribute) | settings.read, settings.write |
| dynatrace_span_capture_rule | /api/v2/settings/objects (schema: builtin:span-capturing) | settings.read, settings.write |
| dynatrace_span_context_propagation | /api/v2/settings/objects (schema: builtin:span-context-propagation) | settings.read, settings.write |
| dynatrace_span_entry_point | /api/v2/settings/objects (schema: builtin:span-entry-points) | settings.read, settings.write |
| dynatrace_span_events | /api/v2/settings/objects (schema: builtin:span-event-attribute) | settings.read, settings.write |
| dynatrace_synthetic_availability | /api/v2/settings/objects (schema: builtin:synthetic.synthetic-availability-settings) | settings.read, settings.write |
| dynatrace_synthetic_location | /api/v1/synthetic/locations | ExternalSyntheticIntegration |
| dynatrace_token_settings | /api/v2/settings/objects (schema: builtin:tokens.token-settings) | settings.read, settings.write |
| dynatrace_transaction_start_filters | /api/v2/settings/objects (schema: builtin:mainframe.txstartfilters) | settings.read, settings.write |
| dynatrace_trello_notification | /api/v2/settings/objects (schema: builtin:problem.notifications) | settings.read, settings.write |
| dynatrace_unified_services_opentel | /api/v2/settings/objects (schema: builtin:unified-services-enablement) | settings.read, settings.write |
| dynatrace_update_windows | /api/v2/settings/objects (schema: builtin:deployment.management.update-windows) | settings.read, settings.write |
| dynatrace_url_based_sampling | /api/v2/settings/objects (schema: builtin:url-based-sampling) | settings.read, settings.write |
| dynatrace_usability_analytics | /api/v2/settings/objects (schema: builtin:usability-analytics) | settings.read, settings.write |
| dynatrace_user | /api/v1.0/onpremise/users | ServiceProviderAPI |
| dynatrace_user_action_metrics | /api/v2/settings/objects (schema: builtin:user-action-custom-metrics) | settings.read, settings.write |
| dynatrace_user_experience_score | /api/v2/settings/objects (schema: builtin:rum.user-experience-score) | settings.read, settings.write |
| dynatrace_user_group | /api/v1.0/onpremise/groups | ServiceProviderAPI |
| dynatrace_user_session_metrics | /api/v2/settings/objects (schema: builtin:custom-metrics) | settings.read, settings.write |
| dynatrace_victor_ops_notification | /api/v2/settings/objects (schema: builtin:problem.notifications) | settings.read, settings.write |
| dynatrace_vmware | /api/v2/settings/objects (schema: builtin:virtualization.vmware) | settings.read, settings.write |
| dynatrace_vmware_anomalies | /api/v2/settings/objects (schema: builtin:anomaly-detection.infrastructure-vmware) | settings.read, settings.write |
| dynatrace_vulnerability_alerting | /api/v2/settings/objects (schema: builtin:appsec.notification-alerting-profile) | securityProblems.read, securityProblems.write |
| dynatrace_vulnerability_code | /api/v2/settings/objects (schema: builtin:appsec.code-level-vulnerability-rule-settings) | securityProblems.read, securityProblems.write |
| dynatrace_vulnerability_settings | /api/v2/settings/objects (schema: builtin:appsec.runtime-vulnerability-detection) | securityProblems.read, securityProblems.write |
| dynatrace_vulnerability_third_party | /api/v2/settings/objects (schema: builtin:appsec.rule-settings) | securityProblems.read, securityProblems.write |
| dynatrace_vulnerability_third_party_attr | /api/v2/settings/objects (schema: builtin:appsec.third-party-vulnerability-rule-settings) | securityProblems.read, securityProblems.write |
| dynatrace_vulnerability_third_party_k8s | /api/v2/settings/objects (schema: builtin:appsec.third-party-vulnerability-kubernetes-label-rule-settings) | securityProblems.read, securityProblems.write |
| dynatrace_web_app_anomalies | /api/v2/settings/objects (schema: builtin:anomaly-detection.rum-web) | settings.read, settings.write |
| dynatrace_web_app_auto_injection | /api/v2/settings/objects (schema: builtin:rum.web.automatic-injection) | settings.read, settings.write |
| dynatrace_web_app_beacon_endpoint | /api/v2/settings/objects (schema: builtin:rum.web.beacon-endpoint) | settings.read, settings.write |
| dynatrace_web_app_beacon_origins | /api/v2/settings/objects (schema: builtin:rum.web.beacon-domain-origins) | settings.read, settings.write |
| dynatrace_web_app_custom_config_properties | /api/v2/settings/objects (schema: builtin:rum.web.custom-configuration-properties) | settings.read, settings.write |
| dynatrace_web_app_custom_prop_restrictions | /api/v2/settings/objects (schema: builtin:rum.web.capture-custom-properties) | settings.read, settings.write |
| dynatrace_web_app_custom_injection | /api/v2/settings/objects (schema: builtin:rum.web.custom-injection-rules) | settings.read, settings.write |
| dynatrace_web_app_enablement | /api/v2/settings/objects (schema: builtin:rum.web.enablement) | settings.read, settings.write |
| dynatrace_web_app_injection_cookie | /api/v2/settings/objects (schema: builtin:rum.web.injection.cookie) | settings.read, settings.write |
| dynatrace_web_app_ip_address_exclusion | /api/v2/settings/objects (schema: builtin:rum.web.ipaddress-exclusion) | settings.read, settings.write |
| dynatrace_web_app_javascript_filename| /api/v2/settings/objects (schema: builtin:rum.web.rum-javascript-file-name) | settings.read, settings.write |
| dynatrace_web_app_javascript_updates | /api/v2/settings/objects (schema: builtin:rum.web.rum-javascript-updates) | settings.read, settings.write |
| dynatrace_web_app_javascript_version | /api/v2/settings/objects (schema: builtin:rum.web.custom-rum-javascript-version) | settings.read, settings.write |
| dynatrace_web_app_key_performance_custom | /api/v2/settings/objects (schema: builtin:rum.web.key-performance-metric-custom-actions) | settings.read, settings.write |
| dynatrace_web_app_key_performance_load | /api/v2/settings/objects (schema: builtin:rum.web.key-performance-metric-load-actions) | settings.read, settings.write |
| dynatrace_web_app_key_performance_xhr | /api/v2/settings/objects (schema: builtin:rum.web.key-performance-metric-xhr-actions) | settings.read, settings.write |
| dynatrace_web_app_manual_insertion | /api/v2/settings/objects (schema: builtin:rum.web.manual-insertion) | settings.read, settings.write |
| dynatrace_web_app_resource_cleanup | /api/v2/settings/objects (schema: builtin:rum.web.resource-cleanup-rules) | settings.read, settings.write |
| dynatrace_web_app_resource_types | /api/v2/settings/objects (schema: builtin:rum.web.resource-types) | settings.read, settings.write |
| dynatrace_web_application | /api/config/v1/applications/web | ReadConfig, WriteConfig |
| dynatrace_webhook_notification | /api/v2/settings/objects (schema: builtin:problem.notifications) | settings.read, settings.write |
| dynatrace_xmatters_notification | /api/v2/settings/objects (schema: builtin:problem.notifications) | settings.read, settings.write |


## OAuth Resources
| Resource name | API endpoint | OAuth permissions |
| ------------- | ----------- | ---------------- |
| dynatrace_automation_business_calendar | /platform/automation/v1/business-calendars | automation:calendars:read, automation:calendars:write |
| dynatrace_automation_scheduling_rule | /platform/automation/v1/scheduling-rules | automation:rules:read, automation:rules:write |
| dynatrace_automation_workflow | /platform/automation/v1/workflows | automation:workflows:read, automation:workflows:write |
| dynatrace_direct_shares | /platform/document/v1/direct-shares | document:direct-shares:read, document:direct-shares:write, document:direct-shares:delete |
| dynatrace_document | /platform/document/v1/documents | document:documents:read, document:documents:write, document:documents:delete, document:trash.documents:delete |
| dynatrace_iam_group | /iam/v1/accounts/{accountUuid}/groups | account-idm-read, account-idm-write |
| dynatrace_iam_permission | /iam/v1/accounts/{accountUuid}/groups/{groupUuid}/permissions | account-idm-read, account-idm-write |
| dynatrace_iam_policy | /iam/v1/repo/{levelType}/{levelId}/policies | iam-policies-management, account-env-read |
| dynatrace_iam_policy_bindings | /iam/v1/repo/{levelType}/{levelId}/bindings/{policyUuid} | iam-policies-management, account-env-read |
| dynatrace_iam_policy_bindings_v2 | /iam/v1/repo/{levelType}/{levelId}/bindings/{policyUuid}{groupUuid} | iam-policies-management, account-env-read |
| dynatrace_iam_policy_boundary | /iam/v1/repo/account/{accountId}/boundaries | iam-policies-management, account-env-read |
| dynatrace_iam_user | /iam/v1/accounts/{accountUuid}/users | account-idm-read, account-idm-write |
| dynatrace_openpipeline_business_events | /platform/openpipeline/v1/configurations/bizevents | openpipeline:configurations:read, openpipeline:configurations:write |
| dynatrace_openpipeline_events | /platform/openpipeline/v1/configurations/events | openpipeline:configurations:read, openpipeline:configurations:write |
| dynatrace_openpipeline_logs | /platform/openpipeline/v1/configurations/logs | openpipeline:configurations:read, openpipeline:configurations:write |
| dynatrace_openpipeline_sdlc_events | /platform/openpipeline/v1/configurations/events.sdlc | openpipeline:configurations:read, openpipeline:configurations:write |
| dynatrace_openpipeline_security_events | /platform/openpipeline/v1/configurations/events.security | openpipeline:configurations:read, openpipeline:configurations:write |
| dynatrace_openpipeline_metrics | /platform/openpipeline/v1/configurations/metrics | openpipeline:configurations:read, openpipeline:configurations:write |
| dynatrace_openpipeline_user_sessions | /platform/openpipeline/v1/configurations/usersessions | openpipeline:configurations:read, openpipeline:configurations:write |
| dynatrace_openpipeline_davis_problems | /platform/openpipeline/v1/configurations/davis.problems | openpipeline:configurations:read, openpipeline:configurations:write |
| dynatrace_openpipeline_davis_events | /platform/openpipeline/v1/configurations/davis.events | openpipeline:configurations:read, openpipeline:configurations:write |
| dynatrace_openpipeline_system_events | /platform/openpipeline/v1/configurations/system.events | openpipeline:configurations:read, openpipeline:configurations:write |
| dynatrace_openpipeline_user_events | /platform/openpipeline/v1/configurations/user.events | openpipeline:configurations:read, openpipeline:configurations:write |
| dynatrace_openpipeline_spans | /platform/openpipeline/v1/configurations/spans | openpipeline:configurations:read, openpipeline:configurations:write |
| dynatrace_platform_bucket | /platform/storage/management/v1/bucket-definitions | storage:bucket-definitions:read, storage:bucket-definitions:write |
| dynatrace_platform_slo | /platform/slo/v1/slos | slo:slos:read, slo:slos:write |
| dynatrace_segment | /platform/storage/filter-segments/v1/filter-segments | storage:filter-segments:read, storage:filter-segments:write, storage:filter-segments:share, storage:filter-segments:delete, storage:filter-segments:admin |