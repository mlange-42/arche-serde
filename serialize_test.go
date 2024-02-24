package archeserde_test

import (
	"fmt"
	"testing"

	archeserde "github.com/mlange-42/arche-serde"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
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

type ChildOf struct {
	Entity ecs.Entity
}

type ChildRelation struct {
	ecs.Relation
	Dummy int
}

type Generic[T any] struct {
	Value T
}

func serialize(opts ...archeserde.Option) ([]byte, ecs.Entity, ecs.Entity, error) {
	w := ecs.NewWorld()

	posId := ecs.ComponentID[Position](&w)
	velId := ecs.ComponentID[Velocity](&w)
	childId := ecs.ComponentID[ChildOf](&w)

	parent := w.NewEntityWith(ecs.Component{ID: posId, Comp: &Position{X: 1, Y: 2}})
	child := w.NewEntityWith(
		ecs.Component{ID: posId, Comp: &Position{X: 3, Y: 4}},
		ecs.Component{ID: velId, Comp: &Velocity{X: 5, Y: 6}},
		ecs.Component{ID: childId, Comp: &ChildOf{Entity: parent}},
	)
	w.NewEntity()

	resId := ecs.ResourceID[Velocity](&w)
	resId2 := ecs.ResourceID[Position](&w)
	w.Resources().Add(resId, &Velocity{X: 1000, Y: 0})
	w.Resources().Add(resId2, &Position{X: 1000, Y: 0})

	js, err := archeserde.Serialize(&w, opts...)
	return js, parent, child, err
}

func TestSerialize(t *testing.T) {
	jsonData, parent, child, err := serialize()

	if err != nil {
		assert.Fail(t, "could not serialize: %s\n", err)
	}

	fmt.Println(string(jsonData))

	w := ecs.NewWorld()
	posId := ecs.ComponentID[Position](&w)
	velId := ecs.ComponentID[Velocity](&w)
	childId := ecs.ComponentID[ChildOf](&w)
	_ = ecs.AddResource[Position](&w, &Position{})
	_ = ecs.AddResource[Velocity](&w, &Velocity{})

	err = archeserde.Deserialize(jsonData, &w)
	if err != nil {
		assert.Fail(t, "could not deserialize: %s\n", err)
	}

	query := w.Query(ecs.All())

	assert.Equal(t, query.Count(), 3)

	query.Next()
	assert.False(t, query.Has(posId))
	assert.False(t, query.Has(velId))

	query.Next()
	assert.True(t, query.Has(posId))
	assert.False(t, query.Has(velId))
	assert.Equal(t, *(*Position)(query.Get(posId)), Position{X: 1, Y: 2})

	query.Next()
	assert.True(t, query.Has(posId))
	assert.True(t, query.Has(velId))
	assert.Equal(t, *(*Position)(query.Get(posId)), Position{X: 3, Y: 4})
	assert.Equal(t, *(*Velocity)(query.Get(velId)), Velocity{X: 5, Y: 6})
	assert.Equal(t, *(*ChildOf)(query.Get(childId)), ChildOf{Entity: parent})

	res := (*Velocity)(ecs.GetResource[Velocity](&w))

	assert.Equal(t, *res, Velocity{X: 1000})

	assert.True(t, w.Alive(parent))
	assert.True(t, w.Alive(child))
}

func TestSerializeSkipEntities(t *testing.T) {
	jsonData, _, _, err := serialize(archeserde.Opts.SkipEntities())

	if err != nil {
		assert.Fail(t, "could not serialize: %s\n", err)
	}

	fmt.Println(string(jsonData))

	w := ecs.NewWorld()
	_ = ecs.AddResource[Position](&w, &Position{})
	_ = ecs.AddResource[Velocity](&w, &Velocity{})

	err = archeserde.Deserialize(jsonData, &w)
	if err != nil {
		assert.Fail(t, "could not deserialize: %s\n", err)
	}

	query := w.Query(ecs.All())

	assert.Equal(t, query.Count(), 0)
	query.Close()

	res := (*Velocity)(ecs.GetResource[Velocity](&w))
	assert.Equal(t, *res, Velocity{X: 1000})
}

func TestSerializeSkipAllComponents(t *testing.T) {
	jsonData, parent, child, err := serialize(archeserde.Opts.SkipAllComponents())

	if err != nil {
		assert.Fail(t, "could not serialize: %s\n", err)
	}

	fmt.Println(string(jsonData))

	w := ecs.NewWorld()
	_ = ecs.AddResource[Position](&w, &Position{})
	_ = ecs.AddResource[Velocity](&w, &Velocity{})

	err = archeserde.Deserialize(jsonData, &w)
	if err != nil {
		assert.Fail(t, "could not deserialize: %s\n", err)
	}

	query := w.Query(ecs.All())

	assert.Equal(t, query.Count(), 3)
	query.Close()

	res := (*Velocity)(ecs.GetResource[Velocity](&w))

	assert.Equal(t, *res, Velocity{X: 1000})

	assert.True(t, w.Alive(parent))
	assert.True(t, w.Alive(child))
}

func TestSerializeSkipComponents(t *testing.T) {
	jsonData, parent, child, err := serialize(archeserde.Opts.SkipComponents(generic.T[Position]()))

	if err != nil {
		assert.Fail(t, "could not serialize: %s\n", err)
	}

	fmt.Println(string(jsonData))

	w := ecs.NewWorld()
	velId := ecs.ComponentID[Velocity](&w)
	childId := ecs.ComponentID[ChildOf](&w)
	_ = ecs.AddResource[Position](&w, &Position{})
	_ = ecs.AddResource[Velocity](&w, &Velocity{})

	err = archeserde.Deserialize(jsonData, &w)
	if err != nil {
		assert.Fail(t, "could not deserialize: %s\n", err)
	}

	query := w.Query(ecs.All())

	assert.Equal(t, query.Count(), 3)

	query.Next()
	assert.False(t, query.Has(velId))

	query.Next()
	assert.False(t, query.Has(velId))

	query.Next()
	assert.True(t, query.Has(velId))
	assert.Equal(t, *(*Velocity)(query.Get(velId)), Velocity{X: 5, Y: 6})
	assert.Equal(t, *(*ChildOf)(query.Get(childId)), ChildOf{Entity: parent})

	res := (*Velocity)(ecs.GetResource[Velocity](&w))

	assert.Equal(t, *res, Velocity{X: 1000})

	assert.True(t, w.Alive(parent))
	assert.True(t, w.Alive(child))
}

func TestSerializeSkipAllResources(t *testing.T) {
	jsonData, _, _, err := serialize(archeserde.Opts.SkipAllResources())

	if err != nil {
		assert.Fail(t, "could not serialize: %s\n", err)
	}

	fmt.Println(string(jsonData))

	w := ecs.NewWorld()
	_ = ecs.ComponentID[Position](&w)
	_ = ecs.ComponentID[Velocity](&w)
	_ = ecs.ComponentID[ChildOf](&w)

	err = archeserde.Deserialize(jsonData, &w)
	if err != nil {
		assert.Fail(t, "could not deserialize: %s\n", err)
	}
}

func TestSerializeSkipResources(t *testing.T) {
	jsonData, _, _, err := serialize(archeserde.Opts.SkipResources(generic.T[Position]()))

	if err != nil {
		assert.Fail(t, "could not serialize: %s\n", err)
	}

	fmt.Println(string(jsonData))

	w := ecs.NewWorld()
	_ = ecs.ComponentID[Position](&w)
	_ = ecs.ComponentID[Velocity](&w)
	_ = ecs.ComponentID[ChildOf](&w)
	_ = ecs.AddResource[Velocity](&w, &Velocity{})

	err = archeserde.Deserialize(jsonData, &w)
	if err != nil {
		assert.Fail(t, "could not deserialize: %s\n", err)
	}

	res := (*Velocity)(ecs.GetResource[Velocity](&w))

	assert.Equal(t, *res, Velocity{X: 1000})
}

func TestSerializeRelation(t *testing.T) {
	w := ecs.NewWorld()

	posId := ecs.ComponentID[Position](&w)
	relId := ecs.ComponentID[ChildRelation](&w)

	parent := w.NewEntityWith(ecs.Component{ID: posId, Comp: &Position{X: 1, Y: 2}})
	child1 := w.NewEntityWith(
		ecs.Component{ID: posId, Comp: &Position{X: 3, Y: 4}},
		ecs.Component{ID: relId, Comp: &ChildRelation{}},
	)
	child2 := w.NewEntityWith(
		ecs.Component{ID: posId, Comp: &Position{X: 5, Y: 6}},
		ecs.Component{ID: relId, Comp: &ChildRelation{}},
	)

	w.Relations().Set(child2, relId, parent)

	jsonData, err := archeserde.Serialize(&w)
	if err != nil {
		assert.Fail(t, "could not serialize: %s\n", err)
	}
	fmt.Println(string(jsonData))

	w = ecs.NewWorld()
	_ = ecs.ComponentID[Position](&w)
	relId = ecs.ComponentID[ChildRelation](&w)

	err = archeserde.Deserialize(jsonData, &w)
	if err != nil {
		assert.Fail(t, "could not deserialize: %s\n", err)
	}

	assert.Equal(t, w.Relations().Get(child1, relId), ecs.Entity{})
	assert.Equal(t, w.Relations().Get(child2, relId), parent)
}

func TestSerializeGeneric(t *testing.T) {
	w := ecs.NewWorld()

	gen1Id := ecs.ComponentID[Generic[int32]](&w)
	gen2Id := ecs.ComponentID[Generic[float32]](&w)

	e1 := w.NewEntityWith(
		ecs.Component{ID: gen1Id, Comp: &Generic[int32]{Value: 1}},
	)
	e2 := w.NewEntityWith(
		ecs.Component{ID: gen2Id, Comp: &Generic[float32]{Value: 2.0}},
	)
	e3 := w.NewEntityWith(
		ecs.Component{ID: gen1Id, Comp: &Generic[int32]{Value: 3}},
		ecs.Component{ID: gen2Id, Comp: &Generic[float32]{Value: 4.0}},
	)

	jsonData, err := archeserde.Serialize(&w)
	if err != nil {
		assert.Fail(t, "could not serialize: %s\n", err)
	}
	fmt.Println(string(jsonData))

	w = ecs.NewWorld()
	_ = ecs.ComponentID[Generic[int32]](&w)
	_ = ecs.ComponentID[Generic[float32]](&w)

	err = archeserde.Deserialize(jsonData, &w)
	if err != nil {
		assert.Fail(t, "could not deserialize: %s\n", err)
	}

	mapper := generic.NewMap2[Generic[int32], Generic[float32]](&w)

	c1, c2 := mapper.Get(e1)
	assert.Equal(t, c1, &Generic[int32]{Value: 1})
	assert.Equal(t, c2, (*Generic[float32])(nil))

	c1, c2 = mapper.Get(e2)
	assert.Equal(t, c1, (*Generic[int32])(nil))
	assert.Equal(t, c2, &Generic[float32]{Value: 2.0})

	c1, c2 = mapper.Get(e3)
	assert.Equal(t, c1, &Generic[int32]{Value: 3})
	assert.Equal(t, c2, &Generic[float32]{Value: 4.0})

	_, _, _ = e1, e2, e3
}
