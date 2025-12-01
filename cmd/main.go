package main

import (
	"context"
	"errors"
	"fmt"
	"gw-currency-wallet/config"
	"gw-currency-wallet/internal/handler"
	"gw-currency-wallet/internal/repository"
	"gw-currency-wallet/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
}

func main() {
	cfg := config.LoadConfig()

	ctx := context.Background()
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		zap.L().Fatal("error to open connect to database")
	}

	r := repository.NewRepository(pool)
	s := service.NewService(r, &cfg.Auth)
	h := handler.NewHandler(s)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: h.Router(),
	}

	go func() {
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.L().Fatal("could not listen")
		}
	}()

	zap.L().Info("server is running")

	<-stop
	zap.L().Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		zap.L().Fatal("server forced to shutdown")
	}

	if pool != nil {
		pool.Close()
	}

	zap.L().Info("server gracefully stopped")
}
