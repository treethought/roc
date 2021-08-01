package main

import (
	"fmt"

	"github.com/treethought/roc"
)

func main() {

	k := roc.NewKernel()

	spaces := roc.LoadSpaces("examples/config.yaml")
	k.Register(spaces...)

	ctx := roc.NewRequestContext("res://hello-world", roc.Source)

	rep, err := k.Dispatch(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(rep)

}
