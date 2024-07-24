package app

import (
	"context"
	"net/http"

	"github.com/andrey67895/go_diplom_first/internal/config"
	"github.com/andrey67895/go_diplom_first/internal/database"
	"github.com/andrey67895/go_diplom_first/internal/helpers"
	"github.com/andrey67895/go_diplom_first/internal/job"
	"github.com/andrey67895/go_diplom_first/internal/transport"
)

func InitServer() {
	config.InitServerConfig()
	if err := database.InitDB(context.Background()); err != nil {
		helpers.TLog.Fatal(err)
	}
	go job.OrdersStatusJob()
	helpers.TLog.Fatal(http.ListenAndServe(config.RunAddress, transport.GetRoutersGophermart()))
}
