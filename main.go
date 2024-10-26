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

	if err != nil {
		http.Error(w, "Ошибка в JSON",http.StatusBadRequest)
		return
	}

	Message = requestBody.Message  
	log.Printf("Received message: %s", Message) 
	fmt.Fprintf(w, "Message received: %s", Message)

}
			
func HandlerGET(w http.ResponseWriter , r *http.Request) {
	if Message == "" {
		http.Error(w, "Message not found" , http.StatusNotFound)
		return
	}

	log.Printf("Received message: %s", Message) 
	fmt.Fprintln(w, "Last received message:" , Message)
}

func HandleDelete( w http.ResponseWriter, r *http.Request) {

	if Message == "" {
		http.Error(w, "Message already empty", http.StatusNotFound)
	}

	Message = ""
	log.Printf("DELETE %s", Message )
	fmt.Fprintln(w, "Данные удалены")
}

func HandlerPut(w http.ResponseWriter, r *http.Request) {
	var requestBody requestBody

	err:= json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Ошибка в JSON", http.StatusBadRequest)
		return
	}

	Message = requestBody.Message
	log.Printf("Update Message: %s" , Message)
	fmt.Fprintf(w, "Message update to %s", Message)

}



func main() {
	// var Message int
	router:= mux.NewRouter()

	router.HandleFunc("/api/hello", HandlerGET).Methods("GET")
	router.HandleFunc("/api/hello", HandlerPost).Methods("POST")
	router.HandleFunc("/api/hello", HandleDelete).Methods("DELETE")
	router.HandleFunc("/api/hello", HandlerPost).Methods("PUT")
	
	log.Println("Server starting on :8080...")
	http.ListenAndServe(":8080", router )
	
}


