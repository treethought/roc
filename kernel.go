package roc

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Kernel struct {
	Spaces   []Space
	receiver chan (RequestContext)
	issuer   chan (RequestContext)
}

func NewKernel() *Kernel {
	return &Kernel{
		Spaces:   []Space{},
		receiver: make(chan RequestContext),
		issuer:   make(chan RequestContext),
	}
}

func (k *Kernel) LoadFromFile(path string) error {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	spaces := []Space{}

	err = yaml.Unmarshal(data, &spaces)
	if err != nil {
		return err
	}

	k.Spaces = append(k.Spaces, spaces...)
	return nil
}

func (k Kernel) startReceiver() {
	for {
		incoming := <-k.receiver
		k.Dispatch(incoming)
	}
}

func (k Kernel) buildResolveRequest(ctx RequestContext) *Request {
	return NewRequest(ctx.Request.Identifier(), Resolve, nil)

}

func (k Kernel) resolveEndpoint(ctx RequestContext) Endpoint {
	c := make(chan (Endpoint))
	for _, s := range k.Spaces {
		go s.MatchEndpoint(ctx, c)
	}

	return <-c
}

func (k Kernel) Dispatch(ctx RequestContext) Representation {
	log.Printf("dispatching request for identifer: %s", ctx.Request.Identifier())

	endpoint := k.resolveEndpoint(ctx)
	log.Printf("resolved to endpoint: %s", endpoint)

	// TODO route verbs to methods
	rep := endpoint.Source(ctx)
	return rep

}
