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

package settings

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func GetEnv(names ...string) string {
	if len(names) == 0 {
		return ""
	}
	for _, name := range names {
		if value := os.Getenv(name); len(value) > 0 {
			return value
		}
	}
	return ""
}

type Credentials struct {
	URL   string
	Token string
	IAM   struct {
		ClientID     string
		AccountID    string
		ClientSecret string
	}
	Automation struct {
		ClientID       string
		ClientSecret   string
		TokenURL       string
		EnvironmentURL string
	}
}

func CreateExportCredentials() (*Credentials, error) {
	environmentURL := os.Getenv("DT_SOURCE_ENV_URL")
	if len(environmentURL) == 0 {
		environmentURL = os.Getenv("DYNATRACE_SOURCE_ENV_URL")
	}
	if environmentURL == "" {
		environmentURL = os.Getenv("DYNATRACE_ENV_URL")
		if environmentURL == "" {
			return nil, errors.New("the environment variable DYNATRACE_ENV_URL or DYNATRACE_SOURCE_ENV_URL needs to be set")
		}
	}
	environmentURL = strings.TrimSuffix(strings.TrimSuffix(environmentURL, " "), "/")
	apiToken := os.Getenv("DT_SOURCE_API_TOKEN")
	if len(apiToken) == 0 {
		apiToken = os.Getenv("DYNATRACE_SOURCE_API_TOKEN")
	}
	if apiToken == "" {
		apiToken = os.Getenv("DYNATRACE_API_TOKEN")
		if apiToken == "" {
			return nil, errors.New("the environment variable DYNATRACE_API_TOKEN or DYNATRACE_SOURCE_API_TOKEN needs to be set")
		}
	}
	automationEnvironmentURL := os.Getenv("DT_AUTOMATION_ENVIRONMENT_URL")
	if len(automationEnvironmentURL) == 0 {
		automationEnvironmentURL = os.Getenv("DYNATRACE_AUTOMATION_ENVIRONMENT_URL")
	}
	automationTokenURL := os.Getenv("DT_AUTOMATION_TOKEN_URL")
	if len(automationTokenURL) == 0 {
		automationTokenURL = os.Getenv("DYNATRACE_AUTOMATION_TOKEN_URL")
	}
	if len(automationEnvironmentURL) == 0 {
		re := regexp.MustCompile(`https:\/\/(.*).(live|apps).dynatrace.com`)
		if match := re.FindStringSubmatch(environmentURL); match != nil && len(match) > 0 {
			automationEnvironmentURL = fmt.Sprintf("https://%s.apps.dynatrace.com", match[1])
			automationTokenURL = "https://sso.dynatrace.com/sso/oauth2/token"
		}
	}
	automationClientID := os.Getenv("DT_AUTOMATION_CLIENT_ID")
	if len(automationClientID) == 0 {
		automationClientID = os.Getenv("DYNATRACE_AUTOMATION_CLIENT_ID")
	}
	automationClientSecret := os.Getenv("DT_AUTOMATION_CLIENT_SECRET")
	if len(automationClientSecret) == 0 {
		automationClientSecret = os.Getenv("DYNATRACE_AUTOMATION_CLIENT_SECRET")
	}
	credentials := &Credentials{
		URL:   environmentURL,
		Token: apiToken,
		IAM: struct {
			ClientID     string
			AccountID    string
			ClientSecret string
		}{
			ClientID:     os.Getenv("DT_CLIENT_ID"),
			AccountID:    os.Getenv("DT_ACCOUNT_ID"),
			ClientSecret: os.Getenv("DT_CLIENT_SECRET"),
		},
		Automation: struct {
			ClientID       string
			ClientSecret   string
			TokenURL       string
			EnvironmentURL string
		}{
			ClientID:       automationClientID,
			ClientSecret:   automationClientSecret,
			EnvironmentURL: automationEnvironmentURL,
			TokenURL:       automationTokenURL,
		},
	}
	return credentials, nil
}

func CreateCredentials() (*Credentials, error) {
	environmentURL := os.Getenv("DYNATRACE_ENV_URL")
	if environmentURL == "" {
		return nil, errors.New("the environment variable DYNATRACE_ENV_URL needs to be set")
	}
	environmentURL = strings.TrimSuffix(strings.TrimSuffix(environmentURL, " "), "/")
	apiToken := os.Getenv("DYNATRACE_API_TOKEN")
	if apiToken == "" {
		return nil, errors.New("the environment variable DYNATRACE_API_TOKEN needs to be set")
	}
	automationEnvironmentURL := os.Getenv("DT_AUTOMATION_ENVIRONMENT_URL")
	automationTokenURL := os.Getenv("DT_AUTOMATION_TOKEN_URL")
	if len(automationEnvironmentURL) == 0 {
		re := regexp.MustCompile(`https:\/\/(.*).(live|apps).dynatrace.com`)
		if match := re.FindStringSubmatch(environmentURL); match != nil && len(match) > 0 {
			automationEnvironmentURL = fmt.Sprintf("https://%s.apps.dynatrace.com", match[1])
			automationTokenURL = "https://sso.dynatrace.com/sso/oauth2/token"
		}
	}
	credentials := &Credentials{
		URL:   environmentURL,
		Token: apiToken,
		IAM: struct {
			ClientID     string
			AccountID    string
			ClientSecret string
		}{
			ClientID:     os.Getenv("DT_CLIENT_ID"),
			AccountID:    os.Getenv("DT_ACCOUNT_ID"),
			ClientSecret: os.Getenv("DT_CLIENT_SECRET"),
		},
		Automation: struct {
			ClientID       string
			ClientSecret   string
			TokenURL       string
			EnvironmentURL string
		}{
			ClientID:       os.Getenv("DT_AUTOMATION_CLIENT_ID"),
			ClientSecret:   os.Getenv("DT_AUTOMATION_CLIENT_SECRET"),
			EnvironmentURL: automationEnvironmentURL,
			TokenURL:       automationTokenURL,
		},
	}
	return credentials, nil
}
