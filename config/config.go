package config

import (
	"encoding/json"
	"os"
)

// Config - Конфигурация приложения
type Config struct {
	DBConnectionString string `json:"db_connection_string"`
	ServerPort         string `json:"server_port"`
}

// LoadConfig - Загрузка конфигурации приложения
func LoadConfig(filename string) (Config, error) {
	var config Config
	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
