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

package goldenstate

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	cfg "github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Read(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	return read(ctx, d, m, "")
}

func read(ctx context.Context, d *schema.ResourceData, m any, indent string) diag.Diagnostics {
	// logging.File.Println(indent, "-- READ --")
	creds, _ := cfg.Credentials(m, cfg.CredValNone)
	var wg sync.WaitGroup
	wg.Add(len(serviceMap))
	for key := range serviceMap {
		go func() {
			defer wg.Done()
			CommonRead(ctx, d, creds, key, indent+"  ")
		}()
	}
	wg.Wait()
	return diag.Diagnostics{}
}

var regexpNameId = regexp.MustCompile(`\[(.*)\]\ (.*)`)

func trimName(name string) string {
	name = strings.TrimSpace(name)
	for strings.HasPrefix(name, "[") {
		name = name[1:]
		name = strings.TrimSpace(name)
	}
	for strings.HasSuffix(name, "]") {
		name = name[:len(name)-1]
		name = strings.TrimSpace(name)
	}
	name = strings.TrimSpace(name)
	if len(name) > 24 {
		name = fmt.Sprintf("%s...", name[:21])
	}
	return name
}

func CommonRead(ctx context.Context, d *schema.ResourceData, creds *settings.Credentials, key export.ResourceType, indent string) error {
	// logging.File.Println(indent, "CommonRead", key)

	idMap := map[string]string{}

	service := serviceMap[key](creds)
	if stubs, err := service.List(ctx); err == nil {
		for _, stub := range stubs {
			idMap[stub.ID] = stub.Name
		}
	}
	// decoder := confighcl.StateDecoderFrom(d, Resource())
	if untypedIDs, ok := d.GetOk(string(key)); ok {
		if idSet, ok := untypedIDs.(*schema.Set); ok {
			for _, untypedID := range idSet.List() {
				if id, ok := untypedID.(string); ok {
					if matches := regexpNameId.FindStringSubmatch(id); len(matches) == 3 {
						idMap[matches[2]] = matches[1]
					} else {
						idMap[id] = id
					}
				}
			}
		}
	}
	ids := []string{}
	for id, name := range idMap {
		if id != name {
			ids = append(ids, fmt.Sprintf("[ %-24s ] %s", trimName(name), id))
		} else {
			ids = append(ids, id)
		}
	}
	d.Set(string(key), ids)
	return nil
}
