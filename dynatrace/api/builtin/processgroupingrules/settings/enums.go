/**
* @license
* Copyright 2026 Dynatrace LLC
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

package processgroupingrules

type IdType string

var IdTypes = struct {
	Custom   IdType
	Existing IdType
}{
	"CUSTOM",
	"EXISTING",
}

type ProcessType string

var ProcessTypes = struct {
	ProcessTypeApacheHttpd   ProcessType
	ProcessTypeGlassfish     ProcessType
	ProcessTypeGo            ProcessType
	ProcessTypeIbmCicsRegion ProcessType
	ProcessTypeIbmImsControl ProcessType
	ProcessTypeIbmImsMpr     ProcessType
	ProcessTypeIisAppPool    ProcessType
	ProcessTypeJava          ProcessType
	ProcessTypeJboss         ProcessType
	ProcessTypeNginx         ProcessType
	ProcessTypeNodeJs        ProcessType
	ProcessTypePhp           ProcessType
	ProcessTypeRuby          ProcessType
	ProcessTypeTomcat        ProcessType
	ProcessTypeWeblogic      ProcessType
	ProcessTypeWebsphere     ProcessType
}{
	"PROCESS_TYPE_APACHE_HTTPD",
	"PROCESS_TYPE_GLASSFISH",
	"PROCESS_TYPE_GO",
	"PROCESS_TYPE_IBM_CICS_REGION",
	"PROCESS_TYPE_IBM_IMS_CONTROL",
	"PROCESS_TYPE_IBM_IMS_MPR",
	"PROCESS_TYPE_IIS_APP_POOL",
	"PROCESS_TYPE_JAVA",
	"PROCESS_TYPE_JBOSS",
	"PROCESS_TYPE_NGINX",
	"PROCESS_TYPE_NODE_JS",
	"PROCESS_TYPE_PHP",
	"PROCESS_TYPE_RUBY",
	"PROCESS_TYPE_TOMCAT",
	"PROCESS_TYPE_WEBLOGIC",
	"PROCESS_TYPE_WEBSPHERE",
}

type ReportItem string

var ReportItems = struct {
	Always ReportItem
	Auto   ReportItem
	Never  ReportItem
}{
	"always",
	"auto",
	"never",
}
