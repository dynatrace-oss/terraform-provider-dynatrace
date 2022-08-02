package common

import (
	"strings"

	tagapi "github.com/dtcookie/dynatrace/api/config/topology/tag"
)

// TagSubsetCheck checks that the input tags are a subset of source tags
// Arguments: source slice of tags, input slice of tags
// Return: true if subset, false if not
func TagSubsetCheck(source []tagapi.Tag, input []tagapi.Tag) bool {
	for _, inputTag := range input {
		found := false
		for _, restTag := range source {
			if restTag.Key == inputTag.Key {
				if restTag.Value == nil && inputTag.Value == nil {
					found = true
					break
				} else if restTag.Value != nil && inputTag.Value != nil && *restTag.Value == *inputTag.Value {
					found = true
					break
				}
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// StringsToTags processes the slice of string tags into a slice of tag structs
// Arguments: slice of string tags, pointer to slice of tag structs
func StringsToTags(tagList []interface{}, tags *[]tagapi.Tag) {
	for _, iTag := range tagList {
		var tag tagapi.Tag
		if strings.Contains(iTag.(string), "=") {
			tagSplit := strings.Split(iTag.(string), "=")
			tag.Key = tagSplit[0]
			tag.Value = &tagSplit[1]
		} else {
			tag.Key = iTag.(string)
		}
		*tags = append(*tags, tag)
	}
}
