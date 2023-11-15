package httpcache

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/address"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	dashboards "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/dashboards/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/cache/tar"
)

type GetV1 struct {
	SchemaID        string
	ServiceSchemaID string
	ID              string
}

func (me *GetV1) Finish(v any) error {
	tarFolder, _, err := tar.NewExisting(CACHE_FOLDER + "/" + strings.ReplaceAll(me.SchemaID, ":", ""))
	if err != nil {
		return err
	}

	if tarFolder != nil {
		stub, data, err := tarFolder.Get(me.ID)
		if err != nil {
			return err
		}
		if stub == nil {
			return rest.Error{Code: 404, Message: fmt.Sprintf("V1_Tar %s not found for %s", me.ID, me.SchemaID)}
		}
		wrapper := struct {
			Downloaded struct {
				ClassidID string          `json:"classicId,omitempty"`
				Value     json.RawMessage `json:"value"`
			} `json:"downloaded"`
		}{}
		if err := json.Unmarshal(data, &wrapper); err != nil {
			return err
		}
		if err := json.Unmarshal(wrapper.Downloaded.Value, &v); err != nil {
			return err
		}

		address.AddToOriginal(address.AddressOriginal{
			TerraformSchemaID: me.ServiceSchemaID,
			OriginalID:        wrapper.Downloaded.ClassidID,
			OriginalSchemaID:  me.SchemaID,
		})
		return nil
	}
	return rest.Error{Code: 404, Message: fmt.Sprintf("V1_Tar Nil %s not found for %s", me.ID, me.SchemaID)}
}

type ListV1 struct {
	SchemaID string
}

func (me *ListV1) Finish(v any) error {

	stubList := &api.StubList{Values: []*api.Stub{}}

	tarFolder, _, err := tar.NewExisting(CACHE_FOLDER + "/" + strings.ReplaceAll(me.SchemaID, ":", ""))
	if err != nil {
		return err
	}

	if tarFolder != nil {
		stubs, err := tarFolder.List()
		if err != nil {
			return err
		}

		for _, stub := range stubs {
			stubList.Values = append(stubList.Values, stub)
		}
	}
	data, err := json.Marshal(stubList)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	return nil
}

type ListMonitorsV1 struct {
	Prefix string
}

func (me *ListMonitorsV1) Finish(v any) error {
	stubList := &api.StubList{Values: []*api.Stub{}}

	tarFolder, _, err := tar.NewExisting(CACHE_FOLDER + "/synthetic-monitor")
	if err != nil {
		return err
	}

	if tarFolder != nil {
		stubs, err := tarFolder.List()
		if err != nil {
			return err
		}

		for _, stub := range stubs {
			if strings.HasPrefix(stub.ID, me.Prefix) {
				stubList.Values = append(stubList.Values, stub)
			}
		}
	}

	monitors := []any{}
	if stubList != nil && len(stubList.Values) > 0 {
		for _, stub := range stubList.Values {
			monitors = append(monitors, map[string]any{"entityId": stub.ID, "name": stub.Name})
		}
	}
	data, err := json.Marshal(map[string]any{"monitors": monitors})
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	return nil
}

type ListPrivateSyntheticLocationsV1 struct{}

func (me *ListPrivateSyntheticLocationsV1) Finish(v any) error {
	stubList := &api.StubList{Values: []*api.Stub{}}

	tarFolder, _, err := tar.NewExisting(CACHE_FOLDER + "/synthetic-location")
	if err != nil {
		return err
	}

	if tarFolder != nil {
		stubs, err := tarFolder.List()
		if err != nil {
			return err
		}

		for _, stub := range stubs {
			stubList.Values = append(stubList.Values, stub)
		}
	}

	locations := []any{}
	if stubList != nil && len(stubList.Values) > 0 {
		for _, stub := range stubList.Values {
			locations = append(locations, map[string]any{"entityId": stub.ID, "name": stub.Name})
		}
	}
	data, err := json.Marshal(map[string]any{"locations": locations})
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	return nil
}

type ListDashboardsV1 struct{}

func (me *ListDashboardsV1) Finish(v any) error {
	stubList := &dashboards.DashboardList{Dashboards: []*dashboards.DashboardStub{}}

	tarFolder, _, err := tar.NewExisting(CACHE_FOLDER + "/dashboard")
	if err != nil {
		return err
	}

	if tarFolder != nil {
		stubs, err := tarFolder.List()
		if err != nil {
			return err
		}

		for _, stub := range stubs {
			stubList.Dashboards = append(stubList.Dashboards, &dashboards.DashboardStub{ID: stub.ID, Name: &stub.Name, Owner: &stub.EntityID})
		}
	}

	data, err := json.Marshal(stubList)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	return nil
}

var logged = map[string]string{}

type NonSupportedV1 struct {
	SchemaID string
}

var lc *sync.Mutex = new(sync.Mutex)

func log(msg string) {
	lc.Lock()
	defer lc.Unlock()
	if _, found := logged[msg]; found {
		return
	}
	fmt.Println(msg)
	logged[msg] = msg

}

func (me *NonSupportedV1) Finish(v any) error {
	log(fmt.Sprintf("Schema %s is not supported yet", me.SchemaID))

	return nil
}
