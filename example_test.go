package archeserde_test

import (
	"fmt"
	"math/rand"

	archeserde "github.com/mlange-42/arche-serde"
	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

const (
	width  = 40
	height = 12
)

type Coords struct {
	X int
	Y int
}

func Example() {
	rng := rand.New(rand.NewSource(42))

	// Create a world.
	world := ecs.NewWorld()

	// Populate the world with entities, components and resources.
	builder := generic.NewMap1[Coords](&world)
	query := builder.NewBatchQ(60)
	for query.Next() {
		coord := query.Get()
		coord.X = rng.Intn(width)
		coord.Y = rng.Intn(height)
	}

	// Print the original world
	fmt.Println("====== Original world ========")
	printWorld(&world)

	// Serialize the world.
	jsonData, err := archeserde.Serialize(&world)
	if err != nil {
		fmt.Printf("could not serialize: %s\n", err)
		return
	}

	// Print the resulting JSON.
	//fmt.Println(string(jsonData))

	// Create a new, empty world.
	newWorld := ecs.NewWorld()

	// Register required components and resources
	_ = ecs.ComponentID[Coords](&newWorld)

	// Deserialize into the new world.
	err = archeserde.Deserialize(jsonData, &newWorld)
	if err != nil {
		fmt.Printf("could not deserialize: %s\n", err)
		return
	}

	// Print the deserialized world
	fmt.Println("====== Deserialized world ========")
	printWorld(&newWorld)
	// Output: ====== Original world ========
	// --------------------------------O-O---O-
	// -----------------------O----------------
	// -O-------------O------OO--------------O-
	// ----O------------------------OOO--------
	// O--------------OO-O---------------------
	// ------------O-----------O---------------
	// --O-------------O-------O---O------O----
	// O-O-----O----OOO-O--O--------------OO---
	// -----------O---OO----O--O------------O--
	// ------------O-----O---------------------
	// ---O---------------O------O--O----------
	// ------O-OO--O---------OO-OOO-----------O
	// ====== Deserialized world ========
	// --------------------------------O-O---O-
	// -----------------------O----------------
	// -O-------------O------OO--------------O-
	// ----O------------------------OOO--------
	// O--------------OO-O---------------------
	// ------------O-----------O---------------
	// --O-------------O-------O---O------O----
	// O-O-----O----OOO-O--O--------------OO---
	// -----------O---OO----O--O------------O--
	// ------------O-----O---------------------
	// ---O---------------O------O--O----------
	// ------O-OO--O---------OO-OOO-----------O
}

func printWorld(world *ecs.World) {
	grid := make([][]rune, height)
	for i := range grid {
		grid[i] = make([]rune, width)
		for j := range grid[i] {
			grid[i][j] = '-'
		}
	}

	filter := generic.NewFilter1[Coords]()
	query := filter.Query(world)

	for query.Next() {
		coords := query.Get()
		grid[coords.Y][coords.X] = 'O'
	}

	for i := 0; i < len(grid); i++ {
		fmt.Println(string(grid[i]))
	}
}
