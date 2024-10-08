/**
* @license
* Copyright 2024 Dynatrace LLC
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

package httpcache

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
)

type client struct {
	client   rest.Client
	schemaID string
}

var CACHE_FOLDER = os.Getenv("DYNATRACE_MIGRATION_CACHE_FOLDER")
var STRICT_CACHE = os.Getenv("DYNATRACE_MIGRATION_CACHE_STRICT") == "true"

func DefaultClient(envURL string, apiToken string, schemaID string) rest.Client {
	restClient := rest.DefaultClient(envURL, apiToken)
	if len(CACHE_FOLDER) > 0 {
		return Client(restClient, schemaID)
	}
	return restClient
}

func Client(c rest.Client, schemaID string) rest.Client {
	return &client{client: c, schemaID: schemaID}
}

var REGEX_SETTINGS_20_LIST, _ = regexp.Compile(`\/api\/v2\/settings\/objects\?schemaIds=([^\&]*)&`)
var REGEX_SETTINGS_20_GET, _ = regexp.Compile(`\/api\/v2\/settings\/objects\/(.*)`)
var REGEX_APPLICATIONS_MOBILE_LIST, _ = regexp.Compile(`\/api\/config\/v1\/applications\/mobile$`)
var REGEX_APPLICATIONS_MOBILE_GET, _ = regexp.Compile(`\/api\/config\/v1\/applications\/mobile\/([^\/]*)$`)
var REGEX_APPLICATIONS_MOBILE_KEY_USER_ACTIONS_GET, _ = regexp.Compile(`\/api\/config\/v1\/applications\/mobile\/([^\/]*)\/keyUserActions$`)
var REGEX_APPLICATIONS_MOBILE_KEY_USER_ACTION_AND_SESSION_PROPERTIES_LIST, _ = regexp.Compile(`\/api\/config\/v1\/applications\/mobile\/([^\/]*)\/userActionAndSessionProperties$`)
var REGEX_APPLICATIONS_MOBILE_KEY_USER_ACTION_AND_SESSION_PROPERTIES_GET, _ = regexp.Compile(`\/api\/config\/v1\/applications\/mobile\/([^\/]*)\/userActionAndSessionProperties/([^\/]*)$`)
var REGEX_APPLICATIONS_WEB_LIST, _ = regexp.Compile(`\/api\/config\/v1\/applications\/web$`)
var REGEX_APPLICATIONS_WEB_GET, _ = regexp.Compile(`\/api\/config\/v1\/applications\/web\/([^\/]*)$`)
var REGEX_APPLICATIONS_WEB_KEY_USER_ACTIONS_GET, _ = regexp.Compile(`\/api\/config\/v1\/applications\/web\/([^\/]*)\/keyUserActions$`)
var REGEX_APPLICATIONS_WEB_ERROR_RULES_GET, _ = regexp.Compile(`\/api\/config\/v1\/applications\/web\/([^\/]*)\/errorRules$`)
var REGEX_APPLICATION_DETECTION_RULES_LIST, _ = regexp.Compile(`\/api\/config\/v1\/applicationDetectionRules$`)
var REGEX_APPLICATION_DETECTION_RULES_GET, _ = regexp.Compile(`\/api\/config\/v1\/applicationDetectionRules\/([^\/]*)$`)
var REGEX_AWS_CREDENTIALS_LIST, _ = regexp.Compile(`\/api\/config\/v1\/aws\/credentials$`)
var REGEX_AWS_CREDENTIALS_GET, _ = regexp.Compile(`\/api\/config\/v1\/aws\/credentials\/([^\/]*)$`)
var REGEX_AZURE_CREDENTIALS_LIST, _ = regexp.Compile(`\/api\/config\/v1\/azure\/credentials$`)
var REGEX_AZURE_CREDENTIALS_GET, _ = regexp.Compile(`\/api\/config\/v1\/azure\/credentials\/([^\/]*)$`)
var REGEX_CUSTOM_SERVICE_NODEJS_LIST, _ = regexp.Compile(`\/api\/config\/v1\/service\/customServices\/nodeJS$`)
var REGEX_CUSTOM_SERVICE_NODEJS_GET, _ = regexp.Compile(`\/api\/config\/v1\/service\/customServices\/nodeJS\/([^\/]*)$`)
var REGEX_CUSTOM_SERVICE_DOTNET_LIST, _ = regexp.Compile(`\/api\/config\/v1\/service\/customServices\/dotNet$`)
var REGEX_CUSTOM_SERVICE_DOTNET_GET, _ = regexp.Compile(`\/api\/config\/v1\/service\/customServices\/dotNet\/([^\/]*)$`)
var REGEX_CUSTOM_SERVICE_GOLANG_LIST, _ = regexp.Compile(`\/api\/config\/v1\/service\/customServices\/go$`)
var REGEX_CUSTOM_SERVICE_GOLANG_GET, _ = regexp.Compile(`\/api\/config\/v1\/service\/customServices\/go\/([^\/]*)$`)
var REGEX_CUSTOM_SERVICE_JAVA_LIST, _ = regexp.Compile(`\/api\/config\/v1\/service\/customServices\/java$`)
var REGEX_CUSTOM_SERVICE_JAVA_GET, _ = regexp.Compile(`\/api\/config\/v1\/service\/customServices\/java\/([^\/]*)$`)
var REGEX_CUSTOM_SERVICE_PHP_LIST, _ = regexp.Compile(`\/api\/config\/v1\/service\/customServices\/php$`)
var REGEX_CUSTOM_SERVICE_PHP_GET, _ = regexp.Compile(`\/api\/config\/v1\/service\/customServices\/php\/([^\/]*)$`)
var REGEX_CALCULATED_METRICS_SERVICE_LIST, _ = regexp.Compile(`\/api\/config\/v1\/calculatedMetrics\/service$`)
var REGEX_CALCULATED_METRICS_SERVICE_GET, _ = regexp.Compile(`\/api\/config\/v1\/calculatedMetrics\/service\/([^\/]*)$`)
var REGEX_CALCULATED_METRICS_SYNTHETIC_LIST, _ = regexp.Compile(`\/api\/config\/v1\/calculatedMetrics\/synthetic$`)
var REGEX_CALCULATED_METRICS_SYNTHETIC_GET, _ = regexp.Compile(`\/api\/config\/v1\/calculatedMetrics\/synthetic\/([^\/]*)$`)
var REGEX_CALCULATED_METRICS_MOBILE_LIST, _ = regexp.Compile(`\/api\/config\/v1\/calculatedMetrics\/mobile$`)
var REGEX_CALCULATED_METRICS_MOBILE_GET, _ = regexp.Compile(`\/api\/config\/v1\/calculatedMetrics\/mobile\/([^\/]*)$`)
var REGEX_CALCULATED_METRICS_RUM_LIST, _ = regexp.Compile(`\/api\/config\/v1\/calculatedMetrics\/rum$`)
var REGEX_CALCULATED_METRICS_RUM_GET, _ = regexp.Compile(`\/api\/config\/v1\/calculatedMetrics\/rum\/([^\/]*)$`)
var REGEX_CONDITIONAL_NAMING_HOST_LIST, _ = regexp.Compile(`\/api\/config\/v1\/conditionalNaming\/host$`)
var REGEX_CONDITIONAL_NAMING_HOST_GET, _ = regexp.Compile(`\/api\/config\/v1\/conditionalNaming\/host\/([^\/]*)$`)
var REGEX_CONDITIONAL_NAMING_PG_LIST, _ = regexp.Compile(`\/api\/config\/v1\/conditionalNaming\/processGroup$`)
var REGEX_CONDITIONAL_NAMING_PG_GET, _ = regexp.Compile(`\/api\/config\/v1\/conditionalNaming\/processGroup\/([^\/]*)$`)
var REGEX_CONDITIONAL_NAMING_SERVICE_LIST, _ = regexp.Compile(`\/api\/config\/v1\/conditionalNaming\/service$`)
var REGEX_CONDITIONAL_NAMING_SERVICE_GET, _ = regexp.Compile(`\/api\/config\/v1\/conditionalNaming\/service\/([^\/]*)$`)
var REGEX_CREDENTIALS_VAULT_LIST, _ = regexp.Compile(`\/api\/v2\/credentials$`)
var REGEX_CREDENTIALS_VAULT_GET, _ = regexp.Compile(`\/api\/v2\/credentials\/([^\/]*)$`)
var REGEX_DASHBOARDS_LIST, _ = regexp.Compile(`\/api\/config\/v1\/dashboards$`)
var REGEX_DASHBOARDS_GET, _ = regexp.Compile(`\/api\/config\/v1\/dashboards\/([^\/]*)$`)
var REGEX_DASHBOARD_SHARING_GET, _ = regexp.Compile(`\/api\/config\/v1\/dashboards\/([^\/]*)\/shareSettings$`)
var REGEX_REQUEST_ATTRIBUTES_LIST, _ = regexp.Compile(`\/api\/config\/v1\/service\/requestAttributes$`)
var REGEX_REQUEST_ATTRIBUTES_GET, _ = regexp.Compile(`\/api\/config\/v1\/service\/requestAttributes\/([^\/?]*)\?includeProcessGroupReferences=true$`)
var REGEX_REQUEST_NAMING_LIST, _ = regexp.Compile(`\/api\/config\/v1\/service\/requestNaming$`)
var REGEX_REQUEST_NAMING_GET, _ = regexp.Compile(`\/api\/config\/v1\/service\/requestNaming\/([^\/]*)$`)
var REGEX_REQUEST_SLO_LIST, _ = regexp.Compile(`\/api\/v2\/slo\?pageSize`)
var REGEX_REQUEST_SLO_GET, _ = regexp.Compile(`\/api\/v2\/slo\/([^\/]*)$`)
var REGEX_NETWORK_ZONES_LIST, _ = regexp.Compile(`\/api\/v2\/networkZones`)
var REGEX_NETWORK_ZONES_GET, _ = regexp.Compile(`\/api\/v2\/networkZones\/([^\/]*)$`)
var REGEX_BROWSER_MONITOR_LIST, _ = regexp.Compile(`\/api\/v1\/synthetic\/monitors\?type=BROWSER$`)
var REGEX_MONITOR_GET, _ = regexp.Compile(`\/api\/v1\/synthetic\/monitors\/([^\/]*)$`)
var REGEX_HTTP_MONITOR_LIST, _ = regexp.Compile(`\/api\/v1\/synthetic\/monitors\?type=HTTP$`)
var REGEX_PRIVATE_SYNTHETIC_LOCATIONS_LIST, _ = regexp.Compile(`\/api\/v1\/synthetic\/locations\?type=PRIVATE$`)
var REGEX_PRIVATE_SYNTHETIC_LOCATIONS_GET, _ = regexp.Compile(`\/api\/v1\/synthetic\/locations\/([^\/]*)$`)
var REGEX_CREDENTIALS_LIST, _ = regexp.Compile(`\/api\/config\/v1\/credentials$`)

func (me *client) Get(ctx context.Context, url string, expectedStatusCodes ...int) rest.Request {
	doGet := func(modSchemaID string, id string) rest.Request {
		return Get(modSchemaID, id, me.schemaID)
	}

	if m := REGEX_SETTINGS_20_LIST.FindStringSubmatch(url); len(m) == 2 {
		return &ListSettings20Request{SchemaID: me.schemaID}
	} else if m := REGEX_SETTINGS_20_GET.FindStringSubmatch(url); len(m) == 2 {
		return &GetSettings20Request{SchemaID: me.schemaID, ID: m[1]}
	} else if m := REGEX_APPLICATIONS_MOBILE_LIST.FindStringSubmatch(url); len(m) == 1 {
		return List("application-mobile")
	} else if m := REGEX_APPLICATIONS_MOBILE_GET.FindStringSubmatch(url); len(m) == 2 {
		return doGet("application-mobile", m[1])
	} else if m := REGEX_APPLICATIONS_MOBILE_KEY_USER_ACTIONS_GET.FindStringSubmatch(url); len(m) == 2 {
		return doGet("application-mobile-key-user-actions", m[1])
	} else if m := REGEX_APPLICATIONS_MOBILE_KEY_USER_ACTION_AND_SESSION_PROPERTIES_LIST.FindStringSubmatch(url); len(m) == 2 {
		return doGet("application-mobile-user-actions-and-session-properties", m[1])
	} else if m := REGEX_APPLICATIONS_MOBILE_KEY_USER_ACTION_AND_SESSION_PROPERTIES_GET.FindStringSubmatch(url); len(m) == 3 {
		return doGet("application-mobile-user-actions-and-session-properties-remote-properties", m[1]+":"+m[2])
	} else if m := REGEX_APPLICATIONS_WEB_LIST.FindStringSubmatch(url); len(m) == 1 {
		return List("application-web")
	} else if m := REGEX_APPLICATIONS_WEB_GET.FindStringSubmatch(url); len(m) == 2 {
		return doGet("application-web", m[1])
	} else if m := REGEX_APPLICATIONS_WEB_KEY_USER_ACTIONS_GET.FindStringSubmatch(url); len(m) == 2 {
		return doGet("application-web-key-user-actions", m[1])
	} else if m := REGEX_APPLICATIONS_WEB_ERROR_RULES_GET.FindStringSubmatch(url); len(m) == 2 {
		return doGet("application-web-error-rules", m[1])
	} else if m := REGEX_APPLICATION_DETECTION_RULES_LIST.FindStringSubmatch(url); len(m) == 1 {
		return List("app-detection-rule")
	} else if m := REGEX_APPLICATION_DETECTION_RULES_GET.FindStringSubmatch(url); len(m) == 2 {
		return doGet("app-detection-rule", m[1])
	} else if m := REGEX_AWS_CREDENTIALS_LIST.FindStringSubmatch(url); len(m) == 1 {
		return List("aws-credentials")
	} else if m := REGEX_AWS_CREDENTIALS_GET.FindStringSubmatch(url); len(m) == 2 {
		return doGet("aws-credentials", m[1])
	} else if m := REGEX_AZURE_CREDENTIALS_LIST.FindStringSubmatch(url); len(m) == 1 {
		return List("azure-credentials")
	} else if m := REGEX_AZURE_CREDENTIALS_GET.FindStringSubmatch(url); len(m) == 2 {
		return doGet("azure-credentials", m[1])
	} else if m := REGEX_CUSTOM_SERVICE_NODEJS_LIST.FindStringSubmatch(url); len(m) == 1 {
		return List("custom-service-nodejs")
	} else if m := REGEX_CUSTOM_SERVICE_NODEJS_GET.FindStringSubmatch(url); len(m) == 2 {
		return doGet("custom-service-nodejs", m[1])
	} else if m := REGEX_CUSTOM_SERVICE_DOTNET_LIST.FindStringSubmatch(url); len(m) == 1 {
		return List("custom-service-dotnet")
	} else if m := REGEX_CUSTOM_SERVICE_DOTNET_GET.FindStringSubmatch(url); len(m) == 2 {
		return doGet("custom-service-dotnet", m[1])
	} else if m := REGEX_CUSTOM_SERVICE_GOLANG_LIST.FindStringSubmatch(url); len(m) == 1 {
		return List("custom-service-go")
	} else if m := REGEX_CUSTOM_SERVICE_GOLANG_GET.FindStringSubmatch(url); len(m) == 2 {
		return doGet("custom-service-go", m[1])
	} else if m := REGEX_CUSTOM_SERVICE_JAVA_LIST.FindStringSubmatch(url); len(m) == 1 {
		return List("custom-service-java")
	} else if m := REGEX_CUSTOM_SERVICE_JAVA_GET.FindStringSubmatch(url); len(m) == 2 {
		return doGet("custom-service-java", m[1])
	} else if m := REGEX_CUSTOM_SERVICE_PHP_LIST.FindStringSubmatch(url); len(m) == 1 {
		return List("custom-service-php")
	} else if m := REGEX_CUSTOM_SERVICE_PHP_GET.FindStringSubmatch(url); len(m) == 2 {
		return doGet("custom-service-php", m[1])
	} else if m := REGEX_CALCULATED_METRICS_SERVICE_LIST.FindStringSubmatch(url); len(m) == 1 {
		return List("calculated-metrics-service")
	} else if m := REGEX_CALCULATED_METRICS_SERVICE_GET.FindStringSubmatch(url); len(m) == 2 {
		return doGet("calculated-metrics-service", m[1])
	} else if m := REGEX_CALCULATED_METRICS_SYNTHETIC_LIST.FindStringSubmatch(url); len(m) == 1 {
		return List("calculated-metrics-synthetic")
	} else if m := REGEX_CALCULATED_METRICS_SYNTHETIC_GET.FindStringSubmatch(url); len(m) == 2 {
		return doGet("calculated-metrics-synthetic", m[1])
	} else if m := REGEX_CALCULATED_METRICS_MOBILE_LIST.FindStringSubmatch(url); len(m) == 1 {
		return List("calculated-metrics-application-mobile")
	} else if m := REGEX_CALCULATED_METRICS_MOBILE_GET.FindStringSubmatch(url); len(m) == 2 {
		return doGet("calculated-metrics-application-mobile", m[1])
	} else if m := REGEX_CALCULATED_METRICS_RUM_LIST.FindStringSubmatch(url); len(m) == 1 {
		return List("calculated-metrics-application-web")
	} else if m := REGEX_CALCULATED_METRICS_RUM_GET.FindStringSubmatch(url); len(m) == 2 {
		return doGet("calculated-metrics-application-web", m[1])
	} else if m := REGEX_CONDITIONAL_NAMING_HOST_LIST.FindStringSubmatch(url); len(m) == 1 {
		return List("conditional-naming-host")
	} else if m := REGEX_CONDITIONAL_NAMING_HOST_GET.FindStringSubmatch(url); len(m) == 2 {
		return doGet("conditional-naming-host", m[1])
	} else if m := REGEX_CONDITIONAL_NAMING_PG_LIST.FindStringSubmatch(url); len(m) == 1 {
		return List("conditional-naming-processgroup")
	} else if m := REGEX_CONDITIONAL_NAMING_PG_GET.FindStringSubmatch(url); len(m) == 2 {
		return doGet("conditional-naming-processgroup", m[1])
	} else if m := REGEX_CONDITIONAL_NAMING_SERVICE_LIST.FindStringSubmatch(url); len(m) == 1 {
		return List("conditional-naming-service")
	} else if m := REGEX_CONDITIONAL_NAMING_SERVICE_GET.FindStringSubmatch(url); len(m) == 2 {
		return doGet("conditional-naming-service", m[1])
	} else if m := REGEX_CREDENTIALS_VAULT_LIST.FindStringSubmatch(url); len(m) == 1 {
		return &ListCredentialsVaultRequest{SchemaID: "credentials-vault"}
	} else if m := REGEX_CREDENTIALS_VAULT_GET.FindStringSubmatch(url); len(m) == 2 {
		return &GetCredentialsVaultRequest{SchemaID: "credentials-vault", ID: m[1], ServiceSchemaID: me.schemaID}
	} else if m := REGEX_DASHBOARDS_LIST.FindStringSubmatch(url); len(m) == 1 {
		return Request(&ListDashboardsV1{})
	} else if m := REGEX_DASHBOARD_SHARING_GET.FindStringSubmatch(url); len(m) == 2 {
		return doGet("dashboard-sharing", m[1])
	} else if m := REGEX_DASHBOARDS_GET.FindStringSubmatch(url); len(m) == 2 {
		return doGet("dashboard", m[1])
	} else if m := REGEX_NETWORK_ZONES_LIST.FindStringSubmatch(url); len(m) == 1 {
		return List("network-zone")
	} else if m := REGEX_NETWORK_ZONES_GET.FindStringSubmatch(url); len(m) == 2 {
		return doGet("network-zone", m[1])
	} else if m := REGEX_REQUEST_ATTRIBUTES_LIST.FindStringSubmatch(url); len(m) == 1 {
		return List("request-attributes")
	} else if m := REGEX_REQUEST_ATTRIBUTES_GET.FindStringSubmatch(url); len(m) == 2 {
		return doGet("request-attributes", m[1])
	} else if m := REGEX_REQUEST_NAMING_LIST.FindStringSubmatch(url); len(m) == 1 {
		return List("request-naming-service")
	} else if m := REGEX_REQUEST_NAMING_GET.FindStringSubmatch(url); len(m) == 2 {
		return doGet("request-naming-service", m[1])
	} else if m := REGEX_REQUEST_SLO_LIST.FindStringSubmatch(url); len(m) == 1 {
		return &ListSLORequest{SchemaID: "slo"}
	} else if m := REGEX_REQUEST_SLO_GET.FindStringSubmatch(url); len(m) == 2 {
		return &GetSLORequest{SchemaID: "slo", ID: m[1], ServiceSchemaID: me.schemaID}
	} else if m := REGEX_BROWSER_MONITOR_LIST.FindStringSubmatch(url); len(m) == 1 {
		return Request(&ListMonitorsV1{"SYNTHETIC_TEST-"})
	} else if m := REGEX_HTTP_MONITOR_LIST.FindStringSubmatch(url); len(m) == 1 {
		return Request(&ListMonitorsV1{"HTTP_CHECK-"})
	} else if m := REGEX_MONITOR_GET.FindStringSubmatch(url); len(m) == 2 {
		return doGet("synthetic-monitor", m[1])
	} else if m := REGEX_PRIVATE_SYNTHETIC_LOCATIONS_LIST.FindStringSubmatch(url); len(m) == 1 {
		return Request(&ListPrivateSyntheticLocationsV1{})
	} else if m := REGEX_PRIVATE_SYNTHETIC_LOCATIONS_GET.FindStringSubmatch(url); len(m) == 2 {
		return doGet("synthetic-location", m[1])
	} else if m := REGEX_CREDENTIALS_LIST.FindStringSubmatch(url); len(m) == 1 {
		return EmptyList("v1:config:credentials")
	}

	if STRICT_CACHE {
		return &rest.StriclyForbidden{Method: "GET", URL: url}
	}
	logging.Debug.Info.Printf("[HTTP_CACHE] [%s] [FAILED] [%s] Could not find a cache match", me.schemaID, url)
	logging.Debug.Warn.Printf("[HTTP_CACHE] [%s] [FAILED] [%s] Could not find a cache match", me.schemaID, url)
	fmt.Printf("\nERROR: Could not find a cache match for: %s, url: %s", me.schemaID, url)
	return me.client.Get(ctx, url, expectedStatusCodes...)
}

func (me *client) Post(ctx context.Context, url string, payload any, expectedStatusCodes ...int) rest.Request {
	return &rest.Forbidden{Method: "POST"}
}

func (me *client) Put(ctx context.Context, url string, payload any, expectedStatusCodes ...int) rest.Request {
	return &rest.Forbidden{Method: "PUT"}
}

func (me *client) Delete(ctx context.Context, url string, expectedStatusCodes ...int) rest.Request {
	return &rest.Forbidden{Method: "DELETE"}
}

func (me *client) Upload(ctx context.Context, url string, reader io.ReadCloser, fileName string, expectedStatusCodes ...int) rest.Request {
	return &rest.Forbidden{Method: "Upload"}
}
