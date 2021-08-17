package roc

import (
	"fmt"

	proto "github.com/treethought/roc/proto/v1"
)

type RequestContext struct {
	m *proto.RequestContext
}

// NewRequestContext creates a new request and context with empty scope
func NewRequestContext(identifier Identifier, verb proto.Verb) *RequestContext {
	req := NewRequest(identifier, verb, "")
	return &RequestContext{
		m: &proto.RequestContext{
			Request: req.m,
			Scope:   &proto.RequestScope{},
		},
	}
}
func (c *RequestContext) Request() *Request {
	req := &Request{m: c.m.Request}
	// log.Error("context type", "type", c.m.ProtoReflect().Descriptor().Name())
	return req
}

func (c *RequestContext) setRequest(req *Request) {
	c.m.Request = req.m
}

// CreateRequest returns a new request that can be issued to obtain resources or use services
func (c *RequestContext) CreateRequest(identifier Identifier) *Request {
	return NewRequest(identifier, proto.Verb_VERB_SOURCE, "")
}

func (c *RequestContext) CreateRequestToEndpoint(endpointID string, args map[string][]string) (*Request, error) {
	eId := NewIdentifier(endpointID)
	mReq := NewRequest(eId, proto.Verb_VERB_META, "")
	metaRep, err := c.IssueRequest(mReq)
	if err != nil {
		log.Error("failed to get endpoint meta", "endpoint", endpointID, "err", err)
		return nil, err
	}

	meta := new(proto.EndpointMeta)
	err = metaRep.To(meta)
	if err != nil {
		log.Error("unexpected meta representation", "endpoint", endpointID, "err", err)
		return nil, err
	}
	identifer := constructIdentifier(meta.GetGrammar(), args)

	return NewRequest(NewIdentifier(identifer), proto.Verb_VERB_SOURCE, ""), nil
}

func (c *RequestContext) Scope() *proto.RequestScope {
	return c.m.GetScope()
}

// InjectSpace adds the given space to the request scope
func (c *RequestContext) InjectSpace(space *proto.Space) {
	log.Debug("injecting space into scope", "space", space.GetIdentifier(), "size", len(space.GetEndpoints()))
	c.m.Scope.Spaces = append(c.m.Scope.Spaces, space)
}

// GetArgument retrieves the identifier for named argument provided by the request.
// To retrieve the representation, use GetArgumentValue
func (c *RequestContext) GetArgument(name string) Identifier {
	arg, ok := c.Request().m.GetArguments()[name]
	if !ok {
		return NewIdentifier("")
	}
	id := arg.Values[0]

	return NewIdentifier(id)
}

// GetArgumentValue sources the identifier of the named argument to obtain it's representation
func (c *RequestContext) GetArgumentValue(name string) (Representation, error) {
	log.Info("sourcing argument value", "name", name)
	identifier := c.GetArgument(name)
	if identifier.String() == "" {
		log.Error("argument identifier is empty", "arg_name", name)
		return NewRepresentation(nil), fmt.Errorf("argument not set")
	}

	rep, err := c.Source(identifier, "")
	if err != nil {
		log.Error("failed to source argument", "arg_name", name, "identifier", identifier.String(), "err", err)
		return NewRepresentation(nil), err
	}
	return rep, nil
}

func (c *RequestContext) injectValueSpace(req *Request) {
	// create dynamic pass-by-value space to hold the representation
	// we then provide the uri to the new dynamic endpoint as the arg value in the request

	defs := []*proto.EndpointMeta{}

	for k, val := range req.m.ArgumentValues {
		// create pbv endpoint to hold representation

		rep := NewRepresentation(val)

		ed := &proto.EndpointMeta{Literal: rep.m}
		pbvEndpoint := NewTransientEndpoint(ed)

		log.Info("created transient pbv endpoint",
			"arg", k, "type", rep.Type(),
			"identifier", pbvEndpoint.Identifier(),
		)

		// set the argument value to the pbv identifier
		id := NewIdentifier(pbvEndpoint.Identifier())
		req.SetArgument(k, id)
		defs = append(defs, pbvEndpoint.Meta())
	}

	if len(defs) > 0 {
		log.Debug("injecting argument value space")
		id := NewIdentifier(fmt.Sprintf("pbv://%s", req.Identifier().String()))
		valSpace := NewSpace(id, defs...)
		c.InjectSpace(valSpace)
	}
}

func (c *RequestContext) IssueRequest(req *Request) (Representation, error) {
	log.Trace("issuing new request", "identifier", req.Identifier().String())

	newReqCtx := NewRequestContext(req.Identifier(), c.Request().m.Verb)
	newReqCtx.setRequest(req)
	newReqCtx.m.Scope = c.m.Scope

	d := NewCoreDispatcher()

	resp, err := d.Dispatch(newReqCtx)
	if err != nil {
		log.Error("failed to dispatch request", "err", err)
		return NewRepresentation(nil), err
	}
	return resp, nil
}

// // Source is a helper method to create and issue a new SOURCE request for the identifier
func (c *RequestContext) Source(identifier Identifier, class RepresentationClass) (Representation, error) {
	req := c.CreateRequest(identifier)
	req.SetRepresentationClass(class.String())
	return c.IssueRequest(req)
}
