package download

import (
	"fmt"
	"os"
	"reflect"
)

func (me DataSourceData) WriteDataSource(dlConfig DownloadConfig, resName string, resFolder string, replacedIDs ReplacedIDs) error {
	var err error
	var file *os.File
	fileName := dlConfig.TargetFolder + "/" + resFolder + "/" + "data_source.tf"
	os.Remove(fileName)
	if file, err = os.Create(fileName); err != nil {
		return err
	}

	for dsName, dataSource := range me {
		if contains(ResourceInfoMap[resName].HardcodedIds, dsName) {
			var writtenIDs []string
			for id, details := range dataSource.RESTMap {
				for _, replacedID := range replacedIDs[resName] {
					if id == replacedID.ID && !contains(writtenIDs, id) {
						if err := me.writer(file, dsName, details); err != nil {
							return err
						}
						writtenIDs = append(writtenIDs, id)
					}
				}
			}
		}
	}
	file.Close()

	return nil
}

func (me DataSourceData) writer(file *os.File, dsName string, detail *DataSourceDetails) error {
	if _, err := file.WriteString(fmt.Sprintf("data \"%s\" \"%s\" {\n", dsName, detail.UniqueName)); err != nil {
		file.Close()
		return err
	}
	for key, value := range detail.Values {
		switch t := value.(type) {
		case string:
			if _, err := file.WriteString(fmt.Sprintf("\t%s = \"%s\"\n", key, t)); err != nil {
				file.Close()
				return err
			}
		default:
			rv := reflect.ValueOf(value)
			switch rv.Kind() {
			case reflect.String:
				if _, err := file.WriteString(fmt.Sprintf("\t%s = \"%s\"\n", key, t)); err != nil {
					file.Close()
					return err
				}
			default:
				panic(fmt.Sprintf(">>>>> type %T not supported yet\n", t))
			}
		}
	}
	if _, err := file.WriteString("}\n\n"); err != nil {
		file.Close()
		return err
	}

	return nil
}

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
