package server

import (
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
)

// Run starts a server for svc on the given channel, and blocks until it
// returns.  The server exit status is reported to the service, and the error
// value is returned.
//
// If the caller does not need the error value and does not want to wait for
// the server to complete, call Run in a goroutine.
func Run(ch channel.Channel, svc Service, opts *jrpc2.ServerOptions) error {
	assigner, err := svc.Assigner()
	if err != nil {
		return err
	}
	srv := jrpc2.NewServer(assigner, opts).Start(ch)
	stat := srv.WaitStatus()
	svc.Finish(assigner, stat)
	return stat.Err
}
