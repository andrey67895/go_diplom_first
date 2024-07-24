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
	"golang.org/x/sync/errgroup"
)

func InitServer() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	var wg sync.WaitGroup

	config.InitServerConfig()
	if err := database.InitDB(ctx); err != nil {
		helpers.TLog.Fatal(err)
	}
	job.OrdersStatusJob(ctx, &wg)
	httpServer := &http.Server{
		Addr:    config.RunAddress,
		Handler: transport.GetRoutersGophermart(),
	}
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return httpServer.Shutdown(context.Background())
	})
	if err := g.Wait(); err != nil {
		helpers.TLog.Info("Server shutdown: ", err.Error(), " :: ", ctx.Err())
	}
}
