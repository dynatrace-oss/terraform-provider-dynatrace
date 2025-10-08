package list

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
)

const (
	NUM_ENTITY_TYPES = 50
	NUM_CUSTOM_TAGS  = 50
	NUM_ENTITIES     = 50
)

type TestClient struct {
	latency     time.Duration // artificial per-request delay inside Finish
	inFlight    int64         // current in-flight
	maxInFlight int64         // max observed in-flight
	total       int64         // total Finish calls
}

func NewTestClient(latency time.Duration) *TestClient {
	return &TestClient{latency: latency}
}

func (client *TestClient) Get(ctx context.Context, url string, expectedStatusCodes ...int) rest.Request {
	return &Request{client: client, ctx: ctx}
}
func (client *TestClient) Post(ctx context.Context, url string, payload any, expectedStatusCodes ...int) rest.Request {
	panic("unsupported operation")
}
func (client *TestClient) Put(ctx context.Context, url string, payload any, expectedStatusCodes ...int) rest.Request {
	panic("unsupported operation")
}
func (client *TestClient) Delete(ctx context.Context, url string, expectedStatusCodes ...int) rest.Request {
	panic("unsupported operation")
}
func (client *TestClient) Credentials() *rest.Credentials {
	panic("unsupported operation")
}

// Accessors for assertions
func (c *TestClient) MaxConcurrent() int64 { return atomic.LoadInt64(&c.maxInFlight) }
func (c *TestClient) Total() int64         { return atomic.LoadInt64(&c.total) }
func (c *TestClient) InFlight() int64      { return atomic.LoadInt64(&c.inFlight) }

type Request struct {
	client *TestClient
	ctx    context.Context
}

func (request *Request) Finish(v ...any) error {
	c := request.client

	// increment in-flight
	cur := atomic.AddInt64(&c.inFlight, 1)

	// update maxInFlight if needed (lock-free CAS loop)
	for {
		max := atomic.LoadInt64(&c.maxInFlight)
		if cur <= max {
			break
		}
		if atomic.CompareAndSwapInt64(&c.maxInFlight, max, cur) {
			break
		}
	}

	// count total requests
	atomic.AddInt64(&c.total, 1)
	defer func() {
		atomic.AddInt64(&c.inFlight, -1)
	}()

	// simulate work/latency so goroutines overlap
	if c.latency > 0 {
		select {
		case <-time.After(c.latency):
		case <-request.ctx.Done():
			return request.ctx.Err()
		}
	}

	if v == nil {
		panic("unsupported operation for v == nil")
	}
	if len(v) == 0 {
		panic("unsupported operation for empty v")
	}
	val := v[0]
	if val == nil {
		panic("unsupported operation for v[0] == nil")
	}

	switch typedValue := val.(type) {
	case *GetEntityTypesResponse:
		typedValue.TotalCount = NUM_ENTITY_TYPES
		typedValue.NextPageKey = ""
		for i := 0; i < typedValue.TotalCount; i++ {
			typedValue.Types = append(typedValue.Types, EntityType{Type: fmt.Sprintf("entity-type-%d", i)})
		}
		return nil
	case *GetCustomTagsResponse:
		typedValue.TotalCount = NUM_CUSTOM_TAGS
		for i := 0; i <= typedValue.TotalCount; i++ {
			typedValue.Tags = append(typedValue.Tags, Tag{
				StringRepresentation: fmt.Sprintf("tag.%d", i),
				Value:                fmt.Sprintf("tag.value.%d", i),
				Key:                  fmt.Sprintf("tag.key.%d", i),
				Context:              fmt.Sprintf("tag.context.%d", i),
			})
		}
		return nil
	case *GETEntitiesResponse:
		typedValue.TotalCount = NUM_ENTITIES
		typedValue.NextPageKey = ""
		for i := 0; i <= typedValue.TotalCount; i++ {
			entity := Entity{ID: fmt.Sprintf("entity.id.%d", i), Name: fmt.Sprintf("entity.name.%d", i)}
			for j := 0; j <= typedValue.TotalCount; j++ {
				entity.Tags = append(entity.Tags, Tag{
					StringRepresentation: fmt.Sprintf("entity.%d.tag.%d", i, j),
					Value:                fmt.Sprintf("entity.%d.tag.value.%d", i, j),
					Key:                  fmt.Sprintf("entity.%d.tag.key.%d", i, j),
					Context:              fmt.Sprintf("entity.%d.tag.context.%d", i, j),
				})
			}
			typedValue.Entities = append(typedValue.Entities, entity)
		}
		return nil
	default:
		panic(fmt.Sprintf("unsupported operation for v type %T", val))
	}
}

func (request *Request) Expect(codes ...int) rest.Request {
	panic("unsupported operation")

}
func (request *Request) OnResponse(onresponse func(resp *http.Response)) rest.Request {
	panic("unsupported operation")
}

func (me *Request) SetHeader(name string, value string) {
	panic("unsupported operation")
}

func TestCustomTagExport(t *testing.T) {
	client := NewTestClient(time.Millisecond * 2)
	_, err := List(t.Context(), client)
	if err != nil {
		t.Error(err)
		return
	}
	if got := client.MaxConcurrent(); got > DefaultMaxConcurrent {
		t.Fatalf("saw %d concurrent requests; want <= 8", got)
	}
}
