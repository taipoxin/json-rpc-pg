package models

import (
	"database/sql"
	"log"
	"os"

	// pg driver
	_ "github.com/lib/pq"
)

// Datastore presenting interface with methods for handlers
type Datastore interface {
	AllPosts() ([]*Post, error)
	GetPost(id int64) (*Post, error)
	AddPost(title string) (int64, error)
	UpdatePost(id int64, title string) (bool, error)
	DeletePost(id int64) (bool, error)
}

// DbHelper presenting helper for sql.DB, implement Datastore
type DbHelper struct {
	*sql.DB
}

// EstablishConnection use DATABASE_TYPE and return DB with sql.DB inside
func EstablishConnection() *DbHelper {
	var db *sql.DB
	dbType := os.Getenv("DATABASE_TYPE")
	switch dbType {
	case "postgres":
		db = connectPostgres()
	default:
		log.Fatal("invalid .env:DATABASE_TYPE, available: postgres")
	}
	return &DbHelper{db}
}

func connectPostgres() *sql.DB {
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")

	user := os.Getenv("PG_USER")
	pass := os.Getenv("PG_PASS")

	dbname := os.Getenv("PG_DB")

	connStr := "postgres://" + user + ":" + pass + "@" + host + ":" + port + "/" + dbname + "?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("connected to mysql db on %s:%s", host, port)

	return db
}
