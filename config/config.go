package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host       string `json:"host"`
	Port       string `json:"port"`
	User       string `json:"username"`
	Password   string `json:"password"`
	Dbname     string `json:"db_name"`
	Ssl        string `json:"ssl"`
	TimeZone   string `json:"time_zone"`
	DbHost     string `json:"db_host"`
	DbPort     string `json:"db_port"`
	TTLMinutes int    `json:"ttl_minutes"`
}

func LoadConfig(path string) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error while open config file: %v", err)
	}
	defer file.Close()

	var cfg Config

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("error to decode config: %v", err)
	}
	return &cfg, nil
}
