package archeserde_test

import (
	"fmt"
	"testing"

	archeserde "github.com/mlange-42/arche-serde"
	"github.com/mlange-42/arche/ecs"
	"github.com/stretchr/testify/assert"
)

func TestDeserializeErrors(t *testing.T) {
	world := ecs.NewWorld()
	_ = ecs.ComponentID[Position](&world)
	_ = ecs.ComponentID[ChildOf](&world)
	_ = ecs.ComponentID[ChildRelation](&world)

	err := archeserde.Deserialize([]byte("{xxx}"), &world)
	assert.Contains(t, err.Error(), "invalid character 'x'")

	err = archeserde.Deserialize([]byte(textOk), &world)
	assert.Contains(t, err.Error(), "component type is not registered")

	world.Reset()
	_ = ecs.ComponentID[Velocity](&world)

	err = archeserde.Deserialize([]byte(textOk), &world)
	assert.Contains(t, err.Error(), "resource type is not registered")

	world.Reset()
	_ = ecs.ResourceID[Velocity](&world)

	err = archeserde.Deserialize([]byte(textOk), &world)
	assert.Contains(t, err.Error(), "resource type registered but nil")

	world.Reset()
	_ = ecs.AddResource(&world, &Velocity{})
	err = archeserde.Deserialize([]byte(textOk), &world)
	assert.Nil(t, err)

	world.Reset()
	_ = ecs.AddResource(&world, &Velocity{})
	err = archeserde.Deserialize([]byte(textErrEntities), &world)
	assert.Contains(t, err.Error(), "world has 2 alive entities")

	world.Reset()
	_ = ecs.AddResource(&world, &Velocity{})
	err = archeserde.Deserialize([]byte(textErrTypes), &world)
	assert.Contains(t, err.Error(), "cannot unmarshal object")

	world.Reset()
	_ = ecs.AddResource(&world, &Velocity{})
	err = archeserde.Deserialize([]byte(textErrComponent), &world)
	assert.Contains(t, err.Error(), "cannot unmarshal array")

	world.Reset()
	_ = ecs.AddResource(&world, &Velocity{})
	err = archeserde.Deserialize([]byte(textErrComponent2), &world)
	fmt.Println(err)
	assert.Contains(t, err.Error(), "cannot unmarshal array")

	world.Reset()
	_ = ecs.AddResource(&world, &Velocity{})
	err = archeserde.Deserialize([]byte(textErrResource), &world)
	fmt.Println(err)
	assert.Contains(t, err.Error(), "cannot unmarshal array")

	world.Reset()
	err = archeserde.Deserialize([]byte(textErrRelation), &world)
	assert.Contains(t, err.Error(), "cannot unmarshal object into Go value of type [2]uint32")
}

const textOk = `{
	"World" : {"Entities":[[0,4294967295],[1,0],[2,0]],"Alive":[1,2],"Next":0,"Available":0},
	"Types" : [
	  "archeserde_test.Velocity",
	  "archeserde_test.ChildOf",
	  "archeserde_test.Position"
	],
	"Components" : [
	  {
		"archeserde_test.Position" : {"X":1,"Y":2}
	  },
	  {
		"archeserde_test.Position" : {"X":3,"Y":4},
		"archeserde_test.Velocity" : {"X":5,"Y":6},
		"archeserde_test.ChildOf" : {"Entity":[1,0]}
	  }
	],
	"Resources" : {
		"archeserde_test.Velocity" : {"X":1000,"Y":0}
	}}`

const textErrEntities = `{
	"World" : {"Entities":[[0,4294967295],[1,0],[2,0]],"Alive":[1,2],"Next":0,"Available":0},
	"Types" : [
		"archeserde_test.Velocity",
		"archeserde_test.ChildOf",
		"archeserde_test.Position"
	],
	"Components" : [
		{
		"archeserde_test.Position" : {"X":3,"Y":4},
		"archeserde_test.Velocity" : {"X":5,"Y":6},
		"archeserde_test.ChildOf" : {"Entity":[1,0]}
		}
	],
	"Resources" : {
		"archeserde_test.Velocity" : {"X":1000,"Y":0}
	}}`

const textErrTypes = `{
	"World" : {"Entities":[[0,4294967295],[1,0],[2,0]],"Alive":[1,2],"Next":0,"Available":0},
	"Types" : {"a": "b"},
	"Components" : [
		{
		  "archeserde_test.Position" : {"X":1,"Y":2}
		},
		{
		"archeserde_test.Position" : {"X":3,"Y":4},
		"archeserde_test.Velocity" : {"X":5,"Y":6},
		"archeserde_test.ChildOf" : {"Entity":[1,0]}
		}
	],
	"Resources" : {
		"archeserde_test.Velocity" : {"X":1000,"Y":0}
	}}`

const textErrComponent = `{
	"World" : {"Entities":[[0,4294967295],[1,0],[2,0]],"Alive":[1,2],"Next":0,"Available":0},
	"Types" : [
		"archeserde_test.Velocity",
		"archeserde_test.ChildOf",
		"archeserde_test.Position"
	],
	"Components" : [
		[],
		{
		"archeserde_test.Position" : {"X":3,"Y":4},
		"archeserde_test.Velocity" : {"X":5,"Y":6},
		"archeserde_test.ChildOf" : {"Entity":[1,0]}
		}
	],
	"Resources" : {
		"archeserde_test.Velocity" : {"X":1000,"Y":0}
	}}`

const textErrComponent2 = `{
	"World" : {"Entities":[[0,4294967295],[1,0],[2,0]],"Alive":[1,2],"Next":0,"Available":0},
	"Types" : [
		"archeserde_test.Velocity",
		"archeserde_test.ChildOf",
		"archeserde_test.Position"
	],
	"Components" : [
		{
		  "archeserde_test.Position" : []
		},
		{
		"archeserde_test.Position" : {"X":3,"Y":4},
		"archeserde_test.Velocity" : {"X":5,"Y":6},
		"archeserde_test.ChildOf" : {"Entity":[1,0]}
		}
	],
	"Resources" : {
		"archeserde_test.Velocity" : {"X":1000,"Y":0}
	}}`

const textErrRelation = `{
	"World" : {"Entities":[[0,4294967295],[1,0],[2,0]],"Alive":[1,2],"Next":0,"Available":0},
	"Types" : [
	  "archeserde_test.Position",
	  "archeserde_test.ChildRelation"
	],
	"Components" : [
	  {
		"archeserde_test.Position" : {"X":1,"Y":2}
	  },
	  {
		"archeserde_test.Position" : {"X":5,"Y":6},
		"arche.relation.Target" : {},
		"archeserde_test.ChildRelation" : {"Dummy":0}
	  }
	],
	"Resources" : {
	}}`

const textErrResource = `{
	"World" : {"Entities":[[0,4294967295]],"Alive":[],"Next":0,"Available":0},
	"Types" : [],
	"Components" : [],
	"Resources" : {
		"archeserde_test.Velocity" : []
	}}`
