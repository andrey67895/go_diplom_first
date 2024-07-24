package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/andrey67895/go_diplom_first/internal/config"
	"github.com/andrey67895/go_diplom_first/internal/database"
	"github.com/andrey67895/go_diplom_first/internal/helpers"
	"github.com/andrey67895/go_diplom_first/internal/job"
	"github.com/andrey67895/go_diplom_first/internal/transport"
)

func InitServer() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	var wg sync.WaitGroup
	wg.Add(10)

	config.InitServerConfig()
	if err := database.InitDB(ctx); err != nil {
		helpers.TLog.Fatal(err)
	}
	job.OrdersStatusJob(ctx, &wg)
	go func() {
		defer wg.Done()
		helpers.TLog.Fatal(http.ListenAndServe(config.RunAddress, transport.GetRoutersGophermart()))
	}()
	select {
	case <-ctx.Done():
		helpers.TLog.Info("Server shutdown: ", ctx.Err())
	}
}
