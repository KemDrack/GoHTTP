package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type requestBody struct {
	Task   string `json:"task"`
	IsDone bool `json:"is_done"` // Именно эти поля в json будут считываться
}

func PostMessage(w http.ResponseWriter, r *http.Request) {
	var requestBody requestBody

	err := json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		http.Error(w, "Error in JSON", http.StatusBadRequest)
		return
	}
	if DB == nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}

	task := Message{Task: requestBody.Task, IsDone: requestBody.IsDone}

	if err := DB.Create(&task).Error; err != nil {
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Задача и прогресс добавлены в БД")

}

func GetMessages(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json") // Чтобы вывод клиенту был в виде JSON
	var messages []Message

	if err := DB.Find(&messages).Error; err != nil {
		http.Error(w, "Failed find text", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&messages); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

}

func PutMessages(w http.ResponseWriter, r *http.Request) {
	// мы должны обработать JSON файл который к нам приходит и декодировать в переменную requestBody
	 // Декодируем данные запроса в структуру requestBody
	var requestBody requestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err!= nil {
		http.Error(w, "Error in JSON", http.StatusBadRequest)
		return
	}
		
 
	// Извлекаем ID из URL
	vars := mux.Vars(r) // используется для извлечения переменных из URL
	idStr := vars["id"] // Функция mux.Vars(r) возвращает map[string]string, содержащую параметры из пути, заданные в маршруте. В данном случае, маршрут /api/messages/{id} задает переменную {id}, поэтому mux.Vars(r)["id"] вернет строковое значение ID из URL.
	id, err := strconv.Atoi(idStr) // Например, если запрос отправляется на /api/messages/1, vars["id"] вернет строку "1".
	if err != nil { // Далее strconv.Atoi(idStr) преобразует строковое значение idStr в целое число. Это важно, потому что gorm ожидает числовой идентификатор для поиска записи по первичному ключу.
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
 
	// Находим запись по ID
	var message Message
	if err := DB.First(&message, id).Error; err != nil { // используется для поиска записи Message с указанным ID в базе данных.
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}
 
	// Обновляем только те поля, которые переданы в JSON
	if requestBody.Task != "" {
		message.Task = requestBody.Task
	}
	
	message.IsDone = requestBody.IsDone
	 
	// Сохраняем изменения
	if err := DB.Save(&message).Error; err != nil {
		http.Error(w, "Failed to update message", http.StatusInternalServerError)
		return
	}
 
	fmt.Fprintf(w, "Запись успешно обновлена")

}

func DeleteMessages(w http.ResponseWriter, r *http.Request) {
	log.Println("Received DELETE request")

	vars := mux.Vars(r)
    idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Invalid ID:", idStr)
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

	if err := DB.Delete(&Message{}, id).Error; err != nil {
		log.Println("Failed to delete message:", err)
        http.Error(w, "Failed to delete message", http.StatusInternalServerError)
        return
    }

	fmt.Fprintf(w, "Запись успешно удалена")
	log.Println("Record deleted successfully for ID:", id)



}

func main() {
	// Вызываем метод InitDB() из файла db.go
	InitDB()

	// Автоматическая миграция модели Message
	DB.AutoMigrate(&Message{}) //  GORM проверяет, существует ли таблица, соответствующая этой структуре, и если нет — создает её с колонками task и is_done.

	router := mux.NewRouter()

	router.HandleFunc("/api/messages", PostMessage).Methods("POST")
	router.HandleFunc("/api/messages", GetMessages).Methods("GET")
	router.HandleFunc("/api/messages/{id}", PutMessages).Methods("PUT")
	router.HandleFunc("/api/messages/{id}", DeleteMessages).Methods("DELETE")

	log.Println("Server is staring with DB on :8080 port...")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println(err)
		return
	}

}
