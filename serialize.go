package archeserde

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/mlange-42/arche/ecs"
)

const targetTag = "arche.relation.Target"

// Serialize an Arche ECS world to JSON.
func Serialize(world *ecs.World) ([]byte, error) {
	builder := strings.Builder{}

	builder.WriteString("{\n")

	if err := serializeWorld(world, &builder); err != nil {
		return nil, err
	}
	builder.WriteString(",\n")

	serializeTypes(world, &builder)
	builder.WriteString(",\n")

	if err := serializeEntities(world, &builder); err != nil {
		return nil, err
	}
	builder.WriteString(",\n")

	if err := serializeComponents(world, &builder); err != nil {
		return nil, err
	}
	builder.WriteString(",\n")

	if err := serializeResources(world, &builder); err != nil {
		return nil, err
	}
	builder.WriteString("}\n")

	return []byte(builder.String()), nil
}

func serializeWorld(world *ecs.World, builder *strings.Builder) error {
	jsonData, err := world.MarshalEntities()
	if err != nil {
		return err
	}
	builder.WriteString(fmt.Sprintf("\"World\" : %s", string(jsonData)))
	return nil
}

func serializeTypes(world *ecs.World, builder *strings.Builder) {
	builder.WriteString("\"Types\" : [\n")

	types := map[ecs.ID]reflect.Type{}
	for i := 0; i < ecs.MaskTotalBits; i++ {
		if info, ok := ecs.ComponentInfo(world, ecs.ID(i)); ok {
			types[ecs.ID(i)] = info.Type
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

	builder.WriteString("]")
}

func serializeEntities(world *ecs.World, builder *strings.Builder) error {

	builder.WriteString("\"Entities\" : [\n")

	query := world.Query(ecs.All())
	lastEntity := query.Count() - 1
	counter := 0
	for query.Next() {
		jsonData, err := json.Marshal(query.Entity())
		if err != nil {
			return err
		}
		builder.WriteString(fmt.Sprintf("    %s", jsonData))
		if counter < lastEntity {
			builder.WriteString(",")
		}
		builder.WriteString("\n")

		counter++
	}
	builder.WriteString("]")

	return nil
}

func serializeComponents(world *ecs.World, builder *strings.Builder) error {

	builder.WriteString("\"Components\" : [\n")

	query := world.Query(ecs.All())
	lastEntity := query.Count() - 1
	counter := 0
	for query.Next() {
		builder.WriteString("  {\n")

		ids := query.Ids()
		last := len(ids) - 1

		for i, id := range ids {
			info, _ := ecs.ComponentInfo(world, id)

			if info.IsRelation {
				target := query.Relation(id)
				builder.WriteString(fmt.Sprintf("    \"%s\" : {\"ID\": %d, \"Gen\": %d},\n", targetTag, target.ID(), target.Gen()))
			}

			comp := query.Get(id)
			value := reflect.NewAt(info.Type, comp).Interface()
			jsonData, err := json.Marshal(value)
			if err != nil {
				return err
			}
			builder.WriteString(fmt.Sprintf("    \"%s\" : ", info.Type.String()))
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
	builder.WriteString("]")

	return nil
}

func serializeResources(world *ecs.World, builder *strings.Builder) error {
	builder.WriteString("\"Resources\" : {\n")

	resTypes := map[ecs.ResID]reflect.Type{}
	for i := 0; i < ecs.MaskTotalBits; i++ {
		if tp, ok := ecs.ResourceType(world, ecs.ResID(i)); ok {
			resTypes[ecs.ResID(i)] = tp
		}
	}

	last := len(resTypes) - 1
	counter := 0
	for id, tp := range resTypes {
		res := world.Resources().Get(id)
		rValue := reflect.ValueOf(res)
		ptr := rValue.UnsafePointer()

		value := reflect.NewAt(tp, ptr).Interface()
		jsonData, err := json.Marshal(value)
		if err != nil {
			return err
		}

		builder.WriteString("    ")
		builder.WriteString(fmt.Sprintf("\"%s\" : ", tp.String()))
		builder.WriteString(string(jsonData))

		if counter < last {
			builder.WriteString(",")
		}
		builder.WriteString("\n")
	}

	builder.WriteString("}")

	return nil
}
