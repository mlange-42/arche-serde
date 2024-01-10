package archeserde_test

import (
	"fmt"
	"testing"

	archeserde "github.com/mlange-42/arche-serde"
	"github.com/mlange-42/arche/ecs"
	"github.com/stretchr/testify/assert"
)

type Position struct {
	X float64
	Y float64
}

type Velocity struct {
	X float64
	Y float64
}

type Parent struct {
	Entity ecs.Entity
}

func TestSerialize(t *testing.T) {
	w := ecs.NewWorld()

	posId := ecs.ComponentID[Position](&w)
	velId := ecs.ComponentID[Velocity](&w)
	parId := ecs.ComponentID[Parent](&w)

	parent := w.NewEntityWith(ecs.Component{ID: posId, Comp: &Position{X: 1, Y: 2}})
	child := w.NewEntityWith(
		ecs.Component{ID: posId, Comp: &Position{X: 3, Y: 4}},
		ecs.Component{ID: velId, Comp: &Velocity{X: 5, Y: 6}},
		ecs.Component{ID: parId, Comp: &Parent{Entity: parent}},
	)

	resId := ecs.ResourceID[Velocity](&w)
	w.Resources().Add(resId, &Velocity{X: 1000, Y: 0})

	jsonData, err := archeserde.Serialize(&w)

	if err != nil {
		assert.Fail(t, "could not serialize: %s\n", err)
	}

	fmt.Println(string(jsonData))

	w = ecs.NewWorld()
	posId = ecs.ComponentID[Position](&w)
	velId = ecs.ComponentID[Velocity](&w)
	parId = ecs.ComponentID[Parent](&w)
	_ = ecs.AddResource[Velocity](&w, &Velocity{})

	err = archeserde.Deserialize(jsonData, &w)
	if err != nil {
		assert.Fail(t, "could not deserialize: %s\n", err)
	}

	query := w.Query(ecs.All())

	assert.Equal(t, query.Count(), 2)

	query.Next()
	assert.True(t, query.Has(posId))
	assert.False(t, query.Has(velId))
	assert.Equal(t, *(*Position)(query.Get(posId)), Position{X: 1, Y: 2})

	query.Next()
	assert.True(t, query.Has(posId))
	assert.True(t, query.Has(velId))
	assert.Equal(t, *(*Position)(query.Get(posId)), Position{X: 3, Y: 4})
	assert.Equal(t, *(*Velocity)(query.Get(velId)), Velocity{X: 5, Y: 6})
	assert.Equal(t, *(*Parent)(query.Get(parId)), Parent{Entity: parent})

	res := (*Velocity)(ecs.GetResource[Velocity](&w))

	assert.Equal(t, *res, Velocity{X: 1000})

	assert.True(t, w.Alive(parent))
	assert.True(t, w.Alive(child))
}
