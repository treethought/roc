package main

import (
	"fmt"
	"time"

	"github.com/treethought/roc"
)

func main() {

	k := roc.NewKernel()

	greeter := roc.NewPhysicalEndpoint("./plugin/greeter/greeter")

	namer := roc.NewPhysicalEndpoint("./plugin/namer/namer")

	space := roc.NewSpace("space://myspace", greeter, namer)

	rstart := time.Now()
	k.Register(space)
	fmt.Println("reg dur")
	fmt.Println(time.Since(rstart).String())

	k.StartDispatcher()

	ctx := roc.NewRequestContext("res://hello-world", roc.Source)
    // ctx.Dispatcher = k.DispatchClient

	start := time.Now()
	rep, err := k.Dispatch(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(rep)
	fmt.Println(time.Since(start).String())

}
