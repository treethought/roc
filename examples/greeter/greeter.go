package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/treethought/roc"
	"github.com/treethought/roc/proto"
)

var log = hclog.New(&hclog.LoggerOptions{
	Color:  hclog.AutoColor,
	Output: os.Stdout,
})

type MyEndpoint struct {
	*roc.Accessor
}

func New() *MyEndpoint {
	return &MyEndpoint{
		Accessor: roc.NewAccessor("greeter"),
	}
}

// Source retrieves representation of resource
func (e *MyEndpoint) Source(ctx *roc.RequestContext) interface{} {
	// get the argument value from the identifier provided as request argument
	name, err := ctx.GetArgumentValue("name")
	if err != nil {
		log.Warn("failed to get name argument", "err", err)
		return "no name?"
	}

	// issue subrequest to upper case the name
	upperID := roc.NewIdentifier("res://toUpper")
	req := ctx.CreateRequest(upperID)
	req.SetArgumentByValue("value", name)
	upped, err := ctx.IssueRequest(req)
	if err != nil {
		log.Error("failed to dispatch subrequest request", "error", err)
	}

	s := new(proto.String)
	err = upped.To(s)
	if err != nil {
		return err
	}

	return fmt.Sprintf("hello world: %s\n", s.Value)

}

func main() {
	endpoint := New()
	roc.Serve(endpoint)

}
