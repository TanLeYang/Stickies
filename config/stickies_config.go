package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	BotName    string
	TgBotToken string
	DbConf     DbConfig
}

type DbConfig struct {
	DbHost     string
	DbPort     string
	DbUser     string
	DbName     string
	DbPassword string
}

func LoadAppConfig() AppConfig {
	if os.Getenv("ENV") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	return AppConfig{
		BotName:    getEnvVar("BOT_NAME"),
		TgBotToken: getEnvVar("TG_BOT_TOKEN"),

		DbConf: DbConfig{
			DbHost:     getEnvVar("DATABASE_HOST"),
			DbPort:     getEnvVar("DATABASE_PORT"),
			DbUser:     getEnvVar("DATABASE_USER"),
			DbName:     getEnvVar("DATABASE_DBNAME"),
			DbPassword: getEnvVar("DATABASE_PASSWORD"),
		},
	}
}

func getEnvVar(varName string) string {
	ev := os.Getenv(varName)
	if ev == "" {
		log.Panicf("environment variable %s not set", varName)
	}

	return ev
}
