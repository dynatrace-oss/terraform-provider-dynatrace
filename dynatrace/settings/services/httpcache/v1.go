package httpcache

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/cache/tar"
)

type GetV1 struct {
	SchemaID string
	ID       string
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
			return &rest.Error{Code: 404, Message: fmt.Sprintf("%s not found", me.ID)}
		}
		wrapper := struct {
			Downloaded struct {
				Value json.RawMessage `json:"value"`
			} `json:"downloaded"`
		}{}
		if err := json.Unmarshal(data, &wrapper); err != nil {
			return err
		}
		if err := json.Unmarshal(wrapper.Downloaded.Value, &v); err != nil {
			return err
		}
		return nil
	}
	return &rest.Error{Code: 404, Message: fmt.Sprintf("%s not found", me.ID)}
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
