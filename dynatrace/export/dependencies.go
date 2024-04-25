/**
* @license
* Copyright 2020 Dynatrace LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package export

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
)

// To allow -target to work with dependencies at an atomic level
var ATOMIC_DEPENDENCIES = os.Getenv("DYNATRACE_ATOMIC_DEPENDENCIES") == "true"

type Dependency interface {
	Replace(environment *Environment, s string, replacingIn ResourceType, resourceId string) (string, []any)
	ResourceType() ResourceType
	DataSourceType() DataSourceType
	IsParent() bool
}

func Coalesce(d Dependency) Dependency {
	switch dep := d.(type) {
	case *entityds:
		return &entityds{Type: dep.Type, Pattern: dep.Pattern, Coalesce: true}
	}
	return d
}

var Dependencies = struct {
	ManagementZone       Dependency
	DashboardLinkID      func(parent bool) Dependency
	HyperLinkDashboardID func() Dependency
	LegacyID             func(resourceType ResourceType) Dependency
	ID                   func(resourceType ResourceType) Dependency
	QuotedID             func(resourceType ResourceType) Dependency
	ResourceID           func(resourceType ResourceType, parent bool) Dependency
	ServiceMethod        Dependency
	Service              Dependency
	HostGroup            Dependency
	Host                 Dependency
	Disk                 Dependency
	ProcessGroup         Dependency
	ProcessGroupInstance Dependency
	RequestAttribute     Dependency
	// CustomApplication    Dependency
	// MobileApplication       Dependency
	DeviceApplicationMethod Dependency
	// Application               Dependency
	ApplicationMethod Dependency
	// SyntheticTest             Dependency
	// HttpCheck                 Dependency
	K8sCluster                Dependency
	CloudApplication          Dependency
	CloudApplicationNamespace Dependency
	EnvironmentActiveGate     Dependency
	Tenant                    Dependency
	GlobalPolicy              Dependency
}{
	ManagementZone:       &mgmzdep{ResourceTypes.ManagementZoneV2},
	DashboardLinkID:      func(parent bool) Dependency { return &dashlinkdep{ResourceTypes.JSONDashboardBase, parent} },
	HyperLinkDashboardID: func() Dependency { return &dashdep{ResourceTypes.JSONDashboardBase, false} },
	LegacyID:             func(resourceType ResourceType) Dependency { return &legacyID{resourceType} },
	ID:                   func(resourceType ResourceType) Dependency { return &iddep{resourceType: resourceType, quoted: false} },
	QuotedID:             func(resourceType ResourceType) Dependency { return &iddep{resourceType: resourceType, quoted: true} },
	ResourceID:           func(resourceType ResourceType, parent bool) Dependency { return &resourceIDDep{resourceType, parent} },
	ServiceMethod:        &entityds{"SERVICE_METHOD", "SERVICE_METHOD-[A-Z0-9]{16}", false},
	Service:              &entityds{"SERVICE", "SERVICE-[A-Z0-9]{16}", false},
	HostGroup:            &entityds{"HOST_GROUP", "HOST_GROUP-[A-Z0-9]{16}", false},
	Host:                 &entityds{"HOST", "HOST-[A-Z0-9]{16}", false},
	Disk:                 &entityds{"DISK", "DISK-[A-Z0-9]{16}", false},
	ProcessGroup:         &entityds{"PROCESS_GROUP", "PROCESS_GROUP-[A-Z0-9]{16}", false},
	ProcessGroupInstance: &entityds{"PROCESS_GROUP_INSTANCE", "PROCESS_GROUP_INSTANCE-[A-Z0-9]{16}", false},
	RequestAttribute:     &reqAttName{ResourceTypes.RequestAttribute},
	// CustomApplication:    &entityds{"CUSTOM_APPLICATION", "CUSTOM_APPLICATION-[A-Z0-9]{16}", false},
	// MobileApplication:       &entityds{"MOBILE_APPLICATION", "MOBILE_APPLICATION-[A-Z0-9]{16}", false},
	DeviceApplicationMethod: &entityds{"DEVICE_APPLICATION_METHOD", "DEVICE_APPLICATION_METHOD-[A-Z0-9]{16}", false},
	// Application:               &entityds{"APPLICATION", "APPLICATION-[A-Z0-9]{16}", false},
	ApplicationMethod: &entityds{"APPLICATION_METHOD", "APPLICATION_METHOD-[A-Z0-9]{16}", false},
	// SyntheticTest:             &entityds{"SYNTHETIC_TEST", "SYNTHETIC_TEST-[A-Z0-9]{16}", false},
	// HttpCheck:                 &entityds{"HTTP_CHECK", "HTTP_CHECK-[A-Z0-9]{16}", false},
	K8sCluster:                &entityds{"KUBERNETES_CLUSTER", "KUBERNETES_CLUSTER-[A-Z0-9]{16}", false},
	CloudApplication:          &entityds{"CLOUD_APPLICATION", "CLOUD_APPLICATION-[A-Z0-9]{16}", false},
	CloudApplicationNamespace: &entityds{"CLOUD_APPLICATION_NAMESPACE", "CLOUD_APPLICATION_NAMESPACE-[A-Z0-9]{16}", false},
	EnvironmentActiveGate:     &entityds{"ENVIRONMENT_ACTIVE_GATE", "ENVIRONMENT_ACTIVE_GATE-[A-Z0-9]{16}", false},
	Tenant:                    &tenantds{},
	GlobalPolicy:              &policyds{},
}

type mgmzdep struct {
	resourceType ResourceType
}

func (me *mgmzdep) IsParent() bool {
	return false
}

func (me *mgmzdep) ResourceType() ResourceType {
	return me.resourceType
}

func (me *mgmzdep) DataSourceType() DataSourceType {
	return ""
}

func (me *mgmzdep) Replace(environment *Environment, s string, replacingIn ResourceType, resourceId string) (string, []any) {
	var replacePattern string
	if ATOMIC_DEPENDENCIES {
		replacePattern = "${var.%s_%s.value.name}"
	} else {
		replacePattern = "${var.%s.%s.name}"
	}
	if environment.Flags.Flat {
		replacePattern = "${%s.%s.name}"
	}
	resources := []any{}
	for _, resource := range environment.Module(me.resourceType).Resources {
		if resource.Status.IsOneOf(ResourceStati.Erronous, ResourceStati.Excluded) {
			continue
		}
		found := false
		resOrDsType := func() string {
			return string(me.resourceType)
		}

		if resource.IsReferencedAsDataSource() {
			resOrDsType = func() string {
				return string(me.resourceType.AsDataSource())
			}
			replacePattern = "${data.%s.%s.name}"
			if environment.Flags.Flat {
				replacePattern = "${data.%s.%s.name}"
			}
		} else {
			if ATOMIC_DEPENDENCIES {
				replacePattern = "${var.%s_%s.value.name}"
			} else {
				replacePattern = "${var.%s.%s.name}"
			}
			if environment.Flags.Flat {
				replacePattern = "${%s.%s.name}"
			}
		}

		escaped := regexp.QuoteMeta(resource.Name)

		if strings.Contains(s, `"managementZone": {`) {
			m1 := regexp.MustCompile(fmt.Sprintf(`"managementZone": {([\s]*[\S ]+[\s]+)"name":(\s)+"%s"([^}]*)}`, escaped))
			replaced := m1.ReplaceAllString(s, fmt.Sprintf(`"managementZone": {$1"name":$2"%s"$3}`, fmt.Sprintf(replacePattern, resOrDsType(), resource.UniqueName)))
			if replaced != s {
				s = replaced
				found = true
			}
		}
		if strings.Contains(s, `management_zones`) {
			m1 := regexp.MustCompile(fmt.Sprintf(`management_zones\s+= \[(.*)\"%s\"(.*)\]`, escaped))
			replaced := m1.ReplaceAllString(s, fmt.Sprintf(`management_zones = [ $1"%s"$2 ]`, fmt.Sprintf(replacePattern, resOrDsType(), resource.UniqueName)))
			if replaced != s {
				s = replaced
				found = true
			}
		}
		if strings.Contains(s, `mzName(`) {
			m1 := regexp.MustCompile(fmt.Sprintf(`mzName\(([~\\\"]+)%s([~\\\"]+)\)`, escaped))
			replaced := m1.ReplaceAllString(s, fmt.Sprintf(`mzName($1%s$2)`, fmt.Sprintf(replacePattern, resOrDsType(), resource.UniqueName)))
			if replaced != s {
				s = replaced
				found = true
			}
		}
		if found {
			resources = append(resources, resource)
		}
	}
	return s, resources
}

type dashlinkdep struct {
	resourceType ResourceType
	parent       bool
}

func (me *dashlinkdep) IsParent() bool {
	return me.parent
}

func (me *dashlinkdep) ResourceType() ResourceType {
	return me.resourceType
}

func (me *dashlinkdep) DataSourceType() DataSourceType {
	return ""
}

func (me *dashlinkdep) Replace(environment *Environment, s string, replacingIn ResourceType, resourceId string) (string, []any) {
	replacePattern := "${%s.%s.id}"
	resources := []any{}

	resourcesMap := map[string]*Resource{}

	resource, found := environment.Module(me.resourceType).Resources[resourceId]
	if found {
		resourcesMap[resourceId] = resource
	} else {
		return s, resources
	}

	for id, resource := range resourcesMap {
		resOrDsType := func() string {
			return string(me.resourceType)
		}
		if resource.Type == replacingIn {
			replacePattern = "${%s.%s.id}"
		} else {
			if resource.IsReferencedAsDataSource() {
				resOrDsType = func() string {
					return string(me.resourceType.AsDataSource())
				}
				replacePattern = "${data.%s.%s.id}"
				if environment.Flags.Flat {
					replacePattern = "${data.%s.%s.id}"
				}
			} else {
				replacePattern = "${%s.%s.id}"
				if environment.Flags.Flat {
					replacePattern = "${%s.%s.id}"
				}
			}
		}
		found := false

		escaped := regexp.QuoteMeta(id)

		m1 := regexp.MustCompile(fmt.Sprintf(`link_id(.*)=(.*)"%s"`, escaped))
		replaced := m1.ReplaceAllString(s, fmt.Sprintf("link_id$1=$2\"%s\"", fmt.Sprintf(replacePattern, resOrDsType(), resource.UniqueName)))
		if replaced != s {
			s = replaced
			found = true
		}
		if found {
			resources = append(resources, resource)
		}
	}
	return s, resources
}

type dashdep struct {
	resourceType ResourceType
	parent       bool
}

func (me *dashdep) IsParent() bool {
	return me.parent
}

func (me *dashdep) ResourceType() ResourceType {
	return me.resourceType
}

func (me *dashdep) DataSourceType() DataSourceType {
	return ""
}

var dashboardHyperLinkIdRegex = regexp.MustCompile(`#dashboard;[^)]*id=([^${;)-]+-[^;)-]+-[^;)-]+-[^;)-]+-[^-;)]+)`)

func (me *dashdep) Replace(environment *Environment, s string, replacingIn ResourceType, resourceId string) (string, []any) {
	resources := []any{}

	resourcesMap := map[string]*Resource{}

	matches := dashboardHyperLinkIdRegex.FindAllStringSubmatch(s, -1)

	for _, match := range matches {
		dashId := match[1]
		resource, found := environment.Module(me.resourceType).Resources[dashId]
		if found {
			resourcesMap[dashId] = resource
		} else {
			return s, resources
		}
	}

	childDescriptor := environment.Module(replacingIn).GetDescriptor()
	isParent := !environment.ChildResourceOverride && childDescriptor.Parent != nil && string(*childDescriptor.Parent) == string(me.resourceType)

	var replacePattern string
	if ATOMIC_DEPENDENCIES {
		replacePattern = "${var.%s_%s.value.id}"
	} else {
		replacePattern = "${var.%s.%s.id}"
	}

	if environment.Flags.Flat || isParent {
		replacePattern = "${%s.%s.id}"
	}

	for id, resource := range resourcesMap {
		resOrDsType := func() string {
			return string(me.resourceType)
		}
		if resource.Type == replacingIn {
			replacePattern = "${%s.%s.id}"
		} else {
			if resource.IsReferencedAsDataSource() {
				resOrDsType = func() string {
					return string(me.resourceType.AsDataSource())
				}
				replacePattern = "${data.%s.%s.id}"
				if environment.Flags.Flat {
					replacePattern = "${data.%s.%s.id}"
				}
			} else {
				if ATOMIC_DEPENDENCIES {
					replacePattern = "${var.%s_%s.value.id}"
				} else {
					replacePattern = "${var.%s.%s.id}"
				}
				if environment.Flags.Flat || isParent {
					replacePattern = "${%s.%s.id}"
				}
			}
		}
		found := false
		if strings.Contains(s, id) {
			s = strings.ReplaceAll(s, id, fmt.Sprintf(replacePattern, resOrDsType(), resource.UniqueName))
			found = true
		}
		if found {
			resources = append(resources, resource)
		}
	}
	return s, resources
}

type legacyID struct {
	resourceType ResourceType
}

func (me *legacyID) IsParent() bool {
	return false
}

func (me *legacyID) ResourceType() ResourceType {
	return me.resourceType
}

func (me *legacyID) DataSourceType() DataSourceType {
	return ""
}

func (me *legacyID) Replace(environment *Environment, s string, replacingIn ResourceType, resourceId string) (string, []any) {
	var replacePattern string
	if ATOMIC_DEPENDENCIES {
		replacePattern = "${var.%s_%s.value.legacy_id}"
	} else {
		replacePattern = "${var.%s.%s.legacy_id}"
	}

	if environment.Flags.Flat {
		replacePattern = "${%s.%s.legacy_id}"
	}

	resources := []any{}

	regexMap := OptimizedKeyRegexV1
	idRegexType := V1_CONFIG_ID_REGEX
	if me.resourceType == ResourceTypes.ManagementZoneV2 {
		regexMap = OptimizedKeyRegexMzLegacy
		idRegexType = LEGACY_MZ_ID_REGEX
	}

	extractedIds := environment.Module(replacingIn).Resources[resourceId].GetExtractedIdsPerRegexType(idRegexType, s, regexMap)

	if len(extractedIds) < 1 {
		return s, resources
	}

	legacyIdMap := environment.Module(me.resourceType).GetLegacyIdMap()

	for extractedId, _ := range extractedIds {
		resource, exists := legacyIdMap[extractedId]

		if exists {
			// pass
		} else {
			continue
		}

		resOrDsType := func() string {
			return string(me.resourceType)
		}
		if resource.IsReferencedAsDataSource() {
			resOrDsType = func() string {
				return string(me.resourceType.AsDataSource())
			}
			replacePattern = "${data.%s.%s.legacy_id}"
			if environment.Flags.Flat {
				replacePattern = "${data.%s.%s.legacy_id}"
			}
		} else {
			if ATOMIC_DEPENDENCIES {
				replacePattern = "${var.%s_%s.value.legacy_id}"
			} else {
				replacePattern = "${var.%s.%s.legacy_id}"
			}
			if environment.Flags.Flat {
				replacePattern = "${%s.%s.legacy_id}"
			}
		}

		s = strings.ReplaceAll(s, resource.LegacyID, fmt.Sprintf(replacePattern, resOrDsType(), resource.UniqueName))
		resources = append(resources, resource)
	}

	return s, resources
}

type optimizedIdDep struct {
	regex            *regexp.Regexp
	containsFunction func(string, string) bool
}

var RegexMutex *sync.Mutex = new(sync.Mutex)

const ENTITY_REGEX = "ENTITY_REGEX"
const V1_CONFIG_ID_REGEX = "V1_CONFIG_ID_REGEX"
const CALC_METRIC_REGEX = "CALC_METRIC_REGEX"
const LEGACY_MZ_ID_REGEX = "LEGACY_MZ_ID"
const NONE = "NONE"

var entityExtractionRegex = regexp.MustCompile(`(((?:[A-Z]+_)?(?:[A-Z]+_)?(?:[A-Z]+_)?[A-Z]+)-[0-9A-Z]{16})`)
var v1ConfigIdRegex = regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
var calcMetricKeyRegex = regexp.MustCompile(`calc:[A-Za-z0-9_.-]*`)

var legacyMzIdRegexFmt = `-{0,1}\d{%d,%d}`
var minLenMzId = 0
var maxLenMzId = 0

var OptimizedKeyRegexV1 = map[string]optimizedIdDep{
	V1_CONFIG_ID_REGEX: {v1ConfigIdRegex, func(text string, id string) bool {
		return true
	}},
}
var OptimizedKeyRegexMzLegacy = map[string]optimizedIdDep{}

var OptimizedKeyRegexId = mergeMaps(OptimizedKeyRegexV1, map[string]optimizedIdDep{
	ENTITY_REGEX: {entityExtractionRegex, func(text string, id string) bool {
		entityIdType := string(id[0:(len(id) - 17)])
		return strings.Contains(text, entityIdType)
	}},
	CALC_METRIC_REGEX: {calcMetricKeyRegex, func(text string, id string) bool {
		return strings.Contains(text, "calc:")
	}},
})

func SetOptimizedRegexModule(module *Module) {
	if module.Type == ResourceTypes.ManagementZoneV2 {
		RegexMutex.Lock()
		defer RegexMutex.Unlock()

		for _, resource := range module.Resources {
			updateMinMaxMzId(resource.LegacyID)
		}
	}
}

func SetOptimizedRegexResource(resource *Resource) {
	if resource.Module.Type == ResourceTypes.ManagementZoneV2 {
		RegexMutex.Lock()
		defer RegexMutex.Unlock()

		updateMinMaxMzId(resource.LegacyID)

	}
}

func updateMinMaxMzId(id string) {

	len := len(id)
	if []byte(id)[0] == '-' {
		len -= 1
	}

	updated := false

	if len < minLenMzId || minLenMzId == 0 {
		minLenMzId = len
		updated = true
	}
	if len > maxLenMzId || maxLenMzId == 0 {
		maxLenMzId = len
		updated = true
	}

	if updated {
		OptimizedKeyRegexMzLegacy = map[string]optimizedIdDep{
			LEGACY_MZ_ID_REGEX: {regexp.MustCompile(fmt.Sprintf(legacyMzIdRegexFmt, minLenMzId, maxLenMzId)), func(text string, id string) bool {
				return true
			}},
		}
	}
}

type iddep struct {
	resourceType ResourceType
	quoted       bool
}

func (me *iddep) IsParent() bool {
	return false
}

func (me *iddep) ResourceType() ResourceType {
	return me.resourceType
}

func (me *iddep) DataSourceType() DataSourceType {
	return ""
}

func (me *iddep) Replace(environment *Environment, s string, replacingIn ResourceType, resourceId string) (string, []any) {
	childDescriptor := environment.Module(replacingIn).GetDescriptor()
	isParent := !environment.ChildResourceOverride && childDescriptor.Parent != nil && string(*childDescriptor.Parent) == string(me.resourceType)

	resources := []any{}

	if len(environment.Module(me.resourceType).Resources) == 0 {
		return s, resources
	}

	idRegexType, isOptimized := environment.Module(me.resourceType).GetDependencyOptimizationInfo()
	extractedIds := environment.Module(replacingIn).Resources[resourceId].GetExtractedIdsPerRegexType(idRegexType, s, OptimizedKeyRegexId)

	if isOptimized && len(extractedIds) < 1 {
		return s, resources
	}

	var replacePattern string
	if ATOMIC_DEPENDENCIES {
		replacePattern = "${var.%s_%s.value.id}"
	} else {
		replacePattern = "${var.%s.%s.id}"
	}
	if environment.Flags.Flat || isParent {
		replacePattern = "${%s.%s.id}"
	}

	getReplaceValues := func(id string) (string, string) {
		valueToReplace := id
		newValue := replacePattern
		if me.quoted {
			valueToReplace = "\"" + id + "\""
			newValue = "\"" + replacePattern + "\""
		}

		return valueToReplace, newValue
	}

	applyReplacement := func(resource *Resource, id string) {

		resOrDsType := func() string {
			return string(me.resourceType)
		}
		if resource.Type == replacingIn {
			replacePattern = "${%s.%s.id}"
		} else {
			if resource.IsReferencedAsDataSource() {
				resOrDsType = func() string {
					return string(me.resourceType.AsDataSource())
				}
				replacePattern = "${data.%s.%s.id}"
				if environment.Flags.Flat {
					replacePattern = "${data.%s.%s.id}"
				}
			} else {
				if ATOMIC_DEPENDENCIES {
					replacePattern = "${var.%s_%s.value.id}"
				} else {
					replacePattern = "${var.%s.%s.id}"
				}
				if environment.Flags.Flat || isParent {
					replacePattern = "${%s.%s.id}"
				}
			}
		}

		valueToReplace, newValue := getReplaceValues(id)
		s = strings.ReplaceAll(s, valueToReplace, fmt.Sprintf(newValue, resOrDsType(), resource.UniqueName))

		resources = append(resources, resource)
	}

	if isOptimized {

		for extractedId := range extractedIds {

			resource, exists := environment.Module(me.resourceType).Resources[extractedId]

			if exists {
				applyReplacement(resource, extractedId)
				delete(environment.Module(replacingIn).Resources[resourceId].ExtractedIdsPerDependencyModule[idRegexType], extractedId)
			}

		}

	} else {

		for id, resource := range environment.Module(me.resourceType).Resources {

			valueToReplace, _ := getReplaceValues(id)
			if strings.Contains(s, valueToReplace) {
				applyReplacement(resource, id)
			}

		}
	}

	return s, resources
}

type resourceIDDep struct {
	resourceType ResourceType
	parent       bool
}

func (me *resourceIDDep) IsParent() bool {
	return me.parent
}

func (me *resourceIDDep) ResourceType() ResourceType {
	return me.resourceType
}

func (me *resourceIDDep) DataSourceType() DataSourceType {
	return ""
}

func (me *resourceIDDep) Replace(environment *Environment, s string, replacingIn ResourceType, resourceId string) (string, []any) {
	replacePattern := "${%s.%s.id}"
	resources := []any{}

	resourcesMap := map[string]*Resource{}

	resource, found := environment.Module(me.resourceType).Resources[resourceId]
	if found {
		resourcesMap[resourceId] = resource
	} else {
		return s, resources
	}

	for id, resource := range resourcesMap {
		resOrDsType := func() string {
			return string(me.resourceType)
		}
		if resource.Type == replacingIn {
			replacePattern = "${%s.%s.id}"
		} else {
			if resource.IsReferencedAsDataSource() {
				resOrDsType = func() string {
					return string(me.resourceType.AsDataSource())
				}
				replacePattern = "${data.%s.%s.id}"
				if environment.Flags.Flat {
					replacePattern = "${data.%s.%s.id}"
				}
			} else {
				replacePattern = "${%s.%s.id}"
				if environment.Flags.Flat {
					replacePattern = "${%s.%s.id}"
				}
			}
		}
		found := false
		if strings.Contains(s, id) {
			s = strings.ReplaceAll(s, id, fmt.Sprintf(replacePattern, resOrDsType(), resource.UniqueName))
			found = true
		}
		if found {

			resources = append(resources, resource)
		}
	}
	return s, resources
}

type tenantds struct {
}

func (me *tenantds) IsParent() bool {
	return false
}

func (me *tenantds) ResourceType() ResourceType {
	return ""
}

func (me *tenantds) DataSourceType() DataSourceType {
	return ""
}

func (me *tenantds) Replace(environment *Environment, s string, replacingIn ResourceType, resourceId string) (string, []any) {
	tenantID := environment.TenantID()
	if len(tenantID) == 0 {
		return s, []any{}
	}
	// when running on HTTP Cache no data sources should get replaced
	// The IDs of these entities are guaranteed to match existing ones
	if len(os.Getenv("DYNATRACE_MIGRATION_CACHE_FOLDER")) > 0 {
		return s, []any{}
	}
	if !strings.Contains(s, tenantID) {
		return s, []any{}
	}
	environment.Module(replacingIn).DataSource("tenant", DataSourceKindTenant)
	s = strings.ReplaceAll(s, tenantID, "${data.dynatrace_tenant.tenant.id}")
	return s, []any{true}
}

type policyds struct {
}

func (me *policyds) IsParent() bool {
	return false
}

func (me *policyds) ResourceType() ResourceType {
	return ""
}

func (me *policyds) DataSourceType() DataSourceType {
	return ""
}

func (me *policyds) Replace(environment *Environment, s string, replacingIn ResourceType, resourceId string) (string, []any) {
	// when running on HTTP Cache no data sources should get replaced
	// The IDs of these entities are guaranteed to match existing ones
	if len(os.Getenv("DYNATRACE_MIGRATION_CACHE_FOLDER")) > 0 {
		return s, []any{}
	}
	if environment.Flags.FlagMigrationOutput {
		return s, []any{}
	}
	found := false
	pattern := `[0-9a-zA-Z]{8}\-[0-9a-zA-Z]{4}\-[0-9a-zA-Z]{4}\-[0-9a-zA-Z]{4}\-[0-9a-zA-Z]{12}#\-\#global#\-\#global`
	m1 := regexp.MustCompile(pattern)
	s = m1.ReplaceAllStringFunc(s, func(id string) string {
		dataSource := environment.Module(replacingIn).DataSource(id, DataSourcePolicy)
		if dataSource == nil {
			return s
		}
		found = true
		return fmt.Sprintf("${data.dynatrace_iam_policy.%s.id}", dataSource.UniqueName)
	})
	if found {
		return s, []any{found}
	} else {
		return s, []any{}
	}
}

type entityds struct {
	Type     string
	Pattern  string
	Coalesce bool
}

func (me *entityds) IsParent() bool {
	return false
}

func (me *entityds) ResourceType() ResourceType {
	return ""
}

func (me *entityds) DataSourceType() DataSourceType {
	return ""
}

func (me *entityds) Replace(environment *Environment, s string, replacingIn ResourceType, resourceId string) (string, []any) {
	// when running on HTTP Cache no data sources should get replaced
	// The IDs of these entities are guaranteed to match existing ones
	if len(os.Getenv("DYNATRACE_MIGRATION_CACHE_FOLDER")) > 0 {
		return s, []any{}
	}
	if environment.Flags.FlagMigrationOutput {
		return s, []any{}
	}
	found := false
	m1 := regexp.MustCompile(me.Pattern)
	s = m1.ReplaceAllStringFunc(s, func(id string) string {
		dataSource := environment.Module(replacingIn).DataSource(id, DataSourceKindEntity)
		if dataSource == nil {
			return s
		}
		found = true
		if me.Coalesce {
			return fmt.Sprintf(`${coalesce(data.dynatrace_entity.%s.id, "%s-0000000000000000")}`, dataSource.ID, me.Type)
		}
		return fmt.Sprintf("${data.dynatrace_entity.%s.id}", dataSource.ID)
	})
	if found {
		return s, []any{found}
	} else {
		return s, []any{}
	}
}

type reqAttName struct {
	resourceType ResourceType
}

func (me *reqAttName) IsParent() bool {
	return false
}

func (me *reqAttName) ResourceType() ResourceType {
	return me.resourceType
}

func (me *reqAttName) DataSourceType() DataSourceType {
	return ""
}

func (me *reqAttName) Replace(environment *Environment, s string, replacingIn ResourceType, resourceId string) (string, []any) {
	var replacePattern string
	if ATOMIC_DEPENDENCIES {
		replacePattern = "${var.%s_%s.value.name}"
	} else {
		replacePattern = "${var.%s.%s.name}"
	}
	if environment.Flags.Flat {
		replacePattern = "${%s.%s.name}"
	}
	resources := []any{}
	for _, resource := range environment.Module(me.resourceType).Resources {
		resOrDsType := func() string {
			return string(me.resourceType)
		}
		if resource.Type == replacingIn {
			replacePattern = "${%s.%s.name}"
		} else {
			if resource.IsReferencedAsDataSource() {
				resOrDsType = func() string {
					return string(me.resourceType.AsDataSource())
				}
				replacePattern = "${data.%s.%s.name}"
				if environment.Flags.Flat {
					replacePattern = "${data.%s.%s.name}"
				}
			} else {
				if ATOMIC_DEPENDENCIES {
					replacePattern = "${var.%s_%s.value.name}"
				} else {
					replacePattern = "${var.%s.%s.name}"
				}
				if environment.Flags.Flat {
					replacePattern = "${%s.%s.name}"
				}
			}
		}
		found := false

		escaped := regexp.QuoteMeta(resource.Name)

		m1 := regexp.MustCompile(fmt.Sprintf("request_attribute(.*)=(.*)\"%s\"", escaped))
		replaced := m1.ReplaceAllString(s, fmt.Sprintf("request_attribute$1=$2\"%s\"", fmt.Sprintf(replacePattern, resOrDsType(), resource.UniqueName)))
		if replaced != s {
			s = replaced
			found = true
		}
		m1 = regexp.MustCompile(fmt.Sprintf("{RequestAttribute:%s}", escaped))
		replaced = m1.ReplaceAllString(s, fmt.Sprintf("{RequestAttribute:%s}", fmt.Sprintf(replacePattern, resOrDsType(), resource.UniqueName)))
		if replaced != s {
			s = replaced
			found = true
		}
		if found {
			resources = append(resources, resource)
		}
	}
	return s, resources
}

func mergeMaps(maps ...map[string]optimizedIdDep) map[string]optimizedIdDep {
	mergedMap := make(map[string]optimizedIdDep)

	for _, m := range maps {
		for key, value := range m {
			mergedMap[key] = value
		}
	}

	return mergedMap
}
