//go:build unit

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

package document

import (
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strings"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	testing2 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const documentsPath = "/platform/document/v1/documents"

// platformClientSet returns a ClientSet with a real platform client pointed at the test server.
func platformClientSet(t *testing.T, serverURL string) *testing2.MockClientSet {
	platformClient, err := rest.CreatePlatformClient(t.Context(), serverURL, &rest.Credentials{
		Platform: rest.PlatformCredentials{EnvironmentURL: serverURL, PlatformToken: "token"},
	})
	require.NoError(t, err)
	return &testing2.MockClientSet{PlatformClientValue: platformClient}
}

// writeDocumentResponse writes the multipart response the get endpoint replies with.
func writeDocumentResponse(t *testing.T, w http.ResponseWriter, metadata string, content string) {
	t.Helper()

	body := &strings.Builder{}
	writer := multipart.NewWriter(body)

	metadataPart, err := writer.CreatePart(textproto.MIMEHeader{"Content-Type": []string{"application/json"}, "Content-Disposition": []string{`form-data; name="metadata"`}})
	require.NoError(t, err)
	_, err = metadataPart.Write([]byte(metadata))
	require.NoError(t, err)

	contentPart, err := writer.CreateFormFile("content", "content.json")
	require.NoError(t, err)
	_, err = contentPart.Write([]byte(content))
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	w.Header().Set("Content-Type", "multipart/form-data; boundary="+writer.Boundary())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body.String()))
}

// TestDataSourceRead tests that the documents data source resolves isReshareable — the one field
// the list endpoint doesn't provide — by fetching every document individually.
func TestDataSourceRead(t *testing.T) {
	// the list endpoint populates description and labels, but never isReshareable
	listResponse := `{"totalCount": 2, "nextPageKey": null, "documents": [
		{"id": "dashboard-id", "name": "the dashboard", "type": "dashboard", "owner": "owner-id", "description": "the description", "labels": ["cloud", "monitoring"]},
		{"id": "notebook-id", "name": "the notebook", "type": "notebook", "owner": "owner-id"}]}`

	// every field except isReshareable is deliberately different from the list payload, so the
	// assertions below pin down which endpoint each value is read from
	getResponses := map[string]string{
		"dashboard-id": `{"id": "dashboard-id", "name": "from get", "type": "from-get", "owner": "from-get", "description": "from get", "labels": ["from-get"], "isReshareable": true}`,
		"notebook-id":  `{"id": "notebook-id", "name": "from get", "type": "from-get", "owner": "from-get"}`,
	}

	var requestedIDs []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		if r.URL.Path == documentsPath {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(listResponse))
			return
		}

		id := strings.TrimPrefix(r.URL.Path, documentsPath+"/")
		metadata, found := getResponses[id]
		require.True(t, found, "unexpected request for document %s", id)
		requestedIDs = append(requestedIDs, id)
		writeDocumentResponse(t, w, metadata, "{}")
	}))
	defer server.Close()

	t.Run("Reads every field from the list endpoint except isReshareable", func(t *testing.T) {
		requestedIDs = nil
		resourceData := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]any{"type": "dashboard"})

		diags := dataSourceRead(t.Context(), resourceData, platformClientSet(t, server.URL))

		require.Empty(t, diags)
		assert.Equal(t, []string{"dashboard-id"}, requestedIDs, "only documents of the requested type must be fetched")
		assert.Equal(t, "documents[dashboard]", resourceData.Id())

		values := resourceData.Get("values").([]any)
		require.Len(t, values, 1)
		value := values[0].(map[string]any)
		assert.Equal(t, "dashboard-id", value["id"])
		assert.Equal(t, "the dashboard", value["name"])
		assert.Equal(t, "dashboard", value["type"])
		assert.Equal(t, "owner-id", value["owner"])
		assert.Equal(t, "the description", value["description"])
		assert.Equal(t, true, value["is_reshareable"])
		assert.ElementsMatch(t, []any{"cloud", "monitoring"}, value["labels"].(*schema.Set).List())
	})

	t.Run("Returns all documents if no type is given", func(t *testing.T) {
		requestedIDs = nil
		resourceData := schema.TestResourceDataRaw(t, DataSource().Schema, map[string]any{})

		diags := dataSourceRead(t.Context(), resourceData, platformClientSet(t, server.URL))

		require.Empty(t, diags)
		assert.ElementsMatch(t, []string{"dashboard-id", "notebook-id"}, requestedIDs)
		assert.Equal(t, "documents", resourceData.Id())
		assert.Len(t, resourceData.Get("values").([]any), 2)
	})
}
