package web

import (
	"encoding/json"
	"log"
	"net/http"

	"currency/usecase"

	"github.com/gorilla/mux"
)

type CurrencyHandler struct {
	Service usecase.CurrencyUsecase
}

func NewCurrencyHandler(service usecase.CurrencyUsecase) *CurrencyHandler {
	return &CurrencyHandler{Service: service}
}

func (handler *CurrencyHandler) SaveCurrency(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	date := params["date"]
	// Call the SaveCurrency method of CurrencyService
	err := handler.Service.SaveCurrency(date)
	if err != nil {
		log.Println("Error saving currency data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Respond to the client
	response := map[string]bool{"success": true}
	json.NewEncoder(w).Encode(response)
}

func (handler *CurrencyHandler) GetCurrency(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	date := params["date"]
	code := params["code"]

	// Call the service layer to get currency data
	currencyData, err := handler.Service.GetCurrency(date, code)
	if err != nil {
		log.Println("Error retrieving currency data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Respond to the client with the retrieved currency data
	json.NewEncoder(w).Encode(currencyData)
}
