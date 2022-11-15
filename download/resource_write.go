package download

import (
	"os"

	"github.com/dtcookie/hcl"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/hclgen"
)

func (me ResourceData) WriteResourceSeparate(dlConfig DownloadConfig, resName string, resFolder string, resources Resources, resNameCnt NameCounter) error {
	var err error
	for _, resource := range resources {
		// var nameCounter NameCounter
		// nameCounter.Replace = func(s string, cnt int) string {
		// 	return fmt.Sprintf("%s_%d", s, cnt)
		// }
		// nameCounter.Numbering()
		var file *os.File
		fileName := dlConfig.TargetFolder + "/" + resFolder + "/" + resFolder + "." + escf(resource.Name) + ".tf"
		// fileName := dlConfig.TargetFolder + "/" + resFolder + "." + escf(resource.Name) + ".tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}

		if dlConfig.CommentedID {
			if err := hclgen.Export(resource.RESTObject, file, resName, resNameCnt.Numbering(escape(resource.Name)), "id = "+resource.ID); err != nil {
				file.Close()
				return err
			}
		} else {
			if err := hclgen.Export(resource.RESTObject, file, resName, escape(resource.Name)); err != nil {
				file.Close()
				return err
			}
		}

		if resName == "dynatrace_dashboard" {
			if err := me.writeDashboardSharing(file, resource.Name); err != nil {
				file.Close()
				return err
			}

		}
		file.Close()
	}

	return nil
}

func (me ResourceData) WriteResourceSingle(mainFile *os.File, dlConfig DownloadConfig, resName string, resFolder string, resources Resources, resNameCnt NameCounter) error {
	// var err error
	// var file *os.File
	// fileName := dlConfig.TargetFolder + "/" + resFolder + "/" + "main.tf"
	// os.Remove(fileName)
	// if file, err = os.Create(fileName); err != nil {
	// 	return err
	// }

	for _, resource := range resources {
		if err := hclgen.Export(resource.RESTObject, mainFile, resName, resNameCnt.Numbering(escape(resource.Name))); err != nil {
			mainFile.Close()
			return err
		}
		if resName == "dynatrace_dashboard" {
			me.writeDashboardSharing(mainFile, resource.Name)
		}
	}
	// file.Close()

	return nil
}

func (me ResourceData) writeDashboardSharing(file *os.File, name string) error {
	var restObject hcl.Marshaler
	var found bool
	for _, resource := range me["dynatrace_dashboard_sharing"] {
		if resource.Name == name {
			restObject = resource.RESTObject
			found = true
			break
		}
	}
	if !found {
		file.Close()
		return nil
	}
	if err := hclgen.Export(restObject, file, "dynatrace_dashboard_sharing", escape(name)); err != nil {
		file.Close()
		return err
	}
	return nil
}
