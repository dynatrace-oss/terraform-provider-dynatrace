package download

type ReplacedIDs map[string][]ReplacedID

type ReplacedID struct {
	ID     string
	RefDS  string
	RefRes string
}

func ProcessDataSourceIDs(resourceData ResourceData, dsData DataSourceData) (ReplacedIDs, error) {
	var replacedIDs = ReplacedIDs{}
	for resName, resources := range resourceData {
		if ResourceInfoMap[resName].HardcodedIds == nil {
			continue
		}
		replacedIDs[resName] = ResourceInfoMap[resName].DsReplaceIds(resources, dsData)
	}

	return replacedIDs, nil
}
