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

// Package settings models the Dynatrace AWS monitoring configuration as a
// fully typed Terraform resource. Internally it is delivered through the
// generic Extensions 2.0 Monitoring Configuration API (extension
// "com.dynatrace.extension.da-aws"), but consumers see only first-class
// attributes — there is intentionally no JSON escape hatch.
package settings

import (
	"encoding/json"
	"sort"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AWSExtensionName is the fully qualified extension name that backs every AWS
// monitoring configuration.
const AWSExtensionName = "com.dynatrace.extension.da-aws"

// DefaultScope is the Settings 2.0 scope used by AWS DAC monitoring configs.
const DefaultScope = "integration-aws"

// Wire-level defaults that match dtctl's `create aws` behavior.
const (
	DefaultActivationContext = "DATA_ACQUISITION"
	DefaultDeploymentScope   = "SINGLE_ACCOUNT"
	DefaultDeploymentMode    = "AUTOMATED"
	DefaultConfigurationMode = "QUICK_START"
)

// Settings is the typed model exposed to Terraform.
type Settings struct {
	Name              string
	Enabled           bool
	ExtensionVersion  string
	ConnectionID      string
	AccountID         string
	Regions           []string
	FeatureSets       []string
	DeploymentRegion  string
	Scope             string
	ActivationContext string
	DeploymentScope   string
	DeploymentMode    string
	ConfigurationMode string
	SmartscapeEnabled bool

	TagFilters         TagFilters
	TagEnrichment      []string
	CloudWatchLogs     *CloudWatchLogsConfig
	CustomNamespaces   CustomNamespaces
	DTLabelEnrichments DTLabelEnrichments
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
			Description: "Version of `com.dynatrace.extension.da-aws` that this configuration targets. Optional — when omitted, the provider resolves the highest semver version installed on the tenant at create time (same behavior as `dtctl create aws`). The resolved value is persisted to state so subsequent plans are stable.",
			Optional:    true,
			Computed:    true,
		},
		"connection_id": {
			Type:        schema.TypeString,
			Description: "ObjectId of the `dynatrace_aws_connection` to use for authentication.",
			Required:    true,
		},
		"account_id": {
			Type:        schema.TypeString,
			Description: "12-digit AWS account id behind `connection_id`.",
			Required:    true,
		},
		"regions": {
			Type:        schema.TypeList,
			Description: "AWS regions to monitor. Used both as `regionFiltering` and as `metricsConfiguration.regions` on the wire.",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"feature_sets": {
			Type:        schema.TypeSet,
			Description: "CloudWatch metric feature sets to enable (e.g. `EC2_essential`). When empty the extension defaults are used.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"deployment_region": {
			Type:        schema.TypeString,
			Description: "AWS region the extension workload runs in. Defaults to the first entry in `regions`.",
			Optional:    true,
			Computed:    true,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "Settings 2.0 scope. Defaults to `integration-aws`. Changing it forces recreation.",
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
		"deployment_scope": {
			Type:        schema.TypeString,
			Description: "Deployment scope. Defaults to `SINGLE_ACCOUNT`.",
			Optional:    true,
			Default:     DefaultDeploymentScope,
		},
		"deployment_mode": {
			Type:        schema.TypeString,
			Description: "Deployment mode. Defaults to `AUTOMATED`.",
			Optional:    true,
			Default:     DefaultDeploymentMode,
		},
		"configuration_mode": {
			Type:        schema.TypeString,
			Description: "Configuration mode. Defaults to `QUICK_START`.",
			Optional:    true,
			Default:     DefaultConfigurationMode,
		},
		"smartscape_enabled": {
			Type:        schema.TypeBool,
			Description: "Whether Smartscape topology mapping is enabled. Defaults to true.",
			Optional:    true,
			Default:     true,
		},
		"tag_filter": {
			Type:        schema.TypeList,
			Description: "Filter monitored resources by AWS tag. Repeat the block to define multiple filters.",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(TagFilter).Schema()},
		},
		"tag_enrichment": {
			Type:        schema.TypeSet,
			Description: "AWS tag keys whose values are copied as Dynatrace tags on monitored entities.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"cloud_watch_logs": {
			Type:        schema.TypeList,
			Description: "CloudWatch Logs ingestion configuration. Omit the block to disable log ingestion.",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(CloudWatchLogsConfig).Schema()},
		},
		"custom_namespace": {
			Type:        schema.TypeList,
			Description: "Additional CloudWatch namespaces to ingest (e.g. AWS/GroundStation or your own custom namespace).",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(CustomNamespace).Schema()},
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
	if err := properties.EncodeAll(map[string]any{
		"name":               me.Name,
		"enabled":            me.Enabled,
		"extension_version":  me.ExtensionVersion,
		"connection_id":      me.ConnectionID,
		"account_id":         me.AccountID,
		"regions":            me.Regions,
		"feature_sets":       featureSets,
		"deployment_region":  me.DeploymentRegion,
		"scope":              me.Scope,
		"activation_context": me.ActivationContext,
		"deployment_scope":   me.DeploymentScope,
		"deployment_mode":    me.DeploymentMode,
		"configuration_mode": me.ConfigurationMode,
		"smartscape_enabled": me.SmartscapeEnabled,
		"tag_enrichment":     me.TagEnrichment,
		"cloud_watch_logs":   me.CloudWatchLogs,
	}); err != nil {
		return err
	}
	if err := properties.EncodeSlice("tag_filter", me.TagFilters); err != nil {
		return err
	}
	if err := properties.EncodeSlice("custom_namespace", me.CustomNamespaces); err != nil {
		return err
	}
	return properties.EncodeSlice("dt_label_enrichment", me.DTLabelEnrichments)
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.DecodeAll(map[string]any{
		"name":               &me.Name,
		"enabled":            &me.Enabled,
		"extension_version":  &me.ExtensionVersion,
		"connection_id":      &me.ConnectionID,
		"account_id":         &me.AccountID,
		"regions":            &me.Regions,
		"feature_sets":       &me.FeatureSets,
		"deployment_region":  &me.DeploymentRegion,
		"scope":              &me.Scope,
		"activation_context": &me.ActivationContext,
		"deployment_scope":   &me.DeploymentScope,
		"deployment_mode":    &me.DeploymentMode,
		"configuration_mode": &me.ConfigurationMode,
		"smartscape_enabled": &me.SmartscapeEnabled,
		"tag_enrichment":     &me.TagEnrichment,
		"cloud_watch_logs":   &me.CloudWatchLogs,
	}); err != nil {
		return err
	}
	if err := decoder.DecodeSlice("tag_filter", &me.TagFilters); err != nil {
		return err
	}
	if err := decoder.DecodeSlice("custom_namespace", &me.CustomNamespaces); err != nil {
		return err
	}
	if err := decoder.DecodeSlice("dt_label_enrichment", &me.DTLabelEnrichments); err != nil {
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
	if me.DeploymentRegion == "" && len(me.Regions) > 0 {
		me.DeploymentRegion = me.Regions[0]
	}
	if me.ActivationContext == "" {
		me.ActivationContext = DefaultActivationContext
	}
	if me.DeploymentScope == "" {
		me.DeploymentScope = DefaultDeploymentScope
	}
	if me.DeploymentMode == "" {
		me.DeploymentMode = DefaultDeploymentMode
	}
	if me.ConfigurationMode == "" {
		me.ConfigurationMode = DefaultConfigurationMode
	}
}

// wirePayload is the on-the-wire shape accepted by
// POST /api/v2/extensions/{name}/monitoringConfigurations.
type wirePayload struct {
	Scope string         `json:"scope"`
	Value map[string]any `json:"value"`
}

// MarshalJSON renders the typed Settings into the JSON body expected by the
// Extensions 2.0 Monitoring Configuration endpoint.
func (me *Settings) MarshalJSON() ([]byte, error) {
	me.applyDefaults()

	cwl := map[string]any{
		"enabled": false,
		"regions": []any{},
	}
	if me.CloudWatchLogs != nil {
		cwl["enabled"] = me.CloudWatchLogs.Enabled
		regs := []any{}
		for _, r := range me.CloudWatchLogs.Regions {
			regs = append(regs, r)
		}
		cwl["regions"] = regs
	}

	namespaces := []any{}
	for _, ns := range me.CustomNamespaces {
		metrics := []any{}
		for _, m := range ns.Metrics {
			metrics = append(metrics, map[string]any{
				"name":         m.Name,
				"unit":         m.Unit,
				"dimensions":   toAnySlice(m.Dimensions),
				"aggregations": toAnySlice(m.Aggregations),
				"type":         m.Type,
			})
		}
		namespaces = append(namespaces, map[string]any{
			"namespace":            ns.Namespace,
			"autoDiscoveryEnabled": ns.AutoDiscoveryEnabled,
			"metrics":              metrics,
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

	tagEnrichment := []any{}
	for _, k := range me.TagEnrichment {
		tagEnrichment = append(tagEnrichment, k)
	}

	awsBlock := map[string]any{
		"deploymentRegion":  me.DeploymentRegion,
		"deploymentScope":   me.DeploymentScope,
		"deploymentMode":    me.DeploymentMode,
		"configurationMode": me.ConfigurationMode,
		"credentials": []map[string]any{
			{
				"enabled":      me.Enabled,
				"description":  me.Name + " - " + me.AccountID,
				"connectionId": me.ConnectionID,
				"accountId":    me.AccountID,
			},
		},
		"namespaces": namespaces,
		"smartscapeConfiguration": map[string]any{
			"enabled": me.SmartscapeEnabled,
		},
		"metricsConfiguration": map[string]any{
			"enabled": true,
			"regions": me.Regions,
		},
		"cloudWatchLogsConfiguration": cwl,
		"regionFiltering":             me.Regions,
		"tagFiltering":                tagFiltering,
		"tagEnrichment":               tagEnrichment,
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
		awsBlock["dtLabelsEnrichment"] = labels
	}

	value := map[string]any{
		"enabled":           me.Enabled,
		"description":       me.Name,
		"version":           me.ExtensionVersion,
		"activationContext": me.ActivationContext,
		"aws":               awsBlock,
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
	aws, _ := wire.Value["aws"].(map[string]any)
	if aws == nil {
		me.applyDefaults()
		return nil
	}

	if v, ok := aws["deploymentRegion"].(string); ok {
		me.DeploymentRegion = v
	}
	if v, ok := aws["deploymentScope"].(string); ok {
		me.DeploymentScope = v
	}
	if v, ok := aws["deploymentMode"].(string); ok {
		me.DeploymentMode = v
	}
	if v, ok := aws["configurationMode"].(string); ok {
		me.ConfigurationMode = v
	}
	if creds, ok := aws["credentials"].([]any); ok && len(creds) > 0 {
		if c, ok := creds[0].(map[string]any); ok {
			if v, ok := c["connectionId"].(string); ok {
				me.ConnectionID = v
			}
			if v, ok := c["accountId"].(string); ok {
				me.AccountID = v
			}
		}
	}
	if mc, ok := aws["metricsConfiguration"].(map[string]any); ok {
		if regs, ok := mc["regions"].([]any); ok {
			me.Regions = me.Regions[:0]
			for _, x := range regs {
				if s, ok := x.(string); ok {
					me.Regions = append(me.Regions, s)
				}
			}
		}
	}
	if sm, ok := aws["smartscapeConfiguration"].(map[string]any); ok {
		if v, ok := sm["enabled"].(bool); ok {
			me.SmartscapeEnabled = v
		}
	}
	if cwl, ok := aws["cloudWatchLogsConfiguration"].(map[string]any); ok {
		tmp := &CloudWatchLogsConfig{}
		if v, ok := cwl["enabled"].(bool); ok {
			tmp.Enabled = v
		}
		if regs, ok := cwl["regions"].([]any); ok {
			for _, x := range regs {
				if s, ok := x.(string); ok {
					tmp.Regions = append(tmp.Regions, s)
				}
			}
		}
		// Avoid plan drift: the API always echoes this block. Only surface it
		// to Terraform state when it carries meaningful (non-default) config.
		if tmp.Enabled || len(tmp.Regions) > 0 {
			me.CloudWatchLogs = tmp
		}
	}
	if tf, ok := aws["tagFiltering"].([]any); ok {
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
	if te, ok := aws["tagEnrichment"].([]any); ok {
		me.TagEnrichment = me.TagEnrichment[:0]
		for _, x := range te {
			if s, ok := x.(string); ok {
				me.TagEnrichment = append(me.TagEnrichment, s)
			}
		}
	}
	if nss, ok := aws["namespaces"].([]any); ok {
		me.CustomNamespaces = me.CustomNamespaces[:0]
		for _, x := range nss {
			m, ok := x.(map[string]any)
			if !ok {
				continue
			}
			ns := &CustomNamespace{}
			if v, ok := m["namespace"].(string); ok {
				ns.Namespace = v
			}
			if v, ok := m["autoDiscoveryEnabled"].(bool); ok {
				ns.AutoDiscoveryEnabled = v
			}
			if metrics, ok := m["metrics"].([]any); ok {
				for _, my := range metrics {
					mm, ok := my.(map[string]any)
					if !ok {
						continue
					}
					cm := &CustomMetric{}
					if v, ok := mm["name"].(string); ok {
						cm.Name = v
					}
					if v, ok := mm["unit"].(string); ok {
						cm.Unit = v
					}
					if v, ok := mm["type"].(string); ok {
						cm.Type = v
					}
					if ds, ok := mm["dimensions"].([]any); ok {
						for _, d := range ds {
							if s, ok := d.(string); ok {
								cm.Dimensions = append(cm.Dimensions, s)
							}
						}
					}
					if ag, ok := mm["aggregations"].([]any); ok {
						for _, a := range ag {
							if s, ok := a.(string); ok {
								cm.Aggregations = append(cm.Aggregations, s)
							}
						}
					}
					ns.Metrics = append(ns.Metrics, cm)
				}
			}
			me.CustomNamespaces = append(me.CustomNamespaces, ns)
		}
	}
	if dtl, ok := aws["dtLabelsEnrichment"].(map[string]any); ok {
		me.DTLabelEnrichments = me.DTLabelEnrichments[:0]
		// Stable order so state-diffs are deterministic.
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
