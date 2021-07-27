package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/treethought/roc"
)

type FakeEndpoint struct {
	roc.Endpoint
}

func (e FakeEndpoint) Source(ctx roc.RequestContext) roc.Representation {
	return "YO"
}

func main() {

	url, err := url.Parse("res://my-resource")
	if err != nil {
		panic(err)
	}

	endpoint := FakeEndpoint{
		Endpoint: roc.Endpoint{
			Grammar: roc.Grammar{
				Base: url,
			},
		},
	}

	space := roc.Space{
		Identifier: "space://myspace",
		Endpoints: []roc.EndpointInteface{
			endpoint,
		},
	}

	k := roc.NewKernel()
	k.Spaces = append(k.Spaces, space)

	ctx := roc.NewRequestContext(context.Background(), "res://my-resource", roc.Sink)

	rep := k.Dispatch(ctx)
	fmt.Println(rep)

}
