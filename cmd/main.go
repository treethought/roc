package main

import (
	"fmt"
	"net/url"

	"github.com/treethought/roc"
)

type FakeEndpoint struct {
	*roc.Accessor
}

func (e FakeEndpoint) Source(request *roc.Request) roc.Representation {
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

	space := &roc.Space{
		Identifier: "space://myspace",
		Endpoints: []roc.Endpoint{
			endpoint,
		},
	}

	k := roc.NewKernel()
	k.Spaces = append(k.Spaces, space)

	request := roc.NewRequest("res://my-resource", roc.Sink, nil)

	rep := k.Dispatch(request)
	fmt.Println(rep)
    k.Serve(8765)

}
