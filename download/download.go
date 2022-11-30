package download

import (
	"fmt"
	"os"
	"strings"
)

func Download(environmentURL string, apiToken string, targetFolder string, refArg bool, comIdArg bool, migrateArg bool, resArgs map[string][]string) bool {
	os.Setenv("dynatrace.secrets", "true")
	var err error
	var ResourceDataMap = ResourceData{}
	var DataSourceDataMap = DataSourceData{}
	var replacedIDs = ReplacedIDs{}

	dlConfig := DownloadConfig{
		EnvironmentURL: environmentURL,
		APIToken:       apiToken,
		TargetFolder:   targetFolder,
		References:     refArg,
		CommentedID:    comIdArg,
		Migrate:        migrateArg,
		ResourceNames:  resArgs,
	}

	if err = ResourceDataMap.ProcessRead(dlConfig); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	run := true
	for run {
		run = false
		if err = ResourceDataMap.RequiresIntervention(dlConfig); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}

		if dlConfig.References {
			if err = DataSourceDataMap.ProcessRead(dlConfig, ResourceDataMap); err != nil {
				fmt.Println(err.Error())
				os.Exit(0)
			}

			if err = ProcessDataSourceIDs(ResourceDataMap, DataSourceDataMap, replacedIDs); err != nil {
				fmt.Println(err.Error())
				os.Exit(0)
			}
		}

		if len(resArgs) > 0 && replacedIDs != nil {
			for _, replacedId := range replacedIDs {
				for _, repIdRes := range replacedId {
					for _, repId := range repIdRes {
						if !repId.Processed {
							run = true
						}
					}
				}
			}
			if run {
				if err = ResourceDataMap.ProcessRepIdRead(dlConfig, replacedIDs); err != nil {
					fmt.Println(err.Error())
					os.Exit(0)
				}
			}
		}
	}

	if err = ProcessWrite(dlConfig, ResourceDataMap, DataSourceDataMap, replacedIDs); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	return true
}

type DownloadConfig struct {
	EnvironmentURL string
	APIToken       string
	TargetFolder   string
	References     bool
	CommentedID    bool
	Migrate        bool
	ResourceNames  map[string][]string
	// DependencyRes  []DependencyResource
}

func ValidateResource(keyVal string) (string, string) {
	res1 := ""
	res2 := ""
	for resName := range ResourceInfoMap {
		if strings.HasPrefix(keyVal, resName) {
			res1 = resName
			if strings.HasPrefix(keyVal, resName+"=") {
				res2 = keyVal[len(resName)+1:]
			}
		}
	}
	return res1, res2
}

func (me DownloadConfig) MatchResource(name string) bool {
	if _, found := me.ResourceNames[name]; found {
		return true
	}
	return false
}

func (me DownloadConfig) MatchDataSource(dsName string, resDataMap ResourceData) bool {
	for resName := range resDataMap {
		if ResourceInfoMap[resName].HardcodedIds == nil {
			continue
		}
		for _, dataSource := range ResourceInfoMap[resName].HardcodedIds {
			if dsName == dataSource {
				return true
			}
		}
	}
	return false
}

func (me DownloadConfig) MatchID(resName string, id string) bool {
	for _, entityId := range me.ResourceNames[resName] {
		if id == entityId {
			return true
		}
	}
	return false
}

// func FindDependencies(resArgs map[string][]string) ([]DependencyResource, int) {
// 	var newResources []DependencyResource
// 	nestedLevel := 0
// 	found := true
// 	for found {
// 		nestedLevel++
// 		found = false
// 		for resName := range resArgs {
// 			if ResourceInfoMap[resName].HardcodedIds != nil {
// 				for _, hcId := range ResourceInfoMap[resName].HardcodedIds {
// 					if _, exists := resArgs[hcId.ResName]; !exists {
// 						resArgs[hcId.ResName] = nil
// 						newResources = append(newResources, DependencyResource{hcId.ResName, nestedLevel})
// 						found = true
// 					}
// 				}
// 			}
// 		}
// 	}

// 	return newResources, nestedLevel
// }

// type DependencyResource struct {
// 	Name  string
// 	Level int
// }
