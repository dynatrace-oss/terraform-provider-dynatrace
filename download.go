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
	"strings"

	"github.com/dtcookie/dynatrace/api/config/alertingprofiles"
	"github.com/dtcookie/dynatrace/api/config/autotags"
	"github.com/dtcookie/dynatrace/api/config/credentials/aws"
	"github.com/dtcookie/dynatrace/api/config/credentials/azure"
	"github.com/dtcookie/dynatrace/api/config/credentials/kubernetes"
	"github.com/dtcookie/dynatrace/api/config/customservices"
	"github.com/dtcookie/dynatrace/api/config/dashboards"
	"github.com/dtcookie/dynatrace/api/config/managementzones"
	"github.com/dtcookie/dynatrace/api/config/notifications"
	"github.com/dtcookie/dynatrace/api/config/requestattributes"
	"github.com/dtcookie/dynatrace/rest"
	"github.com/dtcookie/hcl"
)

func escape(s string) string {
	s = strings.ReplaceAll(s, ":", "_")
	s = strings.ReplaceAll(s, "/", "_")
	return s
}

func importAWSCredentials(targetFolder string, environmentURL string, apiToken string) error {
	rest.Verbose = true
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := aws.NewService(environmentURL+"/api/config/v1", apiToken)
	rest.Verbose = false
	stubList, err := restClient.ListAll()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Values {
		config, err := restClient.Get(stub.ID)
		if err != nil {
			return err
		}
		config.Metadata = nil
		var file *os.File
		fileName := targetFolder + "/" + escape(config.Label) + ".credentials.aws.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		defer file.Close()
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_aws_credentials", config.Label)); err != nil {
			return err
		}
		if err := hcl.Export(config, file); err != nil {
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			return err
		}
	}
	return nil
}

func importAzureCredentials(targetFolder string, environmentURL string, apiToken string) error {
	rest.Verbose = true
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := azure.NewService(environmentURL+"/api/config/v1", apiToken)
	rest.Verbose = false
	stubList, err := restClient.ListAll()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Values {
		config, err := restClient.Get(stub.ID)
		if err != nil {
			return err
		}
		config.Metadata = nil
		var file *os.File
		fileName := targetFolder + "/" + escape(config.Label) + ".credentials.azure.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		defer file.Close()
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_azure_credentials", config.Label)); err != nil {
			return err
		}
		if err := hcl.Export(config, file); err != nil {
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			return err
		}
	}
	return nil
}

func importK8sCredentials(targetFolder string, environmentURL string, apiToken string) error {
	rest.Verbose = true
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := kubernetes.NewService(environmentURL+"/api/config/v1", apiToken)
	rest.Verbose = false
	stubList, err := restClient.ListAll()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Values {
		config, err := restClient.Get(stub.ID)
		if err != nil {
			return err
		}
		config.Metadata = nil
		var file *os.File
		fileName := targetFolder + "/" + escape(config.Label) + ".credentials.k8s.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		defer file.Close()
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_k8s_credentials", config.Label)); err != nil {
			return err
		}
		if err := hcl.Export(config, file); err != nil {
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			return err
		}
	}
	return nil
}

func importNotificationConfigs(targetFolder string, environmentURL string, apiToken string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := notifications.NewService(environmentURL+"/api/config/v1", apiToken)
	rest.Verbose = false
	stubList, err := restClient.ListAll()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Values {
		config, err := restClient.Get(stub.ID)
		if err != nil {
			return err
		}
		exporter := &NotificationExporter{NotificationConfig: config}
		var file *os.File
		fileName := targetFolder + "/" + escape(config.GetName()) + ".notification.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		defer file.Close()
		exporter.ToHCL(file)
	}
	return nil
}

func importManagementZones(targetFolder string, environmentURL string, apiToken string) error {
	rest.Verbose = true
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := managementzones.NewService(environmentURL+"/api/config/v1", apiToken)
	rest.Verbose = false
	stubList, err := restClient.ListAll()
	if err != nil {
		return err
	}
	for _, stub := range stubList {
		config, err := restClient.Get(stub.ID, false)
		if err != nil {
			return err
		}
		config.Metadata = nil
		var file *os.File
		fileName := targetFolder + "/" + escape(config.Name) + ".management_zone.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		defer file.Close()
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_management_zone", config.Name)); err != nil {
			return err
		}
		if err := hcl.Export(config, file); err != nil {
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			return err
		}
	}
	return nil
}

func importAlertingProfiles(targetFolder string, environmentURL string, apiToken string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := alertingprofiles.NewService(environmentURL+"/api/config/v1", apiToken)
	rest.Verbose = false
	stubList, err := restClient.List()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Values {
		alertingProfile, err := restClient.Get(stub.ID)
		if err != nil {
			return err
		}
		exporter := &AlertingProfileExporter{AlertingProfile: alertingProfile}
		var file *os.File
		fileName := targetFolder + "/" + escape(alertingProfile.DisplayName) + ".alerting_profile.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		defer file.Close()
		exporter.ToHCL(file)
	}
	return nil
}

func importAutoTags(targetFolder string, environmentURL string, apiToken string) error {
	rest.Verbose = true
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := autotags.NewService(environmentURL+"/api/config/v1", apiToken)
	rest.Verbose = false
	stubList, err := restClient.ListAll()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Values {
		config, err := restClient.Get(stub.ID)
		if err != nil {
			return err
		}
		config.Metadata = nil
		var file *os.File
		fileName := targetFolder + "/" + escape(config.Name) + ".autotag.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		defer file.Close()
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_autotag", config.Name)); err != nil {
			return err
		}
		if err := hcl.Export(config, file); err != nil {
			return err
		}
		if _, err := file.WriteString("}\n"); err != nil {
			return err
		}
	}
	return nil
}

func importRequestAttributes(targetFolder string, environmentURL string, apiToken string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := requestattributes.NewService(environmentURL+"/api/config/v1", apiToken)
	rest.Verbose = false
	stubList, err := restClient.ListAll()
	if err != nil {
		return err
	}
	for _, stub := range stubList.Values {
		requestAttribute, err := restClient.Get(stub.ID)
		if err != nil {
			return err
		}
		exporter := &RequestAttributeExporter{RequestAttribute: requestAttribute}
		var file *os.File
		fileName := targetFolder + "/" + escape(requestAttribute.Name) + ".request_attribute.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		defer file.Close()
		exporter.ToHCL(file)
	}
	return nil
}

func importDashboards(targetFolder string, environmentURL string, apiToken string) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := dashboards.NewService(environmentURL+"/api/config/v1", apiToken)
	rest.Verbose = false
	dashboards, err := restClient.ListAll()
	if err != nil {
		return err
	}
	for _, dashboardStub := range dashboards.Dashboards {
		dashboard, err := restClient.Get(dashboardStub.ID)
		if err != nil {
			return err
		}
		exporter := &DashboardExporter{Dashboard: dashboard}
		var file *os.File

		fileName := targetFolder + "/" + escape(dashboard.Metadata.Name) + ".dashboard.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		defer file.Close()
		exporter.ToHCL(file)
	}
	return nil
}

func importCustomServices(targetFolder string, environmentURL string, apiToken string) error {
	if err := importCustomServicesTech(targetFolder, environmentURL, apiToken, customservices.Technologies.Java); err != nil {
		return err
	}
	if err := importCustomServicesTech(targetFolder, environmentURL, apiToken, customservices.Technologies.DotNet); err != nil {
		return err
	}
	if err := importCustomServicesTech(targetFolder, environmentURL, apiToken, customservices.Technologies.Go); err != nil {
		return err
	}
	if err := importCustomServicesTech(targetFolder, environmentURL, apiToken, customservices.Technologies.NodeJS); err != nil {
		return err
	}
	if err := importCustomServicesTech(targetFolder, environmentURL, apiToken, customservices.Technologies.PHP); err != nil {
		return err
	}
	return nil
}

func importCustomServicesTech(targetFolder string, environmentURL string, apiToken string, technology customservices.Technology) error {
	os.MkdirAll(targetFolder, os.ModePerm)
	restClient := customservices.NewService(environmentURL+"/api/config/v1", apiToken)
	rest.Verbose = false
	stubList, err := restClient.List(technology)
	if err != nil {
		return err
	}
	for _, stub := range stubList.Values {
		customService, err := restClient.Get(stub.ID, technology, false)
		if err != nil {
			return err
		}
		exporter := &CustomServiceExporter{CustomService: customService}
		var file *os.File
		fileName := targetFolder + "/" + escape(customService.Name) + ".custom_service.tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		defer file.Close()

		exporter.ToHCL(file, string(technology))
	}
	return nil
}
