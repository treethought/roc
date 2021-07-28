package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/treethought/roc"
)

type FakeEndpoint struct {
	*roc.Accessor
}

func (e FakeEndpoint) Source(ctx roc.RequestContext) roc.Representation {
	return "YO"
}

func main() {

	url, err := url.Parse("res://my-resource")
	if err != nil {
		panic(err)
	}

	grammar := roc.Grammar{
		Base: url,
	}

	endpoint := &FakeEndpoint{
		Accessor: roc.NewAccessor(grammar),
	}

	space := roc.Space{
		Identifier: "space://myspace",
		Endpoints: []roc.Endpoint{
			endpoint,
		},
	}

	k := roc.NewKernel()
	k.Spaces = append(k.Spaces, space)

	ctx := roc.NewRequestContext(context.Background(), "res://my-resource", roc.Sink)

	rep := k.Dispatch(ctx)
	fmt.Println(rep)

}
