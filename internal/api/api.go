package api

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/powerman/rpc-codec/jsonrpc2"
	"github.com/taipoxin/json-rpc-pg/internal/api/handlers"
	"github.com/taipoxin/json-rpc-pg/internal/api/models"
)

// Start - run JSON-RPC 2.0 server on param:addr
func Start(addr string) {

	dbContainer := models.EstablishConnection()
	mainHandler := handlers.Main{
		dbContainer,
	}

	server := rpc.NewServer()
	server.Register(&mainHandler)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Println("JSON-RPC 2.0 server is listening on addr", addr)

	// JSON-RPC 2.0 over HTTP
	// var 2 - advanced
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
