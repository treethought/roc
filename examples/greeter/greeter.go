package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/treethought/roc"
	proto "github.com/treethought/roc/proto/v1"
)

var log = hclog.New(&hclog.LoggerOptions{
	Color:  hclog.AutoColor,
	Output: os.Stdout,
})

type MyEndpoint struct {
	*roc.BaseEndpoint
}

func New() *MyEndpoint {
	return &MyEndpoint{
		BaseEndpoint: &roc.BaseEndpoint{},
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

	nameString := new(proto.String)
	err = name.To(nameString)
	if err != nil {
		return err
	}

	// issue subrequest to upper case the name

	// via resource identifer that will match endpoint grammar
	// upperID := roc.NewIdentifier(fmt.Sprintf("active:toUpper+value@%s", name))
	// req := ctx.CreateRequest(upperID)

	// or construct the identifier behind the scenes using the endpoint's ID
	metaID := roc.NewIdentifier("example:upper")
	args := map[string][]string{"value": {nameString.Value}}

	req, err := ctx.CreateRequestToEndpoint(metaID.String(), args)
	if err != nil {
		log.Error("failed to create targetted endpoint request", "endpoint", metaID, "error", err)
		return err
	}

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
