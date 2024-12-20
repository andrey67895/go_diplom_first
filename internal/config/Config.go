package config

import (
	"flag"
	"fmt"
	"os"
)

var DatabaseDsn string
var RunAddress string
var AccrualSystemAddress string

func InitServerConfig() {
	flag.StringVar(&DatabaseDsn, "d", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", `localhost`, 5434, `docker`, `docker`, `postgres`), "Aдрес подключения к базе данных")
	flag.StringVar(&RunAddress, "a", ":8787", "Aдрес и порт запуска сервиса")
	flag.StringVar(&AccrualSystemAddress, "r", ":8080", "адрес системы расчёта начислений")

	flag.Parse()
	if envDatabaseDsn := os.Getenv("DATABASE_URI"); envDatabaseDsn != "" {
		DatabaseDsn = envDatabaseDsn
	}
	if envRunAddress := os.Getenv("RUN_ADDRESS"); envRunAddress != "" {
		RunAddress = envRunAddress
	}
	if envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystemAddress != "" {
		AccrualSystemAddress = envAccrualSystemAddress
	}
}
