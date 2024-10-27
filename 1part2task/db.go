package main

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "log"
)
// переменная, через которую мы будем работать с БД
var DB *gorm.DB // Эта переменная будет использоваться для хранения подключения к базе данных

func InitDB() {
// в dsn вводим данные, которые мы указали при создании контейнера
    dsn := "host=localhost user=postgres password=yourpassword dbname=2task port=5432 sslmode=disable"
    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database: ", err)
    }
}
