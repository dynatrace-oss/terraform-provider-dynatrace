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

package settings20

import (
	"context"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/envutils"
)

type neighbors struct {
	previous, next string
}

func (n neighbors) set(v settings.Settings) {
	if n.previous != "" {
		settings.SetInsertAfter(v, n.previous)
	}
	if n.next != "" {
		settings.SetInsertBefore(v, n.next)
	}
}

type objectOrder map[string]neighbors

func indexOrder(ids []string) objectOrder {
	order := make(objectOrder, len(ids))
	for i, id := range ids {
		if _, seen := order[id]; id == "" || seen {
			continue
		}
		var n neighbors
		if i > 0 {
			n.previous = ids[i-1]
		}
		if i+1 < len(ids) {
			n.next = ids[i+1]
		}
		order[id] = n
	}
	return order
}

type refreshOrderKey struct{}

// RefreshOrder shares a schema's object order across Terraform ReadResource calls.
// Each configured provider gets a new instance. Writes and their readbacks bypass it.
type RefreshOrder struct {
	schemas sync.Map // schema ID -> *orderResult
}

func (r *RefreshOrder) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, refreshOrderKey{}, r)
}

type orderResult struct {
	ready chan struct{}
	order objectOrder
	err   error
}

func orderForRead(ctx context.Context, schemaID string, list func(context.Context) (objectOrder, error)) (objectOrder, error) {
	refresh, _ := ctx.Value(refreshOrderKey{}).(*RefreshOrder)
	// An explicit authentication override can see a different set of objects.
	if refresh == nil || ctx.Value(envutils.DynatraceHTTPOAuthPreference.Key) != nil {
		return list(ctx)
	}

	pending := &orderResult{ready: make(chan struct{})}
	value, loaded := refresh.schemas.LoadOrStore(schemaID, pending)
	if loaded {
		result := value.(*orderResult)
		select {
		case <-result.ready:
			return result.order, result.err
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	pending.order, pending.err = list(ctx)
	if pending.err != nil {
		refresh.schemas.Delete(schemaID)
	}
	close(pending.ready)
	return pending.order, pending.err
}
