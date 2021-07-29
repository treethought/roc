package main

import (
	"fmt"

	"github.com/treethought/roc"
)

func main() {

	k := roc.NewKernel()
	space := roc.NewSpace("space://myspace")

	k.Register(space, "./plugin/my_endpoint")

	request := roc.NewRequest("res://hello-world", roc.Sink, nil)

	rep := k.Dispatch(request)
	fmt.Println(rep)

}
