/*
Package jrpc2 implements a server and a client for the JSON-RPC 2.0 protocol
defined by http://www.jsonrpc.org/specification.

Servers

The *Server type implements a JSON-RPC server. A server communicates with a
client over a channel.Channel, and dispatches client requests to user-defined
method handlers.  Handlers satisfy the jrpc2.Handler interface by exporting a
Handle method with this signature:

   Handle(ctx Context.Context, req *jrpc2.Request) (interface{}, error)

The handler package helps adapt existing functions to this interface.
A server finds the handler for a request by looking up its method name in a
jrpc2.Assigner provided when the server is set up.

For example, suppose we would like to export the following Add function as a
JSON-RPC method:

   // Add returns the sum of a slice of integers.
   func Add(ctx context.Context, values []int) int {
      sum := 0
      for _, v := range values {
         sum += v
      }
      return sum
   }

To convert Add to a jrpc2.Handler, call handler.New, which uses reflection to
lift its argument into the jrpc2.Handler interface:

   h := handler.New(Add)  // h is a jrpc2.Handler that invokes Add

We will advertise this function under the name "Add".  For static assignments
we can use a handler.Map, which finds methods by looking them up in a Go map:

   assigner := handler.Map{
      "Add": handler.New(Add),
   }

Equipped with an Assigner we can now construct a Server:

   srv := jrpc2.NewServer(assigner, nil)  // nil for default options

To serve requests, we need a channel.Channel. Implementations of the Channel
interface handle the framing, transmission, and receipt of JSON messages.  The
channel package provides several common framing disciplines and functions to
wrap them around various input and output streams.  For this example, we'll use
a channel that delimits messages by newlines, and communicates on os.Stdin and
os.Stdout:

   ch := channel.Line(os.Stdin, os.Stdout)
   srv.Start(ch)

Once started, the running server handles incoming requests until the channel
closes, or until it is stopped explicitly by calling srv.Stop(). To wait for
the server to finish, call:

   err := srv.Wait()

This will report the error that led to the server exiting.  The code for this
example is available from cmd/examples/adder/adder.go:

    $ go run cmd/examples/adder/adder.go

Interact with the server by sending JSON-RPC requests on stdin, such as for
example:

   {"jsonrpc":"2.0", "id":1, "method":"Add", "params":[1, 3, 5, 7]}


Clients

The *Client type implements a JSON-RPC client. A client communicates with a
server over a channel.Channel, and is safe for concurrent use by multiple
goroutines. It supports batched requests and may have arbitrarily many pending
requests in flight simultaneously.

To create a client we need a channel:

   import "net"

   conn, err := net.Dial("tcp", "localhost:8080")
   ...
   ch := channel.RawJSON(conn, conn)
   cli := jrpc2.NewClient(ch, nil)  // nil for default options

To send a single RPC, use the Call method:

   rsp, err := cli.Call(ctx, "Add", []int{1, 3, 5, 7})

Call blocks until the response is received. Any error returned by the server,
including cancellation or deadline exceeded, has concrete type *jrpc2.Error.

To issue a batch of requests, use the Batch method:

   rsps, err := cli.Batch(ctx, []jrpc2.Spec{
      {Method: "Math.Add", Params: []int{1, 2, 3}},
      {Method: "Math.Mul", Params: []int{4, 5, 6}},
      {Method: "Math.Max", Params: []int{-1, 5, 3, 0, 1}},
   })

Batch blocks until all the responses are received.  An error from the Batch
call reflects an error in sending the request: The caller must check each
response separately for errors from the server. Responses are returned in the
same order as the Spec values, save that notifications are omitted.

To decode the result from a successful response use its UnmarshalResult method:

   var result int
   if err := rsp.UnmarshalResult(&result); err != nil {
      log.Fatalln("UnmarshalResult:", err)
   }

To close a client and discard all its pending work, call cli.Close().


Notifications

The JSON-RPC protocol also supports notifications.  Notifications differ from
calls in that they are one-way: The client sends them to the server, but the
server does not reply.

Use the Notify method of a jrpc2.Client to send notifications:

   err := cli.Notify(ctx, "Alert", handler.Obj{
      "message": "A fire is burning!",
   })

A notification is complete once it has been sent.

On server, notifications are handled identically to ordinary requests, except
that the return value is discarded once the handler returns. If a handler does
not want to do anything for a notification, it can query the request:

   if req.IsNotification() {
      return 0, nil  // ignore notifications
   }


Services with Multiple Methods

The example above shows a server with one method using handler.New.  To
simplify exporting multiple methods, the handler.Map type collects named
methods:

   mathService := handler.Map{
      "Add": handler.New(Add),
      "Mul": handler.New(Mul),
   }

Maps may be further combined with the handler.ServiceMap type to allow
different services to work together:

   func GetStatus(context.Context) (string, error) {
      return "all is well", nil
   }

   assigner := handler.ServiceMap{
      "Math":   mathService,
      "Status": handler.Map{"Get": handler.New(Status)},
   }

This assigner dispatches "Math.Add" and "Math.Mul" to the arithmetic functions,
and "Status.Get" to the GetStatus function. A ServiceMap splits the method name
on the first period ("."), and you may nest ServiceMaps more deeply if you
require a more complex hierarchy.


Concurrency

A Server issues requests to handlers concurrently, up to the Concurrency limit
given in its ServerOptions. Two requests (either calls or notifications) are
concurrent if they arrive as part of the same batch. In addition, two calls are
concurrent if the time intervals between the arrival of the request objects and
delivery of the response objects overlap.

The server may issue concurrent requests to their handlers in any order.
Otherwise, requests are processed in order of arrival. Notifications, in
particular, can only be concurrent with other requests in the same batch.
This ensures a client that sends a notification can be sure its notification
was fully processed before any subsequent calls are issued.

These rules imply that the client cannot rely on the order of evaluation for
calls that overlap: If the caller needs to ensure that call A completes before
call B starts, it must wait for A to return before invoking B.


Built-in Methods

Per the JSON-RPC 2.0 spec, method names beginning with "rpc." are reserved by
the implementation. By default, a server does not dispatch these methods to its
assigner. In this configuration, the server exports a "rpc.serverInfo" method
taking no parameters and returning a jrpc2.ServerInfo value.

Setting the DisableBuiltin option to true in the ServerOptions removes special
treatment of "rpc." method names, and disables the rpc.serverInfo handler.
When this option is true, method names beginning with "rpc." will be dispatched
to the assigner like any other method.


Server Push

The AllowPush option in jrpc2.ServerOptions allows a server to "push" requests
back to the client. This is a non-standard extension of JSON-RPC used by some
applications such as the Language Server Protocol (LSP). If this feature is
enabled, the server's Notify and Callback methods send requests back to the
client. Otherwise, those methods will report an error:

  if err := s.Notify(ctx, "methodName", params); err == jrpc2.ErrPushUnsupported {
    // server push is not enabled
  }
  if rsp, err := s.Callback(ctx, "methodName", params); err == jrpc2.ErrPushUnsupported {
    // server push is not enabled
  }

A method handler may use jrpc2.ServerFromContext to access the server from its
context, and then invoke these methods on it.

On the client side, the OnNotify and OnCallback options in jrpc2.ClientOptions
provide hooks to which any server requests are delivered, if they are set.
*/
package jrpc2

// Version is the version string for the JSON-RPC protocol understood by this
// implementation, defined at http://www.jsonrpc.org/specification.
const Version = "2.0"
