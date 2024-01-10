package main

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/mlange-42/arche/ecs"
)

// CompA is an example component type
type CompA struct {
	IntValue  int `json:"intValue"`
	BoolValue bool
}

func main() {
	world := ecs.NewWorld()
	compAId := ecs.ComponentID[CompA](&world)

	// Create an entity for testing.
	e1 := world.NewEntityWith(ecs.Component{
		ID:   compAId,
		Comp: &CompA{IntValue: 100, BoolValue: true},
	})

	// Just an example for a single component. There are world.Mask(entity) and query.Mask(entity) to get all component IDs for an entity.
	// Potentially, world.ComponentIds(entity) and query.ComponentIds(entity) could be added.
	component := world.Get(e1, compAId)

	// Arche will need to expose GetType(compId), as component types will be unknown.
	tp := reflect.TypeOf((*CompA)(nil)).Elem()

	// Encode to JSON, without knowing the type
	value := reflect.NewAt(tp, component).Interface()
	fmt.Printf("original comp: %s %v\n", reflect.TypeOf(value), value)

	jsonData, err := json.Marshal(value)
	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)
		return
	}
	fmt.Printf("json data: %s\n", jsonData)

	// Decode the JSON, without knowing the type
	newValue := reflect.New(tp).Interface()

	if err = json.Unmarshal(jsonData, &newValue); err != nil {
		fmt.Printf("could not unmarshal json: %s\n", err)
		return
	}

	fmt.Printf("decoded: %s %v\n", reflect.TypeOf(newValue), newValue)

	// Add the component of unknown type to a new entity
	e2 := world.NewEntityWith(ecs.Component{ID: compAId, Comp: newValue})

	// Get and check the just added component
	compA := (*CompA)(world.Get(e2, compAId))
	fmt.Printf("added component: %v\n", compA)
}
