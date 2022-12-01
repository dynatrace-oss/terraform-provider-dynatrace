package download

type ReplacedIDs map[string]map[string][]*ReplacedID

type ReplacedID struct {
	ID        string
	RefDS     string
	RefRes    string
	NoRes     bool
	Processed bool
}

func ProcessDataSourceIDs(resourceData ResourceData, dsData DataSourceData, replacedIDs ReplacedIDs) error {
	for resName, resources := range resourceData {
		if ResourceInfoMap[resName].HardcodedIds == nil {
			continue
		}

		repIds := ResourceInfoMap[resName].DsReplaceIds(resources, dsData)
		if replacedIDs[resName] == nil {
			replacedIDs[resName] = repIds
		} else {
			for repIdRes, repIdStruct := range repIds {
				replacedIDs[resName][repIdRes] = repIdStruct
			}
		}

	}

	return nil
}
