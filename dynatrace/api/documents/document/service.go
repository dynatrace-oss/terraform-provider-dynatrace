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
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/envutils"

	docclient "github.com/dynatrace/dynatrace-configuration-as-code-core/clients/documents"

	documents "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/documents/document/settings"
)

var supportedDocumentTypes = []docclient.DocumentType{
	docclient.Dashboard,
	docclient.Notebook,
	docclient.Launchpad,
}

func Service(clientSet rest.ClientSet) (settings.CRUDService[*documents.Document], error) {
	platformClient, err := clientSet.PlatformClient()
	if err != nil {
		return nil, err
	}

	return &service{client: docclient.NewClient(platformClient)}, nil
}

type service struct {
	client *docclient.Client
}

func (me *service) Get(ctx context.Context, id string, v *documents.Document) (err error) {
	err = me.get(ctx, id, v)
	if envutils.DTDocumentsIgnoreUnexpectedEOF.Get() && err != nil {
		if strings.Contains(err.Error(), "unexpected EOF") {
			cfg := ctx.Value(settings.ContextKeyStateConfig)
			if stateDocument, ok := cfg.(*documents.Document); ok {
				v.ID = stateDocument.ID
				v.Name = stateDocument.Name
				v.Content = stateDocument.Content
				v.IsPrivate = stateDocument.IsPrivate
				v.Type = stateDocument.Type
				v.Owner = stateDocument.Owner
				v.Version = stateDocument.Version
				return nil
			}
		}
	}
	return err
}

func (me *service) get(ctx context.Context, id string, v *documents.Document) (err error) {
	result, err := me.client.Get(ctx, id)
	if err != nil {
		return err
	}

	v.ID = result.ID
	v.Content = string(result.Data)
	v.IsPrivate = result.IsPrivate
	v.Name = result.Name
	v.Owner = result.Owner
	v.Type = result.Type
	v.Version = result.Version

	return nil
}

func (me *service) SchemaID() string {
	return "document:documents"
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	listResponse, err := me.client.List(ctx, getSupportedDocumentsFilter())
	if err != nil {
		return nil, err
	}
	var stubs api.Stubs
	for _, response := range listResponse.Responses {
		stubs = append(stubs, &api.Stub{ID: response.ID, Name: response.Name, Extra: map[string]any{"type": response.Type, "owner": response.Owner}})
	}

	return stubs, nil
}

// getSupportedDocumentsFilter returns a filter string for supported documents, those of known types and not created by apps or extensions.
func getSupportedDocumentsFilter() string {
	filterBuilder := strings.Builder{}

	// exclude readmade documents created by apps or extensions
	filterBuilder.WriteString("not (originAppId exists) and not (originExtensionId exists) and ")

	// add filter for the supported document types
	filterBuilder.WriteString("type in (")
	for i, documentType := range supportedDocumentTypes {
		if i > 0 {
			filterBuilder.WriteString(", ")
		}
		filterBuilder.WriteString(fmt.Sprintf("'%s'", documentType))
	}
	filterBuilder.WriteString(")")

	return filterBuilder.String()
}

func (me *service) Validate(_ *documents.Document) error {
	return nil // no endpoint for that
}

func (me *service) Create(ctx context.Context, v *documents.Document) (*api.Stub, error) {
	stub, err := me.createPrivate(ctx, v)
	if err != nil {
		return nil, err
	}

	if !v.IsPrivate {
		if err = me.update(ctx, stub.ID, v); err != nil {
			return nil, err
		}
	}
	return stub, nil
}

func (me *service) createPrivate(ctx context.Context, v *documents.Document) (stub *api.Stub, err error) {
	response, err := me.client.Create(ctx, v.Name, v.IsPrivate, v.ID, []byte(v.Content), docclient.DocumentType(v.Type))
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(response.Data, &stub); err != nil {
		return nil, err
	}
	return stub, nil

}

func (me *service) Update(ctx context.Context, id string, v *documents.Document) (err error) {
	return me.update(ctx, id, v)
}

func (me *service) update(ctx context.Context, id string, v *documents.Document) (err error) {
	_, err = me.client.Update(ctx, id, v.Name, v.IsPrivate, []byte(v.Content), docclient.DocumentType(v.Type))
	if err != nil {
		return err
	}

	return nil
}

func (me *service) Delete(ctx context.Context, id string) error {
	_, err := me.client.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (me *service) New() *documents.Document {
	return new(documents.Document)
}
