package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	router:= mux.NewRouter()
	router.HandleFunc("/api/hello", HelloHandler).Methods("GET")
	http.ListenAndServe(":8080", router )

}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
	log.Println("testing")
}