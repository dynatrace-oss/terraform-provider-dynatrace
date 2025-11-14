//go:build unit

package cache_test

import (
	"errors"
	"os"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/cache/tar"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func hide(v any) {}

var testdata = map[string]string{
	"eins": uuid.NewString(),
	"zwei": uuid.NewString(),
	"drei": uuid.NewString(),
	"vier": uuid.NewString(),
}

func TestTarFolder(t *testing.T) {
	os.Remove("reini-war-data.tar")
	folder, _, err := tar.New("reini-war-data")
	hide(folder)
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		os.Remove("reini-war-data.tar")
	}()
	for k, v := range testdata {
		if err := folder.Save(api.Stub{ID: k, Name: k}, []byte(v)); err != nil {
			t.Error(err)
			return
		}
	}
	folder, _, err = tar.New("reini-war-data")
	if err != nil {
		t.Error(err)
		return
	}
	for k, v := range testdata {
		if t.Failed() {
			break
		}
		stub, data, err := folder.Get(k)
		if err != nil {
			t.Error(err)
			return
		}
		assert.Equal(t, k, stub.ID)
		assert.Equal(t, k, stub.Name)
		assert.Equal(t, v, string(data))
	}
	if err := folder.Delete("eins"); err != nil {
		t.Error(err)
		return
	}
	stub, data, err := folder.Get("eins")
	if err != nil {
		t.Error(err)
		return
	}
	if data != nil {
		t.Error(errors.New("data should was expected to be nil"))
		return
	}
	if stub != nil {
		t.Error(errors.New("stub should was expected to be nil"))
		return
	}
}
