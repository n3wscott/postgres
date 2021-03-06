package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"log"

	bindingsql "github.com/mattmoor/bindings/pkg/sql"

	_ "database/sql"

	_ "github.com/lib/pq"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", handler)
	log.Printf("ce-reader: listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Print("ce-reader: received a request")
	db, err := bindingsql.Open(context.TODO(), "postgres")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT first_name, last_name, email FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var first, last, email string
		err = rows.Scan(&first, &last, &email)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%q %q %q\n", first, last, email)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
