package main

import (
	"log"
	"os"

	"github.com/taipoxin/json-rpc-pg/internal/api"
	_ "github.com/taipoxin/json-rpc-pg/internal/demo_templates/httprpc"
	_ "github.com/taipoxin/json-rpc-pg/internal/demo_templates/httprpc2"

	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func main() {
	// httprpc.Call()
	// httprpc2.Call()

	addr := os.Getenv("SERVER_ADDR")
	api.Start(addr)
}
