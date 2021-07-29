package main

import (
	"fmt"

	"github.com/treethought/roc"
)

func main() {

	k := roc.NewKernel()

	endpoint := roc.NewPhysicalEndpoint("./plugin/my_endpoint")

	space := roc.NewSpace("space://myspace", endpoint)

	k.Register(space)

	request := roc.NewRequest("res://hello-world", roc.Sink, nil)

	rep := k.Dispatch(request)
	fmt.Println(rep)

}
