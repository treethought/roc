package roc

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
)

const EndpointTypeAccessor string = "accessor"

// Accessor is a struct implementing the default behavior for an empty EndpointAccessor
// This type is useful for embedding with custom implementations of EndpointAccessor
type Accessor struct {
	// grammar Grammar `yaml:"grammar,omitempty"`
	Name   string
	Logger hclog.Logger
}

func NewAccessor(name string) *Accessor {
	return &Accessor{
		Name: name,
		Logger: hclog.New(&hclog.LoggerOptions{
			Level:       hclog.Debug,
			Output:      os.Stderr,
			JSONFormat:  false,
			Name:        name,
			Color:       hclog.ForceColor,
			DisableTime: true,
		}),
	}
}

func (a *Accessor) Identifier() Identifier {
	path, err := os.Executable()
	if err != nil {
		a.Logger.Error("unable to locate identifier", "error", err)
		return ""
	}

	return Identifier(fmt.Sprintf("accessor://%s", path))
}

func (a *Accessor) SetLogger(l hclog.Logger) {
	a.Logger = l
}

func (e Accessor) Type() string {
	return EndpointTypeAccessor
}

func (e Accessor) String() string {
	return fmt.Sprintf("endpoint://%s", e.Name)
}

func (e Accessor) Source(ctx *RequestContext) Representation {
	return nil
}

func (e Accessor) Sink(ctx *RequestContext) {}

func (e Accessor) New(ctx *RequestContext) Identifier {
	return ""
}
func (e Accessor) Delete(ctx *RequestContext) bool {
	return false
}
func (e Accessor) Exists(ctx *RequestContext) bool {
	return false
}
func (e Accessor) Transrept(ctx *RequestContext) Representation {
	return nil
}

func (e Accessor) Evaluate(ctx *RequestContext) Representation {

	switch ctx.Request.Verb {
	case Source:
		return e.Source(ctx)
	case Sink:
		e.Sink(ctx)
		return nil
	case New:
		return e.New(ctx)
	case Delete:
		return e.Delete(ctx)
	case Exists:
		return e.Exists(ctx)

	default:
		return e.Source(ctx)

	}
}
