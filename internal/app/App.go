package app

import (
	"context"
	"net/http"

	"github.com/andrey67895/go_diplom_first/internal/config"
	"github.com/andrey67895/go_diplom_first/internal/database"
	"github.com/andrey67895/go_diplom_first/internal/helpers"
	"github.com/andrey67895/go_diplom_first/internal/transport"
)

func InitServer() {
	helpers.InitLog()
	config.InitServerConfig()
	database.InitDB(context.Background())
	go helpers.TLog.Fatal(http.ListenAndServe(config.RunAddress, transport.GetRouters()))
	go helpers.TLog.Fatal(http.ListenAndServe(config.AccrualSystemAddress, transport.GetRouters()))
	select {}
}
