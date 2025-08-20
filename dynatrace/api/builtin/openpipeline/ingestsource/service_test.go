/**
* @license
* Copyright 2025 Dynatrace LLC
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

package ingestsource_test

import (
	"context"
	"encoding/json"
	"sync/atomic"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/ingestsource"
	ingestsource2 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/ingestsource/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/processors"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	coreApi "github.com/dynatrace/dynatrace-configuration-as-code-core/api"
)

type genericClientStub struct {
	get    func(ctx context.Context, id string) (settings20.Response, error)
	delete func(ctx context.Context, id string) (settings20.Response, error)
}

func (c *genericClientStub) Get(ctx context.Context, objectID string) (settings20.Response, error) {
	return c.get(ctx, objectID)
}

func (c *genericClientStub) Delete(ctx context.Context, objectID string) (settings20.Response, error) {
	return c.delete(ctx, objectID)
}

type clientForKindStub struct {
	list   func(ctx context.Context) (api.Stubs, error)
	create func(ctx context.Context, v *ingestsource2.IngestSource) (*api.Stub, error)
	update func(ctx context.Context, id string, v *ingestsource2.IngestSource) error
}

func (c *clientForKindStub) List(ctx context.Context) (api.Stubs, error) {
	return c.list(ctx)
}

func (c *clientForKindStub) Create(ctx context.Context, v *ingestsource2.IngestSource) (*api.Stub, error) {
	return c.create(ctx, v)
}

func (c *clientForKindStub) Update(ctx context.Context, objectID string, v *ingestsource2.IngestSource) error {
	return c.update(ctx, objectID, v)
}

func TestService(t *testing.T) {
	bucketName := "bucket"
	pipelineID := "pipeline"
	sampleData := "some sample data"
	matcher := "not true"

	t.Run("Get", func(t *testing.T) {
		t.Run("It gets an ingest source", func(t *testing.T) {
			client := &genericClientStub{
				get: func(ctx context.Context, id string) (settings20.Response, error) {
					ingestSource := ingestsource2.IngestSource{
						Kind:          "events",
						DefaultBucket: &bucketName,
						DisplayName:   "displayName",
						Enabled:       false,
						PathSegment:   "some.path.segment",
						StaticRouting: &ingestsource2.PipelineReference{
							PipelineID:   &pipelineID,
							PipelineType: "custom",
						},
						Processing: &ingestsource2.Processing{
							Processors: []*processors.Processor{
								{
									Enabled:     true,
									Id:          "proc-2",
									Type:        processors.DqlProcessorType,
									Description: "my-proc-2",
									SampleData:  &sampleData,
									Matcher:     &matcher,
									Dql:         &processors.DqlAttributes{Script: "fieldsAdd true"},
								},
							},
						},
					}

					valueBytes, err := ingestSource.MarshalJSON()
					require.NoError(t, err)

					settingsObject := ingestsource.SettingsObject{
						SchemaID: "builtin:openpipeline.events.ingest-sources",
						Value:    valueBytes,
					}

					responseBytes, err := json.Marshal(settingsObject)
					require.NoError(t, err)

					return settings20.Response{
						Response: coreApi.Response{StatusCode: 200, Data: responseBytes},
						ID:       id,
						Items:    nil,
					}, nil
				},
			}
			service := ingestsource.ServiceImpl{GenericSettingsClient: client}
			value := ingestsource2.IngestSource{}
			err := service.Get(t.Context(), "objectID", &value)

			assert.NoError(t, err)
			assert.Equal(t, ingestsource2.IngestSource{
				Kind:          "events",
				DefaultBucket: &bucketName,
				DisplayName:   "displayName",
				Enabled:       false,
				PathSegment:   "some.path.segment",
				StaticRouting: &ingestsource2.PipelineReference{
					PipelineID:   &pipelineID,
					PipelineType: "custom",
				},
				Processing: &ingestsource2.Processing{
					Processors: []*processors.Processor{
						{
							Enabled:     true,
							Id:          "proc-2",
							Type:        processors.DqlProcessorType,
							Description: "my-proc-2",
							SampleData:  &sampleData,
							Matcher:     &matcher,
							Dql:         &processors.DqlAttributes{Script: "fieldsAdd true"},
						},
					},
				},
			}, value)
		})

		t.Run("Errors during get", func(t *testing.T) {
			client := &genericClientStub{
				get: func(ctx context.Context, id string) (settings20.Response, error) {
					return settings20.Response{Response: coreApi.Response{}}, assert.AnError
				},
			}
			service := ingestsource.ServiceImpl{GenericSettingsClient: client}
			err := service.Get(t.Context(), "objectID", &ingestsource2.IngestSource{})
			assert.Error(t, err)
			assert.ErrorIs(t, assert.AnError, err)
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Run("Returns stubs for all object IDs", func(t *testing.T) {
			var clients = make(map[string]ingestsource.SettingsClientForKind)

			clients["events"] = &clientForKindStub{
				list: func(ctx context.Context) (api.Stubs, error) {
					return api.Stubs{
						&api.Stub{ID: "id1"}, &api.Stub{ID: "id2"},
					}, nil
				}}
			clients["events.sdlc"] = &clientForKindStub{
				list: func(ctx context.Context) (api.Stubs, error) {
					return api.Stubs{
						&api.Stub{ID: "id3"},
					}, nil
				}}

			service := ingestsource.ServiceImpl{Credentials: &rest.Credentials{}, SettingsClientsPerKind: clients}
			stubs, _ := service.List(t.Context())
			assert.Len(t, stubs, 3)
			assert.Equal(t, &api.Stub{ID: "id1"}, stubs[0])
			assert.Equal(t, &api.Stub{ID: "id2"}, stubs[1])
			assert.Equal(t, &api.Stub{ID: "id3"}, stubs[2])
		})

		t.Run("Returns error if list fails", func(t *testing.T) {
			var clients = make(map[string]ingestsource.SettingsClientForKind)

			clients["events"] = &clientForKindStub{
				list: func(ctx context.Context) (api.Stubs, error) {
					return nil, assert.AnError
				}}

			service := ingestsource.ServiceImpl{Credentials: &rest.Credentials{}, SettingsClientsPerKind: clients}
			stubs, err := service.List(t.Context())
			assert.ErrorIs(t, err, assert.AnError)
			assert.Equal(t, api.Stubs{}, stubs)
		})

		t.Run("Returns error if getClientForKind fails", func(t *testing.T) {
			service := ingestsource.ServiceImpl{
				Credentials: &rest.Credentials{},
			}
			stubs, err := service.List(t.Context())
			assert.ErrorIs(t, err, rest.NoAPITokenError)
			assert.Equal(t, api.Stubs{}, stubs)
		})
	})

	t.Run("Create", func(t *testing.T) {
		t.Run("It successfully creates", func(t *testing.T) {
			var clients = make(map[string]ingestsource.SettingsClientForKind)
			clients["events"] = &clientForKindStub{
				create: func(ctx context.Context, v *ingestsource2.IngestSource) (*api.Stub, error) {
					return &api.Stub{ID: "object-id"}, nil
				},
			}

			service := ingestsource.ServiceImpl{SettingsClientsPerKind: clients}
			value := ingestsource2.IngestSource{
				Kind:        "events",
				DisplayName: "displayName",
				PathSegment: "my.path.segment",
				Enabled:     true,
			}

			resp, err := service.Create(t.Context(), &value)
			assert.NoError(t, err)
			assert.Equal(t, &api.Stub{ID: "object-id"}, resp)
		})

		t.Run("It errors during create", func(t *testing.T) {
			var clients = make(map[string]ingestsource.SettingsClientForKind)
			clients["events"] = &clientForKindStub{
				create: func(ctx context.Context, v *ingestsource2.IngestSource) (*api.Stub, error) {
					return nil, assert.AnError
				},
			}

			service := ingestsource.ServiceImpl{SettingsClientsPerKind: clients}
			value := ingestsource2.IngestSource{
				Kind:        "events",
				DisplayName: "displayName",
				PathSegment: "my.path.segment",
				Enabled:     true,
			}

			_, err := service.Create(t.Context(), &value)
			assert.ErrorContains(t, err, assert.AnError.Error())
		})

		t.Run("It errors during client creation", func(t *testing.T) {
			service := ingestsource.ServiceImpl{
				Credentials: &rest.Credentials{},
			}
			_, err := service.Create(t.Context(), &ingestsource2.IngestSource{})
			assert.ErrorIs(t, rest.NoAPITokenError, err)
		})
	})

	t.Run("Update", func(t *testing.T) {
		t.Run("It successfully updates", func(t *testing.T) {
			var clients = make(map[string]ingestsource.SettingsClientForKind)
			clients["events"] = &clientForKindStub{
				update: func(ctx context.Context, id string, v *ingestsource2.IngestSource) error {
					return nil
				},
			}

			service := ingestsource.ServiceImpl{SettingsClientsPerKind: clients}
			value := ingestsource2.IngestSource{
				Kind:        "events",
				DisplayName: "displayName",
				PathSegment: "my.path.segment",
				Enabled:     true,
			}

			err := service.Update(t.Context(), "object-id", &value)
			assert.NoError(t, err)
		})

		t.Run("It errors during update", func(t *testing.T) {
			var clients = make(map[string]ingestsource.SettingsClientForKind)
			clients["events"] = &clientForKindStub{
				update: func(ctx context.Context, id string, v *ingestsource2.IngestSource) error {
					return assert.AnError
				},
			}

			service := ingestsource.ServiceImpl{SettingsClientsPerKind: clients}
			value := ingestsource2.IngestSource{
				Kind:        "events",
				DisplayName: "displayName",
				PathSegment: "my.path.segment",
				Enabled:     true,
			}

			err := service.Update(t.Context(), "object-id", &value)
			assert.ErrorContains(t, err, assert.AnError.Error())
		})

		t.Run("It errors during client creation", func(t *testing.T) {
			service := ingestsource.ServiceImpl{
				Credentials: &rest.Credentials{},
			}
			err := service.Update(t.Context(), "", &ingestsource2.IngestSource{})
			assert.ErrorIs(t, rest.NoAPITokenError, err)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("It successfully deletes", func(t *testing.T) {
			var deleteCalled atomic.Int32
			client := &genericClientStub{
				delete: func(ctx context.Context, id string) (settings20.Response, error) {
					deleteCalled.Add(1)
					return settings20.Response{Response: coreApi.Response{StatusCode: 204}}, nil
				},
			}

			service := ingestsource.ServiceImpl{GenericSettingsClient: client}
			err := service.Delete(t.Context(), "objectID")
			assert.NoError(t, err)
			assert.Equal(t, int32(1), deleteCalled.Load())
		})
	})

	t.Run("Service returns a new instance", func(t *testing.T) {
		service := ingestsource.Service(nil)
		assert.IsType(t, &ingestsource.ServiceImpl{}, service)
	})

	t.Run("Returns the schema ID", func(t *testing.T) {
		service := ingestsource.ServiceImpl{}
		assert.Equal(t, "openpipelinev2:ingest-source", service.SchemaID())
	})
}
