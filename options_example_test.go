package archeserde_test

import (
	"fmt"

	archeserde "github.com/mlange-42/arche-serde"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

func Example_options() {
	world := ecs.NewWorld()
	builder := generic.NewMap2[Position, Velocity](&world)
	builder.NewBatch(10)

	// Serialize the world, skipping Velocity.
	jsonData, err := archeserde.Serialize(
		&world,
		archeserde.Opts.SkipComponents(generic.T[Velocity]()),
	)
	if err != nil {
		fmt.Printf("could not serialize: %s\n", err)
		return
	}

	newWorld := ecs.NewWorld()

	// Register required components and resources
	_ = ecs.ComponentID[Position](&newWorld)

	err = archeserde.Deserialize(jsonData, &newWorld)
	if err != nil {
		fmt.Printf("could not deserialize: %s\n", err)
		return
	}
}
