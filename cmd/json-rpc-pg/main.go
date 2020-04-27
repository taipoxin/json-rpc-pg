package main

import (
	"github.com/taipoxin/json-rpc-pg/internal/api"
	_ "github.com/taipoxin/json-rpc-pg/internal/demo_templates/httprpc"
	_ "github.com/taipoxin/json-rpc-pg/internal/demo_templates/httprpc2"
)

func main() {
	// httprpc.Call()
	// httprpc2.Call()

	api.Start(":8080")
}
