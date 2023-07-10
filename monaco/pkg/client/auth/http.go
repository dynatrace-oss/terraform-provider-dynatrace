/*
 * @license
 * Copyright 2023 Dynatrace LLC
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package auth

import (
	"context"
	"log"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var DefaultTokenEndpoint = oauth2.Endpoint{
	TokenURL: "https://sso.dynatrace.com/sso/oauth2/token",
}

// OauthCredentials holds information for authenticating to Dynatrace
// using Oauth2.0 client credential flow
type OauthCredentials struct {
	ClientID     string
	ClientSecret string
	TokenURL     string
	Scopes       []string
}

// NewTokenAuthClient creates a new HTTP client that supports token based authorization
func NewTokenAuthClient(token string) *http.Client {
	if !(strings.HasPrefix(token, "dt0c01.") && strings.Count(token, ".") == 2) {
		log.Println("[WARN] You used an old token format. Please consider switching to the new 1.205+ token format.")
		log.Println("[WARN] More information: https://www.dynatrace.com/support/help/dynatrace-api/basics/dynatrace-api-authentication")
	}
	return &http.Client{Transport: NewTokenAuthTransport(token)}
}

// NewOAuthClient creates a new HTTP client that supports OAuth2 client credentials based authorization
func NewOAuthClient(ctx context.Context, oauthConfig OauthCredentials) *http.Client {

	tokenUrl := oauthConfig.TokenURL
	if tokenUrl == "" {
		log.Printf("[DEBUG] using default token URL %s", DefaultTokenEndpoint.TokenURL)
		tokenUrl = DefaultTokenEndpoint.TokenURL
	}

	config := clientcredentials.Config{
		ClientID:     oauthConfig.ClientID,
		ClientSecret: oauthConfig.ClientSecret,
		TokenURL:     tokenUrl,
		Scopes:       oauthConfig.Scopes,
	}

	return config.Client(ctx)
}
