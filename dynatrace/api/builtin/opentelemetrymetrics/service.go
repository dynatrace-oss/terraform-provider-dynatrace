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

package opentelemetrymetrics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	opentelemetrymetrics "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/opentelemetrymetrics/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/httpcache"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
)

const SchemaVersion = "1.3"
const SchemaID = "builtin:opentelemetry-metrics"

func Service(credentials *settings.Credentials) settings.CRUDService[*opentelemetrymetrics.Settings] {
	return &service{
		service: settings20.Service[*opentelemetrymetrics.Settings](credentials, SchemaID, SchemaVersion),
		client:  httpcache.DefaultClient(credentials.URL, credentials.Token, SchemaID),
	}
}

type service struct {
	service settings.CRUDService[*opentelemetrymetrics.Settings]
	client  rest.Client
}

var mu sync.Mutex

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return me.service.List(ctx)
}

func (me *service) Get(ctx context.Context, id string, v *opentelemetrymetrics.Settings) error {
	mu.Lock()
	defer mu.Unlock()
	stateConfig := getStateConfig(ctx)
	v.Mode = resolveMode(v, stateConfig)
	existingValue, err := me.fetchExistingRecord(ctx)
	if err != nil {
		return err
	}
	if existingValue == nil {
		return rest.Error{Code: 404, Message: "No configuration found"}
	}
	switch v.Mode {
	case opentelemetrymetrics.Modes.Explicit:
		*v = *existingValue
		v.Mode = opentelemetrymetrics.Modes.Explicit
	case opentelemetrymetrics.Modes.Additive:
		if stateConfig == nil {
			*v = *existingValue
			v.Mode = opentelemetrymetrics.Modes.Explicit
		} else {
			// Attributes that aren't maintained within the State need to be excluded
			existingValue.AdditionalAttributes = intersect(existingValue.AdditionalAttributes, stateConfig.AdditionalAttributes)
			existingValue.ToDropAttributes = intersect(existingValue.ToDropAttributes, stateConfig.ToDropAttributes)
			*v = *existingValue
			v.Mode = opentelemetrymetrics.Modes.Additive
		}
	}
	v.Scope = opt.NewString("environment")
	return nil
}

func setFlags(target, existingValue *opentelemetrymetrics.Settings) {
	if existingValue == nil {
		if target.AdditionalAttributesToDimensionEnabled == nil {
			target.AdditionalAttributesToDimensionEnabled = opt.NewBool(true)
		}
		if target.MeterNameToDimensionEnabled == nil {
			target.MeterNameToDimensionEnabled = opt.NewBool(true)
		}
	} else {
		if target.AdditionalAttributesToDimensionEnabled == nil {
			target.AdditionalAttributesToDimensionEnabled = existingValue.AdditionalAttributesToDimensionEnabled
		}
		if target.AdditionalAttributesToDimensionEnabled == nil {
			target.AdditionalAttributesToDimensionEnabled = opt.NewBool(true)
		}
		if target.MeterNameToDimensionEnabled == nil {
			target.MeterNameToDimensionEnabled = existingValue.MeterNameToDimensionEnabled
		}
		if target.MeterNameToDimensionEnabled == nil {
			target.MeterNameToDimensionEnabled = opt.NewBool(true)
		}
	}
}

func toJSON(v any) string {
	data, _ := json.Marshal(v)
	return string(data)
}

func (me *service) Create(ctx context.Context, v *opentelemetrymetrics.Settings) (*api.Stub, error) {
	mu.Lock()
	defer mu.Unlock()
	v.Mode = resolveMode(v, nil)
	existingValue, err := me.fetchExistingRecord(ctx)
	if err != nil {
		return nil, err
	}
	effectiveValue := *v
	setFlags(&effectiveValue, existingValue)
	switch v.Mode {
	case opentelemetrymetrics.Modes.Explicit:
		// nothing to do - the list of attributes we want to apply is exactly what should end up at the remote side
	case opentelemetrymetrics.Modes.Additive:
		if existingValue != nil {
			// Configured Attributes will either REPLACE the existing Attributes or get appended to the end
			effectiveValue.AdditionalAttributes = merge(effectiveValue.AdditionalAttributes, existingValue.AdditionalAttributes)
			effectiveValue.ToDropAttributes = merge(effectiveValue.ToDropAttributes, existingValue.ToDropAttributes)
		}
	}

	stub, err := me.service.Create(ctx, &effectiveValue)
	if err != nil {
		return stub, err
	}
	v.AdditionalAttributesToDimensionEnabled = effectiveValue.AdditionalAttributesToDimensionEnabled
	v.MeterNameToDimensionEnabled = effectiveValue.MeterNameToDimensionEnabled
	v.Scope = opt.NewString("environment")
	return stub, nil
}

func (me *service) Update(ctx context.Context, id string, v *opentelemetrymetrics.Settings) error {
	mu.Lock()
	defer mu.Unlock()
	stateConfig := getStateConfig(ctx)
	v.Mode = resolveMode(v, stateConfig)
	existingValue, err := me.fetchExistingRecord(ctx)
	if err != nil {
		return err
	}
	effectiveValue := *v
	setFlags(&effectiveValue, existingValue)
	switch v.Mode {
	case opentelemetrymetrics.Modes.Explicit:
		// nothing to do - the list of attributes we want to apply is exactly what should end up at the remote side
	case opentelemetrymetrics.Modes.Additive:
		if existingValue != nil {
			// Every Attribute found within the State but NOT within the Config
			// needs to get removed from the Attributes on the Existing Value
			var addAttribsToDelete []*opentelemetrymetrics.AdditionalAttributeItem
			var dropAttribsToDelete []*opentelemetrymetrics.DropAttributeItem
			if stateConfig != nil {
				addAttribsToDelete = remove(stateConfig.AdditionalAttributes, effectiveValue.AdditionalAttributes)
				dropAttribsToDelete = remove(stateConfig.ToDropAttributes, effectiveValue.ToDropAttributes)
			}
			existingValue.AdditionalAttributes = remove(existingValue.AdditionalAttributes, addAttribsToDelete)
			// Configured Attributes will either REPLACE the existing Attributes or get appended to the end
			effectiveValue.AdditionalAttributes = merge(effectiveValue.AdditionalAttributes, existingValue.AdditionalAttributes)

			existingValue.ToDropAttributes = remove(existingValue.ToDropAttributes, dropAttribsToDelete)
			// Configured Attributes will either REPLACE the existing Attributes or get appended to the end
			effectiveValue.ToDropAttributes = merge(effectiveValue.ToDropAttributes, existingValue.ToDropAttributes)
		}

	}
	if err := me.service.Update(ctx, id, &effectiveValue); err != nil {
		return err
	}
	v.AdditionalAttributesToDimensionEnabled = effectiveValue.AdditionalAttributesToDimensionEnabled
	v.MeterNameToDimensionEnabled = effectiveValue.MeterNameToDimensionEnabled
	v.Scope = opt.NewString("environment")
	return nil
}

func (me *service) Delete(ctx context.Context, id string) error {
	mu.Lock()
	defer mu.Unlock()
	stateConfig := getStateConfig(ctx)
	mode := resolveMode(stateConfig, nil)

	existingValue, err := me.fetchExistingRecord(ctx)
	if err != nil {
		return err
	}
	// if someone has deleted our settings already, nothing to do
	if existingValue == nil {
		return nil
	}
	effectiveValue := *existingValue
	if mode == opentelemetrymetrics.Modes.Additive {
		if stateConfig != nil {
			effectiveValue.AdditionalAttributes = remove(effectiveValue.AdditionalAttributes, stateConfig.AdditionalAttributes)
			effectiveValue.ToDropAttributes = remove(effectiveValue.ToDropAttributes, stateConfig.ToDropAttributes)
		}
		// if there are no attributes left, we opt for deleting the whole setting
		if len(effectiveValue.AdditionalAttributes) != 0 || len(effectiveValue.ToDropAttributes) != 0 {
			return me.service.Update(ctx, id, &effectiveValue)
		}
	}

	// if mode is "Explicit" or if there are still attributes remaining
	return me.service.Delete(ctx, id)
}

func (me *service) SchemaID() string {
	return me.service.SchemaID()
}

func getKey(v any) string {
	return reflect.ValueOf(v).Elem().FieldByName("AttributeKey").Interface().(string)
}

func equals(a, b any) bool {
	if a == nil && b != nil {
		return false
	}
	if a != nil && b == nil {
		return false
	}
	if a == nil && b == nil {
		return true
	}
	ka := reflect.ValueOf(a).Elem().FieldByName("AttributeKey").Interface().(string)
	kb := reflect.ValueOf(b).Elem().FieldByName("AttributeKey").Interface().(string)
	return ka == kb
}

func remove[T any](source, elemsToRemove []*T) []*T {
	if len(source) == 0 && len(elemsToRemove) == 0 {
		return source
	}
	result := []*T{}
	for _, elem := range source {
		found := false
		for _, removeElem := range elemsToRemove {
			if equals(elem, removeElem) {
				found = true
				break
			}
		}
		if !found {
			result = append(result, elem)
		}
	}
	return result
}

func intersect[T any](a, b []*T) []*T {
	return remove(a, remove(a, b))
}

func merge[T any](a, b []*T) []*T {
	if len(a) == 0 {
		return b
	}
	if len(b) == 0 {
		return a
	}

	result := append([]*T{}, b...)

	for _, configuredElem := range a {
		found := false
		for idx, existingElem := range result {
			if equals(configuredElem, existingElem) {
				result[idx] = configuredElem
				found = true
				break
			}
		}

		if !found {
			result = append(result, configuredElem)
		}
	}
	return result
}

func getStateConfig(ctx context.Context) *opentelemetrymetrics.Settings {
	var stateConfig *opentelemetrymetrics.Settings
	cfg := ctx.Value(settings.ContextKeyStateConfig)
	if typedConfig, ok := cfg.(*opentelemetrymetrics.Settings); ok {
		stateConfig = typedConfig
	}
	return stateConfig
}

func resolveMode(v, stateConfig *opentelemetrymetrics.Settings) opentelemetrymetrics.Mode {
	if len(v.Mode) != 0 {
		return v.Mode
	}
	if stateConfig == nil || len(stateConfig.Mode) == 0 {
		return opentelemetrymetrics.Modes.Explicit
	}
	return stateConfig.Mode
}

func (me *service) fetchExistingRecord(ctx context.Context) (*opentelemetrymetrics.Settings, error) {
	var err error
	var sol settings20.SettingsObjectList
	req := me.client.Get(ctx, fmt.Sprintf("/api/v2/settings/objects?schemaIds=%s&fields=%s&pageSize=1", url.QueryEscape(me.SchemaID()), url.QueryEscape("objectId,value,scope,schemaVersion")), 200)
	if err = req.Finish(&sol); err != nil {
		return nil, err
	}
	var existingValue *opentelemetrymetrics.Settings
	if len(sol.Items) > 0 {
		if err := json.Unmarshal(sol.Items[0].Value, &existingValue); err != nil {
			return nil, err
		}
	}
	return existingValue, err
}
