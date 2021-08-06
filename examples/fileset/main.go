package main

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/go-hclog"
	"github.com/treethought/roc"
)

var log = hclog.Default()

type Fileset struct {
	*roc.Accessor
}

func New() *Fileset {
	return &Fileset{
		Accessor: roc.NewAccessor("fileset"),
	}
}

// Source retrieves representation of resource
func (e *Fileset) Source(ctx *roc.RequestContext) roc.Representation {
	pathRefs, ok := ctx.Request.Arguments["path"]
	if !ok {
		return fmt.Errorf("no path provided")
	}

	result := ""

	for _, p := range pathRefs {

		path, err := ctx.Source(roc.Identifier(p), nil)
		if err != nil {
			return err
		}

		data, err := ioutil.ReadFile(fmt.Sprint(path))
		if err != nil {
			return err
		}
		result = fmt.Sprintf("%s\n%s", result, string(data))
	}
	return result

}

func main() {

	endpoint := New()
	roc.Serve(endpoint)

}
