package usecase

import (
	"currency/domain"
	"currency/repository"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type currencyService struct {
	Repo           repository.CurrencyRepository
	ExternalClient repository.ExternalService
}

// CurrencyUsecase - Интерфейс сервиса валюты
type CurrencyUsecase interface {
	SaveCurrency(date string) error
	GetCurrency(date, code string) ([]domain.Currency, error)
}

// NewCurrencyUsecase - Инициализация сервиса валюты
func NewCurrencyUsecase(repo repository.CurrencyRepository, externalClient repository.ExternalService) CurrencyUsecase {
	return &currencyService{
		Repo:           repo,
		ExternalClient: externalClient,
	}
}

// SaveCurrency - Сохранение валюты в базу данных
func (service *currencyService) SaveCurrency(date string) error {
	currencyData, err := service.ExternalClient.GetCurrencyData(date)
	if err != nil {
		return err
	}

	go func() {
		err := service.Repo.Save(currencyData)
		if err != nil {
			log.Println("Error saving currency data to database:", err)
		} else {
			log.Println("Currency data saved to database asynchronously")
		}
	}()

	return nil
}

// GetCurrency - Получение валюты из базы данных
func (service *currencyService) GetCurrency(date, code string) ([]domain.Currency, error) {
	return service.Repo.GetCurrency(date, code)
}

// nationalBankClient - Клиент внешних данных
type nationalBankClient struct{}

// NewNationalBankClient - Инициализация клиента внешних данных
func NewNationalBankClient() *nationalBankClient {
	return &nationalBankClient{}
}

// GetCurrencyData - Получение валюты из API National Bank KZ
func (nbc *nationalBankClient) GetCurrencyData(date string) ([]domain.Currency, error) {
	// Construct the URL with the provided date
	url := fmt.Sprintf("https://nationalbank.kz/rss/get_rates.cfm?fdate=%s", date)

	// Make a GET request to the National Bank API
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Decode XML response
	var rates struct {
		Items []struct {
			Title string `xml:"fullname"`
			Code  string `xml:"title"`
			Value string `xml:"description"`
			Date  string `xml:"date"`
		} `xml:"rates>item"`
	}
	if err := xml.Unmarshal(body, &rates); err != nil {
		return nil, err
	}

	// Convert XML data to Currency struct
	var currencies []domain.Currency
	for _, item := range rates.Items {
		value, err := strconv.ParseFloat(item.Value, 64)
		if err != nil {
			return nil, err
		}
		adate, err := time.Parse("02.01.2006", item.Date)
		if err != nil {
			return nil, err
		}
		currencies = append(currencies, domain.Currency{
			Title: item.Title,
			Code:  item.Code,
			Value: value,
			ADate: adate,
		})
	}

	return currencies, nil
}
