package main //import github.com/thanatkrit289/wongnok

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/thanatkrit289/wongnok/internal/api"
	"github.com/thanatkrit289/wongnok/internal/auth"
)

func main() {
	fmt.Println("Wongnok")
	fmt.Println("version: 1.0.0")

	db, err := sql.Open(
		"postgres",
		"postgres://postgres@localhost:5432/wongnok?sslmode=disable",
		)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Println(err)
	}

	server := http.Server{
		Addr: ":8080",
		Handler: api.API{
			Auth: auth.New(db),
		}.Handler(),
	}

	log.Printf("Server listening on %s\n", server.Addr)
	err = server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
