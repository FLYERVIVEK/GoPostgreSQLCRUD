package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/flyervivek/golangpostgree/router"
)

func main() {
	r := router.Router()
	log.Fatal(http.ListenAndServe(":8080", r))
	fmt.Println("Server started serving on port 8080")
}
