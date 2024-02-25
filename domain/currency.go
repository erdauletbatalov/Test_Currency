package domain

import "time"

type Currency struct {
	Title string    `json:"title"`
	Code  string    `json:"code"`
	Value float64   `json:"value"`
	ADate time.Time `json:"a_date"`
}
