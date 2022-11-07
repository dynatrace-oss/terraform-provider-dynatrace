package download

type ReplacedIDs map[string][]string

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

func ProcessResourceIDs(resourceData ResourceData) error {
	for resName := range resourceData {
		if ResourceInfoMap[resName].HardcodedIds == nil {
			continue
		}
		ResourceInfoMap[resName].ResReplaceIds(resName, resourceData)
	}

	return nil
}
