package archeserde

import "github.com/mlange-42/arche/ecs"

type deserializer struct {
	World      ecs.EntityData
	Types      []string
	Entities   []ecs.Entity
	Components []entry
	Resources  map[string]entry
}

type entry struct {
	Bytes []byte
}

func (e *entry) UnmarshalJSON(jsonData []byte) error {
	e.Bytes = jsonData
	return nil
}

func (e *entry) String() string {
	return string(e.Bytes)
}
