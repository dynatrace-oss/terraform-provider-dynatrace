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

package requestattributes

// Capture What to capture from the method.
type Capture string

// Captures offers the known enum values
var Captures = struct {
	Argument        Capture
	ClassName       Capture
	MethodName      Capture
	Occurrences     Capture
	SimpleClassName Capture
	This            Capture
}{
	"ARGUMENT",
	"CLASS_NAME",
	"METHOD_NAME",
	"OCCURRENCES",
	"SIMPLE_CLASS_NAME",
	"THIS",
}

// CapturingAndStorageLocation Specifies the location where the values are captured and stored.
//  Required if the **source** is one of the following: `GET_PARAMETER`, `URI`, `REQUEST_HEADER`, `RESPONSE_HEADER`.
//  Not applicable in other cases.
//  If the **source** value is `REQUEST_HEADER` or `RESPONSE_HEADER`, the `CAPTURE_AND_STORE_ON_BOTH` location is not allowed.
type CapturingAndStorageLocation string

func (me CapturingAndStorageLocation) Ref() *CapturingAndStorageLocation {
	return &me
}

// CapturingAndStorageLocations offers the known enum values
var CapturingAndStorageLocations = struct {
	CaptureAndStoreOnBoth        CapturingAndStorageLocation
	CaptureAndStoreOnClient      CapturingAndStorageLocation
	CaptureAndStoreOnServer      CapturingAndStorageLocation
	CaptureOnClientStoreOnServer CapturingAndStorageLocation
}{
	"CAPTURE_AND_STORE_ON_BOTH",
	"CAPTURE_AND_STORE_ON_CLIENT",
	"CAPTURE_AND_STORE_ON_SERVER",
	"CAPTURE_ON_CLIENT_STORE_ON_SERVER",
}

// SessionAttributeTechnology The technology of the session attribute to capture if the **source** value is `SESSION_ATTRIBUTE`. \n\n Not applicable in other cases.
type SessionAttributeTechnology string

func (me SessionAttributeTechnology) Ref() *SessionAttributeTechnology {
	return &me
}

// SessionAttributeTechnologys offers the known enum values
var SessionAttributeTechnologys = struct {
	ASPNet     SessionAttributeTechnology
	ASPNetCore SessionAttributeTechnology
	Java       SessionAttributeTechnology
}{
	"ASP_NET",
	"ASP_NET_CORE",
	"JAVA",
}

// Technology The technology of the method to capture if the **source** value is `METHOD_PARAM`. \n\n Not applicable in other cases.
type Technology string

func (me Technology) Ref() *Technology {
	return &me
}

// Technologys offers the known enum values
var Technologys = struct {
	DotNet Technology
	Java   Technology
	PHP    Technology
}{
	"DOTNET",
	"JAVA",
	"PHP",
}

// Source The source of the attribute to capture. Works in conjunction with **parameterName** or **methods** and **technology**.
type Source string

// Sources offers the known enum values
var Sources = struct {
	CICSSdk          Source
	ClientIP         Source
	CustomAttribute  Source
	IibLabel         Source
	IibNode          Source
	MethodParam      Source
	PostParameter    Source
	QueryParameter   Source
	RequestHeader    Source
	ResponseHeader   Source
	SessionAttribute Source
	URI              Source
	URIPath          Source
}{
	"CICS_SDK",
	"CLIENT_IP",
	"CUSTOM_ATTRIBUTE",
	"IIB_LABEL",
	"IIB_NODE",
	"METHOD_PARAM",
	"POST_PARAMETER",
	"QUERY_PARAMETER",
	"REQUEST_HEADER",
	"RESPONSE_HEADER",
	"SESSION_ATTRIBUTE",
	"URI",
	"URI_PATH",
}

// IIBNodeType The IBM integration bus node type for which the value is captured.
//  This or `iibMethodNodeCondition` is required if the **source** is: `IIB_NODE`.
//  Not applicable in other cases.
type IIBNodeType string

func (me IIBNodeType) Ref() *IIBNodeType {
	return &me
}

// IIBNodeTypes offers the known enum values
var IIBNodeTypes = struct {
	AggregateControlNode       IIBNodeType
	AggregateReplyNode         IIBNodeType
	AggregateRequestNode       IIBNodeType
	CallableFlowReplyNode      IIBNodeType
	CollectorNode              IIBNodeType
	ComputeNode                IIBNodeType
	DatabaseNode               IIBNodeType
	DecisionServiceNode        IIBNodeType
	DotNetComputeNode          IIBNodeType
	FileReadNode               IIBNodeType
	FilterNode                 IIBNodeType
	FlowOrderNode              IIBNodeType
	GroupCompleteNode          IIBNodeType
	GroupGatherNode            IIBNodeType
	GroupScatterNode           IIBNodeType
	HTTPHeader                 IIBNodeType
	JavaComputeNode            IIBNodeType
	JmsClientReceive           IIBNodeType
	JmsClientReplyNode         IIBNodeType
	JmsHeader                  IIBNodeType
	MqGetNode                  IIBNodeType
	MqOutputNode               IIBNodeType
	PassthruNode               IIBNodeType
	ResetContentDescriptorNode IIBNodeType
	ReSequenceNode             IIBNodeType
	RouteNode                  IIBNodeType
	SAPReplyNode               IIBNodeType
	ScaReplyNode               IIBNodeType
	SecurityPep                IIBNodeType
	SequenceNode               IIBNodeType
	SoapExtractNode            IIBNodeType
	SoapReplyNode              IIBNodeType
	SoapWrapperNode            IIBNodeType
	SrRetrieveEntityNode       IIBNodeType
	SrRetrieveItServiceNode    IIBNodeType
	ThrowNode                  IIBNodeType
	TraceNode                  IIBNodeType
	TryCatchNode               IIBNodeType
	ValidateNode               IIBNodeType
	WsReplyNode                IIBNodeType
	XslMqsiNode                IIBNodeType
}{
	"AGGREGATE_CONTROL_NODE",
	"AGGREGATE_REPLY_NODE",
	"AGGREGATE_REQUEST_NODE",
	"CALLABLE_FLOW_REPLY_NODE",
	"COLLECTOR_NODE",
	"COMPUTE_NODE",
	"DATABASE_NODE",
	"DECISION_SERVICE_NODE",
	"DOT_NET_COMPUTE_NODE",
	"FILE_READ_NODE",
	"FILTER_NODE",
	"FLOW_ORDER_NODE",
	"GROUP_COMPLETE_NODE",
	"GROUP_GATHER_NODE",
	"GROUP_SCATTER_NODE",
	"HTTP_HEADER",
	"JAVA_COMPUTE_NODE",
	"JMS_CLIENT_RECEIVE",
	"JMS_CLIENT_REPLY_NODE",
	"JMS_HEADER",
	"MQ_GET_NODE",
	"MQ_OUTPUT_NODE",
	"PASSTHRU_NODE",
	"RESET_CONTENT_DESCRIPTOR_NODE",
	"RE_SEQUENCE_NODE",
	"ROUTE_NODE",
	"SAP_REPLY_NODE",
	"SCA_REPLY_NODE",
	"SECURITY_PEP",
	"SEQUENCE_NODE",
	"SOAP_EXTRACT_NODE",
	"SOAP_REPLY_NODE",
	"SOAP_WRAPPER_NODE",
	"SR_RETRIEVE_ENTITY_NODE",
	"SR_RETRIEVE_IT_SERVICE_NODE",
	"THROW_NODE",
	"TRACE_NODE",
	"TRY_CATCH_NODE",
	"VALIDATE_NODE",
	"WS_REPLY_NODE",
	"XSL_MQSI_NODE",
}

// Position The position of the extracted string relative to delimiters.
type Position string

// Positions offers the known enum values
var Positions = struct {
	After   Position
	Before  Position
	Between Position
}{
	"AFTER",
	"BEFORE",
	"BETWEEN",
}

// Visibility The visibility of the method to capture.
type Visibility string

// Visibilitys offers the known enum values
var Visibilitys = struct {
	Internal         Visibility
	PackageProtected Visibility
	Private          Visibility
	Protected        Visibility
	Public           Visibility
}{
	"INTERNAL",
	"PACKAGE_PROTECTED",
	"PRIVATE",
	"PROTECTED",
	"PUBLIC",
}

// FileNameMatcher The operator of the comparison.
//
//	If not set, `EQUALS` is used.
type FileNameMatcher string

func (me FileNameMatcher) Ref() *FileNameMatcher {
	return &me
}

// FileNameMatchers offers the known enum values
var FileNameMatchers = struct {
	EndsWith   FileNameMatcher
	Equals     FileNameMatcher
	StartsWith FileNameMatcher
}{
	"ENDS_WITH",
	"EQUALS",
	"STARTS_WITH",
}

// Modifier has no documentation
type Modifier string

// Modifiers offers the known enum values
var Modifiers = struct {
	Abstract Modifier
	Extern   Modifier
	Final    Modifier
	Native   Modifier
	Static   Modifier
}{
	"ABSTRACT",
	"EXTERN",
	"FINAL",
	"NATIVE",
	"STATIC",
}

// DataType The data type of the request attribute.
type DataType string

// DataTypes offers the known enum values
var DataTypes = struct {
	Double  DataType
	Integer DataType
	String  DataType
}{
	"DOUBLE",
	"INTEGER",
	"STRING",
}

// Normalization String values transformation.
//
//	If the **dataType** is not `string`, set the `Original` here.
type Normalization string

// Normalizations offers the known enum values
var Normalizations = struct {
	Original    Normalization
	ToLowerCase Normalization
	ToUpperCase Normalization
}{
	"ORIGINAL",
	"TO_LOWER_CASE",
	"TO_UPPER_CASE",
}

// Aggregation Aggregation type for the request values.
type Aggregation string

// Aggregations offers the known enum values
var Aggregations = struct {
	AllDistinctValues   Aggregation
	Average             Aggregation
	CountDistinctValues Aggregation
	CountValues         Aggregation
	First               Aggregation
	Last                Aggregation
	Maximum             Aggregation
	Minimum             Aggregation
	Sum                 Aggregation
}{
	"ALL_DISTINCT_VALUES",
	"AVERAGE",
	"COUNT_DISTINCT_VALUES",
	"COUNT_VALUES",
	"FIRST",
	"LAST",
	"MAXIMUM",
	"MINIMUM",
	"SUM",
}

// ServiceTechnology Only applies to this service technology.
type ServiceTechnology string

func (me ServiceTechnology) Ref() *ServiceTechnology {
	return &me
}

// ServiceTechnologys offers the known enum values
var ServiceTechnologys = struct {
	ActiveMq                             ServiceTechnology
	ActiveMqArtemis                      ServiceTechnology
	AdoNet                               ServiceTechnology
	AIX                                  ServiceTechnology
	Akka                                 ServiceTechnology
	AmazonRedshift                       ServiceTechnology
	Amqp                                 ServiceTechnology
	ApacheCamel                          ServiceTechnology
	ApacheCassandra                      ServiceTechnology
	ApacheCouchDB                        ServiceTechnology
	ApacheDerby                          ServiceTechnology
	ApacheHTTPClientAsync                ServiceTechnology
	ApacheHTTPClientSync                 ServiceTechnology
	ApacheHTTPServer                     ServiceTechnology
	ApacheKafka                          ServiceTechnology
	ApacheSolr                           ServiceTechnology
	ApacheStorm                          ServiceTechnology
	ApacheSynapse                        ServiceTechnology
	ApacheTomcat                         ServiceTechnology
	Apparmor                             ServiceTechnology
	ApplicationInsightsSdk               ServiceTechnology
	ASPDotNet                            ServiceTechnology
	ASPDotNetCore                        ServiceTechnology
	ASPDotNetCoreSignalr                 ServiceTechnology
	ASPDotNetSignalr                     ServiceTechnology
	AsyncHTTPClient                      ServiceTechnology
	AWSLambda                            ServiceTechnology
	AWSRds                               ServiceTechnology
	AWSService                           ServiceTechnology
	Axis                                 ServiceTechnology
	AzureFunctions                       ServiceTechnology
	AzureServiceBus                      ServiceTechnology
	AzureServiceFabric                   ServiceTechnology
	AzureStorage                         ServiceTechnology
	Boshbpm                              ServiceTechnology
	Citrix                               ServiceTechnology
	CitrixCommon                         ServiceTechnology
	CitrixDesktopDeliveryControllers     ServiceTechnology
	CitrixDirector                       ServiceTechnology
	CitrixLicenseServer                  ServiceTechnology
	CitrixProvisioningServices           ServiceTechnology
	CitrixStorefront                     ServiceTechnology
	CitrixVirtualDeliveryAgent           ServiceTechnology
	CitrixWorkspaceEnvironmentManagement ServiceTechnology
	CloudFoundry                         ServiceTechnology
	CloudFoundryAuctioneer               ServiceTechnology
	CloudFoundryBosh                     ServiceTechnology
	CloudFoundryGorouter                 ServiceTechnology
	Coldfusion                           ServiceTechnology
	Containerd                           ServiceTechnology
	CoreDNS                              ServiceTechnology
	Couchbase                            ServiceTechnology
	Crio                                 ServiceTechnology
	Cxf                                  ServiceTechnology
	Datastax                             ServiceTechnology
	DB2                                  ServiceTechnology
	DiegoCell                            ServiceTechnology
	Docker                               ServiceTechnology
	DotNet                               ServiceTechnology
	DotNetRemoting                       ServiceTechnology
	ElasticSearch                        ServiceTechnology
	Envoy                                ServiceTechnology
	Erlang                               ServiceTechnology
	Etcd                                 ServiceTechnology
	F5Ltm                                ServiceTechnology
	Fsharp                               ServiceTechnology
	Garden                               ServiceTechnology
	Glassfish                            ServiceTechnology
	Go                                   ServiceTechnology
	GraalTruffle                         ServiceTechnology
	Grpc                                 ServiceTechnology
	Grsecurity                           ServiceTechnology
	Hadoop                               ServiceTechnology
	HadoopHdfs                           ServiceTechnology
	HadoopYarn                           ServiceTechnology
	Haproxy                              ServiceTechnology
	Heat                                 ServiceTechnology
	Hessian                              ServiceTechnology
	HornetQ                              ServiceTechnology
	IBMCICSRegion                        ServiceTechnology
	IBMCICSTransactionGateway            ServiceTechnology
	IBMIMSConnectRegion                  ServiceTechnology
	IBMIMSControlRegion                  ServiceTechnology
	IBMIMSMessageProcessingRegion        ServiceTechnology
	IBMIMSSoapGateway                    ServiceTechnology
	IBMIntegrationBus                    ServiceTechnology
	IBMMq                                ServiceTechnology
	IBMMqClient                          ServiceTechnology
	IBMWebshprereApplicationServer       ServiceTechnology
	IBMWebshprereLiberty                 ServiceTechnology
	IIS                                  ServiceTechnology
	IISAppPool                           ServiceTechnology
	Istio                                ServiceTechnology
	Java                                 ServiceTechnology
	JaxWs                                ServiceTechnology
	JBoss                                ServiceTechnology
	JBossEap                             ServiceTechnology
	JdkHTTPServer                        ServiceTechnology
	Jersey                               ServiceTechnology
	Jetty                                ServiceTechnology
	Jruby                                ServiceTechnology
	Jython                               ServiceTechnology
	Kubernetes                           ServiceTechnology
	Libvirt                              ServiceTechnology
	Linkerd                              ServiceTechnology
	Mariadb                              ServiceTechnology
	Memcache                             ServiceTechnology
	MicrosoftSQLServer                   ServiceTechnology
	Mongodb                              ServiceTechnology
	MSSQLClient                          ServiceTechnology
	MuleEsb                              ServiceTechnology
	MySQL                                ServiceTechnology
	MySQLConnector                       ServiceTechnology
	NetflixServo                         ServiceTechnology
	Netty                                ServiceTechnology
	Nginx                                ServiceTechnology
	NodeJs                               ServiceTechnology
	OkHTTPClient                         ServiceTechnology
	OneAgentSdk                          ServiceTechnology
	Opencensus                           ServiceTechnology
	Openshift                            ServiceTechnology
	OpenStackCompute                     ServiceTechnology
	OpenStackController                  ServiceTechnology
	Opentelemetry                        ServiceTechnology
	Opentracing                          ServiceTechnology
	OpenLiberty                          ServiceTechnology
	OracleDatabase                       ServiceTechnology
	OracleWeblogic                       ServiceTechnology
	Owin                                 ServiceTechnology
	Perl                                 ServiceTechnology
	PHP                                  ServiceTechnology
	PHPFpm                               ServiceTechnology
	Play                                 ServiceTechnology
	PostgreSQL                           ServiceTechnology
	PostgreSQLDotNetDataProvider         ServiceTechnology
	PowerDNS                             ServiceTechnology
	Progress                             ServiceTechnology
	Python                               ServiceTechnology
	RabbitMq                             ServiceTechnology
	Redis                                ServiceTechnology
	Resteasy                             ServiceTechnology
	Restlet                              ServiceTechnology
	Riak                                 ServiceTechnology
	Ruby                                 ServiceTechnology
	SagWebmethodsIs                      ServiceTechnology
	SAP                                  ServiceTechnology
	SAPHanadb                            ServiceTechnology
	SAPHybris                            ServiceTechnology
	SAPMaxdb                             ServiceTechnology
	SAPSybase                            ServiceTechnology
	Scala                                ServiceTechnology
	Selinux                              ServiceTechnology
	Sharepoint                           ServiceTechnology
	Spark                                ServiceTechnology
	Spring                               ServiceTechnology
	Sqlite                               ServiceTechnology
	Thrift                               ServiceTechnology
	Tibco                                ServiceTechnology
	TibcoBusinessWorks                   ServiceTechnology
	TibcoEms                             ServiceTechnology
	VarnishCache                         ServiceTechnology
	Vim2                                 ServiceTechnology
	VirtualMachineKvm                    ServiceTechnology
	VirtualMachineQemu                   ServiceTechnology
	Wildfly                              ServiceTechnology
	WindowsContainers                    ServiceTechnology
	Wink                                 ServiceTechnology
	ZeroMq                               ServiceTechnology
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

// Operator Operator comparing the extracted value to the comparison value.
type Operator string

// Operators offers the known enum values
var Operators = struct {
	BeginsWith Operator
	Contains   Operator
	EndsWith   Operator
	Equals     Operator
}{
	"BEGINS_WITH",
	"CONTAINS",
	"ENDS_WITH",
	"EQUALS",
}
