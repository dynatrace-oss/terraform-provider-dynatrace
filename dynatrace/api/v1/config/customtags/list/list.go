package list

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
)

const TIME_FRAME = "now-3y"

type EntityTags struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Tags []Tag  `json:"tags"`
}

func List(client rest.Client) (api.Stubs, error) {
	entityTypes, err := GETEntityTypes(client)
	if err != nil {
		return nil, err
	}
	channels := []chan EntityTags{}
	for _, entityType := range entityTypes {
		results := make(chan EntityTags, 100)
		channels = append(channels, results)
		go GetEntities(entityType, client, results)
	}
	m := map[string]EntityTags{}
	for _, channel := range channels {
		for entityTag := range channel {
			m[entityTag.ID] = entityTag
			// fmt.Println(entityTag.ID)
			// for _, tag := range entityTag.Tags {
			// 	fmt.Println("  -", tag.StringRepresentation)
			// }
		}
	}
	var stubs api.Stubs
	for k, v := range m {
		stubs = append(stubs, &api.Stub{ID: fmt.Sprintf(`entityId(%s)`, k), Name: v.Name})
	}
	return stubs, nil
}

func GetEntities(entityType EntityType, client rest.Client, results chan EntityTags) {
	defer close(results)
	tags, err := entityType.GetCustomTags(client)
	if err != nil {
		panic(err)
	}
	if len(tags) == 0 {
		return
	}
	chanEntities := make(chan Entity, 100)
	go GETEntitiesWithTags(entityType.Type, tags, client, chanEntities)
	for entity := range chanEntities {
		if len(entity.Tags) == 0 {
			continue
		}
		chanEntityTags := make(chan EntityTags, 100)
		go GetCustomTags(entity, client, chanEntityTags)
		for entityTag := range chanEntityTags {
			results <- entityTag
		}
	}
}

func GetCustomTags(entity Entity, client rest.Client, results chan EntityTags) {
	defer close(results)
	tags, err := GETCustomTags(entity.ID, client)
	if err != nil {
		panic(err)
	}
	if len(tags) > 0 {
		results <- EntityTags{ID: entity.ID, Name: entity.Name, Tags: tags}
	}
}
