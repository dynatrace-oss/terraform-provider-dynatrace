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

package documents

import (
	"testing"

	documents "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/documents/document/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetSupportedDocumentsFilter tests that the getSupportedDocumentsFilter function returns a filter less than or equal to 1024 characters in length.
func TestGetSupportedDocumentsFilter(t *testing.T) {
	filter := getSupportedDocumentsFilter()
	assert.LessOrEqual(t, len(filter), 1024, "supported documents filter string is longer than 1024 characters")
}

// TestToMetadata tests that the optional fields are always part of the payload. An update is a
// PATCH, so a field left nil is dropped by the client and read as "leave unchanged" by the API —
// clearing a description or removing the last label has to send an explicit empty value.
func TestToMetadata(t *testing.T) {
	t.Run("Populates every field that was set", func(t *testing.T) {
		meta := toMetadata("document-id", &documents.Document{
			Name:          "the dashboard",
			Type:          "dashboard",
			IsPrivate:     true,
			Description:   "the description",
			Labels:        []string{"cloud", "monitoring"},
			IsReshareable: true,
		})

		assert.Equal(t, "document-id", meta.ID)
		assert.Equal(t, "the dashboard", meta.Name)
		assert.Equal(t, "dashboard", meta.Type)
		assert.True(t, meta.IsPrivate)
		assert.Equal(t, []string{"cloud", "monitoring"}, meta.Labels)
		require.NotNil(t, meta.Description)
		assert.Equal(t, "the description", *meta.Description)
		require.NotNil(t, meta.IsReshareable)
		assert.True(t, *meta.IsReshareable)
	})

	t.Run("Sends explicit empty values for cleared fields", func(t *testing.T) {
		meta := toMetadata("document-id", &documents.Document{Name: "the dashboard", Type: "dashboard"})

		require.NotNil(t, meta.Description, "a nil description is dropped from the PATCH, leaving the old value in place")
		assert.Equal(t, "", *meta.Description)
		require.NotNil(t, meta.Labels, "nil labels are dropped from the PATCH, leaving the old labels in place")
		assert.Empty(t, meta.Labels)
		require.NotNil(t, meta.IsReshareable)
		assert.False(t, *meta.IsReshareable)
	})
}
