//go:build unit

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

package monitors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors"
)

var (
	regularMonitor = &monitors.MonitorCollectionElement{
		Name:     "My HTTP monitor",
		EntityID: "HTTP_CHECK-0000000000000001",
		Type:     monitors.Types.HTTP,
		Enabled:  true,
	}
	synchronizedMonitor = &monitors.MonitorCollectionElement{
		Name:     "Monitor synchronizing credentials with CyberArk Vault (CREDENTIALS_VAULT-1234567890ABCDEF)",
		EntityID: "HTTP_CHECK-0000000000000002",
		Type:     monitors.Types.HTTP,
		Enabled:  true,
	}
)

func toStub(m *monitors.MonitorCollectionElement) *api.Stub {
	return &api.Stub{ID: m.EntityID, Name: m.Name}
}

func TestMonitors_ToStubs(t *testing.T) {
	tests := []struct {
		name     string
		monitors []*monitors.MonitorCollectionElement
		expected api.Stubs
	}{
		{
			name:     "nil list returns empty stubs",
			monitors: nil,
			expected: api.Stubs{},
		},
		{
			name:     "only regular monitors are kept",
			monitors: []*monitors.MonitorCollectionElement{regularMonitor},
			expected: api.Stubs{toStub(regularMonitor)},
		},
		{
			name:     "synchronized monitors are filtered out",
			monitors: []*monitors.MonitorCollectionElement{synchronizedMonitor},
			expected: api.Stubs{},
		},
		{
			name:     "mixed list keeps only regular monitors",
			monitors: []*monitors.MonitorCollectionElement{synchronizedMonitor, regularMonitor},
			expected: api.Stubs{toStub(regularMonitor)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &monitors.Monitors{Monitors: tt.monitors}
			assert.Equal(t, tt.expected, m.ToStubs())
		})
	}
}

func TestSynchronizedMonitors_ToStubs(t *testing.T) {
	tests := []struct {
		name     string
		monitors []*monitors.MonitorCollectionElement
		expected api.Stubs
	}{
		{
			name:     "nil list returns empty stubs",
			monitors: nil,
			expected: api.Stubs{},
		},
		{
			name:     "regular monitors are filtered out",
			monitors: []*monitors.MonitorCollectionElement{regularMonitor},
			expected: api.Stubs{},
		},
		{
			name:     "only synchronized monitors are kept",
			monitors: []*monitors.MonitorCollectionElement{synchronizedMonitor},
			expected: api.Stubs{toStub(synchronizedMonitor)},
		},
		{
			name:     "mixed list keeps only synchronized monitors",
			monitors: []*monitors.MonitorCollectionElement{synchronizedMonitor, regularMonitor},
			expected: api.Stubs{toStub(synchronizedMonitor)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &monitors.SynchronizedMonitors{Monitors: tt.monitors}
			assert.Equal(t, tt.expected, m.ToStubs())
		})
	}
}
