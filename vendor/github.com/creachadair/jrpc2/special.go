package jrpc2

import (
	"context"
)

const (
	rpcServerInfo = "rpc.serverInfo"
)

// CancelRequest instructs s to cancel the pending or in-flight request with
// the specified ID. If no request exists with that ID, this is a no-op.
func (s *Server) CancelRequest(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cancel(id) {
		s.log("Cancelled request %s by client order", id)
	}
}

// methodFunc is a replication of handler.Func redeclared to avert a cycle.
type methodFunc func(context.Context, *Request) (interface{}, error)

func (m methodFunc) Handle(ctx context.Context, req *Request) (interface{}, error) {
	return m(ctx, req)
}

// Handle the special rpc.serverInfo method, that requests server vitals.
func (s *Server) handleRPCServerInfo(context.Context, *Request) (interface{}, error) {
	return s.ServerInfo(), nil
}

// RPCServerInfo calls the built-in rpc.serverInfo method exported by servers.
// It is a convenience wrapper for an invocation of cli.CallResult.
func RPCServerInfo(ctx context.Context, cli *Client) (result *ServerInfo, err error) {
	err = cli.CallResult(ctx, rpcServerInfo, nil, &result)
	return
}
