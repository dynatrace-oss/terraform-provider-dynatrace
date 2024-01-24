/**
* @license
* Copyright 2024 Dynatrace LLC
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

package httpcache

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/address"
	vault "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/credentials/vault/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/cache/tar"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
)

type GetCredentialsVaultRequest struct {
	SchemaID        string
	ServiceSchemaID string
	ID              string
}

func (me *GetCredentialsVaultRequest) Raw() ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (me *GetCredentialsVaultRequest) Finish(vs ...any) error {
	var v any
	if len(vs) > 0 {
		v = vs[0]
	}

	tarFolder, _, err := tar.NewExisting(CACHE_FOLDER + "/" + strings.ReplaceAll(me.SchemaID, ":", ""))
	if err != nil {
		return err
	}

	if tarFolder != nil {
		stub, data, err := tarFolder.Get(me.ID)
		if err != nil {
			return err
		}
		if stub == nil {
			logging.Debug.Info.Printf("[HTTP_CACHE_CredentialsVault_Tar] [%s] [FAILED] [%s] 404 not found", me.SchemaID, me.ID)
			logging.Debug.Warn.Printf("[HTTP_CACHE_CredentialsVault_Tar] [%s] [FAILED] [%s] 404 not found", me.SchemaID, me.ID)
			return rest.Error{Code: 404, Message: fmt.Sprintf("CredentialsVault_Tar %s not found", me.ID)}
		}
		wrapper := struct {
			Downloaded struct {
				ClassidID string          `json:"classicId,omitempty"`
				Value     json.RawMessage `json:"value"`
			} `json:"downloaded"`
		}{}

		if err := json.Unmarshal(data, &wrapper); err != nil {
			return err
		}
		if err := json.Unmarshal(wrapper.Downloaded.Value, &v); err != nil {
			return err
		}

		address.AddToOriginal(address.AddressOriginal{
			TerraformSchemaID: me.ServiceSchemaID,
			OriginalID:        wrapper.Downloaded.ClassidID,
			OriginalSchemaID:  me.SchemaID,
		})
		return nil
	}
	logging.Debug.Info.Printf("[HTTP_CACHE_CredentialsVault_Tar Nil] [%s] [FAILED] [%s] 404 not found", me.SchemaID, me.ID)
	logging.Debug.Warn.Printf("[HTTP_CACHE_CredentialsVault_Tar Nil] [%s] [FAILED] [%s] 404 not found", me.SchemaID, me.ID)
	return rest.Error{Code: 404, Message: fmt.Sprintf("CredentialsVault_Tar Nil %s not found", me.ID)}
}

func (me *GetCredentialsVaultRequest) Expect(codes ...int) rest.Request {
	return me
}

func (me *GetCredentialsVaultRequest) Payload(any) rest.Request {
	return me
}

func (me *GetCredentialsVaultRequest) OnResponse(func(resp *http.Response)) rest.Request {
	return me
}

type ListCredentialsVaultRequest struct {
	SchemaID string
}

func (me *ListCredentialsVaultRequest) Raw() ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (me *ListCredentialsVaultRequest) Finish(vs ...any) error {
	var v any
	if len(vs) > 0 {
		v = vs[0]
	}
	credVault := CredentialsVaultListLocal{Credentials: []vault.CredentialsResponseElement{}}

	tarFolder, _, err := tar.NewExisting(CACHE_FOLDER + "/" + strings.ReplaceAll(me.SchemaID, ":", ""))
	if err != nil {
		return err
	}

	if tarFolder != nil {
		stubs, err := tarFolder.List()
		if err != nil {
			return err
		}

		for _, stub := range stubs {
			_, data, err := tarFolder.Get(stub.ID)
			if err != nil {
				return err
			}
			wrapper := struct {
				Downloaded struct {
					Value vault.CredentialsResponseElement `json:"value"`
				} `json:"downloaded"`
			}{}
			if err := json.Unmarshal(data, &wrapper); err != nil {
				return err
			}
			credVault.Credentials = append(credVault.Credentials, wrapper.Downloaded.Value)
		}
	}
	data, err := json.Marshal(credVault)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	return nil
}

func (me *ListCredentialsVaultRequest) Expect(codes ...int) rest.Request {
	return me
}

func (me *ListCredentialsVaultRequest) Payload(any) rest.Request {
	return me
}

func (me *ListCredentialsVaultRequest) OnResponse(func(resp *http.Response)) rest.Request {
	return me
}

type CredentialsVaultListLocal struct {
	Credentials []vault.CredentialsResponseElement `json:"credentials"` // A list of credentials sets for Synthetic monitors.
	PageSize    *int32                             `json:"pageSize"`
	NextPageKey *string                            `json:"nextPageKey,omitempty"`
	TotalCount  *int64                             `json:"totalCount"`
}
