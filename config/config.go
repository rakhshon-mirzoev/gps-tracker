package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type GeneralConfig struct {
	Admin  Admin       `json:"admin"`
	Minio  MinioConfig `json:"minio"`
	Mobile Mobile      `json:"mobile"`
	DB     DB          `json:"db"`
}

type Admin struct {
	Host       string `json:"host"`
	Port       string `json:"port"`
	TTLMinutes int    `json:"ttl_minutes"`
}

type Mobile struct {
	Host       string `json:"host"`
	Port       string `json:"port"`
	TTLMinutes int    `json:"ttl_minutes"`
}

type MinioConfig struct {
	Endpoint   string `json:"endpoint"`
	AccessKey  string `json:"access_key"`
	SecretKey  string `json:"secret_key"`
	Bucket     string `json:"bucket"`
	UseSSL     bool   `json:"use_ssl"`
	TTLMinutes int    `json:"ttl_minutes"`
}

type DB struct {
	Username string `json:"username"`
	Password string `json:"password"`
	DbName   string `json:"db_name"`
	Ssl      string `json:"ssl"`
	TimeZone string `json:"time_zone"`
	DbHost   string `json:"db_host"`
	DbPort   string `json:"db_port"`
}

func LoadConfig(path string) (*GeneralConfig, error) {
	_ = godotenv.Load()

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error while open config file: %v", err)
	}
	defer file.Close()

	var cfg GeneralConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("error to decode config: %v", err)
	}
	return &cfg, nil
}
