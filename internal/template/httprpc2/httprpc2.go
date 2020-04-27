package httprpc2

// copied from https://morphs.ru/posts/2017/06/15/go-json-rpc-server

/*
curl -H "Content-Type: application/json" -X POST -d \
 '{"jsonrpc": "2.0", "method": "Test.Hello", "params:{"Name":"Mike"}, "id": "1"}' \
 http://127.0.0.1:8080/rpc
*/

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/powerman/rpc-codec/jsonrpc2"
)

type HttpConn struct {
	in  io.Reader
	out io.Writer
}

func (c *HttpConn) Read(p []byte) (n int, err error) {
	return c.in.Read(p)
}

func (c *HttpConn) Write(d []byte) (n int, err error) {
	return c.out.Write(d)
}

func (c *HttpConn) Close() error {
	return nil
}

type Test struct{}

type HelloArgs struct {
	Name string
}

func (test *Test) Hello(args *HelloArgs, result *string) error {
	*result = "Hello " + args.Name
	return nil
}

func Call() {

	server := rpc.NewServer()
	server.Register(&Test{})

	var port string = ":8080"

	listener, err := net.Listen("tcp", port)

	if err != nil {
		panic(err)
	}

	defer listener.Close()

	log.Println("JSON-RPC 2.0 server is listening on port", port)
	http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/rpc" {
			serverCodec := jsonrpc2.NewServerCodec(&HttpConn{in: r.Body, out: w}, server)

			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(200)

			if err1 := server.ServeRequest(serverCodec); err1 != nil {
				http.Error(w, "Error while serving JSON request", 500)
				return
			}
		} else {
			http.Error(w, "Unknown request", 404)
		}
	}))
}
