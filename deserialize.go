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
func Deserialize(jsonData []byte, world *ecs.World) error {
	types := map[ecs.ID]reflect.Type{}
	ids := map[string]ecs.ID{}
	for i := 0; i < ecs.MaskTotalBits; i++ {
		if tp, ok := ecs.ComponentType(world, ecs.ID(i)); ok {
			types[ecs.ID(i)] = tp
			ids[tp.String()] = ecs.ID(i)
		}
	}

	resTypes := map[ecs.ResID]reflect.Type{}
	resIds := map[string]ecs.ResID{}
	for i := 0; i < ecs.MaskTotalBits; i++ {
		if tp, ok := ecs.ResourceType(world, ecs.ResID(i)); ok {
			resTypes[ecs.ResID(i)] = tp
			resIds[tp.String()] = ecs.ResID(i)
		}
	}

	deserial := deserializer{}

	if err := json.Unmarshal(jsonData, &deserial); err != nil {
		return err
	}

	for _, tp := range deserial.Components {
		if _, ok := ids[tp]; !ok {
			return fmt.Errorf("component type is not registered: %s", tp)
		}
	}

	for _, e := range deserial.Entities {
		mp := map[string]entry{}

		if err := json.Unmarshal(e.Bytes, &mp); err != nil {
			return err
		}

		components := []ecs.Component{}
		for tpName, value := range mp {
			id := ids[tpName]
			tp := types[id]

			component := reflect.New(tp).Interface()
			if err := json.Unmarshal(value.Bytes, &component); err != nil {
				return err
			}
			components = append(components, ecs.Component{
				ID:   id,
				Comp: component,
			})
		}
		_ = world.NewEntityWith(components...)
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
