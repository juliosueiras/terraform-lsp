package main

import (
	//"fmt"
	//"github.com/davecgh/go-spew/spew"
	//"github.com/hashicorp/hcl2/hcldec"
	//"github.com/hashicorp/terraform/configs"
	//"strings"
	//"github.com/hashicorp/terraform/configs/configload"
	//"github.com/hashicorp/terraform/lang"
	//"github.com/minamijoyo/tfschema/tfschema"
	//"github.com/sourcegraph/go-lsp"
	"context"
	"crypto/tls"
	"github.com/sourcegraph/jsonrpc2"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"os"
	"strings"
)

func main() {
	//	parser := configs.NewParser(nil)
	//	file, _ := parser.LoadConfigFile("test3.tf")
	//	resourceType := file.ManagedResources[0].Type
	//
	//	provider, err := tfschema.NewClient(strings.Split(resourceType, "_")[0])
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	//
	//	provider_resource, err := provider.GetRawResourceTypeSchema(file.ManagedResources[0].Type)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	//fmt.Println(provider_resource.Block)
	//	//spew.Dump(file.ProviderConfigs)
	//
	//	res2 := provider_resource.Block.DecoderSpec()
	//	res, diags := hcldec.Decode(file.ManagedResources[0].Config, res2, nil)
	//	spew.Dump(res.GoString())
	//	//spew.Dump(res)
	//	spew.Dump(diags)
	//
	//	//spew.Dump(file.ManagedResources[0].Config.Content(nil))
	//
	//	//	for _, resource := range file.ManagedResources {
	//	//		s := &hcl.BodySchema{
	//	//			Blocks: []hcl.BlockHeaderSchema{
	//	//				hcl.BlockHeaderSchema{
	//	//					Type: "list",
	//	//					LabelNames: []string{
	//	//						"Hi",
	//	//					},
	//	//				},
	//	//			},
	//	//		}
	//	//
	//	//		res, diag := resource.Config.Content(s)
	//	//		spew.Dump(res)
	//	//		spew.Dump(diag)
	//	//	}

	// JSONRPC2
	rpc.Response
	newHandler := func() (jsonrpc2.Handler, io.Closer) {
		result := jsonrpc2.HandlerWithError()
		return result, ioutil.NopCloser(strings.NewReader(""))
	}

	listen := func(addr string) (*net.Listener, error) {
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatalf("Could not bind to address %s: %v", addr, err)
			return nil, err
		}
		if os.Getenv("TLS_CERT") != "" && os.Getenv("TLS_KEY") != "" {
			cert, err := tls.X509KeyPair([]byte(os.Getenv("TLS_CERT")), []byte(os.Getenv("TLS_KEY")))
			if err != nil {
				return nil, err
			}
			listener = tls.NewListener(listener, &tls.Config{
				Certificates: []tls.Certificate{cert},
			})
		}
		return &listener, nil
	}

	lis, _ := listen("0.0.0.0:8888")
	var connOpt []jsonrpc2.ConnOpt

	defer (*lis).Close()
	for {
		conn, _ := (*lis).Accept()
		handler, closer := newHandler()

		jsonrpc2Connection := jsonrpc2.NewConn(context.Background(), jsonrpc2.NewBufferedStream(conn, jsonrpc2.VSCodeObjectCodec{}), handler, connOpt...)
		go func() {
			<-jsonrpc2Connection.DisconnectNotify()
			err := closer.Close()
			if err != nil {
				log.Println(err)
			}
		}()
	}
}
