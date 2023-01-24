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
	"os"
	"strings"
)

type Credentials struct {
	URL   string
	Token string
	IAM   struct {
		ClientID     string
		AccountID    string
		ClientSecret string
	}
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
	}
	return credentials, nil
}
