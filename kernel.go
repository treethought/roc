package roc

import "log"

type Kernel struct {
	Spaces   []Space
	receiver chan (RequestContext)
	issuer   chan (RequestContext)
}

func NewKernel() Kernel {
	return Kernel{
		Spaces:   []Space{},
		receiver: make(chan RequestContext),
		issuer:   make(chan RequestContext),
	}
}

func (k Kernel) startReceiver() {
	for {
		incoming := <-k.receiver
		k.Dispatch(incoming)
	}
}

func (k Kernel) buildResolveRequest(ctx RequestContext) Request {
	return Request{
		identifier: ctx.Request.Identifier(),
		verb:       Resolve,
		// TODO: make interface or someting for rep classes
		// representationClass: ClassResolution{},
	}

}

func (k Kernel) resolveEndpoint(ctx RequestContext) EndpointInteface {
	c := make(chan (EndpointInteface))
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
