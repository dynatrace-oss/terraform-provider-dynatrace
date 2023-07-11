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
)

// To allow -target to work with dependencies at an atomic level
var ATOMIC_DEPENDENCIES = os.Getenv("DYNATRACE_ATOMIC_DEPENDENCIES") == "true"

type Dependency interface {
	Replace(environment *Environment, s string, replacingIn ResourceType) (string, []any)
	ResourceType() ResourceType
	DataSourceType() DataSourceType
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
	LegacyID             func(resourceType ResourceType) Dependency
	ID                   func(resourceType ResourceType) Dependency
	ResourceID           func(resourceType ResourceType) Dependency
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
	CloudApplicationNamespace Dependency
	EnvironmentActiveGate     Dependency
	Tenant                    Dependency
}{
	ManagementZone:       &mgmzdep{ResourceTypes.ManagementZoneV2},
	LegacyID:             func(resourceType ResourceType) Dependency { return &legacyID{resourceType} },
	ID:                   func(resourceType ResourceType) Dependency { return &iddep{resourceType} },
	ResourceID:           func(resourceType ResourceType) Dependency { return &resourceIDDep{resourceType} },
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
	CloudApplicationNamespace: &entityds{"CLOUD_APPLICATION_NAMESPACE", "CLOUD_APPLICATION_NAMESPACE-[A-Z0-9]{16}", false},
	EnvironmentActiveGate:     &entityds{"ENVIRONMENT_ACTIVE_GATE", "ENVIRONMENT_ACTIVE_GATE-[A-Z0-9]{16}", false},
	Tenant:                    &tenantds{},
}

type mgmzdep struct {
	resourceType ResourceType
}

func (me *mgmzdep) ResourceType() ResourceType {
	return me.resourceType
}

func (me *mgmzdep) DataSourceType() DataSourceType {
	return ""
}

func (me *mgmzdep) Replace(environment *Environment, s string, replacingIn ResourceType) (string, []any) {
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

		m1 := regexp.MustCompile(fmt.Sprintf(`"managementZone": {([\S\s]*)"name":(.*)"%s"([\S\s]*)}`, resource.Name))
		replaced := m1.ReplaceAllString(s, fmt.Sprintf(`"managementZone": {$1"name":$2"%s"$3}`, fmt.Sprintf(replacePattern, resOrDsType(), resource.UniqueName)))
		if replaced != s {
			s = replaced
			found = true
		}
		m1 = regexp.MustCompile(fmt.Sprintf(`management_zones = \[(.*)\"%s\"(.*)\]`, resource.Name))
		replaced = m1.ReplaceAllString(replaced, fmt.Sprintf(`management_zones = [ $1"%s"$2 ]`, fmt.Sprintf(replacePattern, resOrDsType(), resource.UniqueName)))
		if replaced != s {
			s = replaced
			found = true
		}
		m1 = regexp.MustCompile(fmt.Sprintf(`mzName\((.*)"%s"(.*)\)`, resource.Name))
		replaced = m1.ReplaceAllString(replaced, fmt.Sprintf(`mzName($1"%s"$2)`, fmt.Sprintf(replacePattern, resOrDsType(), resource.UniqueName)))
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

type legacyID struct {
	resourceType ResourceType
}

func (me *legacyID) ResourceType() ResourceType {
	return me.resourceType
}

func (me *legacyID) DataSourceType() DataSourceType {
	return ""
}

func (me *legacyID) Replace(environment *Environment, s string, replacingIn ResourceType) (string, []any) {
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
	for _, resource := range environment.Module(me.resourceType).Resources {
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
		if len(resource.LegacyID) > 0 {
			found := false
			if strings.Contains(s, resource.LegacyID) {
				s = strings.ReplaceAll(s, resource.LegacyID, fmt.Sprintf(replacePattern, resOrDsType(), resource.UniqueName))
				found = true
			}
			if found {
				resources = append(resources, resource)
			}
		}
	}
	return s, resources
}

type iddep struct {
	resourceType ResourceType
}

func (me *iddep) ResourceType() ResourceType {
	return me.resourceType
}

func (me *iddep) DataSourceType() DataSourceType {
	return ""
}

func (me *iddep) Replace(environment *Environment, s string, replacingIn ResourceType) (string, []any) {
	childDescriptor := environment.Module(replacingIn).Descriptor
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
	resources := []any{}
	for id, resource := range environment.Module(me.resourceType).Resources {
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

type resourceIDDep struct {
	resourceType ResourceType
}

func (me *resourceIDDep) ResourceType() ResourceType {
	return me.resourceType
}

func (me *resourceIDDep) DataSourceType() DataSourceType {
	return ""
}

func (me *resourceIDDep) Replace(environment *Environment, s string, replacingIn ResourceType) (string, []any) {
	replacePattern := "${%s.%s.id}"
	resources := []any{}
	for id, resource := range environment.Module(me.resourceType).Resources {
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

func (me *tenantds) ResourceType() ResourceType {
	return ""
}

func (me *tenantds) DataSourceType() DataSourceType {
	return ""
}

func (me *tenantds) Replace(environment *Environment, s string, replacingIn ResourceType) (string, []any) {
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
	environment.Module(replacingIn).DataSource("tenant")
	s = strings.ReplaceAll(s, tenantID, "${data.dynatrace_tenant.tenant.id}")
	return s, []any{true}
}

type entityds struct {
	Type     string
	Pattern  string
	Coalesce bool
}

func (me *entityds) ResourceType() ResourceType {
	return ""
}

func (me *entityds) DataSourceType() DataSourceType {
	return ""
}

func (me *entityds) Replace(environment *Environment, s string, replacingIn ResourceType) (string, []any) {
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
		dataSource := environment.Module(replacingIn).DataSource(id)
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

func (me *reqAttName) ResourceType() ResourceType {
	return me.resourceType
}

func (me *reqAttName) DataSourceType() DataSourceType {
	return ""
}

func (me *reqAttName) Replace(environment *Environment, s string, replacingIn ResourceType) (string, []any) {
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
		m1 := regexp.MustCompile(fmt.Sprintf("request_attribute(.*)=(.*)\"%s\"", resource.Name))
		replaced := m1.ReplaceAllString(s, fmt.Sprintf("request_attribute$1=$2\"%s\"", fmt.Sprintf(replacePattern, resOrDsType(), resource.UniqueName)))
		if replaced != s {
			s = replaced
			found = true
		}
		m1 = regexp.MustCompile(fmt.Sprintf("{RequestAttribute:%s}", resource.Name))
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
