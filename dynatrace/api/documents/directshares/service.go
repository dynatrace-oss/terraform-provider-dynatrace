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

package directshares

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	coreapi "github.com/dynatrace/dynatrace-configuration-as-code-core/api"
	coredirectshares "github.com/dynatrace/dynatrace-configuration-as-code-core/clients/directshares"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	serviceSettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/documents/directshares/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

func Service(credentials *rest.Credentials) settings.CRUDService[*serviceSettings.DirectShare] {
	return &service{credentials}
}

type service struct {
	credentials *rest.Credentials
}

type directShareDTO struct {
	ID         string   `json:"id"`
	DocumentId string   `json:"documentId"`
	Access     []string `json:"access"`
}

type createDirectShareDTO struct {
	DocumentId string         `json:"documentId"`
	Access     string         `json:"access"`
	Recipients []recipientDTO `json:"recipients"`
}

type recipientDTO struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type addDirectShareRecipientsDTO struct {
	Recipients []recipientDTO `json:"recipients"`
}

type removeDirectShareRecipientsDTO struct {
	Ids []string `json:"ids"`
}

func (me *service) client(ctx context.Context) (*coredirectshares.Client, error) {
	platformClient, err := rest.CreatePlatformClient(ctx, me.credentials.OAuth.EnvironmentURL, me.credentials)
	if err != nil {
		return nil, err
	}
	return coredirectshares.NewClient(platformClient), nil
}

func (me *service) Get(ctx context.Context, id string, v *serviceSettings.DirectShare) (err error) {
	client, err := me.client(ctx)
	if err != nil {
		return err
	}

	result, err := client.Get(ctx, id)
	if err != nil {
		apiErr := coreapi.APIError{}
		if errors.As(err, &apiErr) {
			return rest.Error{Code: apiErr.StatusCode, Message: apiErr.Error()}
		}
		return err
	}

	var ds directShareDTO
	if err := json.Unmarshal(result.Data, &ds); err != nil {
		return err
	}

	v.ID = ds.ID
	v.DocumentId = ds.DocumentId
	v.Access = strings.Join(ds.Access, "-")

	recipients, err := me.getRecipients(ctx, client, id)
	if err != nil {
		apiErr := coreapi.APIError{}
		if errors.As(err, &apiErr) {
			return rest.Error{Code: apiErr.StatusCode, Message: apiErr.Error()}
		}
		return err
	}

	v.Recipients = recipients
	return nil
}

func (me *service) getRecipients(ctx context.Context, client *coredirectshares.Client, id string) (serviceSettings.Recipients, error) {
	recipientsResp, err := client.GetRecipients(ctx, id)
	if err != nil {
		return nil, err
	}

	allRecipientObjs := recipientsResp.All()
	recipients := make(serviceSettings.Recipients, 0, len(allRecipientObjs))
	for _, obj := range allRecipientObjs {
		var rec recipientDTO
		if err := json.Unmarshal(obj, &rec); err != nil {
			return nil, err
		}
		recipients = append(recipients, &serviceSettings.Recipient{
			ID:   rec.ID,
			Type: rec.Type,
		})
	}

	return recipients, nil
}

func (me *service) SchemaID() string {
	return "document:direct-shares"
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	var stubs api.Stubs
	return stubs, nil // not implemented
}

func (me *service) Validate(v *serviceSettings.DirectShare) error {
	return nil // no endpoint for that
}

func (me *service) Create(ctx context.Context, v *serviceSettings.DirectShare) (stub *api.Stub, err error) {
	client, err := me.client(ctx)
	if err != nil {
		return nil, err
	}

	recipients := make([]recipientDTO, len(v.Recipients))
	for i, r := range v.Recipients {
		recipients[i] = recipientDTO{
			ID:   r.ID,
			Type: r.Type,
		}
	}

	directShareObj := createDirectShareDTO{
		DocumentId: v.DocumentId,
		Access:     v.Access,
		Recipients: recipients,
	}

	data, err := json.Marshal(directShareObj)
	if err != nil {
		return nil, err
	}

	result, err := client.Create(ctx, data)
	if err != nil {
		apiErr := coreapi.APIError{}
		if errors.As(err, &apiErr) {
			return nil, rest.Error{Code: apiErr.StatusCode, Message: string(apiErr.Body)}
		}
		return nil, err
	}

	var resp struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(result.Data, &resp); err != nil {
		return nil, err
	}
	return &api.Stub{ID: resp.ID}, nil
}

func (me *service) Update(ctx context.Context, id string, v *serviceSettings.DirectShare) (err error) {
	client, err := me.client(ctx)
	if err != nil {
		return err
	}

	remoteRecipients, err := me.getRecipients(ctx, client, id)
	if err != nil {
		apiErr := coreapi.APIError{}
		if errors.As(err, &apiErr) {
			return rest.Error{Code: apiErr.StatusCode, Message: string(apiErr.Body)}
		}
		return err
	}

	if err := me.addRecipients(ctx, client, id, v.Recipients, remoteRecipients); err != nil {
		return err
	}

	if err := me.removeRecipients(ctx, client, id, v.Recipients, remoteRecipients); err != nil {
		return err
	}

	return nil
}

func (me *service) addRecipients(ctx context.Context, client *coredirectshares.Client, id string, recipients serviceSettings.Recipients, remoteRecipients serviceSettings.Recipients) error {
	var add addDirectShareRecipientsDTO
	for _, desired := range recipients {
		if !containsRecipient(remoteRecipients, desired) {
			add.Recipients = append(add.Recipients, recipientDTO{
				ID:   desired.ID,
				Type: desired.Type,
			})
		}
	}

	if len(add.Recipients) == 0 {
		return nil
	}

	data, err := json.Marshal(add)
	if err != nil {
		return err
	}

	if err := client.AddRecipients(ctx, id, data); err != nil {
		apiErr := coreapi.APIError{}
		if errors.As(err, &apiErr) {
			return rest.Error{Code: apiErr.StatusCode, Message: string(apiErr.Body)}
		}
		return err
	}
	return nil
}

func (me *service) removeRecipients(ctx context.Context, client *coredirectshares.Client, id string, recipients serviceSettings.Recipients, remoteRecipients serviceSettings.Recipients) error {
	var remove removeDirectShareRecipientsDTO
	for _, remote := range remoteRecipients {
		if !containsRecipient(recipients, remote) {
			remove.Ids = append(remove.Ids, remote.ID)
		}
	}

	if len(remove.Ids) == 0 {
		return nil
	}

	data, err := json.Marshal(remove)
	if err != nil {
		return err
	}

	if err := client.RemoveRecipients(ctx, id, data); err != nil {
		apiErr := coreapi.APIError{}
		if errors.As(err, &apiErr) {
			return rest.Error{Code: apiErr.StatusCode, Message: string(apiErr.Body)}
		}
		return err
	}
	return nil
}

func containsRecipient(recipients serviceSettings.Recipients, target *serviceSettings.Recipient) bool {
	for _, r := range recipients {
		if r.ID == target.ID && r.Type == target.Type {
			return true
		}
	}
	return false
}

func (me *service) Delete(ctx context.Context, id string) error {
	client, err := me.client(ctx)
	if err != nil {
		return err
	}
	if err := client.Delete(ctx, id); err != nil {
		apiErr := coreapi.APIError{}
		if errors.As(err, &apiErr) {
			return rest.Error{Code: apiErr.StatusCode, Message: string(apiErr.Body)}
		}
		return err
	}
	return nil
}

func (me *service) New() *serviceSettings.DirectShare {
	return new(serviceSettings.DirectShare)
}
