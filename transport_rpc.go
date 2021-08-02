package roc

import (
	"fmt"
	"net/rpc"
)

// TransportRPC is the RPC client implementation of a Transport
// this behaves identical to EndpointRPC, but exposes an Init method to deliver transport scope
type TransportRPC struct {
	EndpointRPC
	client *rpc.Client
}

// Init is a special method for Transport endpoints to deliver their request scope upon intialization
func (e *TransportRPC) Init(scope RequestScope) error {
	var resp error
	err := e.client.Call("Plugin.Init", scope, &resp)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}

// TransportRPCServer is the server side implementation of Transport
// this behaves identical to EndpointRPCServer, but exposes an Init method to deliver transport scope
type TransportRPCServer struct {
	EndpointRPCServer
	// This is the real implementation
	Impl Transport
}

func (s *TransportRPCServer) Init(scope RequestScope, resp *error) error {
	s.Impl.Init(scope)
	*resp = nil
	return nil
}
