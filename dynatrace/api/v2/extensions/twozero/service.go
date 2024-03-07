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

package twozero

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	twozero "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/extensions/twozero/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

func Service(credentials *settings.Credentials) settings.CRUDService[*twozero.Settings] {
	return &service{credentials}
}

type service struct {
	credentials *settings.Credentials
}

func (me *service) Get(id string, v *twozero.Settings) error {
	name, version := splitID(id)
	var extensionsList ExtensionsList
	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	if err := client.Get(fmt.Sprintf("/api/v2/extensions/info?name=%s", url.QueryEscape(name)), 200).Finish(&extensionsList); err != nil {
		return err
	}
	if len(extensionsList.Extensions) == 0 {
		return rest.Error{Code: 404, Message: fmt.Sprintf("extension `%s` is not installed", name)}
	}
	for _, extension := range extensionsList.Extensions {
		if extension.Name != name {
			continue
		}
		// if we want the version to get computed
		if len(version) == 0 {
			v.Name = name
			v.Version = extension.ActiveVersion
			return nil
		}

		// if we're looking for a specific version but a different one is active
		if extension.ActiveVersion != version {
			return rest.Error{Code: 404, Message: fmt.Sprintf("extension `%s` is installed, but in version %s (expected: %s)", name, extension.ActiveVersion, version)}
		}

		// success
		v.Name = name
		v.Version = version
		return nil
	}

	// If none of the extensions within the response matched
	return rest.Error{Code: 404, Message: fmt.Sprintf("extension `%s` is not installed", name)}
}

func (me *service) SchemaID() string {
	return "v2:extensions:twozero"
}

type ExtensionsList struct {
	Extensions []struct {
		Name          string `json:"extensionName"`
		Version       string `json:"version"`
		ActiveVersion string `json:"activeVersion"`
	} `json:"extensions"`
}

func joinID(name string, version string) string {
	return fmt.Sprintf("%s#-#%s", name, version)
}

func splitID(id string) (name string, version string) {
	parts := strings.Split(id, "#-#")
	if len(parts) > 0 {
		name = parts[0]
	}
	if len(parts) > 1 {
		version = parts[1]
	}
	return
}

func (me *service) List() (api.Stubs, error) {
	return me.list("")
}

func (me *service) list(exname string) (api.Stubs, error) {
	var stubs api.Stubs

	var extensionsList ExtensionsList
	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	query := ""
	if len(exname) > 0 {
		query = fmt.Sprintf("?name=%s", url.QueryEscape(exname))
	}
	if err := client.Get(fmt.Sprintf("/api/v2/extensions/info%s", query), 200).Finish(&extensionsList); err != nil {
		return stubs, err
	}
	for _, extension := range extensionsList.Extensions {
		if len(exname) > 0 && extension.Name != exname {
			continue
		}
		name := extension.Name
		version := extension.ActiveVersion
		stubs = append(stubs, &api.Stub{ID: joinID(name, version), Name: joinID(name, version)})
	}

	return stubs, nil
}

func (me *service) Validate(v *twozero.Settings) error {
	return nil // no endpoint for that
}

func (me *service) Create(v *twozero.Settings) (*api.Stub, error) {
	stub, err := me.ensureInstalled(v)
	if err != nil {
		return stub, err
	}
	name, version := splitID(stub.ID)
	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	payload := struct {
		Version string `json:"version"`
	}{
		Version: version,
	}
	return stub, client.Put(fmt.Sprintf("/api/v2/extensions/%s/environmentConfiguration", url.PathEscape(name)), &payload, 200).Finish(nil)
}

func (me *service) ensureInstalled(v *twozero.Settings) (*api.Stub, error) {
	name := v.Name
	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	response := struct {
		Name    string `json:"extensionName"`
		Version string `json:"extensionVersion"`
	}{}
	query := ""
	if len(v.Version) > 0 {
		query = "?version=" + url.QueryEscape(v.Version)
	}
	if err := client.Post(fmt.Sprintf("/api/v2/extensions/%s%s", url.PathEscape(name), query), nil, 200).Finish(&response); err != nil {
		if restErr, ok := err.(rest.Error); ok {
			if (restErr.Code == 400) && restErr.Message == fmt.Sprintf("Extension %s has already been added to environment", name) {
				stubs, err := me.list(name)
				if err != nil {
					return nil, err
				}
				if len(stubs) == 0 {
					return nil, rest.Error{Code: 404, Message: fmt.Sprintf("unable to find existing settings for extension %s", name)}
				}
				return stubs[0], nil
			}
		}
		return &api.Stub{}, err
	}
	return &api.Stub{ID: joinID(response.Name, response.Version), Name: joinID(response.Name, response.Version)}, nil
}

func (me *service) Update(id string, v *twozero.Settings) error {
	return errors.New("an update for this resource is not supposed to happen")
	// name, _ := splitID(id)
	// client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	// payload := struct {
	// 	Version string `json:"version"`
	// }{
	// 	Version: v.Version,
	// }
	// if err := client.Put(fmt.Sprintf("/api/v2/extensions/%s/environmentConfiguration", url.PathEscape(name)), &payload, 200).Finish(nil); err != nil {
	// 	if restErr, ok := err.(rest.Error); ok {
	// 		// Reasoning: We may want to update to a different version, but that version is not necessarily installed
	// 		if restErr.Code == 404 && restErr.Message == fmt.Sprintf("Extension %s, version %s not found or operation forbidden", name, v.Version) {
	// 			if _, err := me.ensureInstalled(v); err != nil {
	// 				return err
	// 			}
	// 			return me.Update(id, v)
	// 		}
	// 	}
	// }
	// return nil
}

func (me *service) Delete(id string) error {
	name, version := splitID(id)
	client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	if err := client.Delete(fmt.Sprintf("/api/v2/extensions/%s/environmentConfiguration", url.PathEscape(name)), 200).Finish(nil); err != nil {
		if restErr, ok := err.(rest.Error); ok {
			if restErr.Code != 404 {
				return err
			}
		}
	}
	if err := client.Delete(fmt.Sprintf("/api/v2/extensions/%s/%s", url.PathEscape(name), url.PathEscape(version)), 200).Finish(nil); err != nil {
		if restErr, ok := err.(rest.Error); ok {
			if restErr.Code != 404 {
				return err
			}
		}
	}
	return nil
}

func (me *service) New() *twozero.Settings {
	return new(twozero.Settings)
}

func (me *service) Name() string {
	return me.SchemaID()
}
