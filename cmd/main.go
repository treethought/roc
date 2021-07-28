package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/treethought/roc"
	"github.com/treethought/roc/endpoint"
	"github.com/treethought/roc/kernel"
	"github.com/treethought/roc/space"
)

type FakeEndpoint struct {
	endpoint.Endpoint
}

func (e FakeEndpoint) Source(ctx roc.RequestContext) roc.Representation {
	return "YO"
}

func main() {

	url, err := url.Parse("res://my-resource")
	if err != nil {
		panic(err)
	}

	ep := FakeEndpoint{
		Endpoint: endpoint.Endpoint{
			Grammar: endpoint.Grammar{
				Base: url,
			},
		},
	}

	space := space.Space{
		Identifier: "space://myspace",
		Endpoints: []endpoint.EndpointInteface{
			ep,
		},
	}

	k := kernel.NewKernel()
	k.Spaces = append(k.Spaces, space)

	ctx := roc.NewRequestContext(context.Background(), "res://my-resource", roc.Sink)

	rep := k.Dispatch(ctx)
	fmt.Println(rep)

}
