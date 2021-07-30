package roc

type RequestContext struct {
	Request *Request

	// TODO set as contetx value instead
	Dispatcher Dispatcher
}

func NewRequestContext(identifier Identifier, verb Verb) *RequestContext {
	req := NewRequest(identifier, verb, nil)
	return &RequestContext{
		// Context: ctx,
		Request: req,
	}
}

// CreateRequest returns a new request that can be issued to obtain resources or use services
func (c *RequestContext) CreateRequest(identifier Identifier) *Request {
	return NewRequest(identifier, Source, nil)
}

func (c *RequestContext) IssueRequest(req *Request) (Representation, error) {
	// newCtx, done := context.WithCancel(c.Context)
	// defer done()

	newReqCtx := NewRequestContext(req.Identifier, c.Request.Verb)
	newReqCtx.Dispatcher = c.Dispatcher

	return c.Dispatcher.Dispatch(newReqCtx)

}

// // Source is a helper method to create and issue a new SOURCE request for the identifier
func (c *RequestContext) Source(identifier Identifier, class RepresentationClass) {
	req := c.CreateRequest(identifier)
	req.SetRepresentationClass(class)
	return
}
