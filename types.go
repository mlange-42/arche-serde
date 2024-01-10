package archeserde

type deserializer struct {
	Components []string
	Entities   []entry
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

type entity struct {
	ID  uint32 `json:".ID"`
	Gen uint32 `json:".Gen"`
}
