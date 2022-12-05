package download

import (
	"fmt"
	"reflect"
	"strings"
)

type MatchFunc func(s string, appInfo *DataSourceDetails) bool

func DefaultMatchFunc(s string, appInfo *DataSourceDetails) bool {
	return s == appInfo.UniqueName
}

func NameMatchFunc(s string, appInfo *DataSourceDetails) bool {
	return s == appInfo.Values["name"]
}

type MatchReplace struct {
	IDReplaceFunc IDReplaceFunc
	MatchFunc     MatchFunc
}

var DefaultMatchReplace = MatchReplace{
	MatchFunc:     DefaultMatchFunc,
	IDReplaceFunc: DefaultIDReplace,
}

var Settings20MatchReplace = MatchReplace{
	MatchFunc:     DefaultMatchFunc,
	IDReplaceFunc: Settings20IDReplace,
}

var NameReplace = MatchReplace{
	MatchFunc:     NameMatchFunc,
	IDReplaceFunc: NameIDReplace,
}

type IDReplaceFunc func(dsName string, appInfo *DataSourceDetails) string

func Settings20IDReplace(dsName string, appInfo *DataSourceDetails) string {
	return "data." + dsName + "." + appInfo.UniqueName + ".settings_20_id"
}

func DefaultIDReplace(dsName string, appInfo *DataSourceDetails) string {
	return "data." + dsName + "." + appInfo.UniqueName + ".id"
}

func NameIDReplace(dsName string, appInfo *DataSourceDetails) string {
	return "data." + dsName + "." + appInfo.UniqueName + ".name"
}

func Replace(resources Resources, dsName string, dataSourceData DataSourceData, replaceIdTemplate ReplacedID, matchReplace ...MatchReplace) map[string][]*ReplacedID {
	ids := map[string][]*ReplacedID{}
	idSet := map[string]string{}
	for _, resource := range resources {
		for _, id := range ReplaceResource(resource, dsName, dataSourceData, matchReplace...) {
			idSet[id] = id
		}
	}
	for id := range idSet {
		replacedId := ReplacedID{ID: id, RefDS: replaceIdTemplate.RefDS, RefRes: replaceIdTemplate.RefRes, Processed: replaceIdTemplate.Processed}
		if len(replaceIdTemplate.RefRes) > 0 {
			ids[replaceIdTemplate.RefRes] = append(ids[replaceIdTemplate.RefRes], &replacedId)
		} else {
			ids[dsName] = append(ids[dsName], &replacedId)
		}
	}
	return ids
}

func ReplaceResource(resource *Resource, dsName string, dataSourceData DataSourceData, matchReplace ...MatchReplace) map[string]string {
	fn := DefaultMatchReplace
	if len(matchReplace) > 0 {
		fn = matchReplace[0]
	}
	rep := replacer{
		resource:       resource,
		dsName:         dsName,
		dataSourceData: dataSourceData,
		MatchReplace:   &fn,
		replaceIDs:     map[string]string{},
	}
	rep.replace(reflect.ValueOf(resource.RESTObject))
	// res := []string{}
	// for k := range rep.replaceIDs {
	// 	res = append(res, k)
	// }
	return rep.replaceIDs
}

type replacer struct {
	resource       *Resource
	dsName         string
	dataSourceData DataSourceData
	// idReplaceFn    IDReplaceFunc
	MatchReplace *MatchReplace
	replaceIDs   map[string]string
}

func (me *replacer) replace(rv reflect.Value) {
	if rv.IsZero() {
		return
	}
	switch rv.Type().Kind() {
	case reflect.Bool:
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	case reflect.Float32, reflect.Float64:
	case reflect.String:
		s := rv.String()
		for id, appInfo := range me.dataSourceData[me.dsName].RESTMap {
			if appInfo.ComputedValues != nil {
				if settings_20_id, ok := appInfo.ComputedValues["settings_20_id"]; ok {
					if me.MatchReplace != nil && me.MatchReplace.IDReplaceFunc != nil {
						testresult := me.MatchReplace.IDReplaceFunc(me.dsName, &DataSourceDetails{Values: map[string]interface{}{}, UniqueName: "unique"})
						if strings.HasSuffix(testresult, "settings_20_id") {
							id = settings_20_id.(string)
						}
					}
				}
			}

			if me.MatchReplace.MatchFunc(s, appInfo) {
				if me.MatchReplace.IDReplaceFunc != nil {
					replacement := me.MatchReplace.IDReplaceFunc(me.dsName, appInfo)
					rv.Set(reflect.ValueOf("HCL-UNQUOTE-" + replacement).Convert(rv.Type()))
					if me.resource.Variables == nil {
						me.resource.Variables = map[string]string{}
					}
					if strings.HasSuffix(replacement, ".name") {
						me.resource.Variables[replacement] = appInfo.Values["name"].(string)
					} else {
						me.resource.Variables[replacement] = id
					}
				}
				me.replaceIDs[id] = id
			} else if me.MatchReplace.IDReplaceFunc != nil {
				searchStr := me.MatchReplace.IDReplaceFunc(me.dsName, appInfo)
				replacement := strings.ReplaceAll(s, id, "${"+searchStr+"}")
				if replacement != s {
					rv.Set(reflect.ValueOf(replacement).Convert(rv.Type()))
					if me.resource.Variables == nil {
						me.resource.Variables = map[string]string{}
					}
					me.resource.Variables[searchStr] = id
					me.replaceIDs[id] = id
				}
			}
		}
	case reflect.Pointer:
		me.replace(rv.Elem())
	case reflect.Map:
	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			me.replace(rv.Index(i))
		}
	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			me.replace(rv.Field(i))
		}
	case reflect.Interface:
		me.replace(rv.Elem())
	default:
		panic(fmt.Sprintf("unsupported type %v (kind: %v)", rv.Type(), rv.Type().Kind()))
	}
}
