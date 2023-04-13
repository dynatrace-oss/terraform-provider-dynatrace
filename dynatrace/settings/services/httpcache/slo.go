package httpcache

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	slo "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/slo/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/cache/tar"
)

type GetSLORequest struct {
	SchemaID string
	ID       string
}

func (me *GetSLORequest) Raw() ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (me *GetSLORequest) Finish(vs ...any) error {
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

func (me *GetSLORequest) Expect(codes ...int) rest.Request {
	return me
}

func (me *GetSLORequest) Payload(any) rest.Request {
	return me
}

func (me *GetSLORequest) OnResponse(func(resp *http.Response)) rest.Request {
	return me
}

type ListSLORequest struct {
	SchemaID string
}

func (me *ListSLORequest) Raw() ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (me *ListSLORequest) Finish(vs ...any) error {
	var v any
	if len(vs) > 0 {
		v = vs[0]
	}
	sol := sloList{SLOs: []*sloListEntry{}}

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
				Downloaded struct {
					Value slo.SLO `json:"value"`
				} `json:"downloaded"`
			}{}
			if err := json.Unmarshal(data, &wrapper); err != nil {
				return err
			}
			sol.SLOs = append(sol.SLOs, &sloListEntry{SLO: wrapper.Downloaded.Value, ID: stub.ID})
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

func (me *ListSLORequest) Expect(codes ...int) rest.Request {
	return me
}

func (me *ListSLORequest) Payload(any) rest.Request {
	return me
}

func (me *ListSLORequest) OnResponse(func(resp *http.Response)) rest.Request {
	return me
}

type sloList struct {
	SLOs        []*sloListEntry `json:"slo"`
	PageSize    *int32          `json:"pageSize"`
	NextPageKey *string         `json:"nextPageKey,omitempty"`
	TotalCount  *int64          `json:"totalCount"`
}

type sloListEntry struct {
	slo.SLO
	ID string `json:"id"`
}
