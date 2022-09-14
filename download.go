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

package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/dtcookie/dynatrace/api/cluster/v2/envs"
	"github.com/dtcookie/dynatrace/api/config/alerting"
	"github.com/dtcookie/dynatrace/api/config/anomalies/applications"
	"github.com/dtcookie/dynatrace/api/config/anomalies/databaseservices"
	"github.com/dtcookie/dynatrace/api/config/anomalies/diskevents"
	"github.com/dtcookie/dynatrace/api/config/anomalies/hosts"
	"github.com/dtcookie/dynatrace/api/config/anomalies/metricevents"
	"github.com/dtcookie/dynatrace/api/config/anomalies/services"
	"github.com/dtcookie/dynatrace/api/config/applications/mobile"
	"github.com/dtcookie/dynatrace/api/config/applications/web"
	"github.com/dtcookie/dynatrace/api/config/autotags"
	"github.com/dtcookie/dynatrace/api/config/credentials/aws"
	"github.com/dtcookie/dynatrace/api/config/credentials/azure"
	"github.com/dtcookie/dynatrace/api/config/credentials/cloudfoundry"
	"github.com/dtcookie/dynatrace/api/config/credentials/kubernetes"
	"github.com/dtcookie/dynatrace/api/config/customservices"
	"github.com/dtcookie/dynatrace/api/config/dashboards"
	"github.com/dtcookie/dynatrace/api/config/dashboards/sharing"
	"github.com/dtcookie/dynatrace/api/config/maintenance"
	"github.com/dtcookie/dynatrace/api/config/managementzones"
	"github.com/dtcookie/dynatrace/api/config/metrics/calculated/service"
	hostnaming "github.com/dtcookie/dynatrace/api/config/naming/hosts"
	processgroupnaming "github.com/dtcookie/dynatrace/api/config/naming/processgroups"
	servicenaming "github.com/dtcookie/dynatrace/api/config/naming/services"
	"github.com/dtcookie/dynatrace/api/config/notifications"
	"github.com/dtcookie/dynatrace/api/config/requestattributes"
	"github.com/dtcookie/dynatrace/api/config/requestnaming"
	"github.com/dtcookie/dynatrace/api/config/synthetic/monitors"
	servicetopology "github.com/dtcookie/dynatrace/api/config/topology/service"
	"github.com/dtcookie/dynatrace/api/config/v2/ibmmq/filters"
	"github.com/dtcookie/dynatrace/api/config/v2/ibmmq/queuemanagers"
	"github.com/dtcookie/dynatrace/api/config/v2/ibmmq/queuesharinggroups"
	"github.com/dtcookie/dynatrace/api/config/v2/keyrequests"
	"github.com/dtcookie/dynatrace/api/config/v2/slo"
	"github.com/dtcookie/dynatrace/api/config/v2/spans/attributes"
	"github.com/dtcookie/dynatrace/api/config/v2/spans/capture"
	"github.com/dtcookie/dynatrace/api/config/v2/spans/ctxprop"
	"github.com/dtcookie/dynatrace/api/config/v2/spans/entrypoints"
	"github.com/dtcookie/dynatrace/api/config/v2/spans/resattr"
	"github.com/dtcookie/dynatrace/rest"
	"github.com/dtcookie/hcl"
	"github.com/google/uuid"
)

func escape(s string) string {
	result := ""
	for _, c := range s {
		if unicode.IsLetter(c) {
			result = result + string(c)
		} else if unicode.IsDigit(c) {
			result = result + string(c)
		} else if c == '-' {
			result = result + string(c)
		} else if c == '_' {
			result = result + string(c)
		} else {
			result = result + "_"
		}
	}
	return result
}

/*
  < (less than)
  > (greater than)
  : (colon - sometimes works, but is actually NTFS Alternate Data Streams)
  " (double quote)
  / (forward slash)
  \ (backslash)
  | (vertical bar or pipe)
  ? (question mark)
  * (asterisk)
*/

var forbiddenFileNameChars = []string{"<", ">", ":", "\"", "/", "|", "?", "*", "	", "\r", "\n", "\f", "\v"}

func escf(s string) string {
	for _, ch := range forbiddenFileNameChars {
		s = strings.ReplaceAll(s, ch, "_")
	}
	return s
}

func escFileName(s string, id string) string {
	if !IsValidUUID(id) && !IsValidID(id) {
		return escf(s)
	}
	return escf(s + "." + id + ".")
}

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func IsValidID(uuid string) bool {
	if uuid == "" {
		return false
	}
	r := regexp.MustCompile("^[A-Z0-9]*$")
	res := r.MatchString(uuid)
	if res {
		return true
	}
	r = regexp.MustCompile(`^[A-Z]*-[A-Z0-9]*$`)
	res = r.MatchString(uuid)
	if res {
		return true
	}
	r = regexp.MustCompile(`^[A-Z]*_[A-Z]*-[A-Z0-9]*$`)
	res = r.MatchString(uuid)
	if res {
		return true
	}
	// if !res {
	// 	fmt.Printf("%v is not a valid ID\n", uuid)
	// }
	return res
}

func ctns(elems []string, elem string) bool {
	if elems == nil {
		return true
	}
	for _, el := range elems {
		if el == elem {
			return true
		}
	}
	return false
}

func importEnvironments(targetFolder string, clusterURL string, apiToken string, argids []string) error {

	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := envs.NewService(clusterURL+"/api/cluster/v2", apiToken)

	environmentList, err := restClient.ListAll()
	if err != nil {
		return err
	}
	for _, environment := range environmentList.Environments {
		if !ctns(argids, *environment.ID) {
			continue
		}
		config, err := restClient.Get(*environment.ID)
		if err != nil {
			return err
		}
		var file *os.File
		fileName := targetFolder + "/" + escFileName(config.Name, *environment.ID) + ".environment.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_environment", escape(config.Name))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importBrowserMonitors(targetFolder string, environmentURL string, apiToken string, argids []string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := monitors.NewService(environmentURL+"/api/v1", apiToken)

	stubList, err := restClient.ListBrowser()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Monitors {
		if !ctns(argids, stub.EntityID) {
			continue
		}
		config, err := restClient.GetBrowser(stub.EntityID)
		if err != nil {
			return err
		}
		var file *os.File
		fileName := targetFolder + "/" + escFileName(config.Name, stub.EntityID) + ".monitor.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_browser_monitor", escape(config.Name))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importHTTPMonitors(targetFolder string, environmentURL string, apiToken string, argids []string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := monitors.NewService(environmentURL+"/api/v1", apiToken)

	stubList, err := restClient.ListHTTP()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Monitors {
		if !ctns(argids, stub.EntityID) {
			continue
		}
		config, err := restClient.GetHTTP(stub.EntityID)
		if err != nil {
			return err
		}
		var file *os.File
		fileName := targetFolder + "/" + escFileName(config.Name, stub.EntityID) + ".monitor.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_http_monitor", escape(config.Name))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

// func importVaultCredentials(targetFolder string, environmentURL string, apiToken string, argids []string) error {

// 	os.MkdirAll(targetFolder, os.ModePerm)
// 	restClient := vault.NewService(environmentURL+"/api/config/v1", apiToken)

// 	stubList, err := restClient.ListAll()
// 	if err != nil {
// 		return err
// 	}
// 	for _, stub := range stubList.Credentials {
// 		if !ctns(argids, *stub.ID) {
// 			continue
// 		}
// 		config, err := restClient.Get(*stub.ID)
// 		if err != nil {
// 			return err
// 		}
// 		var file *os.File
// 		fileName := targetFolder + "/" + escFileName(config.Name, *stub.ID) + ".credentials.vault.tf"
// 		os.Remove(fileName)
// 		if file, err = os.Create(fileName); err != nil {
// 			return err
// 		}
// 		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_credentials", escape(config.Name))); err != nil {
// 			file.Close()
// 			return err
// 		}
// 		if config.GetType() == vault.CredentialsTypes.UsernamePassword {
// 			config.Username = opt.NewString("--redacted--")
// 			config.Password = opt.NewString("--redacted--")
// 		} else if (config.GetType() == vault.CredentialsTypes.Certificate) || (config.GetType() == vault.CredentialsTypes.PublicCertificate) {
// 			config.Certificate = opt.NewString("--redacted--")
// 			config.CertificateFormat = vault.CertificateFormat("--redacted--").Ref()
// 		} else if config.GetType() == vault.CredentialsTypes.Token {
// 			config.Token = opt.NewString("--redacted--")
// 		}
// 		if err := hcl.ExportOpt(config, file); err != nil {
// 			file.Close()
// 			return err
// 		}
// 		if _, err := file.WriteString("}\n"); err != nil {
// 			file.Close()
// 			return err
// 		}
// 		file.Close()
// 	}
// 	return nil
// }

func importAWSCredentials(targetFolder string, environmentURL string, apiToken string, argids []string) error {

	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := aws.NewService(environmentURL+"/api/config/v1", apiToken)

	stubList, err := restClient.ListAll()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Values {
		if !ctns(argids, stub.ID) {
			continue
		}
		config, err := restClient.Get(stub.ID)
		if err != nil {
			return err
		}
		config.Metadata = nil
		var file *os.File
		fileName := targetFolder + "/" + escFileName(config.Label, stub.ID) + ".credentials.aws.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_aws_credentials", escape(config.Label))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importAzureCredentials(targetFolder string, environmentURL string, apiToken string, argids []string) error {

	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := azure.NewService(environmentURL+"/api/config/v1", apiToken)

	stubList, err := restClient.ListAll()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Values {
		if !ctns(argids, stub.ID) {
			continue
		}
		config, err := restClient.Get(stub.ID)
		if err != nil {
			return err
		}
		config.Metadata = nil
		var file *os.File
		fileName := targetFolder + "/" + escFileName(config.Label, stub.ID) + ".credentials.azure.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_azure_credentials", escape(config.Label))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importCloudFoundryCredentials(targetFolder string, environmentURL string, apiToken string, argids []string) error {

	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := cloudfoundry.NewService(environmentURL+"/api/config/v1", apiToken)

	if os.Getenv("DYNATRACE_DEBUG") == "true" {
		rest.Verbose = true
	}

	stubList, err := restClient.ListAll()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Values {
		if !ctns(argids, stub.ID) {
			continue
		}
		config, err := restClient.Get(stub.ID)
		if err != nil {
			return err
		}
		var file *os.File
		fileName := targetFolder + "/" + escFileName(config.Name, stub.ID) + ".credentials.cloudfoundry.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_cloudfoundry_credentials", escape(config.Name))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importK8sCredentials(targetFolder string, environmentURL string, apiToken string, argids []string) error {

	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := kubernetes.NewService(environmentURL+"/api/config/v1", apiToken)

	stubList, err := restClient.ListAll()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Values {
		if !ctns(argids, stub.ID) {
			continue
		}
		config, err := restClient.Get(stub.ID)
		if err != nil {
			return err
		}
		config.Metadata = nil
		var file *os.File
		fileName := targetFolder + "/" + escFileName(config.Label, stub.ID) + ".credentials.k8s.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_k8s_credentials", escape(config.Label))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importNotificationConfigs(targetFolder string, environmentURL string, apiToken string, argids []string) error {

	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := notifications.NewService(environmentURL+"/api/config/v1", apiToken)

	stubList, err := restClient.ListAll()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Values {
		if !ctns(argids, stub.ID) {
			continue
		}
		config, err := restClient.Get(stub.ID)
		if err != nil {
			return err
		}
		var file *os.File
		fileName := targetFolder + "/" + escFileName(config.NotificationConfig.GetName(), stub.ID) + ".notification.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_notification", escape(config.NotificationConfig.GetName()))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExtExport(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importManagementZones(targetFolder string, environmentURL string, apiToken string, argids []string) error {

	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := managementzones.NewService(environmentURL+"/api/config/v1", apiToken)

	stubList, err := restClient.ListAll()
	if err != nil {
		return err
	}
	for _, stub := range stubList {
		if !ctns(argids, stub.ID) {
			continue
		}
		config, err := restClient.Get(stub.ID, false)
		if err != nil {
			return err
		}
		config.Metadata = nil
		var file *os.File
		fileName := targetFolder + "/" + escFileName(config.Name, stub.ID) + ".management_zone.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_management_zone", escape(config.Name))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importAlertingProfiles(targetFolder string, environmentURL string, apiToken string, argids []string) error {

	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := alerting.NewService(environmentURL+"/api/config/v1", apiToken)

	stubList, err := restClient.List()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Values {
		if !ctns(argids, stub.ID) {
			continue
		}
		config, err := restClient.Get(stub.ID)
		if err != nil {
			return err
		}
		config.Metadata = nil
		var file *os.File
		fileName := targetFolder + "/" + escFileName(config.DisplayName, stub.ID) + ".alerting_profile.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_alerting_profile", escape(config.DisplayName))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
	}
	return nil
}

func importAutoTags(targetFolder string, environmentURL string, apiToken string, argids []string) error {

	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := autotags.NewService(environmentURL+"/api/config/v1", apiToken)

	stubList, err := restClient.ListAll()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Values {
		if !ctns(argids, stub.ID) {
			continue
		}
		config, err := restClient.Get(stub.ID)
		if err != nil {
			return err
		}
		config.Metadata = nil
		var file *os.File
		fileName := targetFolder + "/" + escFileName(config.Name, stub.ID) + ".autotag.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_autotag", escape(config.Name))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importMaintenance(targetFolder string, environmentURL string, apiToken string, argids []string) error {

	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := maintenance.NewService(environmentURL+"/api/config/v1", apiToken)

	stubList, err := restClient.ListAll()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Values {
		if !ctns(argids, stub.ID) {
			continue
		}
		config, err := restClient.Get(stub.ID)
		if err != nil {
			return err
		}
		config.Metadata = nil
		var file *os.File
		fileName := targetFolder + "/" + escFileName(config.Name, stub.ID) + ".maintenance.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_maintenance_window", escape(config.Name))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importRequestAttributes(targetFolder string, environmentURL string, apiToken string, argids []string) error {

	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := requestattributes.NewService(environmentURL+"/api/config/v1", apiToken)

	stubList, err := restClient.ListAll()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Values {
		if !ctns(argids, stub.ID) {
			continue
		}
		config, err := restClient.Get(stub.ID)
		if err != nil {
			return err
		}
		config.Metadata = nil
		var file *os.File
		fileName := targetFolder + "/" + escFileName(config.Name, stub.ID) + ".request_attribute.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_request_attribute", escape(config.Name))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importDashboards(targetFolder string, environmentURL string, apiToken string, argids []string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := dashboards.NewService(environmentURL+"/api/config/v1", apiToken)

	stubList, err := restClient.ListAll()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Dashboards {
		if !ctns(argids, stub.ID) {
			continue
		}
		config, err := restClient.Get(stub.ID)
		if err != nil {
			return err
		}
		config.ConfigurationMetadata = nil
		var file *os.File
		fileName := targetFolder + "/" + escFileName(config.Metadata.Name, stub.ID) + ".dashboard.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_dashboard", escape(config.Metadata.Name+"_"+stub.ID))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		shareRestClient := sharing.NewService(environmentURL+"/api/config/v1", apiToken)
		var shareSettings *sharing.DashboardSharing
		if shareSettings, err = shareRestClient.Get(stub.ID); err != nil {
			if !strings.Contains(err.Error(), "Editing or deleting a non user specific dashboard preset is not allowed") {
				file.Close()
				return err
			}
		}
		if shareSettings != nil {
			if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_dashboard_sharing", escape(config.Metadata.Name+"_"+stub.ID))); err != nil {
				file.Close()
				return err
			}
			if err := hcl.ExportOpt(shareSettings, file); err != nil {
				file.Close()
				return err
			}
			if _, err := file.WriteString("}\n"); err != nil {
				file.Close()
				return err
			}
		}
		file.Close()
	}
	return nil
}

func importCustomServices(targetFolder string, environmentURL string, apiToken string, argids []string) error {
	if err := importCustomServicesTech(targetFolder, environmentURL, apiToken, customservices.Technologies.Java, argids); err != nil {
		return err
	}
	if err := importCustomServicesTech(targetFolder, environmentURL, apiToken, customservices.Technologies.DotNet, argids); err != nil {
		return err
	}
	if err := importCustomServicesTech(targetFolder, environmentURL, apiToken, customservices.Technologies.Go, argids); err != nil {
		return err
	}
	if err := importCustomServicesTech(targetFolder, environmentURL, apiToken, customservices.Technologies.NodeJS, argids); err != nil {
		return err
	}
	if err := importCustomServicesTech(targetFolder, environmentURL, apiToken, customservices.Technologies.PHP, argids); err != nil {
		return err
	}
	return nil
}

func importCustomServicesTech(targetFolder string, environmentURL string, apiToken string, technology customservices.Technology, argids []string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := customservices.NewService(environmentURL+"/api/config/v1", apiToken)

	stubList, err := restClient.List(technology)
	if err != nil {
		return err
	}
	for _, stub := range stubList.Values {
		if !ctns(argids, stub.ID) {
			continue
		}
		config, err := restClient.Get(stub.ID, technology, false)
		if err != nil {
			return err
		}
		config.Metadata = nil
		var file *os.File
		fileName := targetFolder + "/" + escFileName(config.Name, stub.ID) + ".custom_service.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_custom_service", escape(config.Name))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importDiskAnomalies(targetFolder string, environmentURL string, apiToken string, argids []string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := diskevents.NewService(environmentURL+"/api/config/v1", apiToken)

	stubList, err := restClient.List()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Values {
		if !ctns(argids, stub.ID) {
			continue
		}
		config, err := restClient.Get(stub.ID)
		if err != nil {
			return err
		}
		var file *os.File
		fileName := targetFolder + "/" + escFileName(config.Name, stub.ID) + ".disk_anomalies.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_disk_anomalies", escape(config.Name))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExtExport(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importMetricAnomalies(targetFolder string, environmentURL string, apiToken string, argids []string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := metricevents.NewService(environmentURL+"/api/config/v1", apiToken)

	stubList, err := restClient.List()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Values {
		if !ctns(argids, stub.ID) {
			continue
		}
		config, err := restClient.Get(stub.ID)
		if err != nil {
			return err
		}
		var file *os.File
		name := config.Name
		if name == "" {
			name = uuid.New().String()
		}
		fileName := targetFolder + "/" + escFileName(config.Name, stub.ID) + ".custom_anomalies.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_custom_anomalies", escape(name))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExtExport(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importDatabaseAnomalies(targetFolder string, environmentURL string, apiToken string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := databaseservices.NewService(environmentURL+"/api/config/v1", apiToken)

	config, err := restClient.Get()
	if err != nil {
		return err
	}
	var file *os.File
	fileName := targetFolder + "/" + "database_anomalies.tf"
	os.Remove(fileName)
	if file, err = os.Create(fileName); err != nil {
		return err
	}
	if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_database_anomalies", "dynatrace_database_anomalies")); err != nil {
		file.Close()
		return err
	}
	if err := hcl.ExtExport(config, file); err != nil {
		file.Close()
		return err
	}
	if _, err := file.WriteString("}\n"); err != nil {
		file.Close()
		return err
	}

	return nil
}

func importHostAnomalies(targetFolder string, environmentURL string, apiToken string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := hosts.NewService(environmentURL+"/api/config/v1", apiToken)

	config, err := restClient.Get()
	if err != nil {
		return err
	}
	var file *os.File
	fileName := targetFolder + "/" + "host_anomalies.tf"
	os.Remove(fileName)
	if file, err = os.Create(fileName); err != nil {
		return err
	}
	if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_host_anomalies", "dynatrace_host_anomalies")); err != nil {
		file.Close()
		return err
	}
	if err := hcl.ExtExport(config, file); err != nil {
		file.Close()
		return err
	}
	if _, err := file.WriteString("}\n"); err != nil {
		file.Close()
		return err
	}
	file.Close()

	return nil
}

func importApplicationAnomalies(targetFolder string, environmentURL string, apiToken string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := applications.NewService(environmentURL+"/api/config/v1", apiToken)

	config, err := restClient.Get()
	if err != nil {
		return err
	}
	var file *os.File
	fileName := targetFolder + "/" + "application_anomalies.tf"
	os.Remove(fileName)
	if file, err = os.Create(fileName); err != nil {
		return err
	}
	if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_application_anomalies", "dynatrace_application_anomalies")); err != nil {
		file.Close()
		return err
	}
	if err := hcl.ExtExport(config, file); err != nil {
		file.Close()
		return err
	}
	if _, err := file.WriteString("}\n"); err != nil {
		file.Close()
		return err
	}
	file.Close()

	return nil
}

func importServiceAnomalies(targetFolder string, environmentURL string, apiToken string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := services.NewService(environmentURL+"/api/config/v1", apiToken)

	config, err := restClient.Get()
	if err != nil {
		return err
	}
	var file *os.File
	fileName := targetFolder + "/" + "service_anomalies.tf"
	os.Remove(fileName)
	if file, err = os.Create(fileName); err != nil {
		return err
	}
	if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_service_anomalies", "dynatrace_service_anomalies")); err != nil {
		file.Close()
		return err
	}
	if err := hcl.ExtExport(config, file); err != nil {
		file.Close()
		return err
	}
	if _, err := file.WriteString("}\n"); err != nil {
		file.Close()
		return err
	}
	file.Close()

	return nil
}

func importCalculatedServiceMetrics(targetFolder string, environmentURL string, apiToken string, argids []string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := service.NewService(environmentURL+"/api/config/v1", apiToken)

	stubList, err := restClient.ListAll()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Values {
		if !ctns(argids, stub.ID) {
			continue
		}
		config, err := restClient.Get(stub.ID)
		if err != nil {
			return err
		}
		var file *os.File
		name := config.Name
		if name == "" {
			name = uuid.New().String()
		}
		fileName := targetFolder + "/" + escFileName(config.Name, stub.ID) + ".calculated_service_metric.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_calculated_service_metric", escape(name))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importMobileApps(targetFolder string, environmentURL string, apiToken string, argids []string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := mobile.NewService(environmentURL+"/api/config/v1", apiToken)

	stubList, err := restClient.List()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Values {
		if !ctns(argids, stub.ID) {
			continue
		}
		config, err := restClient.Get(stub.ID)
		if err != nil {
			return err
		}
		var file *os.File
		name := config.Name
		if name == "" {
			name = uuid.New().String()
		}
		fileName := targetFolder + "/" + escFileName(config.Name, stub.ID) + ".mobile_application.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_mobile_application", escape(name))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importKeyRequests(targetFolder string, environmentURL string, apiToken string, argids []string) error {
	os.MkdirAll(targetFolder, os.ModePerm)

	topRestClient := servicetopology.NewService(environmentURL+"/api/v1", apiToken)
	services, err := topRestClient.List()
	if err != nil {
		return err
	}
	restClient := keyrequests.NewService(environmentURL+"/api/v2", apiToken)

	for _, service := range services {
		keyRequestID, keyRequest, err := restClient.List(service.EntityId)
		if err != nil {
			return err
		}
		if keyRequest == nil {
			continue
		}
		if !ctns(argids, keyRequestID) {
			continue
		}
		var file *os.File
		name := "for " + service.DisplayName
		if name == "" {
			name = uuid.New().String()
		}
		fileName := targetFolder + "/" + escFileName(service.DisplayName, service.EntityId) + ".key_requests.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_key_requests", escape(name))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(keyRequest, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importRequestNamings(targetFolder string, environmentURL string, apiToken string, argids []string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := requestnaming.NewService(environmentURL+"/api/config/v1", apiToken)
	stubList, err := restClient.ListAll()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Values {
		if !ctns(argids, stub.ID) {
			continue
		}
		config, err := restClient.Get(stub.ID)
		if err != nil {
			return err
		}
		var file *os.File
		name := stub.Name
		if name == "" {
			name = uuid.New().String()
		}
		fileName := targetFolder + "/" + escFileName(stub.Name, stub.ID) + ".request_naming.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_request_naming", escape(name))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importApplicationErrorRules(targetFolder string, environmentURL string, apiToken string, argids []string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := web.NewService(environmentURL+"/api/config/v1", apiToken)
	stubList, err := restClient.List()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Values {
		if !ctns(argids, stub.ID) {
			continue
		}
		config, err := restClient.GetErrorRules(stub.ID)
		if err != nil {
			return err
		}
		var file *os.File
		name := stub.Name
		if name == "" {
			name = uuid.New().String()
		}
		fileName := targetFolder + "/" + escFileName(stub.Name, stub.ID) + ".application_error_rules.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_application_error_rules", escape(name))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importApplicationDataPrivacy(targetFolder string, environmentURL string, apiToken string, argids []string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := web.NewService(environmentURL+"/api/config/v1", apiToken)
	stubList, err := restClient.List()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Values {
		if !ctns(argids, stub.ID) {
			continue
		}
		config, err := restClient.GetAppDataPrivacy(stub.ID)
		if err != nil {
			return err
		}
		var file *os.File
		name := stub.Name
		if name == "" {
			name = uuid.New().String()
		}
		fileName := targetFolder + "/" + escFileName(stub.Name, stub.ID) + ".application_data_privacy.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_application_data_privacy", escape(name))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importWebApps(targetFolder string, environmentURL string, apiToken string, argids []string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := web.NewService(environmentURL+"/api/config/v1", apiToken)
	stubList, err := restClient.List()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Values {
		if !ctns(argids, stub.ID) {
			continue
		}
		config, err := restClient.Get(stub.ID)
		if err != nil {
			return err
		}
		var file *os.File
		name := config.Name
		if name == "" {
			name = uuid.New().String()
		}
		fileName := targetFolder + "/" + escFileName(config.Name, stub.ID) + ".web_application.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_web_application", escape(name))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importServiceNamings(targetFolder string, environmentURL string, apiToken string, argids []string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := servicenaming.NewService(environmentURL+"/api/config/v1", apiToken)

	stubList, err := restClient.List()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Values {
		if !ctns(argids, stub.ID) {
			continue
		}
		config, err := restClient.Get(stub.ID)
		if err != nil {
			return err
		}
		var file *os.File
		name := config.Name
		if name == "" {
			name = uuid.New().String()
		}
		fileName := targetFolder + "/" + escFileName(config.Name, stub.ID) + ".service_naming.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_service_naming", escape(name))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importHostNamings(targetFolder string, environmentURL string, apiToken string, argids []string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := hostnaming.NewService(environmentURL+"/api/config/v1", apiToken)

	stubList, err := restClient.List()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Values {
		if !ctns(argids, stub.ID) {
			continue
		}
		config, err := restClient.Get(stub.ID)
		if err != nil {
			return err
		}
		var file *os.File
		name := config.Name
		if name == "" {
			name = uuid.New().String()
		}
		fileName := targetFolder + "/" + escFileName(config.Name, stub.ID) + ".host_naming.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_host_naming", escape(name))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importProcessGroupNamings(targetFolder string, environmentURL string, apiToken string, argids []string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := processgroupnaming.NewService(environmentURL+"/api/config/v1", apiToken)

	stubList, err := restClient.List()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Values {
		if !ctns(argids, stub.ID) {
			continue
		}
		config, err := restClient.Get(stub.ID)
		if err != nil {
			return err
		}
		var file *os.File
		name := config.Name
		if name == "" {
			name = uuid.New().String()
		}
		fileName := targetFolder + "/" + escFileName(config.Name, stub.ID) + ".processgroup_naming.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_processgroup_naming", escape(name))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importSLOs(targetFolder string, environmentURL string, apiToken string, argids []string) error {

	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := slo.NewService(environmentURL+"/api/v2", apiToken)

	ids, err := restClient.List()
	if err != nil {
		return err
	}
	for _, id := range ids {
		if !ctns(argids, id) {
			continue
		}
		config, err := restClient.Get(id)
		if err != nil {
			return err
		}
		var file *os.File
		fileName := targetFolder + "/" + escFileName(config.Name, id) + ".slo.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_slo", escape(config.Name))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importSpanEntryPoints(targetFolder string, environmentURL string, apiToken string, argids []string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := entrypoints.NewService(environmentURL+"/api/v2", apiToken)

	ids, err := restClient.List()
	if err != nil {
		return err
	}
	for _, id := range ids {
		if !ctns(argids, id) {
			continue
		}
		config, err := restClient.Get(id)
		if err != nil {
			return err
		}
		var file *os.File
		fileName := targetFolder + "/" + escFileName(config.EntryPointRule.Name, id) + ".span_entry_point.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_span_entry_point", escape(config.EntryPointRule.Name))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importSpanCaptureRules(targetFolder string, environmentURL string, apiToken string, argids []string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := capture.NewService(environmentURL+"/api/v2", apiToken)

	ids, err := restClient.List()
	if err != nil {
		return err
	}
	for _, id := range ids {
		if !ctns(argids, id) {
			continue
		}
		config, err := restClient.Get(id)
		if err != nil {
			return err
		}
		var file *os.File
		fileName := targetFolder + "/" + escFileName(config.SpanCaptureRule.Name, id) + ".span_capture_rule.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_span_capture_rule", escape(config.SpanCaptureRule.Name))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importSpanContextPropagation(targetFolder string, environmentURL string, apiToken string, argids []string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := ctxprop.NewService(environmentURL+"/api/v2", apiToken)

	ids, err := restClient.List()
	if err != nil {
		return err
	}
	for _, id := range ids {
		if !ctns(argids, id) {
			continue
		}
		config, err := restClient.Get(id)
		if err != nil {
			return err
		}
		var file *os.File
		fileName := targetFolder + "/" + escFileName(config.PropagationRule.Name, id) + ".span_context_propagation.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_span_context_propagation", escape(config.PropagationRule.Name))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importResourceAttributes(targetFolder string, environmentURL string, apiToken string, argids []string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := resattr.NewService(environmentURL+"/api/v2", apiToken)

	ids, err := restClient.List()
	if err != nil {
		return err
	}
	for _, id := range ids {
		if !ctns(argids, id) {
			continue
		}
		config, err := restClient.Get(id)
		if err != nil {
			return err
		}
		var file *os.File
		fileName := targetFolder + "/" + "resource_attributes.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_resource_attributes", "resource_attributes")); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importSpanAttributes(targetFolder string, environmentURL string, apiToken string, argids []string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := attributes.NewService(environmentURL+"/api/v2", apiToken)

	ids, err := restClient.List()
	if err != nil {
		return err
	}
	for _, id := range ids {
		if !ctns(argids, id) {
			continue
		}
		config, err := restClient.Get(id)
		if err != nil {
			return err
		}
		var file *os.File
		fileName := targetFolder + "/" + escFileName(config.Key, "") + ".span_attribute.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_span_attribute", escape(config.Key))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importMQQueueManagers(targetFolder string, environmentURL string, apiToken string, argids []string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := queuemanagers.NewService(environmentURL+"/api/v2", apiToken)

	ids, err := restClient.List()
	if err != nil {
		return err
	}
	for _, id := range ids {
		if !ctns(argids, id) {
			continue
		}
		config, err := restClient.Get(id)
		if err != nil {
			return err
		}
		var file *os.File
		fileName := targetFolder + "/" + escFileName(config.Name, "") + ".queue_manager.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_queue_manager", escape(config.Name))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importMQFilters(targetFolder string, environmentURL string, apiToken string, argids []string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := filters.NewService(environmentURL+"/api/v2", apiToken)

	ids, err := restClient.List()
	if err != nil {
		return err
	}
	for _, id := range ids {
		if !ctns(argids, id) {
			continue
		}
		config, err := restClient.Get(id)
		if err != nil {
			return err
		}
		var file *os.File
		fileName := targetFolder + "/" + "Filters" + ".ibm_mq_filters.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_ibm_mq_filters", "Filters")); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func importMQQueueSharingGroups(targetFolder string, environmentURL string, apiToken string, argids []string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := queuesharinggroups.NewService(environmentURL+"/api/v2", apiToken)

	ids, err := restClient.List()
	if err != nil {
		return err
	}
	for _, id := range ids {
		if !ctns(argids, id) {
			continue
		}
		config, err := restClient.Get(id)
		if err != nil {
			return err
		}
		var file *os.File
		fileName := targetFolder + "/" + escFileName(config.Name, "") + ".queue_sharing_groups.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_queue_sharing_groups", escape(config.Name))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(config, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}
