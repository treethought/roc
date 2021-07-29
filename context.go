package roc

// type RequestContext struct {
// 	context.Context
// 	Request  *Request
// 	Response Response
// 	kChannel chan (*RequestContext)
// }

// func NewRequestContext(ctx context.Context, identifier Identifier, verb Verb) RequestContext {
// 	return RequestContext{
// 		Context: ctx,
// 		Request: &Request{
// 			identifier: identifier,
// 			verb:       verb,
// 		},
// 	}
// }

// func (c *RequestContext) SetResponse(resp Response) {
// 	c.Response = resp
// }

// // CreateRequest returns a new request that can be issued to obtain resources or use services
// func (c *RequestContext) CreateRequest(identifier Identifier) *Request {
// 	return NewRequest(identifier, Sink, nil)
// }

// // func (c *RequestContext) IssueRequest(*Request) Representation {
// // 	if c.kChannel == nil {
// // 		log.Println("Kernel dispatch channel is nil")
// // 		return nil
// // 	}

// // 	c.kChannel <- c

// // }

// // Source is a helper method to create and issue a new SOURCE request for the identifier
// func (c *RequestContext) Source(identifier Identifier) {
// 	return
// }
