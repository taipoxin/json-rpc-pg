package main

import (
	_ "github.com/taipoxin/json-rpc-pg/internal/template/httprpc"
	"github.com/taipoxin/json-rpc-pg/internal/template/httprpc2"
)

func main() {
	// httprpc.Call()

	httprpc2.Call()
}
