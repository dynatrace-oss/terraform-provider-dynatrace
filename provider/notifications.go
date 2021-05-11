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

package provider

import (
	"context"
	"log"
	"reflect"
	"strings"

	"github.com/dtcookie/dynatrace/api/config/notifications"
	"github.com/dtcookie/dynatrace/rest"
	"github.com/dtcookie/dynatrace/terraform"
	"github.com/dtcookie/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type diff struct {
	old string
	new string
}

// NotificationCfg has no documentation
type NotificationCfg struct {
	Config notifications.NotificationConfig `json:"config"` // The name of the notification configuration.
}

type notificationConfigs struct {
	diffs map[string]diff
}

func (nc *notificationConfigs) AttachDiffSuppressFunc(sch *schema.Schema) {
	sch.DiffSuppressFunc = logging.EnableSchemaDiff(func(k, old, new string, d *schema.ResourceData) bool {
		if strings.HasPrefix(k, "metadata.") {
			return true
		}

		if nc.diffs == nil {
			nc.diffs = map[string]diff{}
		}
		if strings.HasSuffix(k, ".#") {
			if old == "0" && new == "1" {
				prefix := k[0 : len(k)-1]
				found := false
				for st := range nc.diffs {
					if strings.HasPrefix(st, prefix) {
						found = true
						break
					}
				}
				if !found {
					return true
				}
			}
		}
		nc.diffs[k] = diff{old, new}
		return false
	})
	if sch.Elem != nil {
		switch typedSchema := sch.Elem.(type) {
		case *schema.Schema:
			nc.AttachDiffSuppressFunc(typedSchema)
		case *schema.Resource:
			nc.AttachDiffSuppressFuncs(typedSchema.Schema)
		}
	}
}

func (nc *notificationConfigs) AttachDiffSuppressFuncs(schemas map[string]*schema.Schema) {
	if schemas == nil {
		return
	}
	for _, sch := range schemas {
		nc.AttachDiffSuppressFunc(sch)
	}
}

func (nc *notificationConfigs) wrap(fn func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		nc.diffs = map[string]diff{}
		result := fn(ctx, d, m)
		nc.diffs = map[string]diff{}
		return result
	}
}

// Resource produces terraform resource definition for Management Zones
func (nc *notificationConfigs) Resource() *schema.Resource {
	resource := terraform.ResourceFor(new(NotificationCfg))
	resource.CreateContext = logging.Enable(nc.wrap(nc.Create))
	resource.UpdateContext = logging.Enable(nc.wrap(nc.Update))
	resource.ReadContext = logging.Enable(nc.wrap(nc.Read))
	resource.DeleteContext = logging.Enable(nc.wrap(nc.Delete))
	nc.AttachDiffSuppressFuncs(resource.Schema)
	resource.Importer = &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext}
	resource.Schema["service_now_notification_config"].DiffSuppressFunc = logging.EnableSchemaDiff(func(k, old, new string, d *schema.ResourceData) bool {
		return strings.HasSuffix(k, ".password") && old == "---terraform---null"
	})
	resource.Schema["jira_notification_config"].DiffSuppressFunc = logging.EnableSchemaDiff(func(k, old, new string, d *schema.ResourceData) bool {
		return strings.HasSuffix(k, ".password") && old == "---terraform---null"
	})
	resource.Schema["ansible_tower_notification_config"].DiffSuppressFunc = logging.EnableSchemaDiff(func(k, old, new string, d *schema.ResourceData) bool {
		return strings.HasSuffix(k, ".password") && old == "---terraform---null"
	})

	return resource
}

// Create expects the configuration of a Management Zone within the given ResourceData
// and send them to the Dynatrace Server in order to create that resource
func (nc *notificationConfigs) Create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if config.HTTPVerbose {
		log.Println("NotificationConfigs.Create(..)")
	}
	var err error
	var resolver terraform.Resolver
	if resolver, err = terraform.NewResolver(d); err != nil {
		return diag.FromErr(err)
	}
	var untypedConfig interface{}
	if untypedConfig, err = resolver.Resolve(reflect.TypeOf(NotificationCfg{})); err != nil {
		return diag.FromErr(err)
	}

	notificationConfig := untypedConfig.(NotificationCfg)
	notificationConfig.Config.SetID(nil)
	conf := m.(*config.ProviderConfiguration)
	notificationConfigService := notifications.NewService(conf.DTenvURL, conf.APIToken)
	rest.Verbose = config.HTTPVerbose
	var notificationConfigStub *notifications.NotificationConfigStub
	if notificationConfigStub, err = notificationConfigService.Create(notificationConfig.Config); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(notificationConfigStub.ID)
	return diag.Diagnostics{}
	// return nc.Read(ctx, d, m)
}

// Update expects the configuration of a Management Zone within the given ResourceData
// and send them to the Dynatrace Server in order to update that resource
func (nc *notificationConfigs) Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if config.HTTPVerbose {
		log.Println("NotificationConfigs.Update(..)")
	}
	var err error
	var resolver terraform.Resolver
	if resolver, err = terraform.NewResolver(d); err != nil {
		return diag.FromErr(err)
	}
	var untypedConfig interface{}
	if untypedConfig, err = resolver.Resolve(reflect.TypeOf(NotificationCfg{})); err != nil {
		return diag.FromErr(err)
	}

	notificationConfig := untypedConfig.(NotificationCfg)
	notificationConfig.Config.SetID(opt.NewString(d.Id()))
	conf := m.(*config.ProviderConfiguration)
	notificationConfigService := notifications.NewService(conf.DTenvURL, conf.APIToken)
	rest.Verbose = config.HTTPVerbose
	if err = notificationConfigService.Update(notificationConfig.Config); err != nil {
		return diag.FromErr(err)
	}
	return nc.Read(ctx, d, m)
}

// Read queries the Dynatrace Server for the configuration of a Management Zone
// identified by the ID within the given ResourceData
func (nc *notificationConfigs) Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if config.HTTPVerbose {
		log.Println("NotificationConfigs.Read(..)")
	}
	var err error
	conf := m.(*config.ProviderConfiguration)
	notificationConfigService := notifications.NewService(conf.DTenvURL, conf.APIToken)
	rest.Verbose = config.HTTPVerbose
	var notificationConfig notifications.NotificationConfig
	if notificationConfig, err = notificationConfigService.Get(d.Id()); err != nil {
		return diag.FromErr(err)
	}
	switch notCfg := notificationConfig.(type) {
	case *notifications.WebHookNotificationConfig:
	case *notifications.EmailNotificationConfig:
	case *notifications.SlackNotificationConfig:
	case *notifications.PagerDutyNotificationConfig:
	case *notifications.TrelloNotificationConfig:
	case *notifications.VictorOpsNotificationConfig:
	case *notifications.XMattersNotificationConfig:
	case *notifications.OpsGenieNotificationConfig:
	case *notifications.HipChatNotificationConfig:
	case *notifications.JiraNotificationConfig:
		notCfg.Password = opt.NewString("---terraform---null")
	case *notifications.AnsibleTowerNotificationConfig:
		notCfg.Password = opt.NewString("---terraform---null")
	case *notifications.ServiceNowNotificationConfig:
		notCfg.Password = opt.NewString("---terraform---null")
	}

	notificationConfig.SetID(nil)
	notificationCfg := NotificationCfg{Config: notificationConfig}
	if err = terraform.ToTerraform(notificationCfg, d); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

// Delete a Management Zone on the Dynatrace Server
func (nc *notificationConfigs) Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if config.HTTPVerbose {
		log.Println("NotificationConfigs.Delete(..)")
	}
	var err error
	conf := m.(*config.ProviderConfiguration)
	notificationConfigService := notifications.NewService(conf.DTenvURL, conf.APIToken)
	rest.Verbose = config.HTTPVerbose
	if err = notificationConfigService.Delete(d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}
