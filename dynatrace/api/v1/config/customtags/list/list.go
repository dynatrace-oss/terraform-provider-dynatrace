package list

import (
	"context"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
)

const (
	TIME_FRAME                                        = "now-3y"
	DYNATRACE_MAX_CONCURRENT_CUSTOM_TAG_LIST_REQUESTS = "DYNATRACE_MAX_CONCURRENT_CUSTOM_TAG_LIST_REQUESTS"
	DefaultMaxConcurrent                              = 4
	MaxConcurrentMinValue                             = 1
	MaxConcurrentMaxValue                             = 20
)

type EntityTags struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Tags []Tag  `json:"tags"`
}

func List(ctx context.Context, client rest.Client) (api.Stubs, error) {
	entityTypes, err := GETEntityTypes(ctx, client)
	if err != nil {
		return nil, err
	}

	maxConcurrent := settings.GetIntEnvCtx(ctx, DYNATRACE_MAX_CONCURRENT_CUSTOM_TAG_LIST_REQUESTS, DefaultMaxConcurrent, MaxConcurrentMinValue, MaxConcurrentMaxValue)

	// limiter shared by *all* downstream calls using the same client
	limiter := make(chan struct{}, maxConcurrent)

	channels := []chan EntityTags{}
	for _, entityType := range entityTypes {
		results := make(chan EntityTags, 100)
		channels = append(channels, results)
		go GetEntities(ctx, entityType, client, limiter, results)
	}
	m := map[string]EntityTags{}
	for _, channel := range channels {
		for entityTag := range channel {
			m[entityTag.ID] = entityTag
		}
	}
	var stubs api.Stubs
	for k, v := range m {
		stubs = append(stubs, &api.Stub{ID: fmt.Sprintf(`entityId(%s)`, k), Name: v.Name})
	}
	return stubs, nil
}

// acquire blocks when maxConcurrent is reached; release frees a slot
func acquire(limiter chan struct{}) { limiter <- struct{}{} }
func release(limiter chan struct{}) { <-limiter }

func GetEntities(ctx context.Context, entityType EntityType, client rest.Client, limiter chan struct{}, results chan EntityTags) {
	defer close(results)

	// Limit this remote call
	acquire(limiter)
	tags, err := entityType.GetCustomTags(ctx, client)
	release(limiter)
	if err != nil {
		logging.Debug.Info.Printf("GetCustomTags(entityType=%s): %v", entityType.Type, err)
		return
	}
	if len(tags) == 0 {
		return
	}
	chanEntities := make(chan Entity, 100)
	// Limit this remote call (spawns its own producer goroutine)
	go func() {
		acquire(limiter)
		defer release(limiter)
		GETEntitiesWithTags(ctx, entityType.Type, tags, client, chanEntities)
	}()

	for entity := range chanEntities {
		if len(entity.Tags) == 0 {
			continue
		}
		chanEntityTags := make(chan EntityTags, 100)
		go GetCustomTags(ctx, entity, client, limiter, chanEntityTags)
		for entityTag := range chanEntityTags {
			results <- entityTag
		}
	}
}

func GetCustomTags(ctx context.Context, entity Entity, client rest.Client, limiter chan struct{}, results chan EntityTags) {
	defer close(results)
	acquire(limiter)
	tags, err := GETCustomTags(ctx, entity.ID, client)
	release(limiter)
	if err != nil {
		logging.Debug.Info.Printf("GETCustomTags(entity=%s): %v", entity.ID, err)
		return
	}
	if len(tags) > 0 {
		results <- EntityTags{ID: entity.ID, Name: entity.Name, Tags: tags}
	}
}
