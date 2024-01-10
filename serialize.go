package archeserde

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/mlange-42/arche/ecs"
)

// Serialize an Arche ECS world to JSON.
func Serialize(world *ecs.World) ([]byte, error) {
	builder := strings.Builder{}

	builder.WriteString("{\"Components\" : [\n")

	types := map[ecs.ID]reflect.Type{}
	for i := 0; i < ecs.MaskTotalBits; i++ {
		if tp, ok := ecs.ComponentType(world, ecs.ID(i)); ok {
			types[ecs.ID(i)] = tp
		}
	}
	maxComp := len(types) - 1
	counter := 0
	for _, tp := range types {
		builder.WriteString(fmt.Sprintf("  \"%s\"", tp.String()))
		if counter < maxComp {
			builder.WriteString(",")
		}
		builder.WriteString("\n")
		counter++
	}

	builder.WriteString("],\n\"Entities\" : [\n")
	query := world.Query(ecs.All())
	lastEntity := query.Count() - 1
	counter = 0
	for query.Next() {
		builder.WriteString("  {\n")

		ids := query.Ids()
		last := len(ids) - 1
		for i, id := range ids {
			tp, _ := ecs.ComponentType(world, id)

			comp := query.Get(id)
			value := reflect.NewAt(tp, comp).Interface()
			jsonData, err := json.Marshal(value)
			if err != nil {
				return nil, err
			}
			builder.WriteString("    ")
			builder.WriteString(fmt.Sprintf("\"%s\" : ", tp.String()))
			builder.WriteString(string(jsonData))
			if i < last {
				builder.WriteString(",")
			}
			builder.WriteString("\n")
		}

		builder.WriteString("  }")
		if counter < lastEntity {
			builder.WriteString(",")
		}
		builder.WriteString("\n")

		counter++
	}
	builder.WriteString("],\n\"Resources\" : {\n")

	resTypes := map[ecs.ResID]reflect.Type{}
	for i := 0; i < ecs.MaskTotalBits; i++ {
		if tp, ok := ecs.ResourceType(world, ecs.ResID(i)); ok {
			resTypes[ecs.ResID(i)] = tp
		}
	}

	last := len(resTypes) - 1
	counter = 0
	for id, tp := range resTypes {
		res := world.Resources().Get(id)
		rValue := reflect.ValueOf(res)
		ptr := rValue.UnsafePointer()

		value := reflect.NewAt(tp, ptr).Interface()
		jsonData, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}

		builder.WriteString("    ")
		builder.WriteString(fmt.Sprintf("\"%s\" : ", tp.String()))
		builder.WriteString(string(jsonData))

		if counter < last {
			builder.WriteString(",")
		}
		builder.WriteString("\n")
	}

	builder.WriteString("}}")

	return []byte(builder.String()), nil
}
