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

package documents

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	docapi "github.com/dynatrace/dynatrace-configuration-as-code-core/api"
	docclient "github.com/dynatrace/dynatrace-configuration-as-code-core/clients/documents"

	documents "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/documents/document/settings"
)

func Service(credentials *rest.Credentials) settings.CRUDService[*documents.Document] {
	return &service{credentials}
}

var IGNORE_UNEXPECTED_EOF = (os.Getenv("DT_DOCUMENTS_IGNORE_UNEXPECTED_EOF") == "true")

type service struct {
	credentials *rest.Credentials
}

func (me *service) client(ctx context.Context) (*docclient.Client, error) {
	platformClient, err := rest.CreatePlatformClient(ctx, me.credentials.OAuth.EnvironmentURL, me.credentials)
	if err != nil {
		return nil, err
	}
	return docclient.NewClient(platformClient), nil
}

func (me *service) Get(ctx context.Context, id string, v *documents.Document) (err error) {
	err = me.get(ctx, id, v)
	if IGNORE_UNEXPECTED_EOF && err != nil {
		if strings.Contains(err.Error(), "unexpected EOF") {
			cfg := ctx.Value(settings.ContextKeyStateConfig)
			if stateDocument, ok := cfg.(*documents.Document); ok {
				v.Name = stateDocument.Name
				v.Content = stateDocument.Content
				v.IsPrivate = stateDocument.IsPrivate
				v.Type = stateDocument.Type
				v.Actor = stateDocument.Actor
				v.Owner = stateDocument.Owner
				v.Version = stateDocument.Version
				v.SchemaVersion = stateDocument.SchemaVersion
				return nil
			}
		}
	}
	return err
}

func (me *service) get(ctx context.Context, id string, v *documents.Document) (err error) {
	client, err := me.client(ctx)
	if err != nil {
		return err
	}
	result, err := client.Get(ctx, id)
	if err != nil {
		if apiError, ok := err.(docapi.APIError); ok {
			return rest.Error{Code: apiError.StatusCode, Message: apiError.Error()}
		}
		return err
	}

	v.Actor = result.Actor
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
	if me == nil {
		return api.Stubs{}, nil
	}
	cl, err := me.client(ctx)
	if err != nil {
		return nil, err
	}

	if cl == nil {
		return api.Stubs{}, nil
	}
	listResponse, err := cl.List(ctx, "")
	if err != nil {
		if apiError, ok := err.(docapi.APIError); ok {
			return nil, rest.Error{Code: apiError.StatusCode, Message: apiError.Error()}
		}
		return nil, err
	}
	var stubs api.Stubs
	for _, response := range listResponse.Responses {
		stubs = append(stubs, &api.Stub{ID: response.ID, Name: response.Name, Extra: map[string]any{"type": response.Type, "owner": response.Owner}})
	}

	return stubs, nil
}

func (me *service) Validate(_ *documents.Document) error {
	return nil // no endpoint for that
}

func (me *service) Create(ctx context.Context, v *documents.Document) (*api.Stub, error) {
	stub, err := me.createPrivate(ctx, v)
	if err != nil {
		if apiError, ok := err.(docapi.APIError); ok {
			return nil, rest.Error{Code: apiError.StatusCode, Message: apiError.Error()}
		}
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
	c, err := me.client(ctx)
	if err != nil {
		return nil, err
	}
	response, err := c.Create(ctx, v.Name, v.IsPrivate, "", []byte(v.Content), docclient.DocumentType(v.Type))
	if err != nil {
		if apiError, ok := err.(docapi.APIError); ok {
			return nil, rest.Error{Code: apiError.StatusCode, Message: string(apiError.Body)}
		}
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
	c, err := me.client(ctx)
	if err != nil {
		return err
	}
	_, err = c.Update(ctx, id, v.Name, v.IsPrivate, []byte(v.Content), docclient.DocumentType(v.Type))
	if err != nil {
		if apiError, ok := err.(docapi.APIError); ok {
			return rest.Error{Code: apiError.StatusCode, Message: string(apiError.Body)}
		}
		return err
	}

	return nil
}

func (me *service) Delete(ctx context.Context, id string) error {
	client, err := me.client(ctx)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, id)
	if err != nil {
		if apiError, ok := err.(docapi.APIError); ok {
			return rest.Error{Code: apiError.StatusCode, Message: string(apiError.Body)}
		}
		return err
	}

	return nil
}

func (me *service) New() *documents.Document {
	return new(documents.Document)
}
