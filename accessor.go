package roc

import (
	"fmt"
	"os"
)

const EndpointTypeAccessor string = "accessor"

// Accessor is a struct implementing the default behavior for an empty EndpointAccessor
// This type is useful for embedding with custom implementations of EndpointAccessor
type Accessor struct {
	BaseEndpoint
	// grammar Grammar `yaml:"grammar,omitempty"`
	Name string
}

func NewAccessor(name string) *Accessor {
	return &Accessor{
		Name:         name,
		BaseEndpoint: BaseEndpoint{},
	}
}

func (a *Accessor) Identifier() Identifier {
	path, err := os.Executable()
	if err != nil {
		log.Error("unable to locate identifier", "error", err)
		return ""
	}

	return Identifier(fmt.Sprintf("accessor://%s", path))
}

func (e Accessor) Type() string {
	return EndpointTypeAccessor
}

func (e Accessor) String() string {
	return fmt.Sprintf("endpoint://%s", e.Name)
}
