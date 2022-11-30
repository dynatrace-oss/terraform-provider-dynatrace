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

		replacedIDs[resName] = ResourceInfoMap[resName].DsReplaceIds(resources, dsData)
	}

	return nil
}
