package main

import "gorm.io/gorm"

type Message struct {
	gorm.Model // добавляем поля ID, CreatedAt, UpdatedAt и DeletedAt
	Text string `json:"text"` // Наш сервер будет ожидать json c полем text
}

//  структура нашего Message для БД, то, какие столбцы будут в нашей БД