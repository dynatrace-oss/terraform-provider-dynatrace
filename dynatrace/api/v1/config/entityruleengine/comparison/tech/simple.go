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

package tech

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Simple The value to compare to.
type Simple struct {
	Type         *SimpleTechType            `json:"type,omitempty"`         // Predefined technology, if technology is not predefined, then the verbatim type must be set
	VerbatimType *string                    `json:"verbatimType,omitempty"` // Non-predefined technology, use for custom technologies.
	Unknowns     map[string]json.RawMessage `json:"-"`
}

func (st *Simple) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "Predefined technology, if technology is not predefined, then the verbatim type must be set.",
			Optional:    true,
		},
		"verbatim_type": {
			Type:        schema.TypeString,
			Description: "Non-predefined technology, use for custom technologies",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider",
			Optional:    true,
		},
	}
}

func (st *Simple) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(st.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("type", st.Type.String()); err != nil {
		return err
	}
	if err := properties.Encode("verbatim_type", st.VerbatimType); err != nil {
		return err
	}
	return nil
}

func (st *Simple) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), st); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &st.Unknowns); err != nil {
			return err
		}
		delete(st.Unknowns, "type")
		delete(st.Unknowns, "verbatim_type")
		if len(st.Unknowns) == 0 {
			st.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("type"); ok {
		st.Type = SimpleTechType(value.(string)).Ref()
	}
	if value, ok := decoder.GetOk("verbatim_type"); ok {
		st.VerbatimType = opt.NewString(value.(string))
	}
	return nil
}

func (st *Simple) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(st.Unknowns) > 0 {
		for k, v := range st.Unknowns {
			m[k] = v
		}
	}
	if st.Type != nil {
		rawMessage, err := json.Marshal(st.Type)
		if err != nil {
			return nil, err
		}
		m["type"] = rawMessage
	}
	if st.VerbatimType != nil {
		rawMessage, err := json.Marshal(st.VerbatimType)
		if err != nil {
			return nil, err
		}
		m["verbatimType"] = rawMessage
	}
	return json.Marshal(m)
}

func (st *Simple) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["type"]; found {
		if err := json.Unmarshal(v, &st.Type); err != nil {
			return err
		}
	}
	if v, found := m["verbatimType"]; found {
		if err := json.Unmarshal(v, &st.VerbatimType); err != nil {
			return err
		}
	}
	delete(m, "verbatimType")
	delete(m, "type")
	if len(m) > 0 {
		st.Unknowns = m
	}
	return nil
}

// SimpleTechType Predefined technology, if technology is not predefined, then the verbatim type must be set
type SimpleTechType string

func (v SimpleTechType) Ref() *SimpleTechType {
	return &v
}

func (v *SimpleTechType) String() string {
	return string(*v)
}

// SimpleTechTypes offers the known enum values
var SimpleTechTypes = struct {
	ActiveMq                             SimpleTechType
	ActiveMqArtemis                      SimpleTechType
	AdoNet                               SimpleTechType
	AIX                                  SimpleTechType
	Akka                                 SimpleTechType
	AmazonRedshift                       SimpleTechType
	Amqp                                 SimpleTechType
	ApacheCamel                          SimpleTechType
	ApacheCassandra                      SimpleTechType
	ApacheCouchDB                        SimpleTechType
	ApacheDerby                          SimpleTechType
	ApacheHTTPClientAsync                SimpleTechType
	ApacheHTTPClientSync                 SimpleTechType
	ApacheHTTPServer                     SimpleTechType
	ApacheKafka                          SimpleTechType
	ApacheSolr                           SimpleTechType
	ApacheStorm                          SimpleTechType
	ApacheSynapse                        SimpleTechType
	ApacheTomcat                         SimpleTechType
	Apparmor                             SimpleTechType
	ApplicationInsightsSdk               SimpleTechType
	ASPDotNet                            SimpleTechType
	ASPDotNetCore                        SimpleTechType
	ASPDotNetCoreSignalr                 SimpleTechType
	ASPDotNetSignalr                     SimpleTechType
	AsyncHTTPClient                      SimpleTechType
	AWSLambda                            SimpleTechType
	AWSRds                               SimpleTechType
	AWSService                           SimpleTechType
	Axis                                 SimpleTechType
	AzureFunctions                       SimpleTechType
	AzureServiceBus                      SimpleTechType
	AzureServiceFabric                   SimpleTechType
	AzureStorage                         SimpleTechType
	Boshbpm                              SimpleTechType
	Citrix                               SimpleTechType
	CitrixCommon                         SimpleTechType
	CitrixDesktopDeliveryControllers     SimpleTechType
	CitrixDirector                       SimpleTechType
	CitrixLicenseServer                  SimpleTechType
	CitrixProvisioningServices           SimpleTechType
	CitrixStorefront                     SimpleTechType
	CitrixVirtualDeliveryAgent           SimpleTechType
	CitrixWorkspaceEnvironmentManagement SimpleTechType
	CloudFoundry                         SimpleTechType
	CloudFoundryAuctioneer               SimpleTechType
	CloudFoundryBosh                     SimpleTechType
	CloudFoundryGorouter                 SimpleTechType
	Coldfusion                           SimpleTechType
	Containerd                           SimpleTechType
	CoreDNS                              SimpleTechType
	Couchbase                            SimpleTechType
	Crio                                 SimpleTechType
	Cxf                                  SimpleTechType
	Datastax                             SimpleTechType
	DB2                                  SimpleTechType
	DiegoCell                            SimpleTechType
	Docker                               SimpleTechType
	DotNet                               SimpleTechType
	DotNetRemoting                       SimpleTechType
	ElasticSearch                        SimpleTechType
	Envoy                                SimpleTechType
	Erlang                               SimpleTechType
	Etcd                                 SimpleTechType
	F5Ltm                                SimpleTechType
	Fsharp                               SimpleTechType
	Garden                               SimpleTechType
	Glassfish                            SimpleTechType
	Go                                   SimpleTechType
	GraalTruffle                         SimpleTechType
	Grpc                                 SimpleTechType
	Grsecurity                           SimpleTechType
	Hadoop                               SimpleTechType
	HadoopHdfs                           SimpleTechType
	HadoopYarn                           SimpleTechType
	Haproxy                              SimpleTechType
	Heat                                 SimpleTechType
	Hessian                              SimpleTechType
	HornetQ                              SimpleTechType
	IBMCICSRegion                        SimpleTechType
	IBMCICSTransactionGateway            SimpleTechType
	IBMIMSConnectRegion                  SimpleTechType
	IBMIMSControlRegion                  SimpleTechType
	IBMIMSMessageProcessingRegion        SimpleTechType
	IBMIMSSoapGateway                    SimpleTechType
	IBMIntegrationBus                    SimpleTechType
	IBMMq                                SimpleTechType
	IBMMqClient                          SimpleTechType
	IBMWebshprereApplicationServer       SimpleTechType
	IBMWebshprereLiberty                 SimpleTechType
	IIS                                  SimpleTechType
	IISAppPool                           SimpleTechType
	Istio                                SimpleTechType
	Java                                 SimpleTechType
	JaxWs                                SimpleTechType
	JBoss                                SimpleTechType
	JBossEap                             SimpleTechType
	JdkHTTPServer                        SimpleTechType
	Jersey                               SimpleTechType
	Jetty                                SimpleTechType
	Jruby                                SimpleTechType
	Jython                               SimpleTechType
	Kubernetes                           SimpleTechType
	Libvirt                              SimpleTechType
	Linkerd                              SimpleTechType
	Mariadb                              SimpleTechType
	Memcache                             SimpleTechType
	MicrosoftSQLServer                   SimpleTechType
	Mongodb                              SimpleTechType
	MSSQLClient                          SimpleTechType
	MuleEsb                              SimpleTechType
	MySQL                                SimpleTechType
	MySQLConnector                       SimpleTechType
	NetflixServo                         SimpleTechType
	Netty                                SimpleTechType
	Nginx                                SimpleTechType
	NodeJs                               SimpleTechType
	OkHTTPClient                         SimpleTechType
	OneAgentSdk                          SimpleTechType
	Opencensus                           SimpleTechType
	Openshift                            SimpleTechType
	OpenStackCompute                     SimpleTechType
	OpenStackController                  SimpleTechType
	Opentelemetry                        SimpleTechType
	Opentracing                          SimpleTechType
	OpenLiberty                          SimpleTechType
	OracleDatabase                       SimpleTechType
	OracleWeblogic                       SimpleTechType
	Owin                                 SimpleTechType
	Perl                                 SimpleTechType
	PHP                                  SimpleTechType
	PHPFpm                               SimpleTechType
	Play                                 SimpleTechType
	PostgreSQL                           SimpleTechType
	PostgreSQLDotNetDataProvider         SimpleTechType
	PowerDNS                             SimpleTechType
	Progress                             SimpleTechType
	Python                               SimpleTechType
	RabbitMq                             SimpleTechType
	Redis                                SimpleTechType
	Resteasy                             SimpleTechType
	Restlet                              SimpleTechType
	Riak                                 SimpleTechType
	Ruby                                 SimpleTechType
	SagWebmethodsIs                      SimpleTechType
	SAP                                  SimpleTechType
	SAPHanadb                            SimpleTechType
	SAPHybris                            SimpleTechType
	SAPMaxdb                             SimpleTechType
	SAPSybase                            SimpleTechType
	Scala                                SimpleTechType
	Selinux                              SimpleTechType
	Sharepoint                           SimpleTechType
	Spark                                SimpleTechType
	Spring                               SimpleTechType
	Sqlite                               SimpleTechType
	Thrift                               SimpleTechType
	Tibco                                SimpleTechType
	TibcoBusinessWorks                   SimpleTechType
	TibcoEms                             SimpleTechType
	VarnishCache                         SimpleTechType
	Vim2                                 SimpleTechType
	VirtualMachineKvm                    SimpleTechType
	VirtualMachineQemu                   SimpleTechType
	Wildfly                              SimpleTechType
	WindowsContainers                    SimpleTechType
	Wink                                 SimpleTechType
	ZeroMq                               SimpleTechType
}{
	"ACTIVE_MQ",
	"ACTIVE_MQ_ARTEMIS",
	"ADO_NET",
	"AIX",
	"AKKA",
	"AMAZON_REDSHIFT",
	"AMQP",
	"APACHE_CAMEL",
	"APACHE_CASSANDRA",
	"APACHE_COUCH_DB",
	"APACHE_DERBY",
	"APACHE_HTTP_CLIENT_ASYNC",
	"APACHE_HTTP_CLIENT_SYNC",
	"APACHE_HTTP_SERVER",
	"APACHE_KAFKA",
	"APACHE_SOLR",
	"APACHE_STORM",
	"APACHE_SYNAPSE",
	"APACHE_TOMCAT",
	"APPARMOR",
	"APPLICATION_INSIGHTS_SDK",
	"ASP_DOTNET",
	"ASP_DOTNET_CORE",
	"ASP_DOTNET_CORE_SIGNALR",
	"ASP_DOTNET_SIGNALR",
	"ASYNC_HTTP_CLIENT",
	"AWS_LAMBDA",
	"AWS_RDS",
	"AWS_SERVICE",
	"AXIS",
	"AZURE_FUNCTIONS",
	"AZURE_SERVICE_BUS",
	"AZURE_SERVICE_FABRIC",
	"AZURE_STORAGE",
	"BOSHBPM",
	"CITRIX",
	"CITRIX_COMMON",
	"CITRIX_DESKTOP_DELIVERY_CONTROLLERS",
	"CITRIX_DIRECTOR",
	"CITRIX_LICENSE_SERVER",
	"CITRIX_PROVISIONING_SERVICES",
	"CITRIX_STOREFRONT",
	"CITRIX_VIRTUAL_DELIVERY_AGENT",
	"CITRIX_WORKSPACE_ENVIRONMENT_MANAGEMENT",
	"CLOUDFOUNDRY",
	"CLOUDFOUNDRY_AUCTIONEER",
	"CLOUDFOUNDRY_BOSH",
	"CLOUDFOUNDRY_GOROUTER",
	"COLDFUSION",
	"CONTAINERD",
	"CORE_DNS",
	"COUCHBASE",
	"CRIO",
	"CXF",
	"DATASTAX",
	"DB2",
	"DIEGO_CELL",
	"DOCKER",
	"DOTNET",
	"DOTNET_REMOTING",
	"ELASTIC_SEARCH",
	"ENVOY",
	"ERLANG",
	"ETCD",
	"F5_LTM",
	"FSHARP",
	"GARDEN",
	"GLASSFISH",
	"GO",
	"GRAAL_TRUFFLE",
	"GRPC",
	"GRSECURITY",
	"HADOOP",
	"HADOOP_HDFS",
	"HADOOP_YARN",
	"HAPROXY",
	"HEAT",
	"HESSIAN",
	"HORNET_Q",
	"IBM_CICS_REGION",
	"IBM_CICS_TRANSACTION_GATEWAY",
	"IBM_IMS_CONNECT_REGION",
	"IBM_IMS_CONTROL_REGION",
	"IBM_IMS_MESSAGE_PROCESSING_REGION",
	"IBM_IMS_SOAP_GATEWAY",
	"IBM_INTEGRATION_BUS",
	"IBM_MQ",
	"IBM_MQ_CLIENT",
	"IBM_WEBSHPRERE_APPLICATION_SERVER",
	"IBM_WEBSHPRERE_LIBERTY",
	"IIS",
	"IIS_APP_POOL",
	"ISTIO",
	"JAVA",
	"JAX_WS",
	"JBOSS",
	"JBOSS_EAP",
	"JDK_HTTP_SERVER",
	"JERSEY",
	"JETTY",
	"JRUBY",
	"JYTHON",
	"KUBERNETES",
	"LIBVIRT",
	"LINKERD",
	"MARIADB",
	"MEMcache",
	"MICROSOFT_SQL_SERVER",
	"MONGODB",
	"MSSQL_CLIENT",
	"MULE_ESB",
	"MYSQL",
	"MYSQL_CONNECTOR",
	"NETFLIX_SERVO",
	"NETTY",
	"NGINX",
	"NODE_JS",
	"OK_HTTP_CLIENT",
	"ONEAGENT_SDK",
	"OPENCENSUS",
	"OPENSHIFT",
	"OPENSTACK_COMPUTE",
	"OPENSTACK_CONTROLLER",
	"OPENTELEMETRY",
	"OPENTRACING",
	"OPEN_LIBERTY",
	"ORACLE_DATABASE",
	"ORACLE_WEBLOGIC",
	"OWIN",
	"PERL",
	"PHP",
	"PHP_FPM",
	"PLAY",
	"POSTGRE_SQL",
	"POSTGRE_SQL_DOTNET_DATA_PROVIDER",
	"POWER_DNS",
	"PROGRESS",
	"PYTHON",
	"RABBIT_MQ",
	"REDIS",
	"RESTEASY",
	"RESTLET",
	"RIAK",
	"RUBY",
	"SAG_WEBMETHODS_IS",
	"SAP",
	"SAP_HANADB",
	"SAP_HYBRIS",
	"SAP_MAXDB",
	"SAP_SYBASE",
	"SCALA",
	"SELINUX",
	"SHAREPOINT",
	"SPARK",
	"SPRING",
	"SQLITE",
	"THRIFT",
	"TIBCO",
	"TIBCO_BUSINESS_WORKS",
	"TIBCO_EMS",
	"VARNISH_CACHE",
	"VIM2",
	"VIRTUAL_MACHINE_KVM",
	"VIRTUAL_MACHINE_QEMU",
	"WILDFLY",
	"WINDOWS_CONTAINERS",
	"WINK",
	"ZERO_MQ",
}
