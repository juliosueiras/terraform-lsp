package jrpc2

import (
	"context"
)

// InboundRequest returns the inbound request associated with the given
// context, or nil if ctx does not have an inbound request. The context passed
// to the handler by *jrpc2.Server will include this value.
//
// This is mainly useful to wrapped server methods that do not have the request
// as an explicit parameter; for direct implementations of Handler.Handle the
// request value returned by InboundRequest will be the same value as was
// passed explicitly.
func InboundRequest(ctx context.Context) *Request {
	if v := ctx.Value(inboundRequestKey{}); v != nil {
		return v.(*Request)
	}
	return nil
}

type inboundRequestKey struct{}

// ServerFromContext returns the server associated with the given context.
// This will be populated on the context passed to request handlers.
// This function is for use by handlers, and will panic for a non-handler context.
//
// It is safe to retain the server and invoke its methods beyond the lifetime
// of the context from which it was extracted; however, a handler must not
// block on the Wait or WaitStatus methods of the server, as the server will
// deadlock waiting for the handler to return.
func ServerFromContext(ctx context.Context) *Server { return ctx.Value(serverKey{}).(*Server) }

type serverKey struct{}
