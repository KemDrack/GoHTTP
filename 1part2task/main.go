package main

import (
	"encoding/json"
	"fmt"

	"log"
	"net/http"

	"github.com/gorilla/mux"
	
)



type requestBody struct {
	Text string `json:"text"`

}


func PostMessage(w http.ResponseWriter, r *http.Request) {
	var requestBody requestBody

	err:= json.NewDecoder(r.Body).Decode(&requestBody)

	if err!= nil {
		http.Error(w, "Error in JSON", http.StatusBadRequest)
		return
	}
	if DB == nil {
        http.Error(w, "Database connection error", http.StatusInternalServerError)
        return
    }
	
	message:= Message{Text: requestBody.Text}

	if err:= DB.Create(&message).Error; err!=nil {
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
        return
	}
	
	fmt.Fprintln(w, "Текст добавлен в БД")
	
}


func GetMessages(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type", "application/json") // Чтобы вывод клиенту был в виде JSON

	var messages []Message

	if err:= DB.Find(&messages).Error; err!= nil {
		http.Error(w,"Failed find text", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&messages); err != nil {
        http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
        return
    }


}


func main() {
	// Вызываем метод InitDB() из файла db.go
	InitDB()
	
	// Автоматическая миграция модели Message
	DB.AutoMigrate(&Message{}) //  Gorm создаёт в базе данных таблицу с именем messages

	router := mux.NewRouter()

	router.HandleFunc("/api/messages", PostMessage).Methods("POST")
	router.HandleFunc("/api/messages", GetMessages).Methods("GET")

	log.Println("Server is staring with DB on :8080 port...")
	http.ListenAndServe(":8080", router)
}