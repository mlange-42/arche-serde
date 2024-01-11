package archeserde

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/mlange-42/arche/ecs"
)

// Deserialize an Arche [ecs.World] from JSON.
//
// The world must be prepared the following way:
//   - All required component types must be registered using [ecs.ComponentID]
//   - All required resources must be added using [ecs.AddResource]
//
// After deserialization, it is not guaranteed that entity iteration order in queries is the same as before.
// More precisely, it should at first be the same as before, but will likely deviate over time from what would
// happen when continuing the original, serialized run.
func Deserialize(jsonData []byte, world *ecs.World) error {
	deserial := deserializer{}
	if err := json.Unmarshal(jsonData, &deserial); err != nil {
		return err
	}

	world.SetEntityData(&deserial.World)

	if err := deserializeComponents(world, &deserial); err != nil {
		return err
	}
	if err := deserializeResources(world, &deserial); err != nil {
		return err
	}

	return nil
}

func deserializeComponents(world *ecs.World, deserial *deserializer) error {
	infos := map[ecs.ID]ecs.CompInfo{}
	ids := map[string]ecs.ID{}
	for i := 0; i < ecs.MaskTotalBits; i++ {
		if info, ok := ecs.ComponentInfo(world, ecs.ID(i)); ok {
			infos[ecs.ID(i)] = info
			ids[info.Type.String()] = ecs.ID(i)
		}
	}

	for _, tp := range deserial.Types {
		if _, ok := ids[tp]; !ok {
			return fmt.Errorf("component type is not registered: %s", tp)
		}
	}

	for i, comps := range deserial.Components {
		entity := deserial.World.Entities[deserial.World.Alive[i]]

		mp := map[string]entry{}

		if err := json.Unmarshal(comps.Bytes, &mp); err != nil {
			return err
		}

		target := ecs.Entity{}
		var targetComp ecs.ID
		components := []ecs.Component{}
		for tpName, value := range mp {
			if tpName == targetTag {
				if err := json.Unmarshal(value.Bytes, &target); err != nil {
					return err
				}
				continue
			}

			id := ids[tpName]
			info := infos[id]

			if info.IsRelation {
				targetComp = id
			}

			component := reflect.New(info.Type).Interface()
			if err := json.Unmarshal(value.Bytes, &component); err != nil {
				return err
			}
			components = append(components, ecs.Component{
				ID:   id,
				Comp: component,
			})
		}
		builder := ecs.NewBuilderWith(world, components...)
		if target.IsZero() {
			builder.Add(entity)
		} else {
			builder = builder.WithRelation(targetComp)
			builder.Add(entity, target)
		}
	}
	return nil
}

func deserializeResources(world *ecs.World, deserial *deserializer) error {
	resTypes := map[ecs.ResID]reflect.Type{}
	resIds := map[string]ecs.ResID{}
	for i := 0; i < ecs.MaskTotalBits; i++ {
		if tp, ok := ecs.ResourceType(world, ecs.ResID(i)); ok {
			resTypes[ecs.ResID(i)] = tp
			resIds[tp.String()] = ecs.ResID(i)
		}
	}

	for tpName, res := range deserial.Resources {
		resID, ok := resIds[tpName]
		if !ok {
			return fmt.Errorf("resource type is not registered: %s", tpName)
		}

		tp := resTypes[resID]
		resource := reflect.New(tp).Interface()
		if err := json.Unmarshal(res.Bytes, &resource); err != nil {
			return err
		}

		resLoc := world.Resources().Get(resID)
		if resLoc == nil {
			return fmt.Errorf("resource type registered but nil: %s", tpName)
		}

		rValue := reflect.ValueOf(resLoc)
		ptr := rValue.UnsafePointer()
		value := reflect.NewAt(tp, ptr).Interface()

		if err := json.Unmarshal(res.Bytes, &value); err != nil {
			return err
		}
	}
	return nil
}
