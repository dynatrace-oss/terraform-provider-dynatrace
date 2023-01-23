package cache_test

import (
	"errors"
	"os"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/cache"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/assert"
	"github.com/google/uuid"
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
	folder, err := cache.NewTarFolder("reini-war-data")
	hide(folder)
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		// os.Remove("reini-war-data.tar")
	}()
	for k, v := range testdata {
		if err := folder.Save(settings.Stub{ID: k, Name: k}, []byte(v)); err != nil {
			t.Error(err)
			return
		}
	}
	folder, err = cache.NewTarFolder("reini-war-data")
	if err != nil {
		t.Error(err)
		return
	}
	assert := assert.New(t)
	for k, v := range testdata {
		if t.Failed() {
			break
		}
		stub, data, err := folder.Get(k)
		if err != nil {
			t.Error(err)
			return
		}
		assert.Equals(k, stub.ID)
		assert.Equals(k, stub.Name)
		assert.Equals(v, string(data))
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
