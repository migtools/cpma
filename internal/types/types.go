package types

import (
	"github.com/fusor/cpma/ocp3"
	"github.com/fusor/cpma/ocp4"
)

type File struct {
	Hostname string
	Path     string
	Content  []byte
	OCP3     ocp3.Cluster
	OCP4     ocp4.Cluster
}
