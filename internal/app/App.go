package app

import (
	"context"
	"net/http"
	"os/exec"

	"github.com/andrey67895/go_diplom_first/internal/config"
	"github.com/andrey67895/go_diplom_first/internal/database"
	"github.com/andrey67895/go_diplom_first/internal/helpers"
	"github.com/andrey67895/go_diplom_first/internal/transport"
)

func InitServer() {
	helpers.InitLog()
	config.InitServerConfig()
	database.InitDB(context.Background())
	go func() {
		cmd := exec.Command("cmd", "/C", "start", config.AccrualSystemAddress)
		err := cmd.Start()
		if err != nil {
			helpers.TLog.Fatal(err.Error())
			return
		}
	}()
	helpers.TLog.Fatal(http.ListenAndServe(config.RunAddress, transport.GetRouters()))

}
