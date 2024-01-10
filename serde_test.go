package archeserde_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/stretchr/testify/assert"
)

// CompA is an example component
type CompA struct {
	IntValue  int  `json:"intValue"`
	BoolValue bool `json:"boolValue"`
}

func TestSerDe(t *testing.T) {
	world := ecs.NewWorld()

	compAId := ecs.ComponentID[CompA](&world)

	e1 := world.NewEntityWith(ecs.Component{ID: compAId, Comp: &CompA{IntValue: 100, BoolValue: true}})

	// Arche will need to expose GetType(compId), as component types will be unknown.
	tp := reflect.TypeOf((*CompA)(nil)).Elem()

	// Just an example. There are world.Mask(entity) and query.Mask(entity) to get all components.
	// Potentially, world.ComponentIds(entity) and query.ComponentIds(entity) could be added.
	data := world.Get(e1, compAId)
	value := reflect.NewAt(tp, data).Interface()

	fmt.Printf("original comp: %s %v\n", reflect.TypeOf(value), value)

	jsonData, err := json.Marshal(value)
	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)
		return
	}
	fmt.Printf("json data: %s\n", jsonData)

	newValue := reflect.New(tp).Interface()

	if err = json.Unmarshal(jsonData, &newValue); err != nil {
		fmt.Printf("could not unmarshal json: %s\n", err)
		return
	}

	fmt.Printf("decoded: %s %v\n", reflect.TypeOf(newValue), newValue)

	e2 := world.NewEntityWith(ecs.Component{ID: compAId, Comp: newValue})
	compA := (*CompA)(world.Get(e2, compAId))

	fmt.Printf("added component: %v\n", compA)

	assert.Equal(t, *compA, CompA{IntValue: 100, BoolValue: true})
}
