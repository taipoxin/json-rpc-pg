package httprpc

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/rpc"

	"github.com/powerman/rpc-codec/jsonrpc2"
)

// A server wishes to export an object of type ExampleSvc:
type ExampleSvc struct{}

// Method with positional params.
func (*ExampleSvc) Sum(vals [2]int, res *int) error {
	*res = vals[0] + vals[1]
	return nil
}

// Method with positional params.
func (*ExampleSvc) SumAll(vals []int, res *int) error {
	for _, v := range vals {
		*res += v
	}
	return nil
}

// Method with named params.
func (*ExampleSvc) MapLen(m map[string]int, res *int) error {
	*res = len(m)
	return nil
}

type NameArg struct{ Fname, Lname string }
type NameRes struct{ Name string }

// Method with named params.
func (*ExampleSvc) FullName(t NameArg, res *NameRes) error {
	*res = NameRes{t.Fname + " " + t.Lname}
	return nil
}

type exampleContextKey string

var RemoteAddrContextKey exampleContextKey = "RemoteAddr"

type NameArgContext struct {
	Fname, Lname string
	jsonrpc2.Ctx
}

// Method with named params and HTTP context.
func (*ExampleSvc) FullName3(t NameArgContext, res *NameRes) error {
	host, _, _ := net.SplitHostPort(jsonrpc2.HTTPRequestFromContext(t.Context()).RemoteAddr)
	fmt.Printf("FullName3(): Remote IP is %s\n", host)
	*res = NameRes{t.Fname + " " + t.Lname}
	return nil
}

// Method returns error with code -32000.
func (*ExampleSvc) Err1(struct{}, *struct{}) error {
	return errors.New("some issue")
}

// Method returns error with code 42.
func (*ExampleSvc) Err2(struct{}, *struct{}) error {
	return jsonrpc2.NewError(42, "some issue")
}

// Method returns error with code 42 and extra error data.
func (*ExampleSvc) Err3(struct{}, *struct{}) error {
	return &jsonrpc2.Error{42, "some issue", []string{"one", "two"}}
}

func Call() {
	// Server export an object of type ExampleSvc.
	rpc.Register(&ExampleSvc{})

	// Server provide a HTTP transport on /rpc endpoint.
	http.Handle("/rpc", jsonrpc2.HTTPHandler(nil))
	lnHTTP, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	defer lnHTTP.Close()
	go http.Serve(lnHTTP, nil)

	// Server end

	fmt.Println(lnHTTP.Addr().String())

	// Client use HTTP transport.
	clientHTTP := jsonrpc2.NewHTTPClient("http://" + lnHTTP.Addr().String() + "/rpc")
	defer clientHTTP.Close()

	// Custom client use HTTP transport.
	clientCustomHTTP := jsonrpc2.NewCustomHTTPClient(
		"http://"+lnHTTP.Addr().String()+"/rpc",
		jsonrpc2.DoerFunc(func(req *http.Request) (*http.Response, error) {
			// Setup custom HTTP client.
			client := &http.Client{}
			// Modify request as needed.
			req.Header.Set("Content-Type", "application/json-rpc")
			return client.Do(req)
		}),
	)
	defer clientCustomHTTP.Close()

	var reply int

	// Synchronous call using positional params and HTTP.
	err = clientHTTP.Call("ExampleSvc.SumAll", []int{3, 5, -2}, &reply) // nolint:ineffassign
	fmt.Printf("SumAll(3,5,-2)=%d\n", reply)

	// Notification using named params and HTTP.
	clientHTTP.Notify("ExampleSvc.FullName", NameArg{"First", "Last"})

	// Synchronous call using named params and HTTP with context.
	clientHTTP.Call("ExampleSvc.FullName3", NameArg{"First", "Last"}, nil)

	err = jsonrpc2.WrapError(clientCustomHTTP.Call("ExampleSvc.Err2", nil, nil))
	if rpcerr := new(jsonrpc2.Error); errors.As(err, &rpcerr) {
		fmt.Printf("Err2(): code=%d msg=%q data=%v\n", rpcerr.Code, rpcerr.Message, rpcerr.Data)
	} else if err != nil {
		fmt.Printf("Err2(): %q\n", err)
	}

	err = clientHTTP.Call("ExampleSvc.Err3", nil, nil)
	if err == rpc.ErrShutdown || err == io.ErrUnexpectedEOF {
		fmt.Printf("Err3(): %q\n", err)
	} else if err != nil {
		rpcerr := jsonrpc2.ServerError(err)
		fmt.Printf("Err3(): code=%d msg=%q data=%v\n", rpcerr.Code, rpcerr.Message, rpcerr.Data)
	}

}
