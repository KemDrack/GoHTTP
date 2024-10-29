package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	
	"github.com/gorilla/mux"
)

var Task string = "Hello World" //  хранит последнее полученное сообщение , полученного от клиента


type requestBody struct {
	Task string `json:"task"` // хранения переданного JSON-поля "message"
}



func HandlerPost(w http.ResponseWriter, r *http.Request) { // ОТВЕТ ЗАПРОС  содержащую всю информацию о запросе (например, метод, заголовки, тело).
	
	var requestBody requestBody
	// decoder:= json.NewDecoder(r.Body)
	err:= json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		http.Error(w, "Ошибка в JSON",http.StatusBadRequest)
		return
	}
	
	Task = requestBody.Task  
	log.Printf("Вы загрузили новую задачу %s", Task) 
	fmt.Fprintf(w, "Message received: %s", Task)

}
			
func HandlerGET(w http.ResponseWriter , r *http.Request) {
	if Task == "" {
		http.Error(w, "Message not found" , http.StatusNotFound)
		return
	}

	log.Printf("Received message: %s", Task) 
	fmt.Fprintln(w, "Вот последняя задача:" , Task)
}

func HandleDelete( w http.ResponseWriter, r *http.Request) {

	if Task == "" {
		http.Error(w, "Message already empty", http.StatusNotFound)
	}

	Task = ""
	log.Printf("DELETE %s", Task )
	fmt.Fprintln(w, "Данные удалены")
}

func HandlerPut(w http.ResponseWriter, r *http.Request) {
	var requestBody requestBody

	err:= json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Ошибка в JSON", http.StatusBadRequest)
		return
	}
	
	Task = requestBody.Task
	log.Printf("Update Message: %s" , Task)
	fmt.Fprintf(w, "Message update to %s", Task)

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


