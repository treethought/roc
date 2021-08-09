package roc

type RequestScope struct {
	Spaces []Space
}

type RequestContext struct {
	Request *Request
	Scope   RequestScope
}

// NewRequestContext creates a new request and context with empty scope
func NewRequestContext(identifier Identifier, verb Verb) *RequestContext {
	req := NewRequest(identifier, verb, nil)
	return &RequestContext{
		Request: req,
		Scope:   RequestScope{},
	}
}

// CreateRequest returns a new request that can be issued to obtain resources or use services
func (c *RequestContext) CreateRequest(identifier Identifier) *Request {
	return NewRequest(identifier, Source, nil)
}

// InjectSpace adds the given space to the request scope
func (c *RequestContext) InjectSpace(space Space) {
	log.Debug("injecting space into scope", "space", space.Identifier)
	c.Scope.Spaces = append(c.Scope.Spaces, space)
}

// GetArgument retrieves the identifier for named argument provided by the request.
// To retrieve the representation, use GetArgumentValue
func (c *RequestContext) GetArgument(name string) Identifier {
	i, ok := c.Request.Arguments[name]
	if !ok {
		return Identifier("")
	}
	return Identifier(i[0])
}

// GetArgumentValue sources the identifier of the named argument to obtain it's representation
func (c *RequestContext) GetArgumentValue(name string) (Representation, error) {
	identifier := c.GetArgument(name)
	if identifier.String() == "" {
		return nil, nil
	}

	rep, err := c.Source(identifier, nil)
	if err != nil {
		return nil, err
	}
	return rep, nil
}

func (c *RequestContext) injectValueSpace(req *Request) {
	// create dynamic pass-by-value space to hold the representation
	// we then provide the uri to the new dynamic endpoint as the arg value in the request

	log.Debug("injecting argument value space")

	defs := []EndpointDefinition{}

	for k, val := range req.argumentValues {
		// create pbv endpoint to hold representation
		pbvEndpoint := NewTransientEndpoint(val[0])
		log.Info("created transient pbv endpoint", "narg", k, "val", val, "identifier", pbvEndpoint.Identifier())

		// set the argument value to the pbv identifier
		req.SetArgument(k, pbvEndpoint.Identifier())
		defs = append(defs, pbvEndpoint.Definition())
	}

	valSpace := NewSpace("transient://", defs...)
	c.InjectSpace(valSpace)
}

func (c *RequestContext) IssueRequest(req *Request) (Representation, error) {
	log.Debug("issuing new request", "identifier", req.Identifier)

	newReqCtx := NewRequestContext(req.Identifier, c.Request.Verb)
	newReqCtx.Request = req
	newReqCtx.Scope = c.Scope

	d := NewCoreDispatcher()

	resp, err := d.Dispatch(newReqCtx)
	if err != nil {
		log.Error("failed to dispatch request", "err", err)
		return nil, err
	}
	return resp, nil
}

// // Source is a helper method to create and issue a new SOURCE request for the identifier
func (c *RequestContext) Source(identifier Identifier, class RepresentationClass) (Representation, error) {
	req := c.CreateRequest(identifier)
	req.SetRepresentationClass(class)
	return c.IssueRequest(req)
}
