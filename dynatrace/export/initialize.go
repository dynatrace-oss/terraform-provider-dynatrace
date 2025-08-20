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

package export

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/cache"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
)

func Initialize(cfgGetter config.Getter) (environment *Environment, err error) {
	flags, tailArgs := createFlags()
	if flags.FlagMigrationOutput && flags.FollowReferences {
		return nil, errors.New("-ref and -migrate are mutually exclusive")
	}
	if flags.FlagMigrationOutput {
		flags.FollowReferences = true
		flags.PersistIDs = true
	}
	if err = ConfigureRESTLog(); err != nil {
		return nil, errors.New("unable to configure log file for REST activity: " + err.Error())
	}

	args := os.Args
	if len(args) == 1 {
		return nil, nil
	}
	if strings.TrimSpace(args[1]) != "-export" {
		return nil, nil
	}

	settings.ExportRunning = true
	os.Setenv("dynatrace.secrets", "true")
	cache.Enable()
	resArgs := map[string][]string{}
	if flags.Exclude {
		for resourceType := range AllResources {
			excludeListed := false
			for _, excludeListedResourceType := range GetExcludeListedResources() {
				if resourceType == excludeListedResourceType {
					excludeListed = true
					break
				}
			}
			if !excludeListed {
				resArgs[string(resourceType)] = []string{}
			}
		}
		for _, idx := range tailArgs {
			key, _ := ValidateResource(idx)
			if len(key) == 0 {
				return nil, fmt.Errorf("unknown resource `%s`", idx)
			}
			delete(resArgs, key)
		}
	} else {
		effectiveTailArgs := map[string]string{}
		excludeMode := false
		for _, idx := range tailArgs {
			if excludeMode {
				key, _ := ValidateResource(idx)
				if len(key) == 0 {
					return nil, fmt.Errorf("unknown resource `%s`", idx)
				}
				delete(effectiveTailArgs, key)
				continue
			}
			if idx == "-exclude" {
				excludeMode = true
				continue
			} else if idx == "*" {
				effectiveTailArgs["*"] = "*"
			} else {
				idx = ToParent(idx)
				effectiveTailArgs[idx] = idx
				key, id := ValidateResource(idx)
				if len(key) == 0 {
					return nil, fmt.Errorf("unknown resource `%s`", idx)
				}

				for _, child := range ResourceType(key).GetChildren() {
					if len(id) == 0 {
						effectiveTailArgs[string(child)] = string(child)
					} else {
						effectiveTailArgs[string(child)+"="+id] = string(child) + "=" + id
					}
				}
			}
		}
		for _, idx := range effectiveTailArgs {
			if idx == "*" {
				for resourceType := range AllResources {
					excludeListed := false
					for _, excludeListedResourceType := range GetExcludeListedResources() {
						if resourceType == excludeListedResourceType {
							excludeListed = true
							break
						}
					}
					if !excludeListed {
						resArgs[string(resourceType)] = nil
					}
				}
			} else {
				key, id := ValidateResource(idx)
				if len(key) == 0 {
					return nil, fmt.Errorf("unknown resource `%s`", idx)
				}
				stored, ok := resArgs[key]
				if ok {
					if stored != nil {
						if id == "" {
							resArgs[key] = nil
						} else {
							stored = append(stored, id)
							resArgs[key] = stored
						}
					}
				} else {
					if id == "" {
						resArgs[key] = nil
					} else {
						stored = []string{id}
						resArgs[key] = stored
					}
				}
			}
		}

		if len(resArgs) == 0 {
			for resourceType := range AllResources {
				excludeListed := false
				for _, excludeListedResourceType := range GetExcludeListedResources() {
					if resourceType == excludeListedResourceType {
						excludeListed = true
						break
					}
				}
				if !excludeListed {
					resArgs[string(resourceType)] = []string{}
				}
			}
		}
	}

	targetFolder := os.Getenv("DYNATRACE_TARGET_FOLDER")
	if targetFolder == "" {
		fmt.Println("The environment variable DYNATRACE_TARGET_FOLDER has not been set - using folder 'configuration' as default")
		targetFolder = "configuration"
	}
	if os.Getenv("DYNATRACE_CLEAN_TARGET_FOLDER") == "true" {
		os.RemoveAll(targetFolder)
	}

	var credentials *rest.Credentials
	configResult, _ := config.ProviderConfigureGeneric(context.Background(), cfgGetter)
	if credentials, err = config.Credentials(configResult, config.CredValExport); err != nil {
		return nil, err
	}

	// If ONLY child resources are getting exported we
	// don't treat them as such. Request from Omar Zaal
	requestingOnlyChildResources := true
	for k := range resArgs {
		if !ResourceType(k).IsChildResource() {
			requestingOnlyChildResources = false
		}
	}
	if requestingOnlyChildResources {
		for _, v := range AllResources {
			v.Parent = nil
		}
	}

	return &Environment{
		OutputFolder:          targetFolder,
		Credentials:           credentials,
		Modules:               map[ResourceType]*Module{},
		Flags:                 flags,
		ResArgs:               resArgs,
		ChildResourceOverride: requestingOnlyChildResources,
	}, nil
}

func createFlags() (flags Flags, tailArgs []string) {
	flag.Bool("export", true, "")
	refArg := flag.Bool("ref", false, "enable data sources and dependencies. mutually exclusive with -migrate")
	dataSourceArg := flag.Bool("datasources", false, "when resolving dependencies eligible resources will be referred as data sources")
	comIdArg := flag.Bool("id", false, "enable commented ids")
	migrateArg := flag.Bool("migrate", false, "enable migration output. mutually exclusive with -ref")
	verbose := flag.Bool("v", false, "enable verbose logging")
	linkArg := flag.Bool("link", false, "enable hard links for .requires_attention and .flawed")
	preview := flag.Bool("preview", false, "preview resource statistics for environment export")
	flat := flag.Bool("flat", false, "prevent creating a module structure")
	importStateV2 := flag.Bool("import-state-v2", false, "deprecated - use `import-state`")
	importState := flag.Bool("import-state", false, "automatically initialize the terraform module and import downloaded resources to the state")
	exclude := flag.Bool("exclude", false, "exclude specified resources")
	skipTerraformInit := flag.Bool("skip-terraform-init", false, "prevent the command line `terraform init` from getting executed after all the configuration files have been created")

	flag.Parse()

	importStateFlag := (importState != nil && *importState == true) || (importStateV2 != nil && *importStateV2 == true)

	return Flags{
		FollowReferences:    *refArg,
		PersistIDs:          *comIdArg,
		FlagMigrationOutput: *migrateArg,
		FlagVerbose:         *verbose,
		FlagHardLinks:       *linkArg,
		FlagPreviewOnly:     *preview,
		Flat:                *flat,
		ImportStateV2:       importStateFlag,
		Exclude:             *exclude,
		DataSources:         *dataSourceArg,
		SkipTerraformInit:   *skipTerraformInit,
	}, flag.Args()
}

func ToParent(keyVal string) string {
	res1 := ""
	res2 := ""
	parts := strings.Split(keyVal, "=")
	keyVal = parts[0]
	for resName := range AllResources {
		if keyVal == string(resName) {
			for resName.IsChildResource() {
				resName = resName.GetParent()
			}
			res1 = string(resName)
			if len(parts) > 1 {
				res2 = parts[1]
			}
		}
	}
	if len(res2) == 0 {
		return res1
	}
	return fmt.Sprintf("%s=%s", res1, res2)
}

func ValidateResource(keyVal string) (string, string) {
	res1 := ""
	res2 := ""
	parts := strings.Split(keyVal, "=")
	keyVal = parts[0]
	for resName := range AllResources {
		if keyVal == string(resName) {
			res1 = string(resName)
			if len(parts) > 1 {
				res2 = parts[1]
			}
		}
	}
	return res1, res2
}

type Flags struct {
	FollowReferences    bool
	PersistIDs          bool
	FlagMigrationOutput bool
	FlagVerbose         bool
	FlagHardLinks       bool
	FlagPreviewOnly     bool
	Flat                bool
	ImportStateV2       bool
	Exclude             bool
	DataSources         bool
	SkipTerraformInit   bool
	Include             bool
}
