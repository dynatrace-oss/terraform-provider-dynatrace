package iam

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
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

var mutex *sync.Mutex = new(sync.Mutex)

var msgInvalidOAuthCredentials = "Invalid OAuth credentials"

const errMsgClientIDMissing = ` No OAuth Client configured. Please specify either one of these environment variables: IAM_CLIENT_ID, DYNATRACE_IAM_CLIENT_ID, DT_IAM_CLIENT_ID, DT_CLIENT_ID, DYNATRACE_CLIENT_ID`
const errMsgAccountIDMissing = ` No Account ID configured. Please specify either one of these environment variables: IAM_ACCOUNT_ID, DYNATRACE_IAM_ACCOUNT_ID, DT_IAM_ACCOUNT_ID, DT_ACCOUNT_ID, DYNATRACE_ACCOUNT_ID`
const errMsgClientSecretMissing = ` No OAuth Client Secret configured. Please specify either one of these environment variables: IAM_CLIENT_SECRET, DYNATRACE_IAM_CLIENT_SECRET, DT_IAM_CLIENT_SECRET, DYNATRACE_CLIENT_SECRET, DT_CLIENT_SECRET`

func getBearer(auth Authenticator, forceNew bool) (string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	clientID := auth.ClientID()
	if len(strings.TrimSpace(clientID)) == 0 {
		return "", errors.New(errMsgClientIDMissing)
	}
	accountID := auth.AccountID()
	if len(strings.TrimSpace(accountID)) == 0 {
		return "", errors.New(errMsgAccountIDMissing)
	}
	clientSecret := auth.ClientSecret()
	if len(strings.TrimSpace(clientSecret)) == 0 {
		return "", errors.New(errMsgClientSecretMissing)
	}

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
		"grant_type=client_credentials&client_id=%s&client_secret=%s",
		url.QueryEscape(auth.ClientID()),
		url.QueryEscape(auth.ClientSecret()),
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
	if os.Getenv("DT_DEBUG_IAM_BEARER") == "true" {
		debugPayloadStr := fmt.Sprintf(
			"grant_type=client_credentials&client_id=%s&client_secret=%s",
			url.QueryEscape(auth.ClientID()),
			url.QueryEscape("<hidden>"),
		)
		rest.Logger.Println("POST https://sso.dynatrace.com/sso/oauth2/token")
		rest.Logger.Println("  " + debugPayloadStr)
		rest.Logger.Println("  -> " + string(body))
	}
	if httpRes.StatusCode == 400 {
		return "", errors.New(msgInvalidOAuthCredentials)
	}
	response := oauthResponse{}
	if err = json.Unmarshal(body, &response); err != nil {
		return "", err
	}
	tokens[auth.ClientID()+auth.AccountID()] = response.AccessToken
	return response.AccessToken, nil
}
