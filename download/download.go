package download

import (
	"fmt"
	"os"
	"strings"
)

func Download(environmentURL string, apiToken string, targetFolder string, refArg bool, comIdArg bool, resArgs map[string][]string) bool {
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
		ResourceNames:  resArgs,
	}

	if err = ResourceDataMap.ProcessRead(dlConfig); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	if dlConfig.References {
		if err = DataSourceDataMap.ProcessRead(dlConfig, ResourceDataMap); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}

		if replacedIDs, err = ProcessDataSourceIDs(ResourceDataMap, DataSourceDataMap); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
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
	ResourceNames  map[string][]string
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
