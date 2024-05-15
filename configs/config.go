package configs

import (
	"os"
	"strconv"
)

type Env struct {
	DBHost         string
	DBPort         int
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLModel     string
	DBMaxOpenConns int
}

var Config *Env

func LoadEnv() *Env {
	if Config != nil {
		return Config
	}

	config := &Env{}

	config.DBHost = getEnv("DB_HOST")
	config.DBPort, _ = strconv.Atoi(getEnv("DB_PORT"))
	config.DBUser = getEnv("DB_USER")
	config.DBPassword = getEnv("DB_PASSWORD")
	config.DBName = getEnv("DB_NAME")
	config.DBSSLModel = getEnv("DB_SSLMODE")
	config.DBMaxOpenConns, _ = strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS"))

	return config
}

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)

	if !exists {
		panic("Environment variable " + key + " is not set")
	}

	return value
}
