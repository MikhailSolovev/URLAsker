package main

import (
	"context"
	"fmt"
	"github.com/MikhailSolovev/URLAsker/internal/config"
	"github.com/MikhailSolovev/URLAsker/internal/models"
	"github.com/MikhailSolovev/URLAsker/internal/services/asker"
	"github.com/MikhailSolovev/URLAsker/internal/storage/pg"
	"github.com/MikhailSolovev/URLAsker/internal/transport/rest"
	"github.com/MikhailSolovev/URLAsker/pkg/logger"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	cfg := config.New()
	if cfg == nil {
		log.Fatal("Can't parse config!")
	}

	lg := logger.New(cfg.ServiceName, cfg.IsPretty, cfg.LogLevel)
	lg.LogInfo(fmt.Sprintf("Logger init for %v with %v level", cfg.ServiceName, cfg.LogLevel))

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GracefulTimeout)
	defer cancel()

	storage, pgPool, err := pg.NewPool(ctx, cfg.PostgresUser, cfg.PostgresPass, cfg.PostgresHost, cfg.PostgresPort,
		cfg.PostgresDB, cfg.PostgresConnTimeout)
	if err != nil {
		lg.LogFatal(fmt.Sprintf("failed to connect to pg due to error: %v", err))
	}
	if err = pgPool.Ping(ctx); err != nil {
		lg.LogFatal(fmt.Sprintf("failed to ping pg due to error: %v", err))
	}

	askerSvc := asker.New(storage, &models.Info{Interval: cfg.AskerInterval, Ticker: time.NewTicker(cfg.AskerInterval),
		URLs: map[string]models.Empty{}})

	lg.LogInfo(fmt.Sprintf("Asker started with interval %v", cfg.AskerInterval.String()))
	ctxAsker, cancelAsker := context.WithCancel(context.Background())
	defer cancelAsker()
	go func() {
		if err = askerSvc.Run(ctxAsker); err != nil {
			lg.LogFatal(fmt.Sprintf("asker crashed due to error: %v", err))
		}
	}()

	router := mux.NewRouter()
	handler := rest.New(router, askerSvc)
	handler.Register()

	listener, err := net.Listen("tcp", cfg.SrvAddr)
	if err != nil {
		lg.LogFatal(err.Error())
	}

	srv := &http.Server{
		Handler:      router,
		WriteTimeout: cfg.SrvWriteTimeout,
		ReadTimeout:  cfg.SrvReadTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	lg.LogInfo(fmt.Sprintf("Server started on %s with: Write timeout: %s Read timeout: %s Idle timeout: %s",
		cfg.SrvAddr, cfg.SrvWriteTimeout, cfg.SrvReadTimeout, cfg.IdleTimeout))
	go func() {
		if err = srv.Serve(listener); err != nil {
			lg.LogInfo(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	if err = srv.Shutdown(ctx); err != nil {
		lg.LogFatal(fmt.Sprintf("shutdown problem due to error: %v", err))
	}
	lg.LogInfo("shutting down")
	os.Exit(0)
}
