package httpcache

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/cache/tar"
)

type GetSettings20Request struct {
	SchemaID string
	ID       string
}

func (me *GetSettings20Request) Raw() ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (me *GetSettings20Request) Finish(vs ...any) error {
	var v any
	if len(vs) > 0 {
		v = vs[0]
	}

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
			return &rest.Error{Code: 404, Message: fmt.Sprintf("%s not found", me.ID)}
		}
		wrapper := struct {
			Downloaded json.RawMessage `json:"downloaded"`
		}{}
		if err := json.Unmarshal(data, &wrapper); err != nil {
			return err
		}
		if err := json.Unmarshal(wrapper.Downloaded, &v); err != nil {
			return err
		}
		return nil
	}
	return &rest.Error{Code: 404, Message: fmt.Sprintf("%s not found", me.ID)}
}

func (me *GetSettings20Request) Expect(codes ...int) rest.Request {
	return me
}

func (me *GetSettings20Request) Payload(any) rest.Request {
	return me
}

func (me *GetSettings20Request) OnResponse(func(resp *http.Response)) rest.Request {
	return me
}

type ListSettings20Request struct {
	SchemaID string
}

func (me *ListSettings20Request) Raw() ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (me *ListSettings20Request) Finish(vs ...any) error {
	var v any
	if len(vs) > 0 {
		v = vs[0]
	}
	sol := SettingsObjectList{Items: []*SettingsObjectListItem{}}

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
			_, data, err := tarFolder.Get(stub.ID)
			if err != nil {
				return err
			}
			wrapper := struct {
				Downloaded SettingsObject `json:"downloaded"`
			}{}
			if err := json.Unmarshal(data, &wrapper); err != nil {
				return err
			}
			if err := json.Unmarshal(wrapper.Downloaded.Value, &v); err != nil {
				return err
			}
			sol.Items = append(sol.Items, &SettingsObjectListItem{ObjectID: stub.ID, SchemaVersion: wrapper.Downloaded.SchemaVersion, Scope: wrapper.Downloaded.Scope, Value: wrapper.Downloaded.Value})
		}
	}
	data, err := json.Marshal(sol)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	return nil
}

func (me *ListSettings20Request) Expect(codes ...int) rest.Request {
	return me
}

func (me *ListSettings20Request) Payload(any) rest.Request {
	return me
}

func (me *ListSettings20Request) OnResponse(func(resp *http.Response)) rest.Request {
	return me
}

type SettingsObjectList struct {
	Items       []*SettingsObjectListItem `json:"items"`
	NextPageKey *string                   `json:"nextPageKey,omitempty"`
}

type SettingsObjectListItem struct {
	ObjectID      string          `json:"objectId"`
	Scope         string          `json:"scope"`
	SchemaVersion string          `json:"schemaVersion"`
	Value         json.RawMessage `json:"value"`
}

type SettingsObject struct {
	SchemaVersion string          `json:"schemaVersion"`
	SchemaID      string          `json:"schemaId"`
	Scope         string          `json:"scope"`
	Value         json.RawMessage `json:"value"`
}
