package roc

import (
	"net/rpc"
)

// Here is an implementation that talks over RPC
type DispatcherRPC struct {
	client *rpc.Client
}

// CanResolve responds affirmatively if the endpoint can handle the request based on the identifier
func (d *DispatcherRPC) Dispatch(ctx *RequestContext) (Representation, error) {
	log.Debug("making dispatch RPC call",
		"identifier", ctx.Request.Identifier,
		"scope_sze", len(ctx.Scope.Spaces),
	)

	var resp Representation
	err := d.client.Call("Plugin.Dispatch", ctx, &resp)
	if err != nil {
		log.Error("failied to make dispatch rpc call", "error", err)
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp, err
}

// Here is the RPC server that DispatchRPC talks to, conforming to
// the requirements of net/rpc
type DispatcherRPCServer struct {
	// This is the real implementation
	Impl Dispatcher
}

func (s *DispatcherRPCServer) Dispatch(ctx *RequestContext, resp *Representation) error {
	log.Debug("received dispatch call",
		"identifier", ctx.Request.Identifier,
		"scope_sze", len(ctx.Scope.Spaces),
	)
	rep, err := s.Impl.Dispatch(ctx)
	*resp = rep
	return err

}
