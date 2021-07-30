package main

import (
	"fmt"
	"time"

	"github.com/treethought/roc"
)

func main() {

	k := roc.NewKernel()

	endpoint := roc.NewPhysicalEndpoint("./plugin/my_endpoint")

	space := roc.NewSpace("space://myspace", endpoint)

	rstart := time.Now()
	k.Register(space)
	fmt.Println("reg dur")
	fmt.Println(time.Since(rstart).String())

	ctx := roc.NewRequestContext("res://hello-world", roc.Source)

	fmt.Println("dispatching request")
	start := time.Now()
	rep, err := k.Dispatch(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(rep)
	fmt.Println(time.Since(start).String())

}
