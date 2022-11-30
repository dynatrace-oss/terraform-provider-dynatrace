package download

type ReplacedIDs map[string]map[string][]*ReplacedID

type ReplacedID struct {
	ID        string
	RefDS     string
	RefRes    string
	Processed bool
}

func ProcessDataSourceIDs(resourceData ResourceData, dsData DataSourceData, replacedIDs ReplacedIDs) error {
	for resName, resources := range resourceData {
		if ResourceInfoMap[resName].HardcodedIds == nil {
			continue
		}
		repIds := ResourceInfoMap[resName].DsReplaceIds(resources, dsData)
		if len(repIds) > 0 {
			replacedIDs[resName] = make(map[string][]*ReplacedID)
			replacedIDs[resName][repIds[0].RefRes] = repIds
		}
	}

	return nil
}
