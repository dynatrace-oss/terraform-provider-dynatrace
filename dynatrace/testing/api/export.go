//go:build integration

/*
 * @license
 * Copyright 2026 Dynatrace LLC
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

package api

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/envutils"
)

func ExportResource(t *testing.T, resourceType string, id string, adminAccess bool) (string, error) {
	tmpDir := t.TempDir()
	t.Setenv(envutils.DynatraceTargetFolder.Key, tmpDir)

	resourceSelector := fmt.Sprintf("%s=%s", resourceType, id)
	adminAccessFlag := fmt.Sprintf("-admin-access=%t", adminAccess)
	args := []string{"path", "-export", "-skip-terraform-init", adminAccessFlag, resourceSelector}

	return tmpDir, dynatrace.Export(args, config.ConfigGetter{Provider: provider.Provider()})
}

func FindExportedResource(targetDir string, resourceType string, identifier string) (bool, error) {
	// <targetDir>\modules\resource_name\res.tf
	resourceFolder := path.Join(targetDir, "modules", strings.TrimPrefix(resourceType, "dynatrace_"))
	resourceFiles, err := os.ReadDir(resourceFolder)
	if err != nil {
		// dir not found => config doesn't exist
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	stringIdentifier := strconv.Quote(identifier) // quotes to find an exact match

	for _, entry := range resourceFiles {
		if entry.IsDir() {
			continue
		}
		content, readErr := os.ReadFile(path.Join(resourceFolder, entry.Name()))
		if readErr != nil {
			return false, readErr
		}

		// Look for "<identifier>"
		if strings.Contains(string(content), stringIdentifier) {
			return true, nil
		}
	}
	return false, nil
}
