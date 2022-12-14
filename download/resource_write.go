package download

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/hclgen"
)

func (me ResourceData) WriteResource(dlConfig DownloadConfig, resName string, resFolder string, resources Resources, resNameCnt NameCounter) error {
	var err error
	for _, resource := range resources {
		if !CanHardLink && len(resource.ReqInter.Type) > 0 {
			continue
		}
		var file *os.File
		fileName := dlConfig.TargetFolder + "/" + resFolder + "/" + resFolder + "." + resource.UniqueName + ".tf"
		fileName = strings.ToLower(fileName)
		for {
			if _, err := os.Stat(fileName); err == nil {
				fileName = strings.TrimSuffix(fileName, ".tf") + "_1.tf"
			} else {
				break
			}
		}

		if file, err = os.Create(fileName); err != nil {
			return err
		}

		if dlConfig.CommentedID {
			resource.Dedup()
			comments := resource.ReqInter.Message

			if dlConfig.CommentedID {
				comments = append(comments, "DEFINE "+resName+"."+resource.UniqueName+"."+"id = "+resource.ID)
				if len(resource.Variables) > 0 {
					for k, v := range resource.Variables {
						comments = append(comments, "DEFINE "+k+" = "+v)
					}

				}
			}
			if len(resource.Variables) > 0 {
				for k, v := range resource.Variables {
					comments = append(comments, "DEFINE "+k+" = "+v)
				}

			}
			if err := hclgen.Export(resource.RESTObject, file, resName, resource.UniqueName, comments...); err != nil {
				file.Close()
				return err
			}
		} else {
			if err := hclgen.Export(resource.RESTObject, file, resName, resource.UniqueName); err != nil {
				file.Close()
				return err
			}
		}

		if resName == "dynatrace_dashboard" || resName == "dynatrace_json_dashboard" {
			if err := me.writeDashboardSharing(file, resource.UniqueName); err != nil {
				file.Close()
				return err
			}
		}
		file.Close()
	}

	return nil
}

func (me ResourceData) writeDashboardSharing(file *os.File, name string) error {
	var found bool
	var resource *Resource
	for _, res := range me["dynatrace_dashboard_sharing"] {
		if res.UniqueName == name {
			resource = res
			found = true
			break
		}
	}
	if !found {
		file.Close()
		return nil
	}
	comments := []string{
		"DEFINE " + "dynatrace_dashboard_sharing" + "." + name + "." + "id = " + resource.ID,
	}
	if len(resource.Variables) > 0 {
		for k, v := range resource.Variables {
			comments = append(comments, "DEFINE "+k+" = "+v)
		}
	}
	if err := hclgen.Export(resource.RESTObject, file, "dynatrace_dashboard_sharing", name, comments...); err != nil {
		file.Close()
		return err
	}
	return nil
}

func (me ResourceData) WriteResReqAttn(dlConfig DownloadConfig) error {
	var err error
	for resName := range InterventionInfoMap {
		if resources, exists := me[resName]; exists {
			for _, resource := range resources {
				if resource.ReqInter.Type == "" {
					continue
				}

				folderName := dlConfig.TargetFolder + "/" + string(resource.ReqInter.Type) + "/" + strings.TrimPrefix(resName, "dynatrace_")
				if _, err := os.Stat(folderName); errors.Is(err, os.ErrNotExist) {
					err := os.MkdirAll(folderName, os.ModePerm)
					if err != nil {
						return err
					}
				}

				fileName := folderName + "/" + strings.TrimPrefix(resName, "dynatrace_") + "." + resource.UniqueName + ".tf"

				if CanHardLink {
					orig, _ := filepath.Abs(dlConfig.TargetFolder + "/" + strings.TrimPrefix(resName, "dynatrace_") + "/" + strings.TrimPrefix(resName, "dynatrace_") + "." + resource.UniqueName + ".tf")
					link, _ := filepath.Abs(fileName)
					os.Link(orig, link)
					continue
				}

				var file *os.File
				os.Remove(fileName)
				if file, err = os.Create(fileName); err != nil {
					return err
				}
				resource.Dedup()

				comments := resource.ReqInter.Message

				if dlConfig.CommentedID {
					comments = append(comments, "DEFINE "+resName+"."+resource.UniqueName+"."+"id = "+resource.ID)
					if len(resource.Variables) > 0 {
						for k, v := range resource.Variables {
							comments = append(comments, "DEFINE "+k+" = "+v)
						}

					}
				}
				if err := hclgen.Export(resource.RESTObject, file, resName, resource.UniqueName, comments...); err != nil {
					file.Close()
					return err
				}

				file.Close()
			}
		}
	}
	return nil
}
