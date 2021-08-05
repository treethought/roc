package roc

type Dispatcher interface {
	Dispatch(ctx *RequestContext) (Representation, error)
}

