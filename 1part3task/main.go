package main

import (
	"encoding/json"
	"fmt"

	"log"
	"net/http"

	"github.com/gorilla/mux"
	
	
)


type requestBody struct {
	Task string `json:"task"`
	IsDone string `json:"progress"` // Именно эти поля в json будут считываться
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
	
	task:= Message{Task: requestBody.Task, IsDone: requestBody.IsDone,}
	
	if err:= DB.Create(&task).Error; err!=nil {
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
        return
	}

		
	fmt.Fprintln(w, "Задача и прогресс добавлены в БД")
	
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

func PutMessages(w http.ResponseWriter, r *http.Request) {
// мы должны обработать JSON файл который к нам приходит и декодировать в переменную requestBody
	var requestBody requestBody
	err:= json.NewDecoder(r.Body).Decode(&requestBody)
	if err!=nil {
		http.Error(w,"error in JSON",http.StatusBadGateway)
		return
	} 
	
	// Используя mux.Vars(r), мы получаем значение id из URL
	vars:= mux.Vars(r)
	id:=vars["id"]
	
	// переменная которая работает со слобцами БД
	var message Message
	
	// ищем запись с указанным ID
	if err:= DB.First(&message, id).Error; err!= nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	// Обновляем поля записи

	if requestBody.Task != "" {
        message.Task = requestBody.Task
    }
    if requestBody.IsDone != "" {
        message.IsDone = requestBody.IsDone
    }
	

	if err:= DB.Save(&message).Error; err!= nil {
		http.Error(w,"Ошибка обновления данных", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Запись успешно обновлена")


}

func DeleteMessages(w http.ResponseWriter, r *http.Request) {


}


func main() {
	// Вызываем метод InitDB() из файла db.go
	InitDB()
	
	// Автоматическая миграция модели Message
	DB.AutoMigrate(&Message{}) //  Gorm создаёт в базе данных таблицу с именем messages

	router := mux.NewRouter()

	router.HandleFunc("/api/messages", PostMessage).Methods("POST")
	router.HandleFunc("/api/messages", GetMessages).Methods("GET")
	router.HandleFunc("/api/messages/{id}", PutMessages).Methods("PUT")

	log.Println("Server is staring with DB on :8080 port...")
	http.ListenAndServe(":8080", router)
	

}