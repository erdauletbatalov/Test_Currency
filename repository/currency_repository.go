// data_access.go

package repository

import (
	"currency/domain"
	"database/sql"
)

type CurrencyRepository interface {
	Save(currency []domain.Currency) error
	GetCurrency(date, code string) ([]domain.Currency, error)
}

type ExternalService interface {
	GetCurrencyData(date string) ([]domain.Currency, error)
}

type SQLCurrencyRepository struct {
	DB *sql.DB
}

func NewSQLCurrencyRepository(db *sql.DB) *SQLCurrencyRepository {
	return &SQLCurrencyRepository{DB: db}
}

func (repo *SQLCurrencyRepository) Save(currencyData []domain.Currency) error {

	for _, currency := range currencyData {
		if _, err := repo.DB.Exec("INSERT INTO R_CURRENCY (TITLE, CODE, VALUE, A_DATE) VALUES (?, ?, ?, ?)",
			currency.Title, currency.Code, currency.Value, currency.ADate); err != nil {
			return err
		}
	}
	return nil
}

func (repo *SQLCurrencyRepository) GetCurrency(date, code string) ([]domain.Currency, error) {
	var currencyData []domain.Currency
	query := "SELECT TITLE, CODE, VALUE, A_DATE FROM R_CURRENCY WHERE A_DATE=?"
	if code != "" {
		query += " AND CODE=?"
	}
	rows, err := repo.DB.Query(query, date, code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var currency domain.Currency
		err := rows.Scan(&currency.Title, &currency.Code, &currency.Value, &currency.ADate)
		if err != nil {
			return nil, err
		}
		currencyData = append(currencyData, currency)
	}
	return currencyData, nil
}
