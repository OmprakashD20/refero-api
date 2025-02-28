package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	AppEnv        string
	Port           string
	DB             DBConfig
}

type DBConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	SSLMode    string
}

var Envs = initConfig()

func initConfig() EnvConfig {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	return EnvConfig{
		AppEnv: *getEnv("APP_ENV"),
		Port:   *getEnv("PORT"),
		DB: DBConfig{
			DBHost:     *getEnv("DB_HOST"),
			DBPort:     *getEnv("DB_PORT"),
			DBUser:     *getEnv("DB_USER"),
			DBPassword: *getEnv("DB_PASS"),
			DBName:     *getEnv("DB_NAME"),
			SSLMode:    *getEnv("DB_SSLMODE"),
		}, 
	}
}

func getEnv(key string) *string {
	if value, ok := os.LookupEnv(key); ok {
		return &value
	}

	log.Fatalf("Environment variable %s is not set", key)
	return nil
}
