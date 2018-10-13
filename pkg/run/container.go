package run

import (
	"fmt"
)

// ContainerSpec defines everything that's needed to run a container.
type ContainerSpec struct {
	Imgref string
	Cmd    []string
	Env    map[string]string
	Mount  Mount
}

// Mount describes a volume mount in a container.
type Mount struct {
	Src   string
	Point string
	Type  string

	// If true, clean up the mounted volume after it's done being used. This
	// means whatever the Docker Engine API means by "removing volumes" when
	// removing a container.
	Cleanup bool
}

// GetEnvArray is a utility function for getting the environment
// variables as an array to be passed to a container config.
func (spec ContainerSpec) GetEnvArray() []string {
	ret := []string{}

	for k, v := range spec.Env {
		ret = append(ret, fmt.Sprintf("%v=%v", k, v))
	}

	return ret
}
