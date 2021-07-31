package main

import (
	"fmt"
	"time"

	"github.com/treethought/roc"
)

func main() {

	k := roc.NewKernel()
	space := roc.NewSpace("space://myspace",
		"./plugin/greeter/greeter",
		"./plugin/namer/namer",
	)

	k.Register(space)

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
