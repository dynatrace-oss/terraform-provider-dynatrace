package httpcache

import (
	"os"
	"regexp"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
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

var REGEX_SETTINGS_20_LIST, _ = regexp.Compile("\\/api\\/v2\\/settings\\/objects\\?schemaIds=([^\\&]*)&")
var REGEX_SETTINGS_20_GET, _ = regexp.Compile("\\/api\\/v2\\/settings\\/objects\\/(.*)")
var REGEX_APPLICATIONS_MOBILE_LIST, _ = regexp.Compile("\\/api\\/config\\/v1\\/applications\\/mobile$")
var REGEX_APPLICATIONS_MOBILE_GET, _ = regexp.Compile("\\/api\\/config\\/v1\\/applications\\/mobile\\/([^\\/]*)$")
var REGEX_APPLICATIONS_MOBILE_KEY_USER_ACTIONS_LIST, _ = regexp.Compile("\\/api\\/config\\/v1\\/applications\\/mobile\\/.*\\/keyUserActions$")
var REGEX_APPLICATIONS_MOBILE_KEY_USER_ACTION_AND_SESSION_PROPERTIES_LIST, _ = regexp.Compile("\\/api\\/config\\/v1\\/applications\\/mobile\\/.*\\/userActionAndSessionProperties$")
var REGEX_APPLICATIONS_WEB_LIST, _ = regexp.Compile("\\/api\\/config\\/v1\\/applications\\/web$")
var REGEX_APPLICATIONS_WEB_GET, _ = regexp.Compile("\\/api\\/config\\/v1\\/applications\\/web\\/([^\\/]*)$")
var REGEX_APPLICATIONS_WEB_KEY_USER_ACTIONS_LIST, _ = regexp.Compile("\\/api\\/config\\/v1\\/applications\\/web\\/.*\\/keyUserActions$")
var REGEX_APPLICATIONS_WEB_KEY_USER_ACTION_AND_SESSION_PROPERTIES_LIST, _ = regexp.Compile("\\/api\\/config\\/v1\\/applications\\/web\\/.*\\/userActionAndSessionProperties$")

func (me *client) Get(url string, expectedStatusCodes ...int) rest.Request {
	if m := REGEX_SETTINGS_20_LIST.FindStringSubmatch(url); len(m) == 2 {
		return &ListSettings20Request{SchemaID: me.schemaID}
	} else if m := REGEX_SETTINGS_20_GET.FindStringSubmatch(url); len(m) == 2 {
		return &GetSettings20Request{SchemaID: me.schemaID, ID: m[1]}
	} else if m := REGEX_APPLICATIONS_MOBILE_LIST.FindStringSubmatch(url); len(m) == 1 {
		return List("application-mobile")
	} else if m := REGEX_APPLICATIONS_MOBILE_GET.FindStringSubmatch(url); len(m) == 2 {
		return Get("application-mobile", m[1])
	} else if m := REGEX_APPLICATIONS_MOBILE_KEY_USER_ACTIONS_LIST.FindStringSubmatch(url); len(m) == 1 {
		return EmptyList(me.schemaID + ":" + "key-user-actions")
	} else if m := REGEX_APPLICATIONS_MOBILE_KEY_USER_ACTION_AND_SESSION_PROPERTIES_LIST.FindStringSubmatch(url); len(m) == 1 {
		return EmptyList(me.schemaID + ":" + "user-actions-and-session-properties")
	} else if m := REGEX_APPLICATIONS_WEB_LIST.FindStringSubmatch(url); len(m) == 1 {
		return List("application-web")
	} else if m := REGEX_APPLICATIONS_WEB_GET.FindStringSubmatch(url); len(m) == 2 {
		return Get("application-web", m[1])
	} else if m := REGEX_APPLICATIONS_WEB_KEY_USER_ACTIONS_LIST.FindStringSubmatch(url); len(m) == 1 {
		return EmptyList(me.schemaID + ":" + "key-user-actions")
	} else if m := REGEX_APPLICATIONS_WEB_KEY_USER_ACTION_AND_SESSION_PROPERTIES_LIST.FindStringSubmatch(url); len(m) == 1 {
		return EmptyList(me.schemaID + ":" + "user-actions-and-session-properties")
	}
	if STRICT_CACHE {
		return &rest.StriclyForbidden{Method: "GET", URL: url}
	}
	return me.client.Get(url, expectedStatusCodes...)
}

func (me *client) Post(url string, payload any, expectedStatusCodes ...int) rest.Request {
	return &rest.Forbidden{Method: "POST"}
}

func (me *client) Put(url string, payload any, expectedStatusCodes ...int) rest.Request {
	return &rest.Forbidden{Method: "PUT"}
}

func (me *client) Delete(url string, expectedStatusCodes ...int) rest.Request {
	return &rest.Forbidden{Method: "DELETE"}
}
