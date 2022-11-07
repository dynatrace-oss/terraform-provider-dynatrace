package download

import (
	"fmt"
	"os"

	"github.com/dtcookie/hcl"
)

func (me ResourceData) WriteResourceSeparate(dlConfig DownloadConfig, resName string, resFolder string, resources Resources) error {
	var err error
	for _, resource := range resources {
		var file *os.File
		fileName := dlConfig.TargetFolder + "/" + resFolder + "/" + escf(resource.Name) + "." + resFolder + ".tf"
		os.Remove(fileName)
		if file, err = os.Create(fileName); err != nil {
			return err
		}
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", resName, Escape(resource.Name))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(resource.RESTObject, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n\n"); err != nil {
			file.Close()
			return err
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

func (me ResourceData) WriteResourceSingle(dlConfig DownloadConfig, resName string, resFolder string, resources Resources) error {
	var err error
	var file *os.File
	fileName := dlConfig.TargetFolder + "/" + resFolder + "/" + "main.tf"
	os.Remove(fileName)
	if file, err = os.Create(fileName); err != nil {
		return err
	}
	for _, resource := range resources {
		if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", resName, Escape(resource.Name))); err != nil {
			file.Close()
			return err
		}
		if err := hcl.ExportOpt(resource.RESTObject, file); err != nil {
			file.Close()
			return err
		}
		if _, err := file.WriteString("}\n\n"); err != nil {
			file.Close()
			return err
		}
		if resName == "dynatrace_dashboard" {
			me.writeDashboardSharing(file, resource.Name)
		}
	}
	file.Close()

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
	if _, err := file.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", "dynatrace_dashboard_sharing", Escape(name))); err != nil {
		file.Close()
		return err
	}
	if err := hcl.ExportOpt(restObject, file); err != nil {
		file.Close()
		return err
	}
	if _, err := file.WriteString("}\n\n"); err != nil {
		file.Close()
		return err
	}
	return nil
}
