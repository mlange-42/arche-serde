package archeserde_test

import (
	"fmt"

	archeserde "github.com/mlange-42/arche-serde"
	"github.com/mlange-42/arche/ecs"
)

func Example() {
	// Create a world.
	world := ecs.NewWorld()

	// Populate the world with entities, components and resources.
	// ...

	// Serialize the world.
	jsonData, err := archeserde.Serialize(&world)
	if err != nil {
		fmt.Printf("could not serialize: %s\n", err)
		return
	}

	// Print the resulting JSON.
	// fmt.Println(string(jsonData))

	// Create a new, empty world.
	newWorld := ecs.NewWorld()

	// Deserialize into the new world.
	err = archeserde.Deserialize(jsonData, &newWorld)
	if err != nil {
		fmt.Printf("could not deserialize: %s\n", err)
		return
	}
	// Output:
}
