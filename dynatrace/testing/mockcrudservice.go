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

// MockCRUDService implements settings.CRUDService for testing
type MockCRUDService[T settings.Settings] struct {
	CreateFunc func(ctx context.Context, v T) (*api.Stub, error)
	GetFunc    func(ctx context.Context, id string, v T) error
	UpdateFunc func(ctx context.Context, id string, v T) error
	ListFunc   func(ctx context.Context) (api.Stubs, error)
	DeleteFunc func(ctx context.Context, id string) error
}

func (m *MockCRUDService[T]) Create(ctx context.Context, v T) (*api.Stub, error) {
	return m.CreateFunc(ctx, v)
}

func (m *MockCRUDService[T]) Get(ctx context.Context, id string, v T) error {
	return m.GetFunc(ctx, id, v)
}

func (m *MockCRUDService[T]) Update(ctx context.Context, id string, v T) error {
	return m.UpdateFunc(ctx, id, v)
}

func (m *MockCRUDService[T]) List(ctx context.Context) (api.Stubs, error) {
	return m.ListFunc(ctx)
}

func (m *MockCRUDService[T]) Delete(ctx context.Context, id string) error {
	return m.DeleteFunc(ctx, id)
}

func (m *MockCRUDService[T]) SchemaID() string { return "" }
