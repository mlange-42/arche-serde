package archeserde

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

// CompB is an example component
type CompB struct {
	FloatValue float32 `json:"floatValue"`
}

func TestSerDe(t *testing.T) {
	world := ecs.NewWorld()

	compAId := ecs.ComponentID[CompA](&world)
	//compBId := ecs.ComponentID[CompB](&world)

	_ = world.NewEntityWith(ecs.Component{ID: compAId, Comp: &CompA{IntValue: 100, BoolValue: true}})

	tp := reflect.TypeOf((*CompA)(nil)).Elem()
	var jsonComp []byte

	query := world.Query(ecs.All())
	for query.Next() {
		data := query.Get(compAId)
		elem := reflect.NewAt(tp, data)
		interf := elem.Interface()
		jsonData, err := json.Marshal(interf)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}
		fmt.Printf("json data: %s\n", jsonData)
		jsonComp = jsonData

		newElem := reflect.New(tp)
		newInterf := newElem.Interface()

		if err = json.Unmarshal(jsonData, &newInterf); err != nil {
			fmt.Printf("could not unmarshal json: %s\n", err)
			return
		}

		fmt.Printf("decoded: %s %v\n", reflect.TypeOf(newInterf), newInterf)
	}

	newElem := reflect.New(tp)
	newInterf := newElem.Interface()

	if err := json.Unmarshal(jsonComp, &newInterf); err != nil {
		fmt.Printf("could not unmarshal json: %s\n", err)
		return
	}
	fmt.Printf("decoded: %s %v\n", reflect.TypeOf(newInterf), newInterf)

	entity := world.NewEntityWith(ecs.Component{ID: compAId, Comp: newInterf})
	compA := (*CompA)(world.Get(entity, compAId))

	fmt.Printf("added component: %v\n", compA)

	assert.Equal(t, *compA, CompA{IntValue: 100, BoolValue: true})
}
