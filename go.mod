module github.com/dynatrace-oss/terraform-provider-dynatrace

go 1.19

require (
	github.com/dtcookie/assert v1.2.0
	github.com/dtcookie/dynatrace/api/accounts/iam v1.0.2
	github.com/dtcookie/dynatrace/api/cluster/v1/groups v1.0.1
	github.com/dtcookie/dynatrace/api/cluster/v1/users v1.0.1
	github.com/dtcookie/dynatrace/api/cluster/v2/envs v1.0.4
	github.com/dtcookie/dynatrace/api/config v1.0.11
	github.com/dtcookie/dynatrace/api/config/alerting v1.0.23
	github.com/dtcookie/dynatrace/api/config/anomalies/applications v1.0.8
	github.com/dtcookie/dynatrace/api/config/anomalies/databaseservices v1.0.8
	github.com/dtcookie/dynatrace/api/config/anomalies/diskevents v1.0.13
	github.com/dtcookie/dynatrace/api/config/anomalies/hosts v1.0.9
	github.com/dtcookie/dynatrace/api/config/anomalies/metricevents v1.0.17
	github.com/dtcookie/dynatrace/api/config/anomalies/services v1.0.14
	github.com/dtcookie/dynatrace/api/config/applications/mobile v1.0.6
	github.com/dtcookie/dynatrace/api/config/applications/web v1.0.17
	github.com/dtcookie/dynatrace/api/config/applications/web/applicationdetectionrules v1.0.3
	github.com/dtcookie/dynatrace/api/config/autotags v1.0.27
	github.com/dtcookie/dynatrace/api/config/credentials/aws v1.0.18
	github.com/dtcookie/dynatrace/api/config/credentials/azure v1.0.14
	github.com/dtcookie/dynatrace/api/config/credentials/cloudfoundry v1.0.3
	github.com/dtcookie/dynatrace/api/config/credentials/kubernetes v1.0.16
	github.com/dtcookie/dynatrace/api/config/credentials/vault v1.0.7
	github.com/dtcookie/dynatrace/api/config/customservices v1.0.22
	github.com/dtcookie/dynatrace/api/config/dashboards v1.0.23
	github.com/dtcookie/dynatrace/api/config/dashboards/sharing v1.0.5
	github.com/dtcookie/dynatrace/api/config/maintenance v1.0.12
	github.com/dtcookie/dynatrace/api/config/managementzones v1.0.20
	github.com/dtcookie/dynatrace/api/config/metrics/calculated/service v1.0.10
	github.com/dtcookie/dynatrace/api/config/naming/hosts v1.0.5
	github.com/dtcookie/dynatrace/api/config/naming/processgroups v1.0.5
	github.com/dtcookie/dynatrace/api/config/naming/services v1.0.5
	github.com/dtcookie/dynatrace/api/config/networkzones v1.0.1
	github.com/dtcookie/dynatrace/api/config/notifications v1.0.13
	github.com/dtcookie/dynatrace/api/config/requestattributes v1.0.14
	github.com/dtcookie/dynatrace/api/config/requestnaming v1.0.2
	github.com/dtcookie/dynatrace/api/config/synthetic/locations v1.0.1
	github.com/dtcookie/dynatrace/api/config/synthetic/monitors v1.0.22
	github.com/dtcookie/dynatrace/api/config/topology/host v1.0.2
	github.com/dtcookie/dynatrace/api/config/topology/process v1.0.3
	github.com/dtcookie/dynatrace/api/config/topology/processgroup v1.0.2
	github.com/dtcookie/dynatrace/api/config/topology/service v1.0.4
	github.com/dtcookie/dynatrace/api/config/topology/tag v1.0.1
	github.com/dtcookie/dynatrace/api/config/v2/alerting v1.0.7
	github.com/dtcookie/dynatrace/api/config/v2/anomalies/frequentissues v1.0.2
	github.com/dtcookie/dynatrace/api/config/v2/anomalies/metricevents v1.0.3
	github.com/dtcookie/dynatrace/api/config/v2/apitokens v1.0.2
	github.com/dtcookie/dynatrace/api/config/v2/entities v1.0.0
	github.com/dtcookie/dynatrace/api/config/v2/ibmmq/filters v1.0.3
	github.com/dtcookie/dynatrace/api/config/v2/ibmmq/imsbridges v1.0.3
	github.com/dtcookie/dynatrace/api/config/v2/ibmmq/queuemanagers v1.0.3
	github.com/dtcookie/dynatrace/api/config/v2/ibmmq/queuesharinggroups v1.0.3
	github.com/dtcookie/dynatrace/api/config/v2/keyrequests v1.0.6
	github.com/dtcookie/dynatrace/api/config/v2/maintenance v1.0.3
	github.com/dtcookie/dynatrace/api/config/v2/managementzones v1.0.3
	github.com/dtcookie/dynatrace/api/config/v2/networkzones v1.0.5
	github.com/dtcookie/dynatrace/api/config/v2/notifications v1.0.8
	github.com/dtcookie/dynatrace/api/config/v2/slo v1.0.9
	github.com/dtcookie/dynatrace/api/config/v2/spans/attributes v1.0.5
	github.com/dtcookie/dynatrace/api/config/v2/spans/capture v1.0.4
	github.com/dtcookie/dynatrace/api/config/v2/spans/ctxprop v1.0.3
	github.com/dtcookie/dynatrace/api/config/v2/spans/entrypoints v1.0.4
	github.com/dtcookie/dynatrace/api/config/v2/spans/resattr v1.0.7
	github.com/dtcookie/dynatrace/rest v1.0.15
	github.com/dtcookie/dynatrace/terraform v1.0.5
	github.com/dtcookie/hcl v1.0.2
	github.com/dtcookie/opt v1.0.0
	github.com/google/uuid v1.3.0
	github.com/hashicorp/hcl/v2 v2.15.0
	github.com/hashicorp/terraform-plugin-docs v0.13.0
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.24.1
	github.com/zclconf/go-cty v1.12.1
)

require (
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.1.1 // indirect
	github.com/Masterminds/sprig/v3 v3.2.2 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-cidr v1.1.0 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/bgentry/speakeasy v0.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dtcookie/dynatrace/api/config/anomalies/common v1.0.12 // indirect
	github.com/dtcookie/dynatrace/api/config/common v1.0.6 // indirect
	github.com/dtcookie/dynatrace/api/config/entityruleengine v1.0.12 // indirect
	github.com/dtcookie/dynatrace/api/config/v2/common v1.0.1 // indirect
	github.com/dtcookie/dynatrace/api/config/v2/spans/match v1.0.2 // indirect
	github.com/dtcookie/gojson v0.9.1 // indirect
	github.com/dtcookie/xjson v1.0.2 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-checkpoint v0.5.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320 // indirect
	github.com/hashicorp/go-hclog v1.4.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-plugin v1.4.6 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hashicorp/hc-install v0.4.0 // indirect
	github.com/hashicorp/logutils v1.0.0 // indirect
	github.com/hashicorp/terraform-exec v0.17.3 // indirect
	github.com/hashicorp/terraform-json v0.14.0 // indirect
	github.com/hashicorp/terraform-plugin-go v0.14.2 // indirect
	github.com/hashicorp/terraform-plugin-log v0.7.0 // indirect
	github.com/hashicorp/terraform-registry-address v0.1.0 // indirect
	github.com/hashicorp/terraform-svchost v0.0.0-20200729002733-f050f53b9734 // indirect
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.13 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mitchellh/cli v1.1.4 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/oklog/run v1.0.0 // indirect
	github.com/posener/complete v1.2.3 // indirect
	github.com/russross/blackfriday v1.6.0 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	github.com/vmihailenco/msgpack/v4 v4.3.12 // indirect
	github.com/vmihailenco/tagparser v0.1.2 // indirect
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d // indirect
	golang.org/x/net v0.4.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/text v0.5.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20200711021454-869866162049 // indirect
	google.golang.org/grpc v1.51.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

replace github.com/dtcookie/dynatrace/api/config/v2/consumption => /Users/adrian.gonciarz/go/src/github.com/dtcookie/dynatrace/api/config/v2/consumption