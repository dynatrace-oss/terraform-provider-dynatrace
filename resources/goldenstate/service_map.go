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

package goldenstate

import (
	dbfeatureflags "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/database/featureflags"
	infraopsfeatureflags "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/infraops/featureflags"
	infraopssettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/app/dynatrace/infraops/settings"
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
	ebpf "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/ebpf/service/discovery"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/eec/local"
	eecremote "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/eec/remote"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/eulasettings"
	networktraffic "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/exclude/network/traffic"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/geosettings"
	hostmonitoring "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/host/monitoring"
	hostmonitoringadvanced "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/host/monitoring/advanced"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/host/monitoring/aixkernelextension"
	hostmonitoringmode "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/host/monitoring/mode"
	hostprocessgroupmonitoring "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/host/processgroups/monitoringstate"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hub/subscriptions"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/ibmmq/imsbridges"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/ibmmq/queuemanagers"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/ibmmq/queuesharinggroup"
	diskedgeanomalydetectors "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/infrastructure/diskedge/anomalydetectors"
	issuetracking "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/issuetracking/integration"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/customlogsourcesettings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/logagentconfiguration"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/logbucketsrules"
	logcustomattributes "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/logcustomattributes"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/logdebugsettings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/logdpprules"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/logevents"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/logsecuritycontextrules"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/logstoragesettings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/schemalesslogmetric"
	sensitivedatamasking "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/sensitivedatamaskingsettings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/timestampconfiguration"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/mainframe/mqfilters"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/mainframe/txmonitoring"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/mainframe/txstartfilters"
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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/overloadprevention"
	rumprocessgroup "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/processgroup"
	rumproviderbreakdown "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/providerbreakdown"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/resourcetimingorigins"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/userexperiencescore"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/beacondomainorigins"
	webappbeaconendpoint "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/beaconendpoint"
	webappcustomconfigproperties "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/customconfigurationproperties"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/customrumjavascriptversion"
	rumwebenablement "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/enablement"
	webappinjectioncookie "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/injection/cookie"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/keyperformancemetric/customactions"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/keyperformancemetric/loadactions"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/keyperformancemetric/xhractions"
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
	networkoutagehandling "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/synthetic/multiprotocol/outagehandling"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/tags/autotagging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/tokens/tokensettings"
	unifiedservicesopentel "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/unifiedservices/enablement"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/unifiedservices/endpointmetrics"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/urlbasedsampling"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/usability/analytics"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/useractioncustommetrics"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/virtualization/vmware"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/customservices"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/reports"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/requestattributes"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/dashboards"
	locations "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/locations/private"
	v2monitors "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/synthetic/monitors"

	v2managementzones "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/managementzones"

	host_naming "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/naming/hosts"
	processgroup_naming "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/naming/processgroups"
	service_naming "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/naming/services"
	networkzone "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/networkzones"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/davis/anomalydetectors"
	envparameters "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/failuredetection/environment/parameters"
	envrules "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/failuredetection/environment/rules"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/failuredetection/service/generalparameters"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/failuredetection/service/httpparameters"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/mobile"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web/detection"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web/errors"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web/keyuseractions"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/credentials/aws"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/credentials/azure"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/requestnaming"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/browser"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/http"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/credentials/vault"

	v2maintenance "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/alerting/maintenancewindow"
	calculated_mobile_metrics "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/metrics/calculated/mobile"
	calculated_service_metrics "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/metrics/calculated/service"
	calculated_synthetic_metrics "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/metrics/calculated/synthetic"
	calculated_web_metrics "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/metrics/calculated/web"
)

var serviceMap = map[export.ResourceType]ServiceFunc{
	export.ResourceTypes.Alerting:                  Wrap(alerting.Service),
	export.ResourceTypes.ManagementZoneV2:          Wrap(v2managementzones.Service),
	export.ResourceTypes.AutoTagV2:                 Wrap(autotagging.Service),
	export.ResourceTypes.RequestAttribute:          Wrap(requestattributes.Service),
	export.ResourceTypes.QueueManager:              Wrap(queuemanagers.Service),
	export.ResourceTypes.IMSBridge:                 Wrap(imsbridges.Service),
	export.ResourceTypes.CustomService:             Wrap(customservices.Service),
	export.ResourceTypes.AWSCredentials:            Wrap(aws.Service),
	export.ResourceTypes.AzureCredentials:          Wrap(azure.Service),
	export.ResourceTypes.SpanCaptureRule:           Wrap(capturing.Service),
	export.ResourceTypes.SpanContextPropagation:    Wrap(contextpropagation.Service),
	export.ResourceTypes.SLOV2:                     Wrap(slov2.Service),
	export.ResourceTypes.WebApplication:            Wrap(web.Service),
	export.ResourceTypes.MobileApplication:         Wrap(mobile.Service),
	export.ResourceTypes.JiraNotification:          Wrap(jira.Service),
	export.ResourceTypes.WebHookNotification:       Wrap(webhook.Service),
	export.ResourceTypes.AnsibleTowerNotification:  Wrap(ansible.Service),
	export.ResourceTypes.EmailNotification:         Wrap(email.Service),
	export.ResourceTypes.OpsGenieNotification:      Wrap(opsgenie.Service),
	export.ResourceTypes.PagerDutyNotification:     Wrap(pagerduty.Service),
	export.ResourceTypes.ServiceNowNotification:    Wrap(servicenow.Service),
	export.ResourceTypes.SlackNotification:         Wrap(slack.Service),
	export.ResourceTypes.TrelloNotification:        Wrap(trello.Service),
	export.ResourceTypes.VictorOpsNotification:     Wrap(victorops.Service),
	export.ResourceTypes.XMattersNotification:      Wrap(xmatters.Service),
	export.ResourceTypes.Maintenance:               Wrap(v2maintenance.Service),
	export.ResourceTypes.MetricEvents:              Wrap(metricevents.Service),
	export.ResourceTypes.KeyRequests:               Wrap(keyrequests.Service),
	export.ResourceTypes.Credentials:               Wrap(vault.Service),
	export.ResourceTypes.CalculatedServiceMetric:   Wrap(calculated_service_metrics.Service),
	export.ResourceTypes.CalculatedWebMetric:       Wrap(calculated_web_metrics.Service),
	export.ResourceTypes.CalculatedMobileMetric:    Wrap(calculated_mobile_metrics.Service),
	export.ResourceTypes.CalculatedSyntheticMetric: Wrap(calculated_synthetic_metrics.Service),
	export.ResourceTypes.HTTPMonitor:               Wrap(http.Service),
	export.ResourceTypes.BrowserMonitor:            Wrap(browser.Service),
	export.ResourceTypes.HostNaming:                Wrap(host_naming.Service),
	export.ResourceTypes.ProcessGroupNaming:        Wrap(processgroup_naming.Service),
	export.ResourceTypes.ServiceNaming:             Wrap(service_naming.Service),
	export.ResourceTypes.RequestNaming:             Wrap(requestnaming.Service),
}

var TodoServiceMap = map[export.ResourceType]ServiceFunc{
	export.ResourceTypes.Dashboard: Wrap(dashboards.Service),
	// export.ResourceTypes.Documents:                           Wrap(documents.Service),
	// export.ResourceTypes.DirectShares:                        Wrap(directshares.Service),
	export.ResourceTypes.ApplicationDetection:                Wrap(detection.Service),
	export.ResourceTypes.ApplicationErrorRules:               Wrap(errors.Service),
	export.ResourceTypes.SyntheticLocation:                   Wrap(locations.Service),
	export.ResourceTypes.QueueSharingGroups:                  Wrap(queuesharinggroup.Service),
	export.ResourceTypes.DDUPool:                             Wrap(ddupool.Service),
	export.ResourceTypes.ProcessGroupAlerting:                Wrap(processgroupalerting.Service),
	export.ResourceTypes.ServiceAnomaliesV2:                  Wrap(service_anomalies_v2.Service),
	export.ResourceTypes.DatabaseAnomaliesV2:                 Wrap(database_anomalies_v2.Service),
	export.ResourceTypes.ProcessMonitoringRule:               Wrap(customprocessmonitoring.Service),
	export.ResourceTypes.DiskAnomaliesV2:                     Wrap(disk_anomalies_v2.Service),
	export.ResourceTypes.DiskSpecificAnomaliesV2:             Wrap(disk_specific_anomalies_v2.Service),
	export.ResourceTypes.HostAnomaliesV2:                     Wrap(host_anomalies_v2.Service),
	export.ResourceTypes.CustomAppAnomalies:                  Wrap(custom_app_anomalies.Service),
	export.ResourceTypes.CustomAppCrashRate:                  Wrap(custom_app_crash_rate.Service),
	export.ResourceTypes.ProcessMonitoring:                   Wrap(processmonitoring.Service),
	export.ResourceTypes.ProcessAvailability:                 Wrap(processavailability.Service),
	export.ResourceTypes.AdvancedProcessGroupDetectionRule:   Wrap(advanceddetectionrule.Service),
	export.ResourceTypes.MobileAppAnomalies:                  Wrap(mobile_app_anomalies.Service),
	export.ResourceTypes.MobileAppCrashRate:                  Wrap(mobile_app_crash_rate.Service),
	export.ResourceTypes.WebAppAnomalies:                     Wrap(web_app_anomalies.Service),
	export.ResourceTypes.MutedRequests:                       Wrap(mutedrequests.Service),
	export.ResourceTypes.ConnectivityAlerts:                  Wrap(connectivityalerts.Service),
	export.ResourceTypes.DeclarativeGrouping:                 Wrap(declarativegrouping.Service),
	export.ResourceTypes.HostMonitoring:                      Wrap(hostmonitoring.Service),
	export.ResourceTypes.HostProcessGroupMonitoring:          Wrap(hostprocessgroupmonitoring.Service),
	export.ResourceTypes.RUMIPLocations:                      Wrap(ipmappings.Service),
	export.ResourceTypes.CustomAppEnablement:                 Wrap(rumcustomenablement.Service),
	export.ResourceTypes.MobileAppEnablement:                 Wrap(rummobileenablement.Service),
	export.ResourceTypes.WebAppEnablement:                    Wrap(rumwebenablement.Service),
	export.ResourceTypes.RUMProcessGroup:                     Wrap(rumprocessgroup.Service),
	export.ResourceTypes.RUMProviderBreakdown:                Wrap(rumproviderbreakdown.Service),
	export.ResourceTypes.UserExperienceScore:                 Wrap(userexperiencescore.Service),
	export.ResourceTypes.WebAppResourceCleanup:               Wrap(webappresourcecleanup.Service),
	export.ResourceTypes.UpdateWindows:                       Wrap(updatewindows.Service),
	export.ResourceTypes.ProcessGroupDetectionFlags:          Wrap(detectionflags.Service),
	export.ResourceTypes.ProcessGroupMonitoring:              Wrap(processgroupmonitoring.Service),
	export.ResourceTypes.ProcessGroupSimpleDetection:         Wrap(processgroupsimpledetection.Service),
	export.ResourceTypes.LogMetrics:                          Wrap(schemalesslogmetric.Service),
	export.ResourceTypes.BrowserMonitorPerformanceThresholds: Wrap(browserperformancethresholds.Service),
	export.ResourceTypes.HttpMonitorPerformanceThresholds:    Wrap(httpperformancethresholds.Service),
	export.ResourceTypes.HttpMonitorCookies:                  Wrap(httpcookies.Service),
	export.ResourceTypes.SessionReplayWebPrivacy:             Wrap(sessionreplaywebprivacy.Service),
	export.ResourceTypes.SessionReplayResourceCapture:        Wrap(sessionreplayresourcecapture.Service),
	export.ResourceTypes.UsabilityAnalytics:                  Wrap(analytics.Service),
	export.ResourceTypes.SyntheticAvailability:               Wrap(availability.Service),
	export.ResourceTypes.BrowserMonitorOutageHandling:        Wrap(browseroutagehandling.Service),
	export.ResourceTypes.HttpMonitorOutageHandling:           Wrap(httpoutagehandling.Service),
	export.ResourceTypes.CloudAppWorkloadDetection:           Wrap(workloaddetection.Service),
	export.ResourceTypes.MainframeTransactionMonitoring:      Wrap(txmonitoring.Service),
	export.ResourceTypes.MonitoredTechnologiesApache:         Wrap(apache.Service),
	export.ResourceTypes.MonitoredTechnologiesDotNet:         Wrap(dotnet.Service),
	export.ResourceTypes.MonitoredTechnologiesGo:             Wrap(golang.Service),
	export.ResourceTypes.MonitoredTechnologiesIIS:            Wrap(iis.Service),
	export.ResourceTypes.MonitoredTechnologiesJava:           Wrap(java.Service),
	export.ResourceTypes.MonitoredTechnologiesNGINX:          Wrap(nginx.Service),
	export.ResourceTypes.MonitoredTechnologiesNodeJS:         Wrap(nodejs.Service),
	export.ResourceTypes.MonitoredTechnologiesOpenTracing:    Wrap(opentracingnative.Service),
	export.ResourceTypes.MonitoredTechnologiesPHP:            Wrap(php.Service),
	export.ResourceTypes.MonitoredTechnologiesVarnish:        Wrap(varnish.Service),
	export.ResourceTypes.MonitoredTechnologiesWSMB:           Wrap(wsmb.Service),
	export.ResourceTypes.ProcessVisibility:                   Wrap(processvisibility.Service),
	export.ResourceTypes.RUMHostHeaders:                      Wrap(hostheaders.Service),
	export.ResourceTypes.RUMIPDetermination:                  Wrap(ipdetermination.Service),
	export.ResourceTypes.TransactionStartFilters:             Wrap(txstartfilters.Service),
	export.ResourceTypes.OneAgentFeatures:                    Wrap(features.Service),
	export.ResourceTypes.RUMOverloadPrevention:               Wrap(overloadprevention.Service),
	export.ResourceTypes.RUMAdvancedCorrelation:              Wrap(resourcetimingorigins.Service),
	export.ResourceTypes.WebAppBeaconOrigins:                 Wrap(beacondomainorigins.Service),
	export.ResourceTypes.WebAppResourceTypes:                 Wrap(resourcetypes.Service),
	export.ResourceTypes.GenericTypes:                        Wrap(generictypes.Service),
	export.ResourceTypes.GenericRelationships:                Wrap(relation.Service),
	export.ResourceTypes.SLONormalization:                    Wrap(normalization.Service),
	export.ResourceTypes.DataPrivacy:                         Wrap(privacy.Service),
	export.ResourceTypes.ServiceFailure:                      Wrap(generalparameters.Service),
	export.ResourceTypes.ServiceHTTPFailure:                  Wrap(httpparameters.Service),
	export.ResourceTypes.DiskOptions:                         Wrap(diskoptions.Service),
	export.ResourceTypes.OSServices:                          Wrap(osservicesmonitoring.Service),
	export.ResourceTypes.ExtensionExecutionController:        Wrap(local.Service),
	export.ResourceTypes.NetTracerTraffic:                    Wrap(traffic.Service),
	export.ResourceTypes.AIXExtension:                        Wrap(aixkernelextension.Service),
	export.ResourceTypes.ActiveGateToken:                     Wrap(activegatetoken.Service),
	export.ResourceTypes.AuditLog:                            Wrap(auditlog.Service),
	export.ResourceTypes.K8sClusterAnomalies:                 Wrap(cluster.Service),
	export.ResourceTypes.K8sNamespaceAnomalies:               Wrap(namespace.Service),
	export.ResourceTypes.K8sNodeAnomalies:                    Wrap(node.Service),
	export.ResourceTypes.K8sWorkloadAnomalies:                Wrap(workload.Service),
	export.ResourceTypes.ContainerBuiltinRule:                Wrap(builtinmonitoringrule.Service),
	export.ResourceTypes.ContainerRule:                       Wrap(monitoringrule.Service),
	export.ResourceTypes.ContainerTechnology:                 Wrap(containertechnology.Service),
	export.ResourceTypes.RemoteEnvironments:                  Wrap(environment.Service),
	export.ResourceTypes.DashboardsGeneral:                   Wrap(dashboardsgeneral.Service),
	export.ResourceTypes.DashboardsPresets:                   Wrap(dashboardspresets.Service),
	export.ResourceTypes.LogProcessing:                       Wrap(logdpprules.Service),
	export.ResourceTypes.LogEvents:                           Wrap(logevents.Service),
	export.ResourceTypes.LogTimestamp:                        Wrap(timestampconfiguration.Service),
	export.ResourceTypes.LogCustomAttribute:                  Wrap(logcustomattributes.Service),
	export.ResourceTypes.LogSensitiveDataMasking:             Wrap(sensitivedatamasking.Service),
	export.ResourceTypes.LogStorage:                          Wrap(logstoragesettings.Service),
	export.ResourceTypes.LogBuckets:                          Wrap(logbucketsrules.Service),
	export.ResourceTypes.LogSecurityContext:                  Wrap(logsecuritycontextrules.Service),
	export.ResourceTypes.EULASettings:                        Wrap(eulasettings.Service),
	export.ResourceTypes.APIDetectionRules:                   Wrap(apidetection.Service),
	export.ResourceTypes.ServiceExternalWebRequest:           Wrap(externalwebrequest.Service),
	export.ResourceTypes.ServiceExternalWebService:           Wrap(externalwebservice.Service),
	export.ResourceTypes.ServiceFullWebRequest:               Wrap(fullwebrequest.Service),
	export.ResourceTypes.ServiceFullWebService:               Wrap(fullwebservice.Service),
	export.ResourceTypes.DashboardsAllowlist:                 Wrap(dashboardsallowlist.Service),
	export.ResourceTypes.FailureDetectionParameters:          Wrap(envparameters.Service),
	export.ResourceTypes.FailureDetectionRules:               Wrap(envrules.Service),
	export.ResourceTypes.LogOneAgent:                         Wrap(logagentconfiguration.Service),
	export.ResourceTypes.IssueTracking:                       Wrap(issuetracking.Service),
	export.ResourceTypes.GeolocationSettings:                 Wrap(geosettings.Service),
	export.ResourceTypes.UserSessionCustomMetrics:            Wrap(custommetrics.Service),
	export.ResourceTypes.CustomUnits:                         Wrap(customunit.Service),
	export.ResourceTypes.DiskAnalytics:                       Wrap(diskanalytics.Service),
	export.ResourceTypes.NetworkTraffic:                      Wrap(networktraffic.Service),
	export.ResourceTypes.TokenSettings:                       Wrap(tokensettings.Service),
	export.ResourceTypes.ExtensionExecutionRemote:            Wrap(eecremote.Service),
	export.ResourceTypes.K8sPVCAnomalies:                     Wrap(pvc.Service),
	export.ResourceTypes.UserActionCustomMetrics:             Wrap(useractioncustommetrics.Service),
	export.ResourceTypes.WebAppJavascriptVersion:             Wrap(customrumjavascriptversion.Service),
	export.ResourceTypes.WebAppJavascriptUpdates:             Wrap(rumjavascriptupdates.Service),
	export.ResourceTypes.OpenTelemetryMetrics:                Wrap(opentelemetrymetrics.Service),
	export.ResourceTypes.ActiveGateUpdates:                   Wrap(activegateupdates.Service),
	export.ResourceTypes.OneAgentDefaultVersion:              Wrap(defaultversion.Service),
	export.ResourceTypes.OneAgentUpdates:                     Wrap(oneagentupdates.Service),
	export.ResourceTypes.OwnershipTeams:                      Wrap(teams.Service),
	export.ResourceTypes.OwnershipConfig:                     Wrap(ownership_config.Service),
	export.ResourceTypes.LogCustomSource:                     Wrap(customlogsourcesettings.Service),
	export.ResourceTypes.Kubernetes:                          Wrap(kubernetesv2.Service),
	export.ResourceTypes.CloudFoundry:                        Wrap(cloudfoundryv2.Service),
	export.ResourceTypes.DiskAnomalyDetectionRules:           Wrap(diskrules.Service),
	export.ResourceTypes.AWSAnomalies:                        Wrap(aws_anomalies.Service),
	export.ResourceTypes.VMwareAnomalies:                     Wrap(vmware_anomalies.Service),
	export.ResourceTypes.BusinessEventsOneAgent:              Wrap(incoming.Service),
	export.ResourceTypes.BusinessEventsBuckets:               Wrap(bizevents_buckets.Service),
	export.ResourceTypes.BusinessEventsMetrics:               Wrap(bizevents_metrics.Service),
	export.ResourceTypes.BusinessEventsProcessing:            Wrap(bizevents_processing.Service),
	export.ResourceTypes.BusinessEventsSecurityContext:       Wrap(bizevents_security.Service),
	export.ResourceTypes.WebAppKeyPerformanceCustom:          Wrap(customactions.Service),
	export.ResourceTypes.WebAppKeyPerformanceLoad:            Wrap(loadactions.Service),
	export.ResourceTypes.WebAppKeyPerformanceXHR:             Wrap(xhractions.Service),
	export.ResourceTypes.BuiltinProcessMonitoring:            Wrap(builtinprocessmonitoringrule.Service),
	export.ResourceTypes.LimitOutboundConnections:            Wrap(allowedoutboundconnections.Service),
	export.ResourceTypes.SpanEvents:                          Wrap(eventattribute.Service),
	export.ResourceTypes.VMware:                              Wrap(vmware.Service),
	export.ResourceTypes.K8sMonitoring:                       Wrap(kubernetesmonitoring.Service),
	// export.ResourceTypes.AutomationWorkflow:                  Wrap(workflows.Service),
	// export.ResourceTypes.AutomationBusinessCalendar:          Wrap(business_calendars.Service),
	// export.ResourceTypes.AutomationSchedulingRule:            Wrap(scheduling_rules.Service),
	export.ResourceTypes.HostMonitoringMode:     Wrap(hostmonitoringmode.Service),
	export.ResourceTypes.HostMonitoringAdvanced: Wrap(hostmonitoringadvanced.Service),
	export.ResourceTypes.IPAddressMasking:       Wrap(ipaddressmasking.Service),
	export.ResourceTypes.UnifiedServicesMetrics: Wrap(endpointmetrics.Service),
	export.ResourceTypes.UnifiedServicesOpenTel: Wrap(unifiedservicesopentel.Service),
	// export.ResourceTypes.PlatformBucket:                Wrap(platformbuckets.Service),
	export.ResourceTypes.KeyUserAction:                Wrap(keyuseractions.Service),
	export.ResourceTypes.UrlBasedSampling:             Wrap(urlbasedsampling.Service),
	export.ResourceTypes.AttributeAllowList:           Wrap(attributeallowlist.Service),
	export.ResourceTypes.AttributeBlockList:           Wrap(attributeblocklist.Service),
	export.ResourceTypes.AttributeMasking:             Wrap(attributemasking.Service),
	export.ResourceTypes.AttributesPreferences:        Wrap(attributespreferences.Service),
	export.ResourceTypes.OneAgentSideMasking:          Wrap(masking.Service),
	export.ResourceTypes.HubSubscriptions:             Wrap(subscriptions.Service),
	export.ResourceTypes.MobileNotifications:          Wrap(mobilenotifications.Service),
	export.ResourceTypes.CrashdumpAnalytics:           Wrap(crashdumpanalytics.Service),
	export.ResourceTypes.AppMonitoring:                Wrap(appmonitoring.Service),
	export.ResourceTypes.GrailSecurityContext:         Wrap(securitycontext.Service),
	export.ResourceTypes.KubernetesApp:                Wrap(kubernetesapp.Service),
	export.ResourceTypes.WebAppBeaconEndpoint:         Wrap(webappbeaconendpoint.Service),
	export.ResourceTypes.WebAppCustomConfigProperties: Wrap(webappcustomconfigproperties.Service),
	export.ResourceTypes.WebAppInjectionCookie:        Wrap(webappinjectioncookie.Service),
	export.ResourceTypes.DatabaseAppFeatureFlags:      Wrap(dbfeatureflags.Service),
	export.ResourceTypes.InfraOpsAppFeatureFlags:      Wrap(infraopsfeatureflags.Service),
	export.ResourceTypes.EBPFServiceDiscovery:         Wrap(ebpf.Service),
	export.ResourceTypes.DavisAnomalyDetectors:        Wrap(anomalydetectors.Service),
	export.ResourceTypes.LogDebugSettings:             Wrap(logdebugsettings.Service),
	export.ResourceTypes.InfraOpsAppSettings:          Wrap(infraopssettings.Service),
	export.ResourceTypes.DiskEdgeAnomalyDetectors:     Wrap(diskedgeanomalydetectors.Service),
	export.ResourceTypes.Reports:                      Wrap(reports.Service),
	export.ResourceTypes.NetworkMonitor:               Wrap(v2monitors.Service),
	export.ResourceTypes.NetworkMonitorOutageHandling: Wrap(networkoutagehandling.Service),
	export.ResourceTypes.NetworkZone:                  Wrap(networkzone.Service),
}

// Settings where only a single instance exists
var SingleConfigServiceMap = map[export.ResourceType]ServiceFunc{
	export.ResourceTypes.IBMMQFilters:   Wrap(mqfilters.Service),
	export.ResourceTypes.FrequentIssues: Wrap(frequentissues.Service),
	export.ResourceTypes.NetworkZones:   Wrap(networkzones.Service),
}

// Freshly provisioned environments are pre-populated with settings that CAN get deleted
var PrePopulatedConfigServiceMap = map[export.ResourceType]ServiceFunc{
	export.ResourceTypes.SpanEntryPoint: Wrap(entrypoints.Service),
}

// Resources that are deprecated
var DeprecatedConfigServiceMap = map[export.ResourceType]ServiceFunc{
	export.ResourceTypes.SpanAttribute:      Wrap(attribute.Service),
	export.ResourceTypes.ResourceAttributes: Wrap(resourceattribute.Service),
}
