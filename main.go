package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var Message string = "Hello World" //  хранит последнее полученное сообщение , полученного от клиента


type requestBody struct {
	Message string `json:"message"` // хранения переданного JSON-поля "message"
}



func HandlerPost(w http.ResponseWriter, r *http.Request) { // ОТВЕТ ЗАПРОС  содержащую всю информацию о запросе (например, метод, заголовки, тело).
	var requestBody requestBody

	err:= json.NewDecoder(r.Body).Decode(&requestBody)
	Message = requestBody.Message  

	if err != nil {
		http.Error(w, "Ошибка в JSON",http.StatusBadRequest)
		return
	}
	log.Printf("Received message: %s", Message) 
	fmt.Fprintf(w, "Message received: %s", Message)
	
}
			

func HandlerGET(w http.ResponseWriter , r *http.Request) {
	if Message == "" {
		http.Error(w, "Message not found" , http.StatusNotFound)
		return
	}
	fmt.Fprintln(w, "Last received message:" , Message)
}



func main() {
	// var Message int
	router:= mux.NewRouter()

	router.HandleFunc("/api/hello", HandlerGET).Methods("GET")
	router.HandleFunc("/api/hello", HandlerPost).Methods("POST")
	
	log.Println("Server starting on :8080...")
	http.ListenAndServe(":8080", router )
	
}


