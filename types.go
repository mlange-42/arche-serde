package archeserde

import "github.com/mlange-42/arche/ecs"

type deserializer struct {
	World      ecs.EntityDump
	Types      []string
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
