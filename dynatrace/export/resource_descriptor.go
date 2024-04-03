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

	dbfeatureflags "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/database/featureflags"
	infraopsfeatureflags "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/infraops/featureflags"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/jiraconnection"
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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/appsec/attackprotectionadvancedconfig"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/appsec/attackprotectionallowlistconfig"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/appsec/attackprotectionsettings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/appsec/codelevelvulnerability"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/appsec/notificationalertingprofile"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/appsec/notificationattackalertingprofile"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/appsec/notificationintegration"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/appsec/rulesettings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/appsec/runtimevulnerabilitydetection"
	kubernetesapp "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/apptransition/kubernetes"
	attributeallowlist "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/attribute/allowlist"
	attributeblocklist "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/attribute/blocklist"
	attributemasking "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/attribute/masking"
	attributespreferences "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/attribute/preferences"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/auditlog"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/availability/processgroupalerting"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/bizevents/http/incoming"
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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/deployment/oneagent/defaultversion"
	oneagentupdates "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/deployment/oneagent/updates"
	diskanalytics "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/disk/analytics/extension"
	diskoptions "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/disk/options"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/dtjavascriptruntime/allowedoutboundconnections"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/dtjavascriptruntime/appmonitoring"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/eec/local"
	eecremote "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/eec/remote"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/eulasettings"
	networktraffic "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/exclude/network/traffic"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/generic"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/geosettings"
	grailmetricsallowall "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/grail/metrics/allowall"
	grailmetricsallowlist "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/grail/metrics/allowlist"
	hostmonitoring "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/host/monitoring"
	hostmonitoringadvanced "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/host/monitoring/advanced"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/host/monitoring/aixkernelextension"
	hostmonitoringmode "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/host/monitoring/mode"
	hostprocessgroupmonitoring "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/host/processgroups/monitoringstate"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hub/subscriptions"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/ibmmq/imsbridges"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/ibmmq/queuemanagers"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/ibmmq/queuesharinggroup"
	issuetracking "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/issuetracking/integration"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/customlogsourcesettings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/logagentconfiguration"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/logbucketsrules"
	logcustomattributes "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/logcustomattributes"
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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/process/builtinprocessmonitoringrule"
	processmonitoring "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/process/monitoring"
	customprocessmonitoring "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/process/monitoring/custom"
	processavailability "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/processavailability"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/processgroup/advanceddetectionrule"
	workloaddetection "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/processgroup/cloudapplication/workloaddetection"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/processgroup/detectionflags"
	processgroupmonitoring "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/processgroup/monitoring/state"
	processgroupsimpledetection "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/processgroup/simpledetectionrule"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/processvisibility"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/remote/environment"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/resourceattribute"
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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/beacondomainorigins"
	webappbeaconendpoint "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/beaconendpoint"
	webappcustomconfigproperties "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/customconfigurationproperties"
	webappcustomerrors "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/customerrors"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/customrumjavascriptversion"
	rumwebenablement "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/enablement"
	webappinjectioncookie "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/injection/cookie"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/keyperformancemetric/customactions"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/keyperformancemetric/loadactions"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/keyperformancemetric/xhractions"
	webapprequesterrors "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/requesterrors"
	webappresourcecleanup "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/resourcecleanuprules"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/resourcetypes"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/rumjavascriptupdates"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/servicedetection/externalwebrequest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/servicedetection/externalwebservice"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/servicedetection/fullwebrequest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/servicedetection/fullwebservice"
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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/bindings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/groups"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/permissions"
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
	processgroup_naming "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/naming/processgroups"
	service_naming "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/naming/services"
	networkzone "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/networkzones"

	envparameters "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/failuredetection/environment/parameters"
	envrules "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/failuredetection/environment/rules"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/failuredetection/service/generalparameters"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/failuredetection/service/httpparameters"
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
)

func NewResourceDescriptor[T settings.Settings](fn func(credentials *settings.Credentials) settings.CRUDService[T], dependencies ...Dependency) ResourceDescriptor {
	return ResourceDescriptor{
		Service: func(credentials *settings.Credentials) settings.CRUDService[settings.Settings] {
			return &settings.GenericCRUDService[T]{Service: cache.CRUD(fn(credentials))}
		},
		protoType:    newSettings(fn),
		Dependencies: dependencies,
	}
}

func NewChildResourceDescriptor[T settings.Settings](fn func(credentials *settings.Credentials) settings.CRUDService[T], parent ResourceType, dependencies ...Dependency) ResourceDescriptor {
	return ResourceDescriptor{
		Service: func(credentials *settings.Credentials) settings.CRUDService[settings.Settings] {
			return &settings.GenericCRUDService[T]{Service: cache.CRUD(fn(credentials))}
		},
		protoType:    newSettings(fn),
		Dependencies: dependencies,
		Parent:       &parent,
	}
}

func newSettings[T settings.Settings](sfn func(credentials *settings.Credentials) settings.CRUDService[T]) T {
	var proto T
	return reflect.New(reflect.TypeOf(proto).Elem()).Interface().(T)
}

type ResourceDescriptor struct {
	Dependencies []Dependency
	Service      func(credentials *settings.Credentials) settings.CRUDService[settings.Settings]
	protoType    settings.Settings
	except       func(id string, name string) bool
	Parent       *ResourceType
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

var AllResources = map[ResourceType]ResourceDescriptor{
	ResourceTypes.Alerting: NewResourceDescriptor(
		alerting.Service,
		Dependencies.LegacyID(ResourceTypes.ManagementZoneV2),
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
	ResourceTypes.Credentials: NewResourceDescriptor(
		vault.Service,
		Dependencies.ID(ResourceTypes.Credentials),
	),
	ResourceTypes.JSONDashboardBase: NewResourceDescriptor(
		jsondashboardsbase.Service,
	),
	ResourceTypes.JSONDashboard: NewChildResourceDescriptor(
		jsondashboards.Service,
		ResourceTypes.JSONDashboardBase,
		Dependencies.DashboardLinkID(true),
		Dependencies.LegacyID(ResourceTypes.ManagementZoneV2),
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
	ResourceTypes.HostNaming:   NewResourceDescriptor(host_naming.Service),
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
	),
	ResourceTypes.ManagementZoneV2: NewResourceDescriptor(v2managementzones.Service),
	ResourceTypes.MetricEvents: NewResourceDescriptor(
		metricevents.Service,
		Dependencies.LegacyID(ResourceTypes.ManagementZoneV2),
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
	ResourceTypes.QueueManager:       NewResourceDescriptor(queuemanagers.Service),
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
	),
	ResourceTypes.ManagementZone: NewResourceDescriptor(managementzonesv1.Service),
	ResourceTypes.Dashboard: NewResourceDescriptor(
		dashboards.Service,
		Dependencies.LegacyID(ResourceTypes.ManagementZoneV2),
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
		Dependencies.ID(ResourceTypes.IAMPermission),
		Dependencies.Tenant,
	),
	ResourceTypes.IAMPermission:       NewResourceDescriptor(permissions.Service),
	ResourceTypes.IAMPolicy:           NewResourceDescriptor(policies.Service),
	ResourceTypes.IAMPolicyBindings:   NewResourceDescriptor(bindings.Service),
	ResourceTypes.IAMPolicyBindingsV2: NewResourceDescriptor(v2bindings.Service),
	ResourceTypes.DDUPool:             NewResourceDescriptor(ddupool.Service),
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
		Dependencies.ID(ResourceTypes.LogProcessing),
	),
	ResourceTypes.LogEvents:                  NewResourceDescriptor(logevents.Service),
	ResourceTypes.LogTimestamp:               NewResourceDescriptor(timestampconfiguration.Service),
	ResourceTypes.LogGrail:                   NewResourceDescriptor(logsongrailactivate.Service),
	ResourceTypes.LogCustomAttribute:         NewResourceDescriptor(logcustomattributes.Service),
	ResourceTypes.LogSensitiveDataMasking:    NewResourceDescriptor(sensitivedatamasking.Service),
	ResourceTypes.LogBuckets:                 NewResourceDescriptor(logbucketsrules.Service),
	ResourceTypes.LogSecurityContext:         NewResourceDescriptor(logsecuritycontextrules.Service),
	ResourceTypes.EULASettings:               NewResourceDescriptor(eulasettings.Service),
	ResourceTypes.APIDetectionRules:          NewResourceDescriptor(apidetection.Service),
	ResourceTypes.ServiceExternalWebRequest:  NewResourceDescriptor(externalwebrequest.Service),
	ResourceTypes.ServiceExternalWebService:  NewResourceDescriptor(externalwebservice.Service),
	ResourceTypes.ServiceFullWebRequest:      NewResourceDescriptor(fullwebrequest.Service),
	ResourceTypes.ServiceFullWebService:      NewResourceDescriptor(fullwebservice.Service),
	ResourceTypes.DashboardsAllowlist:        NewResourceDescriptor(dashboardsallowlist.Service),
	ResourceTypes.FailureDetectionParameters: NewResourceDescriptor(envparameters.Service),
	ResourceTypes.FailureDetectionRules: NewResourceDescriptor(
		envrules.Service,
		Dependencies.ID(ResourceTypes.FailureDetectionParameters),
	),
	ResourceTypes.LogOneAgent:              NewResourceDescriptor(logagentconfiguration.Service),
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
	ResourceTypes.OneAgentDefaultVersion: NewResourceDescriptor(defaultversion.Service),
	ResourceTypes.OneAgentUpdates: NewResourceDescriptor(
		oneagentupdates.Service,
		Coalesce(Dependencies.Host),
		Coalesce(Dependencies.HostGroup),
		Dependencies.ID(ResourceTypes.UpdateWindows),
	),
	ResourceTypes.LogStorage: NewResourceDescriptor(
		logstoragesettings.Service,
		Coalesce(Dependencies.Host),
		Coalesce(Dependencies.HostGroup),
	),
	ResourceTypes.OwnershipTeams:  NewResourceDescriptor(teams.Service),
	ResourceTypes.LogCustomSource: NewResourceDescriptor(customlogsourcesettings.Service),
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
		Dependencies.LegacyID(ResourceTypes.ManagementZoneV2),
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
}

var excludeListedResources = []ResourceType{
	// Officially deprecated resources (EOL)
	ResourceTypes.AlertingProfile,   // Replaced by dynatrace_alerting
	ResourceTypes.CustomAnomalies,   // Replaced by dynatrace_metric_events
	ResourceTypes.MaintenanceWindow, // Replaced by dynatrace_maintenance
	ResourceTypes.Notification,      // Replaced by dynatrace_<type>_notification
	// ResourceTypes.SpanAttribute, // Replaced by dynatrace_attribute_allow_list and dynatrace_attribute_masking. Commenting out of the excludeList temporarily..
	// ResourceTypes.SpanEvents, // Replaced by dynatrace_attribute_allow_list and dynatrace_attribute_masking. Commenting out of the excludeList temporarily..
	// ResourceAttributes, // Replaced by dynatrace_attribute_allow_list and dynatrace_attribute_masking. Commenting out of the excludeList temporarily..

	// Deprecated resources due to better alternatives
	ResourceTypes.ApplicationAnomalies,    // Replaced by dynatrace_web_app_anomalies
	ResourceTypes.ApplicationDataPrivacy,  // Replaced by dynatrace_data_privacy and dynatrace_session_replay_web_privacy
	ResourceTypes.AutoTag,                 // Replaced by dynatrace_autotag_v2
	ResourceTypes.CloudFoundryCredentials, // Replaced by dynatrace_cloud_foundry
	ResourceTypes.Dashboard,               // Replaced by dynatrace_json_dashboard
	ResourceTypes.DatabaseAnomalies,       // Replaced by dynatrace_database_anomalies_v2
	ResourceTypes.DiskEventAnomalies,      // Replaced by dynatrace_disk_anomaly_rules
	ResourceTypes.HostAnomalies,           // Replaced by dynatrace_host_anomalies_v2
	ResourceTypes.KubernetesCredentials,   // Replaced by dynatrace_kubernetes
	ResourceTypes.ManagementZone,          // Replaced by dynatrace_management_zone_v2
	ResourceTypes.ProcessGroupAnomalies,   // Replaced by dynatrace_pg_alerting
	ResourceTypes.ServiceAnomalies,        // Replaced by dynatrace_service_anomalies_v2
	ResourceTypes.SLO,                     // Replaced by dynatrace_slo_v2

	// Resources waiting for full coverage to deprecate v1/v2 counterpart
	ResourceTypes.ApplicationDetectionV2, // Cannot handle ordering, use dynatrace_application_detection_rule
	ResourceTypes.MobileAppRequestErrors, // JS errors missing, use dynatrace_application_error_rules
	ResourceTypes.WebAppCustomErrors,     // JS errors missing, use dynatrace_application_error_rules
	ResourceTypes.WebAppRequestErrors,    // JS errors missing, use dynatrace_application_error_rules

	// Excluded since configuration is under account management
	ResourceTypes.IAMUser,
	ResourceTypes.IAMGroup,
	ResourceTypes.IAMPermission,
	ResourceTypes.IAMPolicy,
	ResourceTypes.IAMPolicyBindings,
	ResourceTypes.IAMPolicyBindingsV2,

	// Cluster Resources
	ResourceTypes.Policy,

	ResourceTypes.UserSettings, // Excluded since it requires a personal token

	// Not included in export - to be discussed
	ResourceTypes.AzureService,
	ResourceTypes.AWSService,
	ResourceTypes.AutomationWorkflow,
	ResourceTypes.AutomationBusinessCalendar,
	ResourceTypes.AutomationSchedulingRule,
	ResourceTypes.AGToken,
	ResourceTypes.MobileAppKeyPerformance,

	// Not included in export - may cause issues for migration use cases
	ResourceTypes.MetricMetadata,
	ResourceTypes.MetricQuery,

	// Not included in default export -  excluded due to potential time consuming execution
	ResourceTypes.CustomTags,
	ResourceTypes.CustomDevice,

	// Deprecated since it is only meant to be used for the initial Logs powered by Grail activation
	ResourceTypes.LogGrail,

	// Excluding AppSec resources from default export since it requires the feature to be activated
	ResourceTypes.AppSecVulnerabilitySettings,
	ResourceTypes.AppSecVulnerabilityThirdParty,
	ResourceTypes.AppSecVulnerabilityCode,
	ResourceTypes.AppSecNotification,
	ResourceTypes.AppSecVulnerabilityAlerting,
	ResourceTypes.AppSecAttackAlerting,
	ResourceTypes.AppSecAttackSettings,
	ResourceTypes.AppSecAttackRules,
	ResourceTypes.AppSecAttackAllowlist,

	// Excluding resources that require apps from Dynatrace Hub
	ResourceTypes.SiteReliabilityGuardian,
	ResourceTypes.JiraForWorkflows,
	ResourceTypes.SlackForWorkflows,

	// Incubator
	ResourceTypes.GenericSetting,
	ResourceTypes.PlatformBucket,
}

var ENABLE_EXPORT_DASHBOARD = os.Getenv("DYNATRACE_ENABLE_EXPORT_DASHBOARD") == "true"

func GetExcludeListedResources() []ResourceType {

	if ENABLE_EXPORT_DASHBOARD {
		return excludeListedResources
	}

	// Excluded due to the potential of a large amount of dashboards
	// Excluded since it is retrieved as a child resource of dynatrace_json_dashboard (ResourceTypes.DashboardSharing)
	return append(excludeListedResources, ResourceTypes.JSONDashboard, ResourceTypes.DashboardSharing)

}

func Service(credentials *settings.Credentials, resourceType ResourceType) settings.CRUDService[settings.Settings] {
	return AllResources[resourceType].Service(credentials)
}

// func DSService(credentials *settings.Credentials, dataSourceType DataSourceType) settings.RService[settings.Settings] {
// 	return AllDataSources[dataSourceType].Service(credentials)
// }
