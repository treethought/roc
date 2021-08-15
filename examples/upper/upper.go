package main

import (
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/treethought/roc"
	proto "github.com/treethought/roc/proto/v1"
)

var log = hclog.Default()

type MyEndpoint struct {
	*roc.Accessor
}

func New() *MyEndpoint {
	return &MyEndpoint{
		Accessor: roc.NewAccessor("upper"),
	}
}

// Source retrieves representation of resource
func (e *MyEndpoint) Source(ctx *roc.RequestContext) interface{} {
	value, err := ctx.GetArgumentValue("value")
	if err != nil {
		log.Error("failed source value", "err", err)
		return err
	}
	s := new(proto.String)
	err = value.To(s)
	if err != nil {
		return err
	}

	// TODO better identifier parsing, and handle this with GetArgumentValue
	if strings.Contains(s.GetValue(), ":/") {
		valRep, err := ctx.Source(roc.NewIdentifier(s.GetValue()), "")
		if err != nil {
			log.Error("failed to source upper identifier")
			return err
		}
		s = new(proto.String)
		err = valRep.To(s)
		if err != nil {
			return err
		}
	}

	return strings.ToUpper(s.Value)

}

func main() {
	endpoint := New()
	roc.Serve(endpoint)

}
