package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/mux"
)

type Currency struct {
	Title string    `json:"title"`
	Code  string    `json:"code"`
	Value float64   `json:"value"`
	ADate time.Time `json:"a_date"`
}

type Config struct {
	Port               string `json:"port"`
	DBConnectionString string `json:"db_connection_string"`
}

type Rates struct {
	XMLName xml.Name `xml:"rates"`
	Date    string   `xml:"date"`
	Items   []Item   `xml:"item"`
}

type Item struct {
	FullName    string  `xml:"fullname"`
	Title       string  `xml:"title"`
	Description float64 `xml:"description"`
}

func getDataFromNationalBank(date string) (currencyData []Currency, err error) {
	// Формирование URL с датой
	url := "https://nationalbank.kz/rss/get_rates.cfm?fdate=" + date

	// Запрос данных из API национального банка
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching data from national bank API:", err)
		return
	}
	defer resp.Body.Close()

	// Чтение тела ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return
	}

	// Декодирование XML
	var rates Rates
	err = xml.Unmarshal(body, &rates)
	if err != nil {
		log.Println("Error decoding XML:", err)
		return
	}
	return
}

func saveCurrency(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	date := params["date"]

	// Получение данных из API национального банка
	currencyData, err := getDataFromNationalBank(date)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Асинхронное сохранение данных в базу данных
	go func() {
		db, err := sql.Open("mssql", config.DBConnectionString)
		if err != nil {
			log.Println("Error opening database connection:", err)
			return
		}
		defer db.Close()

		for _, currency := range currencyData {
			_, err := db.Exec("INSERT INTO R_CURRENCY (TITLE, CODE, VALUE, A_DATE) VALUES (?, ?, ?, ?)",
				currency.Title, currency.Code, currency.Value, currency.ADate)
			if err != nil {
				log.Println("Error saving currency data to database:", err)
			}
		}

		log.Println("Currency data saved to database asynchronously")
	}()

	// Ответ пользователю
	response := map[string]bool{"success": true}
	json.NewEncoder(w).Encode(response)
}

func getCurrency(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	date := params["date"]
	code := params["code"]

	// Получение данных из базы данных
	db, err := sql.Open("mssql", config.DBConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var currencyData []Currency
	query := "SELECT TITLE, CODE, VALUE, A_DATE FROM R_CURRENCY WHERE A_DATE=?"
	if code != "" {
		query += " AND CODE=?"
		rows, err := db.Query(query, date, code)
		if err != nil {
			log.Println("Error querying currency data from database:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var currency Currency
			err := rows.Scan(&currency.Title, &currency.Code, &currency.Value, &currency.ADate)
			if err != nil {
				log.Println("Error scanning currency data:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			currencyData = append(currencyData, currency)
		}
	} else {
		rows, err := db.Query(query, date)
		if err != nil {
			log.Println("Error querying currency data from database:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var currency Currency
			err := rows.Scan(&currency.Title, &currency.Code, &currency.Value, &currency.ADate)
			if err != nil {
				log.Println("Error scanning currency data:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			currencyData = append(currencyData, currency)
		}
	}

	// Ответ пользователю
	json.NewEncoder(w).Encode(currencyData)
}

var config Config

func main() {
	// Загрузка конфигурации из файла
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Error opening config file:", err)
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		log.Fatal("Error decoding config JSON:", err)
	}

	// Создание роутера
	router := mux.NewRouter()

	// Регистрация обработчиков
	router.HandleFunc("/currency/save/{date}", saveCurrency).Methods("GET")
	router.HandleFunc("/currency/{date}/{code}", getCurrency).Methods("GET")

	// Запуск веб-сервера
	http.Handle("/", router)
	port := ":" + config.Port
	log.Fatal(http.ListenAndServe(port, router))
}
