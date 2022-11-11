package download

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/version"
)

func ProcessWrite(dlConfig DownloadConfig, resourceDataMap ResourceData, dataSourceDataMap DataSourceData, replacedIDs ReplacedIDs) error {
	var err error

	if err := os.RemoveAll(dlConfig.TargetFolder); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	os.MkdirAll(dlConfig.TargetFolder, os.ModePerm)
	// var mainFile *os.File
	// fileName := dlConfig.TargetFolder + "/main.tf"
	// if mainFile, err = os.Create(fileName); err != nil {
	// 	return err
	// }

	var dsFile *os.File
	dsFileName := dlConfig.TargetFolder + "/" + ".data_source.tf"
	os.Remove(dsFileName)
	if len(replacedIDs) != 0 {
		if dsFile, err = os.Create(dsFileName); err != nil {
			return err
		}
	}

	var mainFile *os.File
	mainFileName := dlConfig.TargetFolder + "/" + "main.tf"
	if dlConfig.SingleFile {
		if mainFile, err = os.Create(mainFileName); err != nil {
			return err
		}
	}

	for resName, resources := range resourceDataMap {
		if len(resources) == 0 {
			continue
		}
		if resName == "dynatrace_dashboard_sharing" {
			continue
		}
		resFolder := strings.TrimPrefix(resName, "dynatrace_")
		// os.MkdirAll(dlConfig.TargetFolder+"/"+resFolder, os.ModePerm)
		if !dlConfig.SingleFile {
			if err = resourceDataMap.WriteResourceSeparate(dlConfig, resName, resFolder, resources); err != nil {
				return err
			}
		} else {
			// if err = resourceDataMap.WriteResourceSingle(dlConfig, resName, resFolder, resources); err != nil {
			// 	return err
			// }
			if err = resourceDataMap.WriteResourceSingle(mainFile, dlConfig, resName, resFolder, resources); err != nil {
				return err
			}
		}

		// if err := writeNestedProviderFile(dlConfig.TargetFolder, resFolder); err != nil {
		// 	return err
		// }

		if ResourceInfoMap[resName].HardcodedIds != nil && dlConfig.ReplaceIDs == "datasource" {
			// dataSourceDataMap.WriteDataSource(dlConfig, resName, resFolder, replacedIDs)
			dataSourceDataMap.WriteDataSource(dsFile, dlConfig, resName, resFolder, replacedIDs)
		}

		// if _, err := mainFile.WriteString(fmt.Sprintf("module \"%s\" {\n", resFolder)); err != nil {
		// 	mainFile.Close()
		// 	return err
		// }
		// if _, err := mainFile.WriteString(fmt.Sprintf("  source = \"./%s\"\n", resFolder)); err != nil {
		// 	mainFile.Close()
		// 	return err
		// }
		// if _, err := mainFile.WriteString("  providers = {\n    dynatrace = dynatrace.default\n  }\n}\n\n"); err != nil {
		// 	mainFile.Close()
		// 	return err
		// }

	}
	dsFile.Close()
	mainFile.Close()

	if err := writeProviderFile(dlConfig.TargetFolder); err != nil {
		return err
	}

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

func writeProviderFile(targetFolder string) error {
	var err error
	var providerFile *os.File
	fileName := targetFolder + "/providers.tf"
	if providerFile, err = os.Create(fileName); err != nil {
		return err
	}
	content := `terraform {
	required_providers {
		dynatrace = {
		version = "${version}"
		source = "dynatrace-oss/dynatrace"
		}
	}
}
	
# provider "dynatrace" {
#   alias        = "default"
#   dt_env_url   = "https://########.live.dynatrace.com/"
#   dt_api_token = "dt0c01.#########################################################################################"
# }
	`
	content = strings.Replace(content, "${version}", version.Current, 1)
	if _, err := providerFile.WriteString(content); err != nil {
		providerFile.Close()
		return err
	}

	return nil
}
