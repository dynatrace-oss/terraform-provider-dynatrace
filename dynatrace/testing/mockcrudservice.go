//go:build unit

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

package testing

import (
	"context"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

// MockCRUDService implements settings.CRUDService for testing.
type MockCRUDService[T settings.Settings] struct {
	ListFn        func(ctx context.Context) (api.Stubs, error)
	GetFn         func(ctx context.Context, id string, v T) error
	CreateFn      func(ctx context.Context, v T) (*api.Stub, error)
	UpdateFn      func(ctx context.Context, id string, v T) error
	DeleteFn      func(ctx context.Context, id string) error
	SchemaIDValue string
}

func (m *MockCRUDService[T]) List(ctx context.Context) (api.Stubs, error) {
	return m.ListFn(ctx)
}

func (m *MockCRUDService[T]) Get(ctx context.Context, id string, v T) error {
	return m.GetFn(ctx, id, v)
}

func (m *MockCRUDService[T]) SchemaID() string {
	return m.SchemaIDValue
}

func (m *MockCRUDService[T]) Create(ctx context.Context, v T) (*api.Stub, error) {
	return m.CreateFn(ctx, v)
}

func (m *MockCRUDService[T]) Update(ctx context.Context, id string, v T) error {
	return m.UpdateFn(ctx, id, v)
}

func (m *MockCRUDService[T]) Delete(ctx context.Context, id string) error {
	return m.DeleteFn(ctx, id)
}

// MockListFilterService implements settings.CRUDService and settings.ListFilterService for testing.
type MockListFilterService[T settings.Settings] struct {
	MockCRUDService[T]
	ListWithFilterFn func(ctx context.Context, filter string) (api.Stubs, error)
}

func (m *MockListFilterService[T]) ListWithFilter(ctx context.Context, filter string) (api.Stubs, error) {
	return m.ListWithFilterFn(ctx, filter)
}
