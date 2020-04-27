package api

import (
	"log"
	"testing"

	"github.com/powerman/rpc-codec/jsonrpc2"
	"github.com/taipoxin/json-rpc-pg/internal/api/handlers"
)

// test JSON-RPC 2.0
func clientHelloHelper(testStr string, addr string) string {
	clientHTTP := jsonrpc2.NewHTTPClient("http://" + addr + "/rpc")
	defer clientHTTP.Close()
	var reply string

	// async request
	startCall := clientHTTP.Go("Test.Hello",
		handlers.HelloArgs{Name: testStr}, &reply, nil)

	// fetch result
	replyCall := <-startCall.Done
	callRes := *replyCall.Reply.(*string)
	log.Printf("Test.Hello({Name:'%s'})=%s\n", testStr, callRes)
	return callRes
}

func TestStart_Hello(t *testing.T) {
	testStr := "Miker"
	addr := ":8888"

	go Start(addr)

	res := clientHelloHelper(testStr, addr)
	if res != "Hello "+testStr {
		t.Error("RPC not correct")
	}
}
