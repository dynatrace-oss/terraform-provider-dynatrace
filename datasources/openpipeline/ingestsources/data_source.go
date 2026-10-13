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

package ingestsources

import (
	bizevents "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/bizevents/ingestsources"
	davisevents "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/davis/events/ingestsources"
	davisproblems "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/davis/problems/ingestsources"
	events "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/events/ingestsources"
	eventssdlc "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/events/sdlc/ingestsources"
	eventssecurity "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/events/security/ingestsources"
	log "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/logs/ingestsources"
	metrics "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/metrics/ingestsources"
	secevents "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/security/events/ingestsources"
	spans "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/spans/ingestsources"
	systemevents "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/system/events/ingestsources"
	userevents "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/user/events/ingestsources"
	usersessions "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/openpipeline/usersessions/ingestsources"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// The need for the dataSource struct is the unit test
type dataSource struct {
	resource   *schema.Resource
	newService func(clientSet rest.ClientSet) (any, error)
}

func newDataSource[T settings.Settings](serviceFn export.ServiceFn[T]) dataSource {
	return dataSource{
		resource: genericDataSource(serviceFn),
		newService: func(clientSet rest.ClientSet) (any, error) {
			return serviceFn(clientSet)
		},
	}
}

func dataSources() map[string]dataSource {
	return map[string]dataSource{
		"dynatrace_openpipeline_v2_bizevents_ingestsources":       newDataSource(bizevents.Service),
		"dynatrace_openpipeline_v2_davis_events_ingestsources":    newDataSource(davisevents.Service),
		"dynatrace_openpipeline_v2_davis_problems_ingestsources":  newDataSource(davisproblems.Service),
		"dynatrace_openpipeline_v2_events_ingestsources":          newDataSource(events.Service),
		"dynatrace_openpipeline_v2_events_sdlc_ingestsources":     newDataSource(eventssdlc.Service),
		"dynatrace_openpipeline_v2_events_security_ingestsources": newDataSource(eventssecurity.Service),
		"dynatrace_openpipeline_v2_logs_ingestsources":            newDataSource(log.Service),
		"dynatrace_openpipeline_v2_metrics_ingestsources":         newDataSource(metrics.Service),
		"dynatrace_openpipeline_v2_security_events_ingestsources": newDataSource(secevents.Service),
		"dynatrace_openpipeline_v2_spans_ingestsources":           newDataSource(spans.Service),
		"dynatrace_openpipeline_v2_system_events_ingestsources":   newDataSource(systemevents.Service),
		"dynatrace_openpipeline_v2_user_events_ingestsources":     newDataSource(userevents.Service),
		"dynatrace_openpipeline_v2_usersessions_ingestsources":    newDataSource(usersessions.Service),
	}
}

func DataSources() map[string]*schema.Resource {
	sources := dataSources()
	resources := make(map[string]*schema.Resource, len(sources))
	for name, source := range sources {
		resources[name] = source.resource
	}
	return resources
}
