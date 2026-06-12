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

// Package settings models the Dynatrace Azure monitoring configuration as a
// fully typed Terraform resource. Internally it is delivered through the
// generic Extensions 2.0 Monitoring Configuration API (extension
// "com.dynatrace.extension.da-azure"), but consumers see only first-class
// attributes — there is intentionally no JSON escape hatch.
package settings

import (
	"encoding/json"
	"sort"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AzureExtensionName is the fully qualified extension name that backs every
// Azure monitoring configuration.
const AzureExtensionName = "com.dynatrace.extension.da-azure"

// DefaultScope is the Settings 2.0 scope used by Azure DAC monitoring configs.
const DefaultScope = "integration-azure"

// Wire-level defaults that match dtctl's `create azure` behavior.
const (
	DefaultActivationContext         = "DATA_ACQUISITION"
	DefaultConfigurationMode         = "ADVANCED"
	DefaultDeploymentMode            = "AUTOMATED"
	DefaultDeploymentScope           = "SUBSCRIPTION"
	DefaultSubscriptionFilteringMode = "INCLUDE"
	DefaultCredentialType            = "FEDERATED"
)

// Settings is the typed model exposed to Terraform.
type Settings struct {
	Name                      string
	Enabled                   bool
	ExtensionVersion          string
	Scope                     string
	ActivationContext         string
	FeatureSets               []string
	Credentials               Credentials
	Regions                   []string
	SubscriptionFilter        []string
	SubscriptionFilteringMode string
	ConfigurationMode         string
	DeploymentMode            string
	DeploymentScope           string
	TagFilters                TagFilters
	TagEnrichment             []string
	DTLabelEnrichments        DTLabelEnrichments
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
			Description: "Version of `com.dynatrace.extension.da-azure` that this configuration targets. Optional — when omitted at create time, the provider picks the highest semver version installed on the tenant (same behavior as `dtctl create azure monitoring`). The resolved value is persisted to state. On subsequent refreshes the provider reads back whatever version Dynatrace currently reports for this configuration; if the extension was auto-updated (or bumped manually) the new version surfaces as drift in `terraform plan`, but no Terraform-driven update silently re-resolves it. To pin a version, set it explicitly here.",
			Optional:    true,
			Computed:    true,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "Settings 2.0 scope. Defaults to `integration-azure`. Changing it forces recreation.",
			Optional:    true,
			Default:     DefaultScope,
			ForceNew:    true,
		},
		"feature_sets": {
			Type:        schema.TypeSet,
			Description: "Azure feature sets to enable (e.g. `microsoft_compute.virtualmachines_essential`). When empty, the extension defaults are used.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"credential": {
			Type:        schema.TypeList,
			Description: "HAS connection + Service Principal binding. At least one is required. dtctl always writes exactly one, but the API accepts a list.",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(Credential).Schema()},
		},
		"regions": {
			Type:        schema.TypeSet,
			Description: "Azure regions (locations) to monitor, e.g. `eastus`. Empty set = all locations the extension knows about. Maps to `locationFiltering` on the wire.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"subscription_filter": {
			Type:        schema.TypeSet,
			Description: "Subscription GUIDs to include or exclude (per `subscription_filtering_mode`). Empty set means \"all subscriptions reachable by the Service Principal\".",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"subscription_filtering_mode": {
			Type:        schema.TypeString,
			Description: "How to interpret `subscription_filter`. Defaults to `INCLUDE`.",
			Optional:    true,
			Default:     DefaultSubscriptionFilteringMode,
		},
		"tag_filter": {
			Type:        schema.TypeList,
			Description: "Filter monitored resources by Azure tag. Repeat the block to define multiple filters.",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(TagFilter).Schema()},
		},
		"tag_enrichment": {
			Type:        schema.TypeSet,
			Description: "Azure tag keys whose values are copied as Dynatrace tags on monitored entities.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"dt_label_enrichment": {
			Type:        schema.TypeList,
			Description: "Dynatrace labels (`dt.*`) applied to every monitored entity. Each block sets exactly one of `literal` or `tag_key`.",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(DTLabelEnrichment).Schema()},
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	featureSets := append([]string(nil), me.FeatureSets...)
	sort.Strings(featureSets)
	regions := append([]string(nil), me.Regions...)
	sort.Strings(regions)
	subs := append([]string(nil), me.SubscriptionFilter...)
	sort.Strings(subs)
	tagEnrichment := append([]string(nil), me.TagEnrichment...)
	sort.Strings(tagEnrichment)

	if err := properties.EncodeAll(map[string]any{
		"name":                        me.Name,
		"enabled":                     me.Enabled,
		"extension_version":           me.ExtensionVersion,
		"scope":                       me.Scope,
		"feature_sets":                featureSets,
		"regions":                     regions,
		"subscription_filter":         subs,
		"subscription_filtering_mode": me.SubscriptionFilteringMode,
		"tag_enrichment":              tagEnrichment,
	}); err != nil {
		return err
	}
	if err := properties.EncodeSlice("credential", me.Credentials); err != nil {
		return err
	}
	if err := properties.EncodeSlice("tag_filter", me.TagFilters); err != nil {
		return err
	}
	return properties.EncodeSlice("dt_label_enrichment", me.DTLabelEnrichments)
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.DecodeAll(map[string]any{
		"name":                        &me.Name,
		"enabled":                     &me.Enabled,
		"extension_version":           &me.ExtensionVersion,
		"scope":                       &me.Scope,
		"feature_sets":                &me.FeatureSets,
		"regions":                     &me.Regions,
		"subscription_filter":         &me.SubscriptionFilter,
		"subscription_filtering_mode": &me.SubscriptionFilteringMode,
		"tag_enrichment":              &me.TagEnrichment,
	}); err != nil {
		return err
	}
	if err := decoder.DecodeSlice("credential", &me.Credentials); err != nil {
		return err
	}
	if err := decoder.DecodeSlice("tag_filter", &me.TagFilters); err != nil {
		return err
	}
	if err := decoder.DecodeSlice("dt_label_enrichment", &me.DTLabelEnrichments); err != nil {
		return err
	}
	me.applyDefaults()
	return nil
}

// applyDefaults backfills attributes that the user did not set so the wire
// payload always carries a complete, server-accepted configuration. The four
// extension-internal attributes (activation_context, configuration_mode,
// deployment_mode, deployment_scope) are intentionally hidden from the
// user-facing schema and forced to their defaults here — the API may echo
// back different values, but we always re-send the canonical defaults so
// plans stay stable and users cannot override them.
func (me *Settings) applyDefaults() {
	if me.Scope == "" {
		me.Scope = DefaultScope
	}
	if me.SubscriptionFilteringMode == "" {
		me.SubscriptionFilteringMode = DefaultSubscriptionFilteringMode
	}
	me.ActivationContext = DefaultActivationContext
	me.ConfigurationMode = DefaultConfigurationMode
	me.DeploymentMode = DefaultDeploymentMode
	me.DeploymentScope = DefaultDeploymentScope
	for _, c := range me.Credentials {
		if c == nil {
			continue
		}
		if c.Type == "" {
			c.Type = DefaultCredentialType
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
			"enabled":            c.Enabled,
			"description":        c.Description,
			"connectionId":       c.ConnectionID,
			"servicePrincipalId": c.ServicePrincipalID,
			"type":               c.Type,
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

	azureBlock := map[string]any{
		"credentials":               creds,
		"locationFiltering":         toAnySlice(me.Regions),
		"subscriptionFiltering":     toAnySlice(me.SubscriptionFilter),
		"subscriptionFilteringMode": me.SubscriptionFilteringMode,
		"configurationMode":         me.ConfigurationMode,
		"deploymentMode":            me.DeploymentMode,
		"deploymentScope":           me.DeploymentScope,
		"tagFiltering":              tagFiltering,
		"tagEnrichment":             toAnySlice(me.TagEnrichment),
	}

	if len(me.DTLabelEnrichments) > 0 {
		labels := map[string]any{}
		for _, e := range me.DTLabelEnrichments {
			if e.Literal != "" {
				labels[e.Label] = map[string]any{"literal": e.Literal}
			} else {
				labels[e.Label] = map[string]any{"tagKey": e.TagKey}
			}
		}
		azureBlock["dtLabelsEnrichment"] = labels
	}

	value := map[string]any{
		"enabled":           me.Enabled,
		"description":       me.Name,
		"version":           me.ExtensionVersion,
		"activationContext": me.ActivationContext,
		"azure":             azureBlock,
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
	if fs, ok := wire.Value["featureSets"].([]any); ok {
		me.FeatureSets = me.FeatureSets[:0]
		for _, x := range fs {
			if s, ok := x.(string); ok {
				me.FeatureSets = append(me.FeatureSets, s)
			}
		}
	}

	azure, _ := wire.Value["azure"].(map[string]any)
	if azure == nil {
		me.applyDefaults()
		return nil
	}

	if v, ok := azure["subscriptionFilteringMode"].(string); ok {
		me.SubscriptionFilteringMode = v
	}
	if v, ok := azure["configurationMode"].(string); ok {
		me.ConfigurationMode = v
	}
	if v, ok := azure["deploymentMode"].(string); ok {
		me.DeploymentMode = v
	}
	if v, ok := azure["deploymentScope"].(string); ok {
		me.DeploymentScope = v
	}
	if creds, ok := azure["credentials"].([]any); ok {
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
			if v, ok := m["servicePrincipalId"].(string); ok {
				c.ServicePrincipalID = v
			}
			if v, ok := m["type"].(string); ok {
				c.Type = v
			}
			if c.Type == "" {
				// API echoes older configs without `type`; default to FEDERATED
				// so set-comparison stays stable across Reads (no eternal drift).
				c.Type = DefaultCredentialType
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
	if regs, ok := azure["locationFiltering"].([]any); ok {
		me.Regions = me.Regions[:0]
		for _, x := range regs {
			if s, ok := x.(string); ok {
				me.Regions = append(me.Regions, s)
			}
		}
	}
	if subs, ok := azure["subscriptionFiltering"].([]any); ok {
		me.SubscriptionFilter = me.SubscriptionFilter[:0]
		for _, x := range subs {
			if s, ok := x.(string); ok {
				me.SubscriptionFilter = append(me.SubscriptionFilter, s)
			}
		}
	}
	if tf, ok := azure["tagFiltering"].([]any); ok {
		me.TagFilters = me.TagFilters[:0]
		for _, x := range tf {
			if m, ok := x.(map[string]any); ok {
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
				me.TagFilters = append(me.TagFilters, f)
			}
		}
	}
	if te, ok := azure["tagEnrichment"].([]any); ok {
		me.TagEnrichment = me.TagEnrichment[:0]
		for _, x := range te {
			if s, ok := x.(string); ok {
				me.TagEnrichment = append(me.TagEnrichment, s)
			}
		}
	}
	if dtl, ok := azure["dtLabelsEnrichment"].(map[string]any); ok {
		me.DTLabelEnrichments = me.DTLabelEnrichments[:0]
		keys := make([]string, 0, len(dtl))
		for k := range dtl {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			e := &DTLabelEnrichment{Label: k}
			if m, ok := dtl[k].(map[string]any); ok {
				if v, ok := m["literal"].(string); ok {
					e.Literal = v
				}
				if v, ok := m["tagKey"].(string); ok {
					e.TagKey = v
				}
			}
			me.DTLabelEnrichments = append(me.DTLabelEnrichments, e)
		}
	}

	// `namespaces[]` and `eventHubsConfiguration[]` are deliberately ignored —
	// they are API-echo arrays (always returned as empty lists) and surfacing
	// them would produce eternal plan drift. See spec §4 gotcha #1.

	me.applyDefaults()
	return nil
}

func toAnySlice(in []string) []any {
	out := make([]any, 0, len(in))
	for _, s := range in {
		out = append(out, s)
	}
	return out
}
