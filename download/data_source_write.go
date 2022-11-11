package download

import (
	"fmt"
	"os"
)

func (me DataSourceData) WriteDataSource(dsFile *os.File, dlConfig DownloadConfig, resName string, resFolder string, replacedIDs ReplacedIDs) error {
	// var err error
	// var file *os.File
	// fileName := dlConfig.TargetFolder + "/" + resFolder + "/" + ".data_source.tf"
	// os.Remove(fileName)
	// if file, err = os.Create(fileName); err != nil {
	// 	return err
	// }

	for dsName, dataSource := range me {
		if contains(ResourceInfoMap[resName].HardcodedIds, dsName) {
			var writtenIDs []string
			for id, values := range dataSource.RESTMap {
				for _, replacedID := range replacedIDs[resName] {
					if id == replacedID && !contains(writtenIDs, id) {
						if err := me.writer(dsFile, dsName, values); err != nil {
							return err
						}
						writtenIDs = append(writtenIDs, id)
					}
				}
			}
		}
	}
	// file.close()

	return nil
}

func (me DataSourceData) writer(file *os.File, dsName string, values map[string]interface{}) error {
	if _, err := file.WriteString(fmt.Sprintf("data \"%s\" \"%s\" {\n", dsName, Escape(values["name"].(string)))); err != nil {
		file.Close()
		return err
	}
	for key, value := range values {
		switch t := value.(type) {
		case string:
			if _, err := file.WriteString(fmt.Sprintf("\t%s = \"%s\"\n", key, t)); err != nil {
				file.Close()
				return err
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
