package roc

import "fmt"

type RequestScope struct {
	Spaces []Space
    EndpointClients []Endpoint

}

type RequestContext struct {
	Request *Request

	// TODO set as contetx value instead
	// Dispatcher Dispatcher
	Dispatcher DispatcherClient
	Scope RequestScope
}

func NewRequestContext(identifier Identifier, verb Verb) *RequestContext {
	req := NewRequest(identifier, verb, nil)
	return &RequestContext{
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
	// if len(c.Dispatcher.Spaces) == 0 {
	// 	return nil, fmt.Errorf("dispatcher has no spaces to resolve")
	// }

	newReqCtx := NewRequestContext(req.Identifier, c.Request.Verb)
	newReqCtx.Dispatcher = c.Dispatcher
    newReqCtx.Scope = c.Scope
    if len(newReqCtx.Scope.Spaces) == 0 {
        return nil, fmt.Errorf("request scope has no spaces")
    }
    newReqCtx.Scope.EndpointClients = c.Scope.EndpointClients


	return c.Dispatcher.Dispatch(newReqCtx)

}

// // Source is a helper method to create and issue a new SOURCE request for the identifier
func (c *RequestContext) Source(identifier Identifier, class RepresentationClass) (Representation, error) {
	req := c.CreateRequest(identifier)
	req.SetRepresentationClass(class)
	return c.IssueRequest(req)
}
