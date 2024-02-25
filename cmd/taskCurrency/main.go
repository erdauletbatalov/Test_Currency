package main

import (
	"log"
	"net/http"

	"currency/config"
	"currency/delivery/web"
	"currency/repository"
	"currency/repository/mssql"
	"currency/usecase"

	"github.com/gorilla/mux"
)

func main() {
	// Загрузка конфигурации
	config, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	// Инициализация подключения к базе данных
	db, err := mssql.InitSQLDB(config.DBConnectionString)
	if err != nil {
		log.Fatal("Error initializing SQL database:", err)
	}
	defer db.Close()

	// Инициализация репозитория валюты
	repo := repository.NewSQLCurrencyRepository(db)
	externalClient := usecase.NewNationalBankClient()

	// Инициализация сервиса валюты
	service := usecase.NewCurrencyService(repo, externalClient)

	// Инициализация обработчика HTTP
	handler := web.NewCurrencyHandler(service)

	// Настройка маршрутов
	r := mux.NewRouter()
	r.HandleFunc("/currency/save/{date}", handler.SaveCurrency).Methods("GET")
	r.HandleFunc("/currency/save/{date}", handler.GetCurrency).Methods("GET")

	// Запуск сервера
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
