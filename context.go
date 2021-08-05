package roc

import "fmt"

type RequestScope struct {
	Spaces []Space
}

type RequestContext struct {
	Request    *Request
	Dispatcher Dispatcher
	Scope      RequestScope
}

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

func (c *RequestContext) IssueRequest(req *Request) (Representation, error) {
	log.Info("issuing new request")

	newReqCtx := NewRequestContext(req.Identifier, c.Request.Verb)
	newReqCtx.Scope = c.Scope
	newReqCtx.Dispatcher = c.Dispatcher

	if c.Dispatcher == nil {
		return nil, fmt.Errorf("context dispatcher is nil")
	}

	resp, err := c.Dispatcher.Dispatch(newReqCtx)
	if err != nil {
		log.Error("failed to disptach with request context dispatcher", "err", err)
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
