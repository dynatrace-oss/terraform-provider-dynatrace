package iam

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Authenticator interface {
	ClientID() string
	AccountID() string
	ClientSecret() string
}

var tokens = map[string]string{}

type oauthResponse struct {
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
}

func getBearer(auth Authenticator, forceNew bool) (string, error) {
	var httpReq *http.Request
	var httpRes *http.Response
	var body []byte
	var err error

	if !forceNew {
		if token, found := tokens[auth.ClientID()+auth.AccountID()]; found {
			return token, nil
		}
	}

	httpClient := http.DefaultClient

	payloadStr := fmt.Sprintf(
		"grant_type=client_credentials&client_id=%s&client_secret=%s&scope=%s&resource=%s",
		url.QueryEscape(auth.ClientID()),
		url.QueryEscape(auth.ClientSecret()),
		url.QueryEscape("account-idm-read account-idm-write iam:policies:read iam:policies:write iam-policies-management"),
		url.QueryEscape("urn:dtaccount:"+auth.AccountID()),
	)
	payload := strings.NewReader(payloadStr)

	if httpReq, err = http.NewRequest(http.MethodPost, "https://sso.dynatrace.com/sso/oauth2/token", payload); err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if httpRes, err = httpClient.Do(httpReq); err != nil {
		return "", err
	}
	if body, err = io.ReadAll(httpRes.Body); err != nil {
		return "", err
	}
	response := oauthResponse{}
	if err = json.Unmarshal(body, &response); err != nil {
		return "", err
	}
	tokens[auth.ClientID()+auth.AccountID()] = response.AccessToken
	return response.AccessToken, nil
}
