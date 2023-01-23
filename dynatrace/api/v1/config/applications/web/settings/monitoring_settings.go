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

package web

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web/settings/browser"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web/settings/ipaddress"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// MonitoringSettings Real user monitoring settings
type MonitoringSettings struct {
	FetchRequests                    bool                           `json:"fetchRequests"`                              // `fetch()` request capture enabled/disabled
	XmlHttpRequest                   bool                           `json:"xmlHttpRequest"`                             // `XmlHttpRequest` support enabled/disabled
	JavaScriptFrameworkSupport       *JavaScriptFrameworkSupport    `json:"javaScriptFrameworkSupport"`                 // Support of various JavaScript frameworks
	ContentCapture                   *ContentCapture                `json:"contentCapture"`                             // Settings for content capture
	ExcludeXHRRegex                  string                         `json:"excludeXhrRegex"`                            // You can exclude some actions from becoming XHR actions.\n\nPut a regular expression, matching all the required URLs, here.\n\nIf noting specified the feature is disabled
	CorrelationHeaderInclusionRegex  string                         `json:"correlationHeaderInclusionRegex"`            // To enable RUM for XHR calls to AWS Lambda, define a regular expression matching these calls, Dynatrace can then automatically add a custom header (x-dtc) to each such request to the respective endpoints in AWS.\n\nImportant: These endpoints must accept the x-dtc header, or the requests will fail
	InjectionMode                    InjectionMode                  `json:"injectionMode"`                              // Possible valures are `CODE_SNIPPET`, `CODE_SNIPPET_ASYNC`, `INLINE_CODE` and `JAVASCRIPT_TAG`
	AddCrossOriginAnonymousAttribute *bool                          `json:"addCrossOriginAnonymousAttribute,omitempty"` // Add the cross origin = anonymous attribute to capture JavaScript error messages and W3C resource timings
	ScriptTagCacheDurationInHours    *int32                         `json:"scriptTagCacheDurationInHours,omitempty"`    // Time duration for the cache settings
	LibraryFileLocation              *string                        `json:"libraryFileLocation,omitempty"`              // The location of your application’s custom JavaScript library file. \n\n If nothing specified the root directory of your web server is used. \n\n **Required** for auto-injected applications, not supported by agentless applications. Maximum 512 characters.
	MonitoringDataPath               string                         `json:"monitoringDataPath"`                         // The location to send monitoring data from the JavaScript tag.\n\n Specify either a relative or an absolute URL. If you use an absolute URL, data will be sent using CORS. \n\n **Required** for auto-injected applications, optional for agentless applications. Maximum 512 characters.
	CustomConfigurationProperties    string                         `json:"customConfigurationProperties"`              // Additional JavaScript tag properties that are specific to your application. To do this, type key=value pairs separated using a (|) symbol. Maximum 1000 characters.
	ServerRequestPathID              string                         `json:"serverRequestPathId"`                        // Path to identify the server’s request ID. Maximum 150 characters.
	SecureCookieAttribute            bool                           `json:"secureCookieAttribute"`                      // Secure attribute usage for Dynatrace cookies enabled/disabled
	CookiePlacementDomain            string                         `json:"cookiePlacementDomain"`                      // Domain for cookie placement. Maximum 150 characters.
	CacheControlHeaderOptimizations  bool                           `json:"cacheControlHeaderOptimizations"`            // Optimize the value of cache control headers for use with Dynatrace real user monitoring enabled/disabled
	AdvancedJavaScriptTagSettings    *AdvancedJavaScriptTagSettings `json:"advancedJavaScriptTagSettings"`              // Advanced JavaScript tag settings
	BrowserRestrictionSettings       *browser.RestrictionSettings   `json:"browserRestrictionSettings,omitempty"`       // Settings for restricting certain browser type, version, platform and, comparator. It also restricts the mode
	IPAddressRestrictionSettings     *ipaddress.RestrictionSettings `json:"ipAddressRestrictionSettings,omitempty"`     // Settings for restricting certain ip addresses and for introducing subnet mask. It also restricts the mode
	JavaScriptInjectionRules         JavaScriptInjectionRules       `json:"javaScriptInjectionRules,omitempty"`         // Java script injection rules
	AngularPackageName               *string                        `json:"angularPackageName,omitempty"`               // The name of the angular package
}

func (me *MonitoringSettings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"fetch_requests": {
			Type:        schema.TypeBool,
			Description: "`fetch()` request capture enabled/disabled",
			Optional:    true,
		},
		"xml_http_request": {
			Type:        schema.TypeBool,
			Description: "`XmlHttpRequest` support enabled/disabled",
			Optional:    true,
		},
		"javascript_framework_support": {
			Type:        schema.TypeList,
			Description: "Support of various JavaScript frameworks",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(JavaScriptFrameworkSupport).Schema()},
		},
		"content_capture": {
			Type:        schema.TypeList,
			Description: "Settings for content capture",
			Required:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(ContentCapture).Schema()},
		},
		"exclude_xhr_regex": {
			Type:        schema.TypeString,
			Description: "You can exclude some actions from becoming XHR actions.\n\nPut a regular expression, matching all the required URLs, here.\n\nIf noting specified the feature is disabled",
			Optional:    true,
		},
		"correlation_header_inclusion_regex": {
			Type:        schema.TypeString,
			Description: "To enable RUM for XHR calls to AWS Lambda, define a regular expression matching these calls, Dynatrace can then automatically add a custom header (`x-dtc`) to each such request to the respective endpoints in AWS.\n\nImportant: These endpoints must accept the `x-dtc` header, or the requests will fail",
			Optional:    true,
		},
		"injection_mode": {
			Type:        schema.TypeString,
			Description: "Possible valures are `CODE_SNIPPET`, `CODE_SNIPPET_ASYNC`, `INLINE_CODE` and `JAVASCRIPT_TAG`.",
			Required:    true,
		},
		"add_cross_origin_anonymous_attribute": {
			Type:        schema.TypeBool,
			Description: "Add the cross origin = anonymous attribute to capture JavaScript error messages and W3C resource timings",
			Optional:    true,
		},
		"script_tag_cache_duration_in_hours": {
			Type:        schema.TypeInt,
			Description: "Time duration for the cache settings",
			Optional:    true,
		},
		"library_file_location": {
			Type:        schema.TypeString,
			Description: "The location of your application’s custom JavaScript library file. \n\n If nothing specified the root directory of your web server is used. \n\n **Required** for auto-injected applications, not supported by agentless applications. Maximum 512 characters.",
			Optional:    true,
		},
		"monitoring_data_path": {
			Type:        schema.TypeString,
			Description: "The location to send monitoring data from the JavaScript tag.\n\n Specify either a relative or an absolute URL. If you use an absolute URL, data will be sent using CORS. \n\n **Required** for auto-injected applications, optional for agentless applications. Maximum 512 characters.",
			Optional:    true,
		},
		"custom_configuration_properties": {
			Type:        schema.TypeString,
			Description: "The location to send monitoring data from the JavaScript tag.\n\n Specify either a relative or an absolute URL. If you use an absolute URL, data will be sent using CORS. \n\n **Required** for auto-injected applications, optional for agentless applications. Maximum 512 characters.",
			Optional:    true,
		},
		"server_request_path_id": {
			Type:        schema.TypeString,
			Description: "Path to identify the server’s request ID. Maximum 150 characters.",
			Optional:    true,
		},
		"secure_cookie_attribute": {
			Type:        schema.TypeBool,
			Description: "Secure attribute usage for Dynatrace cookies enabled/disabled",
			Optional:    true,
		},
		"cookie_placement_domain": {
			Type:        schema.TypeString,
			Description: "Domain for cookie placement. Maximum 150 characters.",
			Optional:    true,
		},
		"cache_control_header_optimizations": {
			Type:        schema.TypeBool,
			Description: "Optimize the value of cache control headers for use with Dynatrace real user monitoring enabled/disabled",
			Optional:    true,
		},
		"advanced_javascript_tag_settings": {
			Type:        schema.TypeList,
			Description: "Advanced JavaScript tag settings",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(AdvancedJavaScriptTagSettings).Schema()},
		},
		"browser_restriction_settings": {
			Type:        schema.TypeList,
			Description: "Settings for restricting certain browser type, version, platform and, comparator. It also restricts the mode",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(browser.RestrictionSettings).Schema()},
		},
		"ip_address_restriction_settings": {
			Type:        schema.TypeList,
			Description: "Settings for restricting certain ip addresses and for introducing subnet mask. It also restricts the mode",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(ipaddress.RestrictionSettings).Schema()},
		},
		"javascript_injection_rules": {
			Type:        schema.TypeList,
			Description: "Java script injection rules",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(JavaScriptInjectionRules).Schema()},
		},
		"angular_package_name": {
			Type:        schema.TypeString,
			Description: "The name of the angular package",
			Optional:    true,
		},
	}
}

func (me *MonitoringSettings) MarshalHCL(properties hcl.Properties) error {
	if err := properties.EncodeAll(map[string]any{
		"fetch_requests":                       me.FetchRequests,
		"xml_http_request":                     me.XmlHttpRequest,
		"javascript_framework_support":         me.JavaScriptFrameworkSupport,
		"content_capture":                      me.ContentCapture,
		"exclude_xhr_regex":                    me.ExcludeXHRRegex,
		"correlation_header_inclusion_regex":   me.CorrelationHeaderInclusionRegex,
		"injection_mode":                       me.InjectionMode,
		"add_cross_origin_anonymous_attribute": me.AddCrossOriginAnonymousAttribute,
		"script_tag_cache_duration_in_hours":   me.ScriptTagCacheDurationInHours,
		"library_file_location":                me.LibraryFileLocation,
		"monitoring_data_path":                 me.MonitoringDataPath,
		"custom_configuration_properties":      me.CustomConfigurationProperties,
		"server_request_path_id":               me.ServerRequestPathID,
		"secure_cookie_attribute":              me.SecureCookieAttribute,
		"cookie_placement_domain":              me.CookiePlacementDomain,
		"cache_control_header_optimizations":   me.CacheControlHeaderOptimizations,
		"advanced_javascript_tag_settings":     me.AdvancedJavaScriptTagSettings,
		"browser_restriction_settings":         me.BrowserRestrictionSettings,
		"ip_address_restriction_settings":      me.IPAddressRestrictionSettings,
		"javascript_injection_rules":           me.JavaScriptInjectionRules,
		"angular_package_name":                 me.AngularPackageName,
	}); err != nil {
		return err
	}
	if me.BrowserRestrictionSettings != nil {
		if len(me.BrowserRestrictionSettings.BrowserRestrictions) == 0 {
			delete(properties, "browser_restriction_settings")
		}
	}
	if me.IPAddressRestrictionSettings != nil {
		if len(me.IPAddressRestrictionSettings.Restrictions) == 0 {
			delete(properties, "ip_address_restriction_settings")
		}
	}
	return nil
}

func (me *MonitoringSettings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"fetch_requests":                       &me.FetchRequests,
		"xml_http_request":                     &me.XmlHttpRequest,
		"javascript_framework_support":         &me.JavaScriptFrameworkSupport,
		"content_capture":                      &me.ContentCapture,
		"exclude_xhr_regex":                    &me.ExcludeXHRRegex,
		"correlation_header_inclusion_regex":   &me.CorrelationHeaderInclusionRegex,
		"injection_mode":                       &me.InjectionMode,
		"add_cross_origin_anonymous_attribute": &me.AddCrossOriginAnonymousAttribute,
		"script_tag_cache_duration_in_hours":   &me.ScriptTagCacheDurationInHours,
		"library_file_location":                &me.LibraryFileLocation,
		"monitoring_data_path":                 &me.MonitoringDataPath,
		"custom_configuration_properties":      &me.CustomConfigurationProperties,
		"server_request_path_id":               &me.ServerRequestPathID,
		"secure_cookie_attribute":              &me.SecureCookieAttribute,
		"cookie_placement_domain":              &me.CookiePlacementDomain,
		"cache_control_header_optimizations":   &me.CacheControlHeaderOptimizations,
		"advanced_javascript_tag_settings":     &me.AdvancedJavaScriptTagSettings,
		"browser_restriction_settings":         &me.BrowserRestrictionSettings,
		"ip_address_restriction_settings":      &me.IPAddressRestrictionSettings,
		"javascript_injection_rules":           &me.JavaScriptInjectionRules,
		"angular_package_name":                 &me.AngularPackageName,
	})
}
