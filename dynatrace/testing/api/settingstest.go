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

package api

import (
	"encoding/json"
	"io/fs"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/assert"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/google/uuid"
)

type Anonymizer interface {
	Anonymize()
}

func load(path string, v any, randomize string) error {
	var err error
	var data []byte

	if data, err = os.ReadFile(path); err != nil {
		return err
	}

	if loader, ok := v.(settings.Loader); ok {
		return loader.Load([]byte(strings.ReplaceAll(string(data), "${randomize}", randomize)))
	}
	return json.Unmarshal([]byte(strings.ReplaceAll(string(data), "${randomize}", randomize)), v)
}

func TestService[V settings.Settings](t *testing.T, createService func(*settings.Credentials) settings.CRUDService[V]) {
	t.Helper()
	SettingsTest[V]{T: t}.Run(createService)
}

type SettingsTest[V settings.Settings] struct {
	T *testing.T
}

func (st SettingsTest[V]) Run(createService func(*settings.Credentials) settings.CRUDService[V]) {
	st.T.Helper()
	envURL := os.Getenv("DYNATRACE_ENV_URL")
	apiToken := os.Getenv("DYNATRACE_API_TOKEN")
	if envURL == "" || apiToken == "" {
		st.T.Skip("Environment Variables DYNATRACE_ENV_URL and DYNATRACE_API_TOKEN must be specified")
		return
	}

	entries, _ := os.ReadDir("testdata")
	for _, entry := range entries {
		if entry.IsDir() {
			var err error

			folder := path.Join("testdata", entry.Name())

			service := createService(&settings.Credentials{URL: envURL, Token: apiToken})

			randomize := uuid.NewString()

			var entries []fs.DirEntry
			if entries, err = os.ReadDir(folder); err != nil {
				st.T.Error(err)
			}
			allSettings := []V{}

			for _, entry := range entries {
				if !strings.HasSuffix(entry.Name(), ".json") {
					continue
				}
				settingsJSONFile := path.Join(folder, entry.Name())
				settings := settings.NewSettings(service.(settings.RService[V]))
				if err = load(settingsJSONFile, settings, randomize); err != nil {
					st.T.Error(err)
					return
				}
				allSettings = append(allSettings, settings)
			}

			if len(allSettings) == 0 {
				return
			}
			st.T.Run(entry.Name(), func(t *testing.T) {
				t.Helper()
				assert := assert.New(t)
				service := createService(&settings.Credentials{URL: envURL, Token: apiToken})

				var err error
				var stub *settings.Stub

				createSettings := allSettings[0]

				if validator, ok := service.(settings.Validator[V]); ok {
					if err = validator.Validate(createSettings); err != nil {
						assert.Error(err)
						return
					}
				}

				if t.Failed() {
					return
				}

				if stub, err = service.Create(createSettings); err != nil {
					assert.Error(err)
					return
				}

				if !t.Failed() {
					t.Cleanup(func() {
						if service != nil && stub != nil {
							if err = service.Delete(stub.ID); err != nil {
								assert.Error(err)
							}
						}
					})
				}

				remoteSettings := settings.NewSettings(service.(settings.RService[V]))
				if err = service.Get(stub.ID, remoteSettings); err != nil {
					assert.Error(err)
					return
				}

				FillDemoValues(remoteSettings)
				Anonymize(remoteSettings)
				Anonymize(createSettings)
				settings.ClearLegacyID(remoteSettings)

				assert.Equals(createSettings, remoteSettings)

				for idx := 1; idx < len(allSettings); idx++ {
					if t.Failed() {
						return
					}
					updateSettings := allSettings[idx]
					if err = service.Update(stub.ID, updateSettings); err != nil {
						assert.Error(err)
						return
					}

					remoteSettings := settings.NewSettings(service.(settings.RService[V]))
					if err = service.Get(stub.ID, remoteSettings); err != nil {
						assert.Error(err)
						return
					}
					Anonymize(remoteSettings)
					FillDemoValues(remoteSettings)
					Anonymize(updateSettings)
					settings.ClearLegacyID(remoteSettings)
					assert.Equals(updateSettings, remoteSettings)
				}

			})
		}

	}
}

func FillDemoValues(v any) {
	if demoFiller, ok := v.(settings.DemoSettings); ok {
		demoFiller.FillDemoValues()
	}

}

func Anonymize(v any) {
	if anonymizer, ok := v.(Anonymizer); ok {
		anonymizer.Anonymize()
	}

}

type TestAccOptions struct {
	ExpectNonEmptyPlan bool
}

func TestAcc(t *testing.T, opts ...TestAccOptions) {
	var options TestAccOptions
	if len(opts) > 0 {
		options = opts[0]
	}
	t.Helper()
	if v := os.Getenv("TF_ACC"); v == "" {
		t.Skip("TF_ACC has not been set for acceptance tests")
		return
	}
	if v := os.Getenv("DYNATRACE_ENV_URL"); v == "" {
		t.Skip("DYNATRACE_ENV_URL has not been set for acceptance tests")
		return
	}

	if v := os.Getenv("DYNATRACE_API_TOKEN"); v == "" {
		t.Skip("DYNATRACE_API_TOKEN must be set for acceptance tests")
		return
	}
	entries, _ := os.ReadDir("testdata")
	for _, entry := range entries {
		if entry.IsDir() {
			var err error
			folder := path.Join("testdata", entry.Name())
			allFiles := []string{}

			var entries []fs.DirEntry
			if entries, err = os.ReadDir(folder); err != nil {
				t.Error(err)
				return
			}
			for _, entry := range entries {
				if !strings.HasSuffix(entry.Name(), ".tf") {
					continue
				}
				allFiles = append(allFiles, path.Join(folder, entry.Name()))
			}
			for _, file := range allFiles {
				subTestName := strings.TrimPrefix(file, "testdata/")
				t.Run(subTestName, func(t *testing.T) {
					t.Helper()
					var err error

					var content []byte
					if content, err = os.ReadFile(file); err != nil {
						t.Error(err)
						return
					}
					config := string(content)
					name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
					config = strings.ReplaceAll(config, "#name#", name)
					config = strings.ReplaceAll(config, "${randomize}", name)
					provider := provider.Provider()
					testCase := &resource.TestCase{
						ProviderFactories: map[string]func() (*schema.Provider, error){
							"dynatrace": func() (*schema.Provider, error) {
								return provider, nil
							},
						},
						Steps: []resource.TestStep{{Config: config, ExpectNonEmptyPlan: options.ExpectNonEmptyPlan}},
					}
					resource.Test(t, *testCase)
				})
			}
		}
	}
}
