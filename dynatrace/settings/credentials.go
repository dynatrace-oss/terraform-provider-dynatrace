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

const (
	ProdTokenURL   = "https://sso.dynatrace.com/sso/oauth2/token"
	SprintTokenURL = "https://sso.dynatracelabs.com/sso/oauth2/token"
	DevTokenURL    = "https://sso-dev.dynatracelabs.com/sso/oauth2/token"
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
	Cluster struct {
		URL   string
		Token string
	}
}

func getEnv(names ...string) string {
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

func CreateExportCredentials() (*Credentials, error) {
	environmentURL := os.Getenv("DT_SOURCE_ENV_URL")
	if len(environmentURL) == 0 {
		environmentURL = os.Getenv("DYNATRACE_SOURCE_ENV_URL")
	}
	if environmentURL == "" {
		environmentURL = os.Getenv("DYNATRACE_ENV_URL")
	}
	environmentURL = strings.TrimSuffix(strings.TrimSuffix(environmentURL, " "), "/")
	if len(environmentURL) != 0 {
		re := regexp.MustCompile(`https:\/\/(.*).(live|apps).dynatrace.com`)
		if match := re.FindStringSubmatch(environmentURL); len(match) > 0 {
			environmentURL = fmt.Sprintf("https://%s.live.dynatrace.com", match[1])
		}
	}

	apiToken := os.Getenv("DT_SOURCE_API_TOKEN")
	if len(apiToken) == 0 {
		apiToken = os.Getenv("DYNATRACE_SOURCE_API_TOKEN")
	}
	if apiToken == "" {
		apiToken = os.Getenv("DYNATRACE_API_TOKEN")
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
		if match := re.FindStringSubmatch(environmentURL); len(match) > 0 {
			automationEnvironmentURL = fmt.Sprintf("https://%s.apps.dynatrace.com", match[1])
			automationTokenURL = "https://sso.dynatrace.com/sso/oauth2/token"
		}
	}
	if len(automationTokenURL) == 0 {
		if strings.Contains(automationEnvironmentURL, ".live.dynatrace.com") || strings.Contains(automationEnvironmentURL, ".apps.dynatrace.com") {
			automationTokenURL = ProdTokenURL
		} else if strings.Contains(automationEnvironmentURL, ".sprint.dynatracelabs.com") || strings.Contains(automationEnvironmentURL, ".sprint.apps.dynatracelabs.com") {
			automationTokenURL = SprintTokenURL
		} else if strings.Contains(automationEnvironmentURL, ".dev.dynatracelabs.com") || strings.Contains(automationEnvironmentURL, ".dev.apps.dynatracelabs.com") {
			automationTokenURL = DevTokenURL
		}
	}

	client_id := getEnv("DT_CLIENT_ID", "DYNATRACE_CLIENT_ID")
	client_secret := getEnv("DYNATRACE_CLIENT_SECRET", "DT_CLIENT_SECRET")
	account_id := getEnv("DT_ACCOUNT_ID", "DYNATRACE_ACCOUNT_ID")

	iam_client_id := getEnv("IAM_CLIENT_ID", "DYNATRACE_IAM_CLIENT_ID", "DT_IAM_CLIENT_ID", "DT_CLIENT_ID", "DYNATRACE_CLIENT_ID")
	if len(iam_client_id) == 0 {
		iam_client_id = client_id
	}
	iam_account_id := getEnv("IAM_ACCOUNT_ID", "DYNATRACE_IAM_ACCOUNT_ID", "DT_IAM_ACCOUNT_ID", "DT_ACCOUNT_ID", "DYNATRACE_ACCOUNT_ID")
	if len(iam_account_id) == 0 {
		iam_account_id = account_id
	}
	iam_client_secret := getEnv("IAM_CLIENT_SECRET", "DYNATRACE_IAM_CLIENT_SECRET", "DT_IAM_CLIENT_SECRET", "DYNATRACE_CLIENT_SECRET", "DT_CLIENT_SECRET")
	if len(iam_client_secret) == 0 {
		iam_client_secret = client_secret
	}

	automation_client_id := getEnv("AUTOMATION_CLIENT_ID", "DYNATRACE_AUTOMATION_CLIENT_ID", "DT_AUTOMATION_CLIENT_ID", "DT_CLIENT_ID", "DYNATRACE_CLIENT_ID")
	if len(automation_client_id) == 0 {
		automation_client_id = client_id
	}
	automation_client_secret := getEnv("AUTOMATION_CLIENT_SECRET", "DYNATRACE_AUTOMATION_CLIENT_SECRET", "DT_AUTOMATION_CLIENT_SECRET", "DYNATRACE_CLIENT_SECRET", "DT_CLIENT_SECRET")
	if len(automation_client_secret) == 0 {
		automation_client_secret = client_secret
	}

	clusterURL := getEnv("DYNATRACE_CLUSTER_URL", "DT_CLUSTER_URL")
	for strings.HasSuffix(clusterURL, "/") {
		clusterURL = strings.TrimSuffix(clusterURL, "/")
	}
	clusterAPIToken := getEnv("DYNATRACE_CLUSTER_API_TOKEN", "DT_CLUSTER_API_TOKEN")

	if environmentURL == "" && clusterURL == "" && iam_account_id == "" {
		return nil, errors.New("the environment variable DYNATRACE_ENV_URL or DYNATRACE_SOURCE_ENV_URL needs to be set")
	}
	if apiToken == "" && clusterAPIToken == "" && iam_client_id == "" && iam_client_secret == "" {
		return nil, errors.New("the environment variable DYNATRACE_API_TOKEN or DYNATRACE_SOURCE_API_TOKEN needs to be set")
	}

	credentials := &Credentials{
		URL:   environmentURL,
		Token: apiToken,
		IAM: struct {
			ClientID     string
			AccountID    string
			ClientSecret string
		}{
			ClientID:     iam_client_id,
			AccountID:    iam_account_id,
			ClientSecret: iam_client_secret,
		},
		Automation: struct {
			ClientID       string
			ClientSecret   string
			TokenURL       string
			EnvironmentURL string
		}{
			ClientID:       automation_client_id,
			ClientSecret:   automation_client_secret,
			EnvironmentURL: automationEnvironmentURL,
			TokenURL:       automationTokenURL,
		},
		Cluster: struct {
			URL   string
			Token string
		}{
			URL:   clusterURL,
			Token: clusterAPIToken,
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
		if match := re.FindStringSubmatch(environmentURL); len(match) > 0 {
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
