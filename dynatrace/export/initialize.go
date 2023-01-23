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
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/cache"
)

func Initialize() (environment *Environment, err error) {
	flags, tailArgs := createFlags()
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

	os.Setenv("dynatrace.secrets", "true")
	cache.Enable()
	resArgs := map[string][]string{}
	for _, idx := range tailArgs {
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

	if len(resArgs) == 0 {
		for resourceType := range AllResources {
			blackListed := false
			for _, blackListedResourceType := range BlackListedResources {
				if resourceType == blackListedResourceType {
					blackListed = true
					break
				}
			}
			if !blackListed {
				resArgs[string(resourceType)] = []string{}
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

	var credentials *settings.Credentials
	if credentials, err = settings.CreateCredentials(); err != nil {
		return nil, err
	}

	return &Environment{
		OutputFolder: targetFolder,
		Credentials:  credentials,
		Modules:      map[ResourceType]*Module{},
		Flags:        flags,
		ResArgs:      resArgs,
	}, nil
}

func createFlags() (flags Flags, tailArgs []string) {
	flag.Bool("export", true, "")
	refArg := flag.Bool("ref", false, "enable data sources and dependencies")
	comIdArg := flag.Bool("id", false, "enable commented ids")
	migrateArg := flag.Bool("migrate", false, "enable migration output")
	verbose := flag.Bool("v", false, "enable verbose logging")
	linkArg := flag.Bool("link", false, "enable hard links for .requires_attention and .flawed")
	preview := flag.Bool("preview", false, "preview resource statistics for environment export")
	flat := flag.Bool("flat", false, "prevent creating a module structure")
	importState := flag.Bool("import-state", false, "automatically initialize the terraform module and import downloaded resources to the state")

	flag.Parse()
	return Flags{
		FollowReferences:    *refArg,
		PersistIDs:          *comIdArg,
		FlagMigrationOutput: *migrateArg,
		FlagVerbose:         *verbose,
		FlagHardLinks:       *linkArg,
		FlagPreviewOnly:     *preview,
		Flat:                *flat,
		ImportState:         *importState,
	}, flag.Args()
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
	ImportState         bool
}
