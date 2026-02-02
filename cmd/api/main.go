package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/khanalsaroj/typegen-server/internal/pkg/crypto"

	"github.com/khanalsaroj/typegen-server/internal/config"
	"github.com/khanalsaroj/typegen-server/internal/pkg/logger"
	"github.com/khanalsaroj/typegen-server/internal/server"

	"go.uber.org/zap"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	log := logger.New(cfg.App.Environment)
	defer func(log *zap.Logger) {
		err := log.Sync()
		if err != nil {
		}
	}(log)

	log.Info("Starting application",
		zap.String("environment", cfg.App.Environment),
		zap.String("version", cfg.App.Version),
	)

	db, err := connectDatabase(cfg, log)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}

	if err := runMigrations(db, log); err != nil {
		log.Fatal("Failed to run migrations", zap.Error(err))
	}

	cryptoSvc, err := crypto.New(cfg.Security.DbEncryptionKey)
	if err != nil {
		log.Fatal("Failed to initialize crypto service", zap.Error(err))
	}

	srv := server.New(cfg, db, cryptoSvc, log)

	go func() {
		addr := fmt.Sprintf(":%d", cfg.Server.Port)
		log.Info("Starting HTTP server", zap.String("address", addr))

		if err := srv.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exited")
}
