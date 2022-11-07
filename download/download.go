package download

import (
	"fmt"
	"os"
)

type DownloadConfig struct {
	EnvironmentURL string
	APIToken       string
	TargetFolder   string
	ResourceNames  map[string][]string
	SingleFile     bool
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

func Download(environmentURL string, apiToken string, targetFolder string, m map[string][]string, singleFile bool) bool {
	var err error
	var ResourceDataMap = ResourceData{}
	var DataSourceDataMap = DataSourceData{}
	var replacedIDs = ReplacedIDs{}
	dlConfig := DownloadConfig{
		EnvironmentURL: environmentURL,
		APIToken:       apiToken,
		TargetFolder:   targetFolder,
		ResourceNames:  m,
		SingleFile:     singleFile,
	}

	if err = ResourceDataMap.ProcessRead(dlConfig); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	if err = ProcessResourceIDs(ResourceDataMap); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	// Option 2
	// if err = DataSourceDataMap.ProcessRead(dlConfig, ResourceDataMap); err != nil {
	// 	fmt.Println(err.Error())
	// 	os.Exit(0)
	// }

	// if replacedIDs, err = ProcessDataSourceIDs(ResourceDataMap, DataSourceDataMap); err != nil {
	// 	fmt.Println(err.Error())
	// 	os.Exit(0)
	// }

	if err = ProcessWrite(dlConfig, ResourceDataMap, DataSourceDataMap, replacedIDs); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	return true
}
