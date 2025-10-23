/**
* @license
* Copyright 2020 Dynatrace LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package export

import (
	"os"
	"reflect"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/grail/segments"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/permissions"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/openpipeline"
	settingsPermissions "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/settings/objects/permissions"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"

	msentraidconnection "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/azure/connector/microsoftentraidentitydeveloperconnection"
	dbfeatureflags "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/database/featureflags"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/devobs/debugger/gitonprem"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/discovery/coverage/defaultrules"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/discovery/coverage/featureflags"
	githubconnection "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/github/connector/connection"
	gitlabconnection "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/gitlab/connector/connection"
	hubpermissions "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/hub/manage/permissions"
	infraopsfeatureflags "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/infraops/featureflags"
	infraopssettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/infraops/settings"
	jenkinsconnection "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/jenkins/connector/connection"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/jiraconnection"
	k8sautomationconnections "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/kubernetes/connector/connection"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/launchpad"
	ms365emailconnection "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/microsoft365/connector/mail/connection"
	msteamsconnection "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/msteams/connection"
	pagerdutyconnection "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/pagerduty/connection"
	automationcontroller "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/redhat/ansible/automationcontroller/connection"
	edawebhookconnection "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/redhat/ansible/edawebhook/connection"
	servicenowconnection "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/servicenow/connection"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/sitereliabilityguardian"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/slackconnection"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/automation/business_calendars"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/automation/scheduling_rules"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/automation/workflows"
	ddupool "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/accounting/ddulimit"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/activegatetoken"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/alerting/connectivityalerts"
	alerting "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/alerting/profile"
	database_anomalies_v2 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/anomalydetection/databases"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/anomalydetection/diskrules"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/anomalydetection/frequentissues"
	aws_anomalies "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/anomalydetection/infrastructure/aws"
	disk_anomalies_v2 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/anomalydetection/infrastructure/disks"
	disk_specific_anomalies_v2 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/anomalydetection/infrastructure/disks/perdiskoverride"
	host_anomalies_v2 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/anomalydetection/infrastructure/hosts"
	vmware_anomalies "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/anomalydetection/infrastructure/vmware"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/anomalydetection/kubernetes/cluster"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/anomalydetection/kubernetes/namespace"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/anomalydetection/kubernetes/node"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/anomalydetection/kubernetes/pvc"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/anomalydetection/kubernetes/workload"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/anomalydetection/metricevents"
	custom_app_anomalies "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/anomalydetection/rum/custom"
	custom_app_crash_rate "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/anomalydetection/rum/custom/crashrate"
	mobile_app_anomalies "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/anomalydetection/rum/mobile"
	mobile_app_crash_rate "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/anomalydetection/rum/mobile/crashrate"
	web_app_anomalies "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/anomalydetection/rum/web"
	service_anomalies_v2 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/anomalydetection/services"
	apidetection "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/apis/detectionrules"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/appengineregistry/clouddevelopmentenvironments"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/appsec/attackprotectionadvancedconfig"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/appsec/attackprotectionallowlistconfig"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/appsec/attackprotectionsettings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/appsec/codelevelvulnerability"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/appsec/notificationalertingprofile"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/appsec/notificationattackalertingprofile"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/appsec/notificationintegration"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/appsec/rulesettings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/appsec/runtimevulnerabilitydetection"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/appsec/thirdpartyvulnerabilitykuberneteslabelrulesettings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/appsec/thirdpartyvulnerabilityrulesettings"
	kubernetesapp "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/apptransition/kubernetes"
	attributeallowlist "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/attribute/allowlist"
	attributeblocklist "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/attribute/blocklist"
	attributemasking "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/attribute/masking"
	attributespreferences "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/attribute/preferences"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/auditlog"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/automation/approval"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/availability/processgroupalerting"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/bizevents/http/capturingvariants"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/bizevents/http/incoming"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/bizevents/http/outgoing"
	bizevents_buckets "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/bizevents/processing/buckets"
	bizevents_metrics "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/bizevents/processing/metrics"
	bizevents_processing "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/bizevents/processing/pipelines"
	bizevents_security "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/bizevents/security/contextrules"
	cloudfoundryv2 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/cloud/cloudfoundry"
	kubernetesv2 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/cloud/kubernetes"
	kubernetesmonitoring "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/cloud/kubernetes/monitoring"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/container/builtinmonitoringrule"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/container/monitoringrule"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/container/registry"
	containertechnology "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/container/technology"
	crashdumpanalytics "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/crashdump/analytics"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/custommetrics"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/customunit"
	dashboardsgeneral "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/dashboards/general"
	dashboardsallowlist "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/dashboards/image/allowlist"
	dashboardspresets "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/dashboards/presets"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/declarativegrouping"
	activegateupdates "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/deployment/activegate/updates"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/deployment/management/updatewindows"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/deployment/oneagent/defaultmode"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/deployment/oneagent/defaultversion"
	oneagentupdates "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/deployment/oneagent/updates"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/devobs/agentoptin"
	devobsmasking "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/devobs/sensitivedatamasking"
	diskanalytics "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/disk/analytics/extension"
	diskoptions "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/disk/options"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/dtjavascriptruntime/allowedoutboundconnections"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/dtjavascriptruntime/appmonitoring"
	ebpf "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/ebpf/service/discovery"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/eec/local"
	eecremote "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/eec/remote"
	endpointdetectionrules "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/endpointdetectionrules"
	endpointdetectionrulesoptin "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/endpointdetectionrules/optin"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/eulasettings"
	networktraffic "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/exclude/network/traffic"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/failuredetectionrulesets"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/generic"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/geosettings"
	grailmetricsallowall "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/grail/metrics/allowall"
	grailmetricsallowlist "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/grail/metrics/allowlist"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/histogrammetrics"
	hostmonitoring "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/host/monitoring"
	hostmonitoringadvanced "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/host/monitoring/advanced"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/host/monitoring/aixkernelextension"
	hostmonitoringmode "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/host/monitoring/mode"
	hostprocessgroupmonitoring "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/host/processgroups/monitoringstate"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hub/subscriptions"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/awsconnection"

	connections_aws "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/aws"
	connections_aws_role_arn "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/aws/role_arn"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/ibmmq/imsbridges"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/ibmmq/queuemanagers"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/ibmmq/queuesharinggroup"
	diskedgeanomalydetectors "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/infrastructure/diskedge/anomalydetectors"
	issuetracking "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/issuetracking/integration"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/kubernetes/enrichment"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/kubernetes/securityposturemanagement"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/customlogsourcesettings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/logagentconfiguration"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/logagentfeatureflags"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/logbucketsrules"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/logcustomattributes"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/logdebugsettings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/logdpprules"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/logevents"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/logsecuritycontextrules"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/logsongrailactivate"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/logstoragesettings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/schemalesslogmetric"
	sensitivedatamasking "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/sensitivedatamaskingsettings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/timestampconfiguration"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/mainframe/mqfilters"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/mainframe/txmonitoring"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/mainframe/txstartfilters"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/metric/metadata"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/metric/query"
	mobilenotifications "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/mobile/notifications"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/monitoredentities/generic/relation"
	generictypes "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/monitoredentities/generic/type"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/monitoredentities/grail/securitycontext"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/monitoredtechnologies/apache"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/monitoredtechnologies/dotnet"
	golang "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/monitoredtechnologies/go"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/monitoredtechnologies/iis"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/monitoredtechnologies/java"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/monitoredtechnologies/nginx"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/monitoredtechnologies/nodejs"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/monitoredtechnologies/opentracingnative"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/monitoredtechnologies/php"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/monitoredtechnologies/python"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/monitoredtechnologies/varnish"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/monitoredtechnologies/wsmb"
	slov2 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/monitoring/slo"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/monitoring/slo/normalization"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/nettracer/traffic"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/networkzones"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/oneagent/features"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/oneagent/masking"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/opentelemetrymetrics"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/osservicesmonitoring"
	ownership_config "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/ownership/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/ownership/teams"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/preferences/ipaddressmasking"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/preferences/privacy"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/problem/fields"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/problem/notifications"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/problem/notifications/ansible"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/problem/notifications/email"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/problem/notifications/jira"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/problem/notifications/opsgenie"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/problem/notifications/pagerduty"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/problem/notifications/servicenow"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/problem/notifications/slack"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/problem/notifications/trello"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/problem/notifications/victorops"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/problem/notifications/webhook"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/problem/notifications/xmatters"
	problemrecordpropagation "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/problem/record/propagation/rules"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/process/builtinprocessmonitoringrule"
	processmonitoring "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/process/monitoring"
	customprocessmonitoring "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/process/monitoring/custom"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/processavailability"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/processgroup/advanceddetectionrule"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/processgroup/cloudapplication/workloaddetection"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/processgroup/detectionflags"
	processgroupmonitoring "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/processgroup/monitoring/state"
	processgroupsimpledetection "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/processgroup/simpledetectionrule"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/processvisibility"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/remote/environment"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/resourceattribute"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rpcbasedsampling"
	rumcustomenablement "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/custom/enablement"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/hostheaders"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/ipdetermination"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/ipmappings"
	rummobileenablement "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/mobile/enablement"
	mobilekeyperformancemetrics "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/mobile/keyperformancemetrics"
	mobilerequesterrors "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/mobile/requesterrors"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/overloadprevention"
	rumprocessgroup "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/processgroup"
	rumproviderbreakdown "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/providerbreakdown"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/resourcetimingorigins"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/userexperiencescore"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/appdetection"
	webappautoinjection "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/automaticinjection"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/beacondomainorigins"
	webappbeaconendpoint "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/beaconendpoint"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/capturecustomproperties"
	webappcustomconfigproperties "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/customconfigurationproperties"
	webappcustomerrors "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/customerrors"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/custominjectionrules"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/customrumjavascriptversion"
	rumwebenablement "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/enablement"
	webappinjectioncookie "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/injection/cookie"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/ipaddressexclusion"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/keyperformancemetric/customactions"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/keyperformancemetric/loadactions"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/keyperformancemetric/xhractions"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/manualinsertion"
	webapprequesterrors "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/requesterrors"
	webappresourcecleanup "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/resourcecleanuprules"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/resourcetypes"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/rumjavascriptfilename"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/rumjavascriptupdates"
	securitycontextsettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/securitycontext"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/servicedetection/externalwebrequest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/servicedetection/externalwebservice"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/servicedetection/fullwebrequest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/servicedetection/fullwebservice"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/servicedetectionrules"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/servicesplittingrules"
	sessionreplaywebprivacy "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/sessionreplay/web/privacypreferences"
	sessionreplayresourcecapture "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/sessionreplay/web/resourcecapturing"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/settings/keyrequests"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/settings/mutedrequests"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/span/attribute"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/span/capturing"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/span/contextpropagation"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/span/entrypoints"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/span/eventattribute"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/synthetic/availability"
	browseroutagehandling "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/synthetic/browser/outagehandling"
	browserperformancethresholds "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/synthetic/browser/performancethresholds"
	httpcookies "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/synthetic/http/cookies"
	httpoutagehandling "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/synthetic/http/outagehandling"
	httpperformancethresholds "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/synthetic/http/performancethresholds"
	networkoutagehandling "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/synthetic/multiprotocol/outagehandling"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/tags/autotagging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/tokens/tokensettings"
	unifiedservicesopentel "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/unifiedservices/enablement"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/unifiedservices/endpointmetrics"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/urlbasedsampling"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/usability/analytics"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/useractioncustommetrics"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/usersettings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/virtualization/vmware"
	onprempolicybindings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/cluster/v1/bindings"
	onpremusergroups "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/cluster/v1/groups"
	onpremmgmzpermission "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/cluster/v1/permissions/mgmz"
	onprempolicies "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/cluster/v1/policies"
	onpremusers "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/cluster/v1/users"
	managednetworkzones "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/cluster/v2/networkzones"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/service/daviscopilot/dataminingblocklist"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/reports"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/documents/directshares"
	documents "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/documents/document"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/bindings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/boundaries"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/groups"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/policies"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/users"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/v2bindings"
	platformbuckets "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/platform/buckets"
	alertingv1 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/alerting"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/customtags"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/dashboards"
	maintenancev1 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/maintenance"
	managementzonesv1 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/managementzones"
	notificationsv1 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/notifications"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/requestnaming/order"
	locations "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/locations/private"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/activegatetokens"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/customdevice"
	active_version "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/hub/extension/active_version"
	extension_config "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/hub/extension/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/slo"
	v2monitors "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/synthetic/monitors"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/cache"

	v2managementzones "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/managementzones"

	application_anomalies "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/applications"
	database_anomalies "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/databaseservices"
	disk_event_anomalies "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/diskevents"
	host_anomalies "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/hosts"
	custom_anomalies "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/metricevents"
	pg_anomalies "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/processgroups"
	service_anomalies "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/services"

	host_naming "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/naming/hosts"
	host_naming_order "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/naming/hosts/order"
	processgroup_naming "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/naming/processgroups"
	processgroup_naming_order "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/naming/processgroups/order"
	service_naming "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/naming/services"
	service_naming_order "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/naming/services/order"
	networkzone "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/networkzones"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/davis/anomalydetectors"
	envparameters "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/failuredetection/environment/parameters"
	envrules "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/failuredetection/environment/rules"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/failuredetection/service/generalparameters"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/failuredetection/service/httpparameters"
	platformslo "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/slo"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/mobile"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web/dataprivacy"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web/detection"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web/errors"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web/keyuseractions"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/autotags"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/credentials/aws"
	aws_services "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/credentials/aws/services"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/credentials/azure"
	azure_services "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/credentials/azure/services"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/credentials/cloudfoundry"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/credentials/kubernetes"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/customservices"
	customservices_order "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/customservices/order"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/dashboards/sharing"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/jsondashboards"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/jsondashboardsbase"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/requestattributes"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/requestnaming"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/browser"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/http"
	httpscript "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/http/script"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/credentials/vault"

	v2maintenance "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/alerting/maintenancewindow"
	calculated_mobile_metrics "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/metrics/calculated/mobile"
	calculated_service_metrics "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/metrics/calculated/service"
	calculated_synthetic_metrics "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/metrics/calculated/synthetic"
	calculated_web_metrics "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/metrics/calculated/web"

	openpipelinebizeventsingestsources "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/bizevents/ingestsources"
	openpipelinebizeventspipelines "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/bizevents/pipelines"
	openpipelinebizeventsrouting "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/bizevents/routing"
	openpipelinedaviseventsingestsources "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/davis/events/ingestsources"
	openpipelinedaviseventspipelines "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/davis/events/pipelines"
	openpipelinedaviseventsrouting "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/davis/events/routing"
	openpipelinedavisproblemsingestsources "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/davis/problems/ingestsources"
	openpipelinedavisproblemspipelines "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/davis/problems/pipelines"
	openpipelinedavisproblemsrouting "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/davis/problems/routing"
	openpipelineeventsingestsources "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/events/ingestsources"
	openpipelineeventspipelines "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/events/pipelines"
	openpipelineeventsrouting "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/events/routing"
	openpipelineeventssdlcingestsources "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/events/sdlc/ingestsources"
	openpipelineeventssdlcpipelines "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/events/sdlc/pipelines"
	openpipelineeventssdlcrouting "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/events/sdlc/routing"
	openpipelineeventssecurityingestsources "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/events/security/ingestsources"
	openpipelineeventssecuritypipelines "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/events/security/pipelines"
	openpipelineeventssecurityrouting "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/events/security/routing"
	openpipelinelogsingestsources "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/logs/ingestsources"
	openpipelinelogspipelines "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/logs/pipelines"
	openpipelinelogsrouting "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/logs/routing"
	openpipelinemetricsingestsources "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/metrics/ingestsources"
	openpipelinemetricspipelines "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/metrics/pipelines"
	openpipelinemetricsrouting "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/metrics/routing"
	openpipelinesecurityeventsingestsources "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/security/events/ingestsources"
	openpipelinesecurityeventspipelines "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/security/events/pipelines"
	openpipelinesecurityeventsrouting "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/security/events/routing"
	openpipelinespansingestsources "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/spans/ingestsources"
	openpipelinespanspipelines "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/spans/pipelines"
	openpipelinespansrouting "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/spans/routing"
	openpipelinesystemeventsingestsources "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/system/events/ingestsources"
	openpipelinesystemeventspipelines "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/system/events/pipelines"
	openpipelinesystemeventsrouting "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/system/events/routing"
	openpipelineusereventsingestsources "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/user/events/ingestsources"
	openpipelineusereventspipelines "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/user/events/pipelines"
	openpipelineusereventsrouting "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/user/events/routing"
	openpipelineusersessionsingestsources "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/usersessions/ingestsources"
	openpipelineusersessionspipelines "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/usersessions/pipelines"
	openpipelineusersessionsrouting "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/usersessions/routing"
)

func NewResourceDescriptor[T settings.Settings](fn func(credentials *rest.Credentials) settings.CRUDService[T], dependencies ...Dependency) ResourceDescriptor {
	return ResourceDescriptor{
		Service: func(credentials *rest.Credentials) settings.CRUDService[settings.Settings] {
			return &settings.GenericCRUDService[T]{Service: cache.CRUD(fn(credentials))}
		},
		protoType:    newSettings(fn),
		Dependencies: dependencies,
	}
}

func NewResourceDescriptorWithFolderOverride[T settings.Settings](fn func(credentials *rest.Credentials) settings.CRUDService[T], folderName string, dependencies ...Dependency) ResourceDescriptor {
	return ResourceDescriptor{
		Service: func(credentials *rest.Credentials) settings.CRUDService[settings.Settings] {
			return &settings.GenericCRUDService[T]{Service: cache.CRUD(fn(credentials))}
		},
		protoType:    newSettings(fn),
		Dependencies: dependencies,
		FolderName:   folderName,
	}
}

func NewChildResourceDescriptor[T settings.Settings](fn func(credentials *rest.Credentials) settings.CRUDService[T], parent ResourceType, dependencies ...Dependency) ResourceDescriptor {
	return ResourceDescriptor{
		Service: func(credentials *rest.Credentials) settings.CRUDService[settings.Settings] {
			return &settings.GenericCRUDService[T]{Service: cache.CRUD(fn(credentials))}
		},
		protoType:    newSettings(fn),
		Dependencies: dependencies,
		Parent:       &parent,
	}
}

func newSettings[T settings.Settings](sfn func(credentials *rest.Credentials) settings.CRUDService[T]) T {
	var proto T
	return reflect.New(reflect.TypeOf(proto).Elem()).Interface().(T)
}

type ResourceDescriptor struct {
	Dependencies []Dependency
	Service      func(credentials *rest.Credentials) settings.CRUDService[settings.Settings]
	protoType    settings.Settings
	except       func(id string, name string) bool
	Parent       *ResourceType
	FolderName   string
}

func (me ResourceDescriptor) Specify(t notifications.Type) ResourceDescriptor {
	if notification, ok := me.protoType.(*notifications.Notification); ok {
		notification.Type = t
	}
	return me
}

func (me ResourceDescriptor) Except(except func(id string, name string) bool) ResourceDescriptor {
	me.except = except
	return me
}

func (me ResourceDescriptor) NewSettings() settings.Settings {
	res := reflect.New(reflect.TypeOf(me.protoType).Elem()).Interface().(settings.Settings)
	if notification, ok := res.(*notifications.Notification); ok {
		notification.Type = me.protoType.(*notifications.Notification).Type
	}
	return res
}

func (me ResourceDescriptor) HasWeakIDDependencyTo(resType ResourceType) bool {
	for _, dependency := range me.Dependencies {
		if iddependency, ok := dependency.(*iddep); ok {
			if iddependency.resourceType != resType {
				continue
			}
			if !iddependency.onlyNonPostProcessed {
				continue
			}
			return true
		}
	}
	return false
}

func IsSettings20Schema(schemaID string) bool {
	return strings.HasPrefix(schemaID, "builtin:") || strings.HasPrefix(schemaID, "app:")
}

func ContainsInsertAfterAttribute(protoType settings.Settings, schemaID string) bool {
	return IsSettings20Schema(schemaID) && settings.HasInsertAfter(protoType)
}

// AddInsertAfterWeakIDDependencies completes the configured
// Resource Descriptors for every resource that is orderable
// using the mechanism Settings 2.0.
// The `insert_after` attribute allows for ordering settings.
// The export functionality needs to know that these kinds of
// settings are allowed to contain IDs to the same resource type
// in order to replace hardcoded IDs in there.
// `Dependencies.WeakID` takes care of that.
func AddInsertAfterWeakIDDependencies() {
	for resType, descriptor := range AllResources {
		schemaID := descriptor.Service(&rest.Credentials{}).SchemaID()
		if !ContainsInsertAfterAttribute(descriptor.protoType, schemaID) {
			continue
		}
		if descriptor.HasWeakIDDependencyTo(resType) {
			continue
		}
		descriptor.Dependencies = append(descriptor.Dependencies, Dependencies.WeakID(resType))
	}
}

var AllResources = map[ResourceType]ResourceDescriptor{
	ResourceTypes.Alerting: NewResourceDescriptor(
		alerting.Service,
		Dependencies.LegacyID(ResourceTypes.ManagementZoneV2),
		Dependencies.ID(ResourceTypes.ManagementZoneV2),
	),
	ResourceTypes.AnsibleTowerNotification: NewResourceDescriptor(
		ansible.Service,
		Dependencies.ID(ResourceTypes.Alerting),
	).Specify(notifications.Types.AnsibleTower),
	ResourceTypes.ApplicationAnomalies: NewResourceDescriptor(application_anomalies.Service),
	ResourceTypes.ApplicationDataPrivacy: NewResourceDescriptor(
		dataprivacy.Service,
		Dependencies.ID(ResourceTypes.WebApplication),
	),
	ResourceTypes.ApplicationDetection: NewResourceDescriptor(
		detection.Service,
		Dependencies.ID(ResourceTypes.WebApplication),
	),
	ResourceTypes.ApplicationErrorRules: NewResourceDescriptor(
		errors.Service,
		Dependencies.ID(ResourceTypes.WebApplication),
	),
	ResourceTypes.AutoTag: NewResourceDescriptor(
		autotags.Service,
		Coalesce(Dependencies.Service),
		Coalesce(Dependencies.Host),
		Coalesce(Dependencies.HostGroup),
		Coalesce(Dependencies.ProcessGroup),
		Coalesce(Dependencies.ProcessGroupInstance),
	),
	ResourceTypes.AWSCredentials: NewResourceDescriptor(
		aws.Service,
		Dependencies.ID(ResourceTypes.AWSCredentials),
	),
	ResourceTypes.AWSService: NewResourceDescriptor(
		aws_services.Service,
		Dependencies.ID(ResourceTypes.AWSCredentials),
	),
	ResourceTypes.AzureCredentials: NewResourceDescriptor(azure.Service),
	ResourceTypes.AzureService: NewResourceDescriptor(
		azure_services.Service,
		Dependencies.ID(ResourceTypes.AzureCredentials),
	),
	ResourceTypes.BrowserMonitor: NewResourceDescriptor(
		browser.Service,
		Dependencies.ID(ResourceTypes.SyntheticLocation),
		Dependencies.ID(ResourceTypes.WebApplication),
		Dependencies.ID(ResourceTypes.Credentials),
	).Except(func(id string, name string) bool {
		return strings.HasPrefix(name, "Monitor synchronizing credentials with")
	}),
	ResourceTypes.CalculatedServiceMetric: NewResourceDescriptor(
		calculated_service_metrics.Service,
		Dependencies.ManagementZone,
		Dependencies.RequestAttribute,
		Dependencies.Service,
		Dependencies.Host,
		Dependencies.HostGroup,
		Dependencies.ProcessGroup,
		Dependencies.ProcessGroupInstance,
	),
	ResourceTypes.CalculatedWebMetric: NewResourceDescriptor(
		calculated_web_metrics.Service,
		Dependencies.ID(ResourceTypes.WebApplication),
	),
	ResourceTypes.CalculatedMobileMetric: NewResourceDescriptor(
		calculated_mobile_metrics.Service,
		Dependencies.ID(ResourceTypes.MobileApplication),
	),
	ResourceTypes.CalculatedSyntheticMetric: NewResourceDescriptor(
		calculated_synthetic_metrics.Service,
		Dependencies.ID(ResourceTypes.BrowserMonitor),
		Dependencies.ID(ResourceTypes.HTTPMonitor),
	),
	ResourceTypes.CloudFoundryCredentials: NewResourceDescriptor(cloudfoundry.Service),
	ResourceTypes.CustomAnomalies: NewResourceDescriptor(
		custom_anomalies.Service,
		Dependencies.LegacyID(ResourceTypes.ManagementZoneV2),
		Dependencies.ID(ResourceTypes.ManagementZoneV2),
	).Except(func(id string, name string) bool {
		return strings.HasPrefix(id, "builtin:") || strings.HasPrefix(id, "ruxit.") || strings.HasPrefix(id, "dynatrace.") || strings.HasPrefix(id, "custom.remote.python.") || strings.HasPrefix(id, "custom.python.")
	}),
	ResourceTypes.CustomAppAnomalies: NewResourceDescriptor(
		custom_app_anomalies.Service,
		Dependencies.ID(ResourceTypes.MobileApplication),
		Coalesce(Dependencies.DeviceApplicationMethod),
	),
	ResourceTypes.CustomAppCrashRate: NewResourceDescriptor(
		custom_app_crash_rate.Service,
		Dependencies.ID(ResourceTypes.MobileApplication),
	),
	ResourceTypes.MobileAppAnomalies: NewResourceDescriptor(
		mobile_app_anomalies.Service,
		Coalesce(Dependencies.DeviceApplicationMethod),
		Dependencies.ID(ResourceTypes.MobileApplication),
	),
	ResourceTypes.MobileAppCrashRate: NewResourceDescriptor(
		mobile_app_crash_rate.Service,
		Dependencies.ID(ResourceTypes.MobileApplication),
	),
	ResourceTypes.WebAppAnomalies: NewResourceDescriptor(
		web_app_anomalies.Service,
		Dependencies.ID(ResourceTypes.WebApplication),
		Coalesce(Dependencies.ApplicationMethod),
	),
	ResourceTypes.CustomService: NewResourceDescriptor(customservices.Service),
	ResourceTypes.CustomServiceOrder: NewResourceDescriptor(
		customservices_order.Service,
		Dependencies.ID(ResourceTypes.CustomService),
	),
	ResourceTypes.Credentials: NewResourceDescriptor(
		vault.Service,
		Dependencies.ID(ResourceTypes.Credentials),
	),
	ResourceTypes.JSONDashboardBase: NewResourceDescriptorWithFolderOverride(
		jsondashboardsbase.Service,
		"json_dashboard",
	),
	ResourceTypes.Documents: NewResourceDescriptor(
		documents.Service,
	),
	ResourceTypes.DirectShares: NewResourceDescriptor(
		directshares.Service,
	),
	ResourceTypes.OpenPipelineLogs: NewResourceDescriptor(
		openpipeline.LogsService, Dependencies.ID(ResourceTypes.PlatformBucket)),
	ResourceTypes.OpenPipelineEvents: NewResourceDescriptor(
		openpipeline.EventsService, Dependencies.ID(ResourceTypes.PlatformBucket)),
	ResourceTypes.OpenPipelineSecurityEvents: NewResourceDescriptor(
		openpipeline.SecurityEventsService, Dependencies.ID(ResourceTypes.PlatformBucket)),
	ResourceTypes.OpenPipelineBusinessEvents: NewResourceDescriptor(
		openpipeline.BusinessEventsService, Dependencies.ID(ResourceTypes.PlatformBucket)),
	ResourceTypes.OpenPipelineSDLCEvents: NewResourceDescriptor(
		openpipeline.SDLCEventsService, Dependencies.ID(ResourceTypes.PlatformBucket)),
	ResourceTypes.OpenPipelineMetrics: NewResourceDescriptor(
		openpipeline.MetricsService, Dependencies.ID(ResourceTypes.PlatformBucket)),
	ResourceTypes.OpenPipelineUserSessions: NewResourceDescriptor(
		openpipeline.UserSessionsService, Dependencies.ID(ResourceTypes.PlatformBucket)),
	ResourceTypes.OpenPipelineDavisProblems: NewResourceDescriptor(
		openpipeline.DavisProblemsService, Dependencies.ID(ResourceTypes.PlatformBucket)),
	ResourceTypes.OpenPipelineDavisEvents: NewResourceDescriptor(
		openpipeline.DavisEventsService, Dependencies.ID(ResourceTypes.PlatformBucket)),
	ResourceTypes.OpenPipelineSystemEvents: NewResourceDescriptor(
		openpipeline.SystemEventsService, Dependencies.ID(ResourceTypes.PlatformBucket)),
	ResourceTypes.OpenPipelineUserEvents: NewResourceDescriptor(
		openpipeline.UserEventsService, Dependencies.ID(ResourceTypes.PlatformBucket)),
	ResourceTypes.OpenPipelineSpans: NewResourceDescriptor(
		openpipeline.SpansService, Dependencies.ID(ResourceTypes.PlatformBucket)),
	ResourceTypes.JSONDashboard: NewChildResourceDescriptor(
		jsondashboards.Service,
		ResourceTypes.JSONDashboardBase,
		Dependencies.DashboardLinkID(true),
		Dependencies.LegacyID(ResourceTypes.ManagementZoneV2),
		Dependencies.ID(ResourceTypes.ManagementZoneV2),
		Dependencies.ManagementZone,
		// Dependencies.Service,
		Dependencies.LegacyID(ResourceTypes.SLOV2),
		Dependencies.HyperLinkDashboardID(),
		Dependencies.ID(ResourceTypes.WebApplication),
		Dependencies.ID(ResourceTypes.MobileApplication),
		Dependencies.ID(ResourceTypes.SyntheticLocation),
		Dependencies.ID(ResourceTypes.HTTPMonitor),
		Dependencies.ID(ResourceTypes.CalculatedServiceMetric),
		Dependencies.ID(ResourceTypes.CalculatedWebMetric),
		Dependencies.ID(ResourceTypes.CalculatedMobileMetric),
		Dependencies.ID(ResourceTypes.CalculatedSyntheticMetric),
		Dependencies.ID(ResourceTypes.BrowserMonitor),
	),
	ResourceTypes.DashboardSharing: NewChildResourceDescriptor(
		sharing.Service,
		ResourceTypes.JSONDashboardBase,
		Dependencies.ResourceID(ResourceTypes.JSONDashboardBase, true),
		Dependencies.LegacyID(ResourceTypes.ManagementZoneV2),
		Dependencies.ID(ResourceTypes.ManagementZoneV2),
	),
	ResourceTypes.DatabaseAnomalies:  NewResourceDescriptor(database_anomalies.Service),
	ResourceTypes.DiskEventAnomalies: NewResourceDescriptor(disk_event_anomalies.Service),
	ResourceTypes.DiskAnomaliesV2: NewResourceDescriptor(
		disk_anomalies_v2.Service,
		Coalesce(Dependencies.Host),
	),
	ResourceTypes.DiskSpecificAnomaliesV2: NewResourceDescriptor(
		disk_specific_anomalies_v2.Service,
		Coalesce(Dependencies.Disk),
	),
	ResourceTypes.EmailNotification: NewResourceDescriptor(
		email.Service,
		Dependencies.ID(ResourceTypes.Alerting),
	).Specify(notifications.Types.Email),
	ResourceTypes.FrequentIssues: NewResourceDescriptor(frequentissues.Service),
	ResourceTypes.HostAnomalies:  NewResourceDescriptor(host_anomalies.Service),
	ResourceTypes.HostAnomaliesV2: NewResourceDescriptor(
		host_anomalies_v2.Service,
		Coalesce(Dependencies.Host),
		Coalesce(Dependencies.HostGroup),
	),
	ResourceTypes.HTTPMonitor: NewResourceDescriptor(
		http.Service,
		Dependencies.ID(ResourceTypes.SyntheticLocation),
		Dependencies.ID(ResourceTypes.WebApplication),
		Dependencies.ID(ResourceTypes.Credentials),
	),
	ResourceTypes.HostNaming: NewResourceDescriptor(host_naming.Service),
	ResourceTypes.HostNamingOrder: NewResourceDescriptor(
		host_naming_order.Service,
		Dependencies.ID(ResourceTypes.HostNaming),
	),
	ResourceTypes.IBMMQFilters: NewResourceDescriptor(mqfilters.Service),
	ResourceTypes.IMSBridge:    NewResourceDescriptor(imsbridges.Service),
	ResourceTypes.JiraNotification: NewResourceDescriptor(
		jira.Service,
		Dependencies.ID(ResourceTypes.Alerting),
	).Specify(notifications.Types.Jira),
	ResourceTypes.KeyRequests: NewResourceDescriptor(
		keyrequests.Service,
		Coalesce(Dependencies.Service),
		Coalesce(Dependencies.Host),
		Coalesce(Dependencies.HostGroup),
		Coalesce(Dependencies.ProcessGroup),
		Coalesce(Dependencies.ProcessGroupInstance),
	),
	ResourceTypes.KubernetesCredentials: NewResourceDescriptor(kubernetes.Service),
	ResourceTypes.Maintenance: NewResourceDescriptor(
		v2maintenance.Service,
		Dependencies.LegacyID(ResourceTypes.ManagementZoneV2),
		Dependencies.ID(ResourceTypes.ManagementZoneV2),
	),
	ResourceTypes.ManagementZoneV2: NewResourceDescriptor(v2managementzones.Service),
	ResourceTypes.MetricEvents: NewResourceDescriptor(
		metricevents.Service,
		Dependencies.ManagementZone,
		Dependencies.LegacyID(ResourceTypes.ManagementZoneV2),
		Dependencies.ID(ResourceTypes.ManagementZoneV2),
	),
	ResourceTypes.MobileApplication: NewResourceDescriptor(
		mobile.Service,
		Dependencies.ID(ResourceTypes.RequestAttribute),
	),
	ResourceTypes.MutedRequests: NewResourceDescriptor(
		mutedrequests.Service,
		Coalesce(Dependencies.Service),
	),
	ResourceTypes.NetworkZone:  NewResourceDescriptor(networkzone.Service),
	ResourceTypes.NetworkZones: NewResourceDescriptor(networkzones.Service),
	ResourceTypes.OpsGenieNotification: NewResourceDescriptor(
		opsgenie.Service,
		Dependencies.ID(ResourceTypes.Alerting),
	).Specify(notifications.Types.OpsGenie),
	ResourceTypes.PagerDutyNotification: NewResourceDescriptor(
		pagerduty.Service,
		Dependencies.ID(ResourceTypes.Alerting),
	).Specify(notifications.Types.PagerDuty),
	ResourceTypes.ProcessGroupNaming: NewResourceDescriptor(processgroup_naming.Service),
	ResourceTypes.ProcessGroupNamingOrder: NewResourceDescriptor(
		processgroup_naming_order.Service,
		Dependencies.ID(ResourceTypes.ProcessGroupNaming),
	),
	ResourceTypes.QueueManager: NewResourceDescriptor(queuemanagers.Service),
	ResourceTypes.RequestAttribute: NewResourceDescriptor(
		requestattributes.Service,
		Coalesce(Dependencies.Host),
		Coalesce(Dependencies.HostGroup),
		Coalesce(Dependencies.ProcessGroup),
		Coalesce(Dependencies.ProcessGroupInstance),
		Coalesce(Dependencies.Service),
	),
	ResourceTypes.RequestNaming: NewResourceDescriptor(
		requestnaming.Service,
		Dependencies.RequestAttribute,
	),
	ResourceTypes.ResourceAttributes: NewResourceDescriptor(resourceattribute.Service),
	ResourceTypes.ServiceAnomalies:   NewResourceDescriptor(service_anomalies.Service),
	ResourceTypes.ServiceAnomaliesV2: NewResourceDescriptor(
		service_anomalies_v2.Service,
		Coalesce(Dependencies.ServiceMethod),
		Coalesce(Dependencies.Service),
		Coalesce(Dependencies.HostGroup),
	),
	ResourceTypes.ServiceNaming: NewResourceDescriptor(service_naming.Service),
	ResourceTypes.ServiceNamingOrder: NewResourceDescriptor(
		service_naming_order.Service,
		Dependencies.ID(ResourceTypes.ServiceNaming),
	),
	ResourceTypes.ServiceNowNotification: NewResourceDescriptor(
		servicenow.Service,
		Dependencies.ID(ResourceTypes.Alerting),
	).Specify(notifications.Types.ServiceNow),
	ResourceTypes.SlackNotification: NewResourceDescriptor(
		slack.Service,
		Dependencies.ID(ResourceTypes.Alerting),
	).Specify(notifications.Types.Slack),
	ResourceTypes.SLO: NewResourceDescriptor(
		slo.Service,
		Dependencies.ManagementZone,
		Dependencies.LegacyID(ResourceTypes.ManagementZoneV2),
		Dependencies.ID(ResourceTypes.ManagementZoneV2),
		Dependencies.ID(ResourceTypes.CalculatedServiceMetric),
		Dependencies.ID(ResourceTypes.CalculatedWebMetric),
		Dependencies.ID(ResourceTypes.CalculatedMobileMetric),
		Dependencies.ID(ResourceTypes.CalculatedSyntheticMetric),
	),
	ResourceTypes.SpanAttribute:          NewResourceDescriptor(attribute.Service),
	ResourceTypes.SpanCaptureRule:        NewResourceDescriptor(capturing.Service),
	ResourceTypes.SpanContextPropagation: NewResourceDescriptor(contextpropagation.Service),
	ResourceTypes.SpanEntryPoint:         NewResourceDescriptor(entrypoints.Service),
	ResourceTypes.SyntheticLocation:      NewResourceDescriptor(locations.Service),
	ResourceTypes.TrelloNotification: NewResourceDescriptor(
		trello.Service,
		Dependencies.ID(ResourceTypes.Alerting),
	).Specify(notifications.Types.Trello),
	ResourceTypes.VictorOpsNotification: NewResourceDescriptor(
		victorops.Service,
		Dependencies.ID(ResourceTypes.Alerting),
	).Specify(notifications.Types.VictorOps),
	ResourceTypes.WebApplication: NewResourceDescriptor(
		web.Service,
		Dependencies.ID(ResourceTypes.RequestAttribute),
	),
	ResourceTypes.WebHookNotification: NewResourceDescriptor(
		webhook.Service,
		Dependencies.ID(ResourceTypes.Alerting),
	).Specify(notifications.Types.WebHook),
	ResourceTypes.XMattersNotification: NewResourceDescriptor(
		xmatters.Service,
		Dependencies.ID(ResourceTypes.Alerting),
	).Specify(notifications.Types.XMatters),

	ResourceTypes.MaintenanceWindow: NewResourceDescriptor(
		maintenancev1.Service,
		Dependencies.LegacyID(ResourceTypes.ManagementZoneV2),
		Dependencies.ID(ResourceTypes.ManagementZoneV2),
	),
	ResourceTypes.ManagementZone: NewResourceDescriptor(managementzonesv1.Service),
	ResourceTypes.Dashboard: NewResourceDescriptor(
		dashboards.Service,
		Dependencies.LegacyID(ResourceTypes.ManagementZoneV2),
		Dependencies.ID(ResourceTypes.ManagementZoneV2),
		Dependencies.ManagementZone,
		Dependencies.ID(ResourceTypes.SLO),
		Dependencies.ID(ResourceTypes.WebApplication),
		Dependencies.ID(ResourceTypes.SyntheticLocation),
	),
	ResourceTypes.Notification: NewResourceDescriptor(
		notificationsv1.Service,
		Dependencies.LegacyID(ResourceTypes.Alerting),
	),
	ResourceTypes.QueueSharingGroups: NewResourceDescriptor(queuesharinggroup.Service),
	ResourceTypes.AlertingProfile: NewResourceDescriptor(
		alertingv1.Service,
		Dependencies.LegacyID(ResourceTypes.ManagementZoneV2),
		Dependencies.ID(ResourceTypes.ManagementZoneV2),
	),
	ResourceTypes.RequestNamings: NewResourceDescriptor(
		order.Service,
		Dependencies.ID(ResourceTypes.RequestNaming),
	),
	ResourceTypes.IAMUser: NewResourceDescriptor(
		users.Service,
		Dependencies.ID(ResourceTypes.IAMGroup),
	),
	ResourceTypes.IAMGroup: NewResourceDescriptor(
		groups.Service,
		Dependencies.LegacyID(ResourceTypes.ManagementZoneV2),
		Dependencies.ID(ResourceTypes.ManagementZoneV2),
		Dependencies.ID(ResourceTypes.IAMPermission),
		Dependencies.Tenant,
	),
	ResourceTypes.IAMPermission:     NewResourceDescriptor(permissions.Service),
	ResourceTypes.IAMPolicy:         NewResourceDescriptor(policies.Service),
	ResourceTypes.IAMPolicyBindings: NewResourceDescriptor(bindings.Service),
	ResourceTypes.IAMPolicyBindingsV2: NewResourceDescriptor(
		v2bindings.Service,
		Dependencies.ID(ResourceTypes.IAMGroup),
		Dependencies.ID(ResourceTypes.IAMPolicy),
		Dependencies.ID(ResourceTypes.IAMPolicyBoundary),
		Dependencies.GlobalPolicy,
	),
	ResourceTypes.IAMPolicyBoundary: NewResourceDescriptor(boundaries.Service),
	ResourceTypes.DDUPool:           NewResourceDescriptor(ddupool.Service),
	ResourceTypes.ProcessGroupAnomalies: NewResourceDescriptor(
		pg_anomalies.Service,
		Coalesce(Dependencies.ProcessGroup),
		Coalesce(Dependencies.ProcessGroupInstance),
	),
	ResourceTypes.ProcessGroupAlerting: NewResourceDescriptor(
		processgroupalerting.Service,
		Coalesce(Dependencies.ProcessGroup),
	),
	ResourceTypes.DatabaseAnomaliesV2: NewResourceDescriptor(
		database_anomalies_v2.Service,
		Coalesce(Dependencies.ServiceMethod),
		Coalesce(Dependencies.Service),
		Coalesce(Dependencies.HostGroup),
	),
	ResourceTypes.ProcessMonitoringRule: NewResourceDescriptor(
		customprocessmonitoring.Service,
		Coalesce(Dependencies.HostGroup),
	),
	ResourceTypes.ProcessMonitoring: NewResourceDescriptor(
		processmonitoring.Service,
		Coalesce(Dependencies.HostGroup),
	),
	ResourceTypes.ProcessAvailability: NewResourceDescriptor(
		processavailability.Service,
		Coalesce(Dependencies.HostGroup),
		Coalesce(Dependencies.Host),
	),
	ResourceTypes.AdvancedProcessGroupDetectionRule: NewResourceDescriptor(advanceddetectionrule.Service),
	ResourceTypes.ConnectivityAlerts: NewResourceDescriptor(
		connectivityalerts.Service,
		Coalesce(Dependencies.ProcessGroup),
	),
	ResourceTypes.DeclarativeGrouping: NewResourceDescriptor(
		declarativegrouping.Service,
		Coalesce(Dependencies.Host),
		Coalesce(Dependencies.HostGroup),
		Coalesce(Dependencies.K8sCluster),
	),
	ResourceTypes.HostMonitoring: NewResourceDescriptor(
		hostmonitoring.Service,
		Coalesce(Dependencies.Host),
	),
	ResourceTypes.HostProcessGroupMonitoring: NewResourceDescriptor(
		hostprocessgroupmonitoring.Service,
		Coalesce(Dependencies.Host),
		Coalesce(Dependencies.ProcessGroup),
	),
	ResourceTypes.RUMIPLocations: NewResourceDescriptor(ipmappings.Service),
	ResourceTypes.CustomAppEnablement: NewResourceDescriptor(
		rumcustomenablement.Service,
		Dependencies.ID(ResourceTypes.MobileApplication),
	),
	ResourceTypes.MobileAppEnablement: NewResourceDescriptor(
		rummobileenablement.Service,
		Dependencies.ID(ResourceTypes.MobileApplication),
	),
	ResourceTypes.WebAppEnablement: NewResourceDescriptor(
		rumwebenablement.Service,
		Dependencies.ID(ResourceTypes.WebApplication),
	),
	ResourceTypes.RUMProcessGroup: NewResourceDescriptor(
		rumprocessgroup.Service,
		Coalesce(Dependencies.ProcessGroup),
	),
	ResourceTypes.RUMProviderBreakdown:  NewResourceDescriptor(rumproviderbreakdown.Service),
	ResourceTypes.UserExperienceScore:   NewResourceDescriptor(userexperiencescore.Service),
	ResourceTypes.WebAppResourceCleanup: NewResourceDescriptor(webappresourcecleanup.Service),
	ResourceTypes.UpdateWindows:         NewResourceDescriptor(updatewindows.Service),
	ResourceTypes.ProcessGroupDetectionFlags: NewResourceDescriptor(
		detectionflags.Service,
		Coalesce(Dependencies.Host),
		Coalesce(Dependencies.HostGroup),
	),
	ResourceTypes.ProcessGroupMonitoring: NewResourceDescriptor(
		processgroupmonitoring.Service,
		Coalesce(Dependencies.ProcessGroup),
	),
	ResourceTypes.ProcessGroupSimpleDetection: NewResourceDescriptor(processgroupsimpledetection.Service),
	ResourceTypes.LogMetrics:                  NewResourceDescriptor(schemalesslogmetric.Service),
	ResourceTypes.BrowserMonitorPerformanceThresholds: NewResourceDescriptor(
		browserperformancethresholds.Service,
		Dependencies.ID(ResourceTypes.BrowserMonitor),
	),
	ResourceTypes.HttpMonitorPerformanceThresholds: NewResourceDescriptor(
		httpperformancethresholds.Service,
		Dependencies.ID(ResourceTypes.HTTPMonitor),
	),
	ResourceTypes.HttpMonitorCookies: NewResourceDescriptor(
		httpcookies.Service,
		Dependencies.ID(ResourceTypes.HTTPMonitor),
	),
	ResourceTypes.SessionReplayWebPrivacy: NewResourceDescriptor(
		sessionreplaywebprivacy.Service,
		Dependencies.ID(ResourceTypes.WebApplication),
	),
	ResourceTypes.SessionReplayResourceCapture: NewResourceDescriptor(
		sessionreplayresourcecapture.Service,
		Dependencies.ID(ResourceTypes.WebApplication),
	),
	ResourceTypes.UsabilityAnalytics: NewResourceDescriptor(
		analytics.Service,
		Dependencies.ID(ResourceTypes.WebApplication),
	),
	ResourceTypes.SyntheticAvailability: NewResourceDescriptor(availability.Service),
	ResourceTypes.BrowserMonitorOutageHandling: NewResourceDescriptor(
		browseroutagehandling.Service,
		Dependencies.ID(ResourceTypes.BrowserMonitor),
	),
	ResourceTypes.HttpMonitorOutageHandling: NewResourceDescriptor(
		httpoutagehandling.Service,
		Dependencies.ID(ResourceTypes.HTTPMonitor),
	),
	ResourceTypes.CloudAppWorkloadDetection:      NewResourceDescriptor(workloaddetection.Service),
	ResourceTypes.MainframeTransactionMonitoring: NewResourceDescriptor(txmonitoring.Service),
	ResourceTypes.MonitoredTechnologiesApache: NewResourceDescriptor(
		apache.Service,
		Coalesce(Dependencies.Host),
	),
	ResourceTypes.MonitoredTechnologiesDotNet: NewResourceDescriptor(
		dotnet.Service,
		Coalesce(Dependencies.Host),
	),
	ResourceTypes.MonitoredTechnologiesGo: NewResourceDescriptor(
		golang.Service,
		Coalesce(Dependencies.Host),
	),
	ResourceTypes.MonitoredTechnologiesIIS: NewResourceDescriptor(
		iis.Service,
		Coalesce(Dependencies.Host),
	),
	ResourceTypes.MonitoredTechnologiesJava: NewResourceDescriptor(
		java.Service,
		Coalesce(Dependencies.Host),
	),
	ResourceTypes.MonitoredTechnologiesNGINX: NewResourceDescriptor(
		nginx.Service,
		Coalesce(Dependencies.Host),
	),
	ResourceTypes.MonitoredTechnologiesNodeJS: NewResourceDescriptor(
		nodejs.Service,
		Coalesce(Dependencies.Host),
	),
	ResourceTypes.MonitoredTechnologiesOpenTracing: NewResourceDescriptor(
		opentracingnative.Service,
		Coalesce(Dependencies.Host),
	),
	ResourceTypes.MonitoredTechnologiesPHP: NewResourceDescriptor(
		php.Service,
		Coalesce(Dependencies.Host),
	),
	ResourceTypes.MonitoredTechnologiesPython: NewResourceDescriptor(
		python.Service,
		Coalesce(Dependencies.Host),
	),
	ResourceTypes.MonitoredTechnologiesVarnish: NewResourceDescriptor(
		varnish.Service,
		Coalesce(Dependencies.Host),
	),
	ResourceTypes.MonitoredTechnologiesWSMB: NewResourceDescriptor(
		wsmb.Service,
		Coalesce(Dependencies.Host),
	),
	ResourceTypes.ProcessVisibility: NewResourceDescriptor(
		processvisibility.Service,
		Coalesce(Dependencies.Host),
		Coalesce(Dependencies.HostGroup),
	),
	ResourceTypes.RUMHostHeaders:     NewResourceDescriptor(hostheaders.Service),
	ResourceTypes.RUMIPDetermination: NewResourceDescriptor(ipdetermination.Service),
	ResourceTypes.MobileAppRequestErrors: NewResourceDescriptor(
		mobilerequesterrors.Service,
		Dependencies.ID(ResourceTypes.MobileApplication),
	),
	ResourceTypes.TransactionStartFilters: NewResourceDescriptor(txstartfilters.Service),
	ResourceTypes.OneAgentFeatures: NewResourceDescriptor(
		features.Service,
		Coalesce(Dependencies.ProcessGroup),
		Coalesce(Dependencies.ProcessGroupInstance),
		Coalesce(Dependencies.CloudApplication),
		Coalesce(Dependencies.CloudApplicationNamespace),
		Coalesce(Dependencies.K8sCluster),
	),
	ResourceTypes.RUMOverloadPrevention:  NewResourceDescriptor(overloadprevention.Service),
	ResourceTypes.RUMAdvancedCorrelation: NewResourceDescriptor(resourcetimingorigins.Service),
	ResourceTypes.WebAppBeaconOrigins:    NewResourceDescriptor(beacondomainorigins.Service),
	ResourceTypes.WebAppResourceTypes:    NewResourceDescriptor(resourcetypes.Service),
	ResourceTypes.GenericTypes:           NewResourceDescriptor(generictypes.Service),
	ResourceTypes.GenericRelationships:   NewResourceDescriptor(relation.Service),
	ResourceTypes.SLONormalization:       NewResourceDescriptor(normalization.Service),
	ResourceTypes.DataPrivacy: NewResourceDescriptor(
		privacy.Service,
		Dependencies.ID(ResourceTypes.WebApplication),
	),
	ResourceTypes.ServiceFailure: NewResourceDescriptor(
		generalparameters.Service,
		Coalesce(Dependencies.Service),
		Dependencies.ID(ResourceTypes.RequestAttribute),
	),
	ResourceTypes.ServiceHTTPFailure: NewResourceDescriptor(
		httpparameters.Service,
		Coalesce(Dependencies.Service),
	),
	ResourceTypes.DiskOptions: NewResourceDescriptor(
		diskoptions.Service,
		Coalesce(Dependencies.Host),
		Coalesce(Dependencies.HostGroup),
	),
	ResourceTypes.OSServices: NewResourceDescriptor(
		osservicesmonitoring.Service,
		Coalesce(Dependencies.Host),
		Coalesce(Dependencies.HostGroup),
	),
	ResourceTypes.ExtensionExecutionController: NewResourceDescriptor(
		local.Service,
		Coalesce(Dependencies.Host),
		Coalesce(Dependencies.HostGroup),
	),
	ResourceTypes.NetTracerTraffic: NewResourceDescriptor(
		traffic.Service,
		Coalesce(Dependencies.Host),
		Coalesce(Dependencies.HostGroup),
	),
	ResourceTypes.AIXExtension: NewResourceDescriptor(
		aixkernelextension.Service,
		Coalesce(Dependencies.Host),
	),
	ResourceTypes.MetricMetadata:  NewResourceDescriptor(metadata.Service),
	ResourceTypes.MetricQuery:     NewResourceDescriptor(query.Service),
	ResourceTypes.ActiveGateToken: NewResourceDescriptor(activegatetoken.Service),
	ResourceTypes.AGToken:         NewResourceDescriptor(activegatetokens.Service),
	ResourceTypes.AuditLog:        NewResourceDescriptor(auditlog.Service),
	ResourceTypes.K8sClusterAnomalies: NewResourceDescriptor(
		cluster.Service,
		Coalesce(Dependencies.K8sCluster),
	),
	ResourceTypes.K8sNamespaceAnomalies: NewResourceDescriptor(
		namespace.Service,
		Coalesce(Dependencies.CloudApplicationNamespace),
		Coalesce(Dependencies.K8sCluster),
	),
	ResourceTypes.K8sNodeAnomalies: NewResourceDescriptor(
		node.Service,
		Coalesce(Dependencies.K8sCluster),
	),
	ResourceTypes.K8sWorkloadAnomalies: NewResourceDescriptor(
		workload.Service,
		Coalesce(Dependencies.CloudApplicationNamespace),
		Coalesce(Dependencies.K8sCluster),
	),
	ResourceTypes.ContainerBuiltinRule: NewResourceDescriptor(builtinmonitoringrule.Service),
	ResourceTypes.ContainerRule:        NewResourceDescriptor(monitoringrule.Service),
	ResourceTypes.ContainerRegistry:    NewResourceDescriptor(registry.Service),
	ResourceTypes.ContainerTechnology: NewResourceDescriptor(
		containertechnology.Service,
		Coalesce(Dependencies.Host),
		Coalesce(Dependencies.HostGroup),
	),
	ResourceTypes.RemoteEnvironments: NewResourceDescriptor(environment.Service),
	ResourceTypes.WebAppCustomErrors: NewResourceDescriptor(
		webappcustomerrors.Service,
		Dependencies.ID(ResourceTypes.WebApplication),
	),
	ResourceTypes.WebAppRequestErrors: NewResourceDescriptor(
		webapprequesterrors.Service,
		Dependencies.ID(ResourceTypes.WebApplication),
	),
	ResourceTypes.UserSettings:      NewResourceDescriptor(usersettings.Service),
	ResourceTypes.DashboardsGeneral: NewResourceDescriptor(dashboardsgeneral.Service),
	ResourceTypes.DashboardsPresets: NewResourceDescriptor(dashboardspresets.Service),
	ResourceTypes.LogProcessing: NewResourceDescriptor(
		logdpprules.Service,
		// Dependency onto other LogProcessing resources only exists because of `insertAfter`
		// Using `WeakID` here enforces that dependencies are only supposed to get replaced
		// when the LogProcessing rule that is referred to is also a candidate for the export
		Dependencies.WeakID(ResourceTypes.LogProcessing),
	),
	ResourceTypes.LogEvents: NewResourceDescriptor(logevents.Service),
	ResourceTypes.LogTimestamp: NewResourceDescriptor(
		timestampconfiguration.Service,
		Coalesce(Dependencies.Host),
		Coalesce(Dependencies.K8sCluster),
		Coalesce(Dependencies.HostGroup),
	),
	ResourceTypes.LogGrail:           NewResourceDescriptor(logsongrailactivate.Service),
	ResourceTypes.LogCustomAttribute: NewResourceDescriptor(logcustomattributes.Service),
	ResourceTypes.LogSensitiveDataMasking: NewResourceDescriptor(
		sensitivedatamasking.Service,
		Coalesce(Dependencies.Host),
		Coalesce(Dependencies.K8sCluster),
		Coalesce(Dependencies.HostGroup),
	),
	ResourceTypes.LogBuckets:                NewResourceDescriptor(logbucketsrules.Service),
	ResourceTypes.LogSecurityContext:        NewResourceDescriptor(logsecuritycontextrules.Service),
	ResourceTypes.EULASettings:              NewResourceDescriptor(eulasettings.Service),
	ResourceTypes.APIDetectionRules:         NewResourceDescriptor(apidetection.Service),
	ResourceTypes.ServiceExternalWebRequest: NewResourceDescriptor(externalwebrequest.Service),
	ResourceTypes.ServiceExternalWebService: NewResourceDescriptor(externalwebservice.Service),
	ResourceTypes.ServiceFullWebRequest:     NewResourceDescriptor(fullwebrequest.Service),
	ResourceTypes.ServiceFullWebService: NewResourceDescriptor(
		fullwebservice.Service,
		Dependencies.LegacyID(ResourceTypes.ManagementZoneV2),
		Dependencies.ID(ResourceTypes.ManagementZoneV2),
	),
	ResourceTypes.DashboardsAllowlist:        NewResourceDescriptor(dashboardsallowlist.Service),
	ResourceTypes.FailureDetectionParameters: NewResourceDescriptor(envparameters.Service),
	ResourceTypes.FailureDetectionRules: NewResourceDescriptor(
		envrules.Service,
		Dependencies.ID(ResourceTypes.FailureDetectionParameters),
		Dependencies.LegacyID(ResourceTypes.ManagementZoneV2),
		Dependencies.ID(ResourceTypes.ManagementZoneV2),
	),
	ResourceTypes.LogOneAgent: NewResourceDescriptor(
		logagentconfiguration.Service,
		Coalesce(Dependencies.Host),
		Coalesce(Dependencies.K8sCluster),
		Coalesce(Dependencies.HostGroup),
	),
	ResourceTypes.IssueTracking:            NewResourceDescriptor(issuetracking.Service),
	ResourceTypes.GeolocationSettings:      NewResourceDescriptor(geosettings.Service),
	ResourceTypes.UserSessionCustomMetrics: NewResourceDescriptor(custommetrics.Service),
	ResourceTypes.CustomUnits:              NewResourceDescriptor(customunit.Service),
	ResourceTypes.DiskAnalytics:            NewResourceDescriptor(diskanalytics.Service),
	ResourceTypes.NetworkTraffic:           NewResourceDescriptor(networktraffic.Service),
	ResourceTypes.TokenSettings:            NewResourceDescriptor(tokensettings.Service),
	ResourceTypes.ExtensionExecutionRemote: NewResourceDescriptor(
		eecremote.Service,
		Coalesce(Dependencies.EnvironmentActiveGate),
	),
	ResourceTypes.K8sPVCAnomalies: NewResourceDescriptor(
		pvc.Service,
		Coalesce(Dependencies.CloudApplicationNamespace),
		Coalesce(Dependencies.K8sCluster),
	),
	ResourceTypes.UserActionCustomMetrics: NewResourceDescriptor(useractioncustommetrics.Service),
	ResourceTypes.WebAppJavascriptVersion: NewResourceDescriptor(customrumjavascriptversion.Service),
	ResourceTypes.WebAppJavascriptUpdates: NewResourceDescriptor(
		rumjavascriptupdates.Service,
		Dependencies.ID(ResourceTypes.WebApplication),
	),
	ResourceTypes.OpenTelemetryMetrics: NewResourceDescriptor(opentelemetrymetrics.Service),
	ResourceTypes.ActiveGateUpdates: NewResourceDescriptor(
		activegateupdates.Service,
		Coalesce(Dependencies.EnvironmentActiveGate),
	),
	ResourceTypes.OneAgentDefaultMode:    NewResourceDescriptor(defaultmode.Service),
	ResourceTypes.OneAgentDefaultVersion: NewResourceDescriptor(defaultversion.Service),
	ResourceTypes.OneAgentUpdates: NewResourceDescriptor(
		oneagentupdates.Service,
		Coalesce(Dependencies.Host),
		Coalesce(Dependencies.HostGroup),
		Dependencies.ID(ResourceTypes.UpdateWindows),
	),
	ResourceTypes.LogStorage: NewResourceDescriptor(
		logstoragesettings.Service,
		Dependencies.ID(ResourceTypes.LogStorage),
		Coalesce(Dependencies.Host),
		Coalesce(Dependencies.K8sCluster),
		Coalesce(Dependencies.HostGroup),
	),
	ResourceTypes.OwnershipTeams: NewResourceDescriptor(teams.Service),
	ResourceTypes.LogCustomSource: NewResourceDescriptor(
		customlogsourcesettings.Service,
		Coalesce(Dependencies.Host),
		Coalesce(Dependencies.K8sCluster),
		Coalesce(Dependencies.HostGroup),
	),
	ResourceTypes.ApplicationDetectionV2: NewResourceDescriptor(
		appdetection.Service,
		Dependencies.ID(ResourceTypes.WebApplication),
	),
	ResourceTypes.Kubernetes: NewResourceDescriptor(
		kubernetesv2.Service,
		Coalesce(Dependencies.K8sCluster),
	),
	ResourceTypes.CloudFoundry: NewResourceDescriptor(cloudfoundryv2.Service),
	ResourceTypes.DiskAnomalyDetectionRules: NewResourceDescriptor(
		diskrules.Service,
		Coalesce(Dependencies.HostGroup),
	),
	ResourceTypes.AWSAnomalies:    NewResourceDescriptor(aws_anomalies.Service),
	ResourceTypes.VMwareAnomalies: NewResourceDescriptor(vmware_anomalies.Service),
	ResourceTypes.SLOV2: NewResourceDescriptor(
		slov2.Service,
		Dependencies.ManagementZone,
		Dependencies.LegacyID(ResourceTypes.ManagementZoneV2),
		Dependencies.ID(ResourceTypes.ManagementZoneV2),
		Dependencies.ID(ResourceTypes.CalculatedServiceMetric),
		Dependencies.ID(ResourceTypes.CalculatedWebMetric),
		Dependencies.ID(ResourceTypes.CalculatedMobileMetric),
		Dependencies.ID(ResourceTypes.CalculatedSyntheticMetric),
	),
	ResourceTypes.AutoTagV2: NewResourceDescriptor(
		autotagging.Service,
		Coalesce(Dependencies.Service),
		Coalesce(Dependencies.Host),
		Coalesce(Dependencies.HostGroup),
		Coalesce(Dependencies.ProcessGroup),
		Coalesce(Dependencies.ProcessGroupInstance),
	),
	ResourceTypes.BusinessEventsOneAgent: NewResourceDescriptor(
		incoming.Service,
		Coalesce(Dependencies.Host),
		Coalesce(Dependencies.HostGroup),
	),
	ResourceTypes.BusinessEventsOneAgentOutgoing: NewResourceDescriptor(
		outgoing.Service,
		Coalesce(Dependencies.Host),
		Coalesce(Dependencies.HostGroup),
	),
	ResourceTypes.BusinessEventsBuckets:         NewResourceDescriptor(bizevents_buckets.Service),
	ResourceTypes.BusinessEventsMetrics:         NewResourceDescriptor(bizevents_metrics.Service),
	ResourceTypes.BusinessEventsProcessing:      NewResourceDescriptor(bizevents_processing.Service),
	ResourceTypes.BusinessEventsSecurityContext: NewResourceDescriptor(bizevents_security.Service),
	ResourceTypes.WebAppKeyPerformanceCustom: NewResourceDescriptor(
		customactions.Service,
		Dependencies.ID(ResourceTypes.WebApplication),
	),
	ResourceTypes.WebAppKeyPerformanceLoad: NewResourceDescriptor(
		loadactions.Service,
		Dependencies.ID(ResourceTypes.WebApplication),
	),
	ResourceTypes.WebAppKeyPerformanceXHR: NewResourceDescriptor(
		xhractions.Service,
		Dependencies.ID(ResourceTypes.WebApplication),
	),
	ResourceTypes.MobileAppKeyPerformance: NewResourceDescriptor(
		mobilekeyperformancemetrics.Service,
		Dependencies.ID(ResourceTypes.MobileApplication),
	),
	ResourceTypes.OwnershipConfig:          NewResourceDescriptor(ownership_config.Service),
	ResourceTypes.BuiltinProcessMonitoring: NewResourceDescriptor(builtinprocessmonitoringrule.Service),
	ResourceTypes.LimitOutboundConnections: NewResourceDescriptor(allowedoutboundconnections.Service),
	ResourceTypes.SpanEvents:               NewResourceDescriptor(eventattribute.Service),
	ResourceTypes.VMware:                   NewResourceDescriptor(vmware.Service),
	ResourceTypes.CustomDevice:             NewResourceDescriptor(customdevice.Service),
	ResourceTypes.K8sMonitoring: NewResourceDescriptor(
		kubernetesmonitoring.Service,
		Coalesce(Dependencies.K8sCluster),
	),
	ResourceTypes.AutomationWorkflow: NewResourceDescriptor(
		workflows.Service,
		Dependencies.ID(ResourceTypes.AutomationSchedulingRule),
		Dependencies.ID(ResourceTypes.AutomationBusinessCalendar),
	),
	ResourceTypes.AutomationBusinessCalendar: NewResourceDescriptor(business_calendars.Service),
	ResourceTypes.AutomationSchedulingRule: NewResourceDescriptor(
		scheduling_rules.Service,
		Dependencies.ID(ResourceTypes.AutomationSchedulingRule),
		Dependencies.ID(ResourceTypes.AutomationBusinessCalendar),
	),
	ResourceTypes.CustomTags: NewResourceDescriptor(
		customtags.Service,
		Dependencies.ID(ResourceTypes.HTTPMonitor),
		Dependencies.ID(ResourceTypes.BrowserMonitor),
		Dependencies.ID(ResourceTypes.WebApplication),
		Dependencies.ID(ResourceTypes.MobileApplication),
		Coalesce(Dependencies.ProcessGroup),
		Coalesce(Dependencies.ProcessGroupInstance),
		Coalesce(Dependencies.Host),
		Coalesce(Dependencies.HostGroup),
		Coalesce(Dependencies.Service),
	),
	ResourceTypes.HostMonitoringMode: NewResourceDescriptor(
		hostmonitoringmode.Service,
		Coalesce(Dependencies.Host),
	),
	ResourceTypes.IPAddressMasking: NewResourceDescriptor(
		ipaddressmasking.Service,
		Dependencies.ID(ResourceTypes.WebApplication),
		Dependencies.ID(ResourceTypes.MobileApplication),
	),
	ResourceTypes.AppSecVulnerabilitySettings:   NewResourceDescriptor(runtimevulnerabilitydetection.Service),
	ResourceTypes.AppSecVulnerabilityThirdParty: NewResourceDescriptor(rulesettings.Service),
	ResourceTypes.AppSecVulnerabilityCode: NewResourceDescriptor(
		codelevelvulnerability.Service,
		Coalesce(Dependencies.ProcessGroup),
	),
	ResourceTypes.AppSecNotification: NewResourceDescriptor(
		notificationintegration.Service,
		Dependencies.ID(ResourceTypes.AppSecVulnerabilityAlerting),
		Dependencies.ID(ResourceTypes.AppSecAttackAlerting),
	),
	ResourceTypes.AppSecVulnerabilityAlerting: NewResourceDescriptor(
		notificationalertingprofile.Service,
		Dependencies.LegacyID(ResourceTypes.ManagementZoneV2),
		Dependencies.ID(ResourceTypes.ManagementZoneV2),
	),
	ResourceTypes.AppSecAttackAlerting: NewResourceDescriptor(
		notificationattackalertingprofile.Service,
	),
	ResourceTypes.AppSecAttackSettings: NewResourceDescriptor(attackprotectionsettings.Service),
	ResourceTypes.AppSecAttackRules: NewResourceDescriptor(
		attackprotectionadvancedconfig.Service,
		Coalesce(Dependencies.ProcessGroup),
	),
	ResourceTypes.AppSecAttackAllowlist:  NewResourceDescriptor(attackprotectionallowlistconfig.Service),
	ResourceTypes.GenericSetting:         NewResourceDescriptor(generic.Service),
	ResourceTypes.UnifiedServicesMetrics: NewResourceDescriptor(endpointmetrics.Service),
	ResourceTypes.UnifiedServicesOpenTel: NewResourceDescriptor(unifiedservicesopentel.Service),
	ResourceTypes.PlatformBucket:         NewResourceDescriptor(platformbuckets.Service),
	ResourceTypes.KeyUserAction: NewResourceDescriptor(
		keyuseractions.Service,
		Dependencies.ID(ResourceTypes.WebApplication),
	),
	ResourceTypes.UrlBasedSampling: NewResourceDescriptor(
		urlbasedsampling.Service,
		Coalesce(Dependencies.ProcessGroupInstance),
		Coalesce(Dependencies.ProcessGroup),
		Coalesce(Dependencies.CloudApplication),
		Coalesce(Dependencies.CloudApplicationNamespace),
		Coalesce(Dependencies.K8sCluster),
		Coalesce(Dependencies.HostGroup),
	),
	ResourceTypes.HostMonitoringAdvanced: NewResourceDescriptor(
		hostmonitoringadvanced.Service,
		Coalesce(Dependencies.Host),
	),
	ResourceTypes.AttributeAllowList:    NewResourceDescriptor(attributeallowlist.Service),
	ResourceTypes.AttributeBlockList:    NewResourceDescriptor(attributeblocklist.Service),
	ResourceTypes.AttributeMasking:      NewResourceDescriptor(attributemasking.Service),
	ResourceTypes.AttributesPreferences: NewResourceDescriptor(attributespreferences.Service),
	ResourceTypes.OneAgentSideMasking: NewResourceDescriptor(
		masking.Service,
		Coalesce(Dependencies.ProcessGroup),
	),
	ResourceTypes.HubSubscriptions:    NewResourceDescriptor(subscriptions.Service),
	ResourceTypes.MobileNotifications: NewResourceDescriptor(mobilenotifications.Service),
	ResourceTypes.CrashdumpAnalytics: NewResourceDescriptor(
		crashdumpanalytics.Service,
		Coalesce(Dependencies.Host),
		Coalesce(Dependencies.HostGroup),
	),
	ResourceTypes.AppMonitoring:           NewResourceDescriptor(appmonitoring.Service),
	ResourceTypes.GrailSecurityContext:    NewResourceDescriptor(securitycontext.Service),
	ResourceTypes.SiteReliabilityGuardian: NewResourceDescriptor(sitereliabilityguardian.Service),
	ResourceTypes.JiraForWorkflows:        NewResourceDescriptor(jiraconnection.Service),
	ResourceTypes.SlackForWorkflows:       NewResourceDescriptor(slackconnection.Service),
	ResourceTypes.Policy:                  NewResourceDescriptor(onprempolicies.Service),
	ResourceTypes.KubernetesApp: NewResourceDescriptor(
		kubernetesapp.Service,
		Coalesce(Dependencies.K8sCluster),
	),
	ResourceTypes.GrailMetricsAllowall:  NewResourceDescriptor(grailmetricsallowall.Service),
	ResourceTypes.GrailMetricsAllowlist: NewResourceDescriptor(grailmetricsallowlist.Service),
	ResourceTypes.WebAppBeaconEndpoint: NewResourceDescriptor(
		webappbeaconendpoint.Service,
		Dependencies.ID(ResourceTypes.WebApplication),
	),
	ResourceTypes.WebAppCustomConfigProperties: NewResourceDescriptor(
		webappcustomconfigproperties.Service,
		Dependencies.ID(ResourceTypes.WebApplication),
	),
	ResourceTypes.WebAppInjectionCookie: NewResourceDescriptor(
		webappinjectioncookie.Service,
		Dependencies.ID(ResourceTypes.WebApplication),
	),
	ResourceTypes.HTTPMonitorScript: NewResourceDescriptor(
		httpscript.Service,
		Dependencies.ID(ResourceTypes.HTTPMonitor),
	),
	ResourceTypes.UserGroup: NewResourceDescriptor(onpremusergroups.Service),
	ResourceTypes.User: NewResourceDescriptor(
		onpremusers.Service,
		Dependencies.QuotedID(ResourceTypes.UserGroup),
	),
	ResourceTypes.PolicyBinding: NewResourceDescriptor(
		onprempolicybindings.Service,
		Dependencies.ID(ResourceTypes.Policy),
		Dependencies.QuotedID(ResourceTypes.UserGroup),
	),
	ResourceTypes.MgmzPermission: NewResourceDescriptor(
		onpremmgmzpermission.Service,
		Dependencies.ID(ResourceTypes.UserGroup),
	),
	ResourceTypes.ManagedNetworkZones:       NewResourceDescriptor(managednetworkzones.Service),
	ResourceTypes.HubExtensionConfig:        NewResourceDescriptor(extension_config.Service),
	ResourceTypes.HubActiveExtensionVersion: NewResourceDescriptor(active_version.Service),
	ResourceTypes.DatabaseAppFeatureFlags:   NewResourceDescriptor(dbfeatureflags.Service),
	ResourceTypes.InfraOpsAppFeatureFlags:   NewResourceDescriptor(infraopsfeatureflags.Service),
	ResourceTypes.EBPFServiceDiscovery:      NewResourceDescriptor(ebpf.Service),
	ResourceTypes.DavisAnomalyDetectors:     NewResourceDescriptor(anomalydetectors.Service),
	ResourceTypes.LogDebugSettings:          NewResourceDescriptor(logdebugsettings.Service),
	ResourceTypes.InfraOpsAppSettings:       NewResourceDescriptor(infraopssettings.Service),
	ResourceTypes.DiskEdgeAnomalyDetectors:  NewResourceDescriptor(diskedgeanomalydetectors.Service),
	ResourceTypes.Reports: NewResourceDescriptor(
		reports.Service,
		Dependencies.ID(ResourceTypes.JSONDashboard),
	),
	ResourceTypes.NetworkMonitor: NewResourceDescriptor(
		v2monitors.Service,
		Dependencies.ID(ResourceTypes.SyntheticLocation),
	),
	ResourceTypes.NetworkMonitorOutageHandling: NewResourceDescriptor(
		networkoutagehandling.Service,
		Dependencies.ID(ResourceTypes.NetworkMonitor),
	),
	ResourceTypes.HubPermissions: NewResourceDescriptor(hubpermissions.Service),
	ResourceTypes.K8sAutomationConnections: NewResourceDescriptor(
		k8sautomationconnections.Service,
		Dependencies.WeakID(ResourceTypes.K8sAutomationConnections),
	),
	ResourceTypes.WebAppCustomInjectionRules: NewResourceDescriptor(custominjectionrules.Service),
	ResourceTypes.DiscoveryDefaultRules:      NewResourceDescriptor(defaultrules.Service),
	ResourceTypes.DiscoveryFeatureFlags:      NewResourceDescriptor(featureflags.Service),
	ResourceTypes.HistogramMetrics:           NewResourceDescriptor(histogrammetrics.Service),
	ResourceTypes.KubernetesEnrichment: NewResourceDescriptor(
		enrichment.Service,
		Coalesce(Dependencies.K8sCluster),
	),
	ResourceTypes.DevObsGitOnPrem:          NewResourceDescriptor(gitonprem.Service),
	ResourceTypes.AWSAutomationConnections: NewResourceDescriptor(awsconnection.Service),
	ResourceTypes.AWSConnection:            NewResourceDescriptor(connections_aws.Service),
	ResourceTypes.AWSConnectionRoleARN: NewResourceDescriptor(
		connections_aws_role_arn.Service,
		Dependencies.ID(ResourceTypes.AWSConnection),
	),
	ResourceTypes.DevObsAgentOptin: NewResourceDescriptor(
		agentoptin.Service,
		Coalesce(Dependencies.ProcessGroup),
		Coalesce(Dependencies.CloudApplicationNamespace),
		Coalesce(Dependencies.K8sCluster),
	),
	ResourceTypes.DevObsDataMasking:            NewResourceDescriptor(devobsmasking.Service),
	ResourceTypes.DavisCoPilot:                 NewResourceDescriptor(dataminingblocklist.Service),
	ResourceTypes.CloudDevelopmentEnvironments: NewResourceDescriptor(clouddevelopmentenvironments.Service),
	ResourceTypes.KubernetesSPM: NewResourceDescriptor(
		securityposturemanagement.Service,
		Coalesce(Dependencies.K8sCluster),
	),
	ResourceTypes.LogAgentFeatureFlags: NewResourceDescriptor(
		logagentfeatureflags.Service,
		Coalesce(Dependencies.Host),
		Coalesce(Dependencies.HostGroup),
		Coalesce(Dependencies.K8sCluster),
	),
	ResourceTypes.ProblemRecordPropagationRules:   NewResourceDescriptor(problemrecordpropagation.Service),
	ResourceTypes.ProblemFields:                   NewResourceDescriptor(fields.Service),
	ResourceTypes.AutomationControllerConnections: NewResourceDescriptor(automationcontroller.Service),
	ResourceTypes.EventDrivenAnsibleConnections:   NewResourceDescriptor(edawebhookconnection.Service),
	ResourceTypes.ServiceNowConnection:            NewResourceDescriptor(servicenowconnection.Service),
	ResourceTypes.PagerDutyConnection:             NewResourceDescriptor(pagerdutyconnection.Service),
	ResourceTypes.MSTeamsConnection:               NewResourceDescriptor(msteamsconnection.Service),
	ResourceTypes.DefaultLaunchpad:                NewResourceDescriptor(launchpad.Service),
	ResourceTypes.JenkinsConnection:               NewResourceDescriptor(jenkinsconnection.Service),
	ResourceTypes.GitLabConnection:                NewResourceDescriptor(gitlabconnection.Service),
	ResourceTypes.MSEntraIDConnection:             NewResourceDescriptor(msentraidconnection.Service),
	ResourceTypes.GitHubConnection:                NewResourceDescriptor(githubconnection.Service),
	ResourceTypes.Microsoft365EmailConnection:     NewResourceDescriptor(ms365emailconnection.Service),
	ResourceTypes.BusinessEventsCapturingVariants: NewResourceDescriptor(
		capturingvariants.Service,
		Coalesce(Dependencies.Host),
		Coalesce(Dependencies.HostGroup),
	),
	ResourceTypes.WebAppAutoInjection: NewResourceDescriptor(
		webappautoinjection.Service,
		Dependencies.ID(ResourceTypes.WebApplication),
	),
	ResourceTypes.SecurityContext:                   NewResourceDescriptor(securitycontextsettings.Service),
	ResourceTypes.Segments:                          NewResourceDescriptor(segments.Service),
	ResourceTypes.PlatformSLO:                       NewResourceDescriptor(platformslo.Service),
	ResourceTypes.AppSecVulnerabilityThirdPartyK8s:  NewResourceDescriptor(thirdpartyvulnerabilitykuberneteslabelrulesettings.Service),
	ResourceTypes.AppSecVulnerabilityThirdPartyAttr: NewResourceDescriptor(thirdpartyvulnerabilityrulesettings.Service),
	ResourceTypes.WebAppCustomProperties: NewResourceDescriptor(
		capturecustomproperties.Service,
		Dependencies.ID(ResourceTypes.WebApplication),
		Dependencies.ID(ResourceTypes.MobileApplication),
	),
	ResourceTypes.WebAppJavascriptFilename: NewResourceDescriptor(rumjavascriptfilename.Service),
	ResourceTypes.ServiceSplittingRules: NewResourceDescriptor(
		servicesplittingrules.Service,
		Coalesce(Dependencies.CloudApplicationNamespace),
		Coalesce(Dependencies.K8sCluster),
		Coalesce(Dependencies.HostGroup),
	),
	ResourceTypes.ServiceDetectionRules: NewResourceDescriptor(
		servicedetectionrules.Service,
		Coalesce(Dependencies.CloudApplicationNamespace),
		Coalesce(Dependencies.K8sCluster),
		Coalesce(Dependencies.HostGroup),
	),
	ResourceTypes.WebAppIPAddressExclusion: NewResourceDescriptor(
		ipaddressexclusion.Service,
		Dependencies.ID(ResourceTypes.WebApplication),
	),
	ResourceTypes.RPCBasedSampling: NewResourceDescriptor(
		rpcbasedsampling.Service,
		Coalesce(Dependencies.ProcessGroupInstance),
		Coalesce(Dependencies.ProcessGroup),
		Coalesce(Dependencies.CloudApplication),
		Coalesce(Dependencies.CloudApplicationNamespace),
		Coalesce(Dependencies.K8sCluster),
		Coalesce(Dependencies.HostGroup),
	),
	ResourceTypes.WebAppManualInsertion: NewResourceDescriptor(
		manualinsertion.Service,
		Dependencies.ID(ResourceTypes.WebApplication),
	),
	ResourceTypes.SettingsPermissions: NewResourceDescriptor(
		settingsPermissions.Service,

		// OwnerBasedAccessControl enabled resources
		Dependencies.ID(ResourceTypes.PagerDutyConnection),
		Dependencies.ID(ResourceTypes.Microsoft365EmailConnection),
		Dependencies.ID(ResourceTypes.MSEntraIDConnection),
		Dependencies.ID(ResourceTypes.JiraForWorkflows),
		Dependencies.ID(ResourceTypes.JenkinsConnection),
		Dependencies.ID(ResourceTypes.ServiceNowConnection),
		Dependencies.ID(ResourceTypes.AutomationControllerConnections),
		Dependencies.ID(ResourceTypes.EventDrivenAnsibleConnections),
		Dependencies.ID(ResourceTypes.K8sAutomationConnections),
		Dependencies.ID(ResourceTypes.SlackForWorkflows),
		Dependencies.ID(ResourceTypes.GitHubConnection),
		Dependencies.ID(ResourceTypes.GitLabConnection),
		Dependencies.ID(ResourceTypes.MSTeamsConnection),
		Dependencies.ID(ResourceTypes.AWSAutomationConnections),
		Dependencies.ID(ResourceTypes.AWSConnection),

		// OpenPipeline resources
		Dependencies.ID(ResourceTypes.OpenpipelineBizeventsIngestsources),
		Dependencies.ID(ResourceTypes.OpenpipelineBizeventsPipelines),
		Dependencies.ID(ResourceTypes.OpenpipelineDavisEventsIngestsources),
		Dependencies.ID(ResourceTypes.OpenpipelineDavisEventsPipelines),
		Dependencies.ID(ResourceTypes.OpenpipelineDavisProblemsIngestsources),
		Dependencies.ID(ResourceTypes.OpenpipelineDavisProblemsPipelines),
		Dependencies.ID(ResourceTypes.OpenpipelineEventsIngestsources),
		Dependencies.ID(ResourceTypes.OpenpipelineEventsPipelines),
		Dependencies.ID(ResourceTypes.OpenpipelineEventsSdlcIngestsources),
		Dependencies.ID(ResourceTypes.OpenpipelineEventsSdlcPipelines),
		Dependencies.ID(ResourceTypes.OpenpipelineEventsSecurityIngestsources),
		Dependencies.ID(ResourceTypes.OpenpipelineEventsSecurityPipelines),
		Dependencies.ID(ResourceTypes.OpenpipelineLogsIngestsources),
		Dependencies.ID(ResourceTypes.OpenpipelineLogsPipelines),
		Dependencies.ID(ResourceTypes.OpenpipelineMetricsIngestsources),
		Dependencies.ID(ResourceTypes.OpenpipelineMetricsPipelines),
		Dependencies.ID(ResourceTypes.OpenpipelineSecurityEventsIngestsources),
		Dependencies.ID(ResourceTypes.OpenpipelineSecurityEventsPipelines),
		Dependencies.ID(ResourceTypes.OpenpipelineSpansIngestsources),
		Dependencies.ID(ResourceTypes.OpenpipelineSpansPipelines),
		Dependencies.ID(ResourceTypes.OpenpipelineSystemEventsIngestsources),
		Dependencies.ID(ResourceTypes.OpenpipelineSystemEventsPipelines),
		Dependencies.ID(ResourceTypes.OpenpipelineUserEventsIngestsources),
		Dependencies.ID(ResourceTypes.OpenpipelineUserEventsPipelines),
		Dependencies.ID(ResourceTypes.OpenpipelineUsersessionsIngestsources),
		Dependencies.ID(ResourceTypes.OpenpipelineUsersessionsPipelines),

		Dependencies.ID(ResourceTypes.GenericSetting),
	),
	ResourceTypes.FailureDetectionRuleSets: NewResourceDescriptor(
		failuredetectionrulesets.Service,
		Coalesce(Dependencies.CloudApplicationNamespace),
		Coalesce(Dependencies.K8sCluster),
		Coalesce(Dependencies.HostGroup),
	),
	ResourceTypes.EndpointDetectionRulesOptIn: NewResourceDescriptor(endpointdetectionrulesoptin.Service),
	ResourceTypes.EndpointDetectionRules: NewResourceDescriptor(
		endpointdetectionrules.Service,
		Coalesce(Dependencies.CloudApplicationNamespace),
		Coalesce(Dependencies.K8sCluster),
		Coalesce(Dependencies.HostGroup),
	),

	ResourceTypes.OpenpipelineBizeventsIngestsources: NewResourceDescriptor(
		openpipelinebizeventsingestsources.Service,
		Dependencies.ID(ResourceTypes.OpenpipelineBizeventsPipelines),
	),
	ResourceTypes.OpenpipelineBizeventsPipelines: NewResourceDescriptor(openpipelinebizeventspipelines.Service),
	ResourceTypes.OpenpipelineBizeventsRouting: NewResourceDescriptor(
		openpipelinebizeventsrouting.Service,
		Dependencies.ID(ResourceTypes.OpenpipelineBizeventsPipelines),
	),
	ResourceTypes.OpenpipelineDavisEventsIngestsources: NewResourceDescriptor(
		openpipelinedaviseventsingestsources.Service,
		Dependencies.ID(ResourceTypes.OpenpipelineDavisEventsPipelines),
	),
	ResourceTypes.OpenpipelineDavisEventsPipelines: NewResourceDescriptor(openpipelinedaviseventspipelines.Service),
	ResourceTypes.OpenpipelineDavisEventsRouting: NewResourceDescriptor(
		openpipelinedaviseventsrouting.Service,
		Dependencies.ID(ResourceTypes.OpenpipelineDavisEventsPipelines),
	),
	ResourceTypes.OpenpipelineDavisProblemsIngestsources: NewResourceDescriptor(
		openpipelinedavisproblemsingestsources.Service,
		Dependencies.ID(ResourceTypes.OpenpipelineDavisProblemsPipelines),
	),
	ResourceTypes.OpenpipelineDavisProblemsPipelines: NewResourceDescriptor(openpipelinedavisproblemspipelines.Service),
	ResourceTypes.OpenpipelineDavisProblemsRouting: NewResourceDescriptor(
		openpipelinedavisproblemsrouting.Service,
		Dependencies.ID(ResourceTypes.OpenpipelineDavisProblemsPipelines),
	),
	ResourceTypes.OpenpipelineEventsIngestsources: NewResourceDescriptor(
		openpipelineeventsingestsources.Service,
		Dependencies.ID(ResourceTypes.OpenpipelineEventsPipelines),
	),
	ResourceTypes.OpenpipelineEventsPipelines: NewResourceDescriptor(openpipelineeventspipelines.Service),
	ResourceTypes.OpenpipelineEventsRouting: NewResourceDescriptor(
		openpipelineeventsrouting.Service,
		Dependencies.ID(ResourceTypes.OpenpipelineEventsPipelines),
	),
	ResourceTypes.OpenpipelineEventsSdlcIngestsources: NewResourceDescriptor(
		openpipelineeventssdlcingestsources.Service,
		Dependencies.ID(ResourceTypes.OpenpipelineEventsSdlcPipelines),
	),
	ResourceTypes.OpenpipelineEventsSdlcPipelines: NewResourceDescriptor(openpipelineeventssdlcpipelines.Service),
	ResourceTypes.OpenpipelineEventsSdlcRouting: NewResourceDescriptor(
		openpipelineeventssdlcrouting.Service,
		Dependencies.ID(ResourceTypes.OpenpipelineEventsSdlcPipelines),
	),
	ResourceTypes.OpenpipelineEventsSecurityIngestsources: NewResourceDescriptor(
		openpipelineeventssecurityingestsources.Service,
		Dependencies.ID(ResourceTypes.OpenpipelineEventsSecurityPipelines),
	),
	ResourceTypes.OpenpipelineEventsSecurityPipelines: NewResourceDescriptor(openpipelineeventssecuritypipelines.Service),
	ResourceTypes.OpenpipelineEventsSecurityRouting: NewResourceDescriptor(
		openpipelineeventssecurityrouting.Service,
		Dependencies.ID(ResourceTypes.OpenpipelineEventsSecurityPipelines),
	),
	ResourceTypes.OpenpipelineLogsIngestsources: NewResourceDescriptor(
		openpipelinelogsingestsources.Service,
		Dependencies.ID(ResourceTypes.OpenpipelineLogsPipelines),
	),
	ResourceTypes.OpenpipelineLogsPipelines: NewResourceDescriptor(openpipelinelogspipelines.Service),
	ResourceTypes.OpenpipelineLogsRouting: NewResourceDescriptor(
		openpipelinelogsrouting.Service,
		Dependencies.ID(ResourceTypes.OpenpipelineLogsPipelines),
	),
	ResourceTypes.OpenpipelineMetricsIngestsources: NewResourceDescriptor(
		openpipelinemetricsingestsources.Service,
		Dependencies.ID(ResourceTypes.OpenpipelineMetricsPipelines),
	),
	ResourceTypes.OpenpipelineMetricsPipelines: NewResourceDescriptor(openpipelinemetricspipelines.Service),
	ResourceTypes.OpenpipelineMetricsRouting: NewResourceDescriptor(
		openpipelinemetricsrouting.Service,
		Dependencies.ID(ResourceTypes.OpenpipelineMetricsPipelines),
	),
	ResourceTypes.OpenpipelineSecurityEventsIngestsources: NewResourceDescriptor(
		openpipelinesecurityeventsingestsources.Service,
		Dependencies.ID(ResourceTypes.OpenpipelineSecurityEventsPipelines),
	),
	ResourceTypes.OpenpipelineSecurityEventsPipelines: NewResourceDescriptor(openpipelinesecurityeventspipelines.Service),
	ResourceTypes.OpenpipelineSecurityEventsRouting: NewResourceDescriptor(
		openpipelinesecurityeventsrouting.Service,
		Dependencies.ID(ResourceTypes.OpenpipelineSecurityEventsPipelines),
	),
	ResourceTypes.OpenpipelineSpansIngestsources: NewResourceDescriptor(
		openpipelinespansingestsources.Service,
		Dependencies.ID(ResourceTypes.OpenpipelineSpansPipelines),
	),
	ResourceTypes.OpenpipelineSpansPipelines: NewResourceDescriptor(openpipelinespanspipelines.Service),
	ResourceTypes.OpenpipelineSpansRouting: NewResourceDescriptor(
		openpipelinespansrouting.Service,
		Dependencies.ID(ResourceTypes.OpenpipelineSpansPipelines),
	),
	ResourceTypes.OpenpipelineSystemEventsIngestsources: NewResourceDescriptor(
		openpipelinesystemeventsingestsources.Service,
		Dependencies.ID(ResourceTypes.OpenpipelineSystemEventsPipelines),
	),
	ResourceTypes.OpenpipelineSystemEventsPipelines: NewResourceDescriptor(openpipelinesystemeventspipelines.Service),
	ResourceTypes.OpenpipelineSystemEventsRouting: NewResourceDescriptor(
		openpipelinesystemeventsrouting.Service,
		Dependencies.ID(ResourceTypes.OpenpipelineSystemEventsPipelines),
	),
	ResourceTypes.OpenpipelineUserEventsIngestsources: NewResourceDescriptor(
		openpipelineusereventsingestsources.Service,
		Dependencies.ID(ResourceTypes.OpenpipelineUserEventsPipelines),
	),
	ResourceTypes.OpenpipelineUserEventsPipelines: NewResourceDescriptor(openpipelineusereventspipelines.Service),
	ResourceTypes.OpenpipelineUserEventsRouting: NewResourceDescriptor(
		openpipelineusereventsrouting.Service,
		Dependencies.ID(ResourceTypes.OpenpipelineUserEventsPipelines),
	),
	ResourceTypes.OpenpipelineUsersessionsIngestsources: NewResourceDescriptor(
		openpipelineusersessionsingestsources.Service,
		Dependencies.ID(ResourceTypes.OpenpipelineUsersessionsPipelines),
	),
	ResourceTypes.OpenpipelineUsersessionsPipelines: NewResourceDescriptor(openpipelineusersessionspipelines.Service),
	ResourceTypes.OpenpipelineUsersessionsRouting: NewResourceDescriptor(
		openpipelineusersessionsrouting.Service,
		Dependencies.ID(ResourceTypes.OpenpipelineUsersessionsPipelines),
	),
	ResourceTypes.AutomationApproval: NewResourceDescriptor(approval.Service),
}

type ResourceExclusion struct {
	ResourceType ResourceType
	Reason       string
}

type ResourceExclusionGroup struct {
	Reason     string
	Exclusions []ResourceExclusion
}

var excludeListedResourceGroups = []ResourceExclusionGroup{
	// ResourceTypes.SpanAttribute, // Replaced by dynatrace_attribute_allow_list and dynatrace_attribute_masking. Commenting out of the excludeList temporarily..
	// ResourceTypes.SpanEvents, // Replaced by dynatrace_attribute_allow_list and dynatrace_attribute_masking. Commenting out of the excludeList temporarily..
	// ResourceAttributes, // Replaced by dynatrace_attribute_allow_list and dynatrace_attribute_masking. Commenting out of the excludeList temporarily..

	{
		Reason: "Officially deprecated resources (EOL)",
		Exclusions: []ResourceExclusion{
			{ResourceTypes.AlertingProfile, "Replaced by dynatrace_alerting"},
			{ResourceTypes.CustomAnomalies, "Replaced by dynatrace_metric_events"},
			{ResourceTypes.LogGrail, "Only meant to be used for the initial Logs powered by Grail activation"},
			{ResourceTypes.MaintenanceWindow, "Replaced by dynatrace_maintenance"},
			{ResourceTypes.Notification, "Replaced by dynatrace_<type>_notification"},
		},
	},
	{
		Reason: "Deprecated resources due to better alternatives",
		Exclusions: []ResourceExclusion{
			{ResourceTypes.ApplicationAnomalies, "Replaced by dynatrace_web_app_anomalies"},
			{ResourceTypes.ApplicationDataPrivacy, "Replaced by dynatrace_data_privacy and dynatrace_session_replay_web_privacy"},
			{ResourceTypes.AutoTag, "Replaced by dynatrace_autotag_v2"},
			{ResourceTypes.CloudFoundryCredentials, "Replaced by dynatrace_cloud_foundry"},
			{ResourceTypes.Dashboard, "Replaced by dynatrace_json_dashboard"},
			{ResourceTypes.DatabaseAnomalies, "Replaced by dynatrace_database_anomalies_v2"},
			{ResourceTypes.DiskEventAnomalies, "Replaced by dynatrace_disk_anomaly_rules"},
			{ResourceTypes.HostAnomalies, "Replaced by dynatrace_host_anomalies_v2"},
			{ResourceTypes.KubernetesCredentials, "Replaced by dynatrace_kubernetes"},
			{ResourceTypes.ManagementZone, "Replaced by dynatrace_management_zone_v2"},
			{ResourceTypes.ProcessGroupAnomalies, "Replaced by dynatrace_pg_alerting"},
			{ResourceTypes.ServiceAnomalies, "Replaced by dynatrace_service_anomalies_v2"},
			{ResourceTypes.SLO, "Replaced by dynatrace_slo_v2"},
		},
	},
	{
		Reason: "Resources waiting for full coverage to deprecate v1/v2 counterpart",
		Exclusions: []ResourceExclusion{
			{ResourceTypes.ApplicationDetectionV2, "Ordering of Settings is not yet supported. Use `dynatrace_application_detection_rule` instead"},
			{ResourceTypes.MobileAppRequestErrors, "JS errors missing, use dynatrace_application_error_rules"},
			{ResourceTypes.WebAppCustomErrors, "JS errors missing, use dynatrace_application_error_rules"},
			{ResourceTypes.WebAppRequestErrors, "JS errors missing, use dynatrace_application_error_rules"},
		},
	},
	{
		Reason: "Account management requires OAuth2 client and is specific to SaaS",
		Exclusions: []ResourceExclusion{
			{ResourceTypes.IAMUser, ""},
			{ResourceTypes.IAMGroup, ""},
			{ResourceTypes.IAMPermission, ""},
			{ResourceTypes.IAMPolicy, ""},
			{ResourceTypes.IAMPolicyBindings, ""},
			{ResourceTypes.IAMPolicyBindingsV2, ""},
			{ResourceTypes.IAMPolicyBoundary, ""},
		},
	},
	{
		Reason: "Cluster management is specific to Managed",
		Exclusions: []ResourceExclusion{
			{ResourceTypes.Policy, ""},
		},
	},
	{
		Reason: "Requires a personal token",
		Exclusions: []ResourceExclusion{
			{ResourceTypes.UserSettings, ""},
		},
	},
	{
		Reason: "Potential issues for migration use cases",
		Exclusions: []ResourceExclusion{
			{ResourceTypes.MetricMetadata, ""},
			{ResourceTypes.MetricQuery, ""},
		},
	},
	{
		Reason: "Not included in export - to be discussed",
		Exclusions: []ResourceExclusion{
			{ResourceTypes.AzureService, ""},
			{ResourceTypes.AWSService, ""},
			{ResourceTypes.AGToken, ""},
			{ResourceTypes.MobileAppKeyPerformance, ""},
			{ResourceTypes.HTTPMonitorScript, ""},
		},
	},
	{
		Reason: "Requires an OAuth2 client",
		Exclusions: []ResourceExclusion{
			{ResourceTypes.AutomationBusinessCalendar, ""},
			{ResourceTypes.AutomationSchedulingRule, ""},
			{ResourceTypes.AutomationWorkflow, ""},
			{ResourceTypes.Documents, ""},
			{ResourceTypes.PlatformBucket, ""},
			{ResourceTypes.Segments, ""},
			{ResourceTypes.PlatformSLO, ""},
			{ResourceTypes.SettingsPermissions, ""},
		},
	},
	{
		Reason: "Potential time consuming execution",
		Exclusions: []ResourceExclusion{
			{ResourceTypes.CustomTags, ""},
			{ResourceTypes.CustomDevice, ""},
		},
	},
	{
		Reason: "Requires the feature to be activated",
		Exclusions: []ResourceExclusion{
			{ResourceTypes.AppSecVulnerabilitySettings, ""},
			{ResourceTypes.AppSecVulnerabilityThirdParty, ""},
			{ResourceTypes.AppSecVulnerabilityCode, ""},
			{ResourceTypes.AppSecNotification, ""},
			{ResourceTypes.AppSecVulnerabilityAlerting, ""},
			{ResourceTypes.AppSecAttackAlerting, ""},
			{ResourceTypes.AppSecAttackSettings, ""},
			{ResourceTypes.AppSecAttackRules, ""},
			{ResourceTypes.AppSecAttackAllowlist, ""},
			{ResourceTypes.AppSecVulnerabilityThirdPartyK8s, ""},
			{ResourceTypes.AppSecVulnerabilityThirdPartyAttr, ""},
		},
	},
	{
		Reason: "Requires the app from Dynatrace Hub",
		Exclusions: []ResourceExclusion{
			{ResourceTypes.SiteReliabilityGuardian, ""},
			{ResourceTypes.JiraForWorkflows, ""},
			{ResourceTypes.SlackForWorkflows, ""},
		},
	},
	{
		Reason: "Would lead to circular dependencies",
		Exclusions: []ResourceExclusion{
			{ResourceTypes.AWSConnection, "Please export `dynatrace_aws_connection_role_arn` with `-ref` instead"},
		},
	},
	{
		Reason: "Generic resource against any Setting 2.0 schema",
		Exclusions: []ResourceExclusion{
			{ResourceTypes.GenericSetting, ""},
		},
	},
}

var excludeListedResources = genExcludeListedResourceGroups()

func genExcludeListedResourceGroups() []ResourceType {
	result := []ResourceType{}
	for _, eg := range GetExcludeListedResourceGroups() {
		for _, ex := range eg.Exclusions {
			result = append(result, ex.ResourceType)
		}
	}
	return result
}

var ENABLE_EXPORT_DASHBOARD = os.Getenv("DYNATRACE_ENABLE_EXPORT_DASHBOARD") == "true"

func GetExcludeListedResourceGroups() []ResourceExclusionGroup {
	if ENABLE_EXPORT_DASHBOARD {
		return excludeListedResourceGroups
	}

	return append(excludeListedResourceGroups, ResourceExclusionGroup{
		Reason: "Production environments may contain 10k+ dashboards",
		Exclusions: []ResourceExclusion{
			{ResourceTypes.Reports, ""},
			{ResourceTypes.JSONDashboard, ""},
			{ResourceTypes.DashboardSharing, ""},
			{ResourceTypes.JSONDashboardBase, ""},
		},
	})
}

func GetExcludeListedResources() []ResourceType {

	if ENABLE_EXPORT_DASHBOARD {
		return excludeListedResources
	}

	// Excluded due to the potential of a large amount of dashboards
	// Excluded since it is retrieved as a child resource of dynatrace_json_dashboard (ResourceTypes.DashboardSharing)
	return append(excludeListedResources, ResourceTypes.JSONDashboard, ResourceTypes.DashboardSharing, ResourceTypes.JSONDashboardBase)

}

func Service(credentials *rest.Credentials, resourceType ResourceType) settings.CRUDService[settings.Settings] {
	return AllResources[resourceType].Service(credentials)
}

// func DSService(credentials *rest.Credentials, dataSourceType DataSourceType) settings.RService[settings.Settings] {
// 	return AllDataSources[dataSourceType].Service(credentials)
// }
