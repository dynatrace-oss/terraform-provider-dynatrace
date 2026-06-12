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

// Package settings models the Dynatrace GCP monitoring configuration as a
// fully typed Terraform resource. Internally it is delivered through the
// generic Extensions 2.0 Monitoring Configuration API (extension
// "com.dynatrace.extension.da-gcp"), but consumers see only first-class
// attributes — there is intentionally no JSON escape hatch.
package settings

import (
	"encoding/json"
	"sort"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GCPExtensionName is the fully qualified extension name that backs every GCP
// monitoring configuration.
const GCPExtensionName = "com.dynatrace.extension.da-gcp"

// DefaultScope is the Settings 2.0 scope used by GCP DAC monitoring configs.
const DefaultScope = "integration-gcp"

// Wire-level defaults that match dtctl's `create gcp monitoring` behavior.
const (
	DefaultActivationContext = "DATA_ACQUISITION"
)

// Settings is the typed model exposed to Terraform.
type Settings struct {
	Name                       string
	Enabled                    bool
	ExtensionVersion           string
	Scope                      string
	ActivationContext          string
	FeatureSets                []string
	Credentials                Credentials
	Regions                    []string
	ProjectFilter              []string
	FolderFilter               []string
	TagFilters                 TagFilters
	LabelFilters               TagFilters
	TagEnrichment              []string
	LabelEnrichment            []string
	SmartscapeEnabled          bool
	ObservabilityScopesEnabled bool
	ResourceAutodiscovery      ResourceAutodiscoveries
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Human-readable name of the monitoring configuration (written to `description` in the extension config payload).",
			Required:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Whether the monitoring configuration is active. Defaults to true.",
			Optional:    true,
			Default:     true,
		},
		"extension_version": {
			Type:        schema.TypeString,
			Description: "Version of `com.dynatrace.extension.da-gcp` that this configuration targets. Optional — when omitted at create time, the provider picks the highest semver version installed on the tenant (same behavior as `dtctl create gcp monitoring`). The resolved value is persisted to state. On subsequent refreshes the provider reads back whatever version Dynatrace currently reports for this configuration; if the extension was auto-updated (or bumped manually) the new version surfaces as drift in `terraform plan`, but no Terraform-driven update silently re-resolves it. To pin a version, set it explicitly here.",
			Optional:    true,
			Computed:    true,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "Settings 2.0 scope. Defaults to `integration-gcp`. Changing it forces recreation.",
			Optional:    true,
			Default:     DefaultScope,
			ForceNew:    true,
		},
		"activation_context": {
			Type:        schema.TypeString,
			Description: "Extension activation context. Defaults to `DATA_ACQUISITION`.",
			Optional:    true,
			Default:     DefaultActivationContext,
		},
		"feature_sets": {
			Type:        schema.TypeSet,
			Description: "GCP feature sets to enable (e.g. `compute_engine_essential`, `kubernetes_engine_essential`, `sql_essential`). When empty, the extension defaults are used.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"credential": {
			Type:        schema.TypeList,
			Description: "HAS connection + GCP service-account binding. At least one is required. dtctl always writes exactly one, but the API accepts a list.",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(Credential).Schema()},
		},
		"regions": {
			Type:        schema.TypeSet,
			Description: "GCP regions (locations) to monitor, e.g. `us-central1`. Empty set = all locations the extension knows about. Maps to `locationFiltering` on the wire.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"project_filter": {
			Type:        schema.TypeSet,
			Description: "GCP project IDs the customer service account reads. Empty set means \"all projects the SA can impersonate into\". Maps to `projectFiltering` on the wire.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"folder_filter": {
			Type:        schema.TypeSet,
			Description: "GCP folder IDs — fans out to all projects under each folder. Maps to `folderFiltering` on the wire.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"tag_filter": {
			Type:        schema.TypeList,
			Description: "Filter monitored resources by GCP resource-manager tag (`tagKeys/…`). Repeat the block to define multiple filters.",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(TagFilter).Schema()},
		},
		"label_filter": {
			Type:        schema.TypeList,
			Description: "Filter monitored resources by GCP label (classic key/value labels per resource). Distinct from `tag_filter`. Repeat the block to define multiple filters.",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(TagFilter).Schema()},
		},
		"tag_enrichment": {
			Type:        schema.TypeSet,
			Description: "GCP tag keys whose values are copied as Dynatrace tags on monitored entities.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"label_enrichment": {
			Type:        schema.TypeSet,
			Description: "GCP label keys whose values are copied as Dynatrace tags on monitored entities.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"smartscape_enabled": {
			Type:        schema.TypeBool,
			Description: "Whether Smartscape topology mapping is enabled. Defaults to true.",
			Optional:    true,
			Default:     true,
		},
		"observability_scopes_enabled": {
			Type:        schema.TypeBool,
			Description: "Whether observability scopes are enabled. Defaults to false.",
			Optional:    true,
			Default:     false,
		},
		"resource_autodiscovery": {
			Type:        schema.TypeList,
			Description: "Per-resource-type autodiscovery override. Repeat the block once per `resource_type` you want to override.",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(ResourceAutodiscovery).Schema()},
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	featureSets := append([]string(nil), me.FeatureSets...)
	sort.Strings(featureSets)
	regions := append([]string(nil), me.Regions...)
	sort.Strings(regions)
	projects := append([]string(nil), me.ProjectFilter...)
	sort.Strings(projects)
	folders := append([]string(nil), me.FolderFilter...)
	sort.Strings(folders)
	tagEnrichment := append([]string(nil), me.TagEnrichment...)
	sort.Strings(tagEnrichment)
	labelEnrichment := append([]string(nil), me.LabelEnrichment...)
	sort.Strings(labelEnrichment)

	if err := properties.EncodeAll(map[string]any{
		"name":                         me.Name,
		"enabled":                      me.Enabled,
		"extension_version":            me.ExtensionVersion,
		"scope":                        me.Scope,
		"activation_context":           me.ActivationContext,
		"feature_sets":                 featureSets,
		"regions":                      regions,
		"project_filter":               projects,
		"folder_filter":                folders,
		"tag_enrichment":               tagEnrichment,
		"label_enrichment":             labelEnrichment,
		"smartscape_enabled":           me.SmartscapeEnabled,
		"observability_scopes_enabled": me.ObservabilityScopesEnabled,
	}); err != nil {
		return err
	}
	if err := properties.EncodeSlice("credential", me.Credentials); err != nil {
		return err
	}
	if err := properties.EncodeSlice("tag_filter", me.TagFilters); err != nil {
		return err
	}
	if err := properties.EncodeSlice("label_filter", me.LabelFilters); err != nil {
		return err
	}
	return properties.EncodeSlice("resource_autodiscovery", me.ResourceAutodiscovery)
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.DecodeAll(map[string]any{
		"name":                         &me.Name,
		"enabled":                      &me.Enabled,
		"extension_version":            &me.ExtensionVersion,
		"scope":                        &me.Scope,
		"activation_context":           &me.ActivationContext,
		"feature_sets":                 &me.FeatureSets,
		"regions":                      &me.Regions,
		"project_filter":               &me.ProjectFilter,
		"folder_filter":                &me.FolderFilter,
		"tag_enrichment":               &me.TagEnrichment,
		"label_enrichment":             &me.LabelEnrichment,
		"smartscape_enabled":           &me.SmartscapeEnabled,
		"observability_scopes_enabled": &me.ObservabilityScopesEnabled,
	}); err != nil {
		return err
	}
	if err := decoder.DecodeSlice("credential", &me.Credentials); err != nil {
		return err
	}
	if err := decoder.DecodeSlice("tag_filter", &me.TagFilters); err != nil {
		return err
	}
	if err := decoder.DecodeSlice("label_filter", &me.LabelFilters); err != nil {
		return err
	}
	if err := decoder.DecodeSlice("resource_autodiscovery", &me.ResourceAutodiscovery); err != nil {
		return err
	}
	me.applyDefaults()
	return nil
}

// applyDefaults backfills attributes that the user did not set so the wire
// payload always carries a complete, server-accepted configuration.
func (me *Settings) applyDefaults() {
	if me.Scope == "" {
		me.Scope = DefaultScope
	}
	if me.ActivationContext == "" {
		me.ActivationContext = DefaultActivationContext
	}
	for _, c := range me.Credentials {
		if c == nil {
			continue
		}
		if c.Description == "" {
			c.Description = me.Name
		}
	}
}

// wirePayload is the on-the-wire shape accepted by
// POST /platform/extensions/v2/extensions/{name}/monitoringConfigurations.
type wirePayload struct {
	Scope string         `json:"scope"`
	Value map[string]any `json:"value"`
}

// MarshalJSON renders the typed Settings into the JSON body expected by the
// Extensions 2.0 Monitoring Configuration endpoint.
func (me *Settings) MarshalJSON() ([]byte, error) {
	me.applyDefaults()

	creds := []any{}
	for _, c := range me.Credentials {
		if c == nil {
			continue
		}
		creds = append(creds, map[string]any{
			"connectionId":   c.ConnectionID,
			"serviceAccount": c.ServiceAccount,
			"description":    c.Description,
			"enabled":        c.Enabled,
		})
	}

	tagFiltering := []any{}
	for _, f := range me.TagFilters {
		tagFiltering = append(tagFiltering, map[string]any{
			"key":       f.Key,
			"value":     f.Value,
			"condition": f.Condition,
		})
	}

	labelFiltering := []any{}
	for _, f := range me.LabelFilters {
		labelFiltering = append(labelFiltering, map[string]any{
			"key":       f.Key,
			"value":     f.Value,
			"condition": f.Condition,
		})
	}

	resources := []any{}
	for _, r := range me.ResourceAutodiscovery {
		if r == nil {
			continue
		}
		entry := map[string]any{
			"resourceType":         r.ResourceType,
			"autoDiscoveryEnabled": r.AutoDiscoveryEnabled,
		}
		if len(r.ExcludeMetricType) > 0 {
			entry["autodiscoveryExcludeMetricType"] = toAnySlice(r.ExcludeMetricType)
		}
		resources = append(resources, entry)
	}

	googleCloud := map[string]any{
		"credentials":             creds,
		"locationFiltering":       toAnySlice(me.Regions),
		"projectFiltering":        toAnySlice(me.ProjectFilter),
		"folderFiltering":         toAnySlice(me.FolderFilter),
		"tagFiltering":            tagFiltering,
		"labelFiltering":          labelFiltering,
		"tagEnrichment":           toAnySlice(me.TagEnrichment),
		"labelEnrichment":         toAnySlice(me.LabelEnrichment),
		"smartscapeConfiguration": map[string]any{"enabled": me.SmartscapeEnabled},
		"resources":               resources,
	}
	if me.ObservabilityScopesEnabled {
		googleCloud["observabilityScopesEnabled"] = true
	}

	value := map[string]any{
		"enabled":           me.Enabled,
		"description":       me.Name,
		"version":           me.ExtensionVersion,
		"activationContext": me.ActivationContext,
		"googleCloud":       googleCloud,
	}
	if len(me.FeatureSets) > 0 {
		value["featureSets"] = me.FeatureSets
	}

	return json.Marshal(wirePayload{Scope: me.Scope, Value: value})
}

// UnmarshalJSON parses the API response (same shape as the request payload)
// back into the typed Settings.
func (me *Settings) UnmarshalJSON(data []byte) error {
	var wire wirePayload
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	me.Scope = wire.Scope
	if wire.Value == nil {
		me.applyDefaults()
		return nil
	}

	if v, ok := wire.Value["description"].(string); ok {
		me.Name = v
	}
	if v, ok := wire.Value["enabled"].(bool); ok {
		me.Enabled = v
	}
	if v, ok := wire.Value["version"].(string); ok {
		me.ExtensionVersion = v
	}
	if v, ok := wire.Value["activationContext"].(string); ok {
		me.ActivationContext = v
	}
	if fs, ok := wire.Value["featureSets"].([]any); ok && len(fs) > 0 {
		me.FeatureSets = me.FeatureSets[:0]
		for _, x := range fs {
			if s, ok := x.(string); ok {
				me.FeatureSets = append(me.FeatureSets, s)
			}
		}
	} else {
		me.FeatureSets = nil
	}

	gc, _ := wire.Value["googleCloud"].(map[string]any)
	if gc == nil {
		me.applyDefaults()
		return nil
	}

	if creds, ok := gc["credentials"].([]any); ok {
		me.Credentials = me.Credentials[:0]
		for _, x := range creds {
			m, ok := x.(map[string]any)
			if !ok {
				continue
			}
			c := &Credential{Enabled: true}
			if v, ok := m["connectionId"].(string); ok {
				c.ConnectionID = v
			}
			if v, ok := m["serviceAccount"].(string); ok {
				c.ServiceAccount = v
			}
			if v, ok := m["description"].(string); ok {
				c.Description = v
			}
			if v, ok := m["enabled"].(bool); ok {
				c.Enabled = v
			}
			me.Credentials = append(me.Credentials, c)
		}
	}

	me.Regions = readStringList(gc["locationFiltering"])
	me.ProjectFilter = readNonEmptyStringList(gc["projectFiltering"])
	me.FolderFilter = readNonEmptyStringList(gc["folderFiltering"])
	me.TagEnrichment = readNonEmptyStringList(gc["tagEnrichment"])
	me.LabelEnrichment = readNonEmptyStringList(gc["labelEnrichment"])

	if tf, ok := gc["tagFiltering"].([]any); ok && len(tf) > 0 {
		me.TagFilters = decodeTagFilters(tf)
	} else {
		me.TagFilters = nil
	}
	if lf, ok := gc["labelFiltering"].([]any); ok && len(lf) > 0 {
		me.LabelFilters = decodeTagFilters(lf)
	} else {
		me.LabelFilters = nil
	}

	if res, ok := gc["resources"].([]any); ok && len(res) > 0 {
		me.ResourceAutodiscovery = me.ResourceAutodiscovery[:0]
		for _, x := range res {
			m, ok := x.(map[string]any)
			if !ok {
				continue
			}
			r := &ResourceAutodiscovery{}
			if v, ok := m["resourceType"].(string); ok {
				r.ResourceType = v
			}
			if v, ok := m["autoDiscoveryEnabled"].(bool); ok {
				r.AutoDiscoveryEnabled = v
			}
			if exc, ok := m["autodiscoveryExcludeMetricType"].([]any); ok {
				for _, s := range exc {
					if sv, ok := s.(string); ok {
						r.ExcludeMetricType = append(r.ExcludeMetricType, sv)
					}
				}
			}
			me.ResourceAutodiscovery = append(me.ResourceAutodiscovery, r)
		}
	} else {
		me.ResourceAutodiscovery = nil
	}

	// smartscapeConfiguration: default to true if absent, otherwise read enabled.
	me.SmartscapeEnabled = true
	if sc, ok := gc["smartscapeConfiguration"].(map[string]any); ok {
		if v, ok := sc["enabled"].(bool); ok {
			me.SmartscapeEnabled = v
		}
	}

	if v, ok := gc["observabilityScopesEnabled"].(bool); ok {
		me.ObservabilityScopesEnabled = v
	} else {
		me.ObservabilityScopesEnabled = false
	}

	// `featureSetConfiguration[]` is deliberately ignored — it is an API-echo
	// array (always empty in real configs) and surfacing it would produce
	// eternal plan drift. See spec §5 gotcha #3.

	me.applyDefaults()
	return nil
}

func decodeTagFilters(in []any) TagFilters {
	out := TagFilters{}
	for _, x := range in {
		m, ok := x.(map[string]any)
		if !ok {
			continue
		}
		f := &TagFilter{}
		if v, ok := m["key"].(string); ok {
			f.Key = v
		}
		if v, ok := m["value"].(string); ok {
			f.Value = v
		}
		if v, ok := m["condition"].(string); ok {
			f.Condition = v
		}
		out = append(out, f)
	}
	return out
}

// readStringList accepts a []any of strings and returns nil for missing/empty.
func readStringList(in any) []string {
	arr, ok := in.([]any)
	if !ok || len(arr) == 0 {
		return nil
	}
	out := make([]string, 0, len(arr))
	for _, x := range arr {
		if s, ok := x.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

// readNonEmptyStringList behaves like readStringList — kept as a distinct name
// for readability where the API-echo (empty []) trap matters.
func readNonEmptyStringList(in any) []string {
	return readStringList(in)
}

func toAnySlice(in []string) []any {
	out := make([]any, 0, len(in))
	for _, s := range in {
		out = append(out, s)
	}
	return out
}
