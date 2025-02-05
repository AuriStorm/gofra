package app

import (
	"context"
	"errors"
	"gofra/internal/config"
	"gofra/internal/storage"
	"gofra/internal/transport/rest"
	"gofra/internal/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	Server        *http.Server
	AppConfig     config.AppConfig
	StorageConfig config.StorageConfig
}

func New() App {
	return App{}
}

func (a *App) Configure(cliArgs utils.CliArgs) {
	a.AppConfig.MustLoad(cliArgs)
	a.StorageConfig.MustLoad(cliArgs)

	inmemQ := storage.NewInmemoryQueue(a.StorageConfig)
	routing := rest.NewRouting(inmemQ, a.AppConfig)

	mux := http.NewServeMux()
	rest.RegisterAppRoutes(mux, routing)
	a.Server = &http.Server{Addr: a.AppConfig.Addr, Handler: mux}
}

func (a *App) RunServer() {

	go func() {
		if err := a.Server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			// thought use panic here at first, but think os.Exit will be better for CI/CD teardown
			log.Fatalf("something went wrong on server start up: %s", err)
			os.Exit(2)
		}
		log.Println("end up connection serving")
	}()

	// не ну прямо вообще без логов кажется совсем кощунственно...
	log.Println("server started")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	ctx, cancel := context.WithTimeout(context.Background(), a.AppConfig.ShutdownTimeoutSec*time.Second)
	defer cancel()

	if err := a.Server.Shutdown(ctx); err != nil {
		// here is need some workaround to ensure everything has been closed properly tbh
		log.Fatalf("something went wrong on server shutdown: %v", err)
	}

	// TODO have to think about panic reload too

	log.Print("server finished\n\n")
}
