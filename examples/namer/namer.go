package main

import (
	"github.com/hashicorp/go-hclog"
	"github.com/treethought/roc"
	"github.com/treethought/roc/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var log = hclog.Default()

type MyEndpoint struct {
	*roc.Accessor
}

func New() *MyEndpoint {
	return &MyEndpoint{
		Accessor: roc.NewAccessor("namer"),
	}
}

func getAsString(rep roc.Representation) (*proto.String, error) {
	// TODO: this needs to be handled with response vlass behnd the scenes
	any, err := anypb.New(rep)
	if err != nil {
		return nil, err
	}

	m := new(proto.String)
	err = any.UnmarshalTo(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Source retrieves representation of resource
func (e *MyEndpoint) Source(ctx *roc.RequestContext) interface{} {
	log.Debug("Sourcing", "request", ctx.Request)
	ctx.Request().SetRepresentationClass("")

	name, err := ctx.GetArgumentValue("nameArg")
	if err != nil {
		return roc.NewRepresentation(&proto.ErrorMessage{Message: err.Error()})
	}

	m, err := getAsString(name)
	if err != nil {
		return roc.NewRepresentation(&proto.ErrorMessage{Message: err.Error()})
	}

	return m

}

func main() {
	log.Error("STARTING NAMER")

	endpoint := New()
	roc.Serve(endpoint)

}
