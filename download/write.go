package download

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

func ProcessWrite(dlConfig DownloadConfig, resourceDataMap ResourceData, dataSourceDataMap DataSourceData, replacedIDs ReplacedIDs) error {
	var err error

	if err := os.RemoveAll(dlConfig.TargetFolder); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	os.MkdirAll(dlConfig.TargetFolder, os.ModePerm)
	var mainFile *os.File
	fileName := dlConfig.TargetFolder + "/main.tf"
	if mainFile, err = os.Create(fileName); err != nil {
		return err
	}
	for resName, resources := range resourceDataMap {
		if len(resources) == 0 {
			continue
		}
		if resName == "dynatrace_dashboard_sharing" {
			continue
		}
		resFolder := strings.TrimPrefix(resName, "dynatrace_")
		os.MkdirAll(dlConfig.TargetFolder+"/"+resFolder, os.ModePerm)
		if !dlConfig.SingleFile {
			if err = resourceDataMap.WriteResourceSeparate(dlConfig, resName, resFolder, resources); err != nil {
				return err
			}
		} else {
			if err = resourceDataMap.WriteResourceSingle(dlConfig, resName, resFolder, resources); err != nil {
				return err
			}
		}

		if ResourceInfoMap[resName].HardcodedIds != nil {
			dataSourceDataMap.WriteDataSource(dlConfig, resName, resFolder, replacedIDs)
		}

		if _, err := mainFile.WriteString(fmt.Sprintf("module \"%s\" {\n", resFolder)); err != nil {
			mainFile.Close()
			return err
		}
		if _, err := mainFile.WriteString(fmt.Sprintf("\tsource = \"./%s\"\n", resFolder)); err != nil {
			mainFile.Close()
			return err
		}
		if _, err := mainFile.WriteString("}\n\n"); err != nil {
			mainFile.Close()
			return err
		}
	}
	mainFile.Close()

	return nil
}

var forbiddenFileNameChars = []string{"<", ">", ":", "\"", "/", "|", "?", "*", "	", "\r", "\n", "\f", "\v"}

func escf(s string) string {
	for _, ch := range forbiddenFileNameChars {
		s = strings.ReplaceAll(s, ch, "_")
	}
	return s
}

func Escape(s string) string {
	result := ""
	for _, c := range s {
		if unicode.IsLetter(c) {
			result = result + string(c)
		} else if unicode.IsDigit(c) {
			result = result + string(c)
		} else if c == '-' {
			result = result + string(c)
		} else if c == '_' {
			result = result + string(c)
		} else {
			result = result + "_"
		}
	}
	return result
}
