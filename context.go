package roc

import (
	"context"
)

type RequestContext struct {
	context.Context
	Request  *Request
	Response Response
}

func NewRequestContext(ctx context.Context, identifier Identifier, verb Verb) RequestContext {
	return RequestContext{
		Context: ctx,
		Request: &Request{
			identifier: identifier,
			verb:       verb,
		},
	}
}

func (c *RequestContext) SetResponse(resp Response) {
	c.Response = resp
}
