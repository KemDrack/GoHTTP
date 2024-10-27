package main

import (
	
	"net/http"

	
	"github.com/gorilla/mux"
	
)



func main() {
	// Вызываем метод InitDB() из файла db.go
	InitDB()
	
	// Автоматическая миграция модели Message
	DB.AutoMigrate(&Message{}) //  Gorm создаёт в базе данных таблицу с именем messages

	router := mux.NewRouter()
	router.HandleFunc("/api/messages", CreateMessage).Methods("POST")
	router.HandleFunc("/api/messages", GetMessages).Methods("GET")
	http.ListenAndServe(":8080", router)
}

func GetMessages(w http.ResponseWriter, r *http.Request) {

}


func CreateMessage(w http.ResponseWriter, r *http.Request) {

}