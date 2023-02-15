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

package detectionflags

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	AddNodeJsScriptName                    bool    `json:"addNodeJsScriptName"`                    // In older versions, Node.js applications were distinguished based on their directory name, omitting the script name. Changing this setting may change the general handling of Node.js process groups. Leave unchanged if in doubt.
	AutoDetectCassandraClusters            bool    `json:"autoDetectCassandraClusters"`            // Enabling this flag will detect separate Cassandra process groups based on the configured Cassandra cluster name.
	AutoDetectSpringBoot                   bool    `json:"autoDetectSpringBoot"`                   // Enabling this flag will detect Spring Boot process groups based on command line and applications' configuration files.
	AutoDetectTibcoContainerEditionEngines bool    `json:"autoDetectTibcoContainerEditionEngines"` // Enabling this flag will detect separate TIBCO BusinessWorks process groups per engine property file.
	AutoDetectTibcoEngines                 bool    `json:"autoDetectTibcoEngines"`                 // Enabling this flag will detect separate TIBCO BusinessWorks process groups per engine property file.
	AutoDetectWebMethodsIntegrationServer  bool    `json:"autoDetectWebMethodsIntegrationServer"`  // Enabling this flag will detect webMethods Integration Server including specific properties like install root and product name.
	AutoDetectWebSphereLibertyApplication  bool    `json:"autoDetectWebSphereLibertyApplication"`  // Enabling this flag will detect separate WebSphere Liberty process groups based on java command line.
	GroupIBMMQbyInstanceName               bool    `json:"groupIBMMQbyInstanceName"`               // Enable to group and separately analyze the processes of each IBM MQ Queue manager instance. Each process group receives a unique name based on the queue manager instance name.
	IdentifyJbossServerBySystemProperty    bool    `json:"identifyJbossServerBySystemProperty"`    // Enabling this flag will detect the JBoss server name from the system property jboss.server.name=<server-name>, only if -D[Server:<server-name>] is not set.
	IgnoreUniqueIdentifiers                bool    `json:"ignoreUniqueIdentifiers"`                // To determine the unique identity of each detected process, and to generate a unique name for each detected process, Dynatrace evaluates the name of the directory that each process binary is contained within. For application containers like Tomcat and JBoss, Dynatrace evaluates important directories like CATALINA_HOME and JBOSS_HOME for this information. In some automated deployment scenarios such directory names are updated automatically with new version numbers, build numbers, dates, or GUIDs. Enable this setting to ensure that automated directory name changes don't result in Dynatrace registering pre-existing processes as new processes.
	Scope                                  *string `json:"-" scope:"scope"`                        // The scope of this setting (HOST, HOST_GROUP). Omit this property if you want to cover the whole environment.
	ShortLivedProcessesMonitoring          bool    `json:"shortLivedProcessesMonitoring"`          // Enable to monitor CPU and memory usage of short lived processes, otherwise being lost by traditional monitoring. Disabling this flag blocks passing data to cluster only, it does not stop data collection and has no effect on performance.
	SplitOracleDatabasePG                  bool    `json:"splitOracleDatabasePG"`                  // Enable to group and separately analyze the processes of each Oracle DB. Each process group receives a unique name based on the Oracle DB SID.
	SplitOracleListenerPG                  bool    `json:"splitOracleListenerPG"`                  // Enable to group and separately analyze the processes of each Oracle Listener. Each process group receives a unique name based on the Oracle Listener name.
	UseCatalinaBase                        bool    `json:"useCatalinaBase"`                        // By default, Tomcat clusters are identified and named based on the CATALINA_HOME directory name. This setting results in the use of the CATALINA_BASE directory name to identify multiple Tomcat nodes within each Tomcat cluster. If this setting is not enabled, each CATALINA_HOME+CATALINA_BASE combination will be considered a separate Tomcat cluster. In other words, Tomcat clusters can't have multiple nodes on a single host.
	UseDockerContainerName                 bool    `json:"useDockerContainerName"`                 // By default, Dynatrace uses image names as identifiers for individual process groups, with one process-group instance per host. Normally Docker container names can't serve as stable identifiers of process group instances because they are variable and auto-generated. You can however manually assign proper container names to their Docker instances. Such manually-assigned container names can serve as reliable process-group instance identifiers. This flag instructs Dynatrace to use Docker-provided names to distinguish between multiple instances of the same image. If this flag is not applied and you run multiple containers of the same image on the same host, the resulting processes will be consolidated into a single process view. Use this flag with caution!
}

func (me *Settings) Name() string {
	return *me.Scope
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"add_node_js_script_name": {
			Type:        schema.TypeBool,
			Description: "In older versions, Node.js applications were distinguished based on their directory name, omitting the script name. Changing this setting may change the general handling of Node.js process groups. Leave unchanged if in doubt.",
			Required:    true,
		},
		"auto_detect_cassandra_clusters": {
			Type:        schema.TypeBool,
			Description: "Enabling this flag will detect separate Cassandra process groups based on the configured Cassandra cluster name.",
			Required:    true,
		},
		"auto_detect_spring_boot": {
			Type:        schema.TypeBool,
			Description: "Enabling this flag will detect Spring Boot process groups based on command line and applications' configuration files.",
			Required:    true,
		},
		"auto_detect_tibco_container_edition_engines": {
			Type:        schema.TypeBool,
			Description: "Enabling this flag will detect separate TIBCO BusinessWorks process groups per engine property file.",
			Required:    true,
		},
		"auto_detect_tibco_engines": {
			Type:        schema.TypeBool,
			Description: "Enabling this flag will detect separate TIBCO BusinessWorks process groups per engine property file.",
			Required:    true,
		},
		"auto_detect_web_methods_integration_server": {
			Type:        schema.TypeBool,
			Description: "Enabling this flag will detect webMethods Integration Server including specific properties like install root and product name.",
			Required:    true,
		},
		"auto_detect_web_sphere_liberty_application": {
			Type:        schema.TypeBool,
			Description: "Enabling this flag will detect separate WebSphere Liberty process groups based on java command line.",
			Required:    true,
		},
		"group_ibmmqby_instance_name": {
			Type:        schema.TypeBool,
			Description: "Enable to group and separately analyze the processes of each IBM MQ Queue manager instance. Each process group receives a unique name based on the queue manager instance name.",
			Required:    true,
		},
		"identify_jboss_server_by_system_property": {
			Type:        schema.TypeBool,
			Description: "Enabling this flag will detect the JBoss server name from the system property jboss.server.name=<server-name>, only if -D[Server:<server-name>] is not set.",
			Required:    true,
		},
		"ignore_unique_identifiers": {
			Type:        schema.TypeBool,
			Description: "To determine the unique identity of each detected process, and to generate a unique name for each detected process, Dynatrace evaluates the name of the directory that each process binary is contained within. For application containers like Tomcat and JBoss, Dynatrace evaluates important directories like CATALINA_HOME and JBOSS_HOME for this information. In some automated deployment scenarios such directory names are updated automatically with new version numbers, build numbers, dates, or GUIDs. Enable this setting to ensure that automated directory name changes don't result in Dynatrace registering pre-existing processes as new processes.",
			Required:    true,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (HOST, HOST_GROUP). Omit this property if you want to cover the whole environment.",
			Optional:    true,
			Default:     "environment",
		},
		"short_lived_processes_monitoring": {
			Type:        schema.TypeBool,
			Description: "Enable to monitor CPU and memory usage of short lived processes, otherwise being lost by traditional monitoring. Disabling this flag blocks passing data to cluster only, it does not stop data collection and has no effect on performance.",
			Required:    true,
		},
		"split_oracle_database_pg": {
			Type:        schema.TypeBool,
			Description: "Enable to group and separately analyze the processes of each Oracle DB. Each process group receives a unique name based on the Oracle DB SID.",
			Required:    true,
		},
		"split_oracle_listener_pg": {
			Type:        schema.TypeBool,
			Description: "Enable to group and separately analyze the processes of each Oracle Listener. Each process group receives a unique name based on the Oracle Listener name.",
			Required:    true,
		},
		"use_catalina_base": {
			Type:        schema.TypeBool,
			Description: "By default, Tomcat clusters are identified and named based on the CATALINA_HOME directory name. This setting results in the use of the CATALINA_BASE directory name to identify multiple Tomcat nodes within each Tomcat cluster. If this setting is not enabled, each CATALINA_HOME+CATALINA_BASE combination will be considered a separate Tomcat cluster. In other words, Tomcat clusters can't have multiple nodes on a single host.",
			Required:    true,
		},
		"use_docker_container_name": {
			Type:        schema.TypeBool,
			Description: "By default, Dynatrace uses image names as identifiers for individual process groups, with one process-group instance per host. Normally Docker container names can't serve as stable identifiers of process group instances because they are variable and auto-generated. You can however manually assign proper container names to their Docker instances. Such manually-assigned container names can serve as reliable process-group instance identifiers. This flag instructs Dynatrace to use Docker-provided names to distinguish between multiple instances of the same image. If this flag is not applied and you run multiple containers of the same image on the same host, the resulting processes will be consolidated into a single process view. Use this flag with caution!",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"add_node_js_script_name":                     me.AddNodeJsScriptName,
		"auto_detect_cassandra_clusters":              me.AutoDetectCassandraClusters,
		"auto_detect_spring_boot":                     me.AutoDetectSpringBoot,
		"auto_detect_tibco_container_edition_engines": me.AutoDetectTibcoContainerEditionEngines,
		"auto_detect_tibco_engines":                   me.AutoDetectTibcoEngines,
		"auto_detect_web_methods_integration_server":  me.AutoDetectWebMethodsIntegrationServer,
		"auto_detect_web_sphere_liberty_application":  me.AutoDetectWebSphereLibertyApplication,
		"group_ibmmqby_instance_name":                 me.GroupIBMMQbyInstanceName,
		"identify_jboss_server_by_system_property":    me.IdentifyJbossServerBySystemProperty,
		"ignore_unique_identifiers":                   me.IgnoreUniqueIdentifiers,
		"scope":                                       me.Scope,
		"short_lived_processes_monitoring":            me.ShortLivedProcessesMonitoring,
		"split_oracle_database_pg":                    me.SplitOracleDatabasePG,
		"split_oracle_listener_pg":                    me.SplitOracleListenerPG,
		"use_catalina_base":                           me.UseCatalinaBase,
		"use_docker_container_name":                   me.UseDockerContainerName,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"add_node_js_script_name":                     &me.AddNodeJsScriptName,
		"auto_detect_cassandra_clusters":              &me.AutoDetectCassandraClusters,
		"auto_detect_spring_boot":                     &me.AutoDetectSpringBoot,
		"auto_detect_tibco_container_edition_engines": &me.AutoDetectTibcoContainerEditionEngines,
		"auto_detect_tibco_engines":                   &me.AutoDetectTibcoEngines,
		"auto_detect_web_methods_integration_server":  &me.AutoDetectWebMethodsIntegrationServer,
		"auto_detect_web_sphere_liberty_application":  &me.AutoDetectWebSphereLibertyApplication,
		"group_ibmmqby_instance_name":                 &me.GroupIBMMQbyInstanceName,
		"identify_jboss_server_by_system_property":    &me.IdentifyJbossServerBySystemProperty,
		"ignore_unique_identifiers":                   &me.IgnoreUniqueIdentifiers,
		"scope":                                       &me.Scope,
		"short_lived_processes_monitoring":            &me.ShortLivedProcessesMonitoring,
		"split_oracle_database_pg":                    &me.SplitOracleDatabasePG,
		"split_oracle_listener_pg":                    &me.SplitOracleListenerPG,
		"use_catalina_base":                           &me.UseCatalinaBase,
		"use_docker_container_name":                   &me.UseDockerContainerName,
	})
}
