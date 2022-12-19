package download

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

var CanHardLink = true

func DetermineCanHardLink(targetFolder string) {
	var err error
	data := []byte(uuid.New().String())
	var linkedData []byte
	var orig string
	var link string

	if orig, err = filepath.Abs(path.Join(targetFolder, ".hardlink.orig.lck")); err != nil {
		fmt.Println("unable to determine hard link capabilities: " + err.Error())
		CanHardLink = false
		return
	}
	if link, err = filepath.Abs(path.Join(targetFolder, ".hardlink.link.lck")); err != nil {
		fmt.Println("unable to determine hard link capabilities: " + err.Error())
		CanHardLink = false
		return
	}

	if err = os.MkdirAll(targetFolder, os.ModePerm); err != nil {
		fmt.Println("unable to determine hard link capabilities: " + err.Error())
		CanHardLink = false
		return
	}

	if err = os.WriteFile(orig, data, 0644); err != nil {
		fmt.Println("unable to determine hard link capabilities: " + err.Error())
		CanHardLink = false
		return
	}
	if err = os.Link(orig, link); err != nil {
		fmt.Println("unable to determine hard link capabilities: " + err.Error())
		CanHardLink = false
		return
	}
	if err = os.Remove(orig); err != nil {
		fmt.Println("unable to determine hard link capabilities: " + err.Error())
		CanHardLink = false
		return
	}
	if linkedData, err = os.ReadFile(link); err != nil {
		fmt.Println("unable to determine hard link capabilities: " + err.Error())
		CanHardLink = false
		return
	}
	if err = os.Remove(link); err != nil {
		fmt.Println("unable to determine hard link capabilities: " + err.Error())
		CanHardLink = false
		return
	}
	if string(data) != string(linkedData) {
		fmt.Println("... hard link capabilities unstable in this environment")
		CanHardLink = false
	}
	// fmt.Println(".. using hard link capabilities")
}

func Download(environmentURL string, apiToken string, iamClientID string, iamAccountID string, iamClientSecret string, targetFolder string, refArg bool, comIdArg bool, migrateArg bool, excludeArg bool, linkArg bool, verbose bool, preview bool, resArgs map[string][]string) bool {
	if linkArg {
		DetermineCanHardLink(targetFolder)
	} else {
		CanHardLink = false
	}
	os.Setenv("dynatrace.secrets", "true")
	var err error
	var ResourceDataMap = ResourceData{}
	var DataSourceDataMap = DataSourceData{}
	var replacedIDs = ReplacedIDs{}

	dlConfig := DownloadConfig{
		EnvironmentURL:  environmentURL,
		APIToken:        apiToken,
		IAMClientID:     iamClientID,
		IAMAccountID:    iamAccountID,
		IAMClientSecret: iamClientSecret,
		TargetFolder:    targetFolder,
		References:      refArg,
		CommentedID:     comIdArg,
		Migrate:         migrateArg,
		Exclude:         excludeArg,
		ResourceNames:   resArgs,
		Verbose:         verbose,
	}

	if preview {
		if err = PrintPreview(dlConfig); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
		return true
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
	EnvironmentURL  string
	APIToken        string
	TargetFolder    string
	IAMClientID     string
	IAMAccountID    string
	IAMClientSecret string
	References      bool
	CommentedID     bool
	Migrate         bool
	Exclude         bool
	ResourceNames   map[string][]string
	Verbose         bool
}

func ValidateResource(keyVal string) (string, string) {
	res1 := ""
	res2 := ""
	parts := strings.Split(keyVal, "=")
	keyVal = parts[0]
	for resName := range ResourceInfoMap {
		if keyVal == resName {
			res1 = resName
			if len(parts) > 1 {
				res2 = parts[1]
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
