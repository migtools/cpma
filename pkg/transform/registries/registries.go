package registries

// Registries holds a list of Registries
type Registries struct {
	List []string `toml:"registries"`
}

// Validate the registries data
func Validate(e map[string]Registries) int {
	if len(e["block"].List) == 0 && len(e["insecure"].List) == 0 && len(e["search"].List) == 0 {
		return 1
	}
	return 0
}
