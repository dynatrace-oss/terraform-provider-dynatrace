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

package builtinprocessmonitoringrule

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	HostGroupID *string `json:"-" scope:"hostGroupId"` // The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.
	RuleID1     bool    `json:"-1"`                    // (v1.274) Rule id: 1 - Do not monitor processes if PHP script exists
	RuleID2     bool    `json:"-2"`                    // Rule id: 2 - Do not monitor processes if EXE name equals 'php-cgi'
	RuleID3     bool    `json:"-3"`                    // Rule id: 3 - Do monitor processes if ASP.NET Core application path exists
	RuleID4     bool    `json:"-4"`                    // Rule id: 4 - Do monitor processes if EXE name equals 'w3wp.exe'
	RuleID5     bool    `json:"-5"`                    // Rule id: 5 - Do monitor processes if EXE name equals 'caddy'
	RuleID6     bool    `json:"-6"`                    // Rule id: 6 - Do monitor processes if EXE name equals 'influxd'
	RuleID7     bool    `json:"-7"`                    // Rule id: 7 - Do monitor processes if EXE name equals 'adapter'
	RuleID8     bool    `json:"-8"`                    // Rule id: 8 - Do monitor processes if EXE name equals 'auctioneer'
	RuleID9     bool    `json:"-9"`                    // Rule id: 9 - Do monitor processes if EXE name equals 'bbs'
	RuleID10    bool    `json:"-10"`                   // Rule id: 10 - Do monitor processes if EXE name equals 'cc-uploader'
	RuleID11    bool    `json:"-11"`                   // Rule id: 11 - Do monitor processes if EXE name equals 'doppler'
	RuleID12    bool    `json:"-12"`                   // Rule id: 12 - Do monitor processes if EXE name equals 'gorouter'
	RuleID13    bool    `json:"-13"`                   // Rule id: 13 - Do monitor processes if EXE name equals 'locket'
	RuleID14    bool    `json:"-14"`                   // Rule id: 14 - Do monitor processes if EXE name equals 'metron'
	RuleID16    bool    `json:"-16"`                   // Rule id: 16 - Do monitor processes if EXE name equals 'rep'
	RuleID17    bool    `json:"-17"`                   // Rule id: 17 - Do monitor processes if EXE name equals 'route-emitter'
	RuleID18    bool    `json:"-18"`                   // Rule id: 18 - Do monitor processes if EXE name equals 'route-registrar'
	RuleID19    bool    `json:"-19"`                   // Rule id: 19 - Do monitor processes if EXE name equals 'routing-api'
	RuleID20    bool    `json:"-20"`                   // Rule id: 20 - Do monitor processes if EXE name equals 'scheduler'
	RuleID21    bool    `json:"-21"`                   // Rule id: 21 - Do monitor processes if EXE name equals 'silk-daemon'
	RuleID22    bool    `json:"-22"`                   // Rule id: 22 - Do monitor processes if EXE name equals 'switchboard'
	RuleID23    bool    `json:"-23"`                   // Rule id: 23 - Do monitor processes if EXE name equals 'syslogRuleIDdrainRuleIDbinder'
	RuleID24    bool    `json:"-24"`                   // Rule id: 24 - Do monitor processes if EXE name equals 'tps-watcher'
	RuleID25    bool    `json:"-25"`                   // Rule id: 25 - Do monitor processes if EXE name equals 'trafficcontroller'
	RuleID26    bool    `json:"-26"`                   // Rule id: 26 - Do not monitor processes if Node.js application base directory ends with '/nodeRuleIDmodules/prebuild-install'
	RuleID27    bool    `json:"-27"`                   // Rule id: 27 - Do not monitor processes if Node.js application base directory ends with '/nodeRuleIDmodules/npm'
	RuleID28    bool    `json:"-28"`                   // Rule id: 28 - Do not monitor processes if Node.js application base directory ends with '/nodeRuleIDmodules/grunt'
	RuleID29    bool    `json:"-29"`                   // Rule id: 29 - Do not monitor processes if Node.js application base directory ends with '/nodeRuleIDmodules/typescript'
	RuleID32    bool    `json:"-32"`                   // Rule id: 32 - Do not monitor processes if Node.js application base directory ends with '/nodeRuleIDmodules/node-pre-gyp'
	RuleID33    bool    `json:"-33"`                   // Rule id: 33 - Do not monitor processes if Node.js application base directory ends with '/nodeRuleIDmodules/node-gyp'
	RuleID34    bool    `json:"-34"`                   // Rule id: 34 - Do not monitor processes if Node.js application base directory ends with '/nodeRuleIDmodules/gulp-cli'
	RuleID35    bool    `json:"-35"`                   // Rule id: 35 - Do not monitor processes if Node.js script equals 'bin/pm2'
	RuleID36    bool    `json:"-36"`                   // Rule id: 36 - Do not monitor processes if Cloud Foundry application begins with 'apps-manager-js'
	RuleID37    bool    `json:"-37"`                   // Rule id: 37 - Do monitor processes if Cloud Foundry application exists
	RuleID38    bool    `json:"-38"`                   // Rule id: 38 - Do not monitor processes if Kubernetes container name equals 'POD'
	RuleID39    bool    `json:"-39"`                   // Rule id: 39 - Do not monitor processes if Docker stripped image contains 'pause-amd64'
	RuleID40    bool    `json:"-40"`                   // Rule id: 40 - Do monitor processes if Kubernetes namespace exists
	RuleID41    bool    `json:"-41"`                   // Rule id: 41 - Do monitor processes if container name exists
	RuleID43    bool    `json:"-43"`                   // Rule id: 43 - Do not monitor processes if EXE path begins with '/tmp/buildpacks/'
	RuleID44    bool    `json:"-44"`                   // Rule id: 44 - Do not monitor processes if EXE name equals 'oc'
	RuleID45    bool    `json:"-45"`                   // Rule id: 45 - Do not monitor processes if Node.js application equals 'yarn'
	RuleID46    bool    `json:"-46"`                   // Rule id: 46 - Do not monitor processes if EXE path equals '/opt/cni/bin/host-local'
	RuleID47    bool    `json:"-47"`                   // Rule id: 47 - Do not monitor processes if Go Binary Linkage equals 'static'
	RuleID48    bool    `json:"-48"`                   // Rule id: 48 - Do not monitor processes if EXE name begins with 'mqsi'
	RuleID49    bool    `json:"-49"`                   // Rule id: 49 - Do not monitor processes if EXE name equals 'filebeat'
	RuleID50    bool    `json:"-50"`                   // Rule id: 50 - Do not monitor processes if EXE name equals 'metricbeat'
	RuleID51    bool    `json:"-51"`                   // Rule id: 51 - Do not monitor processes if EXE name equals 'packetbeat'
	RuleID52    bool    `json:"-52"`                   // Rule id: 52 - Do not monitor processes if EXE name equals 'auditbeat'
	RuleID53    bool    `json:"-53"`                   // Rule id: 53 - Do not monitor processes if EXE name equals 'heartbeat'
	RuleID54    bool    `json:"-54"`                   // Rule id: 54 - Do not monitor processes if EXE name equals 'functionbeat'
	RuleID55    bool    `json:"-55"`                   // Rule id: 55 - Do not monitor processes if EXE name equals 'grootfs'
	RuleID56    bool    `json:"-56"`                   // Rule id: 56 - Do not monitor processes if EXE name equals 'tardis'
	RuleID57    bool    `json:"-57"`                   // Rule id: 57 - Do not monitor processes if Java JAR file begins with 'org.eclipse.equinox.launcher'
	RuleID58    bool    `json:"-58"`                   // Rule id: 58 - Do not monitor processes if EXE name equals 'calico-node'
	RuleID59    bool    `json:"-59"`                   // Rule id: 59 - Do not monitor processes if EXE name equals 'casclient.exe'
	RuleID60    bool    `json:"-60"`                   // Rule id: 60 - Do not monitor processes if JAR file name equals 'dynatrace_ibm_mq_connector.jar'
	RuleID61    bool    `json:"-61"`                   // Rule id: 61 - Do not monitor processes if EXE name contains 'Agent.Worker'
	RuleID62    bool    `json:"-62"`                   // Rule id: 62 - Do not monitor processes if ASP.NET Core application DLL contains 'Agent.Worker'
	RuleID63    bool    `json:"-63"`                   // Rule id: 63 - Do not monitor processes if EXE name contains 'Agent.Listener'
	RuleID64    bool    `json:"-64"`                   // Rule id: 64 - Do not monitor processes if ASP.NET Core application DLL contains 'Agent.Listener'
	RuleID65    bool    `json:"-65"`                   // Rule id: 65 - Do not monitor processes if EXE name equals 'FlexNetJobExecutorService'
	RuleID66    bool    `json:"-66"`                   // Rule id: 66 - Do not monitor processes if EXE name equals 'FlexNetMaintenanceRemotingService'
	RuleID67    bool    `json:"-67"`                   // Rule id: 67 - Do not monitor processes if EXE path equals '/usr/bin/piper'
	RuleID68    bool    `json:"-68"`                   // Rule id: 68 - Do not monitor processes if Node.js application equals 'corepack'
	RuleID69    bool    `json:"-69"`                   // Rule id: 69 - Do not monitor processes if Kubernetes container name equals 'cassandra-operator'
	RuleID70    bool    `json:"-70"`                   // Rule id: 70 - Do not monitor processes if EXE name contains 'UiPath'
}

func (me *Settings) Name() string {
	return *me.HostGroupID
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"host_group_id": {
			Type:        schema.TypeString,
			Description: "The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.",
			Optional:    true,
			Default:     "environment",
			ForceNew:    true,
		},
		"php_script": {
			Type:        schema.TypeBool,
			Description: "(v1.274) Rule id: 1 - Do not monitor processes if PHP script exists",
			Optional:    true,
			Default:     true,
		},
		"exe_phpcgi": {
			Type:        schema.TypeBool,
			Description: "Rule id: 2 - Do not monitor processes if EXE name equals 'php-cgi'",
			Optional:    true,
			Default:     true,
		},
		"aspnetcore": {
			Type:        schema.TypeBool,
			Description: "Rule id: 3 - Do monitor processes if ASP.NET Core application path exists",
			Optional:    true,
			Default:     true,
		},
		"exe_w3wp": {
			Type:        schema.TypeBool,
			Description: "Rule id: 4 - Do monitor processes if EXE name equals 'w3wp.exe'",
			Optional:    true,
			Default:     true,
		},
		"exe_caddy": {
			Type:        schema.TypeBool,
			Description: "Rule id: 5 - Do monitor processes if EXE name equals 'caddy'",
			Optional:    true,
			Default:     true,
		},
		"exe_influxd": {
			Type:        schema.TypeBool,
			Description: "Rule id: 6 - Do monitor processes if EXE name equals 'influxd'",
			Optional:    true,
			Default:     true,
		},
		"exe_adapter": {
			Type:        schema.TypeBool,
			Description: "Rule id: 7 - Do monitor processes if EXE name equals 'adapter'",
			Optional:    true,
			Default:     true,
		},
		"exe_auctioneer": {
			Type:        schema.TypeBool,
			Description: "Rule id: 8 - Do monitor processes if EXE name equals 'auctioneer'",
			Optional:    true,
			Default:     true,
		},
		"exe_bbs": {
			Type:        schema.TypeBool,
			Description: "Rule id: 9 - Do monitor processes if EXE name equals 'bbs'",
			Optional:    true,
			Default:     true,
		},
		"exe_ccuploader": {
			Type:        schema.TypeBool,
			Description: "Rule id: 10 - Do monitor processes if EXE name equals 'cc-uploader'",
			Optional:    true,
			Default:     true,
		},
		"exe_doppler": {
			Type:        schema.TypeBool,
			Description: "Rule id: 11 - Do monitor processes if EXE name equals 'doppler'",
			Optional:    true,
			Default:     true,
		},
		"exe_gorouter": {
			Type:        schema.TypeBool,
			Description: "Rule id: 12 - Do monitor processes if EXE name equals 'gorouter'",
			Optional:    true,
			Default:     true,
		},
		"exe_locket": {
			Type:        schema.TypeBool,
			Description: "Rule id: 13 - Do monitor processes if EXE name equals 'locket'",
			Optional:    true,
			Default:     true,
		},
		"exe_metron": {
			Type:        schema.TypeBool,
			Description: "Rule id: 14 - Do monitor processes if EXE name equals 'metron'",
			Optional:    true,
			Default:     true,
		},
		"exe_rep": {
			Type:        schema.TypeBool,
			Description: "Rule id: 16 - Do monitor processes if EXE name equals 'rep'",
			Optional:    true,
			Default:     true,
		},
		"exe_routeemitter": {
			Type:        schema.TypeBool,
			Description: "Rule id: 17 - Do monitor processes if EXE name equals 'route-emitter'",
			Optional:    true,
			Default:     true,
		},
		"exe_routeregistrar": {
			Type:        schema.TypeBool,
			Description: "Rule id: 18 - Do monitor processes if EXE name equals 'route-registrar'",
			Optional:    true,
			Default:     true,
		},
		"exe_routingapi": {
			Type:        schema.TypeBool,
			Description: "Rule id: 19 - Do monitor processes if EXE name equals 'routing-api'",
			Optional:    true,
			Default:     true,
		},
		"exe_schedular": {
			Type:        schema.TypeBool,
			Description: "Rule id: 20 - Do monitor processes if EXE name equals 'scheduler'",
			Optional:    true,
			Default:     true,
		},
		"exe_silkdaemon": {
			Type:        schema.TypeBool,
			Description: "Rule id: 21 - Do monitor processes if EXE name equals 'silk-daemon'",
			Optional:    true,
			Default:     true,
		},
		"exe_switchboard": {
			Type:        schema.TypeBool,
			Description: "Rule id: 22 - Do monitor processes if EXE name equals 'switchboard'",
			Optional:    true,
			Default:     true,
		},
		"exe_syslogdrainbinder": {
			Type:        schema.TypeBool,
			Description: "Rule id: 23 - Do monitor processes if EXE name equals 'syslog_drain_binder'",
			Optional:    true,
			Default:     true,
		},
		"exe_tpswatcher": {
			Type:        schema.TypeBool,
			Description: "Rule id: 24 - Do monitor processes if EXE name equals 'tps-watcher'",
			Optional:    true,
			Default:     true,
		},
		"exe_trafficcontroller": {
			Type:        schema.TypeBool,
			Description: "Rule id: 25 - Do monitor processes if EXE name equals 'trafficcontroller'",
			Optional:    true,
			Default:     true,
		},
		"node_prebuildinstall": {
			Type:        schema.TypeBool,
			Description: "Rule id: 26 - Do not monitor processes if Node.js application base directory ends with '/node_modules/prebuild-install'",
			Optional:    true,
			Default:     true,
		},
		"node_npm": {
			Type:        schema.TypeBool,
			Description: "Rule id: 27 - Do not monitor processes if Node.js application base directory ends with '/node_modules/npm'",
			Optional:    true,
			Default:     true,
		},
		"node_grunt": {
			Type:        schema.TypeBool,
			Description: "Rule id: 28 - Do not monitor processes if Node.js application base directory ends with '/node_modules/grunt'",
			Optional:    true,
			Default:     true,
		},
		"node_typescript": {
			Type:        schema.TypeBool,
			Description: "Rule id: 29 - Do not monitor processes if Node.js application base directory ends with '/node_modules/typescript'",
			Optional:    true,
			Default:     true,
		},
		"node_nodepregyp": {
			Type:        schema.TypeBool,
			Description: "Rule id: 32 - Do not monitor processes if Node.js application base directory ends with '/node_modules/node-pre-gyp'",
			Optional:    true,
			Default:     true,
		},
		"node_nodegyp": {
			Type:        schema.TypeBool,
			Description: "Rule id: 33 - Do not monitor processes if Node.js application base directory ends with '/node_modules/node-gyp'",
			Optional:    true,
			Default:     true,
		},
		"node_gulpcli": {
			Type:        schema.TypeBool,
			Description: "Do not monitor processes if Node.js application base directory ends with '/node_modules/gulp-cli'",
			Optional:    true,
			Default:     true,
		},
		"node_binpm2": {
			Type:        schema.TypeBool,
			Description: "Do not monitor processes if Node.js script equals 'bin/pm2'",
			Optional:    true,
			Default:     true,
		},
		"cf_appsmanagerjs": {
			Type:        schema.TypeBool,
			Description: "Do not monitor processes if Cloud Foundry application begins with 'apps-manager-js'",
			Optional:    true,
			Default:     true,
		},
		"cf": {
			Type:        schema.TypeBool,
			Description: "Rule id: 37 - Do monitor processes if Cloud Foundry application exists",
			Optional:    true,
			Default:     true,
		},
		"k8s_containerpod": {
			Type:        schema.TypeBool,
			Description: "Rule id: 38 - Do not monitor processes if Kubernetes container name equals 'POD'",
			Optional:    true,
			Default:     true,
		},
		"docker_pauseamd64": {
			Type:        schema.TypeBool,
			Description: "Rule id: 39 - Do not monitor processes if Docker stripped image contains 'pause-amd64'",
			Optional:    true,
			Default:     true,
		},
		"k8s_namespace": {
			Type:        schema.TypeBool,
			Description: "Rule id: 40 - Do monitor processes if Kubernetes namespace exists",
			Optional:    true,
			Default:     true,
		},
		"container": {
			Type:        schema.TypeBool,
			Description: "Rule id: 41 - Do monitor processes if container name exists",
			Optional:    true,
			Default:     true,
		},
		"exe_tmpbuildpacks": {
			Type:        schema.TypeBool,
			Description: "Rule id: 43 - Do not monitor processes if EXE path begins with '/tmp/buildpacks/'",
			Optional:    true,
			Default:     true,
		},
		"exe_oc": {
			Type:        schema.TypeBool,
			Description: "Rule id: 44 - Do not monitor processes if EXE name equals 'oc'",
			Optional:    true,
			Default:     true,
		},
		"node_yarn": {
			Type:        schema.TypeBool,
			Description: "Rule id: 45 - Do not monitor processes if Node.js application equals 'yarn'",
			Optional:    true,
			Default:     true,
		},
		"exe_optcnibinhostlocal": {
			Type:        schema.TypeBool,
			Description: "Rule id: 46 - Do not monitor processes if EXE path equals '/opt/cni/bin/host-local'",
			Optional:    true,
			Default:     true,
		},
		"go_static": {
			Type:        schema.TypeBool,
			Description: "Rule id: 47 - Do not monitor processes if Go Binary Linkage equals 'static'",
			Optional:    true,
			Default:     true,
		},
		"exe_mqsi": {
			Type:        schema.TypeBool,
			Description: "Rule id: 48 - Do not monitor processes if EXE name begins with 'mqsi'",
			Optional:    true,
			Default:     true,
		},
		"exe_filebeat": {
			Type:        schema.TypeBool,
			Description: "Rule id: 49 - Do not monitor processes if EXE name equals 'filebeat'",
			Optional:    true,
			Default:     true,
		},
		"exe_metricbeat": {
			Type:        schema.TypeBool,
			Description: "Rule id: 50 - Do not monitor processes if EXE name equals 'metricbeat'",
			Optional:    true,
			Default:     true,
		},
		"exe_packetbeat": {
			Type:        schema.TypeBool,
			Description: "Rule id: 51 - Do not monitor processes if EXE name equals 'packetbeat'",
			Optional:    true,
			Default:     true,
		},
		"exe_auditbeat": {
			Type:        schema.TypeBool,
			Description: "Rule id: 52 - Do not monitor processes if EXE name equals 'auditbeat'",
			Optional:    true,
			Default:     true,
		},
		"exe_heartbeat": {
			Type:        schema.TypeBool,
			Description: "Rule id: 53 - Do not monitor processes if EXE name equals 'heartbeat'",
			Optional:    true,
			Default:     true,
		},
		"exe_functionbeat": {
			Type:        schema.TypeBool,
			Description: "Rule id: 54 - Do not monitor processes if EXE name equals 'functionbeat'",
			Optional:    true,
			Default:     true,
		},
		"exe_grootfs": {
			Type:        schema.TypeBool,
			Description: "Rule id: 55 - Do not monitor processes if EXE name equals 'grootfs'",
			Optional:    true,
			Default:     true,
		},
		"exe_tardis": {
			Type:        schema.TypeBool,
			Description: "Rule id: 56 - Do not monitor processes if EXE name equals 'tardis'",
			Optional:    true,
			Default:     true,
		},
		"jar_eclipseequinox": {
			Type:        schema.TypeBool,
			Description: "Rule id: 57 - Do not monitor processes if Java JAR file begins with 'org.eclipse.equinox.launcher'",
			Optional:    true,
			Default:     true,
		},
		"exe_caliconode": {
			Type:        schema.TypeBool,
			Description: "Rule id: 58 - Do not monitor processes if EXE name equals 'calico-node'",
			Optional:    true,
			Default:     true,
		},
		"exe_casclient": {
			Type:        schema.TypeBool,
			Description: "Rule id: 59 - Do not monitor processes if EXE name equals 'casclient.exe'",
			Optional:    true,
			Default:     true,
		},
		"jar_dtibmmqconnector": {
			Type:        schema.TypeBool,
			Description: "Rule id: 60 - Do not monitor processes if JAR file name equals 'dynatrace_ibm_mq_connector.jar'",
			Optional:    true,
			Default:     true,
		},
		"exe_agentworker": {
			Type:        schema.TypeBool,
			Description: "Rule id: 61 - Do not monitor processes if EXE name contains 'Agent.Worker'",
			Optional:    true,
			Default:     true,
		},
		"aspnetcore_agentworker": {
			Type:        schema.TypeBool,
			Description: "Rule id: 62 - Do not monitor processes if ASP.NET Core application DLL contains 'Agent.Worker'",
			Optional:    true,
			Default:     true,
		},
		"exe_agentlistener": {
			Type:        schema.TypeBool,
			Description: "Rule id: 63 - Do not monitor processes if EXE name contains 'Agent.Listener'",
			Optional:    true,
			Default:     true,
		},
		"aspnetcore_agentlistener": {
			Type:        schema.TypeBool,
			Description: "Rule id: 64 - Do not monitor processes if ASP.NET Core application DLL contains 'Agent.Listener'",
			Optional:    true,
			Default:     true,
		},
		"exe_flexnetjobexecutorservice": {
			Type:        schema.TypeBool,
			Description: "Rule id: 65 - Do not monitor processes if EXE name equals 'FlexNetJobExecutorService'",
			Optional:    true,
			Default:     true,
		},
		"exe_flexnetmaintenanceremotingservice": {
			Type:        schema.TypeBool,
			Description: "Rule id: 66 - Do not monitor processes if EXE name equals 'FlexNetMaintenanceRemotingService'",
			Optional:    true,
			Default:     true,
		},
		"exe_userbinpiper": {
			Type:        schema.TypeBool,
			Description: "Rule id: 67 - Do not monitor processes if EXE path equals '/usr/bin/piper'",
			Optional:    true,
			Default:     true,
		},
		"node_corepack": {
			Type:        schema.TypeBool,
			Description: "Rule id: 68 - Do not monitor processes if Node.js application equals 'corepack'",
			Optional:    true,
			Default:     true,
		},
		"k8s_cassandraoperator": {
			Type:        schema.TypeBool,
			Description: "Rule id: 69 - Do not monitor processes if Kubernetes container name equals 'cassandra-operator'",
			Optional:    true,
			Default:     true,
		},
		"exe_uipath": {
			Type:        schema.TypeBool,
			Description: "Rule id: 70 - Do not monitor processes if EXE name contains 'UiPath'",
			Optional:    true,
			Default:     true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"host_group_id":                         me.HostGroupID,
		"php_script":                            me.RuleID1,
		"exe_phpcgi":                            me.RuleID2,
		"aspnetcore":                            me.RuleID3,
		"exe_w3wp":                              me.RuleID4,
		"exe_caddy":                             me.RuleID5,
		"exe_influxd":                           me.RuleID6,
		"exe_adapter":                           me.RuleID7,
		"exe_auctioneer":                        me.RuleID8,
		"exe_bbs":                               me.RuleID9,
		"exe_ccuploader":                        me.RuleID10,
		"exe_doppler":                           me.RuleID11,
		"exe_gorouter":                          me.RuleID12,
		"exe_locket":                            me.RuleID13,
		"exe_metron":                            me.RuleID14,
		"exe_rep":                               me.RuleID16,
		"exe_routeemitter":                      me.RuleID17,
		"exe_routeregistrar":                    me.RuleID18,
		"exe_routingapi":                        me.RuleID19,
		"exe_schedular":                         me.RuleID20,
		"exe_silkdaemon":                        me.RuleID21,
		"exe_switchboard":                       me.RuleID22,
		"exe_syslogdrainbinder":                 me.RuleID23,
		"exe_tpswatcher":                        me.RuleID24,
		"exe_trafficcontroller":                 me.RuleID25,
		"node_prebuildinstall":                  me.RuleID26,
		"node_npm":                              me.RuleID27,
		"node_grunt":                            me.RuleID28,
		"node_typescript":                       me.RuleID29,
		"node_nodepregyp":                       me.RuleID32,
		"node_nodegyp":                          me.RuleID33,
		"node_gulpcli":                          me.RuleID34,
		"node_binpm2":                           me.RuleID35,
		"cf_appsmanagerjs":                      me.RuleID36,
		"cf":                                    me.RuleID37,
		"k8s_containerpod":                      me.RuleID38,
		"docker_pauseamd64":                     me.RuleID39,
		"k8s_namespace":                         me.RuleID40,
		"container":                             me.RuleID41,
		"exe_tmpbuildpacks":                     me.RuleID43,
		"exe_oc":                                me.RuleID44,
		"node_yarn":                             me.RuleID45,
		"exe_optcnibinhostlocal":                me.RuleID46,
		"go_static":                             me.RuleID47,
		"exe_mqsi":                              me.RuleID48,
		"exe_filebeat":                          me.RuleID49,
		"exe_metricbeat":                        me.RuleID50,
		"exe_packetbeat":                        me.RuleID51,
		"exe_auditbeat":                         me.RuleID52,
		"exe_heartbeat":                         me.RuleID53,
		"exe_functionbeat":                      me.RuleID54,
		"exe_grootfs":                           me.RuleID55,
		"exe_tardis":                            me.RuleID56,
		"jar_eclipseequinox":                    me.RuleID57,
		"exe_caliconode":                        me.RuleID58,
		"exe_casclient":                         me.RuleID59,
		"jar_dtibmmqconnector":                  me.RuleID60,
		"exe_agentworker":                       me.RuleID61,
		"aspnetcore_agentworker":                me.RuleID62,
		"exe_agentlistener":                     me.RuleID63,
		"aspnetcore_agentlistener":              me.RuleID64,
		"exe_flexnetjobexecutorservice":         me.RuleID65,
		"exe_flexnetmaintenanceremotingservice": me.RuleID66,
		"exe_userbinpiper":                      me.RuleID67,
		"node_corepack":                         me.RuleID68,
		"k8s_cassandraoperator":                 me.RuleID69,
		"exe_uipath":                            me.RuleID70,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"host_group_id":                         &me.HostGroupID,
		"php_script":                            &me.RuleID1,
		"exe_phpcgi":                            &me.RuleID2,
		"aspnetcore":                            &me.RuleID3,
		"exe_w3wp":                              &me.RuleID4,
		"exe_caddy":                             &me.RuleID5,
		"exe_influxd":                           &me.RuleID6,
		"exe_adapter":                           &me.RuleID7,
		"exe_auctioneer":                        &me.RuleID8,
		"exe_bbs":                               &me.RuleID9,
		"exe_ccuploader":                        &me.RuleID10,
		"exe_doppler":                           &me.RuleID11,
		"exe_gorouter":                          &me.RuleID12,
		"exe_locket":                            &me.RuleID13,
		"exe_metron":                            &me.RuleID14,
		"exe_rep":                               &me.RuleID16,
		"exe_routeemitter":                      &me.RuleID17,
		"exe_routeregistrar":                    &me.RuleID18,
		"exe_routingapi":                        &me.RuleID19,
		"exe_schedular":                         &me.RuleID20,
		"exe_silkdaemon":                        &me.RuleID21,
		"exe_switchboard":                       &me.RuleID22,
		"exe_syslogdrainbinder":                 &me.RuleID23,
		"exe_tpswatcher":                        &me.RuleID24,
		"exe_trafficcontroller":                 &me.RuleID25,
		"node_prebuildinstall":                  &me.RuleID26,
		"node_npm":                              &me.RuleID27,
		"node_grunt":                            &me.RuleID28,
		"node_typescript":                       &me.RuleID29,
		"node_nodepregyp":                       &me.RuleID32,
		"node_nodegyp":                          &me.RuleID33,
		"node_gulpcli":                          &me.RuleID34,
		"node_binpm2":                           &me.RuleID35,
		"cf_appsmanagerjs":                      &me.RuleID36,
		"cf":                                    &me.RuleID37,
		"k8s_containerpod":                      &me.RuleID38,
		"docker_pauseamd64":                     &me.RuleID39,
		"k8s_namespace":                         &me.RuleID40,
		"container":                             &me.RuleID41,
		"exe_tmpbuildpacks":                     &me.RuleID43,
		"exe_oc":                                &me.RuleID44,
		"node_yarn":                             &me.RuleID45,
		"exe_optcnibinhostlocal":                &me.RuleID46,
		"go_static":                             &me.RuleID47,
		"exe_mqsi":                              &me.RuleID48,
		"exe_filebeat":                          &me.RuleID49,
		"exe_metricbeat":                        &me.RuleID50,
		"exe_packetbeat":                        &me.RuleID51,
		"exe_auditbeat":                         &me.RuleID52,
		"exe_heartbeat":                         &me.RuleID53,
		"exe_functionbeat":                      &me.RuleID54,
		"exe_grootfs":                           &me.RuleID55,
		"exe_tardis":                            &me.RuleID56,
		"jar_eclipseequinox":                    &me.RuleID57,
		"exe_caliconode":                        &me.RuleID58,
		"exe_casclient":                         &me.RuleID59,
		"jar_dtibmmqconnector":                  &me.RuleID60,
		"exe_agentworker":                       &me.RuleID61,
		"aspnetcore_agentworker":                &me.RuleID62,
		"exe_agentlistener":                     &me.RuleID63,
		"aspnetcore_agentlistener":              &me.RuleID64,
		"exe_flexnetjobexecutorservice":         &me.RuleID65,
		"exe_flexnetmaintenanceremotingservice": &me.RuleID66,
		"exe_userbinpiper":                      &me.RuleID67,
		"node_corepack":                         &me.RuleID68,
		"k8s_cassandraoperator":                 &me.RuleID69,
		"exe_uipath":                            &me.RuleID70,
	})
}
