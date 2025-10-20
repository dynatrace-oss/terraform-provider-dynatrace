package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const mockToken = "########"

var credential_repo = map[string]Credentials{
	"unconfigured":                           {URL: TestCaseEnvURL},
	"api-token":                              {URL: TestCaseEnvURL, Token: mockToken},
	"oauth":                                  {URL: TestCaseEnvURL, OAuth: OAuthCredentials{ClientID: mockToken, ClientSecret: mockToken}},
	"platform-token":                         {URL: TestCaseEnvURL, OAuth: OAuthCredentials{PlatformToken: mockToken}},
	"api-token-and-oauth":                    {URL: TestCaseEnvURL, Token: mockToken, OAuth: OAuthCredentials{ClientID: mockToken, ClientSecret: mockToken}},
	"api-token-and-platform-token":           {URL: TestCaseEnvURL, Token: mockToken, OAuth: OAuthCredentials{PlatformToken: mockToken}},
	"oauth-and-platform-token":               {URL: TestCaseEnvURL, OAuth: OAuthCredentials{ClientID: mockToken, ClientSecret: mockToken, PlatformToken: mockToken}},
	"api-token-and-oauth-and-platform-token": {URL: TestCaseEnvURL, Token: mockToken, OAuth: OAuthCredentials{PlatformToken: mockToken, ClientID: mockToken, ClientSecret: mockToken}},
}

type testcase struct {
	credentials      creds_with_name
	IsOAuthPreferred bool
	Expected         error
}

func (t testcase) Credentials() *Credentials {
	return &t.credentials.Credentials
}

var expectedAPITokenError = errors.New("No API Token has been specified")
var expectedOAuthCredsError = errors.New("Neither OAuth Credentials nor Platform Token have been specified")
var classicChosen = errors.New("classic")
var platformChosen = errors.New("platform")

type creds_with_name struct {
	Name        string
	Credentials Credentials
}

func credentials(name string) creds_with_name {
	return creds_with_name{Name: name, Credentials: credential_repo[name]}
}

var testcases = []testcase{
	{
		credentials:      credentials("unconfigured"),
		IsOAuthPreferred: false,
		Expected:         expectedAPITokenError,
	},
	{
		credentials:      credentials("unconfigured"),
		IsOAuthPreferred: true,
		Expected:         expectedOAuthCredsError,
	},
	{
		credentials:      credentials("api-token"),
		IsOAuthPreferred: false,
		Expected:         classicChosen,
	},
	{
		credentials:      credentials("api-token"),
		IsOAuthPreferred: true,
		Expected:         classicChosen,
	},
	{
		credentials:      credentials("oauth"),
		IsOAuthPreferred: false,
		Expected:         platformChosen,
	},
	{
		credentials:      credentials("oauth"),
		IsOAuthPreferred: true,
		Expected:         platformChosen,
	},
	{
		credentials:      credentials("platform-token"),
		IsOAuthPreferred: false,
		Expected:         platformChosen,
	},
	{
		credentials:      credentials("platform-token"),
		IsOAuthPreferred: true,
		Expected:         platformChosen,
	},
	{
		credentials:      credentials("oauth-and-platform-token"),
		IsOAuthPreferred: false,
		Expected:         platformChosen,
	},
	{
		credentials:      credentials("oauth-and-platform-token"),
		IsOAuthPreferred: true,
		Expected:         platformChosen,
	},
	{
		credentials:      credentials("api-token-and-oauth"),
		IsOAuthPreferred: false,
		Expected:         classicChosen,
	},
	{
		credentials:      credentials("api-token-and-oauth"),
		IsOAuthPreferred: true,
		Expected:         platformChosen,
	},
	{
		credentials:      credentials("api-token-and-platform-token"),
		IsOAuthPreferred: false,
		Expected:         classicChosen,
	},
	{
		credentials:      credentials("api-token-and-platform-token"),
		IsOAuthPreferred: true,
		Expected:         platformChosen,
	},
	{
		credentials:      credentials("api-token-and-oauth-and-platform-token"),
		IsOAuthPreferred: false,
		Expected:         classicChosen,
	},
	{
		credentials:      credentials("api-token-and-oauth-and-platform-token"),
		IsOAuthPreferred: true,
		Expected:         platformChosen,
	},
}

func TestHybridClient(t *testing.T) {
	for _, testcase := range testcases {
		testCaseName := testcase.credentials.Name
		if testcase.IsOAuthPreferred {
			testCaseName = testCaseName + "/oauth-pref"
		}
		t.Run(testCaseName, func(t *testing.T) {
			t.Parallel()
			ctx := context.WithValue(t.Context(), "DYNATRACE_HTTP_OAUTH_PREFERENCE", testcase.IsOAuthPreferred)
			expect(t, testcase.Expected, HybridClient(testcase.Credentials()).Get(ctx, "").Finish())
		})
	}
}

func TestApiTokenClient(t *testing.T) {
	const (
		endpoint          = "/api/v2/settings/objects"
		activeGatePostfix = "/e/my-env-id"
	)

	t.Run("Correctly transforms an Active Gate URL", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			expectedURL, err := url.JoinPath(activeGatePostfix, endpoint)
			require.NoError(t, err)
			require.Equal(t, expectedURL, r.URL.Path)
			require.Equal(t, expectedURL, r.URL.Path)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("{}"))
		}))
		activeGateURL, err := url.JoinPath(server.URL, activeGatePostfix)
		require.NoError(t, err)

		cred := Credentials{URL: activeGateURL, Token: mockToken}
		client := HybridClient(&cred)

		req := client.Get(t.Context(), endpoint)
		err = req.Finish()
		assert.NoError(t, err)
	})

	t.Run("Errors on empty env URL", func(t *testing.T) {
		cred := Credentials{URL: "", Token: mockToken}
		client := HybridClient(&cred)

		req := client.Get(t.Context(), endpoint)
		err := req.Finish()
		assert.ErrorIs(t, err, NoClassicURLDefinedErr)
	})

	t.Run("Errors on invalid path", func(t *testing.T) {
		cred := Credentials{URL: "my-url", Token: mockToken}
		client := HybridClient(&cred)

		req := client.Get(t.Context(), ":/invalid-url")
		err := req.Finish()
		expectedErr := &url.Error{}
		assert.ErrorAs(t, err, &expectedErr)
	})

	t.Run("Sets correct headers", func(t *testing.T) {
		customHeaderName := "my-header"
		customHeaderValue := "my-value"
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			expectedURL, err := url.JoinPath(activeGatePostfix, endpoint)
			require.NoError(t, err)
			require.Equal(t, expectedURL, r.URL.Path)
			require.Equal(t, fmt.Sprintf("Api-Token %s", mockToken), r.Header.Get("Authorization"))
			require.Equal(t, customHeaderValue, r.Header.Get(customHeaderName))
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("{}"))
		}))
		activeGateURL, err := url.JoinPath(server.URL, activeGatePostfix)
		require.NoError(t, err)

		cred := Credentials{URL: activeGateURL, Token: mockToken}
		client := HybridClient(&cred)

		req := client.Get(t.Context(), endpoint)
		req.SetHeader(customHeaderName, customHeaderValue)
		err = req.Finish()
		assert.NoError(t, err)
	})
}

func expect(t *testing.T, expected error, actual error) {
	if expected == nil {
		if actual == nil {
			return
		}
		actualError := actual.Error()
		if idx := strings.Index(actualError, "."); idx >= 0 {
			actualError = actualError[:idx]
		}
		t.Errorf("expected no error, actual: %s", actualError)
		t.FailNow()
	}
	if actual == nil {
		t.Errorf("expected: '%s...', but no error", expected.Error())
		t.FailNow()
	}
	if strings.HasPrefix(actual.Error(), expected.Error()) {
		return
	}
	actualError := actual.Error()
	if idx := strings.Index(actualError, "."); idx >= 0 {
		actualError = actualError[:idx]
	}
	t.Errorf("expected: '%s...', actual: '%s'", expected.Error(), actualError)
	t.FailNow()
}
