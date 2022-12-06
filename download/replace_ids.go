package download

type ReplacedIDs map[string]map[string][]*ReplacedID

type ReplacedID struct {
	ID        string
	RefDS     string
	RefRes    string
	NoRes     bool
	Processed bool
}

func MergeReplacedIDs(a []*ReplacedID, b []*ReplacedID) []*ReplacedID {
	if len(a) == 0 {
		return b
	}
	if len(b) == 0 {
		return a
	}
	m := map[string]*ReplacedID{}
	for _, e := range a {
		m[e.ID] = e
	}
	for _, e := range b {
		m[e.ID] = e
	}
	result := []*ReplacedID{}
	for _, e := range m {
		result = append(result, e)
	}
	return result
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
