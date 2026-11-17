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

package provider

import (
	"context"
	"sync/atomic"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// NewServer attaches the shared order only to refresh RPCs. A Read called from Create
// or Update must see the order after that write, even in the same provider process.
func NewServer() tfprotov5.ProviderServer {
	return &server{GRPCProviderServer: schema.NewGRPCProviderServer(Provider())}
}

type server struct {
	*schema.GRPCProviderServer
	order atomic.Pointer[settings20.RefreshOrder]
}

func (s *server) ConfigureProvider(ctx context.Context, req *tfprotov5.ConfigureProviderRequest) (*tfprotov5.ConfigureProviderResponse, error) {
	s.order.Store(new(settings20.RefreshOrder))
	return s.GRPCProviderServer.ConfigureProvider(ctx, req)
}

func (s *server) ReadResource(ctx context.Context, req *tfprotov5.ReadResourceRequest) (*tfprotov5.ReadResourceResponse, error) {
	if order := s.order.Load(); order != nil {
		ctx = order.WithContext(ctx)
	}
	return s.GRPCProviderServer.ReadResource(ctx, req)
}

func (s *server) ApplyResourceChange(ctx context.Context, req *tfprotov5.ApplyResourceChangeRequest) (*tfprotov5.ApplyResourceChangeResponse, error) {
	// A write can change another object's neighbors, including writes via dynatrace_generic_setting.
	s.order.Store(nil)
	return s.GRPCProviderServer.ApplyResourceChange(ctx, req)
}
