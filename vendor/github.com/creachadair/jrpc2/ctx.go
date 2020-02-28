package jrpc2

import (
	"context"

	"golang.org/x/xerrors"

	"github.com/creachadair/jrpc2/metrics"
)

// ServerMetrics returns the server metrics collector associated with the given
// context, or nil if ctx doees not have a collector attached.  The context
// passed to a handler by *jrpc2.Server will include this value.
func ServerMetrics(ctx context.Context) *metrics.M {
	if v := ctx.Value(serverMetricsKey{}); v != nil {
		return v.(*metrics.M)
	}
	return nil
}

type serverMetricsKey struct{}

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

// ServerPush posts a server notification to the client. If ctx does not
// contain a server notifier, this reports ErrNotifyUnsupported. The context
// passed to the handler by *jrpc2.Server will support notifications if the
// server was constructed with the AllowPush option set true.
func ServerPush(ctx context.Context, method string, params interface{}) error {
	v := ctx.Value(serverPushKey{})
	if v == nil {
		return ErrNotifyUnsupported
	}
	notify := v.(func(context.Context, string, interface{}) error)
	return notify(ctx, method, params)
}

type serverPushKey struct{}

// ErrNotifyUnsupported is returned by ServerPush if server notifications are
// not enabled in the specified context.
var ErrNotifyUnsupported = xerrors.New("server notifications are not enabled")
