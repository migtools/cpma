package registries

import (
	"errors"
)

// Registries holds a list of Registries
type Registries struct {
	List []string `toml:"registries"`
}

// Validate the registries data
func Validate(e map[string]Registries) error {
	if len(e["block"].List) == 0 && len(e["insecure"].List) == 0 && len(e["search"].List) == 0 {
		return errors.New("no configured registries detected, not generating a cr or report")
	}
	return nil
}
